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

set(CMAKE_SYSTEM_NAME Darwin)
set(CMAKE_HOST_SYSTEM_PROCESSOR arm64)
set(CMAKE_SYSTEM_PROCESSOR arm64)
set(CMAKE_C_COMPILER /usr/bin/clang)
set(PSLOG_NONSHIPPED_ELF_LINKER_FLAGS "")
include("${PSLOG_ROOT}/cmake/pslog_local_runtime.cmake")
if(PSLOG_LOCAL_RUNTIME_FLAGS OR PSLOG_NONSHIPPED_ELF_LINKER_FLAGS OR NOT CMAKE_C_COMPILER STREQUAL "/usr/bin/clang")
    message(FATAL_ERROR "Native macOS runtime policy changed the compiler or injected ELF flags")
endif()
