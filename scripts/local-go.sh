#!/usr/bin/env bash
# Build-time policy only: resulting Go/cgo executables execute directly.
set -euo pipefail
repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cache="$repo_root/build/host/CMakeCache.txt"
read_cache() { sed -n "s/^$1:[^=]*=//p" "$cache" | tail -n 1; }
export CC="$(read_cache CMAKE_C_COMPILER)" CXX="$(read_cache CMAKE_CXX_COMPILER)"
[[ -x "$CC" ]] || { echo 'Run make build-host before building Go/cgo consumers' >&2; exit 1; }
command=${1:?go subcommand required}
shift
flags=()
if [[ $(uname -s) == Linux ]]; then
    runtime=$(read_cache PSLOG_NONSHIPPED_ELF_LINKER_FLAGS)
    [[ -n "$runtime" ]] || { echo 'Missing Bootlin executable runtime flags' >&2; exit 1; }
    export CGO_ENABLED=1
    flags=("-ldflags=-linkmode=external -extldflags '$runtime'")
fi
exec go "$command" "${flags[@]}" "$@"
