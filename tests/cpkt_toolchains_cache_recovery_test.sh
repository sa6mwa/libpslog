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
lock_held=0
ready_checks=0
reject_lock=0
allow_unlocked_ready=0
flock() {
  [[ "$reject_lock" -eq 0 ]] || { printf 'ready toolchain unexpectedly acquired a cache lock\n' >&2; return 1; }
  if [[ "${1:-}" == '-u' ]]; then
    lock_held=0
  else
    lock_held=1
  fi
}
bootlin_ready() {
  ready_checks=$((ready_checks + 1))
  if [[ "$ready_checks" -gt 1 && "$lock_held" -ne 1 && "$allow_unlocked_ready" -ne 1 ]]; then
    printf 'toolchain readiness check ran outside its cache lock\n' >&2
    return 1
  fi
  [[ -f "$1/.ready" ]]
}
download_calls=0
download_file() {
  [[ "$lock_held" -eq 1 ]] || { printf 'toolchain download ran outside its cache lock\n' >&2; return 1; }
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
[[ -f "$CPKT_TOOLCHAIN_CACHE/locks/fixture.lock" ]] || {
  printf 'toolchain cache lock was not created\n' >&2
  exit 1
}
[[ "$lock_held" -eq 0 ]] || {
  printf 'toolchain cache lock was not released after installation\n' >&2
  exit 1
}
reject_lock=1
allow_unlocked_ready=1
install_bootlin fixture
