#include "bench/bench_liblogger.h"

int pslog_bench_liblogger_available(void) { return 0; }

int pslog_bench_liblogger_run(const pslog_bench_liblogger_entry_spec *entries,
                              size_t entry_count, unsigned long long iterations,
                              pslog_bench_liblogger_result *result) {
  (void)entries;
  (void)entry_count;
  (void)iterations;
  (void)result;
  return -1;
}
