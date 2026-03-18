#ifndef PSLOG_BENCH_FIXTURE_H
#define PSLOG_BENCH_FIXTURE_H

#include "pslog.h"

#include <stddef.h>

typedef enum bench_field_kind {
  BENCH_FIELD_NULL = 0,
  BENCH_FIELD_STRING = 1,
  BENCH_FIELD_BOOL = 2,
  BENCH_FIELD_INT64 = 3,
  BENCH_FIELD_UINT64 = 4,
  BENCH_FIELD_FLOAT64 = 5
} bench_field_kind;

typedef struct bench_field_spec {
  const char *key;
  size_t key_len;
  unsigned char trusted_key;
  bench_field_kind kind;
  const char *string_value;
  size_t string_len;
  unsigned char trusted_string;
  unsigned char console_simple_string;
  long signed_value;
  unsigned long unsigned_value;
  double double_value;
  int bool_value;
} bench_field_spec;

typedef struct bench_entry_spec {
  pslog_level level;
  const char *message;
  const bench_field_spec *fields;
  size_t field_count;
} bench_entry_spec;

typedef struct bench_dataset_spec {
  const bench_entry_spec *entries;
  size_t entry_count;
  const bench_field_spec *static_fields;
  size_t static_field_count;
  const bench_entry_spec *dynamic_entries;
  size_t dynamic_entry_count;
} bench_dataset_spec;

typedef struct bench_prepared_entry {
  pslog_level level;
  const char *message;
  pslog_field *fields;
  size_t field_count;
} bench_prepared_entry;

typedef struct bench_prepared_fields {
  pslog_field *fields;
  size_t field_count;
} bench_prepared_fields;

typedef struct bench_prepared_dataset {
  bench_prepared_entry *entries;
  pslog_field *fields;
  size_t entry_count;
  size_t field_count;
} bench_prepared_dataset;

typedef void (*bench_kvfmt_call_fn)(pslog_logger *log);

typedef struct bench_kvfmt_dataset {
  const bench_kvfmt_call_fn *calls;
  size_t entry_count;
} bench_kvfmt_dataset;

size_t bench_dataset_max_fields(const bench_entry_spec *entries,
                                size_t entry_count);
void bench_build_fields(const bench_field_spec *specs, size_t count,
                        pslog_field *fields);
void bench_log_level_fields(pslog_logger *log, pslog_level level,
                            const char *msg, const pslog_field *fields,
                            size_t count);
int bench_prepare_dataset(const bench_entry_spec *entries, size_t entry_count,
                          bench_prepared_dataset *prepared);
int bench_prepare_fields(const bench_field_spec *specs, size_t field_count,
                         bench_prepared_fields *prepared);
void bench_destroy_prepared_dataset(bench_prepared_dataset *prepared);
void bench_destroy_prepared_fields(bench_prepared_fields *prepared);

extern const bench_dataset_spec pslog_bench_production_dataset;
extern const bench_kvfmt_dataset pslog_bench_production_levelf_dataset;
extern const bench_kvfmt_dataset pslog_bench_production_dynamic_levelf_dataset;

#endif
