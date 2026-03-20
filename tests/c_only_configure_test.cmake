if(NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_BINARY_DIR is required")
endif()

if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

if(NOT DEFINED PSLOG_C_COMPILER OR PSLOG_C_COMPILER STREQUAL "")
    message(FATAL_ERROR "PSLOG_C_COMPILER is required")
endif()

set(test_root "${PSLOG_BINARY_DIR}/c-only-configure-test")
set(build_dir "${test_root}/build")
set(fake_cxx "${test_root}/definitely-no-cxx")

file(REMOVE_RECURSE "${test_root}")
file(MAKE_DIRECTORY "${test_root}")

execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -S "${PSLOG_ROOT}"
        -B "${build_dir}"
        -DPSLOG_BUILD_TESTS=OFF
        -DPSLOG_BUILD_BENCHMARKS=OFF
        -DPSLOG_BUILD_FUZZ=OFF
        -DPSLOG_BENCHMARK_WITH_LIBLOGGER=OFF
        -DPSLOG_BENCHMARK_WITH_QUILL=OFF
        -DCMAKE_C_COMPILER=${PSLOG_C_COMPILER}
        -DCMAKE_CXX_COMPILER=${fake_cxx}
    RESULT_VARIABLE configure_result
    OUTPUT_VARIABLE configure_stdout
    ERROR_VARIABLE configure_stderr
)
if(NOT configure_result EQUAL 0)
    message(FATAL_ERROR
        "expected a plain C-only configure to succeed without a working C++ compiler\n"
        "stdout:\n${configure_stdout}\n"
        "stderr:\n${configure_stderr}")
endif()
