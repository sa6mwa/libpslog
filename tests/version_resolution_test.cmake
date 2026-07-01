if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

find_program(PSLOG_GIT_BIN NAMES git)
if(NOT PSLOG_GIT_BIN)
    message(FATAL_ERROR "git is required for version resolution tests")
endif()

set(test_root "${PSLOG_BINARY_DIR}/version-resolution-test")
set(repo_dir "${test_root}/repo")
set(tagged_output "${test_root}/tagged.txt")
set(untagged_output "${test_root}/untagged.txt")

file(REMOVE_RECURSE "${test_root}")
file(MAKE_DIRECTORY "${repo_dir}")
file(WRITE "${repo_dir}/README.md" "version test\n")

execute_process(
    COMMAND "${PSLOG_GIT_BIN}" init
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_init_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_init_result EQUAL 0)
    message(FATAL_ERROR "failed to initialize temporary git repository")
endif()

execute_process(
    COMMAND "${PSLOG_GIT_BIN}" config user.email test@example.com
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_email_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_email_result EQUAL 0)
    message(FATAL_ERROR "failed to configure temporary git email")
endif()

execute_process(
    COMMAND "${PSLOG_GIT_BIN}" config user.name "libpslog version test"
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_name_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_name_result EQUAL 0)
    message(FATAL_ERROR "failed to configure temporary git name")
endif()

execute_process(
    COMMAND "${PSLOG_GIT_BIN}" add README.md
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_add_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_add_result EQUAL 0)
    message(FATAL_ERROR "failed to stage temporary repository contents")
endif()

execute_process(
    COMMAND "${PSLOG_GIT_BIN}" commit -m "initial"
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_commit_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_commit_result EQUAL 0)
    message(FATAL_ERROR "failed to create initial temporary repository commit")
endif()

execute_process(
    COMMAND "${PSLOG_GIT_BIN}" tag v1.2.3
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_tag_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_tag_result EQUAL 0)
    message(FATAL_ERROR "failed to create temporary repository tag")
endif()

execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_VERSION_SOURCE_DIR=${repo_dir}
        -DPSLOG_VERSION_PROBE_OUTPUT=${tagged_output}
        -P ${PSLOG_ROOT}/tests/version_resolution_probe.cmake
    RESULT_VARIABLE tagged_probe_result
)
if(NOT tagged_probe_result EQUAL 0)
    message(FATAL_ERROR "failed to probe tagged version resolution")
endif()

file(READ "${tagged_output}" tagged_probe_contents)
string(STRIP "${tagged_probe_contents}" tagged_probe_contents)
if(NOT tagged_probe_contents STREQUAL "1.2.3|1.2.3|1|2|3")
    message(FATAL_ERROR
        "expected exact v-tag to resolve to 1.2.3, got '${tagged_probe_contents}'")
endif()

set(build_metadata_output "${test_root}/build-metadata.txt")
execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_VERSION_SOURCE_DIR=${repo_dir}
        -DPSLOG_VERSION_OVERRIDE=1.2.3+build.1
        -DPSLOG_VERSION_PROBE_OUTPUT=${build_metadata_output}
        -P ${PSLOG_ROOT}/tests/version_resolution_probe.cmake
    RESULT_VARIABLE build_metadata_probe_result
    ERROR_VARIABLE build_metadata_probe_error
)
if(NOT build_metadata_probe_result EQUAL 0)
    message(FATAL_ERROR
        "build metadata version override failed unexpectedly: ${build_metadata_probe_error}")
endif()

file(READ "${build_metadata_output}" build_metadata_probe_contents)
string(STRIP "${build_metadata_probe_contents}" build_metadata_probe_contents)
if(NOT build_metadata_probe_contents STREQUAL "1.2.3+build.1|1.2.3|1|2|3")
    message(FATAL_ERROR
        "expected build metadata override to resolve to 1.2.3+build.1, got '${build_metadata_probe_contents}'")
endif()

set(nested_source_dir "${repo_dir}/nested-source-with-version")
set(nested_source_version_output "${test_root}/nested-source-version.txt")
file(MAKE_DIRECTORY "${nested_source_dir}")
file(WRITE "${nested_source_dir}/VERSION" "4.5.6\n")
execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_VERSION_SOURCE_DIR=${nested_source_dir}
        -DPSLOG_VERSION_PROBE_OUTPUT=${nested_source_version_output}
        -P ${PSLOG_ROOT}/tests/version_resolution_probe.cmake
    RESULT_VARIABLE nested_source_version_probe_result
)
if(NOT nested_source_version_probe_result EQUAL 0)
    message(FATAL_ERROR "failed to probe nested non-git VERSION resolution")
endif()

file(READ "${nested_source_version_output}" nested_source_version_probe_contents)
string(STRIP "${nested_source_version_probe_contents}" nested_source_version_probe_contents)
if(NOT nested_source_version_probe_contents STREQUAL "4.5.6|4.5.6|4|5|6")
    message(FATAL_ERROR
        "expected nested non-git VERSION file to resolve to 4.5.6, got '${nested_source_version_probe_contents}'")
endif()

file(APPEND "${repo_dir}/README.md" "next commit\n")
execute_process(
    COMMAND "${PSLOG_GIT_BIN}" add README.md
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_add_second_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_add_second_result EQUAL 0)
    message(FATAL_ERROR "failed to stage second temporary repository commit")
endif()

execute_process(
    COMMAND "${PSLOG_GIT_BIN}" commit -m "second"
    WORKING_DIRECTORY "${repo_dir}"
    RESULT_VARIABLE git_second_commit_result
    OUTPUT_QUIET
    ERROR_QUIET
)
if(NOT git_second_commit_result EQUAL 0)
    message(FATAL_ERROR "failed to create second temporary repository commit")
endif()

execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_VERSION_SOURCE_DIR=${repo_dir}
        -DPSLOG_VERSION_PROBE_OUTPUT=${untagged_output}
        -P ${PSLOG_ROOT}/tests/version_resolution_probe.cmake
    RESULT_VARIABLE untagged_probe_result
)
if(NOT untagged_probe_result EQUAL 0)
    message(FATAL_ERROR "failed to probe untagged version resolution")
endif()

file(READ "${untagged_output}" untagged_probe_contents)
string(STRIP "${untagged_probe_contents}" untagged_probe_contents)
if(NOT untagged_probe_contents STREQUAL "0.0.0|0.0.0|0|0|0")
    message(FATAL_ERROR
        "expected untagged commit to resolve to 0.0.0, got '${untagged_probe_contents}'")
endif()

set(source_dir "${test_root}/source-with-version")
set(source_version_output "${test_root}/source-version.txt")
file(MAKE_DIRECTORY "${source_dir}")
file(WRITE "${source_dir}/VERSION" "4.5.6\n")
execute_process(
    COMMAND "${CMAKE_COMMAND}"
        -DPSLOG_ROOT=${PSLOG_ROOT}
        -DPSLOG_VERSION_SOURCE_DIR=${source_dir}
        -DPSLOG_VERSION_PROBE_OUTPUT=${source_version_output}
        -P ${PSLOG_ROOT}/tests/version_resolution_probe.cmake
    RESULT_VARIABLE source_version_probe_result
)
if(NOT source_version_probe_result EQUAL 0)
    message(FATAL_ERROR "failed to probe non-git VERSION resolution")
endif()

file(READ "${source_version_output}" source_version_probe_contents)
string(STRIP "${source_version_probe_contents}" source_version_probe_contents)
if(NOT source_version_probe_contents STREQUAL "4.5.6|4.5.6|4|5|6")
    message(FATAL_ERROR
        "expected non-git VERSION file to resolve to 4.5.6, got '${source_version_probe_contents}'")
endif()
