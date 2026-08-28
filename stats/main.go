package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

const totalKey = "__total__"

type config struct {
	Address       string
	DatabasePath  string
	SiteRoot      string
	AllowedOrigin string
}

type store struct{ db *sql.DB }

func openStore(path string) (*store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, err
	}
	dsn := "file:" + filepath.ToSlash(path) + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &store{db: db}
	if err := s.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *store) migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS page_views (
			path TEXT PRIMARY KEY,
			views INTEGER NOT NULL CHECK (views >= 0)
		);
		CREATE TABLE IF NOT EXISTS metadata (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`)
	return err
}

type viewCounts struct {
	TotalViews   int64  `json:"totalViews"`
	ArticleViews *int64 `json:"articleViews"`
}

func (s *store) increment(ctx context.Context, path string) (viewCounts, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return viewCounts{}, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `INSERT INTO metadata(key,value) VALUES('live_started','1') ON CONFLICT(key) DO UPDATE SET value='1'`); err != nil {
		return viewCounts{}, err
	}
	for _, key := range []string{totalKey, path} {
		if _, err = tx.ExecContext(ctx, `INSERT INTO page_views(path,views) VALUES(?,1) ON CONFLICT(path) DO UPDATE SET views=views+1`, key); err != nil {
			return viewCounts{}, err
		}
	}
	var counts viewCounts
	if err = tx.QueryRowContext(ctx, `SELECT views FROM page_views WHERE path=?`, totalKey).Scan(&counts.TotalViews); err != nil {
		return viewCounts{}, err
	}
	if isArticlePath(path) {
		var article int64
		if err = tx.QueryRowContext(ctx, `SELECT views FROM page_views WHERE path=?`, path).Scan(&article); err != nil {
			return viewCounts{}, err
		}
		counts.ArticleViews = &article
	}
	if err = tx.Commit(); err != nil {
		return viewCounts{}, err
	}
	return counts, nil
}

func (s *store) applyImport(ctx context.Context, counts map[string]int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var marker string
	err = tx.QueryRowContext(ctx, `SELECT value FROM metadata WHERE key IN ('live_started','import_applied') LIMIT 1`).Scan(&marker)
	if err == nil {
		return errors.New("database already contains live or imported data; refusing to import twice")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	var total int64
	paths := make([]string, 0, len(counts))
	for path, n := range counts {
		total += n
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if _, err = tx.ExecContext(ctx, `INSERT INTO page_views(path,views) VALUES(?,?)`, path, counts[path]); err != nil {
			return err
		}
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO page_views(path,views) VALUES(?,?)`, totalKey, total); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO metadata(key,value) VALUES('import_applied',?)`, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return err
	}
	return tx.Commit()
}

func normalizePagePath(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" || !strings.HasPrefix(u.Path, "/") {
		return "", errors.New("invalid path")
	}
	p := u.EscapedPath()
	decoded, err := url.PathUnescape(p)
	segments := "/" + strings.Trim(decoded, "/") + "/"
	if err != nil || strings.Contains(decoded, "\\") || strings.Contains(decoded, "\x00") || strings.Contains(segments, "/../") || strings.Contains(segments, "/./") {
		return "", errors.New("invalid path")
	}
	clean := filepath.ToSlash(filepath.Clean(decoded))
	if clean == "." {
		clean = "/"
	}
	if strings.HasSuffix(decoded, "/") && clean != "/" {
		clean += "/"
	}
	if strings.HasSuffix(clean, "/index.html") {
		clean = strings.TrimSuffix(clean, "index.html")
	}
	if clean == "/index.html" {
		clean = "/"
	}
	return clean, nil
}

func isArticlePath(path string) bool {
	if !strings.HasPrefix(path, "/blog/") || path == "/blog/" || !strings.HasSuffix(path, "/") {
		return false
	}
	return !strings.Contains(strings.TrimSuffix(strings.TrimPrefix(path, "/blog/"), "/"), "/")
}

func isFormalRoute(siteRoot, path string) bool {
	if path == "" || strings.HasPrefix(path, "/prototype/") || strings.HasPrefix(path, "/api/") || !strings.HasSuffix(path, "/") {
		return false
	}
	rel := strings.TrimPrefix(path, "/")
	file := filepath.Join(siteRoot, filepath.FromSlash(rel), "index.html")
	root, err1 := filepath.Abs(siteRoot)
	resolved, err2 := filepath.Abs(file)
	if err1 != nil || err2 != nil || (resolved != root && !strings.HasPrefix(resolved, root+string(os.PathSeparator))) {
		return false
	}
	info, err := os.Stat(resolved)
	return err == nil && !info.IsDir()
}

type bucket struct {
	tokens float64
	last   time.Time
}

type limiter struct {
	mu      sync.Mutex
	buckets map[string]bucket
	now     func() time.Time
}

func newLimiter() *limiter {
	return &limiter{buckets: make(map[string]bucket), now: time.Now}
}

func (l *limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	b := l.buckets[key]
	if b.last.IsZero() {
		b = bucket{tokens: 20, last: now}
	}
	b.tokens += now.Sub(b.last).Seconds()
	if b.tokens > 20 {
		b.tokens = 20
	}
	b.last = now
	if b.tokens < 1 {
		l.buckets[key] = b
		return false
	}
	b.tokens--
	l.buckets[key] = b
	return true
}

type api struct {
	store   *store
	config  config
	limiter *limiter
}

func (a *api) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/stats/health", a.health)
	mux.HandleFunc("POST /api/stats/view", a.view)
	return mux
}

func noStore(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func writeError(w http.ResponseWriter, status int, message string) {
	noStore(w)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (a *api) localProxyRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return err == nil && net.ParseIP(host) != nil && net.ParseIP(host).IsLoopback()
}

func (a *api) health(w http.ResponseWriter, r *http.Request) {
	if !a.localProxyRequest(r) {
		writeError(w, http.StatusForbidden, "local proxy required")
		return
	}
	if err := a.store.db.PingContext(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	noStore(w)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (a *api) view(w http.ResponseWriter, r *http.Request) {
	if !a.localProxyRequest(r) {
		writeError(w, http.StatusForbidden, "local proxy required")
		return
	}
	if origin := r.Header.Get("Origin"); origin != "" && origin != a.config.AllowedOrigin {
		writeError(w, http.StatusForbidden, "origin not allowed")
		return
	}
	if !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		writeError(w, http.StatusUnsupportedMediaType, "application/json required")
		return
	}
	client := strings.TrimSpace(strings.Split(r.Header.Get("X-Real-IP"), ",")[0])
	if client == "" {
		client = "unknown"
	}
	if !a.limiter.allow(client) {
		writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var body struct {
		Path string `json:"path"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	path, err := normalizePagePath(body.Path)
	if err != nil || !isFormalRoute(a.config.SiteRoot, path) {
		writeError(w, http.StatusBadRequest, "unknown page")
		return
	}
	counts, err := a.store.increment(r.Context(), path)
	if err != nil {
		log.Printf("increment %q: %v", path, err)
		writeError(w, http.StatusInternalServerError, "could not record view")
		return
	}
	noStore(w)
	_ = json.NewEncoder(w).Encode(counts)
}

var combinedLog = regexp.MustCompile(`^(\S+) \S+ \S+ \[[^]]+\] "(GET|HEAD) ([^ ]+) HTTP/[^"]+" (\d{3}) \S+ "[^"]*" "([^"]*)"`)
var botAgent = regexp.MustCompile(`(?i)(bot|crawler|spider|slurp|bingpreview|facebookexternalhit|headless|lighthouse|uptimerobot|curl/|wget/)`)

type importReport struct {
	Files         int              `json:"files"`
	Lines         int64            `json:"lines"`
	Imported      int64            `json:"imported"`
	ParseErrors   int64            `json:"parseErrors"`
	Bots          int64            `json:"bots"`
	InvalidStatus int64            `json:"invalidStatus"`
	UnknownRoutes int64            `json:"unknownRoutes"`
	ViewsByPath   map[string]int64 `json:"viewsByPath"`
}

func openLog(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(strings.ToLower(path), ".gz") {
		return f, nil
	}
	gz, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return struct {
		io.Reader
		io.Closer
	}{gz, closerFunc(func() error { _ = gz.Close(); return f.Close() })}, nil
}

type closerFunc func() error

func (f closerFunc) Close() error { return f() }

func scanLogs(patterns []string, siteRoot string) (importReport, error) {
	report := importReport{ViewsByPath: make(map[string]int64)}
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return report, err
		}
		files = append(files, matches...)
	}
	sort.Strings(files)
	for _, path := range files {
		r, err := openLog(path)
		if err != nil {
			return report, fmt.Errorf("open %s: %w", path, err)
		}
		report.Files++
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)
		for scanner.Scan() {
			report.Lines++
			match := combinedLog.FindStringSubmatch(scanner.Text())
			if match == nil {
				report.ParseErrors++
				continue
			}
			status, _ := strconv.Atoi(match[4])
			if status != http.StatusOK && status != http.StatusNotModified {
				report.InvalidStatus++
				continue
			}
			if botAgent.MatchString(match[5]) {
				report.Bots++
				continue
			}
			page, err := normalizePagePath(match[3])
			if err != nil || !isFormalRoute(siteRoot, page) {
				report.UnknownRoutes++
				continue
			}
			report.ViewsByPath[page]++
			report.Imported++
		}
		scanErr := scanner.Err()
		_ = r.Close()
		if scanErr != nil {
			return report, fmt.Errorf("scan %s: %w", path, scanErr)
		}
	}
	if report.Files == 0 {
		return report, errors.New("no log files matched")
	}
	return report, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func runServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	cfg := config{}
	fs.StringVar(&cfg.Address, "address", env("STATS_ADDRESS", "127.0.0.1:8787"), "listen address")
	fs.StringVar(&cfg.DatabasePath, "database", env("STATS_DATABASE", "/var/lib/nabunana-stats/stats.db"), "SQLite database")
	fs.StringVar(&cfg.SiteRoot, "site-root", env("STATS_SITE_ROOT", "/var/www/blog"), "active static site root")
	fs.StringVar(&cfg.AllowedOrigin, "allowed-origin", env("STATS_ALLOWED_ORIGIN", "http://39.108.101.149"), "browser origin")
	if err := fs.Parse(args); err != nil {
		return err
	}
	s, err := openStore(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer s.db.Close()
	server := &http.Server{Addr: cfg.Address, Handler: (&api{store: s, config: cfg, limiter: newLimiter()}).routes(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("nabunana stats listening on %s", cfg.Address)
	return server.ListenAndServe()
}

func runImport(args []string) error {
	fs := flag.NewFlagSet("import-nginx", flag.ContinueOnError)
	database := fs.String("database", env("STATS_DATABASE", "/var/lib/nabunana-stats/stats.db"), "SQLite database")
	siteRoot := fs.String("site-root", env("STATS_SITE_ROOT", "/var/www/blog"), "active static site root")
	apply := fs.Bool("apply", false, "persist the previewed aggregate")
	if err := fs.Parse(args); err != nil {
		return err
	}
	patterns := fs.Args()
	if len(patterns) == 0 {
		patterns = []string{"/var/log/nginx/access.log*"}
	}
	report, err := scanLogs(patterns, *siteRoot)
	if err != nil {
		return err
	}
	output, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(output))
	if !*apply {
		return nil
	}
	s, err := openStore(*database)
	if err != nil {
		return err
	}
	defer s.db.Close()
	if err := s.applyImport(context.Background(), report.ViewsByPath); err != nil {
		return err
	}
	fmt.Printf("imported %d historical page views\n", report.Imported)
	return nil
}

func main() {
	command := "serve"
	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "serve" || args[0] == "import-nginx") {
		command, args = args[0], args[1:]
	}
	var err error
	if command == "import-nginx" {
		err = runImport(args)
	} else {
		err = runServe(args)
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
