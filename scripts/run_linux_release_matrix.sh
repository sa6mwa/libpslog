#!/usr/bin/env bash

set -eu

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf 'missing required command: %s\n' "$1" >&2
        exit 1
    fi
}

require_file() {
    if [ ! -e "$1" ]; then
        printf 'missing required file: %s\n' "$1" >&2
        exit 1
    fi
}

run_target() {
    preset="$1"
    runtime_preset="$2"
    dev_preset="$3"

    printf '\n== %s ==\n' "$preset"
    cmake --preset "$preset"
    cmake --build --preset "$preset"
    ctest --preset "$preset"
    cmake --build --preset "$runtime_preset"
    cmake --build --preset "$dev_preset"
}

require_command cmake
require_command ctest
require_command musl-gcc
require_command aarch64-linux-gnu-gcc
require_command arm-linux-gnueabihf-gcc
require_command aarch64-linux-musl-gcc
require_command arm-linux-musleabihf-gcc
require_command qemu-aarch64
require_command qemu-arm

require_file "$HOME/.local/cross/aarch64-linux-musl/aarch64-linux-musl/lib/ld-musl-aarch64.so.1"
require_file "$HOME/.local/cross/arm-linux-musleabihf/arm-linux-musleabihf/lib/ld-musl-armhf.so.1"

cd "$repo_root"

cmake --preset host
cmake --build --preset package-clean-dist

run_target linux-gnu-release package-runtime-linux-gnu package-dev-linux-gnu
run_target linux-musl-release package-runtime-linux-musl package-dev-linux-musl
run_target aarch64-linux-gnu-release package-runtime-aarch64-linux-gnu package-dev-aarch64-linux-gnu
run_target aarch64-linux-musl-release package-runtime-aarch64-linux-musl package-dev-aarch64-linux-musl
run_target armhf-linux-gnu-release package-runtime-armhf-linux-gnu package-dev-armhf-linux-gnu
run_target armhf-linux-musl-release package-runtime-armhf-linux-musl package-dev-armhf-linux-musl

cmake --build --preset package-single-header
cmake --build --preset package-checksums

printf '\nLinux release matrix completed successfully.\n'
