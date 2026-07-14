#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "${script_dir}/../.." && pwd)
build_dir="${1:-${repo_root}/build/host}"
version="${2:-0.0.0}"
rock_tree="${3:-${repo_root}/build/luarocks}"
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

mkdir -p "${out_dir}"

lua_cflags=""
lua_libs=""
if command -v pkg-config >/dev/null 2>&1; then
    if pkg-config --exists lua5.5; then
        lua_cflags=$(pkg-config --cflags lua5.5)
        lua_libs=$(pkg-config --libs lua5.5)
    fi
fi

if [ -z "${lua_cflags}" ] && [ -f /usr/local/include/lua.h ]; then
    lua_cflags="-I/usr/local/include"
fi
if [ -z "${lua_libs}" ] && [ -f /usr/local/lib/liblua.a ]; then
    lua_libs="/usr/local/lib/liblua.a -lm -ldl"
fi

if [ -z "${lua_cflags}" ] || [ -z "${lua_libs}" ]; then
    printf 'lua interop embedder test: missing Lua C headers or library\n' >&2
    exit 1
fi

"${cc_bin}" -std=c99 -Wall -Wextra -Werror ${cflags} \
    -D_POSIX_C_SOURCE=200809L \
    -DPSLOG_VERSION_STRING=\"${version}\" \
    -I"${repo_root}/include" \
    -I"${build_dir}/generated/include" \
    ${lua_cflags} \
    "${repo_root}/tests/lua_interop_embedder_test.c" \
    "${repo_root}/lua/src/pslog_lua.c" \
    -L"${build_dir}" -lpslog ${lua_libs} -pthread ${ldflags} \
    -o "${test_bin}"

LD_LIBRARY_PATH="${build_dir}:${LD_LIBRARY_PATH:-}" \
DYLD_LIBRARY_PATH="${build_dir}:${DYLD_LIBRARY_PATH:-}" \
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
    -I"${repo_root}/include" \
    -I"${build_dir}/generated/include" \
    ${lua_cflags} \
    "${installed_consumer_src}" \
    "${installed_core}" -L"${build_dir}" -lpslog ${lua_libs} -pthread ${ldflags} \
    -Wl,-rpath,"${build_dir}" \
    -Wl,-rpath,"$(dirname "${installed_core}")" \
    -o "${installed_consumer_bin}"

LD_LIBRARY_PATH="${build_dir}:$(dirname "${installed_core}"):${LD_LIBRARY_PATH:-}" \
DYLD_LIBRARY_PATH="${build_dir}:$(dirname "${installed_core}"):${DYLD_LIBRARY_PATH:-}" \
    "${installed_consumer_bin}"
