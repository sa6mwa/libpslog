if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()
if(NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_BINARY_DIR is required")
endif()
if(NOT DEFINED PSLOG_VERSION)
    message(FATAL_ERROR "PSLOG_VERSION is required")
endif()

set(dist_dir "${PSLOG_ROOT}/dist")
set(work_dir "${PSLOG_BINARY_DIR}/release-privacy-gate")

find_program(PSLOG_STRINGS_BIN NAMES strings)
if(NOT PSLOG_STRINGS_BIN)
    message(FATAL_ERROR "failed to find strings for release privacy gate")
endif()
find_program(PSLOG_GZIP_BIN NAMES gzip)

if(DEFINED PSLOG_PRIVACY_PATHS AND NOT PSLOG_PRIVACY_PATHS STREQUAL "")
    string(REPLACE "|" ";" privacy_paths "${PSLOG_PRIVACY_PATHS}")
else()
    file(GLOB privacy_paths
        "${dist_dir}/libpslog-${PSLOG_VERSION}-*.tar.gz"
        "${dist_dir}/lua-pslog-${PSLOG_VERSION}-*.rockspec"
        "${dist_dir}/lua-pslog-${PSLOG_VERSION}-*.src.rock"
        "${dist_dir}/pslog-${PSLOG_VERSION}.h.gz"
    )
endif()

list(LENGTH privacy_paths privacy_path_count)
if(privacy_path_count EQUAL 0)
    message(FATAL_ERROR "no release artifacts found for privacy gate in ${dist_dir}")
endif()

set(forbidden_literals "${PSLOG_ROOT}" "${PSLOG_BINARY_DIR}")
if(DEFINED ENV{HOME} AND NOT "$ENV{HOME}" STREQUAL "")
    list(APPEND forbidden_literals "$ENV{HOME}")
endif()
string(CONCAT pslog_home_dir_marker "/" "home/")
string(CONCAT pslog_users_dir_marker "/" "Users/")
set(pslog_home_dir_regex "(${pslog_home_dir_marker}|${pslog_users_dir_marker})")

function(pslog_assert_no_private_strings path context)
    if(IS_DIRECTORY "${path}")
        return()
    endif()

    execute_process(
        COMMAND "${PSLOG_STRINGS_BIN}" -a "${path}"
        RESULT_VARIABLE strings_result
        OUTPUT_VARIABLE strings_output
        ERROR_VARIABLE strings_error
    )
    if(NOT strings_result EQUAL 0)
        message(FATAL_ERROR
            "failed to scan ${context} for private paths: ${path}\n${strings_error}")
    endif()

    foreach(forbidden IN LISTS forbidden_literals)
        string(FIND "${strings_output}" "${forbidden}" forbidden_offset)
        if(NOT forbidden STREQUAL "" AND NOT forbidden_offset EQUAL -1)
            message(FATAL_ERROR
                "release artifact leaks private path '${forbidden}' in ${context}: ${path}")
        endif()
    endforeach()

    if(strings_output MATCHES "${pslog_home_dir_regex}")
        message(FATAL_ERROR
            "release artifact contains a hard-coded home directory in ${context}: ${path}")
    endif()
endfunction()

function(pslog_assert_linux_runpath path)
    find_program(PSLOG_READELF_BIN NAMES readelf)
    if(NOT PSLOG_READELF_BIN)
        return()
    endif()

    execute_process(
        COMMAND "${PSLOG_READELF_BIN}" -d "${path}"
        RESULT_VARIABLE readelf_result
        OUTPUT_VARIABLE readelf_output
        ERROR_QUIET
    )
    if(NOT readelf_result EQUAL 0)
        return()
    endif()

    if(readelf_output MATCHES "RPATH|RUNPATH")
        if(NOT readelf_output MATCHES "\\$ORIGIN")
            message(FATAL_ERROR
                "ELF shared library has RPATH/RUNPATH without $ORIGIN: ${path}\n${readelf_output}")
        endif()
        string(FIND "${readelf_output}" "${PSLOG_ROOT}" root_runpath_offset)
        string(FIND "${readelf_output}" "${PSLOG_BINARY_DIR}" build_runpath_offset)
        if(readelf_output MATCHES "${pslog_home_dir_regex}"
           OR NOT root_runpath_offset EQUAL -1
           OR NOT build_runpath_offset EQUAL -1)
            message(FATAL_ERROR
                "ELF shared library RPATH/RUNPATH leaks a private path: ${path}\n${readelf_output}")
        endif()
    endif()
endfunction()

function(pslog_find_darwin_otool out_var)
    set(_pslog_otool_paths "")
    if(DEFINED ENV{OSXCROSS_ROOT} AND NOT "$ENV{OSXCROSS_ROOT}" STREQUAL "")
        list(APPEND _pslog_otool_paths "$ENV{OSXCROSS_ROOT}/bin")
    endif()
    if(DEFINED ENV{HOME} AND NOT "$ENV{HOME}" STREQUAL "")
        list(APPEND _pslog_otool_paths "$ENV{HOME}/.local/cross/osxcross/bin")
    endif()
    find_program(_pslog_otool NAMES
        otool
        arm64-apple-darwin25-otool
        arm64-apple-darwin24-otool
        arm64-apple-darwin23-otool
        arm64-apple-darwin22-otool
        PATHS ${_pslog_otool_paths}
        NO_DEFAULT_PATH
    )
    if(NOT _pslog_otool)
        find_program(_pslog_otool NAMES
            otool
            arm64-apple-darwin25-otool
            arm64-apple-darwin24-otool
            arm64-apple-darwin23-otool
            arm64-apple-darwin22-otool
        )
    endif()
    set(${out_var} "${_pslog_otool}" PARENT_SCOPE)
endfunction()

function(pslog_extract_and_check_archive archive_path context)
    string(MAKE_C_IDENTIFIER "${archive_path}" archive_id)
    set(nested_extract_dir "${work_dir}/${archive_id}.extract")
    file(REMOVE_RECURSE "${nested_extract_dir}")
    file(MAKE_DIRECTORY "${nested_extract_dir}")
    execute_process(
        COMMAND "${CMAKE_COMMAND}" -E tar xf "${archive_path}"
        WORKING_DIRECTORY "${nested_extract_dir}"
        RESULT_VARIABLE nested_extract_result
        ERROR_VARIABLE nested_extract_error
    )
    if(NOT nested_extract_result EQUAL 0)
        message(FATAL_ERROR "failed to extract ${context}: ${archive_path}\n${nested_extract_error}")
    endif()
    pslog_check_extracted_tree("${nested_extract_dir}" "${context}")
endfunction()

function(pslog_assert_darwin_paths path)
    pslog_find_darwin_otool(PSLOG_OTOOL_BIN)
    if(NOT PSLOG_OTOOL_BIN)
        message(FATAL_ERROR "failed to find otool for Darwin release artifact: ${path}")
    endif()

    execute_process(
        COMMAND "${PSLOG_OTOOL_BIN}" -l "${path}"
        RESULT_VARIABLE otool_l_result
        OUTPUT_VARIABLE otool_l_output
        ERROR_QUIET
    )
    if(otool_l_result EQUAL 0 AND otool_l_output MATCHES "LC_RPATH")
        message(FATAL_ERROR
            "Darwin shared library contains LC_RPATH in release artifact: ${path}\n${otool_l_output}")
    endif()

    execute_process(
        COMMAND "${PSLOG_OTOOL_BIN}" -D "${path}"
        RESULT_VARIABLE otool_d_result
        OUTPUT_VARIABLE install_name
        ERROR_QUIET
    )
    if(otool_d_result EQUAL 0
       AND install_name MATCHES "[\r\n]")
        string(REGEX REPLACE "^[^\r\n]*[\r\n]+" "" install_name "${install_name}")
        string(STRIP "${install_name}" install_name)
        if(install_name MATCHES "^(/|${pslog_home_dir_regex}|\\$ORIGIN|@loader_path|@executable_path)")
            message(FATAL_ERROR
                "Darwin shared library install name contains an unexpected path: ${path}\n${install_name}")
        endif()
    endif()
endfunction()

function(pslog_check_extracted_tree tree_root context)
    file(GLOB_RECURSE extracted_entries LIST_DIRECTORIES false "${tree_root}/*")
    foreach(entry IN LISTS extracted_entries)
        pslog_assert_no_private_strings("${entry}" "${context}")
        get_filename_component(entry_name "${entry}" NAME)
        if(entry_name MATCHES "\\.so(\\.|$)|\\.so$")
            pslog_assert_linux_runpath("${entry}")
        elseif(entry_name MATCHES "\\.dylib$")
            pslog_assert_darwin_paths("${entry}")
        elseif(entry_name MATCHES "\\.tar\\.gz$|\\.src\\.rock$")
            pslog_extract_and_check_archive("${entry}" "${context}/${entry_name}")
        elseif(entry_name MATCHES "\\.h\\.gz$")
            if(NOT PSLOG_GZIP_BIN)
                message(FATAL_ERROR "failed to find gzip for release privacy gate")
            endif()
            set(header_path "${entry}.privacy-gate.txt")
            execute_process(
                COMMAND "${PSLOG_GZIP_BIN}" -cd "${entry}"
                OUTPUT_FILE "${header_path}"
                RESULT_VARIABLE nested_gzip_result
                ERROR_VARIABLE nested_gzip_error
            )
            if(NOT nested_gzip_result EQUAL 0)
                message(FATAL_ERROR "failed to decompress ${entry}\n${nested_gzip_error}")
            endif()
            pslog_assert_no_private_strings("${header_path}" "${context}/${entry_name}")
            file(REMOVE "${header_path}")
        endif()
    endforeach()
endfunction()

file(REMOVE_RECURSE "${work_dir}")
file(MAKE_DIRECTORY "${work_dir}")

foreach(artifact IN LISTS privacy_paths)
    if(NOT EXISTS "${artifact}")
        message(FATAL_ERROR "missing release artifact for privacy gate: ${artifact}")
    endif()

    pslog_assert_no_private_strings("${artifact}" "compressed artifact")
    get_filename_component(artifact_name "${artifact}" NAME)
    set(extract_dir "${work_dir}/${artifact_name}.extract")
    file(REMOVE_RECURSE "${extract_dir}")
    file(MAKE_DIRECTORY "${extract_dir}")

    if(artifact_name MATCHES "\\.(tar\\.gz|src\\.rock)$")
        pslog_extract_and_check_archive("${artifact}" "${artifact_name}")
    elseif(artifact_name MATCHES "\\.h\\.gz$")
        if(NOT PSLOG_GZIP_BIN)
            message(FATAL_ERROR "failed to find gzip for release privacy gate")
        endif()
        set(header_path "${extract_dir}/${artifact_name}.txt")
        execute_process(
            COMMAND "${PSLOG_GZIP_BIN}" -cd "${artifact}"
            OUTPUT_FILE "${header_path}"
            RESULT_VARIABLE gzip_result
            ERROR_VARIABLE gzip_error
        )
        if(NOT gzip_result EQUAL 0)
            message(FATAL_ERROR "failed to decompress ${artifact}\n${gzip_error}")
        endif()
        pslog_assert_no_private_strings("${header_path}" "${artifact_name}")
    endif()
endforeach()

file(REMOVE_RECURSE "${work_dir}")
