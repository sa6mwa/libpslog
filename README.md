# libpslog

`libpslog` is a high-performance, zero-dependency structured logger for C.

![libpslog_demo](demo.gif)

`libpslog` is the C port of the Go [pkt.systems/pslog](https://github.com/sa6mwa/pslog) logger and aims to preserve the same overall product shape: console and JSON output, strong structured logging support, careful control over allocations, aggressive benchmarking, and output semantics that are close to the Go implementation. The Go variant is already one of the fastest loggers in that ecosystem; `libpslog` brings the same performance-first design to C while exposing a C-shaped API.

![libpslog elevator pitch](elevatorpitch.gif)

## What It Provides

- Console and JSON loggers from the same API surface.
- Colorized and non-colorized output selected from config, with adaptive tty-aware color by default.
- Structured fields through typed `pslog_field[]` arrays.
- Structured `kvfmt` logging through `tracef/debugf/infof/warnf/errorf`, where `infof` means `message + kvfmt`, not `printf`-formatting the message itself.
- Derived loggers through `with()`, `withf()`, `with_level()`, and `with_level_field()`.
- Environment-driven construction through `pslog_new_from_env()`.
- Explicit trusted-string fast paths for JSON emission.
- Thread-safe shared logging through the same logger tree.
- Zero-allocation hot-path emission for normal log lines, with chunked output for oversized lines instead of truncation.
- Unit tests, fuzzing, pure C benchmarks, and Go-vs-C comparison benchmarks.

## Release Targets

First-release Linux support is:

- `x86_64-linux-gnu`
- `x86_64-linux-musl`
- `aarch64-linux-gnu`
- `aarch64-linux-musl`
- `armhf-linux-gnu`
- `armhf-linux-musl`

All six are wired for:

- configure/build via CMake presets
- tests
- runtime package generation
- dev package generation

For ARM targets, the test presets run under qemu.

## API Overview

Public symbols use the `pslog_` prefix. Preferred usage is through instance methods on `pslog_logger`.

Core structured path:

```c
#include "pslog.h"

pslog_config config;
pslog_logger *log;
pslog_field fields[2];

pslog_default_config(&config);
config.mode = PSLOG_MODE_JSON;
config.color = PSLOG_COLOR_NEVER;
config.output = pslog_output_from_fp(stdout, 0);

log = pslog_new(&config);

fields[0] = pslog_str("service", "api");
fields[1] = pslog_i64("attempt", 3L);
log->info(log, "request handled", fields, 2u);

log->destroy(log);
```

Derived logger path:

```c
pslog_field base[1];
pslog_logger *child;

base[0] = pslog_str("subsystem", "worker");
child = log->with(log, base, 1u);
child = child->with_level_field(child);
child->warn(child, "retrying", NULL, 0u);
child->destroy(child);
```

Derived logger from `kvfmt`:

```c
child = log->withf(log, "service=%s subsystem=%s", "api", "worker");
child->info(child, "request handled", NULL, 0u);
child->destroy(child);
```

`kvfmt` path:

```c
log->infof(log, "request handled", "user=%s code=%d ok=%b ms=%f",
           "alice", 200, 1, 12.34);
```

Direct built-in palette selection:

```c
pslog_default_config(&config);
config.mode = PSLOG_MODE_CONSOLE;
config.color = PSLOG_COLOR_ALWAYS;
config.palette = &pslog_builtin_palette_nord;
```

Additional typed fields include:

- `pslog_bool`
- `pslog_bytes_field`
- `pslog_duration_field`
- `pslog_f64`
- `pslog_i64`
- `pslog_null`
- `pslog_ptr`
- `pslog_time_field`
- `pslog_trusted_str`
- `pslog_u64`

Relevant behavioral notes:

- `PSLOG_LEVEL_NOLEVEL` emits `---` in console mode and `"nolevel"` in JSON mode.
- `log->fatal(...)` and `log->fatalf(...)` log and then exit with failure status.
- `log->panic(...)` and `log->panicf(...)` log and then abort.
- The free-function `pslog_fatal(...)` / `pslog_fatalf(...)` wrappers also terminate.
- The free-function `pslog_panic(...)` / `pslog_panicf(...)` wrappers also abort.
- `with()` returns a derived logger and does not mutate the receiver.
- `close()` closes only owned outputs.
- `pslog_new_from_env()` overlays `LOG_*` settings on top of a seed config.

The public header in [include/pslog.h](include/pslog.h) contains the full API contract and doc comments.

## Examples

The main example is [examples/example.c](examples/example.c). It demonstrates:

- console and JSON loggers
- `log()` and level-specific methods
- `with()`, `with_level()`, and `with_level_field()`
- `withf()` for static `kvfmt`-derived fields
- typed fields
- `infof`/`kvfmt`
- `pslog_new_from_env()`
- palette iteration
- adaptive color behavior

Build it in normal library mode:

```sh
cmake --preset host
cmake --build --preset host
cd examples
cc -I../build/host/generated/include -I../include \
  -o example example.c ../build/host/libpslog.a
./example
```

Build the same example in single-header mode:

```sh
cmake --preset host
cmake --build --preset package-single-header
cd examples
cc -DPSLOG_EXAMPLE_SINGLE_HEADER=1 \
  -I../build/host/generated/include \
  -o example example.c
./example
```

The shipped single-header artifacts are written to `dist/pslog-<version>.h`
and `dist/pslog-<version>.h.gz`. The generated header contains the public API,
the `PSLOG_IMPLEMENTATION` section, the embedded license text, and the
single-header usage notes at the top of the file.

## Build And Test

If you want a simple frontend instead of typing the CMake commands directly, use the repository `Makefile`:

```sh
make build
make test
make test-all
make fuzz
make benchmarks-c
make benchmarks-gobencher
make benchmarks-all
make release
```

Run `make help` for the full target list.

Standard debug build:

```sh
cmake --preset debug
cmake --build --preset debug
ctest --preset debug
```

Address-sanitized run:

```sh
cmake --preset asan
cmake --build --preset asan
ctest --preset asan
```

Fuzzing:

```sh
cmake --preset fuzz
cmake --build --preset fuzz
./build/fuzz/pslog_fuzz -runs=1000
```

## Linux Release Matrix

One-command sweep for the full shipped Linux matrix:

```sh
./scripts/run_linux_release_matrix.sh
```

That script runs, for every shipped Linux target:

- `cmake --preset ...`
- `cmake --build --preset ...`
- `ctest --preset ...`
- runtime package generation
- dev package generation
- final `libpslog-<version>-CHECKSUMS` generation over every shipped file in `dist/`

Toolchain expectations:

- `linux-gnu` cross presets expect distro cross compilers such as `aarch64-linux-gnu-gcc` and `arm-linux-gnueabihf-gcc`.
- `x86_64-linux-musl` expects host `musl-gcc`.
- `aarch64-linux-musl` and `armhf-linux-musl` expect real musl cross compilers such as `aarch64-linux-musl-gcc` and `arm-linux-musleabihf-gcc`.
- musl ARM qemu runs expect the musl loader symlink in the target sysroot to resolve within the prefix, for example `ld-musl-aarch64.so.1 -> libc.so`.

Single-target examples:

```sh
cmake --preset aarch64-linux-gnu-release
cmake --build --preset aarch64-linux-gnu-release
ctest --preset aarch64-linux-gnu-release

cmake --build --preset package-runtime-aarch64-linux-gnu
cmake --build --preset package-dev-aarch64-linux-gnu
```

Artifacts are written to `dist/`.

After the full matrix runs, `dist/` also contains `libpslog-<version>-CHECKSUMS`
with `sha256sum` output for every shipped release file except the checksum file
itself. The listed filenames are bare archive names, not `dist/...` paths.

## Benchmarks

There are two benchmark layers:

- pure C benchmarks in [bench](bench)
- Go-vs-C comparison harness in [gobencher](gobencher)

Useful commands:

```sh
./build/host/pslog_bench 200000 all
./bench/run_rebaseline.sh
```

The benchmark suite covers:

- fixed synthetic workloads
- production-shaped workloads
- prebuilt field-array paths
- per-call field rebuild paths
- `with()`-based production shapes
- `kvfmt` convenience paths
- optional benchmark-only `liblogger` comparison on JSON
- optional benchmark-only `Quill` comparison on JSON

The external comparisons are intentionally scoped:

- `liblogger` is kept only as a simple JSON baseline. It emits JSON, but it is still not a full peer for modern container-style structured JSONL workloads: no `with()`-style persistent fields, no native boolean field type, and a narrower JSON surface than `libpslog`. It is also far slower on the production-shaped benchmark, so this is not just a feature gap.
- `Quill` is kept only as an opt-in negative comparison. Its `JsonSink` is not first-class structured JSON for modern container-style JSONL logging: it does not preserve typed structured fields, it stringifies named arguments internally, and it lacks persistent structured-field attachment. It is also far slower on the production-shaped benchmark once forced into this workload class. If you are evaluating JSON loggers for Kubernetes, serverless, or other container-first workloads, do not treat Quill's `JsonSink` as the same kind of JSON logging product as `libpslog`.

See [bench/README.md](bench/README.md) and [gobencher/README.md](gobencher/README.md) for naming and interpretation details.

## Environment Configuration

`pslog_new_from_env()` reads settings like:

- `LOG_MODE`
- `LOG_LEVEL`
- `LOG_DISABLE_TIMESTAMP`
- `LOG_VERBOSE_FIELDS`
- `LOG_NO_COLOR`
- `LOG_FORCE_COLOR`
- `LOG_PALETTE`
- `LOG_OUTPUT`
- `LOG_OUTPUT_FILE_MODE`
- `LOG_TIME_FORMAT`
- `LOG_UTC`

Example:

```c
pslog_config seed;
pslog_logger *log;

pslog_default_config(&seed);
seed.output = pslog_output_from_fp(stdout, 0);

log = pslog_new_from_env("LOG_", &seed);
log->infof(log, "hello", "key=%s", "value");
log->destroy(log);
```

## Operational Notes

- Shared logging through the same logger tree is thread-safe.
- `infof/debugf/...` do structured `kvfmt`; they are not `printf` message formatters.
- Normal emission stays allocation-free on the hot path for ordinary line sizes.
- Long lines are streamed in chunks instead of truncated.
- Trusted strings are explicit. Untrusted strings are escaped normally.
- The runtime `.so` and static `.a` packages both ship the public header.
