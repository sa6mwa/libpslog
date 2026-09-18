#!/usr/bin/env bash
set -euo pipefail
repo_root=${1:?repo root required}
source "$repo_root/bench/host_baseline.sh"
mkdir -p "$repo_root/build"
scratch=$(mktemp -d "$repo_root/build/host-baseline-test.XXXXXX")
trap 'rm -rf "$scratch"' EXIT
fail() { printf '%s\n' "$*" >&2; exit 1; }
expect_failure() {
    local message=$1
    shift
    if "$@" > "$scratch/stdout" 2> "$scratch/stderr"; then fail "unexpected success: $*"; fi
    grep -F "$message" "$scratch/stderr" >/dev/null || fail "missing error: $message"
}
fixture() {
    mkdir -p "$scratch/baselines/$1"
    printf '%s\n' "$2" > "$scratch/baselines/$1/identity"
    printf 'row ns/op=100\n' > "$scratch/baselines/$1/pure-c-baseline.txt"
    printf 'row 100 c_ns/op\n' > "$scratch/baselines/$1/lua-baseline.txt"
}
a=$(printf a | perf_md5)
b=$(printf b | perf_md5)
c=$(printf c | perf_md5)
[[ $(printf abc | perf_md5) == 900150983cd24fb0d6963f7d28e17f72 ]] || fail 'MD5 mismatch'
expect_failure 'no performance baseline' perf_select_baseline "$scratch/baselines" "$a" "$b"
fixture fallback "nodename-md5 $b"
[[ $(perf_select_baseline "$scratch/baselines" "$a" "$b") == "$scratch/baselines/fallback" ]] || fail fallback
[[ $(perf_capture_baseline_dir "$scratch/baselines" "$a" "$b") == "$scratch/baselines/$a" ]] || fail 'capture overwrites fallback'
fixture precise "fingerprint-md5 $a"
[[ $(perf_select_baseline "$scratch/baselines" "$a" "$b") == "$scratch/baselines/precise" ]] || fail precedence
# An existing directory without a matching identity must not be overwritten.
mkdir "$scratch/baselines/$c"
expect_failure 'without matching identity' perf_capture_baseline_dir "$scratch/baselines" "$c" "$b"
# Multiple registered aliases in one baseline.
printf 'fingerprint-md5 %s\nnodename-md5 %s\n' "$c" "$a" >> "$scratch/baselines/precise/identity"
[[ $(perf_select_baseline "$scratch/baselines" "$c" "$b") == "$scratch/baselines/precise" ]] || fail alias
[[ $(perf_capture_baseline_dir "$scratch/baselines" "$c" "$b") == "$scratch/baselines/precise" ]] || fail 'capture ignores precise alias'
fixture duplicate "fingerprint-md5 $a"
expect_failure ambiguous perf_select_baseline "$scratch/baselines" "$a" "$b"
expect_failure ambiguous perf_capture_baseline_dir "$scratch/baselines" "$a" "$b"
printf 'nodename-md5 %s\n' "$b" > "$scratch/baselines/duplicate/identity"
expect_failure ambiguous perf_select_baseline "$scratch/baselines" "$b" "$b"
# Ambiguous fallbacks do not override a unique precise match.
[[ $(perf_select_baseline "$scratch/baselines" "$a" "$b") == "$scratch/baselines/precise" ]] || fail precedence
: > "$scratch/baselines/precise/lua-baseline.txt"
expect_failure incomplete perf_select_baseline "$scratch/baselines" "$a" "$b"
printf 'fingerprint-md5 invalid\n' > "$scratch/baselines/precise/identity"
expect_failure invalid perf_select_baseline "$scratch/baselines" "$a" "$b"
# Exercise real native discovery, plus deterministic macOS discovery and changes.
perf_host_identity
[[ "$PERF_FINGERPRINT_HASH" =~ ^[0-9a-f]{32}$ && "$PERF_NODENAME_HASH" =~ ^[0-9a-f]{32}$ ]] || fail native
(
    uname() { case "$1" in -n) printf 'fixture-node.example.test\n';; -s) printf 'Darwin\n';; -m) printf 'arm64\n';; *) return 1;; esac; }
    fixture_cpu='Apple M2'
    sysctl() { case "$2" in machdep.cpu.brand_string) printf '%s\n' "$fixture_cpu";; hw.logicalcpu) printf '8\n';; *) return 1;; esac; }
    perf_host_identity
    first=$PERF_FINGERPRINT_HASH
    name=$PERF_NODENAME_HASH
    [[ "$name" == "$(printf fixture-node.example.test | perf_md5)" ]] || fail 'node name hash'
    perf_host_identity
    [[ "$first" == "$PERF_FINGERPRINT_HASH" ]] || fail unstable
    fixture_cpu='Apple M3'
    perf_host_identity
    [[ "$first" != "$PERF_FINGERPRINT_HASH" && "$name" == "$PERF_NODENAME_HASH" ]] || fail 'hardware distinction'
    fixture_cpu=
    expect_failure 'complete benchmark host identity' perf_host_identity
)
# Historical aliases are ordinary baseline metadata; keep this fixture local so
# source archives do not depend on export-excluded benchmark history.
historical=5e6e762c7322034ba38e814aa0aaa8fe
historical_root="$scratch/historical-baselines"
mkdir -p "$historical_root/historical"
printf 'nodename-md5 %s\n' "$historical" > "$historical_root/historical/identity"
printf 'row ns/op=100\n' > "$historical_root/historical/pure-c-baseline.txt"
printf 'row 100 c_ns/op\n' > "$historical_root/historical/lua-baseline.txt"
[[ $(perf_select_baseline "$historical_root" "$a" "$historical") == "$historical_root/historical" ]] || fail historical
printf 'row ns/op=100\n' > "$scratch/base"
printf 'row ns/op=149\n' > "$scratch/current"
"$repo_root/bench/check_perf_baseline.sh" "$scratch/base" "$scratch/current" 0.50 ns/op row
printf 'row ns/op=151\n' > "$scratch/current"
expect_failure 'performance gate failed' "$repo_root/bench/check_perf_baseline.sh" "$scratch/base" "$scratch/current" 0.50 ns/op row
expect_failure 'missing baseline' "$repo_root/bench/check_perf_baseline.sh" "$scratch/base" "$scratch/base" 0.50 c_ns/op row
# The real entry point must stop on an unknown host, before compiler/build work.
(
    uname() { case "$1" in -n) printf 'pslog-nonexistent-test-node-7f02a6\n';; -s) command uname -s;; -m) command uname -m;; *) command uname "$@";; esac; }
    export -f uname
    perf_host_identity
    expect_failure 'no performance baseline' bash "$repo_root/bench/run_perf_gate.sh"
    expect_failure 'requires CC' env CC=/nonexistent bash "$repo_root/bench/run_perf_gate.sh" --freeze-baseline
    [[ ! -e "$repo_root/performance-logs/baselines/$PERF_FINGERPRINT_HASH" ]] || fail 'failed capture published a baseline'
)
printf 'Host baseline tests passed.\n'
