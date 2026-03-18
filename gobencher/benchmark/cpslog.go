package benchmark

/*
#cgo CFLAGS: -O3 -DNDEBUG -I${SRCDIR}/../../include -I${SRCDIR}/../../build/release/generated/include
#cgo LDFLAGS: ${SRCDIR}/../../build/release/libpslog.a

#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "pslog.h"

typedef struct gobench_sink {
	unsigned long long bytes;
	unsigned long long writes;
	char *capture;
	size_t capture_cap;
	size_t capture_len;
} gobench_sink;

typedef struct gobench_result {
	unsigned long long elapsed_ns;
	unsigned long long bytes;
	unsigned long long writes;
	unsigned long long ops;
} gobench_result;

typedef struct gobench_logger {
	gobench_sink *sink;
	pslog_logger *logger;
	int owns_sink;
} gobench_logger;

typedef struct gobench_entry {
	pslog_level level;
	const char *msg;
	const pslog_field *fields;
	size_t field_count;
} gobench_entry;

typedef enum gobench_field_kind {
	GOBENCH_FIELD_NULL = 0,
	GOBENCH_FIELD_STRING = 1,
	GOBENCH_FIELD_BOOL = 2,
	GOBENCH_FIELD_INT64 = 3,
	GOBENCH_FIELD_UINT64 = 4,
	GOBENCH_FIELD_FLOAT64 = 5
} gobench_field_kind;

typedef struct gobench_raw_field {
	const char *key;
	gobench_field_kind kind;
	const char *string_value;
	int bool_value;
	long long int64_value;
	unsigned long long uint64_value;
	double float64_value;
} gobench_raw_field;

typedef struct gobench_raw_entry {
	pslog_level level;
	const char *msg;
	const gobench_raw_field *fields;
	size_t field_count;
} gobench_raw_entry;

typedef struct gobench_palette_snapshot {
	const char *key;
	const char *string_value;
	const char *number;
	const char *boolean;
	const char *null_value;
	const char *trace;
	const char *debug;
	const char *info;
	const char *warn;
	const char *error;
	const char *fatal;
	const char *panic;
	const char *timestamp;
	const char *message;
	const char *message_key;
	const char *reset;
} gobench_palette_snapshot;

static int gobench_sink_write(void *userdata, const char *data, size_t len, size_t *written) {
	gobench_sink *sink = (gobench_sink *)userdata;
	(void)data;
	if (sink != NULL) {
		sink->bytes += (unsigned long long)len;
		sink->writes += 1ull;
		if (sink->capture != NULL && sink->capture_cap > sink->capture_len + 1u) {
			size_t remaining = sink->capture_cap - sink->capture_len - 1u;
			size_t copy_len = len;
			if (copy_len > remaining) {
				copy_len = remaining;
			}
			if (copy_len > 0u) {
				memcpy(sink->capture + sink->capture_len, data, copy_len);
				sink->capture_len += copy_len;
				sink->capture[sink->capture_len] = '\0';
			}
		}
	}
	if (written != NULL) {
		*written = len;
	}
	return 0;
}

static int gobench_bench_sink_write(void *userdata, const char *data, size_t len, size_t *written) {
	gobench_sink *sink = (gobench_sink *)userdata;
	(void)data;
	if (sink != NULL) {
		sink->bytes += (unsigned long long)len;
		sink->writes += 1ull;
	}
	if (written != NULL) {
		*written = len;
	}
	return 0;
}

static int gobench_sink_isatty(void *userdata) {
	(void)userdata;
	return 0;
}

static unsigned long long gobench_monotonic_ns(void) {
#if defined(CLOCK_MONOTONIC)
	struct timespec ts;
	if (clock_gettime(CLOCK_MONOTONIC, &ts) == 0) {
		return ((unsigned long long)ts.tv_sec * 1000000000ull) + (unsigned long long)ts.tv_nsec;
	}
#endif
	return ((unsigned long long)clock() * 1000000000ull) / (unsigned long long)CLOCKS_PER_SEC;
}

static pslog_level gobench_level_from_int(int level) {
	switch (level) {
	case -1:
		return PSLOG_LEVEL_TRACE;
	case 0:
		return PSLOG_LEVEL_DEBUG;
	case 1:
		return PSLOG_LEVEL_INFO;
	case 2:
		return PSLOG_LEVEL_WARN;
	case 3:
		return PSLOG_LEVEL_ERROR;
	case 4:
		return PSLOG_LEVEL_FATAL;
	case 5:
		return PSLOG_LEVEL_PANIC;
	case 6:
		return PSLOG_LEVEL_NOLEVEL;
	default:
		return PSLOG_LEVEL_INFO;
	}
}

static pslog_field gobench_field_null(const char *key) {
	return pslog_null(key);
}

static pslog_field gobench_field_str(const char *key, const char *value) {
	return pslog_str(key, value);
}

static pslog_field gobench_field_bool(const char *key, int value) {
	return pslog_bool(key, value);
}

static pslog_field gobench_field_i64(const char *key, long long value) {
	return pslog_i64(key, (long)value);
}

static pslog_field gobench_field_u64(const char *key, unsigned long long value) {
	return pslog_u64(key, (unsigned long)value);
}

static pslog_field gobench_field_f64(const char *key, double value) {
	return pslog_f64(key, value);
}

static void gobench_field_set_null(pslog_field *dst, const char *key) {
	if (dst != NULL) {
		*dst = gobench_field_null(key);
	}
}

static void gobench_field_set_str(pslog_field *dst, const char *key, const char *value) {
	if (dst != NULL) {
		*dst = gobench_field_str(key, value);
	}
}

static void gobench_field_set_bool(pslog_field *dst, const char *key, int value) {
	if (dst != NULL) {
		*dst = gobench_field_bool(key, value);
	}
}

static void gobench_field_set_i64(pslog_field *dst, const char *key, long long value) {
	if (dst != NULL) {
		*dst = gobench_field_i64(key, value);
	}
}

static void gobench_field_set_u64(pslog_field *dst, const char *key, unsigned long long value) {
	if (dst != NULL) {
		*dst = gobench_field_u64(key, value);
	}
}

static void gobench_field_set_f64(pslog_field *dst, const char *key, double value) {
	if (dst != NULL) {
		*dst = gobench_field_f64(key, value);
	}
}

static void gobench_entry_init(gobench_entry *entry,
                               pslog_level level,
                               const char *msg,
                               const pslog_field *fields,
                               size_t field_count) {
	if (entry == NULL) {
		return;
	}
	entry->level = level;
	entry->msg = msg;
	entry->fields = fields;
	entry->field_count = field_count;
}

static void gobench_raw_field_set_null(gobench_raw_field *dst, const char *key) {
	if (dst == NULL) {
		return;
	}
	dst->key = key;
	dst->kind = GOBENCH_FIELD_NULL;
	dst->string_value = NULL;
	dst->bool_value = 0;
	dst->int64_value = 0;
	dst->uint64_value = 0ull;
	dst->float64_value = 0.0;
}

static void gobench_raw_field_set_str(gobench_raw_field *dst, const char *key, const char *value) {
	if (dst == NULL) {
		return;
	}
	dst->key = key;
	dst->kind = GOBENCH_FIELD_STRING;
	dst->string_value = value;
	dst->bool_value = 0;
	dst->int64_value = 0;
	dst->uint64_value = 0ull;
	dst->float64_value = 0.0;
}

static void gobench_raw_field_set_bool(gobench_raw_field *dst, const char *key, int value) {
	if (dst == NULL) {
		return;
	}
	dst->key = key;
	dst->kind = GOBENCH_FIELD_BOOL;
	dst->string_value = NULL;
	dst->bool_value = value;
	dst->int64_value = 0;
	dst->uint64_value = 0ull;
	dst->float64_value = 0.0;
}

static void gobench_raw_field_set_i64(gobench_raw_field *dst, const char *key, long long value) {
	if (dst == NULL) {
		return;
	}
	dst->key = key;
	dst->kind = GOBENCH_FIELD_INT64;
	dst->string_value = NULL;
	dst->bool_value = 0;
	dst->int64_value = value;
	dst->uint64_value = 0ull;
	dst->float64_value = 0.0;
}

static void gobench_raw_field_set_u64(gobench_raw_field *dst, const char *key, unsigned long long value) {
	if (dst == NULL) {
		return;
	}
	dst->key = key;
	dst->kind = GOBENCH_FIELD_UINT64;
	dst->string_value = NULL;
	dst->bool_value = 0;
	dst->int64_value = 0;
	dst->uint64_value = value;
	dst->float64_value = 0.0;
}

static void gobench_raw_field_set_f64(gobench_raw_field *dst, const char *key, double value) {
	if (dst == NULL) {
		return;
	}
	dst->key = key;
	dst->kind = GOBENCH_FIELD_FLOAT64;
	dst->string_value = NULL;
	dst->bool_value = 0;
	dst->int64_value = 0;
	dst->uint64_value = 0ull;
	dst->float64_value = value;
}

static void gobench_raw_entry_init(gobench_raw_entry *entry,
                                   pslog_level level,
                                   const char *msg,
                                   const gobench_raw_field *fields,
                                   size_t field_count) {
	if (entry == NULL) {
		return;
	}
	entry->level = level;
	entry->msg = msg;
	entry->fields = fields;
	entry->field_count = field_count;
}

static void gobench_build_field(pslog_field *dst, const gobench_raw_field *src) {
	if (dst == NULL || src == NULL) {
		return;
	}
	switch (src->kind) {
	case GOBENCH_FIELD_NULL:
		*dst = pslog_null(src->key);
		break;
	case GOBENCH_FIELD_STRING:
		*dst = pslog_str(src->key, src->string_value);
		break;
	case GOBENCH_FIELD_BOOL:
		*dst = pslog_bool(src->key, src->bool_value);
		break;
	case GOBENCH_FIELD_INT64:
		*dst = pslog_i64(src->key, (long)src->int64_value);
		break;
	case GOBENCH_FIELD_UINT64:
		*dst = pslog_u64(src->key, (unsigned long)src->uint64_value);
		break;
	case GOBENCH_FIELD_FLOAT64:
		*dst = pslog_f64(src->key, src->float64_value);
		break;
	default:
		*dst = pslog_null(src->key);
		break;
	}
}

static void gobench_build_fields(pslog_field *dst, const gobench_raw_field *src, size_t count) {
	size_t i;

	for (i = 0u; i < count; ++i) {
		gobench_build_field(dst + i, src + i);
	}
}

static gobench_logger *gobench_logger_alloc(void) {
	return (gobench_logger *)calloc(1u, sizeof(gobench_logger));
}

static void gobench_palette_snapshot_fill(gobench_palette_snapshot *dst,
                                          const pslog_palette *src) {
	if (dst == NULL || src == NULL) {
		return;
	}
	dst->key = src->key;
	dst->string_value = src->string;
	dst->number = src->number;
	dst->boolean = src->boolean;
	dst->null_value = src->null_value;
	dst->trace = src->trace;
	dst->debug = src->debug;
	dst->info = src->info;
	dst->warn = src->warn;
	dst->error = src->error;
	dst->fatal = src->fatal;
	dst->panic = src->panic;
	dst->timestamp = src->timestamp;
	dst->message = src->message;
	dst->message_key = src->message_key;
	dst->reset = src->reset;
}

static int gobench_palette_lookup(const char *name, gobench_palette_snapshot *out) {
	const pslog_palette *palette;

	if (out == NULL) {
		return -1;
	}
	palette = pslog_palette_by_name(name);
	if (palette == NULL) {
		return -1;
	}
	gobench_palette_snapshot_fill(out, palette);
	return 0;
}

static gobench_logger *gobench_logger_new_with_palette_and_timestamps(
                                                       int mode,
                                                       int color,
                                                       const char *palette_name,
                                                       int timestamps) {
	pslog_config config;
	gobench_logger *wrapper;
	gobench_sink *sink;

	wrapper = gobench_logger_alloc();
	sink = (gobench_sink *)calloc(1u, sizeof(gobench_sink));
	if (wrapper == NULL || sink == NULL) {
		free(wrapper);
		free(sink);
		return NULL;
	}

	pslog_default_config(&config);
	config.mode = mode ? PSLOG_MODE_JSON : PSLOG_MODE_CONSOLE;
	config.color = color ? PSLOG_COLOR_ALWAYS : PSLOG_COLOR_NEVER;
	config.timestamps = timestamps;
	config.min_level = PSLOG_LEVEL_TRACE;
	config.palette = pslog_palette_by_name(palette_name);
	config.output.write = gobench_sink_write;
	config.output.close = NULL;
	config.output.isatty = gobench_sink_isatty;
	config.output.userdata = sink;

	wrapper->logger = pslog_new(&config);
	if (wrapper->logger == NULL) {
		free(wrapper);
		free(sink);
		return NULL;
	}
	wrapper->sink = sink;
	wrapper->owns_sink = 1;
	return wrapper;
}

static gobench_logger *gobench_logger_new_with_palette(int mode,
                                                       int color,
                                                       const char *palette_name) {
	return gobench_logger_new_with_palette_and_timestamps(mode, color,
	                                                     palette_name, 1);
}

static gobench_logger *gobench_logger_new(int mode, int color) {
	return gobench_logger_new_with_palette(mode, color, "default");
}

static gobench_logger *gobench_bench_logger_new_with_palette_and_timestamps(
                                                       int mode,
                                                       int color,
                                                       const char *palette_name,
                                                       int timestamps) {
	pslog_config config;
	gobench_logger *wrapper;
	gobench_sink *sink;

	wrapper = gobench_logger_alloc();
	sink = (gobench_sink *)calloc(1u, sizeof(gobench_sink));
	if (wrapper == NULL || sink == NULL) {
		free(wrapper);
		free(sink);
		return NULL;
	}

	pslog_default_config(&config);
	config.mode = mode ? PSLOG_MODE_JSON : PSLOG_MODE_CONSOLE;
	config.color = color ? PSLOG_COLOR_ALWAYS : PSLOG_COLOR_NEVER;
	config.timestamps = timestamps;
	config.min_level = PSLOG_LEVEL_TRACE;
	config.palette = pslog_palette_by_name(palette_name);
	config.output.write = gobench_bench_sink_write;
	config.output.close = NULL;
	config.output.isatty = NULL;
	config.output.userdata = sink;

	wrapper->logger = pslog_new(&config);
	if (wrapper->logger == NULL) {
		free(wrapper);
		free(sink);
		return NULL;
	}
	wrapper->sink = sink;
	wrapper->owns_sink = 1;
	return wrapper;
}

static gobench_logger *gobench_logger_with(gobench_logger *base,
                                           const pslog_field *fields,
                                           size_t count) {
	gobench_logger *wrapper;

	if (base == NULL || base->logger == NULL) {
		return NULL;
	}

	wrapper = gobench_logger_alloc();
	if (wrapper == NULL) {
		return NULL;
	}
	wrapper->sink = base->sink;
	wrapper->owns_sink = 0;

	if (count == 0u) {
		wrapper->logger = pslog_with(base->logger, NULL, 0u);
		if (wrapper->logger == NULL) {
			free(wrapper);
			return NULL;
		}
		return wrapper;
	}
	wrapper->logger = pslog_with(base->logger, fields, count);
	if (wrapper->logger == NULL) {
		free(wrapper);
		return NULL;
	}
	return wrapper;
}

static void gobench_logger_destroy(gobench_logger *wrapper) {
	if (wrapper == NULL) {
		return;
	}
	if (wrapper->logger != NULL) {
		wrapper->logger->destroy(wrapper->logger);
	}
	if (wrapper->owns_sink && wrapper->sink != NULL) {
		free(wrapper->sink->capture);
		free(wrapper->sink);
	}
	free(wrapper);
}

static void gobench_logger_reset(gobench_logger *wrapper) {
	if (wrapper == NULL || wrapper->sink == NULL) {
		return;
	}
	wrapper->sink->bytes = 0ull;
	wrapper->sink->writes = 0ull;
	wrapper->sink->capture_len = 0u;
	if (wrapper->sink->capture != NULL && wrapper->sink->capture_cap > 0u) {
		wrapper->sink->capture[0] = '\0';
	}
}

static unsigned long long gobench_logger_bytes(gobench_logger *wrapper) {
	if (wrapper == NULL || wrapper->sink == NULL) {
		return 0ull;
	}
	return wrapper->sink->bytes;
}

static unsigned long long gobench_logger_writes(gobench_logger *wrapper) {
	if (wrapper == NULL || wrapper->sink == NULL) {
		return 0ull;
	}
	return wrapper->sink->writes;
}

static int gobench_logger_enable_capture(gobench_logger *wrapper, size_t capacity) {
	char *buffer;

	if (wrapper == NULL || wrapper->sink == NULL || capacity == 0u) {
		return -1;
	}
	buffer = (char *)calloc(capacity, sizeof(char));
	if (buffer == NULL) {
		return -1;
	}
	free(wrapper->sink->capture);
	wrapper->sink->capture = buffer;
	wrapper->sink->capture_cap = capacity;
	wrapper->sink->capture_len = 0u;
	return 0;
}

static const char *gobench_logger_capture_data(gobench_logger *wrapper) {
	if (wrapper == NULL || wrapper->sink == NULL || wrapper->sink->capture == NULL) {
		return "";
	}
	return wrapper->sink->capture;
}

static size_t gobench_logger_capture_len(gobench_logger *wrapper) {
	if (wrapper == NULL || wrapper->sink == NULL) {
		return 0u;
	}
	return wrapper->sink->capture_len;
}

static void gobench_logger_log(gobench_logger *wrapper,
                               pslog_level level,
                               const char *msg,
                               const pslog_field *fields,
                               size_t count) {
	if (wrapper == NULL || wrapper->logger == NULL) {
		return;
	}
	wrapper->logger->log(wrapper->logger, level, msg, fields, count);
}

static void gobench_logger_log_entry(gobench_logger *wrapper, const gobench_entry *entry) {
	if (entry == NULL) {
		return;
	}
	gobench_logger_log(wrapper, entry->level, entry->msg, entry->fields, entry->field_count);
}

static int gobench_logger_run(gobench_logger *wrapper,
                              const gobench_entry *entries,
                              size_t entry_count,
                              size_t iterations,
                              gobench_result *out) {
	unsigned long long start_ns;
	unsigned long long end_ns;
	size_t i;
	size_t index;
	pslog_logger *logger;

	if (wrapper == NULL || wrapper->logger == NULL || entries == NULL || entry_count == 0u || out == NULL) {
		return -1;
	}

	gobench_logger_reset(wrapper);
	logger = wrapper->logger;
	index = 0u;
	start_ns = gobench_monotonic_ns();
	for (i = 0u; i < iterations; ++i) {
		const gobench_entry *entry = entries + index;
		logger->log(logger, entry->level, entry->msg, entry->fields, entry->field_count);
		++index;
		if (index == entry_count) {
			index = 0u;
		}
	}
	end_ns = gobench_monotonic_ns();

	out->elapsed_ns = end_ns - start_ns;
	out->bytes = gobench_logger_bytes(wrapper);
	out->writes = gobench_logger_writes(wrapper);
	out->ops = (unsigned long long)iterations;
	return 0;
}

static int gobench_logger_run_raw(gobench_logger *wrapper,
                                  const gobench_raw_entry *entries,
                                  size_t entry_count,
                                  const gobench_raw_field *static_fields,
                                  size_t static_count,
                                  size_t iterations,
                                  gobench_result *out) {
	unsigned long long start_ns;
	unsigned long long end_ns;
	size_t i;
	size_t index;
	size_t max_fields;
	pslog_field *scratch;
	pslog_field *static_native;
	gobench_logger *active;
	gobench_logger *derived;
	pslog_logger *logger;

	if (wrapper == NULL || wrapper->logger == NULL || entries == NULL || entry_count == 0u || out == NULL) {
		return -1;
	}

	max_fields = 0u;
	for (i = 0u; i < entry_count; ++i) {
		if (entries[i].field_count > max_fields) {
			max_fields = entries[i].field_count;
		}
	}

	scratch = NULL;
	if (max_fields > 0u) {
		scratch = (pslog_field *)calloc(max_fields, sizeof(pslog_field));
		if (scratch == NULL) {
			return -1;
		}
	}

	active = wrapper;
	derived = NULL;
	static_native = NULL;
	if (static_count > 0u) {
		static_native = (pslog_field *)calloc(static_count, sizeof(pslog_field));
		if (static_native == NULL) {
			free(scratch);
			return -1;
		}
		gobench_build_fields(static_native, static_fields, static_count);
		derived = gobench_logger_with(wrapper, static_native, static_count);
		free(static_native);
		if (derived == NULL) {
			free(scratch);
			return -1;
		}
		active = derived;
	}

	gobench_logger_reset(active);
	logger = active->logger;
	index = 0u;
	start_ns = gobench_monotonic_ns();
	for (i = 0u; i < iterations; ++i) {
		const gobench_raw_entry *entry = entries + index;
		if (entry->field_count > 0u) {
			gobench_build_fields(scratch, entry->fields, entry->field_count);
		}
		logger->log(logger, entry->level, entry->msg, scratch, entry->field_count);
		++index;
		if (index == entry_count) {
			index = 0u;
		}
	}
	end_ns = gobench_monotonic_ns();

	out->elapsed_ns = end_ns - start_ns;
	out->bytes = gobench_logger_bytes(active);
	out->writes = gobench_logger_writes(active);
	out->ops = (unsigned long long)iterations;

	if (derived != NULL) {
		gobench_logger_destroy(derived);
	}
	free(scratch);
	return 0;
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	pslog "pkt.systems/pslog"
)

type CPreparedEntry struct {
	ptr *C.gobench_entry
}

type CPreparedFields struct {
	ptr   *C.pslog_field
	count C.size_t
}

type CPreparedData struct {
	Entries     []CPreparedEntry
	entriesPtr  *C.gobench_entry
	fieldsPtr   *C.pslog_field
	pool        map[string]*C.char
	allocations []unsafe.Pointer
}

type CRawField struct {
	ptr *C.gobench_raw_field
}

type CRawFields struct {
	ptr   *C.gobench_raw_field
	count C.size_t
}

type CRawEntry struct {
	ptr *C.gobench_raw_entry
}

type CRawData struct {
	Entries     []CRawEntry
	entriesPtr  *C.gobench_raw_entry
	fieldsPtr   *C.gobench_raw_field
	pool        map[string]*C.char
	allocations []unsafe.Pointer
}

type CLogger struct {
	ptr      *C.gobench_logger
	pool     map[string]*C.char
	pooled   []unsafe.Pointer
	ownsPool bool
	scratch  []C.pslog_field
}

type CPalette struct {
	Key        string
	String     string
	Num        string
	Bool       string
	Nil        string
	Trace      string
	Debug      string
	Info       string
	Warn       string
	Error      string
	Fatal      string
	Panic      string
	Timestamp  string
	Message    string
	MessageKey string
	Reset      string
}

type CRunResult struct {
	Elapsed      time.Duration
	BytesWritten uint64
	Writes       uint64
	Ops          uint64
}

func NewCPreparedData(entries []Entry) *CPreparedData {
	data := &CPreparedData{
		Entries: make([]CPreparedEntry, len(entries)),
		pool:    make(map[string]*C.char),
	}

	if len(entries) == 0 {
		return data
	}

	totalFields := 0
	for _, entry := range entries {
		totalFields += len(entry.Fields)
	}

	data.entriesPtr = (*C.gobench_entry)(data.alloc(C.size_t(len(entries)) * C.size_t(C.sizeof_gobench_entry)))

	var fieldBase unsafe.Pointer
	if totalFields > 0 {
		data.fieldsPtr = (*C.pslog_field)(data.alloc(C.size_t(totalFields) * C.size_t(C.sizeof_pslog_field)))
		fieldBase = unsafe.Pointer(data.fieldsPtr)
	}

	offset := 0
	for i, entry := range entries {
		var fieldsPtr *C.pslog_field

		if len(entry.Fields) > 0 {
			fieldsPtr = (*C.pslog_field)(unsafe.Add(fieldBase, uintptr(offset)*unsafe.Sizeof(*data.fieldsPtr)))
			data.fillNativeFields(fieldsPtr, entry.Fields)
			offset += len(entry.Fields)
		}
		data.Entries[i].ptr = (*C.gobench_entry)(unsafe.Add(unsafe.Pointer(data.entriesPtr), uintptr(i)*unsafe.Sizeof(*data.entriesPtr)))
		C.gobench_entry_init(data.Entries[i].ptr,
			C.gobench_level_from_int(C.int(entry.Level)),
			data.cstring(entry.Message),
			fieldsPtr,
			C.size_t(len(entry.Fields)))
	}

	return data
}

func NewCRawData(entries []Entry) *CRawData {
	data := &CRawData{
		Entries: make([]CRawEntry, len(entries)),
		pool:    make(map[string]*C.char),
	}

	if len(entries) == 0 {
		return data
	}

	totalFields := 0
	for _, entry := range entries {
		totalFields += len(entry.Fields)
	}

	data.entriesPtr = (*C.gobench_raw_entry)(data.alloc(C.size_t(len(entries)) * C.size_t(C.sizeof_gobench_raw_entry)))

	var fieldBase unsafe.Pointer
	if totalFields > 0 {
		data.fieldsPtr = (*C.gobench_raw_field)(data.alloc(C.size_t(totalFields) * C.size_t(C.sizeof_gobench_raw_field)))
		fieldBase = unsafe.Pointer(data.fieldsPtr)
	}

	offset := 0
	for i, entry := range entries {
		var fieldsPtr *C.gobench_raw_field

		if len(entry.Fields) > 0 {
			fieldsPtr = (*C.gobench_raw_field)(unsafe.Add(fieldBase, uintptr(offset)*unsafe.Sizeof(*data.fieldsPtr)))
			data.fillRawFields(fieldsPtr, entry.Fields)
			offset += len(entry.Fields)
		}
		data.Entries[i].ptr = (*C.gobench_raw_entry)(unsafe.Add(unsafe.Pointer(data.entriesPtr), uintptr(i)*unsafe.Sizeof(*data.entriesPtr)))
		C.gobench_raw_entry_init(data.Entries[i].ptr,
			C.gobench_level_from_int(C.int(entry.Level)),
			data.cstring(entry.Message),
			fieldsPtr,
			C.size_t(len(entry.Fields)))
	}

	return data
}

func (d *CPreparedData) PrepareFields(fields []Field) CPreparedFields {
	var prepared CPreparedFields
	if d == nil || len(fields) == 0 {
		return prepared
	}
	prepared.count = C.size_t(len(fields))
	prepared.ptr = (*C.pslog_field)(d.alloc(C.size_t(len(fields)) * C.size_t(C.sizeof_pslog_field)))
	d.fillNativeFields(prepared.ptr, fields)
	return prepared
}

func (d *CRawData) PrepareFields(fields []Field) CRawFields {
	var prepared CRawFields
	if d == nil || len(fields) == 0 {
		return prepared
	}
	prepared.count = C.size_t(len(fields))
	prepared.ptr = (*C.gobench_raw_field)(d.alloc(C.size_t(len(fields)) * C.size_t(C.sizeof_gobench_raw_field)))
	d.fillRawFields(prepared.ptr, fields)
	return prepared
}

func (d *CPreparedData) Close() {
	if d == nil {
		return
	}
	for _, value := range d.pool {
		C.free(unsafe.Pointer(value))
	}
	for _, ptr := range d.allocations {
		C.free(ptr)
	}
	d.pool = nil
	d.allocations = nil
	d.Entries = nil
	d.entriesPtr = nil
	d.fieldsPtr = nil
}

func (d *CRawData) Close() {
	if d == nil {
		return
	}
	for _, value := range d.pool {
		C.free(unsafe.Pointer(value))
	}
	for _, ptr := range d.allocations {
		C.free(ptr)
	}
	d.pool = nil
	d.allocations = nil
	d.Entries = nil
	d.entriesPtr = nil
	d.fieldsPtr = nil
}

func NewCLogger(variant Variant) (*CLogger, error) {
	return newCLoggerWithPaletteAndTimestamps(variant, "default", true)
}

func NewCBenchLogger(variant Variant) (*CLogger, error) {
	return newCBenchLoggerWithPaletteAndTimestamps(variant, "default", true)
}

func NewCLoggerWithPalette(variant Variant, paletteName string) (*CLogger, error) {
	return newCLoggerWithPaletteAndTimestamps(variant, paletteName, true)
}

func newCLoggerWithPaletteAndTimestamps(variant Variant, paletteName string, timestamps bool) (*CLogger, error) {
	mode := C.int(0)
	color := C.int(0)
	timestampMode := C.int(0)
	paletteNameC := C.CString(paletteName)
	if variant.Mode == pslog.ModeStructured {
		mode = 1
	}
	if variant.Color {
		color = 1
	}
	if timestamps {
		timestampMode = 1
	}
	defer C.free(unsafe.Pointer(paletteNameC))
	ptr := C.gobench_logger_new_with_palette_and_timestamps(mode, color, paletteNameC, timestampMode)
	if ptr == nil {
		return nil, fmt.Errorf("create C logger")
	}
	return &CLogger{
		ptr:      ptr,
		pool:     make(map[string]*C.char),
		ownsPool: true,
	}, nil
}

func newCBenchLoggerWithPaletteAndTimestamps(variant Variant, paletteName string, timestamps bool) (*CLogger, error) {
	mode := C.int(0)
	color := C.int(0)
	timestampMode := C.int(0)
	paletteNameC := C.CString(paletteName)
	if variant.Mode == pslog.ModeStructured {
		mode = 1
	}
	if variant.Color {
		color = 1
	}
	if timestamps {
		timestampMode = 1
	}
	defer C.free(unsafe.Pointer(paletteNameC))
	ptr := C.gobench_bench_logger_new_with_palette_and_timestamps(mode, color, paletteNameC, timestampMode)
	if ptr == nil {
		return nil, fmt.Errorf("create C benchmark logger")
	}
	return &CLogger{
		ptr:      ptr,
		pool:     make(map[string]*C.char),
		ownsPool: true,
	}, nil
}

func CPaletteByName(name string) (CPalette, error) {
	var snapshot C.gobench_palette_snapshot
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))

	if C.gobench_palette_lookup(nameC, &snapshot) != 0 {
		return CPalette{}, fmt.Errorf("lookup C palette %q", name)
	}

	return CPalette{
		Key:        C.GoString(snapshot.key),
		String:     C.GoString(snapshot.string_value),
		Num:        C.GoString(snapshot.number),
		Bool:       C.GoString(snapshot.boolean),
		Nil:        C.GoString(snapshot.null_value),
		Trace:      C.GoString(snapshot.trace),
		Debug:      C.GoString(snapshot.debug),
		Info:       C.GoString(snapshot.info),
		Warn:       C.GoString(snapshot.warn),
		Error:      C.GoString(snapshot.error),
		Fatal:      C.GoString(snapshot.fatal),
		Panic:      C.GoString(snapshot.panic),
		Timestamp:  C.GoString(snapshot.timestamp),
		Message:    C.GoString(snapshot.message),
		MessageKey: C.GoString(snapshot.message_key),
		Reset:      C.GoString(snapshot.reset),
	}, nil
}

func (l *CLogger) With(fields []Field) (*CLogger, error) {
	var native []C.pslog_field
	var nativePtr *C.pslog_field

	if l == nil || l.ptr == nil {
		return nil, fmt.Errorf("nil C logger")
	}

	native = l.buildNativeFields(fields)
	if len(native) > 0 {
		nativePtr = &native[0]
	}

	next := C.gobench_logger_with(l.ptr, nativePtr, C.size_t(len(native)))
	runtime.KeepAlive(native)
	if next == nil {
		return nil, fmt.Errorf("derive C logger")
	}
	return &CLogger{
		ptr:      next,
		pool:     l.pool,
		pooled:   l.pooled,
		ownsPool: false,
	}, nil
}

func (l *CLogger) WithPrepared(fields CPreparedFields) (*CLogger, error) {
	if l == nil || l.ptr == nil {
		return nil, fmt.Errorf("nil C logger")
	}
	next := C.gobench_logger_with(l.ptr, fields.ptr, fields.count)
	if next == nil {
		return nil, fmt.Errorf("derive C logger")
	}
	return &CLogger{
		ptr:      next,
		pool:     l.pool,
		pooled:   l.pooled,
		ownsPool: false,
	}, nil
}

func (l *CLogger) Log(entry Entry) {
	var native []C.pslog_field
	var nativePtr *C.pslog_field

	if l == nil || l.ptr == nil {
		return
	}

	native = l.buildNativeFields(entry.Fields)
	if len(native) > 0 {
		nativePtr = &native[0]
	}
	C.gobench_logger_log(l.ptr,
		C.gobench_level_from_int(C.int(entry.Level)),
		l.cstring(entry.Message),
		nativePtr,
		C.size_t(len(native)))
	runtime.KeepAlive(native)
}

func (l *CLogger) LogPrepared(entry CPreparedEntry) {
	if l == nil || l.ptr == nil || entry.ptr == nil {
		return
	}
	C.gobench_logger_log_entry(l.ptr, entry.ptr)
}

func (l *CLogger) RunPrepared(data *CPreparedData, iterations int) (CRunResult, error) {
	var result C.gobench_result

	if l == nil || l.ptr == nil {
		return CRunResult{}, fmt.Errorf("nil C logger")
	}
	if data == nil || data.entriesPtr == nil || len(data.Entries) == 0 {
		return CRunResult{}, fmt.Errorf("empty prepared data")
	}
	if iterations < 0 {
		iterations = 0
	}

	if C.gobench_logger_run(l.ptr, data.entriesPtr, C.size_t(len(data.Entries)), C.size_t(iterations), &result) != 0 {
		return CRunResult{}, fmt.Errorf("run prepared benchmark")
	}

	return CRunResult{
		Elapsed:      time.Duration(result.elapsed_ns),
		BytesWritten: uint64(result.bytes),
		Writes:       uint64(result.writes),
		Ops:          uint64(result.ops),
	}, nil
}

func (l *CLogger) RunRaw(data *CRawData, staticFields CRawFields, iterations int) (CRunResult, error) {
	var result C.gobench_result

	if l == nil || l.ptr == nil {
		return CRunResult{}, fmt.Errorf("nil C logger")
	}
	if data == nil || data.entriesPtr == nil || len(data.Entries) == 0 {
		return CRunResult{}, fmt.Errorf("empty raw data")
	}
	if iterations < 0 {
		iterations = 0
	}

	if C.gobench_logger_run_raw(l.ptr,
		data.entriesPtr,
		C.size_t(len(data.Entries)),
		staticFields.ptr,
		staticFields.count,
		C.size_t(iterations),
		&result) != 0 {
		return CRunResult{}, fmt.Errorf("run raw benchmark")
	}

	return CRunResult{
		Elapsed:      time.Duration(result.elapsed_ns),
		BytesWritten: uint64(result.bytes),
		Writes:       uint64(result.writes),
		Ops:          uint64(result.ops),
	}, nil
}

func (l *CLogger) Reset() {
	if l != nil && l.ptr != nil {
		C.gobench_logger_reset(l.ptr)
	}
}

func (l *CLogger) BytesWritten() uint64 {
	if l == nil || l.ptr == nil {
		return 0
	}
	return uint64(C.gobench_logger_bytes(l.ptr))
}

func (l *CLogger) Writes() uint64 {
	if l == nil || l.ptr == nil {
		return 0
	}
	return uint64(C.gobench_logger_writes(l.ptr))
}

func (l *CLogger) EnableCapture(capacity int) error {
	if l == nil || l.ptr == nil {
		return fmt.Errorf("nil C logger")
	}
	if capacity <= 0 {
		return fmt.Errorf("invalid capture capacity")
	}
	if C.gobench_logger_enable_capture(l.ptr, C.size_t(capacity)) != 0 {
		return fmt.Errorf("enable capture")
	}
	return nil
}

func (l *CLogger) Captured() string {
	if l == nil || l.ptr == nil {
		return ""
	}
	return C.GoStringN(C.gobench_logger_capture_data(l.ptr), C.int(C.gobench_logger_capture_len(l.ptr)))
}

func (l *CLogger) Close() {
	if l != nil && l.ptr != nil {
		C.gobench_logger_destroy(l.ptr)
		l.ptr = nil
	}
	if l != nil && l.ownsPool {
		for _, ptr := range l.pooled {
			C.free(ptr)
		}
		l.pool = nil
		l.pooled = nil
	}
}

func (d *CPreparedData) alloc(size C.size_t) unsafe.Pointer {
	var ptr unsafe.Pointer
	if size == 0 {
		return nil
	}
	ptr = C.malloc(size)
	d.allocations = append(d.allocations, ptr)
	return ptr
}

func (d *CRawData) alloc(size C.size_t) unsafe.Pointer {
	var ptr unsafe.Pointer
	if size == 0 {
		return nil
	}
	ptr = C.malloc(size)
	d.allocations = append(d.allocations, ptr)
	return ptr
}

func (d *CPreparedData) cstring(text string) *C.char {
	if existing, ok := d.pool[text]; ok {
		return existing
	}
	value := C.CString(text)
	d.pool[text] = value
	return value
}

func (d *CRawData) cstring(text string) *C.char {
	if existing, ok := d.pool[text]; ok {
		return existing
	}
	value := C.CString(text)
	d.pool[text] = value
	return value
}

func (l *CLogger) cstring(text string) *C.char {
	if l == nil {
		return nil
	}
	if existing, ok := l.pool[text]; ok {
		return existing
	}
	value := C.CString(text)
	l.pool[text] = value
	l.pooled = append(l.pooled, unsafe.Pointer(value))
	return value
}

func (l *CLogger) buildNativeFields(fields []Field) []C.pslog_field {
	var native []C.pslog_field

	if len(fields) == 0 {
		return nil
	}

	if cap(l.scratch) < len(fields) {
		l.scratch = make([]C.pslog_field, len(fields))
	}
	native = l.scratch[:len(fields)]
	for i, field := range fields {
		key := l.cstring(field.Key)
		dstField := &native[i]
		switch field.Kind {
		case FieldNull:
			C.gobench_field_set_null(dstField, key)
		case FieldString:
			C.gobench_field_set_str(dstField, key, l.cstring(field.String))
		case FieldBool:
			boolValue := 0
			if field.Bool {
				boolValue = 1
			}
			C.gobench_field_set_bool(dstField, key, C.int(boolValue))
		case FieldInt64:
			C.gobench_field_set_i64(dstField, key, C.longlong(field.Int64))
		case FieldUint64:
			C.gobench_field_set_u64(dstField, key, C.ulonglong(field.Uint64))
		case FieldFloat64:
			C.gobench_field_set_f64(dstField, key, C.double(field.Float64))
		default:
			C.gobench_field_set_null(dstField, key)
		}
	}
	return native
}

func (d *CPreparedData) fillNativeFields(dst *C.pslog_field, fields []Field) {
	for i, field := range fields {
		dstField := (*C.pslog_field)(unsafe.Add(unsafe.Pointer(dst), uintptr(i)*unsafe.Sizeof(*dst)))
		key := d.cstring(field.Key)
		switch field.Kind {
		case FieldNull:
			C.gobench_field_set_null(dstField, key)
		case FieldString:
			C.gobench_field_set_str(dstField, key, d.cstring(field.String))
		case FieldBool:
			boolValue := 0
			if field.Bool {
				boolValue = 1
			}
			C.gobench_field_set_bool(dstField, key, C.int(boolValue))
		case FieldInt64:
			C.gobench_field_set_i64(dstField, key, C.longlong(field.Int64))
		case FieldUint64:
			C.gobench_field_set_u64(dstField, key, C.ulonglong(field.Uint64))
		case FieldFloat64:
			C.gobench_field_set_f64(dstField, key, C.double(field.Float64))
		default:
			C.gobench_field_set_null(dstField, key)
		}
	}
}

func (d *CRawData) fillRawFields(dst *C.gobench_raw_field, fields []Field) {
	for i, field := range fields {
		dstField := (*C.gobench_raw_field)(unsafe.Add(unsafe.Pointer(dst), uintptr(i)*unsafe.Sizeof(*dst)))
		key := d.cstring(field.Key)
		switch field.Kind {
		case FieldNull:
			C.gobench_raw_field_set_null(dstField, key)
		case FieldString:
			C.gobench_raw_field_set_str(dstField, key, d.cstring(field.String))
		case FieldBool:
			boolValue := 0
			if field.Bool {
				boolValue = 1
			}
			C.gobench_raw_field_set_bool(dstField, key, C.int(boolValue))
		case FieldInt64:
			C.gobench_raw_field_set_i64(dstField, key, C.longlong(field.Int64))
		case FieldUint64:
			C.gobench_raw_field_set_u64(dstField, key, C.ulonglong(field.Uint64))
		case FieldFloat64:
			C.gobench_raw_field_set_f64(dstField, key, C.double(field.Float64))
		default:
			C.gobench_raw_field_set_null(dstField, key)
		}
	}
}
