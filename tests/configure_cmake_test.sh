#!/usr/bin/env bash
set -euo pipefail

repo_root=${1:?repo root is required}
work_root="$repo_root/build/configure-cmake-test"
host_build="$work_root/host"
lua_build="$work_root/lua"

rm -rf "$work_root"
trap 'rm -rf "$work_root"' EXIT
mkdir -p "$host_build/CMakeFiles"
printf 'CMAKE_C_COMPILER:FILEPATH=/usr/bin/cc\n' > "$host_build/CMakeCache.txt"
printf 'stale\n' > "$host_build/CMakeFiles/stale-marker"

host_output=$("$repo_root/scripts/configure_cmake.sh" \
    --source . --build "${host_build#$repo_root/}" --target host 2>&1)
printf '%s\n' "$host_output" | grep -F 'discarding stale compiler state' >/dev/null
[[ ! -e "$host_build/CMakeFiles/stale-marker" ]]
host_cc=$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' "$host_build/CMakeCache.txt" | tail -n 1)
case "$(uname -s)" in
    Linux) [[ "$host_cc" == *'/c.pkt.systems/toolchains/roots/'* ]] ;;
    Darwin) [[ -n "$host_cc" ]] ;;
    *) exit 1 ;;
esac
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
printf '%s\n' "$lua_output" | grep -F 'discarding stale compiler state' >/dev/null
[[ ! -e "$lua_build/CMakeFiles/stale-marker" ]]
lua_cc=$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' "$lua_build/CMakeCache.txt" | tail -n 1)
case "$(uname -s)" in
    Linux) [[ "$lua_cc" == *'/c.pkt.systems/toolchains/roots/'* ]] ;;
    Darwin) [[ -n "$lua_cc" ]] ;;
    *) exit 1 ;;
esac

printf 'configure_cmake stale-cache recovery tests passed.\n'
