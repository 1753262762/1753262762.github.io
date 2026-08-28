#!/usr/bin/env bash

set -Eeuo pipefail
umask 027

binary=${1:?Usage: install-stats-service.sh BINARY SITE_ORIGIN [NGINX_INCLUDE_DIRECTORY]}
site_origin=${2:?Usage: install-stats-service.sh BINARY SITE_ORIGIN [NGINX_INCLUDE_DIRECTORY]}
nginx_include_directory=${3:-}
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
repo_root=$(cd -- "$script_dir/.." && pwd)

if [[ $EUID -ne 0 ]]; then
  echo "Run this installer as root." >&2
  exit 1
fi
if [[ ! -f $binary ]]; then
  echo "Statistics binary not found: $binary" >&2
  exit 1
fi
if [[ ! $site_origin =~ ^https?://[^/]+$ ]]; then
  echo "SITE_ORIGIN must look like https://example.com (without a trailing slash)." >&2
  exit 1
fi

if ! getent passwd nabunana-stats >/dev/null; then
  useradd --system --home /var/lib/nabunana-stats --shell /usr/sbin/nologin nabunana-stats
fi
install -d -o root -g root -m 755 /opt/nabunana-stats
install -d -o nabunana-stats -g nabunana-stats -m 750 /var/lib/nabunana-stats
install -o root -g root -m 755 "$binary" /opt/nabunana-stats/nabunana-stats
install -o root -g root -m 644 "$repo_root/stats/deploy/nabunana-stats.service" /etc/systemd/system/nabunana-stats.service
install -o root -g root -m 755 "$repo_root/scripts/deploy-stats-release.sh" /usr/local/sbin/deploy-nabunana-stats

cat >/etc/nabunana-stats.env <<EOF
STATS_ADDRESS=127.0.0.1:8787
STATS_DATABASE=/var/lib/nabunana-stats/stats.db
STATS_SITE_ROOT=/var/www/blog
STATS_ALLOWED_ORIGIN=$site_origin
EOF
chmod 640 /etc/nabunana-stats.env
chown root:nabunana-stats /etc/nabunana-stats.env

if [[ -n $nginx_include_directory ]]; then
  if [[ ! -d $nginx_include_directory ]]; then
    echo "Nginx include directory does not exist: $nginx_include_directory" >&2
    exit 1
  fi
  install -o root -g root -m 644 "$repo_root/stats/deploy/nginx-location.conf" "$nginx_include_directory/nabunana-stats.conf"
else
  install -o root -g root -m 644 "$repo_root/stats/deploy/nginx-location.conf" /etc/nginx/snippets/nabunana-stats.conf
  echo "Add 'include /etc/nginx/snippets/nabunana-stats.conf;' inside the existing blog server block." >&2
fi

systemctl daemon-reload
nginx -t
echo "Statistics service files installed but not started. Preview and import Nginx logs before activation; see OPERATIONS.md."
