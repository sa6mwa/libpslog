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
- Optional Lua binding release artifacts shipped as `lua-pslog`.

## Release Targets

First-release binary support is:

- `x86_64-linux-gnu`
- `x86_64-linux-musl`
- `aarch64-linux-gnu`
- `aarch64-linux-musl`
- `armhf-linux-gnu`
- `armhf-linux-musl`
- `arm64-apple-darwin`

All targets are wired for:

- configure/build via CMake presets
- tests where the target can run locally or under qemu
- runtime package generation
- dev package generation

For Linux ARM targets, the test presets run under qemu. These are mandatory
release routes: `make release-matrix` and `make release` fail if `qemu-aarch64`
or `qemu-arm` is unavailable. The Darwin target is a build-and-package target
when the local osxcross toolchain is available.

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
pslog_logger *next;

base[0] = pslog_str("subsystem", "worker");
child = log->with(log, base, 1u);
next = child->with_level_field(child);
child->destroy(child);
child = next;
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
"$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' ../build/host/CMakeCache.txt)" -I../build/host/generated/include -I../include \
  -o example example.c ../build/host/libpslog.a -pthread
./example
```

For an extracted binary SDK, point pkg-config at that SDK and use the compiler
selected for the matching target rather than an ambient host compiler:

```sh
sdk=/path/to/libpslog-<version>-<target>
export PKG_CONFIG_PATH="$sdk/lib/pkgconfig${PKG_CONFIG_PATH:+:$PKG_CONFIG_PATH}"
"$CC" $(pkg-config --cflags pslog) -o example example.c $(pkg-config --libs pslog)
```

```cmake
set(CMAKE_PREFIX_PATH "/path/to/libpslog-<version>-<target>")
find_package(pslog CONFIG REQUIRED)
target_link_libraries(example PRIVATE pslog::pslog)
```

The default CMake target links the shared library. Use `pslog::pslog_static`
when a static link is preferred.

Build the same example in single-header mode:

```sh
cmake --preset host
cmake --build ../build/host --target package-single-header
cd examples
"$(sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' ../build/host/CMakeCache.txt)" -DPSLOG_EXAMPLE_SINGLE_HEADER=1 \
  -I../build/host/generated/include \
  -o example example.c -pthread
./example
```

The release single-header artifact is written to `dist/pslog-<version>.h.gz`.
The uncompressed generated header remains in the build tree at
`build/host/generated/include/pslog_single_header.h` for local compile smoke
tests. The generated header contains the public API, the `PSLOG_IMPLEMENTATION`
section, the embedded license text, and the single-header usage notes at the top
of the file.

## Lua Binding

`libpslog` ships an optional Lua binding as the `lua-pslog` rock package. The
Lua module itself is required as `pslog`:

```lua
local pslog = require("pslog")

local log = pslog.new_json("/tmp/app.log", {
  no_color = true,
  disable_timestamp = true,
}):with("service", "checkout")

log:info("ready", "port", 8080, "ok", true)
```

The Lua binding keeps the same logger product shape as C and Go, but the API is
written for Lua instead of mirroring the C surface literally:

- ordinary Lua scalars are inferred automatically
- structured logging through key/value pairs is the primary call shape
- a single-table field form exists as convenience
- explicit wrappers exist for semantic types Lua does not natively represent
  such as bytes, time, duration, errno, trusted strings, and `u64`
- constructors support stdout/stderr, file handles, file paths, and callback
  sinks

Examples:

```lua
log:info("typed fields",
  "payload", pslog.bytes(string.char(0, 1, 2, 255)),
  "started_at", pslog.time { sec = 1711737600, nsec = 0, offset = 0 },
  "elapsed", pslog.duration_ms(125),
  "errno", pslog.errno(2)
)
```

Local Lua development helpers:

```sh
make lua-rock
make lua-test
```

`make lua-rock` first installs the Bootlin-built public C SDK into the
repo-local `build/lua-sdk/` prefix. The Lua module, C interop checks, and Go
Lua benchmark bridge all consume that installed SDK rather than source-tree
headers or libraries.

Run the Lua examples from the repository root:

```sh
make lua-rock
eval "$(make lua-env)"
lua lua/examples/example.lua
lua lua/examples/basic.lua
lua lua/examples/from_env.lua
lua lua/examples/callback.lua
```

Example entry points live under [`lua/examples/`](lua/examples/):

- [`lua/examples/example.lua`](lua/examples/example.lua)
- [`lua/examples/basic.lua`](lua/examples/basic.lua)
- [`lua/examples/callback.lua`](lua/examples/callback.lua)
- [`lua/examples/from_env.lua`](lua/examples/from_env.lua)

The full Lua API reference lives in [`lua/README.md`](lua/README.md).

`make release` also emits:

- `dist/lua-pslog-<version>.tar.gz`
- `dist/lua-pslog-<version>-1.rockspec`
- `dist/lua-pslog-<version>-1.src.rock`

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
make lua-test
make prerelease
make prerelease-hardening
make release
```

Run `make help` for the full target list.

Standard debug build:

```sh
cmake --preset debug
cmake --build --preset debug
ctest --preset debug
```

Native Valgrind check:

```sh
make valgrind
```

Fuzzing:

```sh
make fuzz
```

## Release Matrix

One-command incremental rehearsal for the full shipped matrix:

```sh
make release-matrix
```

`make prerelease` runs the shared deterministic proof graph without cleaning
generated state first. `make release` cleans first and then runs that exact
same proof graph; it is the final local release gate.
`make prerelease-hardening` adds the explicit long AFL++ fuzz tier to that
proof graph. `make release` writes a tab-separated phase timing report to
`build/release-timings.tsv`; its matrix entries distinguish configure, build,
QEMU-backed test, packaging, source smoke, Lua artifacts, checksums, and
privacy verification. It also separates ordinary C tests, Valgrind, AFL++
smoke, cross tests, Go tests, and the performance gate.

Inspect the latest completed report with:

```sh
column -t -s $'\t' build/release-timings.tsv
```

That script runs, for every shipped Linux target:

- `cmake --preset ...`
- `cmake --build --preset ...`
- `ctest --preset ...`
- runtime package generation
- dev package generation
- `lua-pslog-<version>-1.rockspec` generation
- `lua-pslog-<version>-1.src.rock` generation
- final `libpslog-<version>-CHECKSUMS` generation over every shipped file in `dist/`

When the osxcross arm64 Darwin compiler is available, the script also builds
and packages `arm64-apple-darwin`.

Toolchain expectations:

- Linux presets provision checksum-pinned Bootlin GCC collections through `scripts/cpkt-toolchains.sh`; host compilers and binutils are never selected.
- Native x86_64 CTest, Valgrind, and AFL++ execution launch through the selected Bootlin sysroot loader, so they do not depend on a matching host glibc or musl installation. Native memory checking still uses host Valgrind against a focused Bootlin-built facade test; native x86_64 fuzzing uses the cached AFL++ GCC-plugin wrapper from `scripts/cpkt-aflpp.sh`, which delegates to the same Bootlin collection.
- `clang-format` and `clangd` are host development tools only. `make clangd` checks the native public C consumer with `build/debug/compile_commands.json`, including its public-header surface. Public declarations use Doxygen comments so hover documentation remains useful in clangd. clangd is not a compiler, target-ABI verifier, package check, or release dependency; cross builds, packages, and releases do not invoke it.
- The shared toolchain cache is `${CPKT_TOOLCHAIN_CACHE:-${XDG_CACHE_HOME:-$HOME/.cache}/c.pkt.systems/toolchains}`. Its Bootlin and AFL++ provisioners serialize each collection with a bounded `CPKT_TOOLCHAIN_LOCK_TIMEOUT` (600 seconds by default), and the AFL++ cache identity includes the selected Bootlin collection and sysroot. External dependency archives use `${CPKT_DEPENDENCY_CACHE:-${XDG_CACHE_HOME:-$HOME/.cache}/c.pkt.systems/deps}`. Both survive `make clean`; only disposable extracted dependency state under the repository’s `.cache/` is removed.
- Cross test execution requires `qemu-aarch64` and `qemu-arm`; each uses the matching Bootlin sysroot.
- `arm64-apple-darwin` expects osxcross under `OSXCROSS_ROOT` or `$HOME/.local/cross/osxcross`.

Single-target examples:

```sh
cmake --preset aarch64-linux-gnu-release
cmake --build --preset aarch64-linux-gnu-release
ctest --preset aarch64-linux-gnu-release

cmake --build build/aarch64-linux-gnu-release --target package-archive
```

Each target archive in `dist/` now contains the public headers, the shared
library, and the static library for that target variant.

After the full matrix runs, `dist/` also contains `libpslog-<version>-CHECKSUMS`
with `sha256sum` output for every shipped release file except the checksum file
itself. The listed filenames are bare archive names, not `dist/...` paths.

## Benchmarks

There are two benchmark layers:

- pure C benchmarks in [bench](bench)
- Go-vs-C comparison harness in [gobencher](gobencher)

Useful commands:

```sh
make benchmarks-c
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
