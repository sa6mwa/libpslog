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
LUA_RELEASE_VERSION := $(shell ./lua/scripts/release_version.sh)
LUA_DIST_DIR := $(CURDIR)/dist
LUA_RELEASE_ROCKSPEC := $(LUA_DIST_DIR)/lua-pslog-$(LUA_RELEASE_VERSION)-1.rockspec
LUA_RELEASE_PACK_DIR := $(LUA_DIST_DIR)/.lua-pack
LUA_RELEASE_STAGE_DIR := $(LUA_RELEASE_PACK_DIR)/src
LUA_RELEASE_PACK_ROCKSPEC := $(LUA_RELEASE_PACK_DIR)/lua-pslog-$(LUA_RELEASE_VERSION)-1.rockspec
LUA_RELEASE_ROCK := $(LUA_DIST_DIR)/lua-pslog-$(LUA_RELEASE_VERSION)-1.src.rock
LUA_ROCK_TREE := build/luarocks
LUA_ROCKSPEC := $(LUA_ROCK_TREE)/lua-pslog-$(LUA_RELEASE_VERSION)-1.rockspec
LUA_ROCK_STAMP := $(LUA_ROCK_TREE)/.installed.stamp
LUA_ROCK_BUILD_LOCK := $(LUA_ROCK_TREE)/.build.lock
LUA_ROCK_BUILD_BYPRODUCTS := \
	$(CURDIR)/pslog \
	$(CURDIR)/lua/src/pslog_lua.o \
	$(CURDIR)/src/pslog.o \
	$(CURDIR)/src/pslog_emit_console.o \
	$(CURDIR)/src/pslog_emit_json.o \
	$(CURDIR)/src/pslog_palette.o
LUA_ROCK_SOURCES := \
	lua/lua-pslog.rockspec.in \
	lua/scripts/render_release_rockspec.sh \
	lua/scripts/release_version.sh \
	lua/pslog/init.lua \
	lua/src/pslog_lua.c \
	include/pslog.h \
	src/pslog.c \
	src/pslog_emit_console.c \
	src/pslog_emit_json.c \
	src/pslog_palette.c
GO_PRODUCTION_DATASET := gobencher/benchmark/production_data_generated.go
GO_PRODUCTION_DATASET_SOURCE := bench/bench_production_dataset.c
GO_PRODUCTION_DATASET_TOOL := bench/gen_go_production_dataset.c
GO_CKVFMT_WRAPPERS := gobencher/benchmark/cpslog_kvfmt_generated.go
HOST_GENERATED_VERSION_HEADER := $(CURDIR)/build/host/generated/include/pslog_version.h

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
	prepare-gobencher-data \
	perf-gate \
	cross-build \
	cross-test \
	release-lua-artifacts \
	lua-rock \
	lua-test \
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
		'make lua-rock        Build and install the Lua module into build/luarocks.' \
		'make lua-test        Run the Lua binding smoke tests.' \
		'make release         Run the full Linux release matrix and package generation.' \
		'make clean           Remove build/ and dist/ generated artifacts.' \
		'make clean-dist      Remove dist/ release artifacts.'

build:
	cmake --preset $(DEBUG_PRESET)
	cmake --build --preset $(DEBUG_PRESET)

build-host:
	cmake --preset $(HOST_PRESET)
	cmake --build --preset $(HOST_PRESET)
	$(MAKE) $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)

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

gobencher-tests: build-host lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)
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

benchmarks-gobencher: build-host lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)
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

release-lua-artifacts: $(LUA_RELEASE_ROCKSPEC) $(LUA_RELEASE_ROCK)

$(LUA_RELEASE_PACK_DIR):
	mkdir -p "$(LUA_RELEASE_PACK_DIR)"

$(LUA_DIST_DIR):
	mkdir -p "$(LUA_DIST_DIR)"

$(LUA_RELEASE_STAGE_DIR): Makefile .git/index | $(LUA_RELEASE_PACK_DIR)
	rm -rf "$(LUA_RELEASE_STAGE_DIR)"
	mkdir -p "$(LUA_RELEASE_STAGE_DIR)"
	git archive --format=tar --worktree-attributes HEAD | tar -xf - -C "$(LUA_RELEASE_STAGE_DIR)"
	git -C "$(LUA_RELEASE_STAGE_DIR)" init -q
	git -C "$(LUA_RELEASE_STAGE_DIR)" config user.name "libpslog release pack"
	git -C "$(LUA_RELEASE_STAGE_DIR)" config user.email "release@libpslog.local"
	git -C "$(LUA_RELEASE_STAGE_DIR)" add .
	git -C "$(LUA_RELEASE_STAGE_DIR)" commit -q -m "stage lua release sources"

$(LUA_RELEASE_ROCKSPEC): $(LUA_ROCK_SOURCES) Makefile | $(LUA_DIST_DIR)
	./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "$(LUA_RELEASE_ROCKSPEC)"

$(LUA_RELEASE_PACK_ROCKSPEC): Makefile $(LUA_RELEASE_STAGE_DIR)
	cd "$(LUA_RELEASE_STAGE_DIR)" && ./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "../$(notdir $(LUA_RELEASE_PACK_ROCKSPEC))" "git+file://$(LUA_RELEASE_STAGE_DIR)" ""

$(LUA_RELEASE_ROCK): $(LUA_RELEASE_PACK_ROCKSPEC) $(LUA_RELEASE_ROCKSPEC)
	rm -f "$(LUA_RELEASE_ROCK)"
	cd "$(LUA_RELEASE_STAGE_DIR)" && luarocks pack "../$(notdir $(LUA_RELEASE_PACK_ROCKSPEC))"
	mv "$(LUA_RELEASE_STAGE_DIR)/$(notdir $(LUA_RELEASE_ROCK))" "$(LUA_RELEASE_ROCK)"
	rm -rf "$(LUA_RELEASE_PACK_DIR)"

lua-rock: $(LUA_ROCK_STAMP)

prepare-gobencher-data: build-host

$(LUA_ROCKSPEC): $(LUA_ROCK_SOURCES)
	mkdir -p "$(LUA_ROCK_TREE)"
	./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "$(LUA_ROCKSPEC)" "git+file://$(CURDIR)"

$(LUA_ROCK_STAMP): $(LUA_ROCKSPEC) $(LUA_ROCK_SOURCES)
	flock "$(LUA_ROCK_BUILD_LOCK)" bash -lc 'set -e; luarocks make --tree "$(LUA_ROCK_TREE)" "$(LUA_ROCKSPEC)"; rm -rf $(LUA_ROCK_BUILD_BYPRODUCTS); touch "$(LUA_ROCK_STAMP)"'

lua-test: lua-rock
	eval "$$(luarocks path --tree $(LUA_ROCK_TREE))" && lua lua/tests/test_pslog.lua

$(GO_PRODUCTION_DATASET): $(HOST_GENERATED_VERSION_HEADER) $(GO_PRODUCTION_DATASET_TOOL) $(GO_PRODUCTION_DATASET_SOURCE)
	@tmp_bin="$$(mktemp "$(CURDIR)/.gen_go_production_dataset.XXXXXX")"; \
	rm -f "$$tmp_bin"; \
	cc -std=c99 -O2 -I"$(CURDIR)" -I"$(CURDIR)/include" -I"$(CURDIR)/build/host/generated/include" "$(GO_PRODUCTION_DATASET_TOOL)" "$(GO_PRODUCTION_DATASET_SOURCE)" -o "$$tmp_bin"; \
	"$$tmp_bin" >"$(GO_PRODUCTION_DATASET)"; \
	rm -f "$$tmp_bin"

$(GO_CKVFMT_WRAPPERS): $(GO_PRODUCTION_DATASET) gobencher/cmd/gen_ckvfmt_wrappers/main.go gobencher/benchmark/cpslog_kvfmt.go
	cd gobencher/benchmark && go run ../cmd/gen_ckvfmt_wrappers

release:
	./scripts/run_linux_release_matrix.sh

clean:
	./scripts/clean.sh

clean-dist:
	./scripts/clean.sh --dist-only
