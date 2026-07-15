# Benchmarks

The benchmark surface is split into two layers:

- `pslog_bench` for pure C logger cost
- `gobencher` for Go-vs-C comparison

For a reproducible local rebaseline from the repository root:

```sh
./bench/run_rebaseline.sh
```

That script rebuilds the release binary, runs the matching host test preset,
runs the full pure C benchmark matrix, then runs the Go-vs-C compare suite in
`gobencher`.

For a hard fail-fast performance regression gate:

```sh
make perf-gate
```

The Make target supplies the configured host `CC` and `CXX` to cgo and the Lua
build. For the normal native `host` preset that is the local compiler selected
by CMake, with `c++` as the fallback C++ compiler for Go/cgo. If the host build
is sysroot-backed, the same benchmark runner executes binaries through that
configured loader. Invoke the script directly only when those variables already
name the intended host toolchain.

That gate builds the host release target, runs the host CTest preset, builds the
repo-local Lua rock, then checks:

- pure C `pslog_bench` rows against `performance-logs/pure-c-baseline.txt`
- embedded Lua 5.5 `jsonLua` and table-form rows against
  `performance-logs/lua-baseline.txt`

By default, pure C and Lua rows may be up to `50%` slower than the checked-in
baselines. Override those limits with `PSLOG_PERF_C_TOLERANCE` and
`PSLOG_PERF_LUA_TOLERANCE` when running on a controlled machine. The pure C
run uses `taskset -c 0` when available because the current C baseline was
recorded that way; override the CPU with `PSLOG_PERF_CPU`, or set it to an
empty value to disable pinning. The script still prints an observational
Go-vs-C compare, but Go-relative timing no longer defines the hard pass/fail
contract.

Both the pure C and Go-vs-C benchmark paths regenerate timestamps from the
real clock. The production dataset strips source `ts` fields first, so the
measured logger work includes the actual timestamp path instead of replaying
pre-rendered timestamp text.

The committed source of truth for that production dataset is
`bench_production_dataset.c`. Other benchmark fixture forms are generated from
it as build artifacts when needed.

Committed benchmark artifacts belong in
[performance-logs/](../performance-logs/README.md), not under `build/`. Treat
`performance-logs/pure-c-baseline.txt` and
`performance-logs/lua-baseline.txt` as the current checked-in baselines when
recording new comparison runs.

## Pure C Benchmarks

Run the whole matrix:

```sh
make benchmarks-c
```

Representative focused run:

```sh
./scripts/run_host_binary.sh ./build/host/pslog_bench 200000 \
  console_prod_with_log_fields \
  console_prod_with_level_fields_build \
  console_prod_with_levelf_kvfmt \
  json_prod_with_log_fields \
  json_prod_with_level_fields_build \
  json_prod_with_levelf_kvfmt
```

Naming:

- `*_api`: fixed synthetic dataset with field constructors on every call
- `*_prepared`: fixed synthetic dataset with reused `pslog_field[]`
- `*_prod_log_fields`: production dataset with prebuilt `pslog_field[]` passed to `log(...)`
- `*_prod_with_log_fields`: same, but static fields are attached once via `with()`
- `*_prod_log_fields_build`: production dataset where `pslog_field[]` is rebuilt each call and passed to `log(...)`
- `*_prod_with_log_fields_build`: same, with static fields attached once via `with()`
- `*_prod_level_fields_build`: production dataset where `pslog_field[]` is rebuilt each call and dispatched through `trace/debug/info/warn/error`
- `*_prod_with_level_fields_build`: same, with static fields attached once via `with()`
- `*_prod_levelf_kvfmt`: production dataset through `tracef/debugf/infof/warnf/errorf`
- `*_prod_with_levelf_kvfmt`: same, with static fields attached once via `with()`

Interpretation:

- `*_prod_with_*` is the closest production-style shape for long-lived subsystem loggers with static fields attached once.
- `*_prod_log_fields` and `*_prod_with_log_fields` are the closest logger-core measurements because the `pslog_field[]` data is prebuilt before timing.
- `*_build` names measure the cost of rebuilding `pslog_field[]` on every log call.
- `*_levelf_kvfmt` names measure the convenience `kvfmt` path.

## Go-vs-C Compare

The compare naming used in `gobencher` is:

- `*Go`: upstream Go logger in Go
- `*C`: primary apples-to-apples compare path using one cgo call per benchmark case and prebuilt native `pslog_field` data
- `*Ckvfmt`: one cgo call per benchmark case using the real C `logf`/`kvfmt` public API through generated benchmark-only cgo wrappers
- `*Cffi`: Go calling the C logger once per log call
- `*Cprepared`: Go calling prebuilt native C field data once per log call
- `*Craw`: one cgo call per benchmark case, with dynamic field rebuilding done inside C
- `jsonLua`: embedded Lua 5.5 loading the shipped `pslog` rock from `build/luarocks`, limited to the JSON logger path

See [gobencher/README.md](../gobencher/README.md) for the caveats and the elevator pitch runner.

## Optional liblogger Compare

`liblogger` is benchmark-only. To enable it:

```sh
cmake --preset host -DPSLOG_BENCHMARK_WITH_LIBLOGGER=ON
cmake --build --preset host
./scripts/run_host_binary.sh ./build/host/pslog_bench 200000 liblogger_json liblogger_json_prod
```

This fetches both `liblogger` and its `jansson` dependency only for the benchmark build.

Important caveats:

- this compare is JSON-only
- `liblogger` emits JSON directly, but it is still not a full peer for modern container-style structured JSONL workloads
- `liblogger` does not expose `with()`-style persistent structured fields
- `liblogger` does not expose a native boolean field type on this path
- `liblogger` is also far slower on the production-shaped benchmark, so the issue is not just missing features
- the compare is intentionally limited to non-colored JSON because that is the surface `liblogger` actually exposes

Once the release build was produced with that option, `gobencher` will include
the `jsonLiblogger` compare automatically.

## Optional Quill Compare

`Quill` is benchmark-only. To enable it:

```sh
cmake --preset host -DPSLOG_BENCHMARK_WITH_QUILL=ON
cmake --build --preset host
```

The `gobencher` benchmark suite will include `jsonQuill` automatically, while
the live elevator pitch includes it only when you pass `-include-quill`.

Important caveats:

- this compare is JSON-only
- Quill does not expose `with()`-style persistent structured fields on this path, so static fields are replayed on every event
- Quill's `JsonSink` is not first-class structured JSON for modern container-style JSONL logging: it surfaces named arguments as formatted strings instead of emitting typed structured fields directly
- Quill is also far slower on the production-shaped benchmark once forced into this workload class, so this is not just an API mismatch on paper
- if you are looking for a JSON logger for Kubernetes, serverless, or other container-first workloads, Quill should not be read as offering the same kind of JSON logging product as `libpslog`
- the benchmark adapter reconstructs typed JSON output from Quill's stringified surface so Quill is measured only after paying the full cost of being forced into this class
