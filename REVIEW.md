# libpslog Review Tracker

Use this file to coordinate hardening review passes across the logger
implementation. A checked item means the assigned reviewer completed that pass
and found no actionable findings for that category. If a reviewer finds an
issue, leave the item unchecked and track the finding separately until resolved.

Status convention:

- `- [ ] ...` means not started.
- `- [ ] ... (started)` marks the one review category currently in progress.
- `- [x] ...` means complete with no actionable findings from that review pass.

At most one checklist item may include `(started)` at any time. When a started
review pass completes cleanly, remove `(started)` and check the item. When a
new review pass begins, the reviewer may mark any unchecked item as
`(started)`.

## Core Logger Implementation

- [x] Public C API contract and ABI surface
  - Scope: `include/pslog.h`, `include/pslog_version.h.in`, exported symbols,
    documented ownership and lifetime rules, enum values, struct layout, method
    table shape, fatal/panic behavior, and backward-compatibility expectations.
- [x] Logger construction, configuration, and environment overlay
  - Scope: `pslog_default_config()`, `pslog_new()`, `pslog_new_from_env()`,
    level parsing, mode/color/palette parsing, timestamp options, min-level
    behavior, output path parsing, tee/default output handling, and invalid
    configuration fallback/error behavior.
- [x] Logger object lifecycle and shared-state ownership
  - Scope: root and child logger allocation, `with()`, `withf()`,
    `with_level()`, `with_level_field()`, refcounting, `destroy()`, `close()`,
    owned versus borrowed outputs, error preservation, and double-close safety.
- [x] Output abstraction, write retries, partial writes, and close semantics
  - Scope: `pslog_output`, observed output wrappers, file output initialization,
    callback contracts, short-write retry loops, chunked writes, output locking,
    close ownership, null userdata handling, and actionable failure reporting.
- [x] Buffering, chunking, allocation, and oversized-line behavior
  - Scope: `pslog_buffer`, stack versus heap storage, reserve/growth logic,
    flush behavior, line finalization, zero-allocation hot path claims, long
    line streaming, allocation failure handling, and buffer invariant safety.
- [x] Structured field constructors and scalar rendering
  - Scope: `pslog_null`, `pslog_str`, `pslog_trusted_str`, numeric fields,
    boolean fields, pointer fields, byte fields, time fields, duration fields,
    errno fields, cached key/value lengths, trusted flags, and borrowed-data
    lifetime assumptions.
- [x] Level filtering and dispatch semantics
  - Scope: `pslog_should_log()`, `log()`, level-specific methods, free-function
    wrappers, `PSLOG_LEVEL_NOLEVEL`, `PSLOG_LEVEL_DISABLED`, fatal exit path,
    panic abort path, derived logger thresholds, and optional emitted `loglevel`
    fields.

## Encoding and Formatting

- [x] Console emitter formatting
  - Scope: `src/pslog_emit_console.c`, timestamp/level/message layout, static
    field prefixes, dynamic fields, quoting rules, indentation preservation,
    printable Unicode handling, DEL/control/high-bit behavior, errno rendering,
    duration/time formatting, and Go pslog compatibility expectations.
- [x] JSON emitter formatting and validity
  - Scope: `src/pslog_emit_json.c`, JSON object shape, key ordering, string
    escaping, invalid UTF-8 handling, byte encoding, numeric rendering,
    non-finite float policy, null/bool rendering, trusted-string fast paths, and
    JSONL newline discipline.
- [x] Timestamp and time-zone handling
  - Scope: RFC3339/RFC3339Nano paths, custom `strftime` formats, UTC/local
    behavior, cached timestamp invalidation, nanosecond precision, explicit
    `pslog_time_value` rendering, and thread-safe time formatting.
- [x] Color and palette subsystem
  - Scope: `src/pslog_palette.c`, built-in palette registry, alias lookup,
    console ANSI decoration, JSON color metadata, auto-color `isatty` behavior,
    palette string lengths, long palette literal fallback, and plain/color
    output parity after stripping ANSI.
- [x] `kvfmt` parser, caches, and variadic field path
  - Scope: `infof`/`tracef`/etc., `withf()`, supported format verbs, `%m`
    errno snapshotting, max field limits, long key fallback, pointer/content
    caches, cached console/JSON prefixes, varargs type safety assumptions, and
    behavior for malformed or unsupported format strings.
- [x] Trust, escaping, and injection boundaries
  - Scope: `pslog_string_is_trusted()`, trusted key/value flags, console-simple
    detection, JSON escaping bypass conditions, Lua binding trust caches, ANSI
    injection exposure, newline/control-character handling, and documented
    caller responsibility for trusted fast paths.

## Concurrency, Safety, and Portability

- [x] Thread safety and lock ordering
  - Scope: state mutex, output mutex, shared logger trees, concurrent `withf`
    and `kvfmt` cache use, timestamp cache use, close versus emit races, output
    callback reentrancy expectations, and deadlock risk.
- [x] Memory safety and failure handling
  - Scope: all allocation paths, ownership transfer, cleanup on partial
    construction failure, integer overflow in size calculations, null handling,
    borrowed pointer lifetime assumptions, use-after-free risk in caches, and
    sanitizer/fuzz findings.
- [x] C portability and compiler/platform assumptions
  - Scope: C89 build policy, POSIX feature macros, pthread requirement, 32-bit
    and 64-bit integer behavior, musl/glibc differences, Darwin target behavior,
    endian/alignment assumptions, `strerror_r` variants, and cross-compile
    toolchains.
- [x] Performance-sensitive hot paths
  - Scope: normal log emission, prepared-field emission, `kvfmt` emission,
    prefix caches, double-render cache, timestamp cache, branch/lock costs,
    allocation counts, chunk size effects, and benchmark claims versus measured
    results.

## Bindings and Consumer Surfaces

- [x] Lua binding API and lifecycle
  - Scope: `lua/src/pslog_lua.c`, `lua/pslog/init.lua`, Lua constructors,
    output target ownership, callback output, file handles, wrapper values,
    field collection from pairs/tables, logger userdata GC/close behavior, and
    Lua error messages.
- [x] Lua packaging and examples
  - Scope: `lua/lua-pslog.rockspec.in`, `lua/scripts/*`, `lua/examples/*`,
    Lua README contract, module loading, shipped rock contents, version
    rendering, and consistency with the C logger behavior.
- [x] C examples and documentation contract
  - Scope: `README.md`, `examples/*`, API overview snippets, documented
    guarantees, terminology, build examples, and whether examples exercise the
    supported public surface without relying on internals.

## Build, Packaging, and Release

- [x] CMake build graph and target boundaries
  - Scope: `CMakeLists.txt`, presets, shared/static/test targets, compile
    options, include directories, generated version header, dependency fetches,
    optional benchmark dependencies, fuzz target wiring, install rules, and
    target properties.
- [x] Single-header generation and validation
  - Scope: `cmake/package_single_header.cmake`,
    `cmake/pslog_single_header_preamble.in`, generated alias header, embedded
    implementation, license/version content, formatting, and parity with normal
    library tests.
- [x] Binary/source package archive generation
  - Scope: `cmake/package_archive.cmake`, `package_checksums.cmake`,
    `package_clean_dist.cmake`, package metadata, pkg-config/CMake consumer
    files, library names/sonames, stripping behavior, archive layout, and
    checksum output.
- [x] Release privacy and distribution hygiene
  - Scope: `cmake/check_release_privacy.cmake`, release privacy tests, `dist/`
    contents, generated artifacts, absence of build-host paths/secrets/private
    metadata, and reproducibility-sensitive package contents.
- [x] Cross-target release matrix and scripts
  - Scope: `scripts/run_linux_release_matrix.sh`, `scripts/dist-tree.sh`,
    `scripts/clean.sh`, `cmake/toolchains/*`, qemu-emulated tests, Darwin
    packaging path, target IDs, and clean rebuild behavior.

## Verification Assets

- [x] C unit and integration tests
  - Scope: `tests/test_main.c`, `tests/test_main_single_header.c`,
    CMake-driven tests, fixture quality, observable behavior coverage,
    termination subprocess tests, thread tests, output failure tests, and gaps
    between documented contract and executable assertions.
- [x] Fuzzing harness
  - Scope: `fuzz/pslog_fuzz.c`, input modeling, console/JSON mode coverage,
    color parity assertion, field/value generation, sanitizer compatibility,
    crash triage workflow, and missing fuzz targets for high-risk parsers.
- [x] Benchmark harness and performance gates
  - Scope: `bench/*`, `gobencher/*`, benchmark datasets, generated wrappers,
    optional liblogger/Quill adapters, `run_perf_gate.sh`, `run_rebaseline.sh`,
    timing methodology, baseline artifacts, and benchmark-result
    actionability.
- [x] Performance log evidence and baseline management
  - Scope: `performance-logs/*`, baseline/current comparison files,
    provenance of committed measurements, median/run summaries, and whether
    future perf changes have a clear evidence trail.
- [x] CI-equivalent local verification workflow
  - Scope: documented commands, CMake presets, `ctest`, Lua tests, fuzz smoke,
    package tests, format target, coverage script, release matrix, and whether
    the review process can reproduce the expected quality gates locally.
