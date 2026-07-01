#!/usr/bin/env bash

set -euo pipefail

usage() {
    cat >&2 <<'EOF'
usage: check_perf_baseline.sh BASELINE CURRENT TOLERANCE METRIC LABEL...

Compares benchmark rows in CURRENT against BASELINE and fails if any CURRENT
metric is greater than BASELINE * (1 + TOLERANCE).

Rows are addressed by their first field. Go benchmark worker suffixes such as
BenchmarkFixedCompare/jsonLua-12 are normalized to BenchmarkFixedCompare/jsonLua.
Supported metric forms are "123 ns/op", "123 c_ns/op", and "ns/op=123".
EOF
    exit 2
}

if [ "$#" -lt 5 ]; then
    usage
fi

baseline_file="$1"
current_file="$2"
tolerance="$3"
metric="$4"
shift 4

if [ ! -f "$baseline_file" ]; then
    printf 'performance gate failed: missing baseline file: %s\n' "$baseline_file" >&2
    exit 1
fi
if [ ! -f "$current_file" ]; then
    printf 'performance gate failed: missing current benchmark file: %s\n' "$current_file" >&2
    exit 1
fi

extract_metric() {
    local file="$1"
    local label="$2"

    awk -v want_label="$label" -v want_metric="$metric" '
        function normalize_label(label) {
            sub(/-[0-9]+$/, "", label)
            return label
        }
        function numeric(value) {
            return value ~ /^[0-9]+(\.[0-9]+)?$/
        }
        function emit_metric(value) {
            if (numeric(value)) {
                print value
                found = 1
                exit
            }
        }
        /^[[:space:]]*#/ || NF == 0 {
            next
        }
        normalize_label($1) == want_label {
            for (i = 2; i <= NF; ++i) {
                if ($i == want_metric && i > 1) {
                    emit_metric($(i - 1))
                }
                if (index($i, want_metric "=") == 1) {
                    split($i, parts, "=")
                    emit_metric(parts[2])
                }
            }
        }
        END {
            if (!found) {
                exit 1
            }
        }
    ' "$file"
}

failed=0
for label in "$@"; do
    baseline_value="$(extract_metric "$baseline_file" "$label" || true)"
    current_value="$(extract_metric "$current_file" "$label" || true)"

    if [ -z "$baseline_value" ]; then
        printf 'performance gate failed: missing baseline %s metric for %s in %s\n' \
            "$metric" "$label" "$baseline_file" >&2
        failed=1
        continue
    fi
    if [ -z "$current_value" ]; then
        printf 'performance gate failed: missing current %s metric for %s in %s\n' \
            "$metric" "$label" "$current_file" >&2
        failed=1
        continue
    fi

    awk -v label="$label" \
        -v baseline="$baseline_value" \
        -v current="$current_value" \
        -v tolerance="$tolerance" \
        -v metric="$metric" '
        BEGIN {
            limit = baseline * (1.0 + tolerance)
            if (current <= limit) {
                printf "perf ok: %s %s current=%s baseline=%s limit=%.2f tolerance=%.0f%%\n",
                    label, metric, current, baseline, limit, tolerance * 100.0
                exit 0
            }
            printf "performance gate failed: %s %s current=%s baseline=%s limit=%.2f tolerance=%.0f%%\n",
                label, metric, current, baseline, limit, tolerance * 100.0 > "/dev/stderr"
            exit 1
        }
    ' || failed=1
done

exit "$failed"
