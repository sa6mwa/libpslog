set(archive_name "libpslog-${PSLOG_VERSION}-${PSLOG_TARGET_ID}")
set(package_stage_root "${PSLOG_BINARY_DIR}/package/archive")
set(package_root "${package_stage_root}/${archive_name}")
file(REMOVE_RECURSE "${package_stage_root}")
file(MAKE_DIRECTORY "${package_root}/include")
file(MAKE_DIRECTORY "${package_root}/lib")
file(MAKE_DIRECTORY "${package_root}/share/doc/libpslog")

file(COPY "${PSLOG_PUBLIC_HEADER}" DESTINATION "${package_root}/include")
file(COPY "${PSLOG_PUBLIC_VERSION_HEADER}" DESTINATION "${package_root}/include")
file(COPY "${PSLOG_SHARED_LIB}" DESTINATION "${package_root}/lib")
file(COPY "${PSLOG_STATIC_LIB}" DESTINATION "${package_root}/lib")

set(packaged_shared_lib "${package_root}/lib/${PSLOG_SHARED_LIB_NAME}")
set(packaged_static_lib "${package_root}/lib/${PSLOG_STATIC_LIB_NAME}")
if(DEFINED PSLOG_STRIP AND NOT PSLOG_STRIP STREQUAL "")
    if(PSLOG_TARGET_ID MATCHES "darwin")
        set(strip_shared_flag "-S")
        set(strip_static_flag "-S")
    else()
        set(strip_shared_flag "--strip-unneeded")
        set(strip_static_flag "--strip-debug")
    endif()
    execute_process(
        COMMAND "${PSLOG_STRIP}" "${strip_shared_flag}" "${packaged_shared_lib}"
        RESULT_VARIABLE strip_shared_result
    )
    if(NOT strip_shared_result EQUAL 0)
        message(FATAL_ERROR "failed to strip ${packaged_shared_lib}")
    endif()
    if(NOT PSLOG_TARGET_ID MATCHES "darwin")
        execute_process(
            COMMAND "${PSLOG_STRIP}" "${strip_static_flag}" "${packaged_static_lib}"
            RESULT_VARIABLE strip_static_result
        )
        if(NOT strip_static_result EQUAL 0)
            message(FATAL_ERROR "failed to strip ${packaged_static_lib}")
        endif()
    endif()
endif()

if(DEFINED PSLOG_SHARED_LINK_NAME
   AND DEFINED PSLOG_SHARED_LIB_NAME
   AND NOT PSLOG_SHARED_LINK_NAME STREQUAL ""
   AND NOT PSLOG_SHARED_LINK_NAME STREQUAL PSLOG_SHARED_LIB_NAME)
    file(CREATE_LINK "${PSLOG_SHARED_LIB_NAME}"
         "${package_root}/lib/${PSLOG_SHARED_LINK_NAME}"
         SYMBOLIC)
endif()

if(DEFINED PSLOG_SHARED_SONAME
   AND DEFINED PSLOG_SHARED_LIB_NAME
   AND NOT PSLOG_SHARED_SONAME STREQUAL ""
   AND NOT PSLOG_SHARED_SONAME STREQUAL PSLOG_SHARED_LIB_NAME)
    file(CREATE_LINK "${PSLOG_SHARED_LIB_NAME}"
         "${package_root}/lib/${PSLOG_SHARED_SONAME}"
         SYMBOLIC)
endif()

file(COPY "${PSLOG_ROOT}/LICENSE" DESTINATION "${package_root}/share/doc/libpslog")
file(COPY "${PSLOG_ROOT}/README.md" DESTINATION "${package_root}/share/doc/libpslog")

file(MAKE_DIRECTORY "${PSLOG_ROOT}/dist")
set(archive_base "${PSLOG_ROOT}/dist/${archive_name}.tar")
set(archive "${archive_base}.gz")

find_program(PSLOG_TAR_BIN NAMES tar)
find_program(PSLOG_GZIP_BIN NAMES gzip)
if(NOT PSLOG_TAR_BIN)
    message(FATAL_ERROR "failed to find tar for archive creation")
endif()
if(NOT PSLOG_GZIP_BIN)
    message(FATAL_ERROR "failed to find gzip for archive creation")
endif()

file(REMOVE "${archive_base}" "${archive}")
execute_process(
    COMMAND "${PSLOG_TAR_BIN}" -cf "${archive_base}" --format=gnu --owner=0 --group=0 "${archive_name}"
    WORKING_DIRECTORY "${package_stage_root}"
    RESULT_VARIABLE tar_result
)
if(NOT tar_result EQUAL 0)
    message(FATAL_ERROR "failed to create package tar archive")
endif()

file(REMOVE "${archive}")
execute_process(
    COMMAND "${PSLOG_GZIP_BIN}" -9 -f "${archive_base}"
    RESULT_VARIABLE gzip_result
)
if(NOT gzip_result EQUAL 0)
    message(FATAL_ERROR "failed to gzip package archive")
endif()
