cmake_minimum_required(VERSION 3.21)

include(cmake/pslog_version.cmake)

project(libpslog VERSION ${PSLOG_PROJECT_VERSION_CORE} LANGUAGES C)

set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

set(PSLOG_GENERATED_INCLUDE_DIR "${CMAKE_CURRENT_BINARY_DIR}/generated/include")
file(MAKE_DIRECTORY "${PSLOG_GENERATED_INCLUDE_DIR}")
configure_file(
    "${CMAKE_CURRENT_SOURCE_DIR}/include/pslog_version.h.in"
    "${PSLOG_GENERATED_INCLUDE_DIR}/pslog_version.h"
    @ONLY
)
set(PSLOG_GENERATED_VERSION_HEADER "${PSLOG_GENERATED_INCLUDE_DIR}/pslog_version.h")

function(pslog_configure_c_target target)
    if(UNIX OR APPLE OR CMAKE_SYSTEM_NAME STREQUAL "FreeBSD")
        target_compile_definitions(${target} PRIVATE _POSIX_C_SOURCE=200809L)
    endif()
    if(CMAKE_C_COMPILER_ID MATCHES "Clang|GNU")
        target_compile_options(${target} PRIVATE
            -std=c89
            -Wall
            -Wextra
            -Wpedantic
            -pedantic-errors
        )
    endif()
endfunction()

function(pslog_enable_coverage_if_requested target)
    get_target_property(target_type ${target} TYPE)
    if(PSLOG_ENABLE_COVERAGE AND CMAKE_C_COMPILER_ID MATCHES "Clang|GNU")
        target_compile_definitions(${target} PRIVATE PSLOG_COVERAGE_BUILD=1)
        target_compile_options(${target} PRIVATE --coverage -O0)
        if(NOT target_type STREQUAL "STATIC_LIBRARY")
            target_link_options(${target} PRIVATE --coverage)
        endif()
    endif()
endfunction()

function(pslog_assert_fuzz_toolchain_supported)
    set(fuzz_probe_dir "${CMAKE_CURRENT_BINARY_DIR}/CMakeFiles/fuzz-toolchain-probe")
    set(fuzz_probe_source "${fuzz_probe_dir}/probe.c")
    set(fuzz_probe_binary "${fuzz_probe_dir}/probe")

    file(MAKE_DIRECTORY "${fuzz_probe_dir}")
    file(WRITE "${fuzz_probe_source}" [=[
int LLVMFuzzerTestOneInput(const unsigned char *data, unsigned long size) {
    (void)data;
    (void)size;
    return 0;
}
]=])

    execute_process(
        COMMAND "${CMAKE_C_COMPILER}"
            -v
            -fsanitize=fuzzer,address
            -fno-omit-frame-pointer
            "${fuzz_probe_source}"
            -o "${fuzz_probe_binary}"
        RESULT_VARIABLE fuzz_probe_result
        OUTPUT_VARIABLE fuzz_probe_stdout
        ERROR_VARIABLE fuzz_probe_stderr
    )

    if(NOT fuzz_probe_result EQUAL 0)
        message(FATAL_ERROR
            "PSLOG_BUILD_FUZZ=ON requires a C toolchain that can link -fsanitize=fuzzer,address. "
            "The selected compiler is '${CMAKE_C_COMPILER}'. Install the matching compiler-rt/libFuzzer runtime "
            "for this Clang toolchain or use a Clang toolchain that provides those libraries.\n"
            "clang output:\n${fuzz_probe_stdout}\n${fuzz_probe_stderr}")
    endif()
endfunction()

find_program(CLANG_FORMAT_BIN NAMES clang-format)

file(GLOB_RECURSE PSLOG_FORMAT_SOURCES CONFIGURE_DEPENDS
    ${CMAKE_CURRENT_SOURCE_DIR}/include/*.h
    ${CMAKE_CURRENT_SOURCE_DIR}/src/*.c
    ${CMAKE_CURRENT_SOURCE_DIR}/src/*.h
    ${CMAKE_CURRENT_SOURCE_DIR}/bench/*.c
    ${CMAKE_CURRENT_SOURCE_DIR}/bench/*.h
    ${CMAKE_CURRENT_SOURCE_DIR}/fuzz/*.c
    ${CMAKE_CURRENT_SOURCE_DIR}/fuzz/*.h
    ${CMAKE_CURRENT_SOURCE_DIR}/tests/*.c
    ${CMAKE_CURRENT_SOURCE_DIR}/tests/*.h
    ${CMAKE_CURRENT_SOURCE_DIR}/examples/*.c
    ${CMAKE_CURRENT_SOURCE_DIR}/examples/*.h
)

if(CLANG_FORMAT_BIN)
    add_custom_target(clang-format
        COMMAND ${CLANG_FORMAT_BIN} -i ${PSLOG_FORMAT_SOURCES}
        VERBATIM
    )
else()
    add_custom_target(clang-format
        COMMAND ${CMAKE_COMMAND} -E echo "clang-format not found"
        COMMAND ${CMAKE_COMMAND} -E false
        VERBATIM
    )
endif()

add_custom_target(format DEPENDS clang-format)

option(PSLOG_BUILD_TESTS "Build unit tests" ON)
option(PSLOG_BUILD_BENCHMARKS "Build benchmarks" ON)
option(PSLOG_BUILD_FUZZ "Build fuzz targets" OFF)
option(PSLOG_BENCHMARK_WITH_LIBLOGGER "Enable optional benchmark-only liblogger comparisons" OFF)
option(PSLOG_BENCHMARK_WITH_QUILL "Enable optional benchmark-only Quill comparisons" OFF)
option(PSLOG_ENABLE_COVERAGE "Build test targets with gcov coverage instrumentation" OFF)

if(PSLOG_BUILD_TESTS OR PSLOG_BENCHMARK_WITH_LIBLOGGER OR PSLOG_BENCHMARK_WITH_QUILL)
    include(FetchContent)
endif()

if(PSLOG_BUILD_TESTS OR PSLOG_BENCHMARK_WITH_LIBLOGGER)
    set(JANSSON_BUILD_SHARED_LIBS OFF CACHE BOOL "" FORCE)
    set(JANSSON_BUILD_DOCS OFF CACHE BOOL "" FORCE)
    set(JANSSON_EXAMPLES OFF CACHE BOOL "" FORCE)
    set(JANSSON_WITHOUT_TESTS ON CACHE BOOL "" FORCE)
    set(JANSSON_INSTALL OFF CACHE BOOL "" FORCE)
    FetchContent_Declare(jansson
        GIT_REPOSITORY https://github.com/akheron/jansson.git
        GIT_TAG v2.15.0
    )
    FetchContent_MakeAvailable(jansson)
    set_target_properties(jansson PROPERTIES
        C_STANDARD 99
        C_STANDARD_REQUIRED ON
        C_EXTENSIONS ON
    )
endif()

set(PSLOG_TARGET_ARCH "${CMAKE_SYSTEM_PROCESSOR}" CACHE STRING "Artifact architecture name")
set(PSLOG_TARGET_OS "${CMAKE_SYSTEM_NAME}" CACHE STRING "Artifact operating system name")
set(PSLOG_TARGET_LIBC "" CACHE STRING "Artifact libc suffix for Linux targets (gnu, musl, or empty)")

string(TOLOWER "${PSLOG_TARGET_OS}" PSLOG_TARGET_OS_LOWER)
string(TOLOWER "${PSLOG_TARGET_ARCH}" PSLOG_TARGET_ARCH_LOWER)

if(PSLOG_TARGET_ARCH_LOWER STREQUAL "amd64")
    set(PSLOG_TARGET_ARCH_LOWER "x86_64")
endif()
if(PSLOG_TARGET_ARCH_LOWER STREQUAL "arm64")
    set(PSLOG_TARGET_ARCH_LOWER "aarch64")
endif()

if(PSLOG_TARGET_OS_LOWER STREQUAL "linux" AND NOT PSLOG_TARGET_LIBC)
    set(PSLOG_TARGET_LIBC "gnu")
endif()

if(PSLOG_TARGET_OS_LOWER STREQUAL "darwin")
    set(PSLOG_TARGET_ID "${PSLOG_TARGET_ARCH_LOWER}-darwin")
elseif(PSLOG_TARGET_OS_LOWER STREQUAL "freebsd")
    set(PSLOG_TARGET_ID "${PSLOG_TARGET_ARCH_LOWER}-freebsd")
elseif(PSLOG_TARGET_OS_LOWER STREQUAL "linux")
    if(PSLOG_TARGET_LIBC)
        set(PSLOG_TARGET_ID "${PSLOG_TARGET_ARCH_LOWER}-linux-${PSLOG_TARGET_LIBC}")
    else()
        set(PSLOG_TARGET_ID "${PSLOG_TARGET_ARCH_LOWER}-linux")
    endif()
else()
    set(PSLOG_TARGET_ID "${PSLOG_TARGET_ARCH_LOWER}-${PSLOG_TARGET_OS_LOWER}")
endif()

set(PSLOG_SOURCES
    src/pslog.c
    src/pslog_palette.c
    src/pslog_emit_console.c
    src/pslog_emit_json.c
)

set(PSLOG_BENCH_SUPPORT_SOURCES
    bench/bench_fixture.c
    bench/bench_production_data_generated.c
)

add_library(pslog_shared SHARED ${PSLOG_SOURCES})
add_library(pslog_static STATIC ${PSLOG_SOURCES})

foreach(target pslog_shared pslog_static)
    target_include_directories(${target}
        PUBLIC
            $<BUILD_INTERFACE:${PSLOG_GENERATED_INCLUDE_DIR}>
            $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
            $<INSTALL_INTERFACE:include>
        PRIVATE
            ${CMAKE_CURRENT_SOURCE_DIR}/src
    )
    target_compile_definitions(${target} PRIVATE PSLOG_VERSION_STRING="${PSLOG_RESOLVED_VERSION}")
    set_target_properties(${target} PROPERTIES
        POSITION_INDEPENDENT_CODE ON
        VISIBILITY_INLINES_HIDDEN NO
    )
    target_link_libraries(${target} PRIVATE Threads::Threads)
    pslog_configure_c_target(${target})
    pslog_enable_coverage_if_requested(${target})
endforeach()

if(PSLOG_BUILD_TESTS)
    add_library(pslog_test_static STATIC ${PSLOG_SOURCES})
    target_include_directories(pslog_test_static
        PUBLIC
            $<BUILD_INTERFACE:${PSLOG_GENERATED_INCLUDE_DIR}>
            $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        PRIVATE
            ${CMAKE_CURRENT_SOURCE_DIR}/src
    )
    target_compile_definitions(pslog_test_static PRIVATE
        PSLOG_VERSION_STRING="${PSLOG_RESOLVED_VERSION}"
        PSLOG_TEST_HOOKS=1
    )
    set_target_properties(pslog_test_static PROPERTIES
        POSITION_INDEPENDENT_CODE ON
        VISIBILITY_INLINES_HIDDEN NO
        OUTPUT_NAME pslog_test
    )
    target_link_libraries(pslog_test_static PRIVATE Threads::Threads)
    pslog_configure_c_target(pslog_test_static)
    pslog_enable_coverage_if_requested(pslog_test_static)
endif()

set_target_properties(pslog_shared PROPERTIES
    OUTPUT_NAME pslog
    VERSION ${PROJECT_VERSION}
    SOVERSION ${PROJECT_VERSION_MAJOR}
)

set_target_properties(pslog_static PROPERTIES
    OUTPUT_NAME pslog
)

if(PSLOG_BUILD_TESTS)
    enable_testing()
    add_executable(pslog_tests
        tests/test_main.c
        ${PSLOG_BENCH_SUPPORT_SOURCES}
    )
    target_include_directories(pslog_tests PRIVATE
        ${PSLOG_GENERATED_INCLUDE_DIR}
        include
        ${CMAKE_CURRENT_SOURCE_DIR}
        ${jansson_BINARY_DIR}/include
    )
    target_compile_definitions(pslog_tests PRIVATE PSLOG_TEST_HOOKS=1)
    if(UNIX OR APPLE OR CMAKE_SYSTEM_NAME STREQUAL "FreeBSD")
        target_compile_definitions(pslog_tests PRIVATE _POSIX_C_SOURCE=200809L)
    endif()
    target_link_libraries(pslog_tests PRIVATE pslog_test_static jansson)
    pslog_enable_coverage_if_requested(pslog_tests)
    add_test(NAME pslog_tests COMMAND pslog_tests)
    add_test(NAME package_archives_test
        COMMAND ${CMAKE_COMMAND}
            -DPSLOG_BINARY_DIR=${CMAKE_CURRENT_BINARY_DIR}
            -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
            -DPSLOG_VERSION=${PSLOG_RESOLVED_VERSION}
            -DPSLOG_TARGET_ID=${PSLOG_TARGET_ID}
            -DPSLOG_SHARED_LIB_NAME=$<TARGET_FILE_NAME:pslog_shared>
            -DPSLOG_SHARED_SONAME=$<TARGET_SONAME_FILE_NAME:pslog_shared>
            -P ${CMAKE_CURRENT_SOURCE_DIR}/tests/package_archives_test.cmake
    )
    set_tests_properties(package_archives_test PROPERTIES RUN_SERIAL TRUE)
    add_test(NAME version_resolution_test
        COMMAND ${CMAKE_COMMAND}
            -DPSLOG_BINARY_DIR=${CMAKE_CURRENT_BINARY_DIR}
            -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
            -P ${CMAKE_CURRENT_SOURCE_DIR}/tests/version_resolution_test.cmake
    )
    set_tests_properties(version_resolution_test PROPERTIES RUN_SERIAL TRUE)
    add_test(NAME c_only_configure_test
        COMMAND ${CMAKE_COMMAND}
            -DPSLOG_BINARY_DIR=${CMAKE_CURRENT_BINARY_DIR}
            -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
            -DPSLOG_C_COMPILER=${CMAKE_C_COMPILER}
            -P ${CMAKE_CURRENT_SOURCE_DIR}/tests/c_only_configure_test.cmake
    )
    set_tests_properties(c_only_configure_test PROPERTIES RUN_SERIAL TRUE)
    add_test(NAME clean_script_test
        COMMAND ${CMAKE_COMMAND}
            -DPSLOG_BINARY_DIR=${CMAKE_CURRENT_BINARY_DIR}
            -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
            -P ${CMAKE_CURRENT_SOURCE_DIR}/tests/clean_script_test.cmake
    )
    set_tests_properties(clean_script_test PROPERTIES RUN_SERIAL TRUE)
endif()

if(PSLOG_BUILD_BENCHMARKS)
    if(PSLOG_BENCHMARK_WITH_LIBLOGGER)
        add_library(pslog_bench_liblogger SHARED)
    else()
        add_library(pslog_bench_liblogger SHARED bench/bench_liblogger_stub.c)
    endif()
    target_include_directories(pslog_bench_liblogger PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
    set_target_properties(pslog_bench_liblogger PROPERTIES
        OUTPUT_NAME pslog_bench_liblogger
        POSITION_INDEPENDENT_CODE ON
    )

    if(PSLOG_BENCHMARK_WITH_LIBLOGGER)
        FetchContent_Declare(liblogger
            GIT_REPOSITORY https://github.com/briandowns/liblogger.git
            GIT_TAG 3ecb5c636b5a97fc35d0a6373a3238f1596be82a
        )
        FetchContent_GetProperties(liblogger)
        if(NOT liblogger_POPULATED)
            FetchContent_Populate(liblogger)
        endif()

        target_sources(pslog_bench_liblogger PRIVATE
            bench/bench_liblogger.c
            ${liblogger_SOURCE_DIR}/logger.c
        )
        target_include_directories(pslog_bench_liblogger PRIVATE
            ${liblogger_SOURCE_DIR}
            ${jansson_BINARY_DIR}/include
        )
        target_link_libraries(pslog_bench_liblogger PRIVATE jansson)
        set_source_files_properties(
            bench/bench_liblogger.c
            ${liblogger_SOURCE_DIR}/logger.c
            PROPERTIES COMPILE_OPTIONS "-std=gnu99"
        )
        target_compile_definitions(pslog_bench_liblogger PRIVATE PSLOG_HAVE_LIBLOGGER_BENCH=1)
    endif()

    if(PSLOG_BENCHMARK_WITH_QUILL)
        include(CheckLanguage)
        check_language(CXX)
        if(NOT CMAKE_CXX_COMPILER)
            message(FATAL_ERROR
                "PSLOG_BENCHMARK_WITH_QUILL=ON requires a working C++ compiler for the benchmark-only Quill helper")
        endif()
        enable_language(CXX)

        add_library(pslog_bench_quill SHARED bench/bench_quill.cpp)
        FetchContent_Declare(quill
            GIT_REPOSITORY https://github.com/odygrd/quill.git
            GIT_TAG 7c0ffa54e51c2b8db6ec091c2922aeaf9b3b08c0
        )
        FetchContent_GetProperties(quill)
        if(NOT quill_POPULATED)
            FetchContent_Populate(quill)
        endif()

        target_include_directories(pslog_bench_quill PRIVATE
            ${quill_SOURCE_DIR}/include
        )
        target_link_libraries(pslog_bench_quill PRIVATE Threads::Threads)
        target_compile_features(pslog_bench_quill PRIVATE cxx_std_17)
    else()
        add_library(pslog_bench_quill SHARED bench/bench_quill_stub.c)
    endif()
    target_include_directories(pslog_bench_quill PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
    set_target_properties(pslog_bench_quill PROPERTIES
        OUTPUT_NAME pslog_bench_quill
        POSITION_INDEPENDENT_CODE ON
    )

    add_executable(pslog_bench
        bench/bench_main.c
        ${PSLOG_BENCH_SUPPORT_SOURCES}
    )
    target_include_directories(pslog_bench PRIVATE include ${CMAKE_CURRENT_SOURCE_DIR})
    target_link_libraries(pslog_bench PRIVATE pslog_static)
    pslog_enable_coverage_if_requested(pslog_bench)
    if(PSLOG_BENCHMARK_WITH_LIBLOGGER)
        target_link_libraries(pslog_bench PRIVATE pslog_bench_liblogger)
        target_compile_definitions(pslog_bench PRIVATE PSLOG_HAVE_LIBLOGGER_BENCH=1)
    endif()
endif()

if(PSLOG_BUILD_FUZZ)
    if(NOT CMAKE_C_COMPILER_ID MATCHES "Clang")
        message(FATAL_ERROR "PSLOG_BUILD_FUZZ=ON requires Clang as the C compiler")
    endif()
    pslog_assert_fuzz_toolchain_supported()

    add_executable(pslog_fuzz fuzz/pslog_fuzz.c)
    target_include_directories(pslog_fuzz PRIVATE include)
    target_link_libraries(pslog_fuzz PRIVATE pslog_static)
    pslog_configure_c_target(pslog_fuzz)
    target_compile_options(pslog_fuzz PRIVATE -fsanitize=fuzzer,address)
    target_link_options(pslog_fuzz PRIVATE -fsanitize=fuzzer,address)
endif()

find_program(GCOV_BIN NAMES gcov)

if(PSLOG_BUILD_TESTS AND PSLOG_ENABLE_COVERAGE AND GCOV_BIN)
    add_custom_target(coverage-report
        COMMAND ${CMAKE_CTEST_COMMAND} --test-dir ${CMAKE_CURRENT_BINARY_DIR} --output-on-failure
        COMMAND ${CMAKE_COMMAND} -E env GCOV=${GCOV_BIN}
                ${CMAKE_CURRENT_SOURCE_DIR}/scripts/coverage_report.sh
                ${CMAKE_CURRENT_BINARY_DIR}
                ${CMAKE_CURRENT_BINARY_DIR}/coverage-report
        DEPENDS pslog_tests
        USES_TERMINAL
        VERBATIM
    )
endif()

install(FILES include/pslog.h "${PSLOG_GENERATED_VERSION_HEADER}" DESTINATION include)
install(TARGETS pslog_shared
    LIBRARY DESTINATION lib
    RUNTIME DESTINATION bin
    ARCHIVE DESTINATION lib
)
install(TARGETS pslog_static
    ARCHIVE DESTINATION lib
)

add_custom_target(package-runtime
    COMMAND ${CMAKE_COMMAND}
        -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
        -DPSLOG_BINARY_DIR=${CMAKE_CURRENT_BINARY_DIR}
        -DPSLOG_VERSION=${PSLOG_RESOLVED_VERSION}
        -DPSLOG_TARGET_ID=${PSLOG_TARGET_ID}
        -DPSLOG_PUBLIC_HEADER=${CMAKE_CURRENT_SOURCE_DIR}/include/pslog.h
        -DPSLOG_PUBLIC_VERSION_HEADER=${PSLOG_GENERATED_VERSION_HEADER}
        -DPSLOG_SHARED_LIB=$<TARGET_FILE:pslog_shared>
        -DPSLOG_SHARED_LIB_NAME=$<TARGET_FILE_NAME:pslog_shared>
        -DPSLOG_SHARED_SONAME=$<TARGET_SONAME_FILE_NAME:pslog_shared>
        -P ${CMAKE_CURRENT_SOURCE_DIR}/cmake/package_runtime.cmake
    DEPENDS pslog_shared
    VERBATIM
)

add_custom_target(package-dev
    COMMAND ${CMAKE_COMMAND}
        -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
        -DPSLOG_BINARY_DIR=${CMAKE_CURRENT_BINARY_DIR}
        -DPSLOG_VERSION=${PSLOG_RESOLVED_VERSION}
        -DPSLOG_TARGET_ID=${PSLOG_TARGET_ID}
        -DPSLOG_PUBLIC_HEADER=${CMAKE_CURRENT_SOURCE_DIR}/include/pslog.h
        -DPSLOG_PUBLIC_VERSION_HEADER=${PSLOG_GENERATED_VERSION_HEADER}
        -DPSLOG_STATIC_LIB=$<TARGET_FILE:pslog_static>
        -P ${CMAKE_CURRENT_SOURCE_DIR}/cmake/package_dev.cmake
    DEPENDS pslog_static
    VERBATIM
)

add_custom_target(package-checksums
    COMMAND ${CMAKE_COMMAND}
        -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
        -DPSLOG_VERSION=${PSLOG_RESOLVED_VERSION}
        -P ${CMAKE_CURRENT_SOURCE_DIR}/cmake/package_checksums.cmake
    VERBATIM
)

add_custom_target(package-clean-dist
    COMMAND ${CMAKE_COMMAND}
        -DPSLOG_ROOT=${CMAKE_CURRENT_SOURCE_DIR}
        -P ${CMAKE_CURRENT_SOURCE_DIR}/cmake/package_clean_dist.cmake
    VERBATIM
)
