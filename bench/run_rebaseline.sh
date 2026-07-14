#!/usr/bin/env bash

set -eu

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

compiler_from_host_cache() {
  local key=$1 cache_file="$repo_root/build/host/CMakeCache.txt"
  [[ -f "$cache_file" ]] || return 1
  sed -n "s/^${key}:FILEPATH=//p" "$cache_file" | tail -n 1
}

CC=${CC:-$(compiler_from_host_cache CMAKE_C_COMPILER || true)}
CXX=${CXX:-$(compiler_from_host_cache CMAKE_CXX_COMPILER || true)}
if [[ -z "$CC" || ! -x "$CC" || -z "$CXX" || ! -x "$CXX" ]]; then
  printf 'rebaseline requires configured Bootlin CC and CXX; run make build-host first\n' >&2
  exit 1
fi
export CC CXX

host_executor="${PSLOG_HOST_EXECUTOR:-$repo_root/scripts/run_host_binary.sh}"
if [ ! -x "$host_executor" ]; then
  printf 'rebaseline requires PSLOG_HOST_EXECUTOR to name the native Bootlin sysroot runner\n' >&2
  exit 1
fi

cd "$repo_root"

cmake --preset host \
  -DPSLOG_BENCHMARK_WITH_LIBLOGGER=OFF \
  -DPSLOG_BENCHMARK_WITH_QUILL=OFF
cmake --build --preset host
ctest --preset host

printf '\n== Pure C benchmark rebaseline ==\n'
"$host_executor" ./build/host/pslog_bench 500000 all

printf '\n== Go vs C benchmark compare ==\n'
(
  tmpcache="$(mktemp -d)"
  trap 'rm -rf "$tmpcache"' EXIT
  cd gobencher
  GOCACHE="$tmpcache" go test -exec "$host_executor" ./...
  GOCACHE="$tmpcache" go test -exec "$host_executor" ./benchmark -run '^$' -bench 'Benchmark(Production|Fixed)Compare' -benchmem -benchtime=200ms -count=1
)
