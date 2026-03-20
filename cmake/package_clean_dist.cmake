execute_process(
    COMMAND "${PSLOG_ROOT}/scripts/clean.sh" --dist-only --root "${PSLOG_ROOT}"
    RESULT_VARIABLE clean_result
)

if(NOT clean_result EQUAL 0)
    message(FATAL_ERROR "failed to clean dist under ${PSLOG_ROOT}")
endif()
