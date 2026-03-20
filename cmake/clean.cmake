get_filename_component(pslog_repo_root "${CMAKE_CURRENT_LIST_DIR}/.." ABSOLUTE)

if(NOT DEFINED CLEAN_ROOT)
    set(pslog_clean_root "${pslog_repo_root}")
else()
    get_filename_component(pslog_clean_root "${CLEAN_ROOT}" ABSOLUTE)
endif()

execute_process(
    COMMAND "${pslog_repo_root}/scripts/clean.sh" --root "${pslog_clean_root}"
    RESULT_VARIABLE clean_result
)

if(NOT clean_result EQUAL 0)
    message(FATAL_ERROR "failed to clean ${pslog_clean_root}")
endif()
