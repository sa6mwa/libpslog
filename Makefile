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
	armhf-linux-musl-release \
	$(shell osxcross_bin="$${OSXCROSS_ROOT:-$$HOME/.local/cross/osxcross}/bin"; for cc in "$$osxcross_bin"/arm64-apple-darwin*-clang; do if [ -x "$$cc" ]; then printf '%s' arm64-apple-darwin-release; break; fi; done)

CROSS_RELEASE_PRESETS := \
	aarch64-linux-gnu-release \
	aarch64-linux-musl-release \
	armhf-linux-gnu-release \
	armhf-linux-musl-release
LUA_RELEASE_VERSION := $(shell ./lua/scripts/release_version.sh)
LUA_DIST_DIR := $(CURDIR)/dist
LUA_RELEASE_ROCKSPEC := $(LUA_DIST_DIR)/lua-pslog-$(LUA_RELEASE_VERSION)-1.rockspec
LUA_RELEASE_PACK_DIR := $(LUA_DIST_DIR)/.lua-pack
LUA_RELEASE_STAGE_NAME := lua-pslog-$(LUA_RELEASE_VERSION)
LUA_RELEASE_STAGE_DIR := $(LUA_RELEASE_PACK_DIR)/$(LUA_RELEASE_STAGE_NAME)
LUA_RELEASE_PACK_SOURCE_TAR := $(LUA_RELEASE_PACK_DIR)/$(LUA_RELEASE_STAGE_NAME).tar
LUA_RELEASE_PACK_SOURCE_TARBALL := $(LUA_RELEASE_PACK_SOURCE_TAR).gz
LUA_RELEASE_SOURCE_TARBALL := $(LUA_DIST_DIR)/$(LUA_RELEASE_STAGE_NAME).tar.gz
LUA_RELEASE_PACK_ROCKSPEC := $(LUA_RELEASE_PACK_DIR)/lua-pslog-$(LUA_RELEASE_VERSION)-1.rockspec
LUA_RELEASE_ROCK := $(LUA_DIST_DIR)/lua-pslog-$(LUA_RELEASE_VERSION)-1.src.rock
LUA_ROCK_TREE := build/luarocks
LUA_ROCKSPEC := $(LUA_ROCK_TREE)/lua-pslog-$(LUA_RELEASE_VERSION)-1.rockspec
LUA_ROCK_STAMP := $(LUA_ROCK_TREE)/.installed.stamp
LUA_ROCK_BUILD_LOCK := $(LUA_ROCK_TREE)/.build.lock
LUA_LOCAL_LIBDIR := $(CURDIR)/build/host
LUA_ROCK_BUILD_BYPRODUCTS := \
	$(CURDIR)/pslog \
	$(CURDIR)/lua/src/pslog_lua.o
LUA_ROCK_SOURCES := \
	lua/lua-pslog.rockspec.in \
	lua/scripts/check_binding_boundary.sh \
	lua/scripts/render_release_rockspec.sh \
	lua/scripts/release_version.sh \
	lua/pslog/init.lua \
	lua/src/pslog_lua.c \
	include/pslog_lua.h \
	include/pslog.h
GO_PRODUCTION_DATASET := gobencher/benchmark/production_data_generated.go
GO_PRODUCTION_DATASET_SOURCE := bench/bench_production_dataset.c
GO_PRODUCTION_DATASET_TOOL := bench/gen_go_production_dataset.c
GO_CKVFMT_WRAPPERS := gobencher/benchmark/cpslog_kvfmt_generated.go
HOST_GENERATED_VERSION_HEADER := $(CURDIR)/build/host/generated/include/pslog_version.h

BENCH_ITERS ?= 200000
FUZZ_TIME ?= 30
GO_BENCH_COUNT ?= 1
ELEVATORPITCH_ARGS ?=

.PHONY: \
	help \
	build \
	build-debug \
	build-host \
	build-release \
	format \
	test \
	test-debug \
	test-host \
	test-all \
	asan \
	coverage \
	fuzz \
	fuzz-smoke \
	fuzz-long \
	bench \
	bench-gate \
	benchmarks-c \
	benchmarks-gobencher \
	benchmarks-go \
	benchmarks-all \
	benchmarks \
	elevatorpitch \
	gobencher-tests \
	prepare-gobencher-data \
	perf-gate \
	cross-build \
	cross-test \
	test-cross \
	package \
	package-source \
	package-source-smoke \
	package-single-header \
	package-checksums \
	package-verify \
	verify-release-archives \
	verify-release-privacy \
	release-matrix \
	finalize-slice \
	prerelease \
	prerelease-hardening \
	print-release-version \
	release-lua-artifacts \
	lua-env \
	lua-rock \
	lua-test \
	release \
	clean \
	clean-dist

help:
	@printf '%s\n' \
		'make build           Configure and build the debug preset.' \
		'make build-debug     Alias for make build.' \
		'make build-host      Configure and build the host-native local build.' \
		'make build-release   Configure and build the full shipped release build matrix.' \
		'make format          Run clang-format on repo C/C header sources.' \
		'make test            Run the debug C test suite.' \
		'make test-debug      Alias for make test.' \
		'make test-host       Build and run the host-native CTest suite.' \
		'make test-all        Run C tests, Go gobencher tests, and the Go-vs-C perf gate.' \
		'make asan            Run the ASan/UBSan preset test suite.' \
		'make coverage        Run the coverage preset and generate coverage-report.' \
		'make fuzz            Build fuzz target and run two bounded fuzz passes.' \
		'make fuzz-smoke      Run a short fuzz smoke pass.' \
		'make fuzz-long       Run a longer opt-in fuzz pass (FUZZ_TIME=$(FUZZ_TIME)).' \
		'make bench           Alias for make benchmarks.' \
		'make bench-gate      Alias for make perf-gate.' \
		'make benchmarks-c    Run the pure C benchmark matrix (BENCH_ITERS=$(BENCH_ITERS)).' \
		'make benchmarks-gobencher Run the full gobencher benchmark suite.' \
		'make benchmarks-all  Run both C and gobencher benchmark suites.' \
		'make benchmarks      Alias for make benchmarks-all.' \
		'make elevatorpitch   Run the live Go/C/Lua comparison chart.' \
		'make gobencher-tests Run all Go gobencher tests.' \
		'make perf-gate       Run the Go-vs-C performance gate.' \
		'make cross-build     Build all non-host cross release presets.' \
		'make cross-test      Test all non-host cross release presets.' \
		'make test-cross      Alias for make cross-test.' \
		'make package         Build host C package artifacts, Lua artifacts, and checksums.' \
		'make package-source  Build the source release archive.' \
		'make package-source-smoke Verify the source archive builds outside git.' \
		'make package-single-header Build the single-header release artifact.' \
		'make package-checksums Generate the release checksum manifest.' \
		'make package-verify  Run package archive and release privacy verification.' \
		'make verify-release-archives Alias for make package-verify.' \
		'make verify-release-privacy Run checksum-manifest privacy/relocatability gate.' \
		'make release-matrix  Build/test/package the release target matrix.' \
		'make finalize-slice  Run formatting and focused local verification.' \
		'make prerelease      Run deterministic local pre-release checks.' \
		'make prerelease-hardening Run prerelease plus release matrix.' \
		'make print-release-version Print the version used by release artifacts.' \
		'make lua-rock        Build and install the Lua module into build/luarocks.' \
		'make lua-env         Print shell exports for the repo-local Lua rock.' \
		'make lua-test        Run the Lua binding smoke tests.' \
		'make release         Clean generated state, then run the full release matrix.' \
		'make clean           Remove build/ and dist/ generated artifacts.' \
		'make clean-dist      Remove dist/ release artifacts.'

build:
	cmake --preset $(DEBUG_PRESET)
	cmake --build --preset $(DEBUG_PRESET)

build-debug: build

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

test-debug: test

test-host: build-host
	ctest --preset $(HOST_PRESET) --output-on-failure

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

fuzz-smoke:
	cmake --preset $(FUZZ_PRESET)
	cmake --build --preset $(FUZZ_PRESET)
	./build/fuzz/pslog_fuzz -max_total_time=5 -max_len=256

fuzz-long:
	cmake --preset $(FUZZ_PRESET)
	cmake --build --preset $(FUZZ_PRESET)
	./build/fuzz/pslog_fuzz -max_total_time=$(FUZZ_TIME)

benchmarks-c: build-host
	./build/host/pslog_bench $(BENCH_ITERS) all

benchmarks-gobencher: build-host lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)
	cd gobencher && go test ./benchmark -run '^$$' -bench . -benchmem -count=$(GO_BENCH_COUNT)

benchmarks-go: benchmarks-gobencher

benchmarks-all: benchmarks-c benchmarks-gobencher

benchmarks: benchmarks-all

bench: benchmarks

bench-gate: perf-gate

elevatorpitch: build-host lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)
	export LD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${LD_LIBRARY_PATH:-}" DYLD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${DYLD_LIBRARY_PATH:-}"; eval "$$(luarocks path --tree $(LUA_ROCK_TREE))" && cd gobencher && go run ./cmd/elevatorpitch $(ELEVATORPITCH_ARGS)

cross-build:
	@set -e; for preset in $(CROSS_RELEASE_PRESETS); do \
		cmake --preset "$$preset"; \
		cmake --build --preset "$$preset"; \
	done

cross-test: cross-build
	@set -e; for preset in $(CROSS_RELEASE_PRESETS); do \
		ctest --preset "$$preset" --output-on-failure; \
	done

test-cross: cross-test

package:
	cmake --preset $(HOST_PRESET)
	cmake --build --preset package-archive
	cmake --build --preset package-single-header
	cmake --build --preset package-source
	$(MAKE) release-lua-artifacts
	cmake --build --preset package-checksums

package-source:
	cmake --preset $(HOST_PRESET)
	cmake --build --preset package-source

package-source-smoke:
	cmake --preset $(HOST_PRESET)
	cmake -DPSLOG_ROOT=$(CURDIR) -DPSLOG_BINARY_DIR=$(CURDIR)/build/$(HOST_PRESET) -DPSLOG_VERSION=$(LUA_RELEASE_VERSION) -P tests/source_archive_smoke_test.cmake

package-single-header:
	cmake --preset $(HOST_PRESET)
	cmake --build --preset package-single-header

package-checksums:
	cmake --preset $(HOST_PRESET)
	cmake --build --preset package-checksums

package-verify:
	cmake --preset $(HOST_PRESET)
	ctest --test-dir build/$(HOST_PRESET) -R '^(package_archives_test|release_privacy_gate_test)$$' --output-on-failure
	cmake -DPSLOG_ROOT=$(CURDIR) -DPSLOG_BINARY_DIR=$(CURDIR)/build/$(HOST_PRESET) -DPSLOG_VERSION=$(LUA_RELEASE_VERSION) -P tests/source_archive_smoke_test.cmake
	cmake --build --preset package-privacy-gate

verify-release-archives: package-verify

verify-release-privacy:
	cmake --preset $(HOST_PRESET)
	cmake --build --preset package-privacy-gate

release-matrix:
	./scripts/run_linux_release_matrix.sh

finalize-slice: format build-host
	cmake --preset $(HOST_PRESET)
	ctest --test-dir build/$(HOST_PRESET) -R '^(pslog_tests|pslog_single_header_tests|example_integration_test|public_symbol_visibility_test|darwin_linker_route_test)$$' --output-on-failure

prerelease: format test asan lua-test fuzz-smoke package-verify

prerelease-hardening: prerelease gobencher-tests perf-gate fuzz release-matrix

print-release-version:
	@./lua/scripts/release_version.sh

release-lua-artifacts: $(LUA_RELEASE_ROCKSPEC) $(LUA_RELEASE_SOURCE_TARBALL) $(LUA_RELEASE_ROCK)

$(LUA_RELEASE_PACK_DIR):
	mkdir -p "$(LUA_RELEASE_PACK_DIR)"

$(LUA_DIST_DIR):
	mkdir -p "$(LUA_DIST_DIR)"

$(LUA_RELEASE_STAGE_DIR): Makefile .git/index | $(LUA_RELEASE_PACK_DIR)
	rm -rf "$(LUA_RELEASE_STAGE_DIR)"
	mkdir -p "$(LUA_RELEASE_STAGE_DIR)"
	git archive --format=tar --worktree-attributes HEAD | tar -xf - -C "$(LUA_RELEASE_STAGE_DIR)"
	rm -f "$(LUA_RELEASE_STAGE_DIR)/REVIEW.md"
	printf '%s\n' "$(LUA_RELEASE_VERSION)" >"$(LUA_RELEASE_STAGE_DIR)/VERSION"
	{ cd "$(LUA_RELEASE_STAGE_DIR)" && find . -type f ! -name VERSION ! -name RELEASE_MANIFEST | sed 's#^\./##' | grep -v '^REVIEW\.md$$' | LC_ALL=C sort; printf '%s\n' VERSION RELEASE_MANIFEST; } >"$(LUA_RELEASE_STAGE_DIR)/RELEASE_MANIFEST"

$(LUA_RELEASE_ROCKSPEC): $(LUA_ROCK_SOURCES) Makefile | $(LUA_DIST_DIR)
	./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "$(LUA_RELEASE_ROCKSPEC)"

$(LUA_RELEASE_PACK_ROCKSPEC): Makefile $(LUA_RELEASE_STAGE_DIR)
	cd "$(LUA_RELEASE_STAGE_DIR)" && ./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "../$(notdir $(LUA_RELEASE_PACK_ROCKSPEC))" "file://$(LUA_RELEASE_PACK_SOURCE_TARBALL)" ""

$(LUA_RELEASE_SOURCE_TARBALL): $(LUA_RELEASE_STAGE_DIR) | $(LUA_DIST_DIR)
	rm -f "$(LUA_RELEASE_PACK_SOURCE_TAR)" "$(LUA_RELEASE_PACK_SOURCE_TARBALL)"
	cd "$(LUA_RELEASE_PACK_DIR)" && tar -cf "$(notdir $(LUA_RELEASE_PACK_SOURCE_TAR))" "$(LUA_RELEASE_STAGE_NAME)"
	gzip -9 -f "$(LUA_RELEASE_PACK_SOURCE_TAR)"
	cp "$(LUA_RELEASE_PACK_SOURCE_TARBALL)" "$(LUA_RELEASE_SOURCE_TARBALL)"

$(LUA_RELEASE_ROCK): $(LUA_RELEASE_PACK_ROCKSPEC) $(LUA_RELEASE_ROCKSPEC) $(LUA_RELEASE_SOURCE_TARBALL)
	rm -f "$(LUA_RELEASE_ROCK)"
	cd "$(LUA_RELEASE_PACK_DIR)" && luarocks pack "$(notdir $(LUA_RELEASE_PACK_ROCKSPEC))"
	mv "$(LUA_RELEASE_PACK_DIR)/$(notdir $(LUA_RELEASE_ROCK))" "$(LUA_RELEASE_ROCK)"
	@tmp_dir="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "$$tmp_dir/$(notdir $(LUA_RELEASE_PACK_ROCKSPEC))" "file://$(notdir $(LUA_RELEASE_PACK_SOURCE_TARBALL))" ""; \
	cd "$$tmp_dir" && zip -q -u "$(LUA_RELEASE_ROCK)" "$(notdir $(LUA_RELEASE_PACK_ROCKSPEC))"
	rm -rf "$(LUA_RELEASE_PACK_DIR)"

lua-rock: $(LUA_ROCK_STAMP)

lua-env:
	@printf '%s\n' 'eval "$$(luarocks path --tree $(LUA_ROCK_TREE))"'
	@printf '%s\n' 'export LD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${LD_LIBRARY_PATH:-}"'
	@printf '%s\n' 'export DYLD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${DYLD_LIBRARY_PATH:-}"'

prepare-gobencher-data: build-host

$(LUA_ROCKSPEC): $(LUA_ROCK_SOURCES)
	mkdir -p "$(LUA_ROCK_TREE)"
	./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "$(LUA_ROCKSPEC)" "git+file://$(CURDIR)"

$(LUA_ROCK_STAMP): $(LUA_ROCKSPEC) $(LUA_ROCK_SOURCES) build-host
	flock "$(LUA_ROCK_BUILD_LOCK)" bash -lc 'set -e; luarocks make --tree "$(LUA_ROCK_TREE)" "$(LUA_ROCKSPEC)" LIBPSLOG_INCDIR="$(CURDIR)/include" LIBPSLOG_LIBDIR="$(CURDIR)/build/host"; rm -rf $(LUA_ROCK_BUILD_BYPRODUCTS); touch "$(LUA_ROCK_STAMP)"'

lua-test: lua-rock
	LD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${LD_LIBRARY_PATH:-}" DYLD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${DYLD_LIBRARY_PATH:-}" ./lua/scripts/check_binding_boundary.sh "$(LUA_ROCK_TREE)"
	./lua/scripts/run_interop_embedder_test.sh "$(LUA_LOCAL_LIBDIR)" "$(LUA_RELEASE_VERSION)" "$(LUA_ROCK_TREE)"
	export LD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${LD_LIBRARY_PATH:-}" DYLD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${DYLD_LIBRARY_PATH:-}"; eval "$$(luarocks path --tree $(LUA_ROCK_TREE))" && lua lua/tests/test_pslog.lua

$(GO_PRODUCTION_DATASET): $(HOST_GENERATED_VERSION_HEADER) $(GO_PRODUCTION_DATASET_TOOL) $(GO_PRODUCTION_DATASET_SOURCE)
	@tmp_bin="$$(mktemp "$(CURDIR)/.gen_go_production_dataset.XXXXXX")"; \
	rm -f "$$tmp_bin"; \
	cc -std=c99 -O2 -I"$(CURDIR)" -I"$(CURDIR)/include" -I"$(CURDIR)/build/host/generated/include" "$(GO_PRODUCTION_DATASET_TOOL)" "$(GO_PRODUCTION_DATASET_SOURCE)" -o "$$tmp_bin"; \
	"$$tmp_bin" >"$(GO_PRODUCTION_DATASET)"; \
	rm -f "$$tmp_bin"

$(GO_CKVFMT_WRAPPERS): $(GO_PRODUCTION_DATASET) gobencher/cmd/gen_ckvfmt_wrappers/main.go gobencher/benchmark/cpslog_kvfmt.go
	cd gobencher/benchmark && go run ../cmd/gen_ckvfmt_wrappers

release: clean release-matrix

clean:
	./scripts/clean.sh

clean-dist:
	./scripts/clean.sh --dist-only
