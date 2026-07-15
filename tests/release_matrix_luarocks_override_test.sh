#!/usr/bin/env bash
set -euo pipefail

repo_root=$1
fixture_root=$(mktemp -d)
cleanup() { rm -rf "$fixture_root"; }
trap cleanup EXIT HUP INT TERM

fake_repo="$fixture_root/repo"
fake_bin="$fixture_root/bin"
fake_home="$fixture_root/home"
mkdir -p "$fake_repo/scripts" "$fake_repo/lua/scripts" "$fake_bin" "$fake_home"
cp "$repo_root/scripts/run_linux_release_matrix.sh" "$fake_repo/scripts/"

printf '%s\n' '#!/bin/sh' 'label=$1; shift' 'exec "$@"' >"$fake_repo/scripts/run_timed.sh"
printf '%s\n' '#!/bin/sh' 'printf "0.0.0\\n"' >"$fake_repo/lua/scripts/release_version.sh"
printf '%s\n' '#!/bin/sh' 'exit 0' >"$fake_repo/scripts/verify_release_privacy.sh"
printf '%s\n' '#!/bin/sh' \
  'if [ "$1" = discover ] && [ "$2" = arm64-apple-darwin ]; then printf "%s\n" target=arm64-apple-darwin status=missing note=partial-osxcross; exit 0; fi' \
  'exit 1' >"$fake_repo/scripts/cpkt-toolchains.sh"
printf '%s\n' '#!/bin/sh' \
  'for arg in "$@"; do case "$arg" in *arm64-apple-darwin*) exit 92 ;; esac; done' \
  'exit 0' >"$fake_bin/cmake"
chmod +x "$fake_bin/cmake"
for command in ctest qemu-aarch64 qemu-arm; do
  printf '%s\n' '#!/bin/sh' 'exit 0' >"$fake_bin/$command"
  chmod +x "$fake_bin/$command"
done
printf '%s\n' '#!/bin/sh' \
  'while [ "$#" -gt 0 ]; do case "$1" in *=*) export "$1"; shift ;; *) break ;; esac; done' \
  'exec "$@"' >"$fake_bin/env"
printf '%s\n' '#!/bin/sh' \
  '[ "$1" = "--" ] && shift' \
  'printf "%s\\n" "${1%/*}"' >"$fake_bin/dirname"
printf '%s\n' '#!/bin/sh' \
  '[ "$LUA_ROCKS" = "$EXPECTED_LUA_ROCKS" ] || exit 91' \
  'exit 0' >"$fake_bin/make"
printf '%s\n' '#!/bin/sh' 'exit 0' >"$fake_bin/custom-luarocks"
chmod +x "$fake_repo/scripts/run_timed.sh" "$fake_repo/lua/scripts/release_version.sh" \
  "$fake_repo/scripts/verify_release_privacy.sh" "$fake_repo/scripts/cpkt-toolchains.sh" \
  "$fake_bin/env" "$fake_bin/dirname" "$fake_bin/make" "$fake_bin/custom-luarocks"

mkdir -p "$fake_home/.local/cross/osxcross/bin"
printf '%s\n' '#!/bin/sh' 'exit 0' >"$fake_home/.local/cross/osxcross/bin/arm64-apple-darwin25-clang"
chmod +x "$fake_home/.local/cross/osxcross/bin/arm64-apple-darwin25-clang"

# Deliberately omit a command named "luarocks": the matrix must accept the
# explicit lifecycle override and pass it into the Lua artifact target.
PATH="$fake_bin" HOME="$fake_home" LUA_ROCKS="$fake_bin/custom-luarocks" \
  EXPECTED_LUA_ROCKS="$fake_bin/custom-luarocks" /bin/bash "$fake_repo/scripts/run_linux_release_matrix.sh"
