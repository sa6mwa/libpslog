if(NOT DEFINED PSLOG_BINARY_DIR OR NOT DEFINED PSLOG_ROOT)
    message(FATAL_ERROR "PSLOG_BINARY_DIR and PSLOG_ROOT are required")
endif()

set(timing_file "${PSLOG_BINARY_DIR}/run-timed-test.tsv")
file(REMOVE "${timing_file}")
execute_process(
    COMMAND ${CMAKE_COMMAND} -E env "PKT_TIMING_FILE=${timing_file}"
        "${PSLOG_ROOT}/scripts/run_timed.sh" "timing-fixture" ${CMAKE_COMMAND} -E true
    RESULT_VARIABLE timing_result
    OUTPUT_VARIABLE timing_stdout
    ERROR_VARIABLE timing_stderr
)
if(NOT timing_result EQUAL 0)
    message(FATAL_ERROR "timing wrapper failed: ${timing_stderr}")
endif()
if(NOT timing_stdout MATCHES "PKT_TIMING_BEGIN phase=timing-fixture" OR
   NOT timing_stdout MATCHES "PKT_TIMING_END phase=timing-fixture status=0 elapsed_seconds=[0-9]+")
    message(FATAL_ERROR "timing wrapper did not emit the lifecycle timing contract: ${timing_stdout}")
endif()
if(NOT EXISTS "${timing_file}")
    message(FATAL_ERROR "timing wrapper did not write the requested timing file")
endif()
file(READ "${timing_file}" timing_rows)
if(NOT timing_rows MATCHES "phase[ \t]+status[ \t]+started_epoch[ \t]+finished_epoch[ \t]+elapsed_seconds" OR
   NOT timing_rows MATCHES "timing-fixture[ \t]+0[ \t]+[0-9]+[ \t]+[0-9]+[ \t]+[0-9]+")
    message(FATAL_ERROR "timing file has the wrong lifecycle timing schema: ${timing_rows}")
endif()
file(REMOVE "${timing_file}")

execute_process(
    COMMAND ${CMAKE_COMMAND} -E env "PKT_TIMING_FILE=${timing_file}"
        "${PSLOG_ROOT}/scripts/run_timed.sh" "timing-failure" ${CMAKE_COMMAND} -E false
    RESULT_VARIABLE failure_result
    OUTPUT_VARIABLE failure_stdout
    ERROR_VARIABLE failure_stderr
)
if(failure_result EQUAL 0)
    message(FATAL_ERROR "timing wrapper concealed a failed lifecycle phase")
endif()
if(NOT failure_stdout MATCHES "PKT_TIMING_END phase=timing-failure status=[1-9][0-9]* elapsed_seconds=[0-9]+")
    message(FATAL_ERROR "timing wrapper did not report the failed phase: ${failure_stdout}${failure_stderr}")
endif()
file(READ "${timing_file}" failure_rows)
if(NOT failure_rows MATCHES "timing-failure[ \t]+[1-9][0-9]*[ \t]+[0-9]+[ \t]+[0-9]+[ \t]+[0-9]+")
    message(FATAL_ERROR "timing wrapper did not persist the failed phase status: ${failure_rows}")
endif()
file(REMOVE "${timing_file}")
