# Native x86_64 fuzzing keeps the complete Bootlin collection for the linker,
# binutils, libc, and headers, then replaces only the compiler drivers with
# pinned AFL++ wrappers that delegate back to that collection.
set(CMAKE_SYSTEM_NAME Linux)
set(CMAKE_SYSTEM_PROCESSOR x86_64)
set(CMAKE_TRY_COMPILE_TARGET_TYPE STATIC_LIBRARY)

set(afl_resolver "${CMAKE_CURRENT_LIST_DIR}/../../scripts/cpkt-aflpp.sh")
if(NOT EXISTS "${afl_resolver}")
    message(FATAL_ERROR "missing lifecycle AFL++ resolver: ${afl_resolver}")
endif()
execute_process(
    COMMAND "${afl_resolver}" discover
    RESULT_VARIABLE afl_result
    OUTPUT_VARIABLE afl_description
    ERROR_VARIABLE afl_error
)
if(NOT afl_result EQUAL 0)
    message(FATAL_ERROR "Unable to provision pinned AFL++ GCC-plugin tooling.\n${afl_error}")
endif()
foreach(key cc cxx helper root)
    string(REGEX MATCH "${key}=([^\r\n]+)" afl_match "${afl_description}")
    if(NOT afl_match)
        message(FATAL_ERROR "AFL++ resolver did not report ${key}")
    endif()
    set(afl_${key} "${CMAKE_MATCH_1}")
endforeach()

set(ENV{AFL_PATH} "${afl_helper}")
# Select only the AFL++ compiler drivers before CMake's first compiler
# detection. pslog_bootlin then supplies every linker, binutil, and sysroot
# value from the matching pinned Bootlin collection.
set(PSLOG_BOOTLIN_C_COMPILER_OVERRIDE "${afl_cc}")
set(PSLOG_BOOTLIN_CXX_COMPILER_OVERRIDE "${afl_cxx}")
include("${CMAKE_CURRENT_LIST_DIR}/pslog_bootlin.cmake")
pslog_configure_bootlin_toolchain("x86_64-linux-gnu")
set(PSLOG_TEST_EXECUTOR
    "${CMAKE_CURRENT_LIST_DIR}/../../scripts/run_sysroot_binary.sh;--sysroot;${CMAKE_SYSROOT}"
    CACHE STRING "Bootlin sysroot launcher for native AFL++ test executables" FORCE)
set(PSLOG_TARGET_ARCH x86_64 CACHE STRING "" FORCE)
set(PSLOG_TARGET_OS linux CACHE STRING "" FORCE)
set(PSLOG_TARGET_LIBC gnu CACHE STRING "" FORCE)
set(PSLOG_TARGET_ID x86_64-linux-gnu CACHE STRING "" FORCE)
