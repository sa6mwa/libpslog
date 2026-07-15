#!/usr/bin/env bash
# Run a bounded native AFL++ job against the Bootlin-built fuzz target.
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)

main() {
  local mode=${1:-} seconds=${2:-} description toolchain_description afl_fuzz corpus output failure_input sysroot loader library_path

  case "$mode" in smoke|run|long) ;; *) printf 'usage: %s {smoke|run|long} seconds\n' "$0" >&2; return 2;; esac
  case "$seconds" in ''|*[!0-9]*) printf 'fuzz duration must be a positive integer number of seconds\n' >&2; return 2;; esac
  [[ "$seconds" -gt 0 ]] || { printf 'fuzz duration must be positive\n' >&2; return 2; }

  cd "$repo_root"
  cmake --preset fuzz
  cmake --build --preset fuzz
  description=$(./scripts/cpkt-aflpp.sh discover)
  afl_fuzz=$(sed -n 's/^afl_fuzz=//p' <<<"$description" | tail -n 1)
  [[ -x "$afl_fuzz" ]] || { printf 'pinned AFL++ resolver did not provide afl-fuzz\n' >&2; return 1; }
  toolchain_description=$(./scripts/cpkt-toolchains.sh discover x86_64-linux-gnu)
  sysroot=$(sed -n 's/^sysroot=//p' <<<"$toolchain_description" | tail -n 1)
  [[ -d "$sysroot" ]] || { printf 'pinned Bootlin resolver did not provide the native fuzz sysroot\n' >&2; return 1; }
  loader=$("$repo_root/scripts/run_sysroot_binary.sh" --loader --sysroot "$sysroot")
  library_path="$sysroot/lib:$sysroot/usr/lib:$sysroot/lib64:$sysroot/usr/lib64"

  corpus="$repo_root/fuzz/corpus"
  output="$repo_root/build/fuzz/afl-$mode"
  [[ -d "$corpus" ]] || { printf 'fuzz corpus missing: %s\n' "$corpus" >&2; return 1; }
  case "$output" in "$repo_root"/build/fuzz/*) ;; *) printf 'refusing unsafe fuzz output path\n' >&2; return 1;; esac
  rm -rf "$output"
  # AFL++ inspects the ELF loader first; the instrumented target starts after
  # that loader has entered the selected Bootlin sysroot.
  AFL_NO_UI=1 AFL_SKIP_CPUFREQ=1 AFL_I_DONT_CARE_ABOUT_MISSING_CRASHES=1 AFL_SKIP_BIN_CHECK=1 \
    "$afl_fuzz" -V "$seconds" -i "$corpus" -o "$output" -- \
    "$loader" --library-path "$library_path" "$repo_root/build/fuzz/pslog_fuzz"

  failure_input=$(find "$output" -type f \( -path '*/crashes/id:*' -o -path '*/hangs/id:*' \) -print -quit)
  if [[ -n "$failure_input" ]]; then
    printf 'AFL++ recorded a crashing or hanging input: %s\n' "$failure_input" >&2
    return 1
  fi
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  main "$@"
fi
