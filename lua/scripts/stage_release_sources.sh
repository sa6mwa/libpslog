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
manifest_path="lua/RELEASE_MANIFEST.in"
source_manifest_path="RELEASE_MANIFEST"
use_git=0

[[ -d "$repo_root" ]] || { printf 'Lua release staging repository is missing: %s\n' "$repo_root" >&2; exit 1; }
repo_root=$(CDPATH= cd -- "$repo_root" && pwd -P)
[[ -f "$version_header" ]] || { printf 'Lua release version header is missing: %s\n' "$version_header" >&2; exit 1; }
case "$version" in ''|*[!0-9A-Za-z.+-]*) printf 'invalid Lua release version: %s\n' "$version" >&2; exit 2;; esac

git_toplevel=$(git -C "$repo_root" rev-parse --show-toplevel 2>/dev/null || true)
if [[ -n "$git_toplevel" ]]; then
  git_toplevel=$(CDPATH= cd -- "$git_toplevel" && pwd -P)
fi

if [[ "$git_toplevel" == "$repo_root" ]] &&
   git -C "$repo_root" rev-parse --verify HEAD >/dev/null 2>&1 &&
   git -C "$repo_root" cat-file -e "HEAD:$manifest_path" 2>/dev/null; then
  use_git=1
elif [[ -f "$repo_root/$manifest_path" && -f "$repo_root/$source_manifest_path" && -f "$repo_root/VERSION" ]]; then
  archive_version=$(sed -n '1p' "$repo_root/VERSION")
  if [[ "$archive_version" != "$version" ]]; then
    printf 'Lua release source archive VERSION mismatch: expected %s, found %s\n' "$version" "$archive_version" >&2
    exit 1
  fi
else
  printf 'Lua release staging requires either committed Git metadata or source-archive VERSION, RELEASE_MANIFEST, and %s\n' "$manifest_path" >&2
  exit 1
fi

manifest_contains_entry() {
  local entry=$1
  grep -Fx -- "$entry" "$repo_root/$source_manifest_path" >/dev/null 2>&1
}

copy_manifest_entry() {
  local entry=$1
  if [[ "$use_git" -eq 1 ]]; then
    git -C "$repo_root" cat-file -e "HEAD:$entry" 2>/dev/null || { printf 'Lua release manifest file is missing from HEAD: %s\n' "$entry" >&2; exit 1; }
    git -C "$repo_root" show "HEAD:$entry" >"$stage_dir/$entry"
    entry_mode="$(git -C "$repo_root" ls-tree -r --format='%(objectmode)' HEAD -- "$entry")"
    if [[ "$entry_mode" == "100755" ]]; then
      chmod 755 "$stage_dir/$entry"
    fi
  else
    manifest_contains_entry "$entry" || { printf 'Lua release manifest file is not listed in source archive RELEASE_MANIFEST: %s\n' "$entry" >&2; exit 1; }
    [[ -f "$repo_root/$entry" ]] || { printf 'Lua release manifest file is missing from source archive: %s\n' "$entry" >&2; exit 1; }
    cp "$repo_root/$entry" "$stage_dir/$entry"
    if [[ -x "$repo_root/$entry" ]]; then
      chmod 755 "$stage_dir/$entry"
    fi
  fi
}

mkdir -p "$stage_dir"
staged_entries=()
while IFS= read -r entry || [[ -n "$entry" ]]; do
  [[ -z "$entry" || "$entry" == \#* ]] && continue
  case "$entry" in
    /*|*'..'*|*'//'*|*:* ) printf 'unsafe Lua release manifest path: %s\n' "$entry" >&2; exit 1;;
  esac
  mkdir -p "$stage_dir/$(dirname "$entry")"
  copy_manifest_entry "$entry"
  staged_entries+=("$entry")
done < <(if [[ "$use_git" -eq 1 ]]; then git -C "$repo_root" show "HEAD:$manifest_path"; else cat "$repo_root/$manifest_path"; fi)

mkdir -p "$stage_dir/include"
cp "$version_header" "$stage_dir/include/pslog_version.h"
cp "$stage_dir/lua/README.md" "$stage_dir/README.md"
staged_entries+=("include/pslog_version.h" "README.md" "VERSION" "RELEASE_MANIFEST")

printf '%s\n' "$version" >"$stage_dir/VERSION"
printf '%s\n' "${staged_entries[@]}" | LC_ALL=C sort -u >"$stage_dir/RELEASE_MANIFEST"
