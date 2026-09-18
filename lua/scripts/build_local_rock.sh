#!/usr/bin/env bash
set -euo pipefail
repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
rocks=${1:?LuaRocks command required}
rockspec=${2:?rockspec required}
tree=${3:?tree required}
compiler=${4:?compiler required}
stage="$repo_root/build/lua-rock-source"
mkdir -p "$stage/include" "$stage/lua/src" "$stage/lua/pslog"
cp "$repo_root/include/pslog_lua.h" "$stage/include/"
cp "$repo_root/lua/src/pslog_lua.c" "$stage/lua/src/"
cp "$repo_root/lua/pslog/init.lua" "$stage/lua/pslog/"
link_warnings=-Wl,--fatal-warnings
[[ $(uname -s) != Darwin ]] || link_warnings=-Wl,-fatal_warnings
libflag=$("$rocks" config variables.LIBFLAG)
cd "$stage"
"$rocks" make --tree "$tree" "$rockspec" CC="$compiler" LD="$compiler" \
    CFLAGS='-O2 -fPIC -std=c99 -Wall -Wextra -Werror' LIBFLAG="$libflag $link_warnings" \
    LUA_INCDIR="$repo_root/build/lua-host/include" \
    LIBPSLOG_INCDIR="$repo_root/build/lua-sdk/include" \
    LIBPSLOG_LIBDIR="$repo_root/build/lua-sdk/lib"
