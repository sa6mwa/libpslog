#!/usr/bin/env bash
set -euo pipefail

repo_root=${1:?repo root is required}
cache_file="$repo_root/build/host/CMakeCache.txt"

if [[ ! -f "$cache_file" ]]; then
  printf 'host runner test requires configured host cache: %s\n' "$cache_file" >&2
  exit 1
fi
if sed -n 's/^CMAKE_SYSROOT:[^=]*=//p' "$cache_file" | grep -q .; then
  printf 'SKIP: host preset is sysroot-backed in this configuration\n'
  exit 0
fi

output=$("$repo_root/scripts/run_host_binary.sh" /bin/sh -c 'printf native-host-runner')
if [[ "$output" != native-host-runner ]]; then
  printf 'native host runner produced unexpected output: %s\n' "$output" >&2
  exit 1
fi
