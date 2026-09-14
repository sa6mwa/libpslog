#!/usr/bin/env bash
set -euo pipefail

repo_root=$1
fixture_root=$(mktemp -d "$repo_root/build/afl-resolver-test.XXXXXX")
cleanup() { rm -rf "$fixture_root"; }
trap cleanup EXIT HUP INT TERM

source "$repo_root/scripts/cpkt-aflpp.sh"

# Test the repository resolver's publication readiness and failure cleanup.
ready_root="$fixture_root/ready"
mkdir -p "$ready_root/bin" "$ready_root/lib/afl"
for tool in afl-fuzz afl-showmap cpkt-afl-gcc cpkt-afl-g++; do
  printf '#!/bin/sh\nexit 0\n' > "$ready_root/bin/$tool"
  chmod +x "$ready_root/bin/$tool"
done
touch "$ready_root/lib/afl/afl-gcc-pass.so" "$ready_root/lib/afl/afl-compiler-rt.o" "$ready_root/.cpkt-aflpp-revision-$revision-fixture"
ready "$ready_root" fixture
rm "$ready_root/bin/afl-showmap"
if ready "$ready_root" fixture; then
  echo 'Resolver accepted a missing afl-showmap' >&2
  exit 1
fi
mkdir "$fixture_root/interrupted"
if bash -c 'source "$1/scripts/cpkt-aflpp.sh"; install_cleanup_trap -rf "$2"; exit 7' _ "$repo_root" "$fixture_root/interrupted"; then
  echo 'Cleanup lost the failing exit status' >&2
  exit 1
fi
[[ ! -e "$fixture_root/interrupted" ]]

export CPKT_TOOLCHAIN_CACHE="$fixture_root/cache"
unset CPKT_TOOLCHAIN_LOCK_TIMEOUT
bootlin_description() {
  printf 'target=x86_64-linux-gnu\nroot=%s\nsysroot=%s\ncc=/fixture/cc\ncxx=/fixture/cxx\n' \
    "$fixture_root/cache/roots/x86-64--glibc--stable-2026.08-1" \
    "$fixture_root/cache/roots/x86-64--glibc--stable-2026.08-1/x86_64-buildroot-linux-gnu/sysroot"
}
lock_held=0
ready_calls=0
flock() {
  if [[ "${1:-}" == '-w' ]]; then
    [[ "${2:-}" == 600 ]] || { printf 'AFL++ cache lock did not use the lifecycle timeout\n' >&2; return 1; }
    shift 2
  fi
  if [[ "${1:-}" == '-u' ]]; then
    lock_held=0
  else
    lock_held=1
  fi
}
ready() {
  ready_calls=$((ready_calls + 1))
  if [[ "$ready_calls" -eq 1 ]]; then
    return 1
  fi
  [[ "$lock_held" -eq 1 ]] || {
    printf 'AFL++ readiness was not rechecked under its cache lock\n' >&2
    return 1
  }
  return 0
}
uname() {
  case "${1:-}" in
    -s) printf 'Linux\n' ;;
    -m) printf 'x86_64\n' ;;
    *) command uname "$@" ;;
  esac
}

ensure

[[ "$ready_calls" -eq 2 ]] || {
  printf 'AFL++ cache readiness was not rechecked after locking\n' >&2
  exit 1
}
[[ "$lock_held" -eq 0 ]] || {
  printf 'AFL++ cache lock was not released after a ready recheck\n' >&2
  exit 1
}
[[ -f "$CPKT_TOOLCHAIN_CACHE/locks/aflplusplus-${version}-x86_64-linux-gnu.lock" ]] || {
  printf 'AFL++ versioned cache lock was not created\n' >&2
  exit 1
}
