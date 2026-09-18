foreach(required IN ITEMS PSLOG_BINARY_DIR PSLOG_ROOT)
    if(NOT DEFINED ${required} OR "${${required}}" STREQUAL "")
        message(FATAL_ERROR "${required} is required")
    endif()
endforeach()

set(test_root "${PSLOG_BINARY_DIR}/linux-bootlin-required-test")
file(REMOVE_RECURSE "${test_root}")
if(PSLOG_ALLOW_HOST_COMPILER OR NOT PSLOG_HOST_SYSTEM_PROCESSOR MATCHES "^(x86_64|amd64)$")
    message(STATUS "SKIP: Bootlin default-policy probes require a supported Linux host without the developer compiler override")
else()
    execute_process(
        COMMAND "${CMAKE_COMMAND}" -S "${PSLOG_ROOT}" -B "${test_root}"
            -DPSLOG_ALLOW_HOST_COMPILER=OFF
            -DPSLOG_BUILD_TESTS=OFF
            -DPSLOG_BUILD_BENCHMARKS=OFF
        RESULT_VARIABLE configure_result
        OUTPUT_VARIABLE configure_output
        ERROR_VARIABLE configure_error
    )
    if(NOT configure_result EQUAL 0)
        message(FATAL_ERROR "Default Linux configure failed: ${configure_output}${configure_error}")
    endif()
    file(READ "${test_root}/CMakeCache.txt" cache)
    if(NOT cache MATCHES "PSLOG_BOOTLIN_TOOLCHAIN:BOOL=TRUE" OR cache MATCHES "CMAKE_C_COMPILER:[^=]*=/usr/bin")
        message(FATAL_ERROR "Default Linux configure did not bootstrap Bootlin")
    endif()
    file(STRINGS "${test_root}/CMakeCache.txt" compiler_line REGEX "^CMAKE_C_COMPILER:[^=]*=")
    string(REGEX REPLACE "^[^=]+=" "" selected_compiler "${compiler_line}")
    file(WRITE "${test_root}/spoof.cmake" "set(PSLOG_BOOTLIN_TOOLCHAIN TRUE CACHE BOOL \"\" FORCE)\nset(CMAKE_C_COMPILER \"${selected_compiler}\")\n")
    execute_process(COMMAND "${CMAKE_COMMAND}" -S "${PSLOG_ROOT}" -B "${test_root}/spoof"
        "-DCMAKE_TOOLCHAIN_FILE=${test_root}/spoof.cmake"
        RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error)
    if(result EQUAL 0 OR NOT "${output}${error}" MATCHES "require a pinned Bootlin")
        message(FATAL_ERROR "Spoofed Bootlin marker was not rejected: ${output}${error}")
    endif()

    file(WRITE "${test_root}/fake-host-compiler" "#!/bin/sh\nprintf '/outside-bootlin/ld\\n'\n")
    file(CHMOD "${test_root}/fake-host-compiler" PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE)
    execute_process(COMMAND "${CMAKE_COMMAND}" -S "${PSLOG_ROOT}" -B "${test_root}/host-override"
        -DPSLOG_BOOTLIN_C_COMPILER_OVERRIDE=${test_root}/fake-host-compiler
        RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error)
    if(result EQUAL 0 OR NOT "${output}${error}" MATCHES "Bootlin compiler.*linker")
        message(FATAL_ERROR "Host compiler override was not rejected: ${output}${error}")
    endif()
endif()

execute_process(COMMAND "${CMAKE_COMMAND}" -S "${PSLOG_ROOT}" -B "${test_root}/developer-host-override"
    -DPSLOG_ALLOW_HOST_COMPILER=ON
    -DPSLOG_BUILD_TESTS=OFF
    -DPSLOG_BUILD_BENCHMARKS=OFF
    RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error)
if(NOT result EQUAL 0)
    message(FATAL_ERROR "Developer host compiler override failed: ${output}${error}")
endif()
file(READ "${test_root}/developer-host-override/CMakeCache.txt" override_cache)
if(NOT override_cache MATCHES "PSLOG_ALLOW_HOST_COMPILER:BOOL=ON" OR
   override_cache MATCHES "PSLOG_BOOTLIN_TOOLCHAIN:BOOL=TRUE")
    message(FATAL_ERROR "Developer host compiler override did not bypass Bootlin")
endif()

if(CMAKE_HOST_SYSTEM_NAME STREQUAL "Linux")
    find_program(PSLOG_HOST_CC NAMES cc REQUIRED)
    execute_process(COMMAND "${CMAKE_COMMAND}" -S "${PSLOG_ROOT}" -B "${test_root}/developer-host-coverage"
        -DPSLOG_ALLOW_HOST_COMPILER=ON
        -DPSLOG_ENABLE_COVERAGE=ON
        -DPSLOG_BUILD_TESTS=OFF
        -DPSLOG_BUILD_BENCHMARKS=OFF
        -DCMAKE_C_COMPILER=${PSLOG_HOST_CC}
        RESULT_VARIABLE coverage_result OUTPUT_VARIABLE coverage_output ERROR_VARIABLE coverage_error)
    if(NOT coverage_result EQUAL 0)
        message(FATAL_ERROR
            "Developer host compiler coverage configure failed: ${coverage_output}${coverage_error}")
    endif()
endif()
