#!/usr/bin/env bash

set -eu

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

cd "$repo_root"

cmake --preset release \
  -DPSLOG_BENCHMARK_WITH_LIBLOGGER=OFF \
  -DPSLOG_BENCHMARK_WITH_QUILL=OFF
cmake --build --preset release
ctest --preset debug

printf '\n== Pure C benchmark rebaseline ==\n'
./build/release/pslog_bench 500000 all

printf '\n== Go vs C benchmark compare ==\n'
(
  tmpcache="$(mktemp -d)"
  trap 'rm -rf "$tmpcache"' EXIT
  cd gobencher
  GOCACHE="$tmpcache" go test ./...
  GOCACHE="$tmpcache" go test ./benchmark -run '^$' -bench 'Benchmark(Production|Fixed)Compare' -benchmem -benchtime=200ms -count=1
)
