#!/usr/bin/env bash

set -euo pipefail

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
source "$repo_root/bench/host_baseline.sh"
freeze=0
case "${1:-}" in
    "") ;;
    --freeze-baseline) freeze=1 ;;
    *) printf 'usage: %s [--freeze-baseline]\n' "$0" >&2; exit 2 ;;
esac
[[ $# -le 1 ]] || exit 2
perf_host_identity
baseline_root="$repo_root/performance-logs/baselines"
if [[ "$freeze" == 0 ]]; then
    baseline_dir=$(perf_select_baseline "$baseline_root" "$PERF_FINGERPRINT_HASH" "$PERF_NODENAME_HASH")
else
    baseline_dir=$(perf_capture_baseline_dir "$baseline_root" "$PERF_FINGERPRINT_HASH" "$PERF_NODENAME_HASH")
fi
mkdir -p "$repo_root/build"
scratch=$(mktemp -d "$repo_root/build/perf-gate.XXXXXX")
trap 'rm -rf "$scratch"' EXIT

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

    tmpcache="$(mktemp -d "$scratch/go-cache.XXXXXX")"
    (
        trap 'rm -rf "$tmpcache"' EXIT
        cd "$repo_root/gobencher"
        GOCACHE="$tmpcache" "$repo_root/scripts/local-go.sh" test ./benchmark -run '^$' -bench "$pattern" -benchmem -benchtime="$PSLOG_PERF_GO_BENCHTIME" -count=1
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

if [[ -z "${CC:-}" || ! -x "${CC}" ]]; then
    printf 'perf gate requires CC to name the configured C compiler\n' >&2
    exit 1
fi
if [[ -z "${CXX:-}" || ! -x "${CXX}" ]]; then
    printf 'perf gate requires CXX to name the configured C++ compiler\n' >&2
    exit 1
fi
export CC CXX

PSLOG_PERF_C_ITERS="${PSLOG_PERF_C_ITERS:-200000}"
PSLOG_PERF_C_TOLERANCE="${PSLOG_PERF_C_TOLERANCE:-0.50}"
PSLOG_PERF_LUA_TOLERANCE="${PSLOG_PERF_LUA_TOLERANCE:-0.50}"
PSLOG_PERF_GO_BENCHTIME="${PSLOG_PERF_GO_BENCHTIME:-200ms}"
PSLOG_PERF_CPU="${PSLOG_PERF_CPU-0}"

pure_c_out="$scratch/pure-c.txt"
lua_out="$scratch/lua.txt"
go_compare_out="$scratch/go-compare.txt"
c_baseline="$baseline_dir/pure-c-baseline.txt"
lua_baseline="$baseline_dir/lua-baseline.txt"
if [[ "$freeze" == 1 ]]; then
    # Self-comparison still validates that every required metric was captured.
    c_baseline="$pure_c_out"
    lua_baseline="$lua_out"
fi

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
tmpcache="$(mktemp -d "$scratch/go-cache.XXXXXX")"
(
    trap 'rm -rf "$tmpcache"' EXIT
    cd "$repo_root/gobencher"
    GOCACHE="$tmpcache" "$repo_root/scripts/local-go.sh" test ./benchmark -run '^Test(C(LoggerWithPrepared|LoggerPublicWrites|CPublicPreparedParityFixed|CProductionPreparedOutputParity)|CKVFmt(Fixed|Production)OutputParity|LuaPreparedBenchmarkBridgeMatchesRawRun|LuaPreparedTableBenchmarkBridgeMatchesRawRun)$' -count=1
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

if [[ "$freeze" == 1 ]]; then
    mkdir -p "$baseline_dir"
    for pair in "pure-c:$pure_c_out" "lua:$lua_out"; do
        name=${pair%%:*}
        input=${pair#*:}
        {
            printf '# pslog_perf_artifact: %s-baseline\n' "$name"
            printf '# host: fingerprint-md5 %s\n' "$PERF_FINGERPRINT_HASH"
            printf '# measured_commit: %s\n' "$(git rev-parse HEAD)"
            printf '# worktree: %s\n' "$(if [[ -n $(git status --porcelain --untracked-files=no) ]]; then printf modified; else printf clean; fi)"
            printf '# date: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
            printf '# compiler: %s\n' "$("$CC" --version | head -n 1)"
            printf '# cmake_preset: host; liblogger off; quill off\n'
            printf '# run_count: 1; summary: raw benchmark rows\n'
            printf '# C: iterations=%s; requested_cpu=%s; taskset=%s\n' "$PSLOG_PERF_C_ITERS" "$PSLOG_PERF_CPU" "$(if command -v taskset >/dev/null; then printf available; else printf unavailable; fi)"
            printf '# Lua/Go: benchtime=%s; count=1; pinning=none\n' "$PSLOG_PERF_GO_BENCHTIME"
            cat "$input"
        } > "$scratch/$name-baseline.txt"
        mv "$scratch/$name-baseline.txt" "$baseline_dir/$name-baseline.txt"
    done
    # Preserve manually registered aliases when refreshing this host.
    if [[ ! -f "$baseline_dir/identity" ]]; then
        printf 'fingerprint-md5 %s\n' "$PERF_FINGERPRINT_HASH" > "$baseline_dir/identity"
    fi
    printf '\nBaseline frozen: %s\n' "${baseline_dir##*/}"
else
    printf '\nPerformance gate passed.\n'
fi
