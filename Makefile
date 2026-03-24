SHELL := bash
.DEFAULT_GOAL := help
MAKEFLAGS += --no-builtin-rules

DEBUG_PRESET := debug
HOST_PRESET := host
ASAN_PRESET := asan
FUZZ_PRESET := fuzz
COVERAGE_PRESET := coverage
RELEASE_BUILD_PRESETS := \
	linux-gnu-release \
	linux-musl-release \
	aarch64-linux-gnu-release \
	aarch64-linux-musl-release \
	armhf-linux-gnu-release \
	armhf-linux-musl-release

CROSS_RELEASE_PRESETS := \
	aarch64-linux-gnu-release \
	aarch64-linux-musl-release \
	armhf-linux-gnu-release \
	armhf-linux-musl-release

BENCH_ITERS ?= 200000
FUZZ_TIME ?= 30
GO_BENCH_COUNT ?= 1

.PHONY: \
	help \
	build \
	build-host \
	build-release \
	format \
	test \
	test-all \
	asan \
	coverage \
	fuzz \
	benchmarks-c \
	benchmarks-gobencher \
	benchmarks-go \
	benchmarks-all \
	benchmarks \
	gobencher-tests \
	perf-gate \
	cross-build \
	cross-test \
	release \
	clean \
	clean-dist

help:
	@printf '%s\n' \
		'make build           Configure and build the debug preset.' \
		'make build-host      Configure and build the host-native local build.' \
		'make build-release   Configure and build the full shipped release build matrix.' \
		'make format          Run clang-format on repo C/C header sources.' \
		'make test            Run the debug C test suite.' \
		'make test-all        Run C tests, Go gobencher tests, and the Go-vs-C perf gate.' \
		'make asan            Run the ASan/UBSan preset test suite.' \
		'make coverage        Run the coverage preset and generate coverage-report.' \
		'make fuzz            Build fuzz target and run two bounded fuzz passes.' \
		'make benchmarks-c    Run the pure C benchmark matrix (BENCH_ITERS=$(BENCH_ITERS)).' \
		'make benchmarks-gobencher Run the full gobencher benchmark suite.' \
		'make benchmarks-all  Run both C and gobencher benchmark suites.' \
		'make benchmarks      Alias for make benchmarks-all.' \
		'make gobencher-tests Run all Go gobencher tests.' \
		'make perf-gate       Run the Go-vs-C performance gate.' \
		'make cross-build     Build all non-host cross release presets.' \
		'make cross-test      Test all non-host cross release presets.' \
		'make release         Run the full Linux release matrix and package generation.' \
		'make clean           Remove build/ and dist/ generated artifacts.' \
		'make clean-dist      Remove dist/ release artifacts.'

build:
	cmake --preset $(DEBUG_PRESET)
	cmake --build --preset $(DEBUG_PRESET)

build-host:
	cmake --preset $(HOST_PRESET)
	cmake --build --preset $(HOST_PRESET)

build-release:
	@set -e; for preset in $(RELEASE_BUILD_PRESETS); do \
		cmake --preset "$$preset"; \
		cmake --build --preset "$$preset"; \
	done

format:
	cmake --preset $(DEBUG_PRESET)
	cmake --build --preset format

test: build
	ctest --preset $(DEBUG_PRESET) --output-on-failure

gobencher-tests: build-host
	cd gobencher && go test -a ./...

perf-gate:
	./bench/run_perf_gate.sh

test-all: test gobencher-tests perf-gate

asan:
	cmake --preset $(ASAN_PRESET)
	cmake --build --preset $(ASAN_PRESET)
	ctest --preset $(ASAN_PRESET) --output-on-failure

coverage:
	cmake --preset $(COVERAGE_PRESET)
	cmake --build --preset $(COVERAGE_PRESET)
	ctest --preset $(COVERAGE_PRESET) --output-on-failure
	cmake --build --preset coverage-report

fuzz:
	cmake --preset $(FUZZ_PRESET)
	cmake --build --preset $(FUZZ_PRESET)
	./build/fuzz/pslog_fuzz -max_total_time=$(FUZZ_TIME)
	./build/fuzz/pslog_fuzz -max_total_time=$(FUZZ_TIME) -max_len=256

benchmarks-c: build-host
	./build/host/pslog_bench $(BENCH_ITERS) all

benchmarks-gobencher: build-host
	cd gobencher && go test ./benchmark -run '^$$' -bench . -benchmem -count=$(GO_BENCH_COUNT)

benchmarks-go: benchmarks-gobencher

benchmarks-all: benchmarks-c benchmarks-gobencher

benchmarks: benchmarks-all

cross-build:
	@set -e; for preset in $(CROSS_RELEASE_PRESETS); do \
		cmake --preset "$$preset"; \
		cmake --build --preset "$$preset"; \
	done

cross-test: cross-build
	@set -e; for preset in $(CROSS_RELEASE_PRESETS); do \
		ctest --preset "$$preset" --output-on-failure; \
	done

release:
	./scripts/run_linux_release_matrix.sh

clean:
	./scripts/clean.sh

clean-dist:
	./scripts/clean.sh --dist-only
