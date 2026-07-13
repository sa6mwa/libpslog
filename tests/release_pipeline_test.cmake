if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

function(pslog_make_database target output_var)
    execute_process(
        COMMAND make -pn "${target}"
        WORKING_DIRECTORY "${PSLOG_ROOT}"
        RESULT_VARIABLE result
        OUTPUT_VARIABLE output
        ERROR_VARIABLE error
    )
    if(NOT result EQUAL 0)
        message(FATAL_ERROR "make -pn ${target} failed:\n${output}\n${error}")
    endif()
    set(${output_var} "${output}" PARENT_SCOPE)
endfunction()

pslog_make_database(prerelease prerelease_output)
if(NOT prerelease_output MATCHES "(^|\n)prerelease: release-pipeline(\n|$)")
    message(FATAL_ERROR "prerelease does not invoke the shared release-pipeline")
endif()

pslog_make_database(release release_output)
string(FIND "${release_output}" "$(MAKE) clean" release_clean_recipe_offset)
string(FIND "${release_output}" "$(MAKE) release-pipeline" release_pipeline_recipe_offset)
if(release_clean_recipe_offset LESS 0 OR release_pipeline_recipe_offset LESS release_clean_recipe_offset)
    message(FATAL_ERROR "release must invoke clean before the shared release-pipeline")
endif()

execute_process(
    COMMAND make -n release-pipeline
    WORKING_DIRECTORY "${PSLOG_ROOT}"
    RESULT_VARIABLE test_all_result
    OUTPUT_VARIABLE test_all_output
    ERROR_VARIABLE test_all_error
)
if(NOT test_all_result EQUAL 0)
    message(FATAL_ERROR "make -n release-pipeline failed:\n${test_all_output}\n${test_all_error}")
endif()
if(NOT test_all_output MATCHES "--preset format" OR
   NOT test_all_output MATCHES "aarch64-linux-gnu-release" OR
   NOT test_all_output MATCHES "ctest --preset")
    message(FATAL_ERROR "release-pipeline does not include formatting and the cross-test command graph")
endif()
string(FIND "${test_all_output}" "make format" format_offset)
string(FIND "${test_all_output}" "make test-all" ordinary_offset)
string(FIND "${test_all_output}" "make release-matrix" matrix_offset)
if(format_offset LESS 0 OR ordinary_offset LESS format_offset OR matrix_offset LESS ordinary_offset)
    message(FATAL_ERROR "release-pipeline does not run ordinary checks before release-matrix")
endif()

execute_process(
    COMMAND make -n release
    WORKING_DIRECTORY "${PSLOG_ROOT}"
    RESULT_VARIABLE release_dry_run_result
    OUTPUT_VARIABLE release_dry_run_output
    ERROR_VARIABLE release_dry_run_error
)
if(NOT release_dry_run_result EQUAL 0)
    message(FATAL_ERROR "make -n release failed:\n${release_dry_run_output}\n${release_dry_run_error}")
endif()
string(FIND "${release_dry_run_output}" "./scripts/clean.sh" release_clean_offset)
string(FIND "${release_dry_run_output}" "make format" release_format_offset)
if(release_clean_offset LESS 0 OR release_format_offset LESS release_clean_offset)
    message(FATAL_ERROR "release does not clean before entering the shared release pipeline")
endif()
