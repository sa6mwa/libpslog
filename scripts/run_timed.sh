#!/usr/bin/env bash
# Run one lifecycle phase serially and report its wall-clock duration.
set -euo pipefail

usage() {
    printf 'usage: %s phase command [argument ...]\n' "$0" >&2
    exit 2
}

[[ $# -ge 2 ]] || usage
phase=$1
shift
if [[ "$phase" =~ [[:space:]] || "$phase" == *$'\t'* || "$phase" == *$'\n'* ]]; then
    printf 'run_timed: phase must not contain whitespace: %s\n' "$phase" >&2
    exit 2
fi

started_epoch=$(date +%s)
printf 'PKT_TIMING_BEGIN phase=%s started_epoch=%s\n' "$phase" "$started_epoch"

set +e
"$@"
status=$?
set -e

finished_epoch=$(date +%s)
elapsed_seconds=$((finished_epoch - started_epoch))
printf 'PKT_TIMING_END phase=%s status=%s elapsed_seconds=%s\n' \
    "$phase" "$status" "$elapsed_seconds"

if [[ -n "${PKT_TIMING_FILE:-}" ]]; then
    timing_dir=$(dirname -- "$PKT_TIMING_FILE")
    mkdir -p "$timing_dir"
    if [[ ! -e "$PKT_TIMING_FILE" ]]; then
        printf 'phase\tstatus\tstarted_epoch\tfinished_epoch\telapsed_seconds\n' >"$PKT_TIMING_FILE"
    fi
    printf '%s\t%s\t%s\t%s\t%s\n' \
        "$phase" "$status" "$started_epoch" "$finished_epoch" "$elapsed_seconds" >>"$PKT_TIMING_FILE"
fi

exit "$status"
