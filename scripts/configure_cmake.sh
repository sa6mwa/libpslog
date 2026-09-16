#!/usr/bin/env bash
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)

die() {
    printf 'configure_cmake: %s\n' "$*" >&2
    exit 2
}

host_target() {
    case "$(uname -s):$(uname -m)" in
        Linux:x86_64|Linux:amd64) printf '%s\n' x86_64-linux-gnu ;;
        Linux:aarch64|Linux:arm64) printf '%s\n' aarch64-linux-gnu ;;
        Linux:armv7l) printf '%s\n' armhf-linux-gnu ;;
        Darwin:*) printf '%s\n' none ;;
        *) die "unsupported host: $(uname -s) $(uname -m)" ;;
    esac
}

preset_target() {
    case "$1" in
        debug|debug-lua|host|coverage|profile|format) host_target ;;
        x86_64-linux-gnu-release|valgrind|fuzz) printf '%s\n' x86_64-linux-gnu ;;
        x86_64-linux-musl-release) printf '%s\n' x86_64-linux-musl ;;
        aarch64-linux-gnu-release) printf '%s\n' aarch64-linux-gnu ;;
        aarch64-linux-musl-release) printf '%s\n' aarch64-linux-musl ;;
        armhf-linux-gnu-release) printf '%s\n' armhf-linux-gnu ;;
        armhf-linux-musl-release) printf '%s\n' armhf-linux-musl ;;
        arm64-apple-darwin-release) printf '%s\n' none ;;
        *) die "unknown preset: $1" ;;
    esac
}

cache_compiler() {
    local cache=$1
    sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' "$cache" | tail -n 1
}

cache_toolchain() {
    local cache=$1
    sed -n 's/^CMAKE_TOOLCHAIN_FILE:[^=]*=//p' "$cache" | tail -n 1
}

toolchain_file() {
    case "$1" in
        x86_64-linux-gnu) printf '%s\n' "$repo_root/cmake/toolchains/linux-x86_64-gnu.cmake" ;;
        x86_64-linux-musl) printf '%s\n' "$repo_root/cmake/toolchains/linux-x86_64-musl.cmake" ;;
        aarch64-linux-gnu) printf '%s\n' "$repo_root/cmake/toolchains/linux-aarch64-gnu.cmake" ;;
        aarch64-linux-musl) printf '%s\n' "$repo_root/cmake/toolchains/linux-aarch64-musl.cmake" ;;
        armhf-linux-gnu) printf '%s\n' "$repo_root/cmake/toolchains/linux-armhf-gnu.cmake" ;;
        armhf-linux-musl) printf '%s\n' "$repo_root/cmake/toolchains/linux-armhf-musl.cmake" ;;
        arm64-apple-darwin) printf '%s\n' "$repo_root/cmake/toolchains/arm64-apple-darwin.cmake" ;;
        *) die "no CMake toolchain is available for target: $1" ;;
    esac
}

resolver_compiler() {
    local target=$1
    if [[ "$target" == aflpp ]]; then
        "$repo_root/scripts/cpkt-aflpp.sh" discover | sed -n 's/^cc=//p'
        return
    fi
    "$repo_root/scripts/cpkt-toolchains.sh" discover "$target" | sed -n 's/^cc=//p'
}

stale_cache() {
    local cache=$1 target=$2 expected_toolchain=${3:-} expected current current_toolchain
    [[ -f "$cache" ]] || return 1
    current=$(cache_compiler "$cache")
    [[ -n "$current" ]] || return 0
    expected=$(resolver_compiler "$target")
    [[ -n "$expected" && "$current" == "$expected" ]] || return 0
    if [[ -n "$expected_toolchain" ]]; then
        current_toolchain=$(cache_toolchain "$cache")
        [[ "$current_toolchain" == "$expected_toolchain" ]] || return 0
    fi
    return 1
}

preset=
source_dir=
build_dir=
target=
cmake_args=()
direct_toolchain=
while [[ $# -gt 0 ]]; do
    case "$1" in
        --preset)
            [[ $# -ge 2 ]] || die '--preset requires a name'
            preset=$2
            shift 2
            ;;
        --source)
            [[ $# -ge 2 ]] || die '--source requires a directory'
            source_dir=$2
            shift 2
            ;;
        --build)
            [[ $# -ge 2 ]] || die '--build requires a directory'
            build_dir=$2
            shift 2
            ;;
        --target)
            [[ $# -ge 2 ]] || die '--target requires host or a target id'
            target=$2
            shift 2
            ;;
        --)
            shift
            cmake_args=("$@")
            break
            ;;
        *) die "unknown argument: $1" ;;
    esac
done

if [[ -n "$preset" ]]; then
    [[ -z "$source_dir" && -z "$build_dir" && -z "$target" ]] || die '--preset cannot be combined with --source, --build, or --target'
    source_dir=$repo_root
    build_dir=$repo_root/build/$preset
    target=$(preset_target "$preset")
elif [[ -n "$source_dir" || -n "$build_dir" || -n "$target" ]]; then
    [[ -n "$source_dir" && -n "$build_dir" && -n "$target" ]] || die 'direct configuration requires --source, --build, and --target'
    [[ "$target" != host ]] || target=$(host_target)
    source_dir=$repo_root/${source_dir#./}
    build_dir=$repo_root/${build_dir#./}
    if [[ "$target" != none ]]; then
        direct_toolchain=$(toolchain_file "$target")
    fi
else
    die 'use --preset PRESET or --source DIR --build DIR --target TARGET'
fi

if [[ "$target" != none ]]; then
    if [[ "$target" == x86_64-linux-gnu && "$preset" == fuzz ]]; then
        cache_target=aflpp
    else
        cache_target=$target
        "$repo_root/scripts/cpkt-toolchains.sh" ensure "$target"
    fi
    if stale_cache "$build_dir/CMakeCache.txt" "$cache_target" "$direct_toolchain"; then
        printf 'configure_cmake: discarding stale compiler state in %s\n' "$build_dir" >&2
        cmake -E rm -f "$build_dir/CMakeCache.txt"
        cmake -E rm -rf "$build_dir/CMakeFiles"
    fi
fi

if [[ -n "$preset" ]]; then
    exec cmake --preset "$preset" "${cmake_args[@]}"
fi
if [[ -n "$direct_toolchain" ]]; then
    cmake_args+=("-DCMAKE_TOOLCHAIN_FILE=$direct_toolchain")
fi
exec cmake -S "$source_dir" -B "$build_dir" "${cmake_args[@]}"
