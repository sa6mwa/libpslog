if(NOT DEFINED PSLOG_BINARY_DIR OR NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_BINARY_DIR and PSLOG_ROOT are required")
endif()

set(test_root "${PSLOG_BINARY_DIR}/dependency-cache-test")
file(REMOVE_RECURSE "${test_root}")
file(MAKE_DIRECTORY "${test_root}")
unset(CPKT_DEPENDENCY_CACHE)
unset(CPKT_DEPENDENCY_CACHE CACHE)
set(ENV{CPKT_DEPENDENCY_CACHE} "")
set(ENV{XDG_CACHE_HOME} "")
set(ENV{HOME} "")
include("${PSLOG_ROOT}/cmake/pslog_dependencies.cmake")
if(DEFINED CPKT_DEPENDENCY_CACHE AND NOT CPKT_DEPENDENCY_CACHE STREQUAL "")
    message(FATAL_ERROR "dependency cache was resolved before an archive was requested")
endif()
set(CPKT_DEPENDENCY_CACHE "${test_root}/shared-cache")
set(PSLOG_DEPENDENCY_TEST_ALLOW_FILE_URL TRUE)
set(source_archive "${test_root}/fixture.tar.gz")
file(WRITE "${source_archive}" "verified fixture archive\n")
file(SHA256 "${source_archive}" fixture_sha256)
file(READ "${PSLOG_ROOT}/cmake/pslog_dependencies.cmake" dependency_helper)
if(NOT dependency_helper MATCHES "file\\(LOCK[ \\t]+\"\\$\\{lock_path\\}\"")
    message(FATAL_ERROR "dependency cache helper does not serialize archive writers per digest")
endif()
if(NOT dependency_helper MATCHES "file\\(RENAME[ \\t]+\"\\$\\{temp_path\\}\"[ \\t]+\"\\$\\{archive_path\\}\"")
    message(FATAL_ERROR "dependency cache helper does not atomically publish verified archives")
endif()
if(dependency_helper MATCHES "EXPECTED_HASH" OR
   NOT dependency_helper MATCHES "foreach\\(download_attempt[ \\t]+RANGE[ \\t]+1[ \\t]+3\\)" OR
   NOT dependency_helper MATCHES "file\\(SHA256[ \\t]+\"\\$\\{temp_path\\}\"[ \\t]+downloaded_digest\\)")
    message(FATAL_ERROR "dependency downloads must retry and explicitly verify temporary archives before publication")
endif()
if(NOT dependency_helper MATCHES "set\\(stage_lock_path[ \\t]+\"\\$\\{CMAKE_SOURCE_DIR\\}/\\.cache/deps/locks/" OR
   NOT dependency_helper MATCHES "file\\(LOCK[ \\t]+\"\\$\\{stage_lock_path\\}\"[ \\t]+GUARD[ \\t]+FUNCTION" OR
   NOT dependency_helper MATCHES "file\\(SHA256[ \\t]+\"\\$\\{ARG_ARCHIVE\\}\"[ \\t]+staged_archive_digest")
    message(FATAL_ERROR "dependency staging does not hold a per-stage lock through extraction")
endif()
pslog_acquire_verified_archive(
    COMPONENT fixture URL "file://${source_archive}" SHA256 "${fixture_sha256}" OUTPUT first_archive)
if(NOT EXISTS "${first_archive}")
    message(FATAL_ERROR "cache miss did not publish an archive")
endif()
file(SHA256 "${first_archive}" first_digest)
if(NOT first_digest STREQUAL fixture_sha256)
    message(FATAL_ERROR "published archive is not checksum verified")
endif()

file(REMOVE "${source_archive}")
pslog_acquire_verified_archive(
    COMPONENT fixture URL "file://${source_archive}" SHA256 "${fixture_sha256}" OUTPUT offline_archive)
if(NOT offline_archive STREQUAL first_archive)
    message(FATAL_ERROR "offline cache hit did not reuse the verified archive")
endif()

file(WRITE "${first_archive}" "corrupt cache entry\n")
file(WRITE "${source_archive}" "verified fixture archive\n")
pslog_acquire_verified_archive(
    COMPONENT fixture URL "file://${source_archive}" SHA256 "${fixture_sha256}" OUTPUT repaired_archive)
file(SHA256 "${repaired_archive}" repaired_digest)
if(NOT repaired_digest STREQUAL fixture_sha256)
    message(FATAL_ERROR "corrupt cache entry was not rejected and repaired")
endif()

set(stage_failure_script "${test_root}/stage-failure.cmake")
file(WRITE "${stage_failure_script}"
"set(CMAKE_SOURCE_DIR \"${test_root}/stage-project\")\n"
"include(\"${PSLOG_ROOT}/cmake/pslog_dependencies.cmake\")\n"
"pslog_stage_dependency_archive(COMPONENT fixture ARCHIVE \"${repaired_archive}\" SHA256 \"0000000000000000000000000000000000000000000000000000000000000000\" OUTPUT staged)\n")
execute_process(
    COMMAND "${CMAKE_COMMAND}" -P "${stage_failure_script}"
    RESULT_VARIABLE stage_failure_result
    OUTPUT_VARIABLE stage_failure_output
    ERROR_VARIABLE stage_failure_error
)
if(stage_failure_result EQUAL 0 OR
   NOT stage_failure_error MATCHES "dependency archive checksum mismatch before staging")
    message(FATAL_ERROR
        "dependency staging accepted a mismatched archive\n"
        "stdout:\n${stage_failure_output}\n"
        "stderr:\n${stage_failure_error}")
endif()
file(REMOVE_RECURSE "${test_root}")
