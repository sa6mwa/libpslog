#!/usr/bin/env bash
set -euo pipefail

repo_root=${1:?repo root is required}
cache_file="$repo_root/build/host/CMakeCache.txt"

if [[ ! -f "$cache_file" ]]; then
  printf 'rebaseline compiler test requires configured host cache: %s\n' "$cache_file" >&2
  exit 1
fi

host_cc=$(sed -n 's/^CMAKE_C_COMPILER:FILEPATH=//p' "$cache_file" | tail -n 1)
if [[ -z "$host_cc" || ! -x "$host_cc" ]]; then
  printf 'configured host C compiler is unavailable: %s\n' "${host_cc:-<empty>}" >&2
  exit 1
fi

cxx=${CXX:-$(command -v c++ 2>/dev/null || true)}
if [[ -z "$cxx" || ! -x "$cxx" ]]; then
  printf 'SKIP: no C++ compiler available for rebaseline fallback\n'
  exit 0
fi

output=$(env -u CXX PSLOG_REBASELINE_VALIDATE_ONLY=1 "$repo_root/bench/run_rebaseline.sh")
printf '%s\n' "$output" | grep -F "CC=$host_cc" >/dev/null
printf '%s\n' "$output" | grep -F "CXX=$cxx" >/dev/null
