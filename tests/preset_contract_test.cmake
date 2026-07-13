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
foreach(target_id IN ITEMS
        x86_64-linux-gnu x86_64-linux-musl aarch64-linux-gnu aarch64-linux-musl
        armhf-linux-gnu armhf-linux-musl arm64-apple-darwin)
    if(NOT presets MATCHES "\"PSLOG_TARGET_ID\": \"${target_id}\"")
        message(FATAL_ERROR "release preset does not declare target ID: ${target_id}")
    endif()
endforeach()
