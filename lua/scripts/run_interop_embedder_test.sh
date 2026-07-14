#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "${script_dir}/../.." && pwd)
build_dir="${1:-${repo_root}/build/host}"
version="${2:-0.0.0}"
rock_tree="${3:-${repo_root}/build/luarocks}"
sdk_prefix="${4:-${repo_root}/build/lua-sdk}"
out_dir="${repo_root}/build/lua-interop"
test_bin="${out_dir}/pslog_lua_interop_tests"
installed_consumer_src="${out_dir}/installed_rock_consumer.c"
installed_consumer_bin="${out_dir}/installed_rock_consumer"
if [ -n "${CC:-}" ]; then
    cc_bin="${CC}"
else
    cache_file="${build_dir}/CMakeCache.txt"
    if [ ! -f "${cache_file}" ]; then
        printf 'lua interop embedder test: missing configured Bootlin compiler cache: %s\n' "${cache_file}" >&2
        exit 1
    fi
    cc_bin=$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' "${cache_file}" | tail -n 1)
fi
if [ -z "${cc_bin}" ] || [ ! -x "${cc_bin}" ]; then
    printf 'lua interop embedder test: configured C compiler is unavailable: %s\n' "${cc_bin:-<empty>}" >&2
    exit 1
fi
cflags="${CFLAGS:-}"
ldflags="${LDFLAGS:-}"
lua_stage_dir="${repo_root}/build/lua-host"
lua_include_dir="${lua_stage_dir}/include"
lua_lib_dir="${lua_stage_dir}/lib"
sdk_include_dir="${sdk_prefix}/include"
sdk_lib_dir="${sdk_prefix}/lib"

mkdir -p "${out_dir}"

if [ ! -f "${lua_include_dir}/lua.h" ] || [ ! -f "${lua_lib_dir}/liblua.a" ]; then
    printf 'lua interop embedder test: missing staged Lua 5.5 C development files; run make lua-rock first\n' >&2
    exit 1
fi
if [ ! -f "${sdk_include_dir}/pslog.h" ] || [ ! -e "${sdk_lib_dir}/libpslog.so" ]; then
    printf 'lua interop embedder test: missing installed local libpslog SDK; run make lua-rock first\n' >&2
    exit 1
fi
lua_cflags="-I${lua_include_dir}"
lua_libs="-L${lua_lib_dir} -llua -lm -ldl"

"${cc_bin}" -std=c99 -Wall -Wextra -Werror ${cflags} \
    -D_POSIX_C_SOURCE=200809L \
    -DPSLOG_VERSION_STRING=\"${version}\" \
    -I"${sdk_include_dir}" \
    -I"${repo_root}/include" \
    ${lua_cflags} \
    "${repo_root}/tests/lua_interop_embedder_test.c" \
    "${repo_root}/lua/src/pslog_lua.c" \
    -L"${sdk_lib_dir}" -lpslog ${lua_libs} -pthread ${ldflags} \
    -o "${test_bin}"

LD_LIBRARY_PATH="${sdk_lib_dir}:${LD_LIBRARY_PATH:-}" \
DYLD_LIBRARY_PATH="${sdk_lib_dir}:${DYLD_LIBRARY_PATH:-}" \
    "${test_bin}"

installed_header=$(find "${rock_tree}/share/lua" -name pslog_lua.h -type f | head -n 1)
if [ -z "${installed_header}" ]; then
    printf 'lua interop embedder test: installed rock is missing pslog_lua.h\n' >&2
    exit 1
fi
installed_header_dir=$(dirname "${installed_header}")
installed_core=$(find "${rock_tree}/lib/lua" -path '*/pslog/core.so' -type f | head -n 1)
if [ -z "${installed_core}" ]; then
    printf 'lua interop embedder test: installed rock is missing pslog/core.so\n' >&2
    exit 1
fi

installed_consumer_src="${repo_root}/tests/lua_interop_installed_consumer.c"

"${cc_bin}" -std=c99 -Wall -Wextra -Werror ${cflags} \
    -I"${installed_header_dir}" \
    -I"${sdk_include_dir}" \
    -I"${repo_root}/include" \
    ${lua_cflags} \
    "${installed_consumer_src}" \
    "${installed_core}" -L"${sdk_lib_dir}" -lpslog ${lua_libs} -pthread ${ldflags} \
    -Wl,-rpath,"${sdk_lib_dir}" \
    -Wl,-rpath,"$(dirname "${installed_core}")" \
    -o "${installed_consumer_bin}"

LD_LIBRARY_PATH="${sdk_lib_dir}:$(dirname "${installed_core}"):${LD_LIBRARY_PATH:-}" \
DYLD_LIBRARY_PATH="${sdk_lib_dir}:$(dirname "${installed_core}"):${DYLD_LIBRARY_PATH:-}" \
    "${installed_consumer_bin}"
