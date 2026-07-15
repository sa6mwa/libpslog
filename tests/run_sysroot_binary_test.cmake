if(NOT DEFINED PSLOG_BINARY_DIR OR NOT DEFINED PSLOG_ROOT OR NOT DEFINED PSLOG_BINARY)
    message(FATAL_ERROR "PSLOG_BINARY_DIR, PSLOG_ROOT, and PSLOG_BINARY are required")
endif()

set(runner "${PSLOG_ROOT}/scripts/run_sysroot_binary.sh")
set(extra_library_path "${PSLOG_BINARY_DIR}/run-sysroot-extra-libs")
file(MAKE_DIRECTORY "${extra_library_path}")

execute_process(
    COMMAND "${runner}" --loader --build-dir "${PSLOG_BINARY_DIR}"
        --library-path "${extra_library_path}"
    RESULT_VARIABLE loader_result
    OUTPUT_VARIABLE loader_output
    ERROR_VARIABLE loader_error
    OUTPUT_STRIP_TRAILING_WHITESPACE
)
if(NOT loader_result EQUAL 0 OR NOT EXISTS "${loader_output}")
    message(FATAL_ERROR
        "sysroot runner did not resolve the native loader with an explicit library path\n"
        "${loader_error}")
endif()

execute_process(
    COMMAND "${runner}" --build-dir "${PSLOG_BINARY_DIR}"
        --library-path "${extra_library_path}" "${PSLOG_BINARY}"
    RESULT_VARIABLE run_result
    OUTPUT_VARIABLE run_output
    ERROR_VARIABLE run_error
)
if(NOT run_result EQUAL 0)
    message(FATAL_ERROR
        "sysroot runner did not execute the native target with an explicit library path\n"
        "stdout:\n${run_output}\n"
        "stderr:\n${run_error}")
endif()
