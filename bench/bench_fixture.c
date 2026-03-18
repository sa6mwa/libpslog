#include "bench/bench_fixture.h"

#include <stdlib.h>
#include <string.h>

size_t bench_dataset_max_fields(const bench_entry_spec *entries,
                                size_t entry_count) {
  size_t i;
  size_t max_fields;

  max_fields = 0u;
  if (entries == NULL) {
    return 0u;
  }
  for (i = 0u; i < entry_count; ++i) {
    if (entries[i].field_count > max_fields) {
      max_fields = entries[i].field_count;
    }
  }
  return max_fields;
}

static pslog_field bench_build_field(const bench_field_spec *spec) {
  pslog_field field;

  if (spec == NULL) {
    return pslog_null("");
  }

  field.key = spec->key;
  field.key_len = spec->key_len;
  field.value_len = 0u;
  field.trusted_key = spec->trusted_key;
  field.trusted_value = 0u;
  field.console_simple_value = 0u;

  switch (spec->kind) {
  case BENCH_FIELD_NULL:
    field.type = PSLOG_FIELD_NULL;
    field.as.pointer_value = NULL;
    return field;
  case BENCH_FIELD_STRING:
    field.type = PSLOG_FIELD_STRING;
    field.value_len = spec->string_len;
    field.trusted_value = spec->trusted_string;
    field.console_simple_value = spec->console_simple_string;
    field.as.string_value =
        spec->string_value != NULL ? spec->string_value : "";
    return field;
  case BENCH_FIELD_BOOL:
    field.type = PSLOG_FIELD_BOOL;
    field.as.bool_value = spec->bool_value ? 1 : 0;
    return field;
  case BENCH_FIELD_INT64:
    field.type = PSLOG_FIELD_SIGNED;
    field.as.signed_value = spec->signed_value;
    return field;
  case BENCH_FIELD_UINT64:
    field.type = PSLOG_FIELD_UNSIGNED;
    field.as.unsigned_value = spec->unsigned_value;
    return field;
  case BENCH_FIELD_FLOAT64:
    field.type = PSLOG_FIELD_DOUBLE;
    field.as.double_value = spec->double_value;
    return field;
  default:
    field.type = PSLOG_FIELD_NULL;
    field.as.pointer_value = NULL;
    return field;
  }
}

void bench_build_fields(const bench_field_spec *specs, size_t count,
                        pslog_field *fields) {
  size_t i;

  if (specs == NULL || fields == NULL) {
    return;
  }
  for (i = 0u; i < count; ++i) {
    fields[i] = bench_build_field(&specs[i]);
  }
}

void bench_log_level_fields(pslog_logger *log, pslog_level level,
                            const char *msg, const pslog_field *fields,
                            size_t count) {
  if (log == NULL) {
    return;
  }
  switch (level) {
  case PSLOG_LEVEL_TRACE:
    log->trace(log, msg, fields, count);
    break;
  case PSLOG_LEVEL_DEBUG:
    log->debug(log, msg, fields, count);
    break;
  case PSLOG_LEVEL_INFO:
    log->info(log, msg, fields, count);
    break;
  case PSLOG_LEVEL_WARN:
    log->warn(log, msg, fields, count);
    break;
  case PSLOG_LEVEL_ERROR:
    log->error(log, msg, fields, count);
    break;
  case PSLOG_LEVEL_NOLEVEL:
    log->log(log, PSLOG_LEVEL_NOLEVEL, msg, fields, count);
    break;
  default:
    log->log(log, level, msg, fields, count);
    break;
  }
}

int bench_prepare_dataset(const bench_entry_spec *entries, size_t entry_count,
                          bench_prepared_dataset *prepared) {
  size_t i;
  size_t total_fields;
  size_t offset;

  if (prepared == NULL) {
    return -1;
  }
  memset(prepared, 0, sizeof(*prepared));
  if (entries == NULL || entry_count == 0u) {
    return 0;
  }

  total_fields = 0u;
  for (i = 0u; i < entry_count; ++i) {
    total_fields += entries[i].field_count;
  }

  prepared->entries =
      (bench_prepared_entry *)calloc(entry_count, sizeof(bench_prepared_entry));
  if (prepared->entries == NULL) {
    return -1;
  }
  if (total_fields > 0u) {
    prepared->fields = (pslog_field *)calloc(total_fields, sizeof(pslog_field));
    if (prepared->fields == NULL) {
      free(prepared->entries);
      memset(prepared, 0, sizeof(*prepared));
      return -1;
    }
  }

  prepared->entry_count = entry_count;
  prepared->field_count = total_fields;
  offset = 0u;
  for (i = 0u; i < entry_count; ++i) {
    prepared->entries[i].level = entries[i].level;
    prepared->entries[i].message = entries[i].message;
    prepared->entries[i].field_count = entries[i].field_count;
    if (entries[i].field_count > 0u) {
      prepared->entries[i].fields = prepared->fields + offset;
      bench_build_fields(entries[i].fields, entries[i].field_count,
                         prepared->entries[i].fields);
      offset += entries[i].field_count;
    }
  }
  return 0;
}

int bench_prepare_fields(const bench_field_spec *specs, size_t field_count,
                         bench_prepared_fields *prepared) {
  if (prepared == NULL) {
    return -1;
  }
  memset(prepared, 0, sizeof(*prepared));
  if (specs == NULL || field_count == 0u) {
    return 0;
  }

  prepared->fields = (pslog_field *)calloc(field_count, sizeof(pslog_field));
  if (prepared->fields == NULL) {
    return -1;
  }
  bench_build_fields(specs, field_count, prepared->fields);
  prepared->field_count = field_count;
  return 0;
}

void bench_destroy_prepared_dataset(bench_prepared_dataset *prepared) {
  if (prepared == NULL) {
    return;
  }
  free(prepared->fields);
  free(prepared->entries);
  memset(prepared, 0, sizeof(*prepared));
}

void bench_destroy_prepared_fields(bench_prepared_fields *prepared) {
  if (prepared == NULL) {
    return;
  }
  free(prepared->fields);
  memset(prepared, 0, sizeof(*prepared));
}
