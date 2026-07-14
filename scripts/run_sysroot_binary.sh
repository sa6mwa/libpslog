#!/usr/bin/env bash
# Run a native Bootlin-built ELF through the dynamic loader in its selected sysroot.
set -euo pipefail

usage() {
  printf 'usage: %s [--loader] --sysroot DIR PROGRAM [argument ...]\n' "$0" >&2
  printf '       %s [--loader] --build-dir DIR PROGRAM [argument ...]\n' "$0" >&2
  exit 2
}

sysroot=
print_loader=false
if [[ "${1:-}" == "--loader" ]]; then
  print_loader=true
  shift
fi
case "${1:-}" in
  --sysroot) [[ $# -ge 2 ]] || usage; sysroot=$2; shift 2 ;;
  --build-dir)
    [[ $# -ge 2 ]] || usage
    cache_file=$2/CMakeCache.txt
    [[ -f "$cache_file" ]] || { printf 'run_sysroot_binary: missing CMake cache: %s\n' "$cache_file" >&2; exit 1; }
    sysroot=$(sed -n 's/^CMAKE_SYSROOT:[^=]*=//p' "$cache_file" | tail -n 1)
    [[ -n "$sysroot" ]] || { printf 'run_sysroot_binary: CMAKE_SYSROOT is missing from %s\n' "$cache_file" >&2; exit 1; }
    shift 2
    ;;
  *) usage ;;
esac

if [[ "$print_loader" == false ]]; then
  [[ $# -ge 1 ]] || usage
fi
sysroot=$(CDPATH= cd -- "$sysroot" && pwd -P)

for loader in \
  "$sysroot/lib/ld-linux-x86-64.so.2" \
  "$sysroot/lib/ld-musl-x86_64.so.1"; do
  if [[ -x "$loader" ]]; then
    if [[ "$print_loader" == true ]]; then
      printf '%s\n' "$loader"
      exit 0
    fi
    library_path="$sysroot/lib:$sysroot/usr/lib:$sysroot/lib64:$sysroot/usr/lib64"
    exec "$loader" --library-path "$library_path" "$@"
  fi
done

printf 'run_sysroot_binary: no supported native x86_64 loader under %s\n' "$sysroot" >&2
exit 1
