#include "bench/bench_quill.h"

int pslog_bench_quill_available(void) { return 0; }

int pslog_bench_quill_run(const pslog_bench_quill_entry_spec *entries,
                          size_t entry_count, unsigned long long iterations,
                          pslog_bench_quill_result *result) {
  (void)entries;
  (void)entry_count;
  (void)iterations;
  (void)result;
  return -1;
}

int pslog_bench_quill_render(const pslog_bench_quill_entry_spec *entry,
                             char *buffer, size_t capacity, size_t *written) {
  (void)entry;
  (void)buffer;
  (void)capacity;
  (void)written;
  return -1;
}
