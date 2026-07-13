if(NOT DEFINED PSLOG_ROOT OR NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_ROOT and PSLOG_BINARY_DIR are required")
endif()

set(test_root "${PSLOG_BINARY_DIR}/discover-target-tools-test")
set(build_dir "${test_root}/build")
set(bin_dir "${test_root}/bin")
set(path_bin "${test_root}/path-bin")
file(REMOVE_RECURSE "${test_root}")
file(MAKE_DIRECTORY "${build_dir}" "${bin_dir}" "${path_bin}")

function(fake_tool path)
    file(WRITE "${path}" "#!/usr/bin/env sh\nexit 0\n")
    file(CHMOD "${path}" PERMISSIONS
        OWNER_READ OWNER_WRITE OWNER_EXECUTE
        GROUP_READ GROUP_EXECUTE
        WORLD_READ WORLD_EXECUTE)
endfunction()

foreach(tool IN ITEMS x86_64-linux-gcc configured-strip configured-readelf)
    fake_tool("${bin_dir}/${tool}")
endforeach()
file(WRITE "${build_dir}/CMakeCache.txt"
"CMAKE_C_COMPILER:FILEPATH=${bin_dir}/x86_64-linux-gcc\n"
"CMAKE_STRIP:FILEPATH=${bin_dir}/configured-strip\n"
"CMAKE_READELF:FILEPATH=${bin_dir}/configured-readelf\n")
execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/discover_target_tools.sh"
        --build-dir "${build_dir}" --target-id x86_64-linux-gnu
    RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error
)
if(NOT result EQUAL 0)
    message(FATAL_ERROR "configured target tool discovery failed:\n${error}")
endif()
foreach(expected IN ITEMS
        "CC=${bin_dir}/x86_64-linux-gcc"
        "STRIP=${bin_dir}/configured-strip"
        "READELF=${bin_dir}/configured-readelf")
    if(NOT output MATCHES "${expected}")
        message(FATAL_ERROR "missing configured tool '${expected}':\n${output}")
    endif()
endforeach()

set(darwin_bin "${test_root}/darwin-bin")
set(darwin_build "${test_root}/darwin-build")
file(MAKE_DIRECTORY "${darwin_bin}" "${darwin_build}")
foreach(tool IN ITEMS clang ld strip install_name_tool otool)
    fake_tool("${darwin_bin}/arm64-apple-darwin25-${tool}")
endforeach()
file(WRITE "${darwin_build}/CMakeCache.txt"
"CMAKE_C_COMPILER:FILEPATH=${darwin_bin}/arm64-apple-darwin25-clang\n"
"PSLOG_OSXCROSS_HOST:STRING=arm64-apple-darwin25\n")
execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/discover_target_tools.sh"
        --build-dir "${darwin_build}" --target-id arm64-apple-darwin
    RESULT_VARIABLE darwin_result OUTPUT_VARIABLE darwin_output ERROR_VARIABLE darwin_error
)
if(NOT darwin_result EQUAL 0 OR
   NOT darwin_output MATCHES "OTOOL=${darwin_bin}/arm64-apple-darwin25-otool" OR
   NOT darwin_output MATCHES "INSTALL_NAME_TOOL=${darwin_bin}/arm64-apple-darwin25-install_name_tool")
    message(FATAL_ERROR "Darwin target-prefixed sibling lookup failed:\n${darwin_output}\n${darwin_error}")
endif()

set(sibling_bin "${test_root}/sibling-bin")
set(sibling_build "${test_root}/sibling-build")
file(MAKE_DIRECTORY "${sibling_bin}" "${sibling_build}")
foreach(tool IN ITEMS cc strip readelf)
    fake_tool("${sibling_bin}/${tool}")
endforeach()
file(WRITE "${sibling_build}/CMakeCache.txt"
"CMAKE_C_COMPILER:FILEPATH=${sibling_bin}/cc\n")
execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/discover_target_tools.sh"
        --build-dir "${sibling_build}" --target-id x86_64-linux-gnu
    RESULT_VARIABLE sibling_result OUTPUT_VARIABLE sibling_output ERROR_VARIABLE sibling_error
)
if(NOT sibling_result EQUAL 0 OR NOT sibling_output MATCHES "STRIP=${sibling_bin}/strip" OR
   NOT sibling_output MATCHES "READELF=${sibling_bin}/readelf")
    message(FATAL_ERROR "unprefixed compiler sibling lookup failed:\n${sibling_output}\n${sibling_error}")
endif()

set(path_build "${test_root}/path-build")
set(path_compiler_bin "${test_root}/path-compiler-bin")
file(MAKE_DIRECTORY "${path_build}" "${path_compiler_bin}")
fake_tool("${path_compiler_bin}/cc")
fake_tool("${path_bin}/strip")
fake_tool("${path_bin}/readelf")
file(WRITE "${path_build}/CMakeCache.txt"
"CMAKE_C_COMPILER:FILEPATH=${path_compiler_bin}/cc\n")
execute_process(
    COMMAND "${CMAKE_COMMAND}" -E env "PATH=${path_bin}:$ENV{PATH}"
        "${PSLOG_ROOT}/scripts/discover_target_tools.sh"
        --build-dir "${path_build}" --target-id fake-linux-gnu
    RESULT_VARIABLE path_result OUTPUT_VARIABLE path_output ERROR_VARIABLE path_error
)
if(NOT path_result EQUAL 0 OR NOT path_output MATCHES "STRIP=${path_bin}/strip" OR
   NOT path_output MATCHES "READELF=${path_bin}/readelf")
    message(FATAL_ERROR "PATH fallback lookup failed:\n${path_output}\n${path_error}")
endif()

file(REMOVE "${darwin_bin}/arm64-apple-darwin25-otool")
file(WRITE "${darwin_build}/CMakeCache.txt"
"CMAKE_C_COMPILER:FILEPATH=${darwin_bin}/arm64-apple-darwin25-clang\n"
"PSLOG_OSXCROSS_HOST:STRING=arm64-apple-darwin25\n"
"CPKT_OTOOL:FILEPATH=/usr/bin/strip\n")
execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/discover_target_tools.sh"
        --build-dir "${darwin_build}" --target-id arm64-apple-darwin
    RESULT_VARIABLE host_tool_result OUTPUT_VARIABLE host_tool_output ERROR_VARIABLE host_tool_error
)
if(host_tool_result EQUAL 0 OR NOT host_tool_error MATCHES "refusing known host OTOOL")
    message(FATAL_ERROR "Darwin discovery accepted a known host tool:\n${host_tool_output}\n${host_tool_error}")
endif()

file(REMOVE_RECURSE "${test_root}")
