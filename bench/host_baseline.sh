#!/usr/bin/env bash
# Shared host identity and baseline selection. No plaintext node name is persisted.

perf_md5() {
    if command -v md5sum >/dev/null 2>&1; then
        md5sum | awk '{print $1}'
    elif command -v md5 >/dev/null 2>&1; then
        md5 -q
    else
        printf 'host baselines require md5sum (Linux) or md5 (macOS)\n' >&2
        return 1
    fi
}

perf_host_identity() {
    local node_name os arch cpu count
    node_name=$(uname -n) || return 1
    os=$(uname -s) || return 1
    arch=$(uname -m) || return 1
    case "$os" in
        Linux)
            cpu=$(LC_ALL=C awk -F ': *' '/^(model name|Hardware|Processor)[[:space:]]*:/ {print $2; exit}' /proc/cpuinfo)
            count=$(getconf _NPROCESSORS_ONLN) || return 1
            ;;
        Darwin)
            cpu=$(sysctl -n machdep.cpu.brand_string) || return 1
            count=$(sysctl -n hw.logicalcpu) || return 1
            ;;
        *) printf 'unsupported benchmark host OS: %s\n' "$os" >&2; return 1 ;;
    esac
    if [[ -z "$node_name" || -z "$cpu" || ! "$count" =~ ^[1-9][0-9]*$ ]]; then
        printf 'cannot determine complete benchmark host identity\n' >&2
        return 1
    fi
    PERF_NODENAME_HASH=$(printf '%s' "$node_name" | perf_md5) || return 1
    PERF_FINGERPRINT_HASH=$(printf 'v1\nhost=%s\nos=%s\narch=%s\ncpu=%s\nlogical_cpus=%s\n' \
        "$node_name" "$os" "$arch" "$cpu" "$count" | perf_md5) || return 1
}

perf_select_baseline() {
    local root=$1 fingerprint=$2 node_name_hash=$3 scope=${4:-all} kind wanted identity match
    # Search all precise aliases before considering any node-name alias.
    for kind in fingerprint-md5 nodename-md5; do
        [[ "$scope" != fingerprint-only || "$kind" != nodename-md5 ]] || break
        if [[ "$kind" == fingerprint-md5 ]]; then wanted=$fingerprint; else wanted=$node_name_hash; fi
        match=
        for identity in "$root"/*/identity; do
            [[ -f "$identity" ]] || continue
            if ! awk 'NF && $1 !~ /^#/ {if (NF != 2 || ($1 != "fingerprint-md5" && $1 != "nodename-md5") || length($2) != 32 || $2 ~ /[^0-9a-f]/) exit 1}' "$identity"; then
                printf 'invalid baseline identity: %s\n' "$identity" >&2
                return 1
            fi
            if awk -v kind="$kind" -v hash="$wanted" '$1 == kind && $2 == hash {found=1} END {exit !found}' "$identity"; then
                if [[ -n "$match" ]]; then
                    printf 'ambiguous %s baseline match: %s and %s\n' "$kind" "$match" "${identity%/identity}" >&2
                    return 1
                fi
                match=${identity%/identity}
            fi
        done
        if [[ -n "$match" ]]; then
            if [[ ! -s "$match/pure-c-baseline.txt" || ! -s "$match/lua-baseline.txt" ]]; then
                printf 'incomplete host baseline: %s\n' "$match" >&2
                return 1
            fi
            printf 'Performance baseline: %s (%s match)\n' "${match##*/}" "$kind" >&2
            printf '%s\n' "$match"
            return 0
        fi
    done
    [[ "$scope" != fingerprint-only ]] || return 2
    printf 'no performance baseline for fingerprint-md5 %s / nodename-md5 %s; run make bench-freeze-baseline\n' "$fingerprint" "$node_name_hash" >&2
    return 1
}

perf_capture_baseline_dir() {
    local root=$1 fingerprint=$2 node_name_hash=$3 match status
    if match=$(perf_select_baseline "$root" "$fingerprint" "$node_name_hash" fingerprint-only); then
        printf '%s\n' "$match"
    else
        status=$?
        [[ "$status" == 2 ]] || return "$status"
        if [[ -e "$root/$fingerprint" ]]; then
            printf 'baseline directory exists without matching identity: %s\n' "$root/$fingerprint" >&2
            return 1
        fi
        printf '%s/%s\n' "$root" "$fingerprint"
    fi
}
