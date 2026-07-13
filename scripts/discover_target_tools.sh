#!/usr/bin/env bash
# Report target inspection/mutation tools from the CMake build that produced an artifact.
set -euo pipefail

usage() {
  printf 'usage: %s --build-dir DIR --target-id TARGET_ID\n' "$0" >&2
  exit 2
}

build_dir=
target_id=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --build-dir) [[ $# -ge 2 ]] || usage; build_dir=$2; shift 2 ;;
    --target-id) [[ $# -ge 2 ]] || usage; target_id=$2; shift 2 ;;
    *) usage ;;
  esac
done
[[ -n "$build_dir" && -n "$target_id" ]] || usage
cache_file="$build_dir/CMakeCache.txt"
[[ -f "$cache_file" ]] || { printf 'discover_target_tools: missing CMake cache: %s\n' "$cache_file" >&2; exit 1; }

cache_value() {
  sed -n "s|^$1:[^=]*=||p" "$cache_file" | tail -n 1
}
required_cache_value() {
  local value
  value=$(cache_value "$1")
  [[ -n "$value" && "$value" != *-NOTFOUND ]] || {
    printf 'discover_target_tools: %s is unavailable in %s for %s\n' "$1" "$cache_file" "$target_id" >&2
    exit 1
  }
  [[ -x "$value" ]] || {
    printf 'discover_target_tools: configured %s is not executable for %s: %s\n' "$1" "$target_id" "$value" >&2
    exit 1
  }
  printf '%s\n' "$value"
}

cc=$(required_cache_value CMAKE_C_COMPILER)
strip=$(required_cache_value CMAKE_STRIP)
readelf=$(cache_value CMAKE_READELF)
install_name_tool=$(cache_value CMAKE_INSTALL_NAME_TOOL)
otool=$(cache_value CPKT_OTOOL)

if [[ "$target_id" == *darwin* ]]; then
  install_name_tool=$(required_cache_value CMAKE_INSTALL_NAME_TOOL)
  otool=$(required_cache_value CPKT_OTOOL)
else
  readelf=$(required_cache_value CMAKE_READELF)
fi

printf 'CC=%q\nSTRIP=%q\nINSTALL_NAME_TOOL=%q\nOTOOL=%q\nREADELF=%q\nTARGET_ID=%q\n' \
  "$cc" "$strip" "$install_name_tool" "$otool" "$readelf" "$target_id"
