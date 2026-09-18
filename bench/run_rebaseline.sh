#!/usr/bin/env bash

set -eu

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

compiler_from_host_cache() {
  local key=$1 cache_file="$repo_root/build/host/CMakeCache.txt" value
  [[ -f "$cache_file" ]] || return 1
  value=$(sed -n "s/^${key}:[^=]*=//p" "$cache_file" | tail -n 1)
  [[ -n "$value" ]] || return 1
  printf '%s\n' "$value"
}

fallback_cxx() {
  local host_override
  host_override=$(compiler_from_host_cache PSLOG_ALLOW_HOST_COMPILER || true)
  case "$(uname -s):${host_override}" in
    Darwin:*|*:1|*:ON|*:TRUE|*:YES) command -v c++ 2>/dev/null || true ;;
  esac
}

CC=${CC:-$(compiler_from_host_cache CMAKE_C_COMPILER || true)}
CXX=${CXX:-$(compiler_from_host_cache CMAKE_CXX_COMPILER || fallback_cxx)}
if [[ -z "$CC" || ! -x "$CC" || -z "$CXX" || ! -x "$CXX" ]]; then
  printf 'rebaseline requires configured host CC and available CXX; run make build-host first\n' >&2
  exit 1
fi
export CC CXX

if [[ "${PSLOG_REBASELINE_VALIDATE_ONLY:-}" = 1 ]]; then
  printf 'CC=%s\nCXX=%s\n' "$CC" "$CXX"
  exit 0
fi

exec "$repo_root/bench/run_perf_gate.sh" --freeze-baseline
