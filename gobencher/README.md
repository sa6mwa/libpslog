# gobencher

`gobencher` is the Go-side comparison harness for upstream Go `pslog` versus
`libpslog`.

It exists for three jobs:

- compare the Go and C implementations on the same datasets
- compare the embedded Lua binding against the C implementation on the same
  datasets
- render the live elevator pitch chart that puts the Go, C, and embedded-Lua
  variants side by side where that comparison is meaningful

## Important Disclaimer

This suite is intentionally an **attempt** at comparing the Go and C implementations. It is not a recommendation to call `libpslog` from Go in production.

The primary `*C` comparison path crosses the cgo boundary once per benchmark case and uses prebuilt native `pslog_field` data. That is the closest counterpart to the Go benchmark setup because the Go side also prebuilds `Entry.Keyvals` outside the timed loop. Real Go usage of a C logger would usually cross cgo on every log call, which adds large overhead. In practice:

- use the Go logger from Go
- use the C logger from C

## Prerequisite

The root Make targets prepare every required native, Lua, and generated input:

```sh
make gobencher-tests
# or: make benchmarks-gobencher
```

The cgo bridge links against the repo-local installed SDK at
`../build/lua-sdk/lib/libpslog.a` and uses its public headers from
`../build/lua-sdk/include`.
The embedded-Lua compare path loads the shipped rock from `../build/luarocks`.
Do not invoke raw Go commands until one of those root targets has prepared the
selected Bootlin compiler and generated benchmark inputs.

## Useful Commands

```sh
make gobencher-tests
make benchmarks-gobencher
make -C .. elevatorpitch
make -C .. elevatorpitch ELEVATORPITCH_ARGS='-duration=3s -interval=100ms -limit=512'
make -C .. elevatorpitch ELEVATORPITCH_ARGS='-include-quill'
```

From the repository root, `./bench/run_rebaseline.sh` runs the pure C benchmark matrix first and then this Go-vs-C compare suite.

If you want a fail-fast gate instead of an observational rebaseline, run:

```sh
./bench/run_perf_gate.sh
```

That script uses a fresh temporary `GOCACHE` so stale cgo objects do not make the C compare path look broken or artificially slow after header or ABI changes.

The committed source of truth for the production benchmark dataset lives in
`../bench/bench_production_dataset.c`. The Go-side
`benchmark/production_data_generated.go` file is generated from that C fixture
for benchmark-oriented builds so the Go and C runners still use the same input
without carrying two committed copies in the repository.

The benchmark parsers strip source `ts` fields from that dataset. Both the Go
and C benchmark logger constructors then regenerate timestamps from the real
clock using the logger defaults, matching upstream `pslog` benchmark and
elevatorpitch behavior.

## Benchmark Naming

- `*Go`: native Go `pslog`.
- `*C`: primary apples-to-apples compare path using one cgo call per benchmark case and prebuilt native `pslog_field` data.
- `*Ckvfmt`: one cgo call per benchmark case using the real C `logf` public API through generated benchmark-only cgo wrappers. On production data, static fields still go through `with()` once and only the dynamic fields use `kvfmt`.
- `*Cffi`: Go calling into `libpslog` once per log call, so cgo boundary cost is included.
- `*Cprepared`: Go calling prepared native C field data once per log call.
- `*Craw`: one cgo call per benchmark case, with the hot loop and dynamic field rebuilding done inside C.
- `jsonLua`: embedded Lua 5.5 loading the shipped `pslog` rock from
  `build/luarocks`, using the JSON logger with a callback sink inside the Lua
  VM.
- `jsonLiblogger`: optional JSON-only compare against `briandowns/liblogger`, available automatically when the benchmark helper library was built with `-DPSLOG_BENCHMARK_WITH_LIBLOGGER=ON`.
- `jsonQuill`: optional JSON-only compare against Quill, available automatically when the benchmark helper library was built with `-DPSLOG_BENCHMARK_WITH_QUILL=ON`.

## How To Read The Results

- Treat `*C` as the headline compare path.
- Treat `*Ckvfmt` as the honest convenience-API compare for C `logf`/`kvfmt` against Go variadic keyvals.
- Treat `*Cffi` as a diagnostic showing why per-call cgo is not a meaningful deployment model for logger-core comparison.
- Treat `*Craw` as a diagnostic for C-side rebuilding semantics with the cgo cost paid once per benchmark case.
- Treat `jsonLua` as an embedding compare, not a pure logger-core compare. It
  measures the shipped Lua API plus the cost of running inside an embedded Lua
  VM, while still using `libpslog` underneath.
- Treat `jsonLiblogger` as a benchmark-only JSON comparison showing a simpler JSON logger, not a full peer for modern container-style structured JSONL workloads: Unix-seconds timestamps, no `with()` path for persistent structured fields, no color path, no native boolean field type, and much worse performance on the production-shaped benchmark.
- Treat `jsonQuill` as a benchmark-only JSON comparison showing why Quill is not a first-class structured-JSON logger for modern container-style JSONL workloads: no `with()` path for persistent structured fields, `JsonSink` named arguments that are stringified internally instead of preserved as typed fields, a benchmark adapter that reconstructs typed JSON output so Quill pays the full cost of being forced into this class, and much worse performance on the production-shaped benchmark.

## Elevator Pitch

`make elevatorpitch` from the repository root renders the live comparison chart
after preparing the local C library, generated benchmark fixtures, and Lua rock.
The default bar set compares:

- `jsonGo` / `jsonC`
- `jsonLua`
- `coljsonGo` / `coljsonC`
- `consoleGo` / `consoleC`
- `colconGo` / `colconC`

If the optional `Quill` helper was built, `jsonQuill` is available only when `-include-quill` is passed.

`jsonLua` is the only Lua lane in the live chart. The embedded comparison is
intentionally limited to the shipped JSON logger path.
