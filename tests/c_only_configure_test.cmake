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
        -DCMAKE_TOOLCHAIN_FILE=${PSLOG_ROOT}/cmake/toolchains/host.cmake
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
file(READ "${build_dir}/CMakeCache.txt" configured_cache)
if(CMAKE_HOST_SYSTEM_NAME STREQUAL "Linux" AND NOT configured_cache MATCHES "PSLOG_BOOTLIN_TOOLCHAIN:BOOL=TRUE")
    message(FATAL_ERROR "Linux C-only configure did not select the required Bootlin collection")
endif()

file(GLOB_RECURSE cxx_discovery "${build_dir}/CMakeFiles/*/CMakeCXXCompiler.cmake")
if(cxx_discovery OR configure_stdout MATCHES "CXX compiler identification")
    message(FATAL_ERROR "C-only configure unexpectedly enabled C++ discovery")
endif()

if(CMAKE_HOST_SYSTEM_NAME STREQUAL "Linux")
    execute_process(COMMAND "${CMAKE_COMMAND}" --build "${build_dir}" --target pslog_runtime_probe
        COMMAND_ERROR_IS_FATAL ANY)
    execute_process(COMMAND "${CMAKE_CTEST_COMMAND}" --test-dir "${build_dir}"
        -R runtime_mapping_and_exec_test --output-on-failure COMMAND_ERROR_IS_FATAL ANY)
endif()
