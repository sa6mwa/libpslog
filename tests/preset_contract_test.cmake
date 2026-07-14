if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

file(READ "${PSLOG_ROOT}/CMakePresets.json" presets)
foreach(required_preset IN ITEMS
        base debug debug-lua valgrind fuzz
        x86_64-linux-gnu-release x86_64-linux-musl-release
        aarch64-linux-gnu-release aarch64-linux-musl-release
        armhf-linux-gnu-release armhf-linux-musl-release
        arm64-apple-darwin-release)
    if(NOT presets MATCHES "\"name\": \"${required_preset}\"")
        message(FATAL_ERROR "lifecycle preset is missing: ${required_preset}")
    endif()
endforeach()
if(NOT presets MATCHES "\"CMAKE_EXPORT_COMPILE_COMMANDS\": \"ON\"")
    message(FATAL_ERROR "base preset does not enable compile commands")
endif()
if(NOT presets MATCHES "\"PSLOG_DEPENDENCY_MODE\": \"bundled-sdk\"")
    message(FATAL_ERROR "base preset does not pin bundled-sdk dependency mode")
endif()
if(NOT presets MATCHES "\"PSLOG_BUILD_LUA\": \"OFF\"")
    message(FATAL_ERROR "base preset must keep staged Lua CTest opt-in")
endif()
file(READ "${PSLOG_ROOT}/CMakeLists.txt" cmake_lists)
if(cmake_lists MATCHES "pkg_check_modules\\(PSLOG_LUA" OR
   cmake_lists MATCHES "PSLOG_LUA_INCLUDE_DIRS" OR
   cmake_lists MATCHES "PSLOG_LUA_LIBRARY_DIRS")
    message(FATAL_ERROR "Lua CTest must use the repository-local staged Lua inputs, not pkg-config host paths")
endif()
foreach(required_lua_stage IN ITEMS PSLOG_LUA_STAGE_ROOT PSLOG_LUA_SDK_PREFIX)
    if(NOT cmake_lists MATCHES "${required_lua_stage}")
        message(FATAL_ERROR "Lua CTest staging contract is missing: ${required_lua_stage}")
    endif()
endforeach()
foreach(target_id IN ITEMS
        x86_64-linux-gnu x86_64-linux-musl aarch64-linux-gnu aarch64-linux-musl
        armhf-linux-gnu armhf-linux-musl arm64-apple-darwin)
    if(NOT presets MATCHES "\"PSLOG_TARGET_ID\": \"${target_id}\"")
        message(FATAL_ERROR "release preset does not declare target ID: ${target_id}")
    endif()
endforeach()
