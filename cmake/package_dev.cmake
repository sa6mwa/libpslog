set(package_root "${PSLOG_BINARY_DIR}/package/dev")
file(REMOVE_RECURSE "${package_root}")
file(MAKE_DIRECTORY "${package_root}/include")
file(MAKE_DIRECTORY "${package_root}/lib")
file(MAKE_DIRECTORY "${package_root}/share/doc/libpslog")
file(COPY "${PSLOG_PUBLIC_HEADER}" DESTINATION "${package_root}/include")
file(COPY "${PSLOG_PUBLIC_VERSION_HEADER}" DESTINATION "${package_root}/include")
file(COPY "${PSLOG_STATIC_LIB}" DESTINATION "${package_root}/lib")
file(COPY "${PSLOG_ROOT}/LICENSE" DESTINATION "${package_root}/share/doc/libpslog")
file(COPY "${PSLOG_ROOT}/README.md" DESTINATION "${package_root}/share/doc/libpslog")
file(MAKE_DIRECTORY "${PSLOG_ROOT}/dist")
set(archive_base "${PSLOG_ROOT}/dist/libpslog-${PSLOG_VERSION}-${PSLOG_TARGET_ID}-dev.tar")
set(archive "${archive_base}.gz")
find_program(PSLOG_TAR_BIN NAMES tar)
find_program(PSLOG_GZIP_BIN NAMES gzip)
if(NOT PSLOG_TAR_BIN)
    message(FATAL_ERROR "failed to find tar for dev archive creation")
endif()
if(NOT PSLOG_GZIP_BIN)
    message(FATAL_ERROR "failed to find gzip for dev archive creation")
endif()
file(REMOVE "${archive_base}" "${archive}")
execute_process(
    COMMAND "${PSLOG_TAR_BIN}" -cf "${archive_base}" --format=gnu .
    WORKING_DIRECTORY "${package_root}"
    RESULT_VARIABLE tar_result
)
if(NOT tar_result EQUAL 0)
    message(FATAL_ERROR "failed to create dev tar archive")
endif()
file(REMOVE "${archive}")
execute_process(
    COMMAND "${PSLOG_GZIP_BIN}" -9 -f "${archive_base}"
    RESULT_VARIABLE gzip_result
)
if(NOT gzip_result EQUAL 0)
    message(FATAL_ERROR "failed to gzip dev archive")
endif()
