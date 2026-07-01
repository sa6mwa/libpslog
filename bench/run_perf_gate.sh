#!/usr/bin/env bash

set -euo pipefail

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf 'missing required command: %s\n' "$1" >&2
        exit 1
    fi
}

run_go_bench() {
    local pattern="$1"
    local out_file="$2"
    local tmpcache

    tmpcache="$(mktemp -d)"
    (
        trap 'rm -rf "$tmpcache"' EXIT
        cd "$repo_root/gobencher"
        GOCACHE="$tmpcache" go test ./benchmark -run '^$' -bench "$pattern" -benchmem -benchtime="$PSLOG_PERF_GO_BENCHTIME" -count=1
    ) | tee "$out_file"
}

run_maybe_pinned() {
    if [ -n "$PSLOG_PERF_CPU" ] && command -v taskset >/dev/null 2>&1; then
        taskset -c "$PSLOG_PERF_CPU" "$@"
    else
        "$@"
    fi
}

require_command cmake
require_command ctest
require_command go
require_command awk
require_command mktemp

PSLOG_PERF_C_ITERS="${PSLOG_PERF_C_ITERS:-200000}"
PSLOG_PERF_C_TOLERANCE="${PSLOG_PERF_C_TOLERANCE:-0.50}"
PSLOG_PERF_LUA_TOLERANCE="${PSLOG_PERF_LUA_TOLERANCE:-0.50}"
PSLOG_PERF_GO_BENCHTIME="${PSLOG_PERF_GO_BENCHTIME:-200ms}"
PSLOG_PERF_CPU="${PSLOG_PERF_CPU:-0}"

c_baseline="$repo_root/performance-logs/pure-c-baseline.txt"
lua_baseline="$repo_root/performance-logs/lua-baseline.txt"
pure_c_out="$(mktemp)"
lua_out="$(mktemp)"
go_compare_out="$(mktemp)"
trap 'rm -f "$pure_c_out" "$lua_out" "$go_compare_out"' EXIT

cd "$repo_root"

cmake --preset host \
  -DPSLOG_BENCHMARK_WITH_LIBLOGGER=OFF \
  -DPSLOG_BENCHMARK_WITH_QUILL=OFF
cmake --build --preset host
ctest --preset host
make lua-rock

printf '\n== pure C regression gate ==\n'
run_maybe_pinned ./build/host/pslog_bench "$PSLOG_PERF_C_ITERS" all | tee "$pure_c_out"
"$repo_root/bench/check_perf_baseline.sh" "$c_baseline" "$pure_c_out" \
  "$PSLOG_PERF_C_TOLERANCE" ns/op \
  console_api \
  console_prepared \
  consolecolor_api \
  consolecolor_prepared \
  json_api \
  json_prepared \
  jsoncolor_api \
  jsoncolor_prepared \
  console_prod_log_fields \
  console_prod_with_log_fields \
  console_prod_log_fields_build \
  console_prod_with_log_fields_build \
  console_prod_level_fields_build \
  console_prod_with_level_fields_build \
  console_prod_levelf_kvfmt \
  console_prod_with_levelf_kvfmt \
  consolecolor_prod_log_fields \
  consolecolor_prod_with_log_fields \
  consolecolor_prod_log_fields_build \
  consolecolor_prod_with_log_fields_build \
  consolecolor_prod_level_fields_build \
  consolecolor_prod_with_level_fields_build \
  consolecolor_prod_levelf_kvfmt \
  consolecolor_prod_with_levelf_kvfmt \
  json_prod_log_fields \
  json_prod_with_log_fields \
  json_prod_log_fields_build \
  json_prod_with_log_fields_build \
  json_prod_level_fields_build \
  json_prod_with_level_fields_build \
  json_prod_levelf_kvfmt \
  json_prod_with_levelf_kvfmt \
  jsoncolor_prod_log_fields \
  jsoncolor_prod_with_log_fields \
  jsoncolor_prod_log_fields_build \
  jsoncolor_prod_with_log_fields_build \
  jsoncolor_prod_level_fields_build \
  jsoncolor_prod_with_level_fields_build \
  jsoncolor_prod_levelf_kvfmt \
  jsoncolor_prod_with_levelf_kvfmt

printf '\n== gobencher C/Lua smoke tests ==\n'
tmpcache="$(mktemp -d)"
(
    trap 'rm -rf "$tmpcache"' EXIT
    cd "$repo_root/gobencher"
    GOCACHE="$tmpcache" go test ./benchmark -run '^Test(C(LoggerWithPrepared|LoggerPublicWrites|CPublicPreparedParityFixed|CProductionPreparedOutputParity)|CKVFmt(Fixed|Production)OutputParity|LuaPreparedBenchmarkBridgeMatchesRawRun|LuaPreparedTableBenchmarkBridgeMatchesRawRun)$' -count=1
)

printf '\n== Lua regression gate ==\n'
run_go_bench 'Benchmark(ProductionCompare|FixedCompare|LuaTableForm)' "$lua_out"
"$repo_root/bench/check_perf_baseline.sh" "$lua_baseline" "$lua_out" \
  "$PSLOG_PERF_LUA_TOLERANCE" c_ns/op \
  BenchmarkProductionCompare/jsonLua \
  BenchmarkFixedCompare/jsonLua \
  BenchmarkLuaTableForm/Production \
  BenchmarkLuaTableForm/Fixed

printf '\n== observational Go-vs-C compare ==\n'
run_go_bench 'Benchmark(Production|Fixed)Compare/(jsonGo|jsonC|jsonCkvfmt|jsoncolorGo|jsoncolorC|consoleGo|consoleC|consolecolorGo|consolecolorC)$' "$go_compare_out"

printf '\nPerformance gate passed.\n'
