#!/usr/bin/env bash
# Stage the standalone Lua facade source package from its explicit manifest.
set -euo pipefail

if [[ $# -ne 4 ]]; then
  printf 'usage: %s REPOSITORY_ROOT STAGE_DIR VERSION GENERATED_VERSION_HEADER\n' "$0" >&2
  exit 2
fi

repo_root=$1
stage_dir=$2
version=$3
version_header=$4
manifest="$repo_root/lua/RELEASE_MANIFEST.in"

[[ -d "$repo_root" ]] || { printf 'Lua release staging repository is missing: %s\n' "$repo_root" >&2; exit 1; }
[[ -f "$manifest" ]] || { printf 'Lua release manifest is missing: %s\n' "$manifest" >&2; exit 1; }
[[ -f "$version_header" ]] || { printf 'Lua release version header is missing: %s\n' "$version_header" >&2; exit 1; }
git -C "$repo_root" rev-parse --verify HEAD >/dev/null 2>&1 || { printf 'Lua release staging requires a committed HEAD: %s\n' "$repo_root" >&2; exit 1; }
case "$version" in ''|*[!0-9A-Za-z.+-]*) printf 'invalid Lua release version: %s\n' "$version" >&2; exit 2;; esac

mkdir -p "$stage_dir"
staged_entries=()
while IFS= read -r entry || [[ -n "$entry" ]]; do
  [[ -z "$entry" || "$entry" == \#* ]] && continue
  case "$entry" in
    /*|*'..'*|*'//'*|*:* ) printf 'unsafe Lua release manifest path: %s\n' "$entry" >&2; exit 1;;
  esac
  git -C "$repo_root" cat-file -e "HEAD:$entry" 2>/dev/null || { printf 'Lua release manifest file is missing from HEAD: %s\n' "$entry" >&2; exit 1; }
  mkdir -p "$stage_dir/$(dirname "$entry")"
  git -C "$repo_root" show "HEAD:$entry" >"$stage_dir/$entry"
  entry_mode="$(git -C "$repo_root" ls-tree -r --format='%(objectmode)' HEAD -- "$entry")"
  if [[ "$entry_mode" == "100755" ]]; then
    chmod 755 "$stage_dir/$entry"
  fi
  staged_entries+=("$entry")
done <"$manifest"

mkdir -p "$stage_dir/include"
cp "$version_header" "$stage_dir/include/pslog_version.h"
cp "$stage_dir/lua/README.md" "$stage_dir/README.md"
staged_entries+=("include/pslog_version.h" "README.md" "VERSION" "RELEASE_MANIFEST")

printf '%s\n' "$version" >"$stage_dir/VERSION"
printf '%s\n' "${staged_entries[@]}" | LC_ALL=C sort -u >"$stage_dir/RELEASE_MANIFEST"
