#ifndef PSLOG_BENCH_QUILL_H
#define PSLOG_BENCH_QUILL_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef enum pslog_bench_quill_level {
  PSLOG_BENCH_QUILL_TRACE = -1,
  PSLOG_BENCH_QUILL_DEBUG = 0,
  PSLOG_BENCH_QUILL_INFO = 1,
  PSLOG_BENCH_QUILL_WARN = 2,
  PSLOG_BENCH_QUILL_ERROR = 3,
  PSLOG_BENCH_QUILL_FATAL = 4,
  PSLOG_BENCH_QUILL_PANIC = 5,
  PSLOG_BENCH_QUILL_NOLEVEL = 6
} pslog_bench_quill_level;

typedef enum pslog_bench_quill_field_kind {
  PSLOG_BENCH_QUILL_FIELD_NULL = 0,
  PSLOG_BENCH_QUILL_FIELD_STRING = 1,
  PSLOG_BENCH_QUILL_FIELD_BOOL = 2,
  PSLOG_BENCH_QUILL_FIELD_INT64 = 3,
  PSLOG_BENCH_QUILL_FIELD_UINT64 = 4,
  PSLOG_BENCH_QUILL_FIELD_FLOAT64 = 5
} pslog_bench_quill_field_kind;

typedef struct pslog_bench_quill_field_spec {
  const char *key;
  pslog_bench_quill_field_kind kind;
  const char *string_value;
  long long signed_value;
  unsigned long long unsigned_value;
  double double_value;
  int bool_value;
} pslog_bench_quill_field_spec;

typedef struct pslog_bench_quill_entry_spec {
  pslog_bench_quill_level level;
  const char *message;
  const pslog_bench_quill_field_spec *fields;
  size_t field_count;
} pslog_bench_quill_entry_spec;

typedef struct pslog_bench_quill_result {
  unsigned long long elapsed_ns;
  unsigned long long bytes;
  unsigned long long ops;
} pslog_bench_quill_result;

int pslog_bench_quill_available(void);
int pslog_bench_quill_run(const pslog_bench_quill_entry_spec *entries,
                          size_t entry_count, unsigned long long iterations,
                          pslog_bench_quill_result *result);
int pslog_bench_quill_render(const pslog_bench_quill_entry_spec *entry,
                             char *buffer, size_t capacity, size_t *written);

#ifdef __cplusplus
}
#endif

#endif
