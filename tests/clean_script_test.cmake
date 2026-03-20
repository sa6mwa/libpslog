set(clean_test_root "${PSLOG_BINARY_DIR}/clean-script-test-root")
set(unrelated_file "${clean_test_root}/keep.txt")

file(REMOVE_RECURSE "${clean_test_root}")
file(MAKE_DIRECTORY "${clean_test_root}/build/nested")
file(MAKE_DIRECTORY "${clean_test_root}/dist")
file(WRITE "${clean_test_root}/build/nested/build.txt" "build\n")
file(WRITE "${clean_test_root}/dist/artifact.txt" "dist\n")
file(WRITE "${unrelated_file}" "keep\n")

execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/clean.sh" --root "${clean_test_root}"
    RESULT_VARIABLE clean_result
)

if(NOT clean_result EQUAL 0)
    message(FATAL_ERROR "full clean failed for ${clean_test_root}")
endif()

if(EXISTS "${clean_test_root}/build")
    message(FATAL_ERROR "full clean did not remove build/ under ${clean_test_root}")
endif()

if(EXISTS "${clean_test_root}/dist")
    message(FATAL_ERROR "full clean did not remove dist/ under ${clean_test_root}")
endif()

if(NOT EXISTS "${unrelated_file}")
    message(FATAL_ERROR "full clean removed unrelated files under ${clean_test_root}")
endif()

file(MAKE_DIRECTORY "${clean_test_root}/build")
file(MAKE_DIRECTORY "${clean_test_root}/dist")
file(WRITE "${clean_test_root}/build/recreated.txt" "build\n")
file(WRITE "${clean_test_root}/dist/recreated.txt" "dist\n")

execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/clean.sh" --dist-only --root "${clean_test_root}"
    RESULT_VARIABLE clean_dist_result
)

if(NOT clean_dist_result EQUAL 0)
    message(FATAL_ERROR "dist-only clean failed for ${clean_test_root}")
endif()

if(NOT EXISTS "${clean_test_root}/build")
    message(FATAL_ERROR "dist-only clean unexpectedly removed build/ under ${clean_test_root}")
endif()

if(EXISTS "${clean_test_root}/dist")
    message(FATAL_ERROR "dist-only clean did not remove dist/ under ${clean_test_root}")
endif()

if(NOT EXISTS "${unrelated_file}")
    message(FATAL_ERROR "dist-only clean removed unrelated files under ${clean_test_root}")
endif()

file(MAKE_DIRECTORY "${clean_test_root}/build")
file(MAKE_DIRECTORY "${clean_test_root}/dist")
file(WRITE "${clean_test_root}/build/from-cmake.txt" "build\n")
file(WRITE "${clean_test_root}/dist/from-cmake.txt" "dist\n")

execute_process(
    COMMAND ${CMAKE_COMMAND}
        -DCLEAN_ROOT=${clean_test_root}
        -P ${PSLOG_ROOT}/cmake/clean.cmake
    RESULT_VARIABLE clean_cmake_result
)

if(NOT clean_cmake_result EQUAL 0)
    message(FATAL_ERROR "cmake clean wrapper failed for ${clean_test_root}")
endif()

if(EXISTS "${clean_test_root}/build")
    message(FATAL_ERROR "cmake clean wrapper did not remove build/ under ${clean_test_root}")
endif()

if(EXISTS "${clean_test_root}/dist")
    message(FATAL_ERROR "cmake clean wrapper did not remove dist/ under ${clean_test_root}")
endif()

if(NOT EXISTS "${unrelated_file}")
    message(FATAL_ERROR "cmake clean wrapper removed unrelated files under ${clean_test_root}")
endif()

file(REMOVE_RECURSE "${clean_test_root}")
