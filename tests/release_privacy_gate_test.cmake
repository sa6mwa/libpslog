if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()
if(NOT DEFINED PSLOG_BINARY_DIR)
    message(FATAL_ERROR "PSLOG_BINARY_DIR is required")
endif()

set(test_root "${PSLOG_BINARY_DIR}/release-privacy-gate-test")
set(clean_artifact "${test_root}/clean.rockspec")
set(leaky_artifact "${test_root}/leaky.rockspec")
set(generic_home_artifact "${test_root}/generic-home.rockspec")

file(REMOVE_RECURSE "${test_root}")
file(MAKE_DIRECTORY "${test_root}")

file(WRITE "${clean_artifact}"
    "package = \"lua-pslog\"\n"
    "source = { url = \"file://lua-pslog-9.9.9.tar.gz\" }\n")
execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_BINARY_DIR=${PSLOG_BINARY_DIR}
        -DPSLOG_VERSION=9.9.9
        -DPSLOG_PRIVACY_PATHS=${clean_artifact}
        -P ${PSLOG_ROOT}/cmake/check_release_privacy.cmake
    RESULT_VARIABLE clean_result
    OUTPUT_VARIABLE clean_output
    ERROR_VARIABLE clean_error
)
if(NOT clean_result EQUAL 0)
    message(FATAL_ERROR
        "release privacy gate rejected clean relative source URL\n"
        "stdout:\n${clean_output}\n"
        "stderr:\n${clean_error}")
endif()

file(WRITE "${leaky_artifact}"
    "package = \"lua-pslog\"\n"
    "source = { url = \"git+file://${PSLOG_ROOT}\" }\n")
execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_BINARY_DIR=${PSLOG_BINARY_DIR}
        -DPSLOG_VERSION=9.9.9
        -DPSLOG_PRIVACY_PATHS=${leaky_artifact}
        -P ${PSLOG_ROOT}/cmake/check_release_privacy.cmake
    RESULT_VARIABLE leaky_result
    OUTPUT_VARIABLE leaky_output
    ERROR_VARIABLE leaky_error
)
if(leaky_result EQUAL 0)
    message(FATAL_ERROR "release privacy gate accepted an artifact containing the repo path")
endif()
if(NOT "${leaky_output}${leaky_error}" MATCHES "leaks private path")
    message(FATAL_ERROR
        "release privacy gate failure did not explain the private path leak\n"
        "stdout:\n${leaky_output}\n"
        "stderr:\n${leaky_error}")
endif()

string(CONCAT generic_home_url "file:///" "home/builder/src")
file(WRITE "${generic_home_artifact}"
    "package = \"lua-pslog\"\n"
    "source = { url = \"${generic_home_url}\" }\n")
execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_BINARY_DIR=${PSLOG_BINARY_DIR}
        -DPSLOG_VERSION=9.9.9
        -DPSLOG_PRIVACY_PATHS=${generic_home_artifact}
        -P ${PSLOG_ROOT}/cmake/check_release_privacy.cmake
    RESULT_VARIABLE generic_home_result
    OUTPUT_VARIABLE generic_home_output
    ERROR_VARIABLE generic_home_error
)
if(generic_home_result EQUAL 0)
    message(FATAL_ERROR "release privacy gate accepted a formatted /home path leak")
endif()
if(NOT "${generic_home_output}${generic_home_error}" MATCHES "hard-coded home directory")
    message(FATAL_ERROR
        "release privacy gate failure did not explain the generic home path leak\n"
        "stdout:\n${generic_home_output}\n"
        "stderr:\n${generic_home_error}")
endif()
