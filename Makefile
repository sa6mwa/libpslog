SHELL := bash
.DEFAULT_GOAL := help
MAKEFLAGS += --no-builtin-rules
.NOTPARALLEL:

DEBUG_PRESET := debug
HOST_PRESET := host
VALGRIND_PRESET := valgrind
FUZZ_PRESET := fuzz
COVERAGE_PRESET := coverage
RELEASE_BUILD_PRESETS := \
	x86_64-linux-gnu-release \
	x86_64-linux-musl-release \
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
LUA_HOST_INCLUDE_DIR := $(shell luarocks config variables.LUA_INCDIR)
LUA_HOST_LIB_DIR := $(shell luarocks config variables.LUA_LIBDIR)
LUA_STAGED_ROOT := build/lua-host
LUA_STAGED_INCLUDE_DIR := $(LUA_STAGED_ROOT)/include
LUA_STAGED_LIB_DIR := $(LUA_STAGED_ROOT)/lib
LUA_STAGED_LUA_DEPS := $(LUA_STAGED_ROOT)/.staged.stamp
LUA_SDK_PREFIX := $(CURDIR)/build/lua-sdk
LUA_SDK_INCLUDE_DIR := $(LUA_SDK_PREFIX)/include
LUA_SDK_LIB_DIR := $(LUA_SDK_PREFIX)/lib
LUA_SDK_STAMP := $(LUA_SDK_PREFIX)/.installed.stamp
LUA_LOCAL_LIBDIR := $(LUA_SDK_LIB_DIR)
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
HOST_C_COMPILER = $(shell sed -n 's/^CMAKE_C_COMPILER:[^=]*=//p' build/host/CMakeCache.txt | head -n 1)
HOST_CXX_COMPILER = $(shell sed -n 's/^CMAKE_CXX_COMPILER:[^=]*=//p' build/host/CMakeCache.txt | head -n 1)

BENCH_ITERS ?= 200000
FUZZ_TIME ?= 30
FUZZ_LONG_TIME ?= 300
GO_BENCH_COUNT ?= 1
ELEVATORPITCH_ARGS ?=
TIMED := ./scripts/run_timed.sh
RELEASE_TIMING_FILE := $(CURDIR)/build/release-timings.tsv

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
	deps-debug \
	deps-release \
	deps-cross \
	clangd \
	valgrind \
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
	release-pipeline \
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
		'make deps-debug      Provision the pinned native Bootlin collection.' \
		'make deps-release    Provision all pinned Linux release collections.' \
		'make deps-cross      Provision the pinned non-host Linux collections.' \
		'make build-debug     Alias for make build.' \
		'make build-host      Configure and build the host-native local build.' \
		'make build-release   Configure and build the full shipped release build matrix.' \
		'make format          Run clang-format on repo C/C header sources.' \
		'make clangd          Check native debug sources with host clangd.' \
		'make test            Run the debug C test suite.' \
		'make test-debug      Alias for make test.' \
		'make test-host       Build and run the host-native CTest suite.' \
		'make test-all        Run C tests, native hardening, Go gobencher tests, and the Go-vs-C perf gate.' \
		'make valgrind        Run the native Valgrind memory-check gate.' \
		'make coverage        Run the coverage preset and generate coverage-report.' \
		'make fuzz            Run a bounded native AFL++ fuzz job.' \
		'make fuzz-smoke      Run a short fuzz smoke pass.' \
		'make fuzz-long       Run a longer native AFL++ fuzz job (FUZZ_LONG_TIME=$(FUZZ_LONG_TIME)).' \
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
		'make release         Clean, run prerelease proof, and write build/release-timings.tsv.' \
		'make prerelease-hardening Run prerelease plus the long AFL++ fuzz pass.' \
		'make print-release-version Print the version used by release artifacts.' \
		'make release-lua-artifacts Build the standalone Lua release artifacts.' \
		'make lua-rock        Install a local C SDK and build the Lua module into build/luarocks.' \
		'make lua-env         Print shell exports for the repo-local Lua rock.' \
		'make lua-test        Run the Lua binding smoke tests.' \
		'make prepare-gobencher-data Regenerate deterministic Go benchmark inputs.' \
		'make clean           Remove build/ and dist/ generated artifacts.' \
		'make clean-dist      Remove dist/ release artifacts.'

build:
	cmake --preset $(DEBUG_PRESET)
	cmake --build --preset $(DEBUG_PRESET)

deps-debug:
	./scripts/cpkt-toolchains.sh ensure x86_64-linux-gnu

deps-release:
	./scripts/cpkt-toolchains.sh ensure all

deps-cross:
	@set -e; for target in aarch64-linux-gnu aarch64-linux-musl armhf-linux-gnu armhf-linux-musl; do \
		./scripts/cpkt-toolchains.sh ensure "$$target"; \
	done

build-debug: build

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

clangd: build
	./scripts/clangd_check.sh

test: build
	ctest --preset $(DEBUG_PRESET) --output-on-failure

test-debug: test

test-host: build-host
	ctest --preset $(HOST_PRESET) --output-on-failure

gobencher-tests: build-host lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)
	cd gobencher && CC="$(HOST_C_COMPILER)" CXX="$(HOST_CXX_COMPILER)" go test -a ./...

perf-gate: build-host lua-rock
	CC="$(HOST_C_COMPILER)" CXX="$(HOST_CXX_COMPILER)" ./bench/run_perf_gate.sh

test-all: test valgrind fuzz-smoke cross-test gobencher-tests perf-gate

valgrind:
	cmake --preset $(VALGRIND_PRESET)
	cmake --build --preset $(VALGRIND_PRESET)
	valgrind --leak-check=full --track-origins=yes --error-exitcode=1 ./build/$(VALGRIND_PRESET)/pslog_valgrind_facade_tests

coverage:
	cmake --preset $(COVERAGE_PRESET)
	cmake --build --preset $(COVERAGE_PRESET)
	ctest --preset $(COVERAGE_PRESET) --output-on-failure
	cmake --build --preset coverage-report

fuzz:
	./scripts/fuzz.sh run $(FUZZ_TIME)

fuzz-smoke:
	./scripts/fuzz.sh smoke 5

fuzz-long:
	./scripts/fuzz.sh long $(FUZZ_LONG_TIME)

benchmarks-c: build-host
	./build/host/pslog_bench $(BENCH_ITERS) all

benchmarks-gobencher: build-host lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)
	cd gobencher && CC="$(HOST_C_COMPILER)" CXX="$(HOST_CXX_COMPILER)" go test ./benchmark -run '^$$' -bench . -benchmem -count=$(GO_BENCH_COUNT)

benchmarks-go: benchmarks-gobencher

benchmarks-all: benchmarks-c benchmarks-gobencher

benchmarks: benchmarks-all

bench: benchmarks

bench-gate: perf-gate

elevatorpitch: build-host lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)
	export LD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${LD_LIBRARY_PATH:-}" DYLD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${DYLD_LIBRARY_PATH:-}"; eval "$$(luarocks path --tree $(LUA_ROCK_TREE))" && cd gobencher && CC="$(HOST_C_COMPILER)" CXX="$(HOST_CXX_COMPILER)" go run ./cmd/elevatorpitch $(ELEVATORPITCH_ARGS)

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
	cmake --build build/$(HOST_PRESET) --target package-clean-dist
	cmake --build build/$(HOST_PRESET) --target package-archive
	cmake --build build/$(HOST_PRESET) --target package-single-header
	cmake --build build/$(HOST_PRESET) --target package-source
	$(MAKE) release-lua-artifacts
	cmake --build build/$(HOST_PRESET) --target package-checksums

package-source:
	cmake --preset $(HOST_PRESET)
	cmake --build build/$(HOST_PRESET) --target package-source

package-source-smoke:
	cmake --preset $(HOST_PRESET)
	cmake -DPSLOG_ROOT=$(CURDIR) -DPSLOG_BINARY_DIR=$(CURDIR)/build/$(HOST_PRESET) -DPSLOG_VERSION=$(LUA_RELEASE_VERSION) -DPSLOG_TOOLCHAIN_RELATIVE=cmake/toolchains/linux-x86_64-gnu.cmake -P tests/source_archive_smoke_test.cmake

package-single-header:
	cmake --preset $(HOST_PRESET)
	cmake --build build/$(HOST_PRESET) --target package-single-header

package-checksums:
	cmake --preset $(HOST_PRESET)
	cmake --build build/$(HOST_PRESET) --target package-checksums

package-verify:
	cmake --preset $(HOST_PRESET)
	ctest --test-dir build/$(HOST_PRESET) -R '^(package_archives_test|release_privacy_gate_test)$$' --output-on-failure
	cmake -DPSLOG_ROOT=$(CURDIR) -DPSLOG_BINARY_DIR=$(CURDIR)/build/$(HOST_PRESET) -DPSLOG_VERSION=$(LUA_RELEASE_VERSION) -DPSLOG_TOOLCHAIN_RELATIVE=cmake/toolchains/linux-x86_64-gnu.cmake -P tests/source_archive_smoke_test.cmake
	cmake --build build/$(HOST_PRESET) --target package-privacy-gate

verify-release-archives: package-verify

verify-release-privacy:
	cmake --preset $(HOST_PRESET)
	./scripts/verify_release_privacy.sh --build-dir build/$(HOST_PRESET) --target-id x86_64-linux-gnu

release-matrix:
	./scripts/run_linux_release_matrix.sh

finalize-slice: format clangd build-host
	cmake --preset $(HOST_PRESET)
	ctest --test-dir build/$(HOST_PRESET) -R '^(pslog_tests|pslog_single_header_tests|example_integration_test|public_symbol_visibility_test|darwin_linker_route_test)$$' --output-on-failure

release-pipeline:
	$(TIMED) format $(MAKE) format
	$(TIMED) test-all $(MAKE) test-all
	$(TIMED) lua-test $(MAKE) lua-test
	$(TIMED) release-matrix $(MAKE) release-matrix

prerelease: release-pipeline

prerelease-hardening: prerelease fuzz-long

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

prepare-gobencher-data: lua-rock $(GO_PRODUCTION_DATASET) $(GO_CKVFMT_WRAPPERS)

$(LUA_ROCKSPEC): $(LUA_ROCK_SOURCES)
	mkdir -p "$(LUA_ROCK_TREE)"
	./lua/scripts/render_release_rockspec.sh "$(LUA_RELEASE_VERSION)" "$(LUA_ROCKSPEC)" "git+file://$(CURDIR)"

$(LUA_STAGED_LUA_DEPS):
	@test -n "$(LUA_HOST_INCLUDE_DIR)" || { echo 'Lua 5.5 development headers are required through LuaRocks'; exit 1; }
	@test -n "$(LUA_HOST_LIB_DIR)" || { echo 'Lua 5.5 development library is required through LuaRocks'; exit 1; }
	test -f "$(LUA_HOST_INCLUDE_DIR)/lua.h" || { echo "missing Lua 5.5 header: $(LUA_HOST_INCLUDE_DIR)/lua.h"; exit 1; }
	test -f "$(LUA_HOST_LIB_DIR)/liblua.a" || { echo "missing Lua 5.5 library: $(LUA_HOST_LIB_DIR)/liblua.a"; exit 1; }
	rm -rf "$(LUA_STAGED_ROOT)"
	mkdir -p "$(LUA_STAGED_INCLUDE_DIR)" "$(LUA_STAGED_LIB_DIR)"
	cp -a "$(LUA_HOST_INCLUDE_DIR)/." "$(LUA_STAGED_INCLUDE_DIR)/"
	cp "$(LUA_HOST_LIB_DIR)/liblua.a" "$(LUA_STAGED_LIB_DIR)/"
	touch "$@"

$(LUA_SDK_STAMP): build-host
	rm -rf "$(LUA_SDK_PREFIX)"
	cmake --install build/$(HOST_PRESET) --prefix "$(LUA_SDK_PREFIX)"
	test -f "$(LUA_SDK_INCLUDE_DIR)/pslog.h"
	test -f "$(LUA_SDK_INCLUDE_DIR)/pslog_version.h"
	test -e "$(LUA_SDK_LIB_DIR)/libpslog.so"
	touch "$@"

$(LUA_ROCK_STAMP): $(LUA_ROCKSPEC) $(LUA_ROCK_SOURCES) $(LUA_STAGED_LUA_DEPS) $(LUA_SDK_STAMP)
	flock "$(LUA_ROCK_BUILD_LOCK)" env CC="$(HOST_C_COMPILER)" CXX="$(HOST_CXX_COMPILER)" bash -lc 'set -e; luarocks make --tree "$(LUA_ROCK_TREE)" "$(LUA_ROCKSPEC)" LUA_INCDIR="$(CURDIR)/$(LUA_STAGED_INCLUDE_DIR)" LIBPSLOG_INCDIR="$(LUA_SDK_INCLUDE_DIR)" LIBPSLOG_LIBDIR="$(LUA_SDK_LIB_DIR)"; rm -rf $(LUA_ROCK_BUILD_BYPRODUCTS); touch "$(LUA_ROCK_STAMP)"'

lua-test: lua-rock
	LD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${LD_LIBRARY_PATH:-}" DYLD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${DYLD_LIBRARY_PATH:-}" ./lua/scripts/check_binding_boundary.sh "$(LUA_ROCK_TREE)"
	CC="$(HOST_C_COMPILER)" CXX="$(HOST_CXX_COMPILER)" ./lua/scripts/run_interop_embedder_test.sh "$(CURDIR)/build/$(HOST_PRESET)" "$(LUA_RELEASE_VERSION)" "$(LUA_ROCK_TREE)" "$(LUA_SDK_PREFIX)"
	export LD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${LD_LIBRARY_PATH:-}" DYLD_LIBRARY_PATH="$(LUA_LOCAL_LIBDIR):$${DYLD_LIBRARY_PATH:-}"; eval "$$(luarocks path --tree $(LUA_ROCK_TREE))" && lua lua/tests/test_pslog.lua

$(GO_PRODUCTION_DATASET): $(HOST_GENERATED_VERSION_HEADER) $(GO_PRODUCTION_DATASET_TOOL) $(GO_PRODUCTION_DATASET_SOURCE)
	@tmp_bin="$$(mktemp "$(CURDIR)/.gen_go_production_dataset.XXXXXX")"; tmp_output="$$(mktemp "$(CURDIR)/.gen_go_production_dataset_output.XXXXXX")"; \
	trap 'rm -f "$$tmp_bin" "$$tmp_output"' EXIT; \
	rm -f "$$tmp_bin"; \
	"$(HOST_C_COMPILER)" -std=c99 -O2 -I"$(CURDIR)" -I"$(CURDIR)/include" -I"$(CURDIR)/build/host/generated/include" "$(GO_PRODUCTION_DATASET_TOOL)" "$(GO_PRODUCTION_DATASET_SOURCE)" -o "$$tmp_bin"; \
	"$$tmp_bin" >"$$tmp_output"; \
	mv "$$tmp_output" "$(GO_PRODUCTION_DATASET)"

$(GO_CKVFMT_WRAPPERS): $(GO_PRODUCTION_DATASET) $(LUA_SDK_STAMP) gobencher/cmd/gen_ckvfmt_wrappers/main.go gobencher/benchmark/cpslog_kvfmt.go
	cd gobencher/benchmark && CC="$(HOST_C_COMPILER)" CXX="$(HOST_CXX_COMPILER)" go run ../cmd/gen_ckvfmt_wrappers

release:
	PKT_TIMING_FILE="$(RELEASE_TIMING_FILE)" $(TIMED) release-clean $(MAKE) clean
	PKT_TIMING_FILE="$(RELEASE_TIMING_FILE)" $(TIMED) release-pipeline $(MAKE) release-pipeline
	@printf 'Release timings: %s\n' "$(RELEASE_TIMING_FILE)"

clean:
	./scripts/clean.sh

clean-dist:
	./scripts/clean.sh --dist-only
