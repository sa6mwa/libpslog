set(runtime_archive "${PSLOG_ROOT}/dist/libpslog-${PSLOG_VERSION}-${PSLOG_TARGET_ID}.tar.gz")
set(dev_archive "${PSLOG_ROOT}/dist/libpslog-${PSLOG_VERSION}-${PSLOG_TARGET_ID}-dev.tar.gz")
set(single_header "${PSLOG_ROOT}/dist/pslog-${PSLOG_VERSION}.h")
set(single_header_gz "${single_header}.gz")
set(checksums_file "${PSLOG_ROOT}/dist/libpslog-${PSLOG_VERSION}-CHECKSUMS")

file(REMOVE "${runtime_archive}" "${dev_archive}" "${single_header}" "${single_header_gz}" "${checksums_file}")

execute_process(
    COMMAND "${CMAKE_COMMAND}" --build "${PSLOG_BINARY_DIR}" --target package-runtime package-dev package-single-header
    RESULT_VARIABLE build_result
)
if(NOT build_result EQUAL 0)
    message(FATAL_ERROR "failed to build package archives and single-header artifact")
endif()

execute_process(
    COMMAND "${CMAKE_COMMAND}" --build "${PSLOG_BINARY_DIR}" --target package-checksums
    RESULT_VARIABLE checksum_build_result
)
if(NOT checksum_build_result EQUAL 0)
    message(FATAL_ERROR "failed to build package checksums")
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

if(NOT EXISTS "${single_header}")
    message(FATAL_ERROR "missing single-header artifact: ${single_header}")
endif()
if(NOT EXISTS "${single_header_gz}")
    message(FATAL_ERROR "missing gzipped single-header artifact: ${single_header_gz}")
endif()
if(NOT EXISTS "${checksums_file}")
    message(FATAL_ERROR "missing checksums file: ${checksums_file}")
endif()

file(READ "${single_header}" single_header_text)
if(NOT single_header_text MATCHES "PSLOG_IMPLEMENTATION")
    message(FATAL_ERROR "single-header artifact is missing PSLOG_IMPLEMENTATION section")
endif()
if(NOT single_header_text MATCHES "Artifact: pslog-${PSLOG_VERSION}\\.h")
    message(FATAL_ERROR "single-header artifact is missing versioned artifact metadata")
endif()
if(NOT single_header_text MATCHES "MIT License")
    message(FATAL_ERROR "single-header artifact is missing embedded license text")
endif()

file(READ "${checksums_file}" checksums_text)
if(checksums_text MATCHES "(^|\n)[0-9a-f]+  pslog-${PSLOG_VERSION}\\.h(\n|$)")
    message(FATAL_ERROR "checksums file unexpectedly includes pslog-${PSLOG_VERSION}.h")
endif()
if(NOT checksums_text MATCHES "(^|\n)[0-9a-f]+  pslog-${PSLOG_VERSION}\\.h\\.gz(\n|$)")
    message(FATAL_ERROR "checksums file is missing pslog-${PSLOG_VERSION}.h.gz")
endif()
