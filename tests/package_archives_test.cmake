set(runtime_archive "${PSLOG_ROOT}/dist/libpslog-${PSLOG_VERSION}-${PSLOG_TARGET_ID}.tar.gz")
set(dev_archive "${PSLOG_ROOT}/dist/libpslog-${PSLOG_VERSION}-${PSLOG_TARGET_ID}-dev.tar.gz")

file(REMOVE "${runtime_archive}" "${dev_archive}")

execute_process(
    COMMAND "${CMAKE_COMMAND}" --build "${PSLOG_BINARY_DIR}" --target package-runtime package-dev
    RESULT_VARIABLE build_result
)
if(NOT build_result EQUAL 0)
    message(FATAL_ERROR "failed to build package archives")
endif()

function(assert_archive_layout archive_path)
    if(NOT EXISTS "${archive_path}")
        message(FATAL_ERROR "missing archive: ${archive_path}")
    endif()

    execute_process(
        COMMAND "${CMAKE_COMMAND}" -E tar tf "${archive_path}"
        RESULT_VARIABLE tar_result
        OUTPUT_VARIABLE archive_listing
    )
    if(NOT tar_result EQUAL 0)
        message(FATAL_ERROR "failed to list archive contents: ${archive_path}")
    endif()

    if(NOT archive_listing MATCHES "(^|\n)(\\./)?share/doc/libpslog/LICENSE(\n|$)")
        message(FATAL_ERROR "archive missing share/doc/libpslog/LICENSE: ${archive_path}")
    endif()
    if(NOT archive_listing MATCHES "(^|\n)(\\./)?share/doc/libpslog/README.md(\n|$)")
        message(FATAL_ERROR "archive missing share/doc/libpslog/README.md: ${archive_path}")
    endif()
    if(NOT archive_listing MATCHES "(^|\n)(\\./)?include/pslog.h(\n|$)")
        message(FATAL_ERROR "archive missing include/pslog.h: ${archive_path}")
    endif()
    if(NOT archive_listing MATCHES "(^|\n)(\\./)?include/pslog_version.h(\n|$)")
        message(FATAL_ERROR "archive missing include/pslog_version.h: ${archive_path}")
    endif()
    if(archive_listing MATCHES "(^|\n)(\\./)?share/libpslog(/|\n|$)")
        message(FATAL_ERROR "archive still contains legacy share/libpslog path: ${archive_path}")
    endif()
    if(archive_listing MATCHES "(^|\n)(\\./)?share/doc/libpslog/demo.gif(\n|$)")
        message(FATAL_ERROR "archive unexpectedly contains share/doc/libpslog/demo.gif: ${archive_path}")
    endif()
    if(archive_listing MATCHES "(^|\n)(\\./)?share/doc/libpslog/elevatorpitch.gif(\n|$)")
        message(FATAL_ERROR "archive unexpectedly contains share/doc/libpslog/elevatorpitch.gif: ${archive_path}")
    endif()
    if(DEFINED PSLOG_SHARED_SONAME
       AND NOT PSLOG_SHARED_SONAME STREQUAL ""
       AND NOT PSLOG_SHARED_SONAME STREQUAL PSLOG_SHARED_LIB_NAME
       AND archive_path STREQUAL runtime_archive)
        string(REPLACE "." "\\." shared_soname_regex "${PSLOG_SHARED_SONAME}")
        if(NOT archive_listing MATCHES "(^|\n)(\\./)?lib/${shared_soname_regex}(\n|$)")
            message(FATAL_ERROR "runtime archive missing shared-library SONAME entry lib/${PSLOG_SHARED_SONAME}: ${archive_path}")
        endif()
    endif()

    file(READ "${archive_path}" archive_xfl HEX OFFSET 8 LIMIT 1)
    string(TOLOWER "${archive_xfl}" archive_xfl)
    if(NOT archive_xfl STREQUAL "02")
        message(FATAL_ERROR "archive is not using gzip maximum compression header: ${archive_path}")
    endif()
endfunction()

assert_archive_layout("${runtime_archive}")
assert_archive_layout("${dev_archive}")
