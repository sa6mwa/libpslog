#!/usr/bin/env bash
set -euo pipefail

output=${1:?output path required}
shift
if [ "$#" -eq 0 ]; then
    printf 'usage: %s OUTPUT SAMPLE...\n' "$0" >&2
    exit 2
fi

awk '
{
    name = $1
    value = ""
    for (field = 1; field <= NF; field++) {
        if ($field ~ /^ns\/op=/) {
            value = $field
            sub(/^ns\/op=/, "", value)
            break
        }
    }
    if (name == "" || value == "") {
        printf "invalid benchmark row in %s: %s\n", FILENAME, $0 > "/dev/stderr"
        exit 2
    }
    if (!(name in seen)) {
        seen[name] = ++count
        order[count] = name
    }
    if (!(name in best) || value + 0 < best[name]) {
        best[name] = value + 0
        row[name] = $0
    }
}
END {
    for (ordinal = 1; ordinal <= count; ordinal++) {
        print row[order[ordinal]]
    }
}
' "$@" >"$output"
