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
