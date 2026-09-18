if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

# The host preset is shared by Linux and macOS. Simulate macOS script-mode
# evaluation to prove it leaves compiler selection to CMake rather than trying
# to provision a Linux Bootlin collection.
set(CMAKE_HOST_SYSTEM_NAME "Darwin")
include("${PSLOG_ROOT}/cmake/toolchains/host.cmake")
if(DEFINED PSLOG_BOOTLIN_TOOLCHAIN)
    message(FATAL_ERROR "macOS host toolchain must not configure Bootlin")
endif()

set(arm_policy "${PSLOG_ROOT}/build/host-toolchain-arm-policy.cmake")
file(WRITE "${arm_policy}"
    "set(CMAKE_HOST_SYSTEM_NAME Linux)\n"
    "set(CMAKE_HOST_SYSTEM_PROCESSOR aarch64)\n"
    "set(PSLOG_ALLOW_HOST_COMPILER OFF)\n"
    "include(\"${PSLOG_ROOT}/cmake/toolchains/host.cmake\")\n")
execute_process(COMMAND "${CMAKE_COMMAND}" -P "${arm_policy}"
    RESULT_VARIABLE arm_result OUTPUT_VARIABLE arm_output ERROR_VARIABLE arm_error)
if(arm_result EQUAL 0 OR NOT "${arm_output}${arm_error}" MATCHES "PSLOG_ALLOW_HOST_COMPILER")
    message(FATAL_ERROR "Native ARM Linux host policy was not rejected clearly: ${arm_output}${arm_error}")
endif()

file(WRITE "${arm_policy}"
    "set(CMAKE_HOST_SYSTEM_NAME Linux)\n"
    "set(CMAKE_HOST_SYSTEM_PROCESSOR aarch64)\n"
    "set(PSLOG_ALLOW_HOST_COMPILER ON)\n"
    "include(\"${PSLOG_ROOT}/cmake/toolchains/host.cmake\")\n"
    "if(DEFINED PSLOG_BOOTLIN_TOOLCHAIN)\n"
    "  message(FATAL_ERROR \"host compiler override configured Bootlin\")\n"
    "endif()\n")
execute_process(COMMAND "${CMAKE_COMMAND}" -P "${arm_policy}"
    RESULT_VARIABLE override_result OUTPUT_VARIABLE override_output ERROR_VARIABLE override_error)
if(NOT override_result EQUAL 0)
    message(FATAL_ERROR "Native ARM Linux host override failed: ${override_output}${override_error}")
endif()
file(REMOVE "${arm_policy}")
unset(PSLOG_ALLOW_HOST_COMPILER CACHE)
unset(PSLOG_ALLOW_HOST_COMPILER)
unset(PSLOG_HOST_COMPILER_OVERRIDE CACHE)
unset(PSLOG_HOST_COMPILER_OVERRIDE)

set(CMAKE_SYSTEM_NAME Darwin)
set(CMAKE_HOST_SYSTEM_PROCESSOR arm64)
set(CMAKE_SYSTEM_PROCESSOR arm64)
set(CMAKE_C_COMPILER /usr/bin/clang)
set(PSLOG_NONSHIPPED_ELF_LINKER_FLAGS "")
include("${PSLOG_ROOT}/cmake/pslog_local_runtime.cmake")
if(PSLOG_LOCAL_RUNTIME_FLAGS OR PSLOG_NONSHIPPED_ELF_LINKER_FLAGS OR NOT CMAKE_C_COMPILER STREQUAL "/usr/bin/clang")
    message(FATAL_ERROR "Native macOS runtime policy changed the compiler or injected ELF flags")
endif()
