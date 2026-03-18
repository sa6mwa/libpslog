#ifndef PSLOG_BENCH_LIBLOGGER_H
#define PSLOG_BENCH_LIBLOGGER_H

#include <stddef.h>

typedef enum pslog_bench_liblogger_level {
  PSLOG_BENCH_LIBLOGGER_TRACE = -1,
  PSLOG_BENCH_LIBLOGGER_DEBUG = 0,
  PSLOG_BENCH_LIBLOGGER_INFO = 1,
  PSLOG_BENCH_LIBLOGGER_WARN = 2,
  PSLOG_BENCH_LIBLOGGER_ERROR = 3,
  PSLOG_BENCH_LIBLOGGER_FATAL = 4,
  PSLOG_BENCH_LIBLOGGER_PANIC = 5,
  PSLOG_BENCH_LIBLOGGER_NOLEVEL = 6
} pslog_bench_liblogger_level;

typedef enum pslog_bench_liblogger_field_kind {
  PSLOG_BENCH_LIBLOGGER_FIELD_NULL = 0,
  PSLOG_BENCH_LIBLOGGER_FIELD_STRING = 1,
  PSLOG_BENCH_LIBLOGGER_FIELD_BOOL = 2,
  PSLOG_BENCH_LIBLOGGER_FIELD_INT64 = 3,
  PSLOG_BENCH_LIBLOGGER_FIELD_UINT64 = 4,
  PSLOG_BENCH_LIBLOGGER_FIELD_FLOAT64 = 5
} pslog_bench_liblogger_field_kind;

typedef struct pslog_bench_liblogger_field_spec {
  const char *key;
  pslog_bench_liblogger_field_kind kind;
  const char *string_value;
  long signed_value;
  unsigned long unsigned_value;
  double double_value;
  int bool_value;
} pslog_bench_liblogger_field_spec;

typedef struct pslog_bench_liblogger_entry_spec {
  pslog_bench_liblogger_level level;
  const char *message;
  const pslog_bench_liblogger_field_spec *fields;
  size_t field_count;
} pslog_bench_liblogger_entry_spec;

typedef struct pslog_bench_liblogger_result {
  unsigned long long elapsed_ns;
  unsigned long long bytes;
  unsigned long long ops;
} pslog_bench_liblogger_result;

int pslog_bench_liblogger_available(void);
int pslog_bench_liblogger_run(const pslog_bench_liblogger_entry_spec *entries,
                              size_t entry_count, unsigned long long iterations,
                              pslog_bench_liblogger_result *result);

#endif
