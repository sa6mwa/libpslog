if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

find_program(PSLOG_MAKE_BIN NAMES make REQUIRED)
foreach(make_args IN ITEMS "-s|help" "-n|build" "-n|clean")
    string(REPLACE "|" ";" make_args "${make_args}")
    execute_process(
        COMMAND ${CMAKE_COMMAND} -E env "PATH=/usr/bin:/bin"
            "${PSLOG_MAKE_BIN}" ${make_args}
        WORKING_DIRECTORY "${PSLOG_ROOT}"
        RESULT_VARIABLE make_result
        OUTPUT_VARIABLE make_output
        ERROR_VARIABLE make_error
    )
    if(NOT make_result EQUAL 0)
        message(FATAL_ERROR
            "C-only make command failed without LuaRocks: ${make_args}\n${make_output}\n${make_error}")
    endif()
    if(NOT make_error STREQUAL "")
        message(FATAL_ERROR
            "C-only make command emitted LuaRocks diagnostics: ${make_args}\n${make_error}")
    endif()
endforeach()
