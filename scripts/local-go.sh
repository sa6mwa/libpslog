#!/usr/bin/env bash
# Build-time policy only: resulting Go/cgo executables execute directly.
set -euo pipefail
repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cache="$repo_root/build/host/CMakeCache.txt"
read_cache() { sed -n "s/^$1:[^=]*=//p" "$cache" | tail -n 1; }
host_compiler_override=$(read_cache PSLOG_ALLOW_HOST_COMPILER)
cc=$(read_cache CMAKE_C_COMPILER)
cxx=$(read_cache CMAKE_CXX_COMPILER)
if [[ -z "$cxx" && "$host_compiler_override" =~ ^(1|ON|TRUE|YES)$ ]]; then
    cxx=$(command -v c++ 2>/dev/null || true)
fi
export CC="$cc" CXX="$cxx"
[[ -x "$CC" ]] || { echo 'Run make build-host before building Go/cgo consumers' >&2; exit 1; }
command=${1:?go subcommand required}
shift
flags=()
if [[ $(uname -s) == Linux ]]; then
    runtime=$(read_cache PSLOG_NONSHIPPED_ELF_LINKER_FLAGS)
    export CGO_ENABLED=1
    if [[ -n "$runtime" ]]; then
        flags=("-ldflags=-linkmode=external -extldflags '$runtime'")
    elif [[ ! "$host_compiler_override" =~ ^(1|ON|TRUE|YES)$ ]]; then
        echo 'Missing Bootlin executable runtime flags' >&2
        exit 1
    fi
fi
if [[ -n ${flags[@]+present} ]]; then
    exec go "$command" "${flags[@]}" "$@"
fi
exec go "$command" "$@"
