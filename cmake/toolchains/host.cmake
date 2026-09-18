# Bootlin publishes x86-64-hosted compilers. A developer may explicitly opt
# out for a host-target build; release and cross-target toolchains never use
# this escape hatch.
if(NOT DEFINED PSLOG_ALLOW_HOST_COMPILER AND DEFINED ENV{PSLOG_ALLOW_HOST_COMPILER})
    set(PSLOG_ALLOW_HOST_COMPILER "$ENV{PSLOG_ALLOW_HOST_COMPILER}" CACHE BOOL
        "Allow the host compiler for developer host-target builds")
endif()
if(PSLOG_ALLOW_HOST_COMPILER)
    set(PSLOG_ALLOW_HOST_COMPILER ON CACHE BOOL
        "Allow the host compiler for developer host-target builds" FORCE)
    list(APPEND CMAKE_TRY_COMPILE_PLATFORM_VARIABLES PSLOG_ALLOW_HOST_COMPILER)
    list(REMOVE_DUPLICATES CMAKE_TRY_COMPILE_PLATFORM_VARIABLES)
    set(PSLOG_HOST_COMPILER_OVERRIDE TRUE CACHE BOOL "Developer host compiler override" FORCE)
    message(STATUS "libpslog: developer host compiler override is enabled")
elseif(CMAKE_HOST_SYSTEM_NAME STREQUAL "Linux")
    if(CMAKE_HOST_SYSTEM_PROCESSOR MATCHES "^(x86_64|amd64)$")
        set(bootlin_host_target x86_64-linux-gnu)
    else()
        message(FATAL_ERROR
            "Native Linux host architecture ${CMAKE_HOST_SYSTEM_PROCESSOR} is unsupported by the "
            "published Bootlin collections, whose compiler executables require x86-64. "
            "For a developer-only host build, set PSLOG_ALLOW_HOST_COMPILER=1.")
    endif()
    include("${CMAKE_CURRENT_LIST_DIR}/pslog_bootlin.cmake")
    pslog_configure_bootlin_toolchain("${bootlin_host_target}")
elseif(NOT CMAKE_HOST_SYSTEM_NAME STREQUAL "Darwin")
    message(FATAL_ERROR "libpslog supports Linux with Bootlin or native macOS builds; host is ${CMAKE_HOST_SYSTEM_NAME}")
endif()
