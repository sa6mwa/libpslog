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
