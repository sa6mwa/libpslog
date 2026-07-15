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
if(presets MATCHES "\"name\": \"debug\",[^\n]*\n[ \t]*\"inherits\": \"base\",[^\n]*\n[ \t]*\"toolchainFile\"" OR
   presets MATCHES "\"name\": \"host\",[^\n]*\n[ \t]*\"inherits\": \"base\",[^\n]*\n[ \t]*\"toolchainFile\"")
    message(FATAL_ERROR "standard debug and host presets must remain host-native and must not force a cross or Bootlin toolchain")
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
if(NOT cmake_lists MATCHES "CMAKE_SHARED_LIBRARY_PREFIX}pslog\\$\\{CMAKE_SHARED_LIBRARY_SUFFIX")
    message(FATAL_ERROR "Lua SDK shared library path must use the platform shared-library suffix")
endif()
file(READ "${PSLOG_ROOT}/Makefile" makefile_text)
if(makefile_text MATCHES "test -e \"\\$\\(LUA_SDK_LIB_DIR\\)/libpslog\\.so\"[^\n]*$")
    message(FATAL_ERROR "Make Lua SDK check must not require only the Linux .so suffix")
endif()
file(READ "${PSLOG_ROOT}/lua/scripts/run_interop_embedder_test.sh" lua_interop_runner)
if(NOT lua_interop_runner MATCHES "libpslog\\.dylib")
    message(FATAL_ERROR "Lua interop runner must accept the Darwin libpslog.dylib SDK library")
endif()
if(NOT lua_interop_runner MATCHES "uname -s" OR
   NOT lua_interop_runner MATCHES "Darwin\\)" OR
   lua_interop_runner MATCHES "lua_libs=\"[^\n\"]*-ldl")
    message(FATAL_ERROR "Lua interop runner must select the dynamic-loader library per host platform")
endif()
foreach(target_id IN ITEMS
        x86_64-linux-gnu x86_64-linux-musl aarch64-linux-gnu aarch64-linux-musl
        armhf-linux-gnu armhf-linux-musl arm64-apple-darwin)
    if(NOT presets MATCHES "\"PSLOG_TARGET_ID\": \"${target_id}\"")
        message(FATAL_ERROR "release preset does not declare target ID: ${target_id}")
    endif()
endforeach()
foreach(native_toolchain IN ITEMS linux-x86_64-gnu.cmake linux-x86_64-musl.cmake)
    file(READ "${PSLOG_ROOT}/cmake/toolchains/${native_toolchain}" native_toolchain_text)
    if(NOT native_toolchain_text MATCHES "PSLOG_TEST_EXECUTOR" OR
       NOT native_toolchain_text MATCHES "run_sysroot_binary\\.sh")
        message(FATAL_ERROR "native Bootlin toolchain is missing its sysroot runtime executor: ${native_toolchain}")
    endif()
endforeach()

file(READ "${PSLOG_ROOT}/scripts/clangd_check.sh" clangd_gate)
if(NOT clangd_gate MATCHES "database_dir=.*build/debug")
    message(FATAL_ERROR "clangd gate must use only the native debug compile database")
endif()
if(NOT clangd_gate MATCHES "source_file=.*examples/example\\.c")
    message(FATAL_ERROR "clangd gate must validate the native public C consumer")
endif()
if(NOT clangd_gate MATCHES "native public-consumer command is missing")
    message(FATAL_ERROR "clangd gate must reject a missing public-consumer compilation command")
endif()
if(clangd_gate MATCHES "build/(host|x86_64-linux|aarch64|armhf|arm64-apple-darwin)")
    message(FATAL_ERROR "clangd gate must not select a host-release or cross-target compile database")
endif()
if(NOT cmake_lists MATCHES "add_executable\\(pslog_clangd_public_consumer EXCLUDE_FROM_ALL examples/example\\.c\\)" OR
   NOT cmake_lists MATCHES "target_link_libraries\\(pslog_clangd_public_consumer PRIVATE pslog_static Threads::Threads\\)")
    message(FATAL_ERROR "clangd public consumer is missing its exact public CMake compilation target")
endif()
file(READ "${PSLOG_ROOT}/cmake/toolchains/linux-aflpp.cmake" afl_toolchain)
if(NOT afl_toolchain MATCHES "set\\(PSLOG_BOOTLIN_C_COMPILER_OVERRIDE \"\\$\\{afl_cc\\}\"\\)" OR
   NOT afl_toolchain MATCHES "pslog_configure_bootlin_toolchain\\(\"x86_64-linux-gnu\"\\)" OR
   NOT afl_toolchain MATCHES "PSLOG_TEST_EXECUTOR" OR
   NOT afl_toolchain MATCHES "run_sysroot_binary\\.sh")
    message(FATAL_ERROR "AFL++ wrappers must select before Bootlin toolchain setup and retain its native sysroot executor")
endif()
file(READ "${PSLOG_ROOT}/cmake/toolchains/pslog_bootlin.cmake" bootlin_toolchain)
if(bootlin_toolchain MATCHES "IS_EXECUTABLE")
    message(FATAL_ERROR "Bootlin toolchain validation must remain compatible with the CMake 3.21 minimum")
endif()
if(NOT lua_interop_runner MATCHES "CMAKE_SYSROOT" OR
   NOT lua_interop_runner MATCHES "run_sysroot_binary\\.sh" OR
   NOT lua_interop_runner MATCHES "--library-path" OR
   NOT lua_interop_runner MATCHES "LD_LIBRARY_PATH=.*sdk_lib_dir")
    message(FATAL_ERROR "Lua interop consumers must use the selected sysroot loader when configured and direct SDK library paths for native host builds")
endif()
