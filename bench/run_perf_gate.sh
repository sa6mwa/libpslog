#!/usr/bin/env bash

set -euo pipefail

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf 'missing required command: %s\n' "$1" >&2
        exit 1
    fi
}

extract_ns() {
    awk 'match($0, /[0-9]+(\.[0-9]+)? ns\/op/) { print substr($0, RSTART, RLENGTH-6); exit }'
}

assert_lt() {
    local lhs="$1"
    local rhs="$2"
    local label="$3"

    awk -v lhs="$lhs" -v rhs="$rhs" 'BEGIN { exit !(lhs + 0 < rhs + 0) }' || {
        printf 'performance gate failed: %s (lhs=%s rhs=%s)\n' "$label" "$lhs" "$rhs" >&2
        exit 1
    }
}

run_go_bench() {
    local pattern="$1"
    local out_file="$2"
    local tmpcache

    tmpcache="$(mktemp -d)"
    (
        trap 'rm -rf "$tmpcache"' EXIT
        cd "$repo_root/gobencher"
        GOCACHE="$tmpcache" go test ./benchmark -run '^$' -bench "$pattern" -benchmem -benchtime=200ms -count=1
    ) | tee "$out_file"
}

require_command cmake
require_command ctest
require_command go
require_command awk
require_command mktemp

cd "$repo_root"

cmake --preset host \
  -DPSLOG_BENCHMARK_WITH_LIBLOGGER=OFF \
  -DPSLOG_BENCHMARK_WITH_QUILL=OFF
cmake --build --preset host
ctest --preset debug

printf '\n== gobencher C path smoke tests ==\n'
tmpcache="$(mktemp -d)"
(
    trap 'rm -rf "$tmpcache"' EXIT
    cd "$repo_root/gobencher"
    GOCACHE="$tmpcache" go test ./benchmark -run '^Test(C(LoggerWithPrepared|LoggerPublicWrites|CPublicPreparedParityFixed|CProductionPreparedOutputParity)|CKVFmt(Fixed|Production)OutputParity)$' -count=1
)

prod_out="$(mktemp)"
fixed_out="$(mktemp)"
trap 'rm -f "$prod_out" "$fixed_out"' EXIT

printf '\n== gobencher production compare ==\n'
run_go_bench 'BenchmarkProductionCompare/(jsonGo|jsonC|jsonCkvfmt|jsoncolorGo|jsoncolorC|consoleGo|consoleC|consolecolorGo|consolecolorC)$' "$prod_out"

printf '\n== gobencher fixed compare ==\n'
run_go_bench 'BenchmarkFixedCompare/(jsonGo|jsonC|jsonCkvfmt|jsoncolorGo|jsoncolorC|consoleGo|consoleC|consolecolorGo|consolecolorC)$' "$fixed_out"

prod_json_go="$(grep 'BenchmarkProductionCompare/jsonGo-' "$prod_out" | extract_ns)"
prod_json_c="$(grep 'BenchmarkProductionCompare/jsonC-' "$prod_out" | extract_ns)"
prod_json_ckvfmt="$(grep 'BenchmarkProductionCompare/jsonCkvfmt-' "$prod_out" | extract_ns)"
prod_jsoncolor_go="$(grep 'BenchmarkProductionCompare/jsoncolorGo-' "$prod_out" | extract_ns)"
prod_jsoncolor_c="$(grep 'BenchmarkProductionCompare/jsoncolorC-' "$prod_out" | extract_ns)"
prod_console_go="$(grep 'BenchmarkProductionCompare/consoleGo-' "$prod_out" | extract_ns)"
prod_console_c="$(grep 'BenchmarkProductionCompare/consoleC-' "$prod_out" | extract_ns)"
prod_consolecolor_go="$(grep 'BenchmarkProductionCompare/consolecolorGo-' "$prod_out" | extract_ns)"
prod_consolecolor_c="$(grep 'BenchmarkProductionCompare/consolecolorC-' "$prod_out" | extract_ns)"

fixed_json_go="$(grep 'BenchmarkFixedCompare/jsonGo-' "$fixed_out" | extract_ns)"
fixed_json_c="$(grep 'BenchmarkFixedCompare/jsonC-' "$fixed_out" | extract_ns)"
fixed_json_ckvfmt="$(grep 'BenchmarkFixedCompare/jsonCkvfmt-' "$fixed_out" | extract_ns)"
fixed_jsoncolor_go="$(grep 'BenchmarkFixedCompare/jsoncolorGo-' "$fixed_out" | extract_ns)"
fixed_jsoncolor_c="$(grep 'BenchmarkFixedCompare/jsoncolorC-' "$fixed_out" | extract_ns)"
fixed_console_go="$(grep 'BenchmarkFixedCompare/consoleGo-' "$fixed_out" | extract_ns)"
fixed_console_c="$(grep 'BenchmarkFixedCompare/consoleC-' "$fixed_out" | extract_ns)"
fixed_consolecolor_go="$(grep 'BenchmarkFixedCompare/consolecolorGo-' "$fixed_out" | extract_ns)"
fixed_consolecolor_c="$(grep 'BenchmarkFixedCompare/consolecolorC-' "$fixed_out" | extract_ns)"

assert_lt "$prod_json_c" "$prod_json_go" "production jsonC < jsonGo"
assert_lt "$prod_json_ckvfmt" "$prod_json_go" "production jsonCkvfmt < jsonGo"
assert_lt "$prod_jsoncolor_c" "$prod_jsoncolor_go" "production jsoncolorC < jsoncolorGo"
assert_lt "$prod_console_c" "$prod_console_go" "production consoleC < consoleGo"
assert_lt "$prod_consolecolor_c" "$prod_consolecolor_go" "production consolecolorC < consolecolorGo"

assert_lt "$fixed_json_c" "$fixed_json_go" "fixed jsonC < jsonGo"
assert_lt "$fixed_json_ckvfmt" "$fixed_json_go" "fixed jsonCkvfmt < jsonGo"
assert_lt "$fixed_jsoncolor_c" "$fixed_jsoncolor_go" "fixed jsoncolorC < jsoncolorGo"
assert_lt "$fixed_console_c" "$fixed_console_go" "fixed consoleC < consoleGo"
assert_lt "$fixed_consolecolor_c" "$fixed_consolecolor_go" "fixed consolecolorC < consolecolorGo"

printf '\nPerformance gate passed.\n'
