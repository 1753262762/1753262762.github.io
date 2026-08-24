#!/usr/bin/env bash

set -Eeuo pipefail
umask 022

archive=${1:?Usage: deploy-release.sh ARCHIVE TARGET RELEASE_ID}
target=${2:?Usage: deploy-release.sh ARCHIVE TARGET RELEASE_ID}
release_id=${3:?Usage: deploy-release.sh ARCHIVE TARGET RELEASE_ID}

case "$target" in
  /*) ;;
  *) echo "Deploy target must be an absolute path: $target" >&2; exit 2 ;;
esac

if [[ "$target" == "/" || "$target" == "/var" || "$target" == "/var/www" ]]; then
  echo "Refusing unsafe deploy target: $target" >&2
  exit 2
fi

if [[ ! "$release_id" =~ ^[0-9a-f]{40}-[1-9][0-9]*$ ]]; then
  echo "Invalid release id: $release_id" >&2
  exit 2
fi

parent=$(dirname -- "$target")
site_name=$(basename -- "$target")
releases="$parent/.${site_name}-releases"
backups="$parent/.${site_name}-backups"
release="$releases/$release_id"
staging="$release.staging"
next_link="$parent/.${site_name}.next-$release_id"

cleanup() {
  rm -f -- "$archive" "$next_link"
  if [[ -d "$staging" ]]; then
    rm -rf -- "$staging"
  fi
  if [[ -n "${legacy_backup:-}" && ! -e "$target" && ! -L "$target" && -d "$legacy_backup" ]]; then
    mv -- "$legacy_backup" "$target"
    echo "Activation failed; restored $target from $legacy_backup" >&2
  fi
}
trap cleanup EXIT

mkdir -p -- "$releases" "$backups"
if [[ -e "$release" ]]; then
  echo "Release already exists: $release" >&2
  exit 1
fi

mkdir -- "$staging"
tar -xzf "$archive" -C "$staging"

if [[ ! -f "$staging/index.html" ]]; then
  echo "Build archive does not contain index.html" >&2
  exit 1
fi

# MP3 files are intentionally not committed to Git. Carry them forward from the
# active server release so a code deployment cannot silently erase the library.
audio_source="$target/media/music"
audio_destination="$staging/media/music"
mkdir -p -- "$audio_destination"
if [[ -d "$audio_source" ]]; then
  (
    cd -- "$audio_source"
    find . -type f -name '*.mp3' -exec cp --parents -t "$audio_destination" -- {} +
  )
fi

audio_count=$(find "$audio_destination" -type f -name '*.mp3' 2>/dev/null | wc -l)
if [[ "$audio_count" -eq 0 ]]; then
  echo "No MP3 files were found in the active site; refusing to publish a broken music library." >&2
  echo "Upload the MP3 library to $target/media/music once, then rerun this workflow." >&2
  exit 1
fi

mv -- "$staging" "$release"
printf '%s\n' "$release_id" > "$release/RELEASE"
ln -s -- "$release" "$next_link"

if [[ -L "$target" ]]; then
  previous_release=$(readlink -f -- "$target" || true)
  if [[ -n "$previous_release" ]]; then
    printf '%s\n' "$previous_release" > "$release/PREVIOUS_RELEASE"
  fi
  mv -Tf -- "$next_link" "$target"
elif [[ -d "$target" ]]; then
  backup="$backups/legacy-$(date -u +%Y%m%dT%H%M%SZ)"
  mv -- "$target" "$backup"
  legacy_backup="$backup"
  printf '%s\n' "$backup" > "$release/PREVIOUS_RELEASE"
  mv -T -- "$next_link" "$target"
  legacy_backup=""
elif [[ -e "$target" ]]; then
  echo "Deploy target exists but is not a directory or symlink: $target" >&2
  exit 1
else
  mv -T -- "$next_link" "$target"
fi

echo "Activated $release with $audio_count MP3 files"
