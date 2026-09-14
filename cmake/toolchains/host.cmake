# Use the pinned Bootlin collection for every local Linux configuration.  On
# macOS, leave compiler selection to the host exactly as CMake normally does.
if(CMAKE_HOST_SYSTEM_NAME STREQUAL "Linux")
    execute_process(
        COMMAND uname -m
        RESULT_VARIABLE host_arch_result
        OUTPUT_VARIABLE host_arch
        OUTPUT_STRIP_TRAILING_WHITESPACE
    )
    if(NOT host_arch_result EQUAL 0)
        message(FATAL_ERROR "Unable to determine the Linux host architecture for the Bootlin toolchain")
    endif()
    if(host_arch MATCHES "^(x86_64|amd64)$")
        set(bootlin_host_target x86_64-linux-gnu)
    elseif(host_arch MATCHES "^(aarch64|arm64)$")
        set(bootlin_host_target aarch64-linux-gnu)
    elseif(host_arch MATCHES "^armv7l$")
        set(bootlin_host_target armhf-linux-gnu)
    else()
        message(FATAL_ERROR "No pinned Bootlin collection supports Linux host architecture: ${host_arch}")
    endif()
    include("${CMAKE_CURRENT_LIST_DIR}/pslog_bootlin.cmake")
    pslog_configure_bootlin_toolchain("${bootlin_host_target}")
elseif(NOT CMAKE_HOST_SYSTEM_NAME STREQUAL "Darwin")
    message(FATAL_ERROR "libpslog supports Linux with Bootlin or native macOS builds; host is ${CMAKE_HOST_SYSTEM_NAME}")
endif()
