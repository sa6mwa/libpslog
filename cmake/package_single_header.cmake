if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()
if(NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_BINARY_DIR is required")
endif()
if(NOT DEFINED PSLOG_VERSION)
    message(FATAL_ERROR "PSLOG_VERSION is required")
endif()
if(NOT DEFINED PSLOG_PUBLIC_HEADER)
    message(FATAL_ERROR "PSLOG_PUBLIC_HEADER is required")
endif()
if(NOT DEFINED PSLOG_PUBLIC_VERSION_HEADER)
    message(FATAL_ERROR "PSLOG_PUBLIC_VERSION_HEADER is required")
endif()
if(NOT DEFINED PSLOG_INTERNAL_HEADER)
    message(FATAL_ERROR "PSLOG_INTERNAL_HEADER is required")
endif()
if(NOT DEFINED PSLOG_SINGLE_HEADER_BUILD)
    message(FATAL_ERROR "PSLOG_SINGLE_HEADER_BUILD is required")
endif()
if(NOT DEFINED PSLOG_SINGLE_HEADER_BUILD_ALIAS)
    message(FATAL_ERROR "PSLOG_SINGLE_HEADER_BUILD_ALIAS is required")
endif()
if(NOT DEFINED PSLOG_SINGLE_HEADER_DIST)
    message(FATAL_ERROR "PSLOG_SINGLE_HEADER_DIST is required")
endif()
if(NOT DEFINED PSLOG_SINGLE_HEADER_DIST_GZ)
    message(FATAL_ERROR "PSLOG_SINGLE_HEADER_DIST_GZ is required")
endif()
if(NOT DEFINED PSLOG_SINGLE_HEADER_SOURCES)
    message(FATAL_ERROR "PSLOG_SINGLE_HEADER_SOURCES is required")
endif()
if(NOT DEFINED PSLOG_CLANG_FORMAT_BIN OR PSLOG_CLANG_FORMAT_BIN STREQUAL "")
    message(FATAL_ERROR "PSLOG_CLANG_FORMAT_BIN is required to format the single-header artifact")
endif()

string(REPLACE "|" ";" PSLOG_SINGLE_HEADER_SOURCES "${PSLOG_SINGLE_HEADER_SOURCES}")

function(pslog_read_normalized out_var path)
    file(READ "${path}" content)
    string(REPLACE "\r\n" "\n" content "${content}")
    set(${out_var} "${content}" PARENT_SCOPE)
endfunction()

function(pslog_strip_guard out_var guard_macro content)
    set(stripped "${content}")
    string(REGEX REPLACE "^#ifndef ${guard_macro}\n#define ${guard_macro}\n\n?" "" stripped "${stripped}")
    string(REGEX REPLACE "\n#endif[ \t]*\n?$" "\n" stripped "${stripped}")
    set(${out_var} "${stripped}" PARENT_SCOPE)
endfunction()

pslog_read_normalized(version_header_raw "${PSLOG_PUBLIC_VERSION_HEADER}")
pslog_strip_guard(version_header "PSLOG_VERSION_H" "${version_header_raw}")

pslog_read_normalized(public_header_raw "${PSLOG_PUBLIC_HEADER}")
string(REPLACE "#include \"pslog_version.h\"\n" "" public_header_raw "${public_header_raw}")
pslog_strip_guard(public_header "PSLOG_H" "${public_header_raw}")

pslog_read_normalized(internal_header_raw "${PSLOG_INTERNAL_HEADER}")
string(REPLACE "#include \"pslog.h\"\n\n" "" internal_header_raw "${internal_header_raw}")
string(REPLACE "#include \"pslog.h\"\n" "" internal_header_raw "${internal_header_raw}")
pslog_strip_guard(internal_header "PSLOG_INTERNAL_H" "${internal_header_raw}")

pslog_read_normalized(license_text "${PSLOG_ROOT}/LICENSE")
string(REPLACE "\"" "\\\"" license_text "${license_text}")

pslog_read_normalized(preamble_template "${PSLOG_ROOT}/cmake/pslog_single_header_preamble.in")
get_filename_component(pslog_single_header_filename "${PSLOG_SINGLE_HEADER_DIST}" NAME)
set(PSLOG_LICENSE_TEXT "${license_text}")
set(PSLOG_SINGLE_HEADER_FILENAME "${pslog_single_header_filename}")
string(CONFIGURE "${preamble_template}" preamble @ONLY)

set(single_header "${preamble}\n\n#ifndef PSLOG_H\n#define PSLOG_H\n\n")
string(APPEND single_header "/* Inlined version header */\n")
string(APPEND single_header "${version_header}\n")
string(APPEND single_header "/* Public API */\n")
string(APPEND single_header "${public_header}\n")
string(APPEND single_header "#endif /* PSLOG_H */\n\n")
string(APPEND single_header "#ifdef PSLOG_IMPLEMENTATION\n")
string(APPEND single_header "#ifndef PSLOG_IMPLEMENTATION_ONCE\n")
string(APPEND single_header "#define PSLOG_IMPLEMENTATION_ONCE\n\n")
string(APPEND single_header "/* Internal definitions */\n")
string(APPEND single_header "${internal_header}\n")

foreach(source_path IN LISTS PSLOG_SINGLE_HEADER_SOURCES)
    pslog_read_normalized(source_content "${source_path}")
    string(REPLACE "#include \"pslog_internal.h\"\n\n" "" source_content "${source_content}")
    string(REPLACE "#include \"pslog_internal.h\"\n" "" source_content "${source_content}")
    string(REPLACE "#include \"pslog.h\"\n\n" "" source_content "${source_content}")
    string(REPLACE "#include \"pslog.h\"\n" "" source_content "${source_content}")
    get_filename_component(source_name "${source_path}" NAME)
    string(APPEND single_header "\n/* Source: ${source_name} */\n")
    string(APPEND single_header "${source_content}\n")
endforeach()

string(APPEND single_header "\n#endif /* PSLOG_IMPLEMENTATION_ONCE */\n")
string(APPEND single_header "#endif /* PSLOG_IMPLEMENTATION */\n")

get_filename_component(build_dir "${PSLOG_SINGLE_HEADER_BUILD}" DIRECTORY)
get_filename_component(build_alias_dir "${PSLOG_SINGLE_HEADER_BUILD_ALIAS}" DIRECTORY)
get_filename_component(dist_dir "${PSLOG_SINGLE_HEADER_DIST}" DIRECTORY)
file(MAKE_DIRECTORY "${build_dir}" "${build_alias_dir}" "${dist_dir}")
file(WRITE "${PSLOG_SINGLE_HEADER_BUILD}" "${single_header}")
file(WRITE "${PSLOG_SINGLE_HEADER_BUILD_ALIAS}" "${single_header}")
file(WRITE "${PSLOG_SINGLE_HEADER_DIST}" "${single_header}")

execute_process(
    COMMAND "${PSLOG_CLANG_FORMAT_BIN}"
        -i
        "${PSLOG_SINGLE_HEADER_BUILD}"
        "${PSLOG_SINGLE_HEADER_BUILD_ALIAS}"
        "${PSLOG_SINGLE_HEADER_DIST}"
    RESULT_VARIABLE clang_format_result
)
if(NOT clang_format_result EQUAL 0)
    message(FATAL_ERROR "failed to clang-format single-header artifact")
endif()

find_program(PSLOG_GZIP_BIN NAMES gzip)
if(NOT PSLOG_GZIP_BIN)
    message(FATAL_ERROR "failed to find gzip for single-header artifact creation")
endif()

file(REMOVE "${PSLOG_SINGLE_HEADER_DIST_GZ}")
execute_process(
    COMMAND "${PSLOG_GZIP_BIN}" -9 -c "${PSLOG_SINGLE_HEADER_DIST}"
    OUTPUT_FILE "${PSLOG_SINGLE_HEADER_DIST_GZ}"
    RESULT_VARIABLE gzip_result
)
if(NOT gzip_result EQUAL 0)
    message(FATAL_ERROR "failed to gzip single-header artifact")
endif()
