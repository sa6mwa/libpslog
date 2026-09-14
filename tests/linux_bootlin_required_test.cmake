foreach(required IN ITEMS PSLOG_BINARY_DIR PSLOG_ROOT)
    if(NOT DEFINED ${required} OR "${${required}}" STREQUAL "")
        message(FATAL_ERROR "${required} is required")
    endif()
endforeach()

set(test_root "${PSLOG_BINARY_DIR}/linux-bootlin-required-test")
file(REMOVE_RECURSE "${test_root}")
execute_process(
    COMMAND "${CMAKE_COMMAND}" -S "${PSLOG_ROOT}" -B "${test_root}"
        -DPSLOG_BUILD_TESTS=OFF
        -DPSLOG_BUILD_BENCHMARKS=OFF
    RESULT_VARIABLE configure_result
    OUTPUT_VARIABLE configure_output
    ERROR_VARIABLE configure_error
)
if(configure_result EQUAL 0)
    message(FATAL_ERROR "Linux configure unexpectedly succeeded without a Bootlin toolchain")
endif()
set(configure_log "${configure_output}${configure_error}")
if(NOT configure_log MATCHES "before compiler discovery")
    message(FATAL_ERROR "Linux configure did not explain the missing Bootlin toolchain:\n${configure_log}")
endif()
if(configure_log MATCHES "C compiler identification")
    message(FATAL_ERROR "Linux configure selected a compiler before rejecting the missing Bootlin toolchain")
endif()
