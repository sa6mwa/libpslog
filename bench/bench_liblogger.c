#define _GNU_SOURCE

#include "bench/bench_liblogger.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "logger.h"

typedef struct pslog_bench_liblogger_sink {
  FILE *stream;
  char *buffer;
  size_t size;
} pslog_bench_liblogger_sink;

static unsigned long long pslog_bench_liblogger_now_ns(void) {
#if defined(CLOCK_MONOTONIC)
  struct timespec ts;

  if (clock_gettime(CLOCK_MONOTONIC, &ts) == 0) {
    return ((unsigned long long)ts.tv_sec * 1000000000ull) +
           (unsigned long long)ts.tv_nsec;
  }
#endif
  return ((unsigned long long)clock() * 1000000000ull) /
         (unsigned long long)CLOCKS_PER_SEC;
}

static int pslog_bench_liblogger_sink_open(pslog_bench_liblogger_sink *sink) {
  if (sink == NULL) {
    return -1;
  }
  memset(sink, 0, sizeof(*sink));
#if defined(__GLIBC__) || defined(__linux__) || defined(__FreeBSD__) ||        \
    defined(__APPLE__)
  sink->stream = open_memstream(&sink->buffer, &sink->size);
#endif
  if (sink->stream == NULL) {
    sink->stream = tmpfile();
  }
  return sink->stream != NULL ? 0 : -1;
}

static void pslog_bench_liblogger_sink_close(pslog_bench_liblogger_sink *sink) {
  if (sink == NULL) {
    return;
  }
  if (sink->stream != NULL) {
    fclose(sink->stream);
  }
  free(sink->buffer);
  memset(sink, 0, sizeof(*sink));
}

static unsigned long long
pslog_bench_liblogger_sink_bytes(pslog_bench_liblogger_sink *sink) {
  long pos;

  if (sink == NULL || sink->stream == NULL) {
    return 0ull;
  }
  fflush(sink->stream);
  if (sink->buffer != NULL) {
    return (unsigned long long)sink->size;
  }
  pos = ftell(sink->stream);
  if (pos < 0) {
    return 0ull;
  }
  return (unsigned long long)pos;
}

static char *
pslog_bench_liblogger_level_name(pslog_bench_liblogger_level level) {
  switch (level) {
  case PSLOG_BENCH_LIBLOGGER_TRACE:
    return S_LOG_TRACE;
  case PSLOG_BENCH_LIBLOGGER_DEBUG:
    return S_LOG_DEBUG;
  case PSLOG_BENCH_LIBLOGGER_INFO:
    return S_LOG_INFO;
  case PSLOG_BENCH_LIBLOGGER_WARN:
    return S_LOG_WARN;
  case PSLOG_BENCH_LIBLOGGER_ERROR:
    return S_LOG_ERROR;
  default:
    return S_LOG_INFO;
  }
}

static struct s_log_field_t *pslog_bench_liblogger_build_field(
    const pslog_bench_liblogger_field_spec *spec) {
  if (spec == NULL) {
    return NULL;
  }
  switch (spec->kind) {
  case PSLOG_BENCH_LIBLOGGER_FIELD_STRING:
    return s_log_string(spec->key,
                        spec->string_value != NULL ? spec->string_value : "");
  case PSLOG_BENCH_LIBLOGGER_FIELD_BOOL:
    return s_log_string(spec->key, spec->bool_value ? "true" : "false");
  case PSLOG_BENCH_LIBLOGGER_FIELD_INT64:
    return s_log_int64(spec->key, (long long)spec->signed_value);
  case PSLOG_BENCH_LIBLOGGER_FIELD_UINT64:
    return s_log_uint64(spec->key, (unsigned long long)spec->unsigned_value);
  case PSLOG_BENCH_LIBLOGGER_FIELD_FLOAT64:
    return s_log_double(spec->key, spec->double_value);
  case PSLOG_BENCH_LIBLOGGER_FIELD_NULL:
  default:
    return s_log_string(spec->key, "null");
  }
}

static void pslog_bench_liblogger_reallog(char *level,
                                          struct s_log_field_t **fields,
                                          size_t field_count) {
  switch (field_count) {
  case 1u:
    reallog(level, fields[0], NULL);
    return;
  case 2u:
    reallog(level, fields[0], fields[1], NULL);
    return;
  case 3u:
    reallog(level, fields[0], fields[1], fields[2], NULL);
    return;
  case 4u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], NULL);
    return;
  case 5u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4], NULL);
    return;
  case 6u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], NULL);
    return;
  case 7u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], NULL);
    return;
  case 8u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], NULL);
    return;
  case 9u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], NULL);
    return;
  case 10u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], NULL);
    return;
  case 11u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], fields[10],
            NULL);
    return;
  case 12u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], fields[10],
            fields[11], NULL);
    return;
  case 13u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], fields[10],
            fields[11], fields[12], NULL);
    return;
  case 14u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], fields[10],
            fields[11], fields[12], fields[13], NULL);
    return;
  case 15u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], fields[10],
            fields[11], fields[12], fields[13], fields[14], NULL);
    return;
  case 16u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], fields[10],
            fields[11], fields[12], fields[13], fields[14], fields[15], NULL);
    return;
  case 17u:
    reallog(level, fields[0], fields[1], fields[2], fields[3], fields[4],
            fields[5], fields[6], fields[7], fields[8], fields[9], fields[10],
            fields[11], fields[12], fields[13], fields[14], fields[15],
            fields[16], NULL);
    return;
  default:
    return;
  }
}

int pslog_bench_liblogger_available(void) { return 1; }

int pslog_bench_liblogger_run(const pslog_bench_liblogger_entry_spec *entries,
                              size_t entry_count, unsigned long long iterations,
                              pslog_bench_liblogger_result *result) {
  pslog_bench_liblogger_sink sink;
  unsigned long long start_ns;
  unsigned long long end_ns;
  unsigned long long i;
  size_t index;
  size_t max_fields;
  size_t slot;
  struct s_log_field_t **native_fields;

  if (entries == NULL || entry_count == 0u || result == NULL) {
    return -1;
  }

  if (pslog_bench_liblogger_sink_open(&sink) != 0) {
    return -1;
  }
  s_log_init(sink.stream);

  max_fields = 0u;
  for (slot = 0u; slot < entry_count; ++slot) {
    if (entries[slot].field_count > max_fields) {
      max_fields = entries[slot].field_count;
    }
  }

  native_fields =
      (struct s_log_field_t **)calloc(max_fields + 1u, sizeof(*native_fields));
  if (native_fields == NULL) {
    pslog_bench_liblogger_sink_close(&sink);
    return -1;
  }
  if (max_fields + 1u > 17u) {
    free(native_fields);
    pslog_bench_liblogger_sink_close(&sink);
    return -1;
  }

  index = 0u;
  start_ns = pslog_bench_liblogger_now_ns();
  for (i = 0ull; i < iterations; ++i) {
    const pslog_bench_liblogger_entry_spec *entry;

    entry = &entries[index];
    native_fields[0] =
        s_log_string("msg", entry->message != NULL ? entry->message : "");
    for (slot = 0u; slot < entry->field_count; ++slot) {
      native_fields[slot + 1u] =
          pslog_bench_liblogger_build_field(&entry->fields[slot]);
    }
    pslog_bench_liblogger_reallog(
        pslog_bench_liblogger_level_name(entry->level), native_fields,
        entry->field_count + 1u);

    ++index;
    if (index == entry_count) {
      index = 0u;
    }
  }
  end_ns = pslog_bench_liblogger_now_ns();

  result->elapsed_ns = end_ns - start_ns;
  result->bytes = pslog_bench_liblogger_sink_bytes(&sink);
  result->ops = iterations;

  free(native_fields);
  pslog_bench_liblogger_sink_close(&sink);
  return 0;
}
