#!/usr/bin/env bash
set -euo pipefail

repo_root=$1
fixture_root=$(mktemp -d)
cleanup() { rm -rf "$fixture_root"; }
trap cleanup EXIT HUP INT TERM

source "$repo_root/scripts/cpkt-toolchains.sh"

archive_source="$fixture_root/fixture.tar.xz"
stage_root="$fixture_root/stage/fixture"
mkdir -p "$stage_root/bin"
touch "$stage_root/.ready"
tar -C "$fixture_root/stage" -cJf "$archive_source" fixture
expected_sha=$(sha256_file "$archive_source")
export CPKT_TOOLCHAIN_CACHE="$fixture_root/cache"

bootlin_values() {
  printf 'test|fixture|%s|fake|sysroot|%s\n' "$expected_sha" "$CPKT_TOOLCHAIN_CACHE/roots/fixture"
}
bootlin_ready() { [[ -f "$1/.ready" ]]; }
download_calls=0
download_file() {
  download_calls=$((download_calls + 1))
  cp "$archive_source" "$2"
}

mkdir -p "$CPKT_TOOLCHAIN_CACHE/archives"
printf 'corrupt cached archive\n' >"$CPKT_TOOLCHAIN_CACHE/archives/fixture.tar.xz"
install_bootlin fixture

[[ "$download_calls" -eq 1 ]] || { printf 'corrupt cache entry was not reacquired\n' >&2; exit 1; }
[[ "$(sha256_file "$CPKT_TOOLCHAIN_CACHE/archives/fixture.tar.xz")" == "$expected_sha" ]] || {
  printf 'reacquired cache entry was not verified\n' >&2
  exit 1
}
[[ -f "$CPKT_TOOLCHAIN_CACHE/roots/fixture/.ready" ]] || {
  printf 'reacquired toolchain archive was not extracted\n' >&2
  exit 1
}
