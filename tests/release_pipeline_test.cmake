if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

# This is deliberately structural: make -n still executes recursive $(MAKE)
# recipes, which makes a clean-tree release-contract test materialize outputs.
file(READ "${PSLOG_ROOT}/Makefile" makefile)

if(NOT makefile MATCHES "prerelease: release-pipeline")
    message(FATAL_ERROR "prerelease does not invoke the shared release-pipeline")
endif()
if(NOT makefile MATCHES "release:\n[ \t]*PKT_TIMING_FILE=.*release-clean \\$\\(MAKE\\) clean\n[ \t]*PKT_TIMING_FILE=.*release-pipeline \\$\\(MAKE\\) release-pipeline")
    message(FATAL_ERROR "release must clean before invoking the shared release-pipeline")
endif()

string(FIND "${makefile}" "release-pipeline:" pipeline_offset)
string(FIND "${makefile}" "$(MAKE) format" format_offset)
string(FIND "${makefile}" "$(MAKE) test-all" ordinary_offset)
string(FIND "${makefile}" "$(MAKE) release-matrix" matrix_offset)
if(pipeline_offset LESS 0 OR format_offset LESS pipeline_offset OR
   ordinary_offset LESS format_offset OR matrix_offset LESS ordinary_offset)
    message(FATAL_ERROR "release-pipeline does not run ordinary checks before release-matrix")
endif()

if(NOT makefile MATCHES "prerelease-hardening: prerelease fuzz-long")
    message(FATAL_ERROR "prerelease-hardening must include the explicit long AFL++ tier")
endif()

if(NOT makefile MATCHES "RELEASE_TIMING_FILE := \\$\\(CURDIR\\)/build/release-timings\\.tsv" OR
   NOT makefile MATCHES "TIMED := ./scripts/run_timed\\.sh")
    message(FATAL_ERROR "release must expose the lifecycle timing output through scripts/run_timed.sh")
endif()
if(NOT makefile MATCHES "test-all:\n[ \t]*\\$\\(TIMED\\) test \\$\\(MAKE\\) test\n[ \t]*\\$\\(TIMED\\) valgrind \\$\\(MAKE\\) valgrind" OR
   NOT makefile MATCHES "\\$\\(TIMED\\) fuzz-smoke \\$\\(MAKE\\) fuzz-smoke" OR
   NOT makefile MATCHES "\\$\\(TIMED\\) cross-test \\$\\(MAKE\\) cross-test" OR
   NOT makefile MATCHES "\\$\\(TIMED\\) gobencher-tests \\$\\(MAKE\\) gobencher-tests" OR
   NOT makefile MATCHES "\\$\\(TIMED\\) perf-gate \\$\\(MAKE\\) perf-gate")
    message(FATAL_ERROR "test-all must expose its serial release timing phases")
endif()
if(NOT makefile MATCHES "HOST_BINARY_RUNNER := \\$\\(CURDIR\\)/scripts/run_host_binary\\.sh" OR
   NOT makefile MATCHES "go test -exec \"\\$\\(HOST_BINARY_RUNNER\\)\"" OR
   NOT makefile MATCHES "PSLOG_HOST_EXECUTOR=\"\\$\\(HOST_BINARY_RUNNER\\)\"")
    message(FATAL_ERROR "Bootlin-linked Go and performance gates must run through the configured host sysroot runner")
endif()
if(NOT makefile MATCHES "LUA_HOST_INTERPRETER = \\$\\(shell \\$\\(LUA_ROCKS\\) config variables\\.LUA" OR
   NOT makefile MATCHES "LUA_STAGED_IDENTITY := \\$\\(LUA_STAGED_ROOT\\)/\\.luarocks-identity" OR
   NOT makefile MATCHES "\\$\\(LUA_STAGED_LUA_DEPS\\): \\$\\(LUA_STAGED_IDENTITY\\)" OR
   NOT makefile MATCHES "\\./scripts/sha256_files\\.sh \"\\$\\(LUA_HOST_INCLUDE_DIR\\)/lua\\.h\" \"\\$\\(LUA_HOST_LIB_DIR\\)/liblua\\.a\"" OR
   NOT makefile MATCHES "\\./scripts/with_lock\\.sh \"\\$\\(LUA_ROCK_BUILD_LOCK\\)\" env CC=" OR
   NOT makefile MATCHES "\\\"\\$\\(LUA_HOST_INTERPRETER\\)\\\" lua/tests/test_pslog\\.lua")
    message(FATAL_ERROR "Lua gates must track the selected LuaRocks ABI inputs and run with the selected Lua interpreter")
endif()
if(makefile MATCHES "(^|\n)[^\n]*flock \"\\$\\(LUA_ROCK_BUILD_LOCK\\)\"" OR
   makefile MATCHES "(^|\n)[^\n]*sha256sum \"\\$\\(LUA_HOST_INCLUDE_DIR\\)/lua\\.h\"")
    message(FATAL_ERROR "Lua gates must avoid GNU/Linux-only flock and sha256sum commands")
endif()
file(READ "${PSLOG_ROOT}/scripts/verify_release_privacy.sh" privacy_gate)
if(NOT privacy_gate MATCHES "scripts/sha256_files\\.sh\" --check" OR
   privacy_gate MATCHES "(^|\n)[^\n]*sha256sum -c")
    message(FATAL_ERROR "release privacy checksum verification must use the portable checksum helper")
endif()
file(READ "${PSLOG_ROOT}/tests/package_archives_test.cmake" package_archives_test)
if(NOT package_archives_test MATCHES "scripts/sha256_files\\.sh\" --check" OR
   package_archives_test MATCHES "(^|\n)[^\n]*sha256sum -c")
    message(FATAL_ERROR "package archive checksum verification must use the portable checksum helper")
endif()
if(NOT makefile MATCHES "loader=\\\"\\$\\$\\(\\./scripts/run_sysroot_binary\\.sh --loader --build-dir" OR
   NOT makefile MATCHES "valgrind --leak-check=full --track-origins=yes --error-exitcode=1" OR
   NOT makefile MATCHES "\\\"\\$\\$loader\\\" --library-path")
    message(FATAL_ERROR "Valgrind must run the selected Bootlin facade process directly through its sysroot loader")
endif()
