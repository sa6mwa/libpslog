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
if(NOT release_output MATCHES "(^|\n)release: clean release-pipeline(\n|$)")
    message(FATAL_ERROR "release must clean before invoking release-pipeline")
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
