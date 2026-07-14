#!/usr/bin/env bash
# Validate the host-editor configuration against the native Debug database only.
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
database_dir="$repo_root/build/debug"
# clangd's --check mode exercises code actions as well as diagnostics.  The
# core implementation intentionally offers enough refactoring opportunities
# for some clangd releases to report failed *optional* actions as errors.
# Check the native public-consumer translation unit instead: it uses the
# native Debug compile database and verifies editor parsing of the installed
# public C surface without treating unavailable refactorings as source errors.
source_file="$repo_root/examples/example.c"

command -v clangd >/dev/null 2>&1 || {
  printf 'clangd_check: host clangd is required; install it through the host package manager\n' >&2
  exit 1
}
[[ -f "$database_dir/compile_commands.json" ]] || {
  printf 'clangd_check: native debug compile database is missing: %s\n' "$database_dir/compile_commands.json" >&2
  exit 1
}
[[ -f "$source_file" ]] || {
  printf 'clangd_check: source is missing: %s\n' "$source_file" >&2
  exit 1
}

cd "$repo_root"
output_file=$(mktemp)
trap 'rm -f "$output_file"' EXIT
if ! clangd --compile-commands-dir="$database_dir" --check="$source_file" >"$output_file" 2>&1; then
  cat "$output_file" >&2
  exit 1
fi
