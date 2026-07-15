#!/usr/bin/env bash
# Run a native Bootlin-built ELF through the dynamic loader in its selected sysroot.
set -euo pipefail

usage() {
  printf 'usage: %s [--loader] [--library-path DIR]... --sysroot DIR PROGRAM [argument ...]\n' "$0" >&2
  printf '       %s [--loader] [--library-path DIR]... --build-dir DIR PROGRAM [argument ...]\n' "$0" >&2
  exit 2
}

sysroot=
print_loader=false
library_paths=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --loader) print_loader=true; shift ;;
    --library-path)
      [[ $# -ge 2 ]] || usage
      [[ -d "$2" ]] || { printf 'run_sysroot_binary: library path is missing or not a directory: %s\n' "$2" >&2; exit 1; }
      library_paths+=("$(CDPATH= cd -- "$2" && pwd -P)")
      shift 2
      ;;
    --sysroot)
      [[ $# -ge 2 ]] || usage
      [[ -z "$sysroot" ]] || usage
      sysroot=$2
      shift 2
      ;;
    --)
      shift
      break
      ;;
    --build-dir)
      [[ $# -ge 2 ]] || usage
      [[ -z "$sysroot" ]] || usage
      cache_file=$2/CMakeCache.txt
      [[ -f "$cache_file" ]] || { printf 'run_sysroot_binary: missing CMake cache: %s\n' "$cache_file" >&2; exit 1; }
      sysroot=$(sed -n 's/^CMAKE_SYSROOT:[^=]*=//p' "$cache_file" | tail -n 1)
      [[ -n "$sysroot" ]] || { printf 'run_sysroot_binary: CMAKE_SYSROOT is missing from %s\n' "$cache_file" >&2; exit 1; }
      shift 2
      ;;
    *) break ;;
  esac
done

[[ -n "$sysroot" ]] || usage

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
    for extra_library_path in "${library_paths[@]}"; do
      library_path="$extra_library_path:$library_path"
    done
    exec "$loader" --library-path "$library_path" "$@"
  fi
done

printf 'run_sysroot_binary: no supported native x86_64 loader under %s\n' "$sysroot" >&2
exit 1
