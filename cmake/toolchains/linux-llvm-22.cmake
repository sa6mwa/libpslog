# LLVM is reserved for compiler-rt consumers: ASan, MSan, and libFuzzer.
set(CMAKE_SYSTEM_NAME Linux)
set(CMAKE_SYSTEM_PROCESSOR x86_64)
set(CMAKE_TRY_COMPILE_TARGET_TYPE STATIC_LIBRARY)

set(resolver "${CMAKE_CURRENT_LIST_DIR}/../../scripts/cpkt-llvm.sh")
execute_process(
    COMMAND "${resolver}" ensure
    RESULT_VARIABLE resolver_result
    OUTPUT_VARIABLE resolver_stdout
    ERROR_VARIABLE resolver_stderr
)
if(NOT resolver_result EQUAL 0)
    message(FATAL_ERROR "Unable to provision pinned LLVM 22.1.6.\n${resolver_stdout}${resolver_stderr}")
endif()
execute_process(
    COMMAND "${resolver}" discover
    RESULT_VARIABLE resolver_result
    OUTPUT_VARIABLE resolver_description
    ERROR_VARIABLE resolver_stderr
)
if(NOT resolver_result EQUAL 0)
    message(FATAL_ERROR "Unable to inspect pinned LLVM 22.1.6.\n${resolver_stderr}")
endif()
foreach(key cc cxx ld ar ranlib strip nm objcopy objdump addr2line readelf)
    string(REGEX MATCH "${key}=([^\r\n]+)" resolver_match "${resolver_description}")
    if(NOT resolver_match)
        message(FATAL_ERROR "LLVM resolver did not report ${key}")
    endif()
    set(llvm_${key} "${CMAKE_MATCH_1}")
endforeach()
set(CMAKE_C_COMPILER "${llvm_cc}" CACHE FILEPATH "" FORCE)
set(CMAKE_CXX_COMPILER "${llvm_cxx}" CACHE FILEPATH "" FORCE)
set(CMAKE_LINKER "${llvm_ld}" CACHE FILEPATH "" FORCE)
set(CMAKE_AR "${llvm_ar}" CACHE FILEPATH "" FORCE)
set(CMAKE_RANLIB "${llvm_ranlib}" CACHE FILEPATH "" FORCE)
set(CMAKE_STRIP "${llvm_strip}" CACHE FILEPATH "" FORCE)
set(CMAKE_NM "${llvm_nm}" CACHE FILEPATH "" FORCE)
set(CMAKE_OBJCOPY "${llvm_objcopy}" CACHE FILEPATH "" FORCE)
set(CMAKE_OBJDUMP "${llvm_objdump}" CACHE FILEPATH "" FORCE)
set(CMAKE_ADDR2LINE "${llvm_addr2line}" CACHE FILEPATH "" FORCE)
set(CMAKE_READELF "${llvm_readelf}" CACHE FILEPATH "" FORCE)
set(PSLOG_TARGET_ARCH x86_64 CACHE STRING "" FORCE)
set(PSLOG_TARGET_OS linux CACHE STRING "" FORCE)
set(PSLOG_TARGET_LIBC gnu CACHE STRING "" FORCE)
set(PSLOG_TARGET_ID x86_64-linux-gnu CACHE STRING "" FORCE)
