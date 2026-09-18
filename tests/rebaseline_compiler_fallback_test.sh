#!/usr/bin/env bash
set -euo pipefail

repo_root=${1:?repo root is required}
cache_file="$repo_root/build/host/CMakeCache.txt"

if [[ ! -f "$cache_file" ]]; then
  printf 'rebaseline compiler test requires configured host cache: %s\n' "$cache_file" >&2
  exit 1
fi

host_cc=$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' "$cache_file" | tail -n 1)
if [[ -z "$host_cc" || ! -x "$host_cc" ]]; then
  printf 'configured host C compiler is unavailable: %s\n' "${host_cc:-<empty>}" >&2
  exit 1
fi

cxx=$(sed -n 's/^CMAKE_CXX_COMPILER:[^=]*=//p' "$cache_file" | tail -n 1)
host_override=$(sed -n 's/^PSLOG_ALLOW_HOST_COMPILER:[^=]*=//p' "$cache_file" | tail -n 1)
if [[ -z "$cxx" && "$host_override" =~ ^(1|ON|TRUE|YES)$ ]]; then
  cxx=$(command -v c++ 2>/dev/null || true)
fi
if [[ -z "$cxx" || ! -x "$cxx" ]]; then
  printf 'SKIP: no C++ compiler available for rebaseline fallback\n'
  exit 0
fi

output=$(env -u CXX PSLOG_REBASELINE_VALIDATE_ONLY=1 "$repo_root/bench/run_rebaseline.sh")
printf '%s\n' "$output" | grep -F "CC=$host_cc" >/dev/null
printf '%s\n' "$output" | grep -F "CXX=$cxx" >/dev/null
make_cxx=$(make -C "$repo_root" -s --eval 'print-host-cxx: ; @printf "%s\\n" "$(HOST_CXX_COMPILER)"' print-host-cxx)
[[ "$make_cxx" == "$cxx" ]] || {
    printf 'Makefile host C++ compiler mismatch: %s\n' "$make_cxx" >&2
    exit 1
}

scratch=$(mktemp -d "$repo_root/build/local-go-compiler-test.XXXXXX")
trap 'rm -rf "$scratch"' EXIT
printf '%s\n' '#!/usr/bin/env bash' 'printf "%s\n" "$CXX" > "$PSLOG_LOCAL_GO_CXX"' > "$scratch/go"
chmod +x "$scratch/go"
PSLOG_LOCAL_GO_CXX="$scratch/cxx" PATH="$scratch:$PATH" "$repo_root/scripts/local-go.sh" version
[[ "$(cat "$scratch/cxx")" == "$cxx" ]] || {
    printf 'local-go host C++ compiler mismatch: %s\n' "$(cat "$scratch/cxx")" >&2
    exit 1
}

# A Darwin Go invocation has no private Linux runtime flags.  It must forward
# no empty argument when the system Bash 3.2 nounset compatibility path runs.
printf '%s\n' '#!/usr/bin/env bash' \
    'case "$1" in -s) printf "Darwin\\n" ;; *) exec /usr/bin/uname "$@" ;; esac' \
    > "$scratch/uname"
printf '%s\n' '#!/usr/bin/env bash' \
    'printf "%s\\n" "$@" > "$PSLOG_LOCAL_GO_ARGS"' \
    > "$scratch/go"
chmod +x "$scratch/uname" "$scratch/go"
PSLOG_LOCAL_GO_ARGS="$scratch/darwin-args" PATH="$scratch:$PATH" "$repo_root/scripts/local-go.sh" version
cmp -s <(printf '%s\n' version) "$scratch/darwin-args" || {
    printf 'local-go did not preserve an empty runtime flag list on Darwin\n' >&2
    exit 1
}
