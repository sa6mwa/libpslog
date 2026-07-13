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
    archive_target="$2"

    printf '\n== %s ==\n' "$preset"
    cmake --preset "$preset"
    cmake --build --preset "$preset"
    ctest --preset "$preset"
    cmake --build "build/$preset" --target "$archive_target"
}

run_build_only_target() {
    preset="$1"
    archive_target="$2"

    printf '\n== %s ==\n' "$preset"
    cmake --preset "$preset"
    cmake --build --preset "$preset"
    cmake --build "build/$preset" --target "$archive_target"
}

require_command cmake
require_command ctest
require_command make
require_command luarocks
require_command qemu-aarch64
require_command qemu-arm

darwin_toolchain=""
darwin_bin="${OSXCROSS_ROOT:-$HOME/.local/cross/osxcross}/bin"
for cc in "$darwin_bin"/arm64-apple-darwin*-clang; do
    if [ -x "$cc" ]; then
        darwin_toolchain="$cc"
        break
    fi
done

cd "$repo_root"

cmake --preset host
cmake --build build/host --target package-clean-dist

run_target x86_64-linux-gnu-release package-archive
run_target x86_64-linux-musl-release package-archive
run_target aarch64-linux-gnu-release package-archive
run_target aarch64-linux-musl-release package-archive
run_target armhf-linux-gnu-release package-archive
run_target armhf-linux-musl-release package-archive
if [ -x "$darwin_toolchain" ]; then
    run_build_only_target arm64-apple-darwin-release package-archive
else
    printf '\n== arm64-apple-darwin-release ==\n'
    printf 'Skipping Darwin release target: osxcross arm64 clang not available under %s\n' "$darwin_bin"
fi

cmake --build build/host --target package-single-header
cmake --build build/host --target package-source
cmake -DPSLOG_ROOT="$repo_root" -DPSLOG_BINARY_DIR="$repo_root/build/host" -DPSLOG_VERSION="$(./lua/scripts/release_version.sh)" -DPSLOG_TOOLCHAIN_RELATIVE=cmake/toolchains/linux-x86_64-gnu.cmake -P tests/source_archive_smoke_test.cmake
make release-lua-artifacts
cmake --build build/host --target package-checksums
if [ -n "$darwin_toolchain" ]; then
    ./scripts/verify_release_privacy.sh --build-dir build/host --target-id x86_64-linux-gnu --darwin-build-dir build/arm64-apple-darwin-release
else
    ./scripts/verify_release_privacy.sh --build-dir build/host --target-id x86_64-linux-gnu
fi

printf '\nRelease matrix completed successfully.\n'
