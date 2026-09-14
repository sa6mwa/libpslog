#!/usr/bin/env sh
set -eu

tree="${1:-build/luarocks}"
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "${script_dir}/../.." && pwd)
public_header="${2:-${repo_root}/include/pslog.h}"
lua_public_header="${3:-${repo_root}/include/pslog_lua.h}"
core_so="${tree}/lib/lua/5.5/pslog/core.so"

case "$(uname -s 2>/dev/null || printf unknown)" in
    Linux) ;;
    *)
        printf 'lua binding boundary check: skipped on non-ELF host\n' >&2
        exit 0
        ;;
esac

if [ ! -f "${core_so}" ]; then
    printf 'lua binding boundary check: missing %s\n' "${core_so}" >&2
    exit 1
fi

cache="${repo_root}/build/host/CMakeCache.txt"
nm_bin=$(sed -n 's/^CMAKE_NM:[^=]*=//p' "${cache}" | tail -n 1)
readelf_bin=$(sed -n 's/^CMAKE_READELF:[^=]*=//p' "${cache}" | tail -n 1)
if [ ! -x "${nm_bin}" ] || [ ! -x "${readelf_bin}" ]; then
    printf 'lua boundary check requires configured Bootlin inspection tools\n' >&2
    exit 1
fi

if ! "${readelf_bin}" -h "${core_so}" >/dev/null 2>&1; then
    printf 'lua binding boundary check: skipped; %s is not an ELF shared object\n' "${core_so}" >&2
    exit 0
fi

if ! "${nm_bin}" -D --defined-only "${core_so}" | grep -Eq '[[:space:]]luaopen_pslog_core$'; then
    printf 'lua binding boundary check: %s does not export luaopen_pslog_core\n' "${core_so}" >&2
    exit 1
fi

if "${nm_bin}" -D --defined-only "${core_so}" | awk '{ print $3 }' | grep -Ev '^pslog_lua_' | grep -Eq '^pslog(_|$)'; then
    printf 'lua binding boundary check: %s defines libpslog symbols\n' "${core_so}" >&2
    "${nm_bin}" -D --defined-only "${core_so}" | awk '{ print $3 }' | grep -Ev '^pslog_lua_' | grep -E '^pslog(_|$)' >&2
    exit 1
fi

if [ ! -f "${public_header}" ]; then
    printf 'lua binding boundary check: missing public header %s\n' "${public_header}" >&2
    exit 1
fi
if [ ! -f "${lua_public_header}" ]; then
    printf 'lua binding boundary check: missing Lua interop header %s\n' "${lua_public_header}" >&2
    exit 1
fi

public_symbols=$(mktemp)
referenced_symbols=$(mktemp)
private_symbols=$(mktemp)
trap 'rm -f "${public_symbols}" "${referenced_symbols}" "${private_symbols}"' EXIT HUP INT TERM

{
    sed -nE 's/.*PSLOG_API[[:space:]].*[ *]((pslog|pslog_[A-Za-z0-9_]+))[[:space:]]*\(.*/\1/p' "${public_header}"
    sed -nE 's/.*PSLOG_API[[:space:]]+extern[[:space:]].*[ *]((pslog|pslog_[A-Za-z0-9_]+))[[:space:]]*;.*/\1/p' "${public_header}"
    sed -nE 's/.*PSLOG_API[[:space:]].*[ *]((pslog_lua_[A-Za-z0-9_]+))[[:space:]]*\(.*/\1/p' "${lua_public_header}"
} | sort -u > "${public_symbols}"

"${nm_bin}" -D --undefined-only "${core_so}" |
    awk '{ print $NF }' |
    grep -E '^pslog(_|$)' |
    sort -u > "${referenced_symbols}" || true

grep -Fvx -f "${public_symbols}" "${referenced_symbols}" > "${private_symbols}" || true
if [ -s "${private_symbols}" ]; then
    printf 'lua binding boundary check: %s references private libpslog symbols\n' "${core_so}" >&2
    cat "${private_symbols}" >&2
    exit 1
fi

if ! "${readelf_bin}" -d "${core_so}" | grep -Eq 'Shared library: \[libpslog\.so(\.[0-9]+)*\]'; then
    printf 'lua binding boundary check: %s does not declare a libpslog.so dependency\n' "${core_so}" >&2
    "${readelf_bin}" -d "${core_so}" >&2
    exit 1
fi
