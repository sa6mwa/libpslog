foreach(required IN ITEMS
        PSLOG_BINARY PSLOG_READELF PSLOG_BOOTLIN_INTERPRETER
        PSLOG_BOOTLIN_SYSROOT PSLOG_BOOTLIN_RPATH)
    if(NOT DEFINED ${required} OR "${${required}}" STREQUAL "")
        message(FATAL_ERROR "${required} is required")
    endif()
endforeach()

execute_process(
    COMMAND "${PSLOG_READELF}" -l "${PSLOG_BINARY}"
    RESULT_VARIABLE interpreter_result
    OUTPUT_VARIABLE interpreter_output
    ERROR_VARIABLE interpreter_error
)
if(NOT interpreter_result EQUAL 0)
    message(FATAL_ERROR "Unable to inspect ELF interpreter: ${interpreter_error}")
endif()
if(NOT interpreter_output MATCHES "Requesting program interpreter: ${PSLOG_BOOTLIN_INTERPRETER}")
    message(FATAL_ERROR
        "Non-shipped ELF executable does not pin the selected Bootlin interpreter.\n${interpreter_output}")
endif()

execute_process(
    COMMAND "${PSLOG_READELF}" -d "${PSLOG_BINARY}"
    RESULT_VARIABLE dynamic_result
    OUTPUT_VARIABLE dynamic_output
    ERROR_VARIABLE dynamic_error
)
if(NOT dynamic_result EQUAL 0)
    message(FATAL_ERROR "Unable to inspect ELF dynamic metadata: ${dynamic_error}")
endif()
string(FIND "${dynamic_output}" "${PSLOG_BOOTLIN_RPATH}" rpath_offset)
if(rpath_offset EQUAL -1)
    message(FATAL_ERROR
        "Non-shipped ELF executable does not carry its Bootlin DT_RPATH.\n${dynamic_output}")
endif()

execute_process(
    COMMAND "${PSLOG_BOOTLIN_INTERPRETER}" --list "${PSLOG_BINARY}"
    RESULT_VARIABLE runtime_result
    OUTPUT_VARIABLE runtime_output
    ERROR_VARIABLE runtime_error
)
if(NOT runtime_result EQUAL 0)
    message(FATAL_ERROR "Bootlin loader could not resolve ${PSLOG_BINARY}: ${runtime_error}")
endif()
string(FIND "${runtime_output}${runtime_error}" "${PSLOG_BOOTLIN_SYSROOT}/" sysroot_runtime_offset)
if(sysroot_runtime_offset EQUAL -1)
    message(FATAL_ERROR
        "Bootlin loader resolution did not use its selected sysroot.\n${runtime_output}${runtime_error}")
endif()

execute_process(
    COMMAND "${PSLOG_BINARY}"
    RESULT_VARIABLE direct_run_result
)
if(NOT direct_run_result EQUAL 0)
    message(FATAL_ERROR "Bootlin-pinned executable did not run directly: ${direct_run_result}")
endif()
