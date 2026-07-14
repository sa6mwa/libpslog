#!/usr/bin/env bash
set -euo pipefail

repo_root=$1
fixture_root=$(mktemp -d)
cleanup() { rm -rf "$fixture_root"; }
trap cleanup EXIT HUP INT TERM

source "$repo_root/scripts/fuzz.sh"
repo_root=$fixture_root
mkdir -p "$repo_root/scripts" "$repo_root/fuzz/corpus" "$repo_root/build/fuzz"
touch "$repo_root/fuzz/corpus/seed"

cmake() { :; }

make_fake_fuzzer() {
  local mode=$1
  printf '%s\n' '#!/usr/bin/env bash' 'set -euo pipefail' \
    'output=' \
    'while [[ $# -gt 0 ]]; do' \
    '  if [[ "$1" == "-o" ]]; then output=$2; shift 2; continue; fi' \
    '  shift' \
    'done' \
    'mkdir -p "$output/default/crashes" "$output/default/hangs"' \
    "${mode}" >"$repo_root/fake-afl-fuzz"
  chmod +x "$repo_root/fake-afl-fuzz"
}

printf '%s\n' '#!/usr/bin/env bash' 'printf "afl_fuzz=%s\\n" "$(dirname "$0")/../fake-afl-fuzz"' \
  >"$repo_root/scripts/cpkt-aflpp.sh"
chmod +x "$repo_root/scripts/cpkt-aflpp.sh"

make_fake_fuzzer 'touch "$output/default/crashes/id:000000,sig:06"'
if main smoke 1; then
  printf 'fuzz gate accepted a recorded crash\n' >&2
  exit 1
fi

make_fake_fuzzer ':'
main smoke 1
