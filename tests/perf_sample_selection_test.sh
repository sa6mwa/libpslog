#!/usr/bin/env bash
set -euo pipefail

repo_root=${1:?repository root required}
mkdir -p "$repo_root/build"
scratch=$(mktemp -d "$repo_root/build/perf-sample-selection-test.XXXXXX")
trap 'rm -rf "$scratch"' EXIT HUP INT TERM

cat >"$scratch/slow.txt" <<'EOF'
console_prepared iterations=1 ns/op=390.00 bytes/op=72.94
json_prepared iterations=1 ns/op=260.00 bytes/op=127.94
EOF
cat >"$scratch/fast.txt" <<'EOF'
console_prepared iterations=1 ns/op=204.00 bytes/op=72.94
json_prepared iterations=1 ns/op=277.00 bytes/op=127.94
EOF

"$repo_root/bench/select_perf_samples.sh" "$scratch/selected.txt" \
  "$scratch/slow.txt" "$scratch/fast.txt"
expected=$'console_prepared iterations=1 ns/op=204.00 bytes/op=72.94\njson_prepared iterations=1 ns/op=260.00 bytes/op=127.94'
actual=$(cat "$scratch/selected.txt")
if [ "$actual" != "$expected" ]; then
    printf 'sample selection did not retain the fastest value for each metric\n' >&2
    exit 1
fi
