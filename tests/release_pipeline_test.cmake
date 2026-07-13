if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

# This is deliberately structural: make -n still executes recursive $(MAKE)
# recipes, which makes a clean-tree release-contract test materialize outputs.
file(READ "${PSLOG_ROOT}/Makefile" makefile)

if(NOT makefile MATCHES "prerelease: release-pipeline")
    message(FATAL_ERROR "prerelease does not invoke the shared release-pipeline")
endif()
if(NOT makefile MATCHES "release:\n[ \t]*\\$\\(MAKE\\) clean\n[ \t]*\\$\\(MAKE\\) release-pipeline")
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
