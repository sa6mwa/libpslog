#!/usr/bin/env bash
# Report target inspection and mutation tools from the CMake build that produced
# an artifact.  Configured state wins; ambient PATH is only the final fallback.
set -euo pipefail

usage() {
  printf 'usage: %s --build-dir DIR --target-id TARGET_ID [--cc PATH] [--strip PATH] [--readelf PATH] [--install-name-tool PATH] [--otool PATH] [--format shell|cmake]\n' "$0" >&2
  exit 2
}

build_dir= target_id= output_format=shell
override_cc=${PSLOG_TARGET_CC:-}
override_strip=${PSLOG_TARGET_STRIP:-}
override_readelf=${PSLOG_TARGET_READELF:-}
override_install_name_tool=${PSLOG_TARGET_INSTALL_NAME_TOOL:-}
override_otool=${PSLOG_TARGET_OTOOL:-}
while [[ $# -gt 0 ]]; do
  case "$1" in
    --build-dir) [[ $# -ge 2 ]] || usage; build_dir=$2; shift 2 ;;
    --target-id) [[ $# -ge 2 ]] || usage; target_id=$2; shift 2 ;;
    --cc) [[ $# -ge 2 ]] || usage; override_cc=$2; shift 2 ;;
    --strip) [[ $# -ge 2 ]] || usage; override_strip=$2; shift 2 ;;
    --readelf) [[ $# -ge 2 ]] || usage; override_readelf=$2; shift 2 ;;
    --install-name-tool) [[ $# -ge 2 ]] || usage; override_install_name_tool=$2; shift 2 ;;
    --otool) [[ $# -ge 2 ]] || usage; override_otool=$2; shift 2 ;;
    --format) [[ $# -ge 2 ]] || usage; output_format=$2; shift 2 ;;
    *) usage ;;
  esac
done
[[ -n "$build_dir" && -n "$target_id" ]] || usage
[[ "$output_format" == shell || "$output_format" == cmake ]] || usage

cache_file="$build_dir/CMakeCache.txt"
[[ -f "$cache_file" ]] || { printf 'discover_target_tools: missing CMake cache: %s\n' "$cache_file" >&2; exit 1; }

cache_value() {
  sed -n "s|^$1:[^=]*=||p" "$cache_file" | tail -n 1
}

require_executable() {
  local name=$1 value=$2
  [[ -n "$value" && "$value" != *-NOTFOUND && -x "$value" ]] || {
    printf 'discover_target_tools: configured %s is unavailable for %s: %s\n' "$name" "$target_id" "${value:-<empty>}" >&2
    exit 1
  }
}

cc=${override_cc:-$(cache_value CMAKE_C_COMPILER)}
require_executable CMAKE_C_COMPILER "$cc"
compiler_dir=$(CDPATH= cd -- "$(dirname -- "$cc")" && pwd)
compiler_name=$(basename -- "$cc")
host_prefix=$(cache_value PSLOG_OSXCROSS_HOST)
if [[ -z "$host_prefix" || "$host_prefix" == *-NOTFOUND ]]; then
  host_prefix=${CPKT_OSXCROSS_HOST:-}
fi
if [[ -z "$host_prefix" ]]; then
  host_prefix=${compiler_name%-clang}
  host_prefix=${host_prefix%-gcc}
  host_prefix=${host_prefix%-cc}
fi

is_known_host_tool() {
  local tool=$1
  [[ "$target_id" != x86_64-linux-gnu && "$tool" =~ ^/usr/bin/(strip|readelf|otool|install_name_tool)$ ]]
}

lookup_tool() {
  local label=$1 cache_key=$2 override=$3 tool_name=$4 required=$5 candidate= source=
  if [[ -n "$override" ]]; then
    candidate=$override; source=override
  fi
  if [[ -z "$candidate" ]]; then
    candidate=$(cache_value "$cache_key")
    [[ "$candidate" == *-NOTFOUND ]] && candidate=
    [[ -n "$candidate" ]] && source=cmake-cache
  fi
  if [[ -z "$candidate" ]]; then
    for candidate in "$compiler_dir/$host_prefix-$tool_name" "$compiler_dir/$target_id-$tool_name"; do
      [[ -x "$candidate" ]] && { source=target-prefixed-sibling; break; }
      candidate=
    done
  fi
  if [[ -z "$candidate" && -x "$compiler_dir/$tool_name" ]]; then
    candidate="$compiler_dir/$tool_name"; source=unprefixed-sibling
  fi
  if [[ -z "$candidate" ]]; then
    for candidate in "${host_prefix}-${tool_name}" "${target_id}-${tool_name}"; do
      candidate=$(command -v "$candidate" 2>/dev/null || true)
      [[ -n "$candidate" ]] && { source=target-prefixed-path; break; }
    done
  fi
  if [[ -z "$candidate" ]]; then
    candidate=$(command -v "$tool_name" 2>/dev/null || true)
    [[ -n "$candidate" ]] && source=path
  fi
  if [[ -n "$candidate" && ! -x "$candidate" ]]; then
    printf 'discover_target_tools: %s from %s is not executable for %s: %s\n' "$label" "$source" "$target_id" "$candidate" >&2
    exit 1
  fi
  if [[ -n "$candidate" ]] && is_known_host_tool "$candidate"; then
    printf 'discover_target_tools: refusing known host %s for cross target %s: %s (checked override, CMake cache, compiler siblings, and PATH)\n' "$label" "$target_id" "$candidate" >&2
    exit 1
  fi
  if [[ -z "$candidate" && "$required" == required ]]; then
    printf 'discover_target_tools: required %s unavailable for %s; checked override, CMake cache, %s siblings, and PATH\n' "$label" "$target_id" "$compiler_dir" >&2
    exit 1
  fi
  printf '%s' "$candidate"
}

strip=$(lookup_tool STRIP CMAKE_STRIP "$override_strip" strip required)
if [[ "$target_id" == *darwin* ]]; then
  install_name_tool=$(lookup_tool INSTALL_NAME_TOOL CMAKE_INSTALL_NAME_TOOL "$override_install_name_tool" install_name_tool optional)
  otool=$(lookup_tool OTOOL CPKT_OTOOL "$override_otool" otool required)
  readelf=
else
  install_name_tool=$(lookup_tool INSTALL_NAME_TOOL CMAKE_INSTALL_NAME_TOOL "$override_install_name_tool" install_name_tool optional)
  otool=$(lookup_tool OTOOL CPKT_OTOOL "$override_otool" otool optional)
  readelf=$(lookup_tool READELF CMAKE_READELF "$override_readelf" readelf required)
fi

if [[ "$output_format" == cmake ]]; then
  printf 'set(PSLOG_DISCOVERED_CC [==[%s]==])\n' "$cc"
  printf 'set(PSLOG_DISCOVERED_STRIP [==[%s]==])\n' "$strip"
  printf 'set(PSLOG_DISCOVERED_INSTALL_NAME_TOOL [==[%s]==])\n' "$install_name_tool"
  printf 'set(PSLOG_DISCOVERED_OTOOL [==[%s]==])\n' "$otool"
  printf 'set(PSLOG_DISCOVERED_READELF [==[%s]==])\n' "$readelf"
  printf 'set(PSLOG_DISCOVERED_TARGET_ID [==[%s]==])\n' "$target_id"
  printf 'set(PSLOG_DISCOVERED_TARGET_HOST_PREFIX [==[%s]==])\n' "$host_prefix"
else
  printf 'CC=%q\nSTRIP=%q\nINSTALL_NAME_TOOL=%q\nOTOOL=%q\nREADELF=%q\nTARGET_ID=%q\nTARGET_HOST_PREFIX=%q\n' \
    "$cc" "$strip" "$install_name_tool" "$otool" "$readelf" "$target_id" "$host_prefix"
fi
