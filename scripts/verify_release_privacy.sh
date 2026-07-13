#!/usr/bin/env bash
# Verify checksum-listed artifacts with the inspection tools from their producing builds.
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
usage() {
  printf 'usage: %s --build-dir DIR --target-id TARGET_ID [--darwin-build-dir DIR]\n' "$0" >&2
  exit 2
}

build_dir=
target_id=
darwin_build_dir=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --build-dir) [[ $# -ge 2 ]] || usage; build_dir=$2; shift 2 ;;
    --target-id) [[ $# -ge 2 ]] || usage; target_id=$2; shift 2 ;;
    --darwin-build-dir) [[ $# -ge 2 ]] || usage; darwin_build_dir=$2; shift 2 ;;
    *) usage ;;
  esac
done
[[ -n "$build_dir" && -n "$target_id" ]] || usage

eval "$("$repo_root/scripts/discover_target_tools.sh" --build-dir "$build_dir" --target-id "$target_id")"
linux_readelf=$READELF
if [[ -n "$darwin_build_dir" ]]; then
  eval "$("$repo_root/scripts/discover_target_tools.sh" --build-dir "$darwin_build_dir" --target-id arm64-apple-darwin)"
fi

version=$("$repo_root/lua/scripts/release_version.sh")
checksum_file="$repo_root/dist/libpslog-${version}-CHECKSUMS"
[[ -f "$checksum_file" ]] || {
  printf 'verify_release_privacy: checksum manifest is missing: %s\n' "$checksum_file" >&2
  exit 1
}
(cd "$repo_root/dist" && sha256sum -c "$(basename -- "$checksum_file")")

exec cmake \
  -DPSLOG_ROOT="$repo_root" \
  -DPSLOG_BINARY_DIR="$build_dir" \
  -DPSLOG_VERSION="$version" \
  -DPSLOG_READELF="$linux_readelf" \
  -DPSLOG_OTOOL="$OTOOL" \
  -P "$repo_root/cmake/check_release_privacy.cmake"
