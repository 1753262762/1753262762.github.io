#!/usr/bin/env bash

set -Eeuo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
test_root=$(mktemp -d)
trap 'rm -rf -- "$test_root"' EXIT

target="$test_root/blog"
mkdir -p "$target/media/music/album"
printf 'old site\n' > "$target/index.html"
printf 'audio\n' > "$target/media/music/album/01.mp3"

make_archive() {
  local label=$1
  local build="$test_root/build-$label"
  local archive="$test_root/site-$label.tar.gz"

  mkdir -p "$build/media/music/album"
  printf '%s\n' "$label" > "$build/index.html"
  printf 'lyrics\n' > "$build/media/music/album/01.lrc"
  tar -czf "$archive" -C "$build" .
  printf '%s\n' "$archive"
}

release_one="aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-1"
archive_one=$(make_archive first)
bash "$script_dir/deploy-release.sh" "$archive_one" "$target" "$release_one"

test -L "$target"
test "$(cat "$target/index.html")" = first
test -f "$target/media/music/album/01.mp3"
test -f "$target/media/music/album/01.lrc"
test "$(find "$test_root/.blog-backups" -mindepth 1 -maxdepth 1 -type d | wc -l)" -eq 1

release_two="bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb-1"
archive_two=$(make_archive second)
previous=$(readlink -f "$target")
bash "$script_dir/deploy-release.sh" "$archive_two" "$target" "$release_two"

test -L "$target"
test "$(cat "$target/index.html")" = second
test -f "$target/media/music/album/01.mp3"
test "$(cat "$target/PREVIOUS_RELEASE")" = "$previous"

echo "deployment release test passed"
