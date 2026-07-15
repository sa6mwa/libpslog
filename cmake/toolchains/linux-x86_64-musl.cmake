set(CMAKE_SYSTEM_NAME Linux)
set(CMAKE_SYSTEM_PROCESSOR x86_64)
set(CMAKE_TRY_COMPILE_TARGET_TYPE STATIC_LIBRARY)
include("${CMAKE_CURRENT_LIST_DIR}/pslog_bootlin.cmake")
pslog_configure_bootlin_toolchain("x86_64-linux-musl")
set(PSLOG_TEST_EXECUTOR
    "${CMAKE_CURRENT_LIST_DIR}/../../scripts/run_sysroot_binary.sh;--sysroot;${CMAKE_SYSROOT}"
    CACHE STRING "Bootlin sysroot launcher for native x86_64 test executables" FORCE)
set(PSLOG_TARGET_ARCH x86_64 CACHE STRING "" FORCE)
set(PSLOG_TARGET_OS linux CACHE STRING "" FORCE)
set(PSLOG_TARGET_LIBC musl CACHE STRING "" FORCE)
set(PSLOG_TARGET_ID x86_64-linux-musl CACHE STRING "" FORCE)
