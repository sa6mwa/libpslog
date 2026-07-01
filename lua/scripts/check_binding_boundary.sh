#!/usr/bin/env sh
set -eu

tree="${1:-build/luarocks}"
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "${script_dir}/../.." && pwd)
public_header="${2:-${repo_root}/include/pslog.h}"
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

for tool in nm readelf ldd; do
    if ! command -v "${tool}" >/dev/null 2>&1; then
        printf 'lua binding boundary check: skipped; missing inspection tool: %s\n' "${tool}" >&2
        exit 0
    fi
done

if ! readelf -h "${core_so}" >/dev/null 2>&1; then
    printf 'lua binding boundary check: skipped; %s is not an ELF shared object\n' "${core_so}" >&2
    exit 0
fi

if ! nm -D --defined-only "${core_so}" | grep -Eq '[[:space:]]luaopen_pslog_core$'; then
    printf 'lua binding boundary check: %s does not export luaopen_pslog_core\n' "${core_so}" >&2
    exit 1
fi

if nm -D --defined-only "${core_so}" | awk '{ print $3 }' | grep -Eq '^pslog(_|$)'; then
    printf 'lua binding boundary check: %s defines libpslog symbols\n' "${core_so}" >&2
    nm -D --defined-only "${core_so}" | awk '{ print $3 }' | grep -E '^pslog(_|$)' >&2
    exit 1
fi

if [ ! -f "${public_header}" ]; then
    printf 'lua binding boundary check: missing public header %s\n' "${public_header}" >&2
    exit 1
fi

public_symbols=$(mktemp)
referenced_symbols=$(mktemp)
private_symbols=$(mktemp)
trap 'rm -f "${public_symbols}" "${referenced_symbols}" "${private_symbols}"' EXIT HUP INT TERM

{
    sed -nE 's/.*PSLOG_API[[:space:]].*[ *]((pslog|pslog_[A-Za-z0-9_]+))[[:space:]]*\(.*/\1/p' "${public_header}"
    sed -nE 's/.*PSLOG_API[[:space:]]+extern[[:space:]].*[ *]((pslog|pslog_[A-Za-z0-9_]+))[[:space:]]*;.*/\1/p' "${public_header}"
} | sort -u > "${public_symbols}"

nm -D --undefined-only "${core_so}" |
    awk '{ print $NF }' |
    grep -E '^pslog(_|$)' |
    sort -u > "${referenced_symbols}" || true

grep -Fvx -f "${public_symbols}" "${referenced_symbols}" > "${private_symbols}" || true
if [ -s "${private_symbols}" ]; then
    printf 'lua binding boundary check: %s references private libpslog symbols\n' "${core_so}" >&2
    cat "${private_symbols}" >&2
    exit 1
fi

if ! readelf -d "${core_so}" | grep -Eq 'Shared library: \[libpslog\.so(\.[0-9]+)*\]'; then
    printf 'lua binding boundary check: %s does not declare a libpslog.so dependency\n' "${core_so}" >&2
    readelf -d "${core_so}" >&2
    exit 1
fi

if ldd "${core_so}" | grep -F 'not found' >/dev/null 2>&1; then
    printf 'lua binding boundary check: unresolved shared library dependency in %s\n' "${core_so}" >&2
    ldd "${core_so}" >&2
    exit 1
fi
