#!/usr/bin/env bash
# Run a host-built program. Bootlin-backed host builds execute through their
# configured sysroot loader; native host builds execute directly.
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
cache_file="$repo_root/build/host/CMakeCache.txt"
sysroot=
if [[ -f "$cache_file" ]]; then
  sysroot=$(sed -n 's/^CMAKE_SYSROOT:[^=]*=//p' "$cache_file" | tail -n 1)
fi

if [[ -n "$sysroot" ]]; then
  exec "$repo_root/scripts/run_sysroot_binary.sh" \
    --build-dir "$repo_root/build/host" \
    "$@"
fi

if [[ $# -lt 1 ]]; then
  printf 'run_host_binary: missing program\n' >&2
  exit 2
fi
exec "$@"
