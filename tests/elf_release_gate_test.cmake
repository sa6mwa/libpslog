cmake_minimum_required(VERSION 3.21)
set(work "${PSLOG_BINARY_DIR}/elf-release-gate-test")
file(REMOVE_RECURSE "${work}")
file(MAKE_DIRECTORY "${work}")
file(WRITE "${work}/main.c" "int main(void) { return 0; }\n")
file(WRITE "${work}/inspect.cmake"
    "cmake_minimum_required(VERSION 3.21)\ninclude(\"${PSLOG_ROOT}/cmake/pslog_elf_privacy.cmake\")\npslog_assert_linux_runpath(\"\${artifact}\")\n")
function(check name expected)
    execute_process(COMMAND "${PSLOG_C_COMPILER}" "${work}/main.c"
        -Wl,--build-id=none ${ARGN} -o "${work}/${name}"
        RESULT_VARIABLE result ERROR_VARIABLE error)
    if(NOT result EQUAL 0)
        message(FATAL_ERROR "Cannot build ${name}: ${error}")
    endif()
    execute_process(COMMAND "${CMAKE_COMMAND}" "-Dartifact=${work}/${name}"
        "-DPSLOG_READELF=${PSLOG_READELF}" -P "${work}/inspect.cmake"
        RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error)
    if(expected STREQUAL "pass")
        if(NOT result EQUAL 0)
            message(FATAL_ERROR "Clean ELF rejected: ${name}: ${output}${error}")
        endif()
    elseif(result EQUAL 0 OR NOT "${output}${error}" MATCHES "${expected}")
        message(FATAL_ERROR "Contaminated ELF not rejected correctly: ${name}: ${output}${error}")
    endif()
endfunction()
check(system pass)
check(static pass -static)
check(origin pass "-Wl,-rpath,$ORIGIN/../lib")
check(custom-interpreter "non-system interpreter" "-Wl,--dynamic-linker,/opt/custom/bootlin/sysroot/lib/ld-linux-x86-64.so.2")
check(mixed-rpath "non-relocatable RPATH" "-Wl,-rpath,$ORIGIN:/opt/custom/lib")
check(empty-rpath "non-relocatable RPATH" "-Wl,-rpath,$ORIGIN:")
check(blank-rpath "non-relocatable RPATH" "-Wl,-rpath,")
check(path-soname "dependency path" -shared -Wl,-soname,/opt/custom/dependency.so)
check(path-needed "dependency path" -Wl,--no-as-needed "${work}/path-soname")
file(COPY_FILE "${PSLOG_BOOTLIN_INTERPRETER}" "${work}/renamed-loader")
execute_process(COMMAND "${CMAKE_COMMAND}" "-Dartifact=${work}/renamed-loader"
    "-DPSLOG_READELF=${PSLOG_READELF}" -P "${work}/inspect.cmake"
    RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error)
if(result EQUAL 0 OR NOT "${output}${error}" MATCHES "bundles a libc loader")
    message(FATAL_ERROR "Release guard accepted a renamed Bootlin loader: ${output}${error}")
endif()
check(fake-loader "bundles a libc loader" -shared -Wl,-soname,ld-linux-x86-64.so.2)
if(NOT PSLOG_C_COMPILER MATCHES "musl")
execute_process(COMMAND "${PSLOG_C_COMPILER}" "${PSLOG_ROOT}/tests/runtime_probe.c"
    "-DPSLOG_RUNTIME_ROOT=\"/not-the-host-runtime\"" -o "${work}/host-runtime"
    COMMAND_ERROR_IS_FATAL ANY)
execute_process(COMMAND "${work}/host-runtime" RESULT_VARIABLE host_result ERROR_VARIABLE host_error)
if(NOT host_result EQUAL 2 OR NOT host_error MATCHES "Unexpected runtime mapping")
    message(FATAL_ERROR "Runtime probe failed to reject actual host libc resolution: ${host_result}: ${host_error}")
endif()
endif()
check(missing-runtime "non-system interpreter" "-Wl,--dynamic-linker,${work}/missing-loader")
execute_process(COMMAND "${work}/missing-runtime" RESULT_VARIABLE missing_result)
if(missing_result EQUAL 0)
    message(FATAL_ERROR "Executable unexpectedly ran with a missing interpreter")
endif()
# Exercise the actual release entry point, including nested source rocks.
foreach(payload custom-interpreter fake-loader mixed-rpath)
    execute_process(COMMAND "${CMAKE_COMMAND}" -E tar czf "${work}/inner.tar.gz" "${payload}"
        WORKING_DIRECTORY "${work}" COMMAND_ERROR_IS_FATAL ANY)
    execute_process(COMMAND "${CMAKE_COMMAND}" -E tar cf "${work}/outer.src.rock" inner.tar.gz
        WORKING_DIRECTORY "${work}" COMMAND_ERROR_IS_FATAL ANY)
    foreach(artifact "${work}/${payload}" "${work}/outer.src.rock")
        execute_process(COMMAND "${CMAKE_COMMAND}"
            "-DPSLOG_ROOT=${PSLOG_ROOT}" "-DPSLOG_BINARY_DIR=${work}" -DPSLOG_VERSION=0.0.0
            "-DPSLOG_PRIVACY_PATHS=${artifact}" "-DPSLOG_READELF=${PSLOG_READELF}"
            -P "${PSLOG_ROOT}/cmake/check_release_privacy.cmake"
            RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error)
        if(result EQUAL 0 OR NOT "${output}${error}" MATCHES "${payload}")
            message(FATAL_ERROR "Release entry point accepted ${payload}: ${output}${error}")
        endif()
    endforeach()
endforeach()
set(fake_sysroot "${work}/fake-sysroot")
file(MAKE_DIRECTORY "${fake_sysroot}/lib")
file(WRITE "${work}/runtime-policy.cmake"
    "set(CMAKE_SYSTEM_NAME Linux)\nset(CMAKE_HOST_SYSTEM_NAME Linux)\n"
    "set(CMAKE_SYSTEM_PROCESSOR x86_64)\nset(CMAKE_HOST_SYSTEM_PROCESSOR x86_64)\n"
    "set(CMAKE_C_COMPILER_TARGET x86_64-linux-gnu)\n"
    "set(PSLOG_BOOTLIN_SYSROOT \"${fake_sysroot}\")\n"
    "set(CMAKE_READELF \"${PSLOG_READELF}\")\n"
    "include(\"${PSLOG_ROOT}/cmake/pslog_local_runtime.cmake\")\n")
foreach(expected "missing its interpreter" "architecture mismatch")
    execute_process(COMMAND "${CMAKE_COMMAND}" -P "${work}/runtime-policy.cmake"
        RESULT_VARIABLE result OUTPUT_VARIABLE output ERROR_VARIABLE error)
    if(result EQUAL 0 OR NOT "${output}${error}" MATCHES "${expected}")
        message(FATAL_ERROR "Runtime policy accepted ${expected}: ${output}${error}")
    endif()
    file(WRITE "${fake_sysroot}/lib/ld-linux-x86-64.so.2" "invalid interpreter\n")
endforeach()
set(PSLOG_READELF "${work}/missing-readelf")
check(missing-inspector "Cannot inspect ELF")
set(PSLOG_READELF /usr/bin/true)
check(empty-inspection "Cannot inspect ELF")
