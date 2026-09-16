#!/usr/bin/env bash
set -euo pipefail

repo_root=${1:?repo root is required}
work_root="$repo_root/build/configure-cmake-test"
host_build="$work_root/host"
lua_build="$work_root/lua"
host_platform=$(uname -s)

linux_host_target() {
    case "$(uname -m)" in
        x86_64|amd64) printf '%s\n' x86_64-linux-gnu ;;
        aarch64|arm64) printf '%s\n' aarch64-linux-gnu ;;
        armv7l) printf '%s\n' armhf-linux-gnu ;;
        *) printf 'unsupported Linux host architecture: %s\n' "$(uname -m)" >&2; exit 1 ;;
    esac
}

assert_configured_compiler() {
    local name=$1 compiler=$2 expected
    case "$host_platform" in
        Linux)
            expected=$("$repo_root/scripts/cpkt-toolchains.sh" discover "$(linux_host_target)" |
                sed -n 's/^cc=//p' | tail -n 1)
            [[ -n "$expected" && "$compiler" == "$expected" ]] || {
                printf '%s compiler does not match the resolved Bootlin compiler: %s\n' "$name" "$compiler" >&2
                exit 1
            }
            ;;
        Darwin) [[ -n "$compiler" ]] ;;
        *) printf 'unsupported host platform: %s\n' "$host_platform" >&2; exit 1 ;;
    esac
}

assert_cache_policy() {
    local name=$1 build_dir=$2 output=$3
    case "$host_platform" in
        Linux)
            printf '%s\n' "$output" | grep -F 'discarding stale compiler state' >/dev/null
            [[ ! -e "$build_dir/CMakeFiles/stale-marker" ]]
            ;;
        Darwin)
            if printf '%s\n' "$output" | grep -F 'discarding stale compiler state' >/dev/null; then
                printf 'configure_cmake discarded native macOS cache state for %s\n' "$name" >&2
                exit 1
            fi
            [[ -e "$build_dir/CMakeFiles/stale-marker" ]] || {
                printf 'configure_cmake did not preserve native macOS cache state for %s\n' "$name" >&2
                exit 1
            }
            ;;
    esac
}

rm -rf "$work_root"
trap 'rm -rf "$work_root"' EXIT
mkdir -p "$host_build/CMakeFiles"
printf 'CMAKE_C_COMPILER:FILEPATH=/usr/bin/cc\n' > "$host_build/CMakeCache.txt"
printf 'stale\n' > "$host_build/CMakeFiles/stale-marker"

host_output=$("$repo_root/scripts/configure_cmake.sh" \
    --source . --build "${host_build#$repo_root/}" --target host 2>&1)
assert_cache_policy host "$host_build" "$host_output"
host_cc=$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' "$host_build/CMakeCache.txt" | tail -n 1)
assert_configured_compiler host "$host_cc"
host_recheck=$("$repo_root/scripts/configure_cmake.sh" \
    --source . --build "${host_build#$repo_root/}" --target host 2>&1)
if printf '%s\n' "$host_recheck" | grep -F 'discarding stale compiler state' >/dev/null; then
    printf 'configure_cmake discarded a matching host compiler cache\n' >&2
    exit 1
fi

mkdir -p "$lua_build/CMakeFiles"
printf 'CMAKE_C_COMPILER:FILEPATH=/usr/bin/cc\n' > "$lua_build/CMakeCache.txt"
printf 'stale\n' > "$lua_build/CMakeFiles/stale-marker"
lua_output=$("$repo_root/scripts/configure_cmake.sh" \
    --source cmake/lua --build "${lua_build#$repo_root/}" --target host -- \
    -DCMAKE_BUILD_TYPE=Release 2>&1)
assert_cache_policy lua "$lua_build" "$lua_output"
lua_cc=$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' "$lua_build/CMakeCache.txt" | tail -n 1)
assert_configured_compiler lua "$lua_cc"

fuzz_repo="$work_root/fuzz-repo"
fuzz_build="$fuzz_repo/build/fuzz"
mkdir -p "$fuzz_repo/scripts" "$fuzz_build/CMakeFiles" "$fuzz_repo/fuzz-bin"
cp "$repo_root/scripts/configure_cmake.sh" "$fuzz_repo/scripts/"
printf '%s\n' '#!/usr/bin/env bash' \
    "printf 'cc=%s\\n' '$fuzz_repo/fuzz-bin/cpkt-afl-gcc'" \
    > "$fuzz_repo/scripts/cpkt-aflpp.sh"
printf '%s\n' '#!/usr/bin/env bash' 'exit 1' > "$fuzz_repo/scripts/cpkt-toolchains.sh"
printf '%s\n' '#!/usr/bin/env bash' 'exit 0' > "$fuzz_repo/cmake"
chmod +x "$fuzz_repo/scripts/"*.sh "$fuzz_repo/cmake"
printf 'CMAKE_C_COMPILER:FILEPATH=%s\n' "$fuzz_repo/fuzz-bin/cpkt-afl-gcc" > "$fuzz_build/CMakeCache.txt"
printf 'matching\n' > "$fuzz_build/CMakeFiles/cache-marker"
fuzz_output=$(PATH="$fuzz_repo:$PATH" "$fuzz_repo/scripts/configure_cmake.sh" --preset fuzz 2>&1)
if printf '%s\n' "$fuzz_output" | grep -F 'discarding stale compiler state' >/dev/null ||
   [[ ! -e "$fuzz_build/CMakeFiles/cache-marker" ]]; then
    printf 'configure_cmake discarded a matching AFL++ compiler cache\n' >&2
    exit 1
fi

direct_repo="$work_root/direct-repo"
direct_args="$direct_repo/cmake-arguments"
direct_bin="$direct_repo/bin"
mkdir -p "$direct_repo/scripts" "$direct_repo/cmake/toolchains" "$direct_bin"
cp "$repo_root/scripts/configure_cmake.sh" "$direct_repo/scripts/"
printf '%s\n' '#!/usr/bin/env bash' \
    'case "$1" in ensure) exit 0 ;; discover) printf "cc=/toolchains/aarch64-gcc\\n" ;; esac' \
    > "$direct_repo/scripts/cpkt-toolchains.sh"
printf '%s\n' '#!/usr/bin/env bash' \
    "printf '%s\\n' \"\$@\" > '$direct_args'" \
    > "$direct_bin/cmake"
chmod +x "$direct_repo/scripts/"*.sh "$direct_bin/cmake"
PATH="$direct_bin:$PATH" "$direct_repo/scripts/configure_cmake.sh" \
    --source . --build build/direct-aarch64 --target aarch64-linux-gnu -- \
    -DCMAKE_TOOLCHAIN_FILE=/wrong-toolchain
expected_toolchain="-DCMAKE_TOOLCHAIN_FILE=$direct_repo/cmake/toolchains/linux-aarch64-gnu.cmake"
[[ "$(tail -n 1 "$direct_args")" == "$expected_toolchain" ]] || {
    printf 'direct configuration did not select the requested aarch64 toolchain\n' >&2
    exit 1
}

printf 'configure_cmake stale-cache recovery tests passed.\n'
