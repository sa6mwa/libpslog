if(NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_ROOT is required")
endif()

set(rockspec "${PSLOG_ROOT}/lua/lua-pslog.rockspec.in")
if(NOT EXISTS "${rockspec}")
    message(FATAL_ERROR "missing rockspec template: ${rockspec}")
endif()

file(READ "${rockspec}" content)

foreach(private_source
        "src/pslog.c"
        "src/pslog_emit_console.c"
        "src/pslog_emit_json.c"
        "src/pslog_palette.c")
    if(content MATCHES "\"${private_source}\"")
        message(FATAL_ERROR "Lua rockspec must not compile private libpslog source: ${private_source}")
    endif()
endforeach()

if(content MATCHES "\"src\"")
    message(FATAL_ERROR "Lua rockspec must not add the private src directory to include paths")
endif()

foreach(required_fragment
        "external_dependencies"
        "LIBPSLOG"
        "header = \"pslog.h\""
        "library = \"pslog\""
        "\"lua/src/pslog_lua.c\"")
    if(NOT content MATCHES "${required_fragment}")
        message(FATAL_ERROR "Lua rockspec missing required public libpslog linkage fragment: ${required_fragment}")
    endif()
endforeach()

if(NOT content MATCHES "libraries[ \t\r\n]*=[ \t\r\n]*\\{[^}]*\"pslog\"")
    message(FATAL_ERROR "Lua rockspec must link pslog.core against libpslog")
endif()
