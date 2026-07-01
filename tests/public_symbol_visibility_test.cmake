if(NOT DEFINED PSLOG_PUBLIC_HEADER OR NOT EXISTS "${PSLOG_PUBLIC_HEADER}")
    message(FATAL_ERROR "PSLOG_PUBLIC_HEADER is required")
endif()
if(NOT DEFINED PSLOG_SHARED_LIB OR NOT EXISTS "${PSLOG_SHARED_LIB}")
    message(FATAL_ERROR "PSLOG_SHARED_LIB is required")
endif()
if(NOT DEFINED PSLOG_NM OR PSLOG_NM STREQUAL "")
    message(FATAL_ERROR "PSLOG_NM is required")
endif()

file(READ "${PSLOG_PUBLIC_HEADER}" pslog_public_content)
string(REGEX REPLACE "[\r\n]+[ \t]+" " " pslog_public_content
       "${pslog_public_content}")
string(REGEX MATCHALL "PSLOG_API[^;{]*[;(]" pslog_public_lines
       "${pslog_public_content}")
file(STRINGS "${PSLOG_PUBLIC_HEADER}" pslog_public_palette_lines
     REGEX "^PSLOG_API extern const pslog_palette ")
list(APPEND pslog_public_lines ${pslog_public_palette_lines})
set(pslog_expected_symbols "")
foreach(pslog_line IN LISTS pslog_public_lines)
    string(REGEX MATCH "pslog_[A-Za-z0-9_]+[ \t]*\\(" pslog_symbol_match
           "${pslog_line}")
    if(pslog_symbol_match)
        string(REGEX REPLACE "[ \t]*\\($" "" pslog_symbol
               "${pslog_symbol_match}")
        list(APPEND pslog_expected_symbols "${pslog_symbol}")
    else()
        string(REGEX MATCH
               "extern const pslog_palette[ \t]+(pslog_[A-Za-z0-9_]+)"
               pslog_palette_match "${pslog_line}")
        if(pslog_palette_match)
            string(REGEX REPLACE
                   ".*extern const pslog_palette[ \t]+(pslog_[A-Za-z0-9_]+).*"
                   "\\1" pslog_symbol "${pslog_palette_match}")
            list(APPEND pslog_expected_symbols "${pslog_symbol}")
        endif()
    endif()
endforeach()
list(REMOVE_DUPLICATES pslog_expected_symbols)

execute_process(
    COMMAND "${PSLOG_NM}" -D --defined-only "${PSLOG_SHARED_LIB}"
    RESULT_VARIABLE pslog_nm_result
    OUTPUT_VARIABLE pslog_nm_output
    ERROR_VARIABLE pslog_nm_error
)
if(NOT pslog_nm_result EQUAL 0)
    execute_process(
        COMMAND "${PSLOG_NM}" -gU "${PSLOG_SHARED_LIB}"
        RESULT_VARIABLE pslog_nm_result
        OUTPUT_VARIABLE pslog_nm_output
        ERROR_VARIABLE pslog_nm_error
    )
endif()
if(NOT pslog_nm_result EQUAL 0)
    message(FATAL_ERROR
        "Unable to inspect shared-library symbols with ${PSLOG_NM}:\n"
        "${pslog_nm_error}")
endif()

string(REGEX MATCHALL "pslog_[A-Za-z0-9_]+" pslog_exported_symbols
       "${pslog_nm_output}")
list(REMOVE_DUPLICATES pslog_exported_symbols)

set(pslog_extra_symbols "")
foreach(pslog_symbol IN LISTS pslog_exported_symbols)
    list(FIND pslog_expected_symbols "${pslog_symbol}" pslog_found_index)
    if(pslog_found_index EQUAL -1)
        list(APPEND pslog_extra_symbols "${pslog_symbol}")
    endif()
endforeach()

set(pslog_missing_symbols "")
foreach(pslog_symbol IN LISTS pslog_expected_symbols)
    list(FIND pslog_exported_symbols "${pslog_symbol}" pslog_found_index)
    if(pslog_found_index EQUAL -1)
        list(APPEND pslog_missing_symbols "${pslog_symbol}")
    endif()
endforeach()

if(pslog_extra_symbols OR pslog_missing_symbols)
    list(JOIN pslog_extra_symbols "\n  " pslog_extra_text)
    list(JOIN pslog_missing_symbols "\n  " pslog_missing_text)
    message(FATAL_ERROR
        "Shared-library public symbol surface does not match ${PSLOG_PUBLIC_HEADER}\n"
        "Extra exported symbols:\n  ${pslog_extra_text}\n"
        "Declared but missing symbols:\n  ${pslog_missing_text}")
endif()
