# Resolve a complete, pinned Bootlin collection before project() discovers any
# compiler state. Linux projects must not fall back to host compilers or tools.

function(pslog_configure_bootlin_toolchain target_id)
    set(resolver "${CMAKE_CURRENT_LIST_DIR}/../../scripts/cpkt-toolchains.sh")
    if(NOT EXISTS "${resolver}")
        message(FATAL_ERROR "missing lifecycle toolchain resolver: ${resolver}")
    endif()

    execute_process(
        COMMAND "${resolver}" ensure "${target_id}"
        RESULT_VARIABLE resolver_result
        OUTPUT_VARIABLE resolver_stdout
        ERROR_VARIABLE resolver_stderr
    )
    if(NOT resolver_result EQUAL 0)
        message(FATAL_ERROR
            "Unable to provision the pinned Bootlin toolchain for ${target_id}.\n"
            "${resolver_stdout}${resolver_stderr}")
    endif()

    execute_process(
        COMMAND "${resolver}" discover "${target_id}"
        RESULT_VARIABLE resolver_result
        OUTPUT_VARIABLE resolver_description
        ERROR_VARIABLE resolver_stderr
    )
    if(NOT resolver_result EQUAL 0)
        message(FATAL_ERROR
            "Unable to inspect the pinned Bootlin toolchain for ${target_id}.\n"
            "${resolver_stderr}")
    endif()

    foreach(key cc cxx ld ar ranlib strip nm objcopy objdump addr2line readelf sysroot target_triple root)
        string(REGEX MATCH "${key}=([^\r\n]+)" resolver_match "${resolver_description}")
        if(NOT resolver_match)
            message(FATAL_ERROR "Bootlin resolver did not report ${key} for ${target_id}")
        endif()
        set(bootlin_${key} "${CMAKE_MATCH_1}")
    endforeach()

    set(CMAKE_C_COMPILER "${bootlin_cc}" CACHE FILEPATH "" FORCE)
    set(CMAKE_CXX_COMPILER "${bootlin_cxx}" CACHE FILEPATH "" FORCE)
    set(CMAKE_C_COMPILER_TARGET "${bootlin_target_triple}" CACHE STRING "" FORCE)
    set(CMAKE_CXX_COMPILER_TARGET "${bootlin_target_triple}" CACHE STRING "" FORCE)
    set(CMAKE_LINKER "${bootlin_ld}" CACHE FILEPATH "" FORCE)
    set(CMAKE_AR "${bootlin_ar}" CACHE FILEPATH "" FORCE)
    set(CMAKE_RANLIB "${bootlin_ranlib}" CACHE FILEPATH "" FORCE)
    set(CMAKE_STRIP "${bootlin_strip}" CACHE FILEPATH "" FORCE)
    set(CMAKE_NM "${bootlin_nm}" CACHE FILEPATH "" FORCE)
    set(CMAKE_OBJCOPY "${bootlin_objcopy}" CACHE FILEPATH "" FORCE)
    set(CMAKE_OBJDUMP "${bootlin_objdump}" CACHE FILEPATH "" FORCE)
    set(CMAKE_ADDR2LINE "${bootlin_addr2line}" CACHE FILEPATH "" FORCE)
    set(CMAKE_READELF "${bootlin_readelf}" CACHE FILEPATH "" FORCE)
    set(CMAKE_SYSROOT "${bootlin_sysroot}" CACHE PATH "" FORCE)
    set(CMAKE_FIND_ROOT_PATH "${bootlin_sysroot}" "${bootlin_root}" CACHE STRING "" FORCE)
    set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER CACHE STRING "" FORCE)
    set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY CACHE STRING "" FORCE)
    set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY CACHE STRING "" FORCE)
    set(CMAKE_FIND_ROOT_PATH_MODE_PACKAGE ONLY CACHE STRING "" FORCE)
endfunction()
