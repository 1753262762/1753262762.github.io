package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func testSite(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{"index.html", "blog/index.html", "blog/hello/index.html", "prototype/index.html"} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("ok"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func testStore(t *testing.T) *store {
	t.Helper()
	s, err := openStore(filepath.Join(t.TempDir(), "stats.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.db.Close() })
	return s
}

func TestNormalizeAndRoute(t *testing.T) {
	root := testSite(t)
	cases := map[string]string{"/": "/", "/index.html": "/", "/blog/hello/?x=1": "/blog/hello/"}
	for input, want := range cases {
		got, err := normalizePagePath(input)
		if err != nil || got != want {
			t.Fatalf("normalize %q = %q, %v", input, got, err)
		}
	}
	if !isFormalRoute(root, "/blog/hello/") {
		t.Fatal("current article should be formal")
	}
	if isFormalRoute(root, "/prototype/") {
		t.Fatal("prototype must be excluded")
	}
	if isFormalRoute(root, "/missing/") {
		t.Fatal("missing route must be excluded")
	}
}

func TestConcurrentIncrement(t *testing.T) {
	s := testStore(t)
	const workers = 24
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := s.increment(context.Background(), "/blog/hello/"); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	counts, err := s.increment(context.Background(), "/")
	if err != nil {
		t.Fatal(err)
	}
	if counts.TotalViews != workers+1 {
		t.Fatalf("total=%d", counts.TotalViews)
	}
	var article int64
	if err := s.db.QueryRow(`SELECT views FROM page_views WHERE path='/blog/hello/'`).Scan(&article); err != nil {
		t.Fatal(err)
	}
	if article != workers {
		t.Fatalf("article=%d", article)
	}
}

func TestStorePersistsAcrossReopen(t *testing.T) {
	database := filepath.Join(t.TempDir(), "stats.db")
	s, err := openStore(database)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.increment(context.Background(), "/"); err != nil {
		t.Fatal(err)
	}
	if err := s.db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = openStore(database)
	if err != nil {
		t.Fatal(err)
	}
	defer s.db.Close()
	var total int64
	if err := s.db.QueryRow(`SELECT views FROM page_views WHERE path=?`, totalKey).Scan(&total); err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Fatalf("persisted total=%d", total)
	}
}

func TestViewHandler(t *testing.T) {
	root := testSite(t)
	s := testStore(t)
	a := &api{store: s, config: config{SiteRoot: root, AllowedOrigin: "https://example.test"}, limiter: newLimiter()}
	req := httptest.NewRequest(http.MethodPost, "/api/stats/view", strings.NewReader(`{"path":"/blog/hello/"}`))
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://example.test")
	req.Header.Set("X-Real-IP", "203.0.113.2")
	res := httptest.NewRecorder()
	a.routes().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	var counts viewCounts
	if err := json.Unmarshal(res.Body.Bytes(), &counts); err != nil {
		t.Fatal(err)
	}
	if counts.TotalViews != 1 || counts.ArticleViews == nil || *counts.ArticleViews != 1 {
		t.Fatalf("counts=%+v", counts)
	}
}

func TestViewHandlerRejectsOriginAndPrototype(t *testing.T) {
	root := testSite(t)
	s := testStore(t)
	a := &api{store: s, config: config{SiteRoot: root, AllowedOrigin: "https://example.test"}, limiter: newLimiter()}

	request := func(origin, body string) int {
		req := httptest.NewRequest(http.MethodPost, "/api/stats/view", strings.NewReader(body))
		req.RemoteAddr = "127.0.0.1:1234"
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", origin)
		res := httptest.NewRecorder()
		a.routes().ServeHTTP(res, req)
		return res.Code
	}
	if status := request("https://evil.test", `{"path":"/"}`); status != http.StatusForbidden {
		t.Fatalf("bad origin status=%d", status)
	}
	if status := request("https://example.test", `{"path":"/prototype/"}`); status != http.StatusBadRequest {
		t.Fatalf("prototype status=%d", status)
	}
}

func TestLimiter(t *testing.T) {
	l := newLimiter()
	now := time.Unix(0, 0)
	l.now = func() time.Time { return now }
	for i := 0; i < 20; i++ {
		if !l.allow("ip") {
			t.Fatalf("request %d rejected", i)
		}
	}
	if l.allow("ip") {
		t.Fatal("burst should be exhausted")
	}
	now = now.Add(time.Second)
	if !l.allow("ip") {
		t.Fatal("one token should refill per second")
	}
}

func TestImportAndRefuseSecondImport(t *testing.T) {
	root := testSite(t)
	logPath := filepath.Join(t.TempDir(), "access.log")
	lines := strings.Join([]string{
		`203.0.113.1 - - [28/Aug/2026:12:00:00 +0800] "GET /blog/hello/ HTTP/1.1" 200 10 "-" "Mozilla/5.0"`,
		`203.0.113.1 - - [28/Aug/2026:12:00:01 +0800] "GET /prototype/ HTTP/1.1" 200 10 "-" "Mozilla/5.0"`,
		`203.0.113.1 - - [28/Aug/2026:12:00:02 +0800] "GET / HTTP/1.1" 200 10 "-" "Googlebot"`,
		`bad line`,
	}, "\n")
	if err := os.WriteFile(logPath, []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := scanLogs([]string{logPath}, root)
	if err != nil {
		t.Fatal(err)
	}
	if report.Imported != 1 || report.Bots != 1 || report.UnknownRoutes != 1 || report.ParseErrors != 1 {
		t.Fatalf("report=%+v", report)
	}
	s := testStore(t)
	if err := s.applyImport(context.Background(), report.ViewsByPath); err != nil {
		t.Fatal(err)
	}
	if err := s.applyImport(context.Background(), report.ViewsByPath); err == nil {
		t.Fatal("second import should fail")
	}
}

func TestCompressedLogImport(t *testing.T) {
	root := testSite(t)
	logPath := filepath.Join(t.TempDir(), "access.log.1.gz")
	file, err := os.Create(logPath)
	if err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	_, err = writer.Write([]byte(`203.0.113.1 - - [28/Aug/2026:12:00:00 +0800] "HEAD /blog/hello/ HTTP/1.1" 304 0 "-" "Mozilla/5.0"`))
	if closeErr := writer.Close(); err == nil {
		err = closeErr
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		t.Fatal(err)
	}
	report, err := scanLogs([]string{logPath}, root)
	if err != nil {
		t.Fatal(err)
	}
	if report.Imported != 1 || report.ViewsByPath["/blog/hello/"] != 1 {
		t.Fatalf("report=%+v", report)
	}
}
