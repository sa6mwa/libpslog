if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()
if(NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_BINARY_DIR is required")
endif()
if(NOT DEFINED PSLOG_VERSION)
    message(FATAL_ERROR "PSLOG_VERSION is required")
endif()

if(DEFINED PSLOG_DIST_DIR AND NOT PSLOG_DIST_DIR STREQUAL "")
    set(dist_dir "${PSLOG_DIST_DIR}")
else()
    set(dist_dir "${PSLOG_ROOT}/dist")
endif()
set(work_dir "${PSLOG_BINARY_DIR}/release-privacy-gate")

find_program(PSLOG_STRINGS_BIN NAMES strings)
if(NOT PSLOG_STRINGS_BIN)
    message(FATAL_ERROR "failed to find strings for release privacy gate")
endif()
find_program(PSLOG_GZIP_BIN NAMES gzip)

if(DEFINED PSLOG_PRIVACY_PATHS AND NOT PSLOG_PRIVACY_PATHS STREQUAL "")
    string(REPLACE "|" ";" privacy_paths "${PSLOG_PRIVACY_PATHS}")
else()
    set(checksums_path "${dist_dir}/libpslog-${PSLOG_VERSION}-CHECKSUMS")
    if(NOT EXISTS "${checksums_path}")
        message(FATAL_ERROR "missing release checksum manifest: ${checksums_path}")
    endif()
    set(privacy_paths "${checksums_path}")
    set(checksum_entries "${checksums_path}")
    file(READ "${checksums_path}" checksums_text)
    string(REPLACE "\n" ";" checksum_lines "${checksums_text}")
    foreach(checksum_line IN LISTS checksum_lines)
        string(STRIP "${checksum_line}" checksum_line)
        if(checksum_line STREQUAL "")
            continue()
        endif()
        if(NOT checksum_line MATCHES "^[0-9A-Fa-f]+[ \t]+\\*?(.+)$")
            message(FATAL_ERROR "invalid checksum manifest line: ${checksum_line}")
        endif()
        set(checksum_entry "${CMAKE_MATCH_1}")
        if(checksum_entry MATCHES "^/" OR checksum_entry MATCHES "\\.\\.")
            message(FATAL_ERROR
                "checksum manifest contains non-relocatable artifact path: ${checksum_entry}")
        endif()
        set(checksum_path "${dist_dir}/${checksum_entry}")
        list(APPEND privacy_paths "${checksum_path}")
        list(APPEND checksum_entries "${checksum_path}")
    endforeach()

    file(GLOB dist_paths "${dist_dir}/*")
    foreach(release_path IN LISTS dist_paths)
        if(IS_DIRECTORY "${release_path}")
            continue()
        endif()
        list(FIND checksum_entries "${release_path}" checksum_index)
        if(checksum_index EQUAL -1)
            message(FATAL_ERROR
                "dist file is not listed in checksum manifest: ${release_path}")
        endif()
    endforeach()
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
    if(DEFINED PSLOG_READELF AND NOT PSLOG_READELF STREQUAL "")
        set(PSLOG_READELF_BIN "${PSLOG_READELF}")
    else()
        find_program(PSLOG_READELF_BIN NAMES readelf)
    endif()
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
    if(DEFINED PSLOG_OTOOL AND NOT PSLOG_OTOOL STREQUAL "")
        set(${out_var} "${PSLOG_OTOOL}" PARENT_SCOPE)
        return()
    endif()

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

function(pslog_assert_darwin_relative_rpath rpath path)
    if(rpath MATCHES "^@loader_path(/|$)" OR rpath MATCHES "^@executable_path(/|$)")
        return()
    endif()
    message(FATAL_ERROR
        "Darwin artifact contains non-relocatable LC_RPATH '${rpath}': ${path}")
endfunction()

function(pslog_assert_darwin_dependency_path dependency path)
    if(dependency MATCHES "^@rpath/" OR dependency MATCHES "^@loader_path/" OR dependency MATCHES "^@executable_path/")
        return()
    endif()
    if(dependency MATCHES "^/usr/lib/" OR dependency MATCHES "^/System/Library/")
        return()
    endif()
    if(dependency MATCHES "^/")
        message(FATAL_ERROR
            "Darwin artifact contains non-system absolute dependency '${dependency}': ${path}")
    endif()
    message(FATAL_ERROR
        "Darwin artifact contains non-relocatable dependency '${dependency}': ${path}")
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

function(pslog_assert_release_manifest_entries_exist manifest_path context)
    get_filename_component(manifest_dir "${manifest_path}" DIRECTORY)
    file(READ "${manifest_path}" manifest_text)
    string(REPLACE "\n" ";" manifest_entries "${manifest_text}")
    foreach(manifest_entry IN LISTS manifest_entries)
        string(STRIP "${manifest_entry}" manifest_entry)
        if(manifest_entry STREQUAL "")
            continue()
        endif()
        if(manifest_entry MATCHES "^/" OR manifest_entry MATCHES "(^|/)\\.git(/|$)" OR manifest_entry MATCHES "\\.\\.")
            message(FATAL_ERROR
                "release manifest contains unsafe path '${manifest_entry}' in ${context}: ${manifest_path}")
        endif()
        if(NOT EXISTS "${manifest_dir}/${manifest_entry}")
            message(FATAL_ERROR
                "release manifest lists missing file '${manifest_entry}' in ${context}: ${manifest_path}")
        endif()
    endforeach()
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
        ERROR_VARIABLE otool_l_error
    )
    if(NOT otool_l_result EQUAL 0)
        message(FATAL_ERROR
            "failed to inspect Darwin load commands with otool -l: ${path}\n${otool_l_error}")
    endif()
    string(REGEX MATCHALL "path [^ \r\n]+ \\(offset [0-9]+\\)" rpath_lines "${otool_l_output}")
    foreach(rpath_line IN LISTS rpath_lines)
        string(REGEX REPLACE "^path ([^ \r\n]+) \\(offset [0-9]+\\)$" "\\1" rpath "${rpath_line}")
        pslog_assert_darwin_relative_rpath("${rpath}" "${path}")
    endforeach()

    execute_process(
        COMMAND "${PSLOG_OTOOL_BIN}" -D "${path}"
        RESULT_VARIABLE otool_d_result
        OUTPUT_VARIABLE install_name
        ERROR_VARIABLE otool_d_error
    )
    if(NOT otool_d_result EQUAL 0)
        message(FATAL_ERROR
            "failed to inspect Darwin install name with otool -D: ${path}\n${otool_d_error}")
    endif()
    if(NOT install_name MATCHES "[\r\n]")
        message(FATAL_ERROR
            "Darwin otool -D output did not include an install name: ${path}\n${install_name}")
    endif()
    if(NOT install_name MATCHES "(^|\n)[ \t]*(@rpath/[^\r\n]+)")
        message(FATAL_ERROR
            "Darwin shared library install name is not @rpath-relative: ${path}\n${install_name}")
    endif()
    set(install_name "${CMAKE_MATCH_2}")
    if(NOT install_name MATCHES "^@rpath/")
        message(FATAL_ERROR
            "Darwin shared library install name is not @rpath-relative: ${path}\n${install_name}")
    endif()

    execute_process(
        COMMAND "${PSLOG_OTOOL_BIN}" -L "${path}"
        RESULT_VARIABLE otool_dependency_result
        OUTPUT_VARIABLE otool_dependency_output
        ERROR_VARIABLE otool_dependency_error
    )
    if(NOT otool_dependency_result EQUAL 0)
        message(FATAL_ERROR
            "failed to inspect Darwin dependencies with otool -L: ${path}\n${otool_dependency_error}")
    endif()
    string(REPLACE "\n" ";" dependency_lines "${otool_dependency_output}")
    set(first_dependency_line TRUE)
    foreach(dependency_line IN LISTS dependency_lines)
        if(first_dependency_line)
            set(first_dependency_line FALSE)
            continue()
        endif()
        string(STRIP "${dependency_line}" dependency_line)
        if(dependency_line STREQUAL "")
            continue()
        endif()
        string(REGEX REPLACE "[ \t]+\\(.*$" "" dependency "${dependency_line}")
        pslog_assert_darwin_dependency_path("${dependency}" "${path}")
    endforeach()
endfunction()

function(pslog_check_extracted_tree tree_root context)
    file(GLOB_RECURSE extracted_entries LIST_DIRECTORIES false "${tree_root}/*")
    foreach(entry IN LISTS extracted_entries)
        pslog_assert_no_private_strings("${entry}" "${context}")
        get_filename_component(entry_name "${entry}" NAME)
        if(entry_name STREQUAL "RELEASE_MANIFEST")
            pslog_assert_release_manifest_entries_exist("${entry}" "${context}")
        elseif(entry_name MATCHES "\\.so(\\.|$)|\\.so$")
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
    elseif(artifact_name MATCHES "\\.dylib$")
        pslog_assert_darwin_paths("${artifact}")
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
