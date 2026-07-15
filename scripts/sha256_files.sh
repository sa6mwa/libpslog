#!/usr/bin/env sh
# Print sha256 checksums for one or more files using tools available on both
# GNU/Linux and stock macOS hosts.
set -eu

if [ "$#" -lt 1 ]; then
    printf 'usage: %s FILE [FILE ...]\n' "$0" >&2
    printf '       %s --check MANIFEST\n' "$0" >&2
    exit 2
fi

if command -v sha256sum >/dev/null 2>&1; then
    if [ "$1" = "--check" ]; then
        shift
        exec sha256sum -c "$@"
    fi
    exec sha256sum "$@"
fi

if command -v shasum >/dev/null 2>&1; then
    if [ "$1" = "--check" ]; then
        shift
        exec shasum -a 256 -c "$@"
    fi
    exec shasum -a 256 "$@"
fi

printf 'sha256_files: need sha256sum or shasum in PATH\n' >&2
exit 1
