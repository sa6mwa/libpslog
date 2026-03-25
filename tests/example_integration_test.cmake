if(NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_BINARY_DIR is required")
endif()
if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()
if(NOT DEFINED PSLOG_C_COMPILER OR PSLOG_C_COMPILER STREQUAL "")
    message(FATAL_ERROR "PSLOG_C_COMPILER is required")
endif()
if(NOT DEFINED PSLOG_THREAD_LIBS)
    set(PSLOG_THREAD_LIBS "")
endif()
if(NOT DEFINED PSLOG_CROSSCOMPILING)
    set(PSLOG_CROSSCOMPILING FALSE)
endif()
if(NOT DEFINED PSLOG_CROSSCOMPILING_EMULATOR)
    set(PSLOG_CROSSCOMPILING_EMULATOR "")
endif()

set(test_root "${PSLOG_BINARY_DIR}/example-integration-test")
set(example_source "${PSLOG_ROOT}/examples/example.c")
set(library_binary "${test_root}/example-library")
set(single_header_binary "${test_root}/example-single-header")

file(REMOVE_RECURSE "${test_root}")
file(MAKE_DIRECTORY "${test_root}")

set(example_run_command)
if(PSLOG_CROSSCOMPILING)
    string(REPLACE "|" ";" pslog_cross_emulator "${PSLOG_CROSSCOMPILING_EMULATOR}")
    list(LENGTH pslog_cross_emulator pslog_cross_emulator_len)
    if(pslog_cross_emulator_len EQUAL 0)
        message(STATUS "Skipping example execution for cross build without emulator")
    else()
        set(example_run_command ${pslog_cross_emulator})
    endif()
endif()

execute_process(
    COMMAND "${CMAKE_COMMAND}" --build "${PSLOG_BINARY_DIR}" --target pslog_static package-single-header
    RESULT_VARIABLE build_result
)
if(NOT build_result EQUAL 0)
    message(FATAL_ERROR "failed to build pslog_static and package-single-header targets")
endif()

separate_arguments(pslog_thread_libs NATIVE_COMMAND "${PSLOG_THREAD_LIBS}")

set(library_compile_command
    "${PSLOG_C_COMPILER}"
    "-I${PSLOG_BINARY_DIR}/generated/include"
    "-I${PSLOG_ROOT}/include"
    "-o" "${library_binary}"
    "${example_source}"
    "${PSLOG_BINARY_DIR}/libpslog.a"
)
list(APPEND library_compile_command ${pslog_thread_libs})

execute_process(
    COMMAND ${library_compile_command}
    RESULT_VARIABLE library_compile_result
    OUTPUT_VARIABLE library_compile_stdout
    ERROR_VARIABLE library_compile_stderr
)
if(NOT library_compile_result EQUAL 0)
    message(FATAL_ERROR
        "failed to compile examples/example.c against libpslog.a\n"
        "stdout:\n${library_compile_stdout}\n"
        "stderr:\n${library_compile_stderr}")
endif()

list(LENGTH example_run_command example_run_command_len)
if(example_run_command_len EQUAL 0)
    execute_process(
        COMMAND "${library_binary}"
        RESULT_VARIABLE library_run_result
    )
    if(NOT library_run_result EQUAL 0)
        message(FATAL_ERROR "compiled library-mode example exited with status ${library_run_result}")
    endif()
else()
    execute_process(
        COMMAND ${example_run_command} "${library_binary}"
        RESULT_VARIABLE library_run_result
    )
    if(NOT library_run_result EQUAL 0)
        message(FATAL_ERROR "compiled library-mode example exited with status ${library_run_result}")
    endif()
endif()

set(single_header_compile_command
    "${PSLOG_C_COMPILER}"
    "-DPSLOG_EXAMPLE_SINGLE_HEADER=1"
    "-I${PSLOG_BINARY_DIR}/generated/include"
    "-o" "${single_header_binary}"
    "${example_source}"
)
list(APPEND single_header_compile_command ${pslog_thread_libs})

execute_process(
    COMMAND ${single_header_compile_command}
    RESULT_VARIABLE single_header_compile_result
    OUTPUT_VARIABLE single_header_compile_stdout
    ERROR_VARIABLE single_header_compile_stderr
)
if(NOT single_header_compile_result EQUAL 0)
    message(FATAL_ERROR
        "failed to compile examples/example.c against generated single header\n"
        "stdout:\n${single_header_compile_stdout}\n"
        "stderr:\n${single_header_compile_stderr}")
endif()

if(example_run_command_len EQUAL 0)
    execute_process(
        COMMAND "${single_header_binary}"
        RESULT_VARIABLE single_header_run_result
    )
    if(NOT single_header_run_result EQUAL 0)
        message(FATAL_ERROR "compiled single-header example exited with status ${single_header_run_result}")
    endif()
else()
    execute_process(
        COMMAND ${example_run_command} "${single_header_binary}"
        RESULT_VARIABLE single_header_run_result
    )
    if(NOT single_header_run_result EQUAL 0)
        message(FATAL_ERROR "compiled single-header example exited with status ${single_header_run_result}")
    endif()
endif()

file(REMOVE_RECURSE "${test_root}")
