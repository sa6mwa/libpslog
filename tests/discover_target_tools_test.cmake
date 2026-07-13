if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()
if(NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_BINARY_DIR is required")
endif()

set(test_root "${PSLOG_BINARY_DIR}/discover-target-tools-test")
set(build_dir "${test_root}/build")
set(bin_dir "${test_root}/bin")
file(REMOVE_RECURSE "${test_root}")
file(MAKE_DIRECTORY "${build_dir}" "${bin_dir}")
foreach(tool IN ITEMS x86_64-linux-gcc x86_64-linux-strip x86_64-linux-readelf)
    file(WRITE "${bin_dir}/${tool}" "#!/usr/bin/env sh\nexit 0\n")
    file(CHMOD "${bin_dir}/${tool}" PERMISSIONS
        OWNER_READ OWNER_WRITE OWNER_EXECUTE
        GROUP_READ GROUP_EXECUTE
        WORLD_READ WORLD_EXECUTE)
endforeach()
file(WRITE "${build_dir}/CMakeCache.txt"
"CMAKE_C_COMPILER:FILEPATH=${bin_dir}/x86_64-linux-gcc\n"
"CMAKE_STRIP:FILEPATH=${bin_dir}/x86_64-linux-strip\n"
"CMAKE_READELF:FILEPATH=${bin_dir}/x86_64-linux-readelf\n")

execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/discover_target_tools.sh"
        --build-dir "${build_dir}"
        --target-id x86_64-linux-gnu
    RESULT_VARIABLE result
    OUTPUT_VARIABLE output
    ERROR_VARIABLE error
)
if(NOT result EQUAL 0)
    message(FATAL_ERROR "target tool discovery failed:\n${error}")
endif()
foreach(expected IN ITEMS
        "CC=${bin_dir}/x86_64-linux-gcc"
        "STRIP=${bin_dir}/x86_64-linux-strip"
        "READELF=${bin_dir}/x86_64-linux-readelf")
    if(NOT output MATCHES "${expected}")
        message(FATAL_ERROR "missing discovered tool '${expected}':\n${output}")
    endif()
endforeach()
file(REMOVE "${build_dir}/CMakeCache.txt")
execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/discover_target_tools.sh"
        --build-dir "${build_dir}"
        --target-id aarch64-linux-gnu
    RESULT_VARIABLE missing_result
)
if(missing_result EQUAL 0)
    message(FATAL_ERROR "target tool discovery accepted a missing configured cache")
endif()
file(REMOVE_RECURSE "${test_root}")
