if(NOT DEFINED PSLOG_BINARY_DIR OR NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_BINARY_DIR and PSLOG_ROOT are required")
endif()
cmake_policy(SET CMP0007 NEW)

execute_process(
    COMMAND make release-lua-artifacts
    WORKING_DIRECTORY "${PSLOG_ROOT}"
    RESULT_VARIABLE artifact_result
    OUTPUT_VARIABLE artifact_output
    ERROR_VARIABLE artifact_error
)
if(NOT artifact_result EQUAL 0)
    message(FATAL_ERROR
        "failed to build Lua release artifacts\n${artifact_output}${artifact_error}")
endif()

execute_process(
    COMMAND "${PSLOG_ROOT}/lua/scripts/release_version.sh"
    WORKING_DIRECTORY "${PSLOG_ROOT}"
    RESULT_VARIABLE version_result
    OUTPUT_VARIABLE version
)
if(NOT version_result EQUAL 0)
    message(FATAL_ERROR "failed to resolve Lua release version")
endif()
string(STRIP "${version}" version)
set(source_archive "${PSLOG_ROOT}/dist/lua-pslog-${version}.tar.gz")
set(source_rock "${PSLOG_ROOT}/dist/lua-pslog-${version}-1.src.rock")
set(source_root "${PSLOG_BINARY_DIR}/lua-release-artifacts-test/lua-pslog-${version}")
set(extract_root "${PSLOG_BINARY_DIR}/lua-release-artifacts-test")
file(REMOVE_RECURSE "${extract_root}")
file(MAKE_DIRECTORY "${extract_root}")

execute_process(
    COMMAND "${CMAKE_COMMAND}" -E tar xzf "${source_archive}"
    WORKING_DIRECTORY "${extract_root}"
    RESULT_VARIABLE extract_result
    ERROR_VARIABLE extract_error
)
if(NOT extract_result EQUAL 0)
    message(FATAL_ERROR "failed to extract Lua source package: ${extract_error}")
endif()

file(GLOB_RECURSE staged_files RELATIVE "${source_root}" LIST_DIRECTORIES false
    "${source_root}/*")
list(SORT staged_files)
set(expected_files
    LICENSE
    README.md
    RELEASE_MANIFEST
    VERSION
    include/pslog.h
    include/pslog_lua.h
    include/pslog_version.h
    lua/README.md
    lua/RELEASE_MANIFEST.in
    lua/lua-pslog.rockspec.in
    lua/pslog/init.lua
    lua/scripts/release_version.sh
    lua/scripts/render_release_rockspec.sh
    lua/scripts/stage_release_sources.sh
    lua/src/pslog_lua.c
)
list(SORT expected_files)
if(NOT staged_files STREQUAL expected_files)
    string(JOIN "\n" actual_listing ${staged_files})
    string(JOIN "\n" expected_listing ${expected_files})
    message(FATAL_ERROR
        "Lua source package does not match its explicit minimal manifest\n"
        "actual:\n${actual_listing}\nexpected:\n${expected_listing}")
endif()

file(READ "${source_root}/RELEASE_MANIFEST" release_manifest)
string(REPLACE "\n" ";" manifest_files "${release_manifest}")
list(REMOVE_ITEM manifest_files "")
list(SORT manifest_files)
if(NOT manifest_files STREQUAL expected_files)
    message(FATAL_ERROR "Lua source package RELEASE_MANIFEST does not match its payload")
endif()
file(READ "${source_root}/VERSION" staged_version)
string(STRIP "${staged_version}" staged_version)
if(NOT staged_version STREQUAL version)
    message(FATAL_ERROR "Lua source package VERSION does not match its artifact version")
endif()

execute_process(
    COMMAND "${CMAKE_COMMAND}" -E tar tf "${source_rock}"
    RESULT_VARIABLE rock_list_result
    OUTPUT_VARIABLE rock_listing
    ERROR_VARIABLE rock_list_error
)
if(NOT rock_list_result EQUAL 0)
    message(FATAL_ERROR "failed to inspect Lua source rock: ${rock_list_error}")
endif()
string(FIND "${rock_listing}" "lua-pslog-${version}.tar.gz" source_archive_entry)
string(FIND "${rock_listing}" "lua-pslog-${version}-1.rockspec" rockspec_entry)
if(source_archive_entry EQUAL -1 OR rockspec_entry EQUAL -1)
    message(FATAL_ERROR "Lua source rock is missing its rendered rockspec or staged source package")
endif()
set(rock_extract_root "${extract_root}/source-rock")
file(MAKE_DIRECTORY "${rock_extract_root}")
execute_process(
    COMMAND "${CMAKE_COMMAND}" -E tar xf "${source_rock}"
    WORKING_DIRECTORY "${rock_extract_root}"
    RESULT_VARIABLE rock_extract_result
    ERROR_VARIABLE rock_extract_error
)
if(NOT rock_extract_result EQUAL 0)
    message(FATAL_ERROR "failed to extract Lua source rock: ${rock_extract_error}")
endif()
file(SHA256 "${source_archive}" source_archive_sha256)
file(SHA256 "${rock_extract_root}/lua-pslog-${version}.tar.gz" embedded_archive_sha256)
if(NOT source_archive_sha256 STREQUAL embedded_archive_sha256)
    message(FATAL_ERROR "Lua source rock does not contain the staged Lua source package")
endif()

set(provenance_root "${extract_root}/provenance-fixture")
file(MAKE_DIRECTORY "${provenance_root}/lua")
file(WRITE "${provenance_root}/lua/RELEASE_MANIFEST.in"
    "lua/README.md\n"
    "tracked.txt\n")
file(WRITE "${provenance_root}/lua/README.md" "committed Lua documentation\n")
file(WRITE "${provenance_root}/tracked.txt" "committed source\n")
file(WRITE "${provenance_root}/pslog_version.h" "#define PSLOG_VERSION_STRING \"9.9.9\"\n")
execute_process(
    COMMAND git init --quiet
    WORKING_DIRECTORY "${provenance_root}"
    RESULT_VARIABLE fixture_git_init_result
)
execute_process(COMMAND git config user.email lifecycle-test@example.invalid
    WORKING_DIRECTORY "${provenance_root}" RESULT_VARIABLE fixture_git_email_result)
execute_process(COMMAND git config user.name lifecycle-test
    WORKING_DIRECTORY "${provenance_root}" RESULT_VARIABLE fixture_git_name_result)
execute_process(COMMAND git add . WORKING_DIRECTORY "${provenance_root}"
    RESULT_VARIABLE fixture_git_add_result)
execute_process(COMMAND git commit --quiet -m fixture
    WORKING_DIRECTORY "${provenance_root}" RESULT_VARIABLE fixture_git_commit_result)
if(NOT fixture_git_init_result EQUAL 0 OR NOT fixture_git_email_result EQUAL 0
        OR NOT fixture_git_name_result EQUAL 0 OR NOT fixture_git_add_result EQUAL 0
        OR NOT fixture_git_commit_result EQUAL 0)
    message(FATAL_ERROR "failed to create the Lua release provenance fixture")
endif()
file(WRITE "${provenance_root}/tracked.txt" "uncommitted source\n")
execute_process(
    COMMAND bash "${PSLOG_ROOT}/lua/scripts/stage_release_sources.sh"
        "${provenance_root}" "${provenance_root}/stage" "9.9.9"
        "${provenance_root}/pslog_version.h"
    RESULT_VARIABLE provenance_result
    ERROR_VARIABLE provenance_error
)
if(NOT provenance_result EQUAL 0)
    message(FATAL_ERROR "Lua release provenance staging failed: ${provenance_error}")
endif()
file(READ "${provenance_root}/stage/tracked.txt" staged_tracked_source)
if(NOT staged_tracked_source STREQUAL "committed source\n")
    message(FATAL_ERROR "Lua release staging used uncommitted working-tree content")
endif()
file(REMOVE_RECURSE "${extract_root}")
