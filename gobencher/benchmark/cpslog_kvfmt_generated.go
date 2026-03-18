package benchmark

/*
#include <stddef.h>
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

typedef struct gobench_kvfmt_arg {
	int kind;
	const char *string_value;
	int bool_value;
	long signed_value;
	unsigned long unsigned_value;
	double double_value;
} gobench_kvfmt_arg;

typedef struct gobench_kvfmt_entry {
	pslog_level level;
	const char *msg;
	const char *kvfmt;
	int signature_id;
	void (*call)(pslog_logger *logger, const struct gobench_kvfmt_entry *entry);
	const gobench_kvfmt_arg *args;
	size_t arg_count;
} gobench_kvfmt_entry;

static unsigned long long gobench_kvfmt_monotonic_ns(void) {
#if defined(CLOCK_MONOTONIC)
	struct timespec ts;
	if (clock_gettime(CLOCK_MONOTONIC, &ts) == 0) {
		return ((unsigned long long)ts.tv_sec * 1000000000ull) + (unsigned long long)ts.tv_nsec;
	}
#endif
	return ((unsigned long long)clock() * 1000000000ull) / (unsigned long long)CLOCKS_PER_SEC;
}

static void gobench_kvfmt_sink_reset(gobench_sink *sink) {
	if (sink == NULL) {
		return;
	}
	sink->bytes = 0ull;
	sink->writes = 0ull;
	sink->capture_len = 0u;
	if (sink->capture != NULL && sink->capture_cap > 0u) {
		sink->capture[0] = '\0';
	}
}

static void gobench_logger_tracef_sig_0(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].bool_value, args[1].string_value);
}

static void gobench_logger_debugf_sig_0(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].bool_value, args[1].string_value);
}

static void gobench_logger_infof_sig_0(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].bool_value, args[1].string_value);
}

static void gobench_logger_warnf_sig_0(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].bool_value, args[1].string_value);
}

static void gobench_logger_errorf_sig_0(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].bool_value, args[1].string_value);
}

static void gobench_logger_logf_sig_0(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].bool_value, args[1].string_value);
}

static void gobench_logger_tracef_sig_1(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].bool_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_debugf_sig_1(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].bool_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_infof_sig_1(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].bool_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_warnf_sig_1(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].bool_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_errorf_sig_1(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].bool_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_logf_sig_1(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].bool_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_tracef_sig_2(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].bool_value, args[2].string_value, args[3].bool_value, args[4].string_value);
}

static void gobench_logger_debugf_sig_2(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].bool_value, args[2].string_value, args[3].bool_value, args[4].string_value);
}

static void gobench_logger_infof_sig_2(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].bool_value, args[2].string_value, args[3].bool_value, args[4].string_value);
}

static void gobench_logger_warnf_sig_2(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].bool_value, args[2].string_value, args[3].bool_value, args[4].string_value);
}

static void gobench_logger_errorf_sig_2(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].bool_value, args[2].string_value, args[3].bool_value, args[4].string_value);
}

static void gobench_logger_logf_sig_2(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].bool_value, args[2].string_value, args[3].bool_value, args[4].string_value);
}

static void gobench_logger_tracef_sig_3(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_debugf_sig_3(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_infof_sig_3(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_warnf_sig_3(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_errorf_sig_3(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_logf_sig_3(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_tracef_sig_4(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_debugf_sig_4(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_infof_sig_4(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_warnf_sig_4(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_errorf_sig_4(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_logf_sig_4(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_tracef_sig_5(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_debugf_sig_5(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_infof_sig_5(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_warnf_sig_5(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_errorf_sig_5(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_logf_sig_5(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].bool_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_tracef_sig_6(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].double_value);
}

static void gobench_logger_debugf_sig_6(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].double_value);
}

static void gobench_logger_infof_sig_6(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].double_value);
}

static void gobench_logger_warnf_sig_6(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].double_value);
}

static void gobench_logger_errorf_sig_6(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].double_value);
}

static void gobench_logger_logf_sig_6(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].double_value);
}

static void gobench_logger_tracef_sig_7(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].signed_value, args[4].signed_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_debugf_sig_7(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].signed_value, args[4].signed_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_infof_sig_7(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].signed_value, args[4].signed_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_warnf_sig_7(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].signed_value, args[4].signed_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_errorf_sig_7(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].signed_value, args[4].signed_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_logf_sig_7(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].signed_value, args[4].signed_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_tracef_sig_8(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value);
}

static void gobench_logger_debugf_sig_8(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value);
}

static void gobench_logger_infof_sig_8(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value);
}

static void gobench_logger_warnf_sig_8(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value);
}

static void gobench_logger_errorf_sig_8(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value);
}

static void gobench_logger_logf_sig_8(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].bool_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value);
}

static void gobench_logger_tracef_sig_9(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].signed_value, args[3].string_value);
}

static void gobench_logger_debugf_sig_9(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].signed_value, args[3].string_value);
}

static void gobench_logger_infof_sig_9(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].signed_value, args[3].string_value);
}

static void gobench_logger_warnf_sig_9(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].signed_value, args[3].string_value);
}

static void gobench_logger_errorf_sig_9(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].signed_value, args[3].string_value);
}

static void gobench_logger_logf_sig_9(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].signed_value, args[3].string_value);
}

static void gobench_logger_tracef_sig_10(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].signed_value, args[5].string_value, args[6].signed_value, args[7].signed_value, args[8].string_value);
}

static void gobench_logger_debugf_sig_10(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].signed_value, args[5].string_value, args[6].signed_value, args[7].signed_value, args[8].string_value);
}

static void gobench_logger_infof_sig_10(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].signed_value, args[5].string_value, args[6].signed_value, args[7].signed_value, args[8].string_value);
}

static void gobench_logger_warnf_sig_10(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].signed_value, args[5].string_value, args[6].signed_value, args[7].signed_value, args[8].string_value);
}

static void gobench_logger_errorf_sig_10(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].signed_value, args[5].string_value, args[6].signed_value, args[7].signed_value, args[8].string_value);
}

static void gobench_logger_logf_sig_10(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].signed_value, args[5].string_value, args[6].signed_value, args[7].signed_value, args[8].string_value);
}

static void gobench_logger_tracef_sig_11(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value);
}

static void gobench_logger_debugf_sig_11(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value);
}

static void gobench_logger_infof_sig_11(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value);
}

static void gobench_logger_warnf_sig_11(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value);
}

static void gobench_logger_errorf_sig_11(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value);
}

static void gobench_logger_logf_sig_11(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value);
}

static void gobench_logger_tracef_sig_12(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_debugf_sig_12(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_infof_sig_12(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_warnf_sig_12(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_errorf_sig_12(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_logf_sig_12(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_tracef_sig_13(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_debugf_sig_13(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_infof_sig_13(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_warnf_sig_13(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_errorf_sig_13(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_logf_sig_13(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_tracef_sig_14(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_debugf_sig_14(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_infof_sig_14(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_warnf_sig_14(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_errorf_sig_14(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_logf_sig_14(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_tracef_sig_15(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value);
}

static void gobench_logger_debugf_sig_15(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value);
}

static void gobench_logger_infof_sig_15(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value);
}

static void gobench_logger_warnf_sig_15(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value);
}

static void gobench_logger_errorf_sig_15(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value);
}

static void gobench_logger_logf_sig_15(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value);
}

static void gobench_logger_tracef_sig_16(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_debugf_sig_16(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_infof_sig_16(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_warnf_sig_16(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_errorf_sig_16(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_logf_sig_16(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_tracef_sig_17(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_debugf_sig_17(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_infof_sig_17(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_warnf_sig_17(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_errorf_sig_17(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_logf_sig_17(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_tracef_sig_18(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_debugf_sig_18(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_infof_sig_18(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_warnf_sig_18(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_errorf_sig_18(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_logf_sig_18(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].bool_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_tracef_sig_19(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_debugf_sig_19(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_infof_sig_19(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_warnf_sig_19(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_errorf_sig_19(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_logf_sig_19(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_tracef_sig_20(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_debugf_sig_20(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_infof_sig_20(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_warnf_sig_20(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_errorf_sig_20(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_logf_sig_20(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_tracef_sig_21(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_debugf_sig_21(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_infof_sig_21(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_warnf_sig_21(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_errorf_sig_21(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_logf_sig_21(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_tracef_sig_22(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_debugf_sig_22(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_infof_sig_22(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_warnf_sig_22(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_errorf_sig_22(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_logf_sig_22(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_tracef_sig_23(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_debugf_sig_23(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_infof_sig_23(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_warnf_sig_23(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_errorf_sig_23(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_logf_sig_23(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_tracef_sig_24(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_debugf_sig_24(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_infof_sig_24(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_warnf_sig_24(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_errorf_sig_24(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_logf_sig_24(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_tracef_sig_25(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_debugf_sig_25(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_infof_sig_25(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_warnf_sig_25(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_errorf_sig_25(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_logf_sig_25(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].signed_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_tracef_sig_26(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value);
}

static void gobench_logger_debugf_sig_26(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value);
}

static void gobench_logger_infof_sig_26(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value);
}

static void gobench_logger_warnf_sig_26(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value);
}

static void gobench_logger_errorf_sig_26(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value);
}

static void gobench_logger_logf_sig_26(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value);
}

static void gobench_logger_tracef_sig_27(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_debugf_sig_27(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_infof_sig_27(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_warnf_sig_27(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_errorf_sig_27(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_logf_sig_27(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_tracef_sig_28(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_debugf_sig_28(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_infof_sig_28(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_warnf_sig_28(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_errorf_sig_28(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_logf_sig_28(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].bool_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_tracef_sig_29(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_debugf_sig_29(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_infof_sig_29(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_warnf_sig_29(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_errorf_sig_29(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_logf_sig_29(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_tracef_sig_30(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_debugf_sig_30(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_infof_sig_30(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_warnf_sig_30(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_errorf_sig_30(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_logf_sig_30(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_tracef_sig_31(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_debugf_sig_31(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_infof_sig_31(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_warnf_sig_31(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_errorf_sig_31(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_logf_sig_31(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_tracef_sig_32(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_debugf_sig_32(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_infof_sig_32(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_warnf_sig_32(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_errorf_sig_32(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_logf_sig_32(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_tracef_sig_33(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_debugf_sig_33(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_infof_sig_33(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_warnf_sig_33(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_errorf_sig_33(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_logf_sig_33(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].bool_value, args[4].bool_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_tracef_sig_34(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_debugf_sig_34(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_infof_sig_34(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_warnf_sig_34(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_errorf_sig_34(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_logf_sig_34(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].signed_value);
}

static void gobench_logger_tracef_sig_35(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_debugf_sig_35(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_infof_sig_35(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_warnf_sig_35(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_errorf_sig_35(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_logf_sig_35(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value);
}

static void gobench_logger_tracef_sig_36(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].signed_value);
}

static void gobench_logger_debugf_sig_36(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].signed_value);
}

static void gobench_logger_infof_sig_36(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].signed_value);
}

static void gobench_logger_warnf_sig_36(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].signed_value);
}

static void gobench_logger_errorf_sig_36(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].signed_value);
}

static void gobench_logger_logf_sig_36(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].signed_value);
}

static void gobench_logger_tracef_sig_37(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].bool_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_debugf_sig_37(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].bool_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_infof_sig_37(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].bool_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_warnf_sig_37(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].bool_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_errorf_sig_37(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].bool_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_logf_sig_37(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].bool_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_tracef_sig_38(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_debugf_sig_38(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_infof_sig_38(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_warnf_sig_38(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_errorf_sig_38(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_logf_sig_38(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value);
}

static void gobench_logger_tracef_sig_39(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_debugf_sig_39(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_infof_sig_39(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_warnf_sig_39(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_errorf_sig_39(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_logf_sig_39(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].signed_value);
}

static void gobench_logger_tracef_sig_40(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_debugf_sig_40(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_infof_sig_40(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_warnf_sig_40(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_errorf_sig_40(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_logf_sig_40(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_tracef_sig_41(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].bool_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_debugf_sig_41(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].bool_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_infof_sig_41(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].bool_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_warnf_sig_41(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].bool_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_errorf_sig_41(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].bool_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_logf_sig_41(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].bool_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_tracef_sig_42(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_debugf_sig_42(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_infof_sig_42(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_warnf_sig_42(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_errorf_sig_42(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_logf_sig_42(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_tracef_sig_43(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_debugf_sig_43(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_infof_sig_43(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_warnf_sig_43(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_errorf_sig_43(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_logf_sig_43(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_tracef_sig_44(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_debugf_sig_44(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_infof_sig_44(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_warnf_sig_44(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_errorf_sig_44(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_logf_sig_44(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_tracef_sig_45(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_debugf_sig_45(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_infof_sig_45(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_warnf_sig_45(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_errorf_sig_45(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_logf_sig_45(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_tracef_sig_46(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_debugf_sig_46(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_infof_sig_46(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_warnf_sig_46(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_errorf_sig_46(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_logf_sig_46(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_tracef_sig_47(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_debugf_sig_47(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_infof_sig_47(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_warnf_sig_47(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_errorf_sig_47(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_logf_sig_47(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_tracef_sig_48(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_debugf_sig_48(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_infof_sig_48(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_warnf_sig_48(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_errorf_sig_48(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_logf_sig_48(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_tracef_sig_49(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_debugf_sig_49(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_infof_sig_49(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_warnf_sig_49(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_errorf_sig_49(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_logf_sig_49(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].signed_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value);
}

static void gobench_logger_tracef_sig_50(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value);
}

static void gobench_logger_debugf_sig_50(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value);
}

static void gobench_logger_infof_sig_50(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value);
}

static void gobench_logger_warnf_sig_50(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value);
}

static void gobench_logger_errorf_sig_50(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value);
}

static void gobench_logger_logf_sig_50(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value);
}

static void gobench_logger_tracef_sig_51(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_debugf_sig_51(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_infof_sig_51(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_warnf_sig_51(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_errorf_sig_51(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_logf_sig_51(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_tracef_sig_52(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_debugf_sig_52(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_infof_sig_52(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_warnf_sig_52(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_errorf_sig_52(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_logf_sig_52(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_tracef_sig_53(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_debugf_sig_53(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_infof_sig_53(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_warnf_sig_53(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_errorf_sig_53(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_logf_sig_53(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_tracef_sig_54(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_debugf_sig_54(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_infof_sig_54(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_warnf_sig_54(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_errorf_sig_54(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_logf_sig_54(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_tracef_sig_55(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_debugf_sig_55(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_infof_sig_55(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_warnf_sig_55(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_errorf_sig_55(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_logf_sig_55(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_tracef_sig_56(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_debugf_sig_56(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_infof_sig_56(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_warnf_sig_56(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_errorf_sig_56(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_logf_sig_56(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_tracef_sig_57(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value, args[11].string_value, args[12].string_value, args[13].string_value);
}

static void gobench_logger_debugf_sig_57(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value, args[11].string_value, args[12].string_value, args[13].string_value);
}

static void gobench_logger_infof_sig_57(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value, args[11].string_value, args[12].string_value, args[13].string_value);
}

static void gobench_logger_warnf_sig_57(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value, args[11].string_value, args[12].string_value, args[13].string_value);
}

static void gobench_logger_errorf_sig_57(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value, args[11].string_value, args[12].string_value, args[13].string_value);
}

static void gobench_logger_logf_sig_57(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value, args[11].string_value, args[12].string_value, args[13].string_value);
}

static void gobench_logger_tracef_sig_58(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_debugf_sig_58(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_infof_sig_58(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_warnf_sig_58(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_errorf_sig_58(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_logf_sig_58(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].signed_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_tracef_sig_59(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value);
}

static void gobench_logger_debugf_sig_59(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value);
}

static void gobench_logger_infof_sig_59(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value);
}

static void gobench_logger_warnf_sig_59(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value);
}

static void gobench_logger_errorf_sig_59(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value);
}

static void gobench_logger_logf_sig_59(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value);
}

static void gobench_logger_tracef_sig_60(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_debugf_sig_60(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_infof_sig_60(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_warnf_sig_60(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_errorf_sig_60(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_logf_sig_60(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].bool_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_tracef_sig_61(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value);
}

static void gobench_logger_debugf_sig_61(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value);
}

static void gobench_logger_infof_sig_61(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value);
}

static void gobench_logger_warnf_sig_61(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value);
}

static void gobench_logger_errorf_sig_61(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value);
}

static void gobench_logger_logf_sig_61(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value);
}

static void gobench_logger_tracef_sig_62(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_debugf_sig_62(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_infof_sig_62(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_warnf_sig_62(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_errorf_sig_62(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_logf_sig_62(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_tracef_sig_63(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_debugf_sig_63(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_infof_sig_63(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_warnf_sig_63(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_errorf_sig_63(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_logf_sig_63(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].signed_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_tracef_sig_64(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_debugf_sig_64(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_infof_sig_64(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_warnf_sig_64(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_errorf_sig_64(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_logf_sig_64(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_tracef_sig_65(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_debugf_sig_65(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_infof_sig_65(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_warnf_sig_65(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_errorf_sig_65(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_logf_sig_65(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_tracef_sig_66(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_debugf_sig_66(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_infof_sig_66(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_warnf_sig_66(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_errorf_sig_66(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_logf_sig_66(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].signed_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].signed_value);
}

static void gobench_logger_tracef_sig_67(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_debugf_sig_67(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_infof_sig_67(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_warnf_sig_67(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_errorf_sig_67(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_logf_sig_67(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value);
}

static void gobench_logger_tracef_sig_68(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_debugf_sig_68(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_infof_sig_68(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_warnf_sig_68(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_errorf_sig_68(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_logf_sig_68(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].signed_value);
}

static void gobench_logger_tracef_sig_69(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_debugf_sig_69(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_infof_sig_69(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_warnf_sig_69(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_errorf_sig_69(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_logf_sig_69(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].signed_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].string_value, args[13].string_value, args[14].string_value, args[15].signed_value);
}

static void gobench_logger_tracef_sig_70(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_debugf_sig_70(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_infof_sig_70(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_warnf_sig_70(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_errorf_sig_70(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_logf_sig_70(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value);
}

static void gobench_logger_tracef_sig_71(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_debugf_sig_71(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_infof_sig_71(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_warnf_sig_71(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_errorf_sig_71(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_logf_sig_71(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].signed_value, args[6].string_value);
}

static void gobench_logger_tracef_sig_72(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_debugf_sig_72(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_infof_sig_72(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_warnf_sig_72(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_errorf_sig_72(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_logf_sig_72(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value);
}

static void gobench_logger_tracef_sig_73(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value);
}

static void gobench_logger_debugf_sig_73(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value);
}

static void gobench_logger_infof_sig_73(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value);
}

static void gobench_logger_warnf_sig_73(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value);
}

static void gobench_logger_errorf_sig_73(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value);
}

static void gobench_logger_logf_sig_73(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value);
}

static void gobench_logger_tracef_sig_74(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_debugf_sig_74(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_infof_sig_74(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_warnf_sig_74(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_errorf_sig_74(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_logf_sig_74(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].signed_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].signed_value);
}

static void gobench_logger_tracef_sig_75(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_debugf_sig_75(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_infof_sig_75(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_warnf_sig_75(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_errorf_sig_75(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_logf_sig_75(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value);
}

static void gobench_logger_tracef_sig_76(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_debugf_sig_76(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_infof_sig_76(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_warnf_sig_76(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_errorf_sig_76(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_logf_sig_76(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value);
}

static void gobench_logger_tracef_sig_77(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value);
}

static void gobench_logger_debugf_sig_77(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value);
}

static void gobench_logger_infof_sig_77(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value);
}

static void gobench_logger_warnf_sig_77(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value);
}

static void gobench_logger_errorf_sig_77(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value);
}

static void gobench_logger_logf_sig_77(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value);
}

static void gobench_logger_tracef_sig_78(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_debugf_sig_78(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_infof_sig_78(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_warnf_sig_78(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_errorf_sig_78(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_logf_sig_78(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].signed_value, args[9].string_value, args[10].string_value, args[11].string_value, args[12].signed_value);
}

static void gobench_logger_tracef_sig_79(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_debugf_sig_79(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_infof_sig_79(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_warnf_sig_79(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_errorf_sig_79(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_logf_sig_79(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value);
}

static void gobench_logger_tracef_sig_80(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_debugf_sig_80(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_infof_sig_80(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_warnf_sig_80(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_errorf_sig_80(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_logf_sig_80(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value);
}

static void gobench_logger_tracef_sig_81(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_debugf_sig_81(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_infof_sig_81(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_warnf_sig_81(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_errorf_sig_81(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_logf_sig_81(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value);
}

static void gobench_logger_tracef_sig_82(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->tracef(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_debugf_sig_82(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->debugf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_infof_sig_82(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->infof(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_warnf_sig_82(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->warnf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_errorf_sig_82(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->errorf(logger, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

static void gobench_logger_logf_sig_82(pslog_logger *logger,
				      const gobench_kvfmt_entry *entry) {
	const char *msg = entry->msg;
	const char *kvfmt = entry->kvfmt;
	const gobench_kvfmt_arg *args = entry->args;
	logger->logf(logger, entry->level, msg, kvfmt, args[0].string_value, args[1].string_value, args[2].string_value, args[3].string_value, args[4].string_value, args[5].string_value, args[6].string_value, args[7].string_value, args[8].string_value, args[9].string_value, args[10].string_value, args[11].string_value);
}

int gobench_kvfmt_entry_set_signature(gobench_kvfmt_entry *entry, int signature_id) {
	if (entry == NULL) {
		return -1;
	}
	entry->signature_id = signature_id;
	switch (entry->level) {
	case PSLOG_LEVEL_TRACE:
		switch (signature_id) {
		case 0:
			entry->call = gobench_logger_tracef_sig_0;
			return 0;
		case 1:
			entry->call = gobench_logger_tracef_sig_1;
			return 0;
		case 2:
			entry->call = gobench_logger_tracef_sig_2;
			return 0;
		case 3:
			entry->call = gobench_logger_tracef_sig_3;
			return 0;
		case 4:
			entry->call = gobench_logger_tracef_sig_4;
			return 0;
		case 5:
			entry->call = gobench_logger_tracef_sig_5;
			return 0;
		case 6:
			entry->call = gobench_logger_tracef_sig_6;
			return 0;
		case 7:
			entry->call = gobench_logger_tracef_sig_7;
			return 0;
		case 8:
			entry->call = gobench_logger_tracef_sig_8;
			return 0;
		case 9:
			entry->call = gobench_logger_tracef_sig_9;
			return 0;
		case 10:
			entry->call = gobench_logger_tracef_sig_10;
			return 0;
		case 11:
			entry->call = gobench_logger_tracef_sig_11;
			return 0;
		case 12:
			entry->call = gobench_logger_tracef_sig_12;
			return 0;
		case 13:
			entry->call = gobench_logger_tracef_sig_13;
			return 0;
		case 14:
			entry->call = gobench_logger_tracef_sig_14;
			return 0;
		case 15:
			entry->call = gobench_logger_tracef_sig_15;
			return 0;
		case 16:
			entry->call = gobench_logger_tracef_sig_16;
			return 0;
		case 17:
			entry->call = gobench_logger_tracef_sig_17;
			return 0;
		case 18:
			entry->call = gobench_logger_tracef_sig_18;
			return 0;
		case 19:
			entry->call = gobench_logger_tracef_sig_19;
			return 0;
		case 20:
			entry->call = gobench_logger_tracef_sig_20;
			return 0;
		case 21:
			entry->call = gobench_logger_tracef_sig_21;
			return 0;
		case 22:
			entry->call = gobench_logger_tracef_sig_22;
			return 0;
		case 23:
			entry->call = gobench_logger_tracef_sig_23;
			return 0;
		case 24:
			entry->call = gobench_logger_tracef_sig_24;
			return 0;
		case 25:
			entry->call = gobench_logger_tracef_sig_25;
			return 0;
		case 26:
			entry->call = gobench_logger_tracef_sig_26;
			return 0;
		case 27:
			entry->call = gobench_logger_tracef_sig_27;
			return 0;
		case 28:
			entry->call = gobench_logger_tracef_sig_28;
			return 0;
		case 29:
			entry->call = gobench_logger_tracef_sig_29;
			return 0;
		case 30:
			entry->call = gobench_logger_tracef_sig_30;
			return 0;
		case 31:
			entry->call = gobench_logger_tracef_sig_31;
			return 0;
		case 32:
			entry->call = gobench_logger_tracef_sig_32;
			return 0;
		case 33:
			entry->call = gobench_logger_tracef_sig_33;
			return 0;
		case 34:
			entry->call = gobench_logger_tracef_sig_34;
			return 0;
		case 35:
			entry->call = gobench_logger_tracef_sig_35;
			return 0;
		case 36:
			entry->call = gobench_logger_tracef_sig_36;
			return 0;
		case 37:
			entry->call = gobench_logger_tracef_sig_37;
			return 0;
		case 38:
			entry->call = gobench_logger_tracef_sig_38;
			return 0;
		case 39:
			entry->call = gobench_logger_tracef_sig_39;
			return 0;
		case 40:
			entry->call = gobench_logger_tracef_sig_40;
			return 0;
		case 41:
			entry->call = gobench_logger_tracef_sig_41;
			return 0;
		case 42:
			entry->call = gobench_logger_tracef_sig_42;
			return 0;
		case 43:
			entry->call = gobench_logger_tracef_sig_43;
			return 0;
		case 44:
			entry->call = gobench_logger_tracef_sig_44;
			return 0;
		case 45:
			entry->call = gobench_logger_tracef_sig_45;
			return 0;
		case 46:
			entry->call = gobench_logger_tracef_sig_46;
			return 0;
		case 47:
			entry->call = gobench_logger_tracef_sig_47;
			return 0;
		case 48:
			entry->call = gobench_logger_tracef_sig_48;
			return 0;
		case 49:
			entry->call = gobench_logger_tracef_sig_49;
			return 0;
		case 50:
			entry->call = gobench_logger_tracef_sig_50;
			return 0;
		case 51:
			entry->call = gobench_logger_tracef_sig_51;
			return 0;
		case 52:
			entry->call = gobench_logger_tracef_sig_52;
			return 0;
		case 53:
			entry->call = gobench_logger_tracef_sig_53;
			return 0;
		case 54:
			entry->call = gobench_logger_tracef_sig_54;
			return 0;
		case 55:
			entry->call = gobench_logger_tracef_sig_55;
			return 0;
		case 56:
			entry->call = gobench_logger_tracef_sig_56;
			return 0;
		case 57:
			entry->call = gobench_logger_tracef_sig_57;
			return 0;
		case 58:
			entry->call = gobench_logger_tracef_sig_58;
			return 0;
		case 59:
			entry->call = gobench_logger_tracef_sig_59;
			return 0;
		case 60:
			entry->call = gobench_logger_tracef_sig_60;
			return 0;
		case 61:
			entry->call = gobench_logger_tracef_sig_61;
			return 0;
		case 62:
			entry->call = gobench_logger_tracef_sig_62;
			return 0;
		case 63:
			entry->call = gobench_logger_tracef_sig_63;
			return 0;
		case 64:
			entry->call = gobench_logger_tracef_sig_64;
			return 0;
		case 65:
			entry->call = gobench_logger_tracef_sig_65;
			return 0;
		case 66:
			entry->call = gobench_logger_tracef_sig_66;
			return 0;
		case 67:
			entry->call = gobench_logger_tracef_sig_67;
			return 0;
		case 68:
			entry->call = gobench_logger_tracef_sig_68;
			return 0;
		case 69:
			entry->call = gobench_logger_tracef_sig_69;
			return 0;
		case 70:
			entry->call = gobench_logger_tracef_sig_70;
			return 0;
		case 71:
			entry->call = gobench_logger_tracef_sig_71;
			return 0;
		case 72:
			entry->call = gobench_logger_tracef_sig_72;
			return 0;
		case 73:
			entry->call = gobench_logger_tracef_sig_73;
			return 0;
		case 74:
			entry->call = gobench_logger_tracef_sig_74;
			return 0;
		case 75:
			entry->call = gobench_logger_tracef_sig_75;
			return 0;
		case 76:
			entry->call = gobench_logger_tracef_sig_76;
			return 0;
		case 77:
			entry->call = gobench_logger_tracef_sig_77;
			return 0;
		case 78:
			entry->call = gobench_logger_tracef_sig_78;
			return 0;
		case 79:
			entry->call = gobench_logger_tracef_sig_79;
			return 0;
		case 80:
			entry->call = gobench_logger_tracef_sig_80;
			return 0;
		case 81:
			entry->call = gobench_logger_tracef_sig_81;
			return 0;
		case 82:
			entry->call = gobench_logger_tracef_sig_82;
			return 0;
		default:
			entry->call = NULL;
			return -1;
		}
	case PSLOG_LEVEL_DEBUG:
		switch (signature_id) {
		case 0:
			entry->call = gobench_logger_debugf_sig_0;
			return 0;
		case 1:
			entry->call = gobench_logger_debugf_sig_1;
			return 0;
		case 2:
			entry->call = gobench_logger_debugf_sig_2;
			return 0;
		case 3:
			entry->call = gobench_logger_debugf_sig_3;
			return 0;
		case 4:
			entry->call = gobench_logger_debugf_sig_4;
			return 0;
		case 5:
			entry->call = gobench_logger_debugf_sig_5;
			return 0;
		case 6:
			entry->call = gobench_logger_debugf_sig_6;
			return 0;
		case 7:
			entry->call = gobench_logger_debugf_sig_7;
			return 0;
		case 8:
			entry->call = gobench_logger_debugf_sig_8;
			return 0;
		case 9:
			entry->call = gobench_logger_debugf_sig_9;
			return 0;
		case 10:
			entry->call = gobench_logger_debugf_sig_10;
			return 0;
		case 11:
			entry->call = gobench_logger_debugf_sig_11;
			return 0;
		case 12:
			entry->call = gobench_logger_debugf_sig_12;
			return 0;
		case 13:
			entry->call = gobench_logger_debugf_sig_13;
			return 0;
		case 14:
			entry->call = gobench_logger_debugf_sig_14;
			return 0;
		case 15:
			entry->call = gobench_logger_debugf_sig_15;
			return 0;
		case 16:
			entry->call = gobench_logger_debugf_sig_16;
			return 0;
		case 17:
			entry->call = gobench_logger_debugf_sig_17;
			return 0;
		case 18:
			entry->call = gobench_logger_debugf_sig_18;
			return 0;
		case 19:
			entry->call = gobench_logger_debugf_sig_19;
			return 0;
		case 20:
			entry->call = gobench_logger_debugf_sig_20;
			return 0;
		case 21:
			entry->call = gobench_logger_debugf_sig_21;
			return 0;
		case 22:
			entry->call = gobench_logger_debugf_sig_22;
			return 0;
		case 23:
			entry->call = gobench_logger_debugf_sig_23;
			return 0;
		case 24:
			entry->call = gobench_logger_debugf_sig_24;
			return 0;
		case 25:
			entry->call = gobench_logger_debugf_sig_25;
			return 0;
		case 26:
			entry->call = gobench_logger_debugf_sig_26;
			return 0;
		case 27:
			entry->call = gobench_logger_debugf_sig_27;
			return 0;
		case 28:
			entry->call = gobench_logger_debugf_sig_28;
			return 0;
		case 29:
			entry->call = gobench_logger_debugf_sig_29;
			return 0;
		case 30:
			entry->call = gobench_logger_debugf_sig_30;
			return 0;
		case 31:
			entry->call = gobench_logger_debugf_sig_31;
			return 0;
		case 32:
			entry->call = gobench_logger_debugf_sig_32;
			return 0;
		case 33:
			entry->call = gobench_logger_debugf_sig_33;
			return 0;
		case 34:
			entry->call = gobench_logger_debugf_sig_34;
			return 0;
		case 35:
			entry->call = gobench_logger_debugf_sig_35;
			return 0;
		case 36:
			entry->call = gobench_logger_debugf_sig_36;
			return 0;
		case 37:
			entry->call = gobench_logger_debugf_sig_37;
			return 0;
		case 38:
			entry->call = gobench_logger_debugf_sig_38;
			return 0;
		case 39:
			entry->call = gobench_logger_debugf_sig_39;
			return 0;
		case 40:
			entry->call = gobench_logger_debugf_sig_40;
			return 0;
		case 41:
			entry->call = gobench_logger_debugf_sig_41;
			return 0;
		case 42:
			entry->call = gobench_logger_debugf_sig_42;
			return 0;
		case 43:
			entry->call = gobench_logger_debugf_sig_43;
			return 0;
		case 44:
			entry->call = gobench_logger_debugf_sig_44;
			return 0;
		case 45:
			entry->call = gobench_logger_debugf_sig_45;
			return 0;
		case 46:
			entry->call = gobench_logger_debugf_sig_46;
			return 0;
		case 47:
			entry->call = gobench_logger_debugf_sig_47;
			return 0;
		case 48:
			entry->call = gobench_logger_debugf_sig_48;
			return 0;
		case 49:
			entry->call = gobench_logger_debugf_sig_49;
			return 0;
		case 50:
			entry->call = gobench_logger_debugf_sig_50;
			return 0;
		case 51:
			entry->call = gobench_logger_debugf_sig_51;
			return 0;
		case 52:
			entry->call = gobench_logger_debugf_sig_52;
			return 0;
		case 53:
			entry->call = gobench_logger_debugf_sig_53;
			return 0;
		case 54:
			entry->call = gobench_logger_debugf_sig_54;
			return 0;
		case 55:
			entry->call = gobench_logger_debugf_sig_55;
			return 0;
		case 56:
			entry->call = gobench_logger_debugf_sig_56;
			return 0;
		case 57:
			entry->call = gobench_logger_debugf_sig_57;
			return 0;
		case 58:
			entry->call = gobench_logger_debugf_sig_58;
			return 0;
		case 59:
			entry->call = gobench_logger_debugf_sig_59;
			return 0;
		case 60:
			entry->call = gobench_logger_debugf_sig_60;
			return 0;
		case 61:
			entry->call = gobench_logger_debugf_sig_61;
			return 0;
		case 62:
			entry->call = gobench_logger_debugf_sig_62;
			return 0;
		case 63:
			entry->call = gobench_logger_debugf_sig_63;
			return 0;
		case 64:
			entry->call = gobench_logger_debugf_sig_64;
			return 0;
		case 65:
			entry->call = gobench_logger_debugf_sig_65;
			return 0;
		case 66:
			entry->call = gobench_logger_debugf_sig_66;
			return 0;
		case 67:
			entry->call = gobench_logger_debugf_sig_67;
			return 0;
		case 68:
			entry->call = gobench_logger_debugf_sig_68;
			return 0;
		case 69:
			entry->call = gobench_logger_debugf_sig_69;
			return 0;
		case 70:
			entry->call = gobench_logger_debugf_sig_70;
			return 0;
		case 71:
			entry->call = gobench_logger_debugf_sig_71;
			return 0;
		case 72:
			entry->call = gobench_logger_debugf_sig_72;
			return 0;
		case 73:
			entry->call = gobench_logger_debugf_sig_73;
			return 0;
		case 74:
			entry->call = gobench_logger_debugf_sig_74;
			return 0;
		case 75:
			entry->call = gobench_logger_debugf_sig_75;
			return 0;
		case 76:
			entry->call = gobench_logger_debugf_sig_76;
			return 0;
		case 77:
			entry->call = gobench_logger_debugf_sig_77;
			return 0;
		case 78:
			entry->call = gobench_logger_debugf_sig_78;
			return 0;
		case 79:
			entry->call = gobench_logger_debugf_sig_79;
			return 0;
		case 80:
			entry->call = gobench_logger_debugf_sig_80;
			return 0;
		case 81:
			entry->call = gobench_logger_debugf_sig_81;
			return 0;
		case 82:
			entry->call = gobench_logger_debugf_sig_82;
			return 0;
		default:
			entry->call = NULL;
			return -1;
		}
	case PSLOG_LEVEL_INFO:
		switch (signature_id) {
		case 0:
			entry->call = gobench_logger_infof_sig_0;
			return 0;
		case 1:
			entry->call = gobench_logger_infof_sig_1;
			return 0;
		case 2:
			entry->call = gobench_logger_infof_sig_2;
			return 0;
		case 3:
			entry->call = gobench_logger_infof_sig_3;
			return 0;
		case 4:
			entry->call = gobench_logger_infof_sig_4;
			return 0;
		case 5:
			entry->call = gobench_logger_infof_sig_5;
			return 0;
		case 6:
			entry->call = gobench_logger_infof_sig_6;
			return 0;
		case 7:
			entry->call = gobench_logger_infof_sig_7;
			return 0;
		case 8:
			entry->call = gobench_logger_infof_sig_8;
			return 0;
		case 9:
			entry->call = gobench_logger_infof_sig_9;
			return 0;
		case 10:
			entry->call = gobench_logger_infof_sig_10;
			return 0;
		case 11:
			entry->call = gobench_logger_infof_sig_11;
			return 0;
		case 12:
			entry->call = gobench_logger_infof_sig_12;
			return 0;
		case 13:
			entry->call = gobench_logger_infof_sig_13;
			return 0;
		case 14:
			entry->call = gobench_logger_infof_sig_14;
			return 0;
		case 15:
			entry->call = gobench_logger_infof_sig_15;
			return 0;
		case 16:
			entry->call = gobench_logger_infof_sig_16;
			return 0;
		case 17:
			entry->call = gobench_logger_infof_sig_17;
			return 0;
		case 18:
			entry->call = gobench_logger_infof_sig_18;
			return 0;
		case 19:
			entry->call = gobench_logger_infof_sig_19;
			return 0;
		case 20:
			entry->call = gobench_logger_infof_sig_20;
			return 0;
		case 21:
			entry->call = gobench_logger_infof_sig_21;
			return 0;
		case 22:
			entry->call = gobench_logger_infof_sig_22;
			return 0;
		case 23:
			entry->call = gobench_logger_infof_sig_23;
			return 0;
		case 24:
			entry->call = gobench_logger_infof_sig_24;
			return 0;
		case 25:
			entry->call = gobench_logger_infof_sig_25;
			return 0;
		case 26:
			entry->call = gobench_logger_infof_sig_26;
			return 0;
		case 27:
			entry->call = gobench_logger_infof_sig_27;
			return 0;
		case 28:
			entry->call = gobench_logger_infof_sig_28;
			return 0;
		case 29:
			entry->call = gobench_logger_infof_sig_29;
			return 0;
		case 30:
			entry->call = gobench_logger_infof_sig_30;
			return 0;
		case 31:
			entry->call = gobench_logger_infof_sig_31;
			return 0;
		case 32:
			entry->call = gobench_logger_infof_sig_32;
			return 0;
		case 33:
			entry->call = gobench_logger_infof_sig_33;
			return 0;
		case 34:
			entry->call = gobench_logger_infof_sig_34;
			return 0;
		case 35:
			entry->call = gobench_logger_infof_sig_35;
			return 0;
		case 36:
			entry->call = gobench_logger_infof_sig_36;
			return 0;
		case 37:
			entry->call = gobench_logger_infof_sig_37;
			return 0;
		case 38:
			entry->call = gobench_logger_infof_sig_38;
			return 0;
		case 39:
			entry->call = gobench_logger_infof_sig_39;
			return 0;
		case 40:
			entry->call = gobench_logger_infof_sig_40;
			return 0;
		case 41:
			entry->call = gobench_logger_infof_sig_41;
			return 0;
		case 42:
			entry->call = gobench_logger_infof_sig_42;
			return 0;
		case 43:
			entry->call = gobench_logger_infof_sig_43;
			return 0;
		case 44:
			entry->call = gobench_logger_infof_sig_44;
			return 0;
		case 45:
			entry->call = gobench_logger_infof_sig_45;
			return 0;
		case 46:
			entry->call = gobench_logger_infof_sig_46;
			return 0;
		case 47:
			entry->call = gobench_logger_infof_sig_47;
			return 0;
		case 48:
			entry->call = gobench_logger_infof_sig_48;
			return 0;
		case 49:
			entry->call = gobench_logger_infof_sig_49;
			return 0;
		case 50:
			entry->call = gobench_logger_infof_sig_50;
			return 0;
		case 51:
			entry->call = gobench_logger_infof_sig_51;
			return 0;
		case 52:
			entry->call = gobench_logger_infof_sig_52;
			return 0;
		case 53:
			entry->call = gobench_logger_infof_sig_53;
			return 0;
		case 54:
			entry->call = gobench_logger_infof_sig_54;
			return 0;
		case 55:
			entry->call = gobench_logger_infof_sig_55;
			return 0;
		case 56:
			entry->call = gobench_logger_infof_sig_56;
			return 0;
		case 57:
			entry->call = gobench_logger_infof_sig_57;
			return 0;
		case 58:
			entry->call = gobench_logger_infof_sig_58;
			return 0;
		case 59:
			entry->call = gobench_logger_infof_sig_59;
			return 0;
		case 60:
			entry->call = gobench_logger_infof_sig_60;
			return 0;
		case 61:
			entry->call = gobench_logger_infof_sig_61;
			return 0;
		case 62:
			entry->call = gobench_logger_infof_sig_62;
			return 0;
		case 63:
			entry->call = gobench_logger_infof_sig_63;
			return 0;
		case 64:
			entry->call = gobench_logger_infof_sig_64;
			return 0;
		case 65:
			entry->call = gobench_logger_infof_sig_65;
			return 0;
		case 66:
			entry->call = gobench_logger_infof_sig_66;
			return 0;
		case 67:
			entry->call = gobench_logger_infof_sig_67;
			return 0;
		case 68:
			entry->call = gobench_logger_infof_sig_68;
			return 0;
		case 69:
			entry->call = gobench_logger_infof_sig_69;
			return 0;
		case 70:
			entry->call = gobench_logger_infof_sig_70;
			return 0;
		case 71:
			entry->call = gobench_logger_infof_sig_71;
			return 0;
		case 72:
			entry->call = gobench_logger_infof_sig_72;
			return 0;
		case 73:
			entry->call = gobench_logger_infof_sig_73;
			return 0;
		case 74:
			entry->call = gobench_logger_infof_sig_74;
			return 0;
		case 75:
			entry->call = gobench_logger_infof_sig_75;
			return 0;
		case 76:
			entry->call = gobench_logger_infof_sig_76;
			return 0;
		case 77:
			entry->call = gobench_logger_infof_sig_77;
			return 0;
		case 78:
			entry->call = gobench_logger_infof_sig_78;
			return 0;
		case 79:
			entry->call = gobench_logger_infof_sig_79;
			return 0;
		case 80:
			entry->call = gobench_logger_infof_sig_80;
			return 0;
		case 81:
			entry->call = gobench_logger_infof_sig_81;
			return 0;
		case 82:
			entry->call = gobench_logger_infof_sig_82;
			return 0;
		default:
			entry->call = NULL;
			return -1;
		}
	case PSLOG_LEVEL_WARN:
		switch (signature_id) {
		case 0:
			entry->call = gobench_logger_warnf_sig_0;
			return 0;
		case 1:
			entry->call = gobench_logger_warnf_sig_1;
			return 0;
		case 2:
			entry->call = gobench_logger_warnf_sig_2;
			return 0;
		case 3:
			entry->call = gobench_logger_warnf_sig_3;
			return 0;
		case 4:
			entry->call = gobench_logger_warnf_sig_4;
			return 0;
		case 5:
			entry->call = gobench_logger_warnf_sig_5;
			return 0;
		case 6:
			entry->call = gobench_logger_warnf_sig_6;
			return 0;
		case 7:
			entry->call = gobench_logger_warnf_sig_7;
			return 0;
		case 8:
			entry->call = gobench_logger_warnf_sig_8;
			return 0;
		case 9:
			entry->call = gobench_logger_warnf_sig_9;
			return 0;
		case 10:
			entry->call = gobench_logger_warnf_sig_10;
			return 0;
		case 11:
			entry->call = gobench_logger_warnf_sig_11;
			return 0;
		case 12:
			entry->call = gobench_logger_warnf_sig_12;
			return 0;
		case 13:
			entry->call = gobench_logger_warnf_sig_13;
			return 0;
		case 14:
			entry->call = gobench_logger_warnf_sig_14;
			return 0;
		case 15:
			entry->call = gobench_logger_warnf_sig_15;
			return 0;
		case 16:
			entry->call = gobench_logger_warnf_sig_16;
			return 0;
		case 17:
			entry->call = gobench_logger_warnf_sig_17;
			return 0;
		case 18:
			entry->call = gobench_logger_warnf_sig_18;
			return 0;
		case 19:
			entry->call = gobench_logger_warnf_sig_19;
			return 0;
		case 20:
			entry->call = gobench_logger_warnf_sig_20;
			return 0;
		case 21:
			entry->call = gobench_logger_warnf_sig_21;
			return 0;
		case 22:
			entry->call = gobench_logger_warnf_sig_22;
			return 0;
		case 23:
			entry->call = gobench_logger_warnf_sig_23;
			return 0;
		case 24:
			entry->call = gobench_logger_warnf_sig_24;
			return 0;
		case 25:
			entry->call = gobench_logger_warnf_sig_25;
			return 0;
		case 26:
			entry->call = gobench_logger_warnf_sig_26;
			return 0;
		case 27:
			entry->call = gobench_logger_warnf_sig_27;
			return 0;
		case 28:
			entry->call = gobench_logger_warnf_sig_28;
			return 0;
		case 29:
			entry->call = gobench_logger_warnf_sig_29;
			return 0;
		case 30:
			entry->call = gobench_logger_warnf_sig_30;
			return 0;
		case 31:
			entry->call = gobench_logger_warnf_sig_31;
			return 0;
		case 32:
			entry->call = gobench_logger_warnf_sig_32;
			return 0;
		case 33:
			entry->call = gobench_logger_warnf_sig_33;
			return 0;
		case 34:
			entry->call = gobench_logger_warnf_sig_34;
			return 0;
		case 35:
			entry->call = gobench_logger_warnf_sig_35;
			return 0;
		case 36:
			entry->call = gobench_logger_warnf_sig_36;
			return 0;
		case 37:
			entry->call = gobench_logger_warnf_sig_37;
			return 0;
		case 38:
			entry->call = gobench_logger_warnf_sig_38;
			return 0;
		case 39:
			entry->call = gobench_logger_warnf_sig_39;
			return 0;
		case 40:
			entry->call = gobench_logger_warnf_sig_40;
			return 0;
		case 41:
			entry->call = gobench_logger_warnf_sig_41;
			return 0;
		case 42:
			entry->call = gobench_logger_warnf_sig_42;
			return 0;
		case 43:
			entry->call = gobench_logger_warnf_sig_43;
			return 0;
		case 44:
			entry->call = gobench_logger_warnf_sig_44;
			return 0;
		case 45:
			entry->call = gobench_logger_warnf_sig_45;
			return 0;
		case 46:
			entry->call = gobench_logger_warnf_sig_46;
			return 0;
		case 47:
			entry->call = gobench_logger_warnf_sig_47;
			return 0;
		case 48:
			entry->call = gobench_logger_warnf_sig_48;
			return 0;
		case 49:
			entry->call = gobench_logger_warnf_sig_49;
			return 0;
		case 50:
			entry->call = gobench_logger_warnf_sig_50;
			return 0;
		case 51:
			entry->call = gobench_logger_warnf_sig_51;
			return 0;
		case 52:
			entry->call = gobench_logger_warnf_sig_52;
			return 0;
		case 53:
			entry->call = gobench_logger_warnf_sig_53;
			return 0;
		case 54:
			entry->call = gobench_logger_warnf_sig_54;
			return 0;
		case 55:
			entry->call = gobench_logger_warnf_sig_55;
			return 0;
		case 56:
			entry->call = gobench_logger_warnf_sig_56;
			return 0;
		case 57:
			entry->call = gobench_logger_warnf_sig_57;
			return 0;
		case 58:
			entry->call = gobench_logger_warnf_sig_58;
			return 0;
		case 59:
			entry->call = gobench_logger_warnf_sig_59;
			return 0;
		case 60:
			entry->call = gobench_logger_warnf_sig_60;
			return 0;
		case 61:
			entry->call = gobench_logger_warnf_sig_61;
			return 0;
		case 62:
			entry->call = gobench_logger_warnf_sig_62;
			return 0;
		case 63:
			entry->call = gobench_logger_warnf_sig_63;
			return 0;
		case 64:
			entry->call = gobench_logger_warnf_sig_64;
			return 0;
		case 65:
			entry->call = gobench_logger_warnf_sig_65;
			return 0;
		case 66:
			entry->call = gobench_logger_warnf_sig_66;
			return 0;
		case 67:
			entry->call = gobench_logger_warnf_sig_67;
			return 0;
		case 68:
			entry->call = gobench_logger_warnf_sig_68;
			return 0;
		case 69:
			entry->call = gobench_logger_warnf_sig_69;
			return 0;
		case 70:
			entry->call = gobench_logger_warnf_sig_70;
			return 0;
		case 71:
			entry->call = gobench_logger_warnf_sig_71;
			return 0;
		case 72:
			entry->call = gobench_logger_warnf_sig_72;
			return 0;
		case 73:
			entry->call = gobench_logger_warnf_sig_73;
			return 0;
		case 74:
			entry->call = gobench_logger_warnf_sig_74;
			return 0;
		case 75:
			entry->call = gobench_logger_warnf_sig_75;
			return 0;
		case 76:
			entry->call = gobench_logger_warnf_sig_76;
			return 0;
		case 77:
			entry->call = gobench_logger_warnf_sig_77;
			return 0;
		case 78:
			entry->call = gobench_logger_warnf_sig_78;
			return 0;
		case 79:
			entry->call = gobench_logger_warnf_sig_79;
			return 0;
		case 80:
			entry->call = gobench_logger_warnf_sig_80;
			return 0;
		case 81:
			entry->call = gobench_logger_warnf_sig_81;
			return 0;
		case 82:
			entry->call = gobench_logger_warnf_sig_82;
			return 0;
		default:
			entry->call = NULL;
			return -1;
		}
	case PSLOG_LEVEL_ERROR:
		switch (signature_id) {
		case 0:
			entry->call = gobench_logger_errorf_sig_0;
			return 0;
		case 1:
			entry->call = gobench_logger_errorf_sig_1;
			return 0;
		case 2:
			entry->call = gobench_logger_errorf_sig_2;
			return 0;
		case 3:
			entry->call = gobench_logger_errorf_sig_3;
			return 0;
		case 4:
			entry->call = gobench_logger_errorf_sig_4;
			return 0;
		case 5:
			entry->call = gobench_logger_errorf_sig_5;
			return 0;
		case 6:
			entry->call = gobench_logger_errorf_sig_6;
			return 0;
		case 7:
			entry->call = gobench_logger_errorf_sig_7;
			return 0;
		case 8:
			entry->call = gobench_logger_errorf_sig_8;
			return 0;
		case 9:
			entry->call = gobench_logger_errorf_sig_9;
			return 0;
		case 10:
			entry->call = gobench_logger_errorf_sig_10;
			return 0;
		case 11:
			entry->call = gobench_logger_errorf_sig_11;
			return 0;
		case 12:
			entry->call = gobench_logger_errorf_sig_12;
			return 0;
		case 13:
			entry->call = gobench_logger_errorf_sig_13;
			return 0;
		case 14:
			entry->call = gobench_logger_errorf_sig_14;
			return 0;
		case 15:
			entry->call = gobench_logger_errorf_sig_15;
			return 0;
		case 16:
			entry->call = gobench_logger_errorf_sig_16;
			return 0;
		case 17:
			entry->call = gobench_logger_errorf_sig_17;
			return 0;
		case 18:
			entry->call = gobench_logger_errorf_sig_18;
			return 0;
		case 19:
			entry->call = gobench_logger_errorf_sig_19;
			return 0;
		case 20:
			entry->call = gobench_logger_errorf_sig_20;
			return 0;
		case 21:
			entry->call = gobench_logger_errorf_sig_21;
			return 0;
		case 22:
			entry->call = gobench_logger_errorf_sig_22;
			return 0;
		case 23:
			entry->call = gobench_logger_errorf_sig_23;
			return 0;
		case 24:
			entry->call = gobench_logger_errorf_sig_24;
			return 0;
		case 25:
			entry->call = gobench_logger_errorf_sig_25;
			return 0;
		case 26:
			entry->call = gobench_logger_errorf_sig_26;
			return 0;
		case 27:
			entry->call = gobench_logger_errorf_sig_27;
			return 0;
		case 28:
			entry->call = gobench_logger_errorf_sig_28;
			return 0;
		case 29:
			entry->call = gobench_logger_errorf_sig_29;
			return 0;
		case 30:
			entry->call = gobench_logger_errorf_sig_30;
			return 0;
		case 31:
			entry->call = gobench_logger_errorf_sig_31;
			return 0;
		case 32:
			entry->call = gobench_logger_errorf_sig_32;
			return 0;
		case 33:
			entry->call = gobench_logger_errorf_sig_33;
			return 0;
		case 34:
			entry->call = gobench_logger_errorf_sig_34;
			return 0;
		case 35:
			entry->call = gobench_logger_errorf_sig_35;
			return 0;
		case 36:
			entry->call = gobench_logger_errorf_sig_36;
			return 0;
		case 37:
			entry->call = gobench_logger_errorf_sig_37;
			return 0;
		case 38:
			entry->call = gobench_logger_errorf_sig_38;
			return 0;
		case 39:
			entry->call = gobench_logger_errorf_sig_39;
			return 0;
		case 40:
			entry->call = gobench_logger_errorf_sig_40;
			return 0;
		case 41:
			entry->call = gobench_logger_errorf_sig_41;
			return 0;
		case 42:
			entry->call = gobench_logger_errorf_sig_42;
			return 0;
		case 43:
			entry->call = gobench_logger_errorf_sig_43;
			return 0;
		case 44:
			entry->call = gobench_logger_errorf_sig_44;
			return 0;
		case 45:
			entry->call = gobench_logger_errorf_sig_45;
			return 0;
		case 46:
			entry->call = gobench_logger_errorf_sig_46;
			return 0;
		case 47:
			entry->call = gobench_logger_errorf_sig_47;
			return 0;
		case 48:
			entry->call = gobench_logger_errorf_sig_48;
			return 0;
		case 49:
			entry->call = gobench_logger_errorf_sig_49;
			return 0;
		case 50:
			entry->call = gobench_logger_errorf_sig_50;
			return 0;
		case 51:
			entry->call = gobench_logger_errorf_sig_51;
			return 0;
		case 52:
			entry->call = gobench_logger_errorf_sig_52;
			return 0;
		case 53:
			entry->call = gobench_logger_errorf_sig_53;
			return 0;
		case 54:
			entry->call = gobench_logger_errorf_sig_54;
			return 0;
		case 55:
			entry->call = gobench_logger_errorf_sig_55;
			return 0;
		case 56:
			entry->call = gobench_logger_errorf_sig_56;
			return 0;
		case 57:
			entry->call = gobench_logger_errorf_sig_57;
			return 0;
		case 58:
			entry->call = gobench_logger_errorf_sig_58;
			return 0;
		case 59:
			entry->call = gobench_logger_errorf_sig_59;
			return 0;
		case 60:
			entry->call = gobench_logger_errorf_sig_60;
			return 0;
		case 61:
			entry->call = gobench_logger_errorf_sig_61;
			return 0;
		case 62:
			entry->call = gobench_logger_errorf_sig_62;
			return 0;
		case 63:
			entry->call = gobench_logger_errorf_sig_63;
			return 0;
		case 64:
			entry->call = gobench_logger_errorf_sig_64;
			return 0;
		case 65:
			entry->call = gobench_logger_errorf_sig_65;
			return 0;
		case 66:
			entry->call = gobench_logger_errorf_sig_66;
			return 0;
		case 67:
			entry->call = gobench_logger_errorf_sig_67;
			return 0;
		case 68:
			entry->call = gobench_logger_errorf_sig_68;
			return 0;
		case 69:
			entry->call = gobench_logger_errorf_sig_69;
			return 0;
		case 70:
			entry->call = gobench_logger_errorf_sig_70;
			return 0;
		case 71:
			entry->call = gobench_logger_errorf_sig_71;
			return 0;
		case 72:
			entry->call = gobench_logger_errorf_sig_72;
			return 0;
		case 73:
			entry->call = gobench_logger_errorf_sig_73;
			return 0;
		case 74:
			entry->call = gobench_logger_errorf_sig_74;
			return 0;
		case 75:
			entry->call = gobench_logger_errorf_sig_75;
			return 0;
		case 76:
			entry->call = gobench_logger_errorf_sig_76;
			return 0;
		case 77:
			entry->call = gobench_logger_errorf_sig_77;
			return 0;
		case 78:
			entry->call = gobench_logger_errorf_sig_78;
			return 0;
		case 79:
			entry->call = gobench_logger_errorf_sig_79;
			return 0;
		case 80:
			entry->call = gobench_logger_errorf_sig_80;
			return 0;
		case 81:
			entry->call = gobench_logger_errorf_sig_81;
			return 0;
		case 82:
			entry->call = gobench_logger_errorf_sig_82;
			return 0;
		default:
			entry->call = NULL;
			return -1;
		}
	default:
		switch (signature_id) {
		case 0:
			entry->call = gobench_logger_logf_sig_0;
			return 0;
		case 1:
			entry->call = gobench_logger_logf_sig_1;
			return 0;
		case 2:
			entry->call = gobench_logger_logf_sig_2;
			return 0;
		case 3:
			entry->call = gobench_logger_logf_sig_3;
			return 0;
		case 4:
			entry->call = gobench_logger_logf_sig_4;
			return 0;
		case 5:
			entry->call = gobench_logger_logf_sig_5;
			return 0;
		case 6:
			entry->call = gobench_logger_logf_sig_6;
			return 0;
		case 7:
			entry->call = gobench_logger_logf_sig_7;
			return 0;
		case 8:
			entry->call = gobench_logger_logf_sig_8;
			return 0;
		case 9:
			entry->call = gobench_logger_logf_sig_9;
			return 0;
		case 10:
			entry->call = gobench_logger_logf_sig_10;
			return 0;
		case 11:
			entry->call = gobench_logger_logf_sig_11;
			return 0;
		case 12:
			entry->call = gobench_logger_logf_sig_12;
			return 0;
		case 13:
			entry->call = gobench_logger_logf_sig_13;
			return 0;
		case 14:
			entry->call = gobench_logger_logf_sig_14;
			return 0;
		case 15:
			entry->call = gobench_logger_logf_sig_15;
			return 0;
		case 16:
			entry->call = gobench_logger_logf_sig_16;
			return 0;
		case 17:
			entry->call = gobench_logger_logf_sig_17;
			return 0;
		case 18:
			entry->call = gobench_logger_logf_sig_18;
			return 0;
		case 19:
			entry->call = gobench_logger_logf_sig_19;
			return 0;
		case 20:
			entry->call = gobench_logger_logf_sig_20;
			return 0;
		case 21:
			entry->call = gobench_logger_logf_sig_21;
			return 0;
		case 22:
			entry->call = gobench_logger_logf_sig_22;
			return 0;
		case 23:
			entry->call = gobench_logger_logf_sig_23;
			return 0;
		case 24:
			entry->call = gobench_logger_logf_sig_24;
			return 0;
		case 25:
			entry->call = gobench_logger_logf_sig_25;
			return 0;
		case 26:
			entry->call = gobench_logger_logf_sig_26;
			return 0;
		case 27:
			entry->call = gobench_logger_logf_sig_27;
			return 0;
		case 28:
			entry->call = gobench_logger_logf_sig_28;
			return 0;
		case 29:
			entry->call = gobench_logger_logf_sig_29;
			return 0;
		case 30:
			entry->call = gobench_logger_logf_sig_30;
			return 0;
		case 31:
			entry->call = gobench_logger_logf_sig_31;
			return 0;
		case 32:
			entry->call = gobench_logger_logf_sig_32;
			return 0;
		case 33:
			entry->call = gobench_logger_logf_sig_33;
			return 0;
		case 34:
			entry->call = gobench_logger_logf_sig_34;
			return 0;
		case 35:
			entry->call = gobench_logger_logf_sig_35;
			return 0;
		case 36:
			entry->call = gobench_logger_logf_sig_36;
			return 0;
		case 37:
			entry->call = gobench_logger_logf_sig_37;
			return 0;
		case 38:
			entry->call = gobench_logger_logf_sig_38;
			return 0;
		case 39:
			entry->call = gobench_logger_logf_sig_39;
			return 0;
		case 40:
			entry->call = gobench_logger_logf_sig_40;
			return 0;
		case 41:
			entry->call = gobench_logger_logf_sig_41;
			return 0;
		case 42:
			entry->call = gobench_logger_logf_sig_42;
			return 0;
		case 43:
			entry->call = gobench_logger_logf_sig_43;
			return 0;
		case 44:
			entry->call = gobench_logger_logf_sig_44;
			return 0;
		case 45:
			entry->call = gobench_logger_logf_sig_45;
			return 0;
		case 46:
			entry->call = gobench_logger_logf_sig_46;
			return 0;
		case 47:
			entry->call = gobench_logger_logf_sig_47;
			return 0;
		case 48:
			entry->call = gobench_logger_logf_sig_48;
			return 0;
		case 49:
			entry->call = gobench_logger_logf_sig_49;
			return 0;
		case 50:
			entry->call = gobench_logger_logf_sig_50;
			return 0;
		case 51:
			entry->call = gobench_logger_logf_sig_51;
			return 0;
		case 52:
			entry->call = gobench_logger_logf_sig_52;
			return 0;
		case 53:
			entry->call = gobench_logger_logf_sig_53;
			return 0;
		case 54:
			entry->call = gobench_logger_logf_sig_54;
			return 0;
		case 55:
			entry->call = gobench_logger_logf_sig_55;
			return 0;
		case 56:
			entry->call = gobench_logger_logf_sig_56;
			return 0;
		case 57:
			entry->call = gobench_logger_logf_sig_57;
			return 0;
		case 58:
			entry->call = gobench_logger_logf_sig_58;
			return 0;
		case 59:
			entry->call = gobench_logger_logf_sig_59;
			return 0;
		case 60:
			entry->call = gobench_logger_logf_sig_60;
			return 0;
		case 61:
			entry->call = gobench_logger_logf_sig_61;
			return 0;
		case 62:
			entry->call = gobench_logger_logf_sig_62;
			return 0;
		case 63:
			entry->call = gobench_logger_logf_sig_63;
			return 0;
		case 64:
			entry->call = gobench_logger_logf_sig_64;
			return 0;
		case 65:
			entry->call = gobench_logger_logf_sig_65;
			return 0;
		case 66:
			entry->call = gobench_logger_logf_sig_66;
			return 0;
		case 67:
			entry->call = gobench_logger_logf_sig_67;
			return 0;
		case 68:
			entry->call = gobench_logger_logf_sig_68;
			return 0;
		case 69:
			entry->call = gobench_logger_logf_sig_69;
			return 0;
		case 70:
			entry->call = gobench_logger_logf_sig_70;
			return 0;
		case 71:
			entry->call = gobench_logger_logf_sig_71;
			return 0;
		case 72:
			entry->call = gobench_logger_logf_sig_72;
			return 0;
		case 73:
			entry->call = gobench_logger_logf_sig_73;
			return 0;
		case 74:
			entry->call = gobench_logger_logf_sig_74;
			return 0;
		case 75:
			entry->call = gobench_logger_logf_sig_75;
			return 0;
		case 76:
			entry->call = gobench_logger_logf_sig_76;
			return 0;
		case 77:
			entry->call = gobench_logger_logf_sig_77;
			return 0;
		case 78:
			entry->call = gobench_logger_logf_sig_78;
			return 0;
		case 79:
			entry->call = gobench_logger_logf_sig_79;
			return 0;
		case 80:
			entry->call = gobench_logger_logf_sig_80;
			return 0;
		case 81:
			entry->call = gobench_logger_logf_sig_81;
			return 0;
		case 82:
			entry->call = gobench_logger_logf_sig_82;
			return 0;
		default:
			entry->call = NULL;
			return -1;
		}
	}
}

int gobench_logger_log_kvfmt(void *logger_ptr,
				 const gobench_kvfmt_entry *entry) {
	gobench_logger *wrapper;

	wrapper = (gobench_logger *)logger_ptr;
	if (wrapper == NULL || wrapper->logger == NULL || entry == NULL || entry->call == NULL) {
		return -1;
	}
	entry->call(wrapper->logger, entry);
	return 0;
}

int gobench_logger_run_kvfmt(void *logger_ptr,
				 const gobench_kvfmt_entry *entries,
				 size_t entry_count,
				 const pslog_field *static_fields,
				 size_t static_count,
				 size_t iterations,
				 gobench_result *out) {
	unsigned long long start_ns;
	unsigned long long end_ns;
	gobench_logger *wrapper;
	pslog_logger *active_logger;
	pslog_logger *derived;
	size_t i;
	size_t index;

	wrapper = (gobench_logger *)logger_ptr;
	if (wrapper == NULL || wrapper->logger == NULL || entries == NULL || entry_count == 0u || out == NULL) {
		return -1;
	}

	active_logger = wrapper->logger;
	derived = NULL;
	if (static_count > 0u) {
		derived = wrapper->logger->with(wrapper->logger, static_fields, static_count);
		if (derived == NULL) {
			return -1;
		}
		active_logger = derived;
	}

	gobench_kvfmt_sink_reset(wrapper->sink);
	index = 0u;
	start_ns = gobench_kvfmt_monotonic_ns();
	for (i = 0u; i < iterations; ++i) {
		const gobench_kvfmt_entry *entry = entries + index;
		if (entry->call == NULL) {
			if (derived != NULL) {
				derived->destroy(derived);
			}
			return -1;
		}
		entry->call(active_logger, entry);
		++index;
		if (index == entry_count) {
			index = 0u;
		}
	}
	end_ns = gobench_kvfmt_monotonic_ns();

	out->elapsed_ns = end_ns - start_ns;
	out->bytes = wrapper->sink != NULL ? wrapper->sink->bytes : 0ull;
	out->writes = wrapper->sink != NULL ? wrapper->sink->writes : 0ull;
	out->ops = (unsigned long long)iterations;

	if (derived != NULL) {
		derived->destroy(derived);
	}
	return 0;
}
*/
import "C"

func init() {
	ckvfmtSignatureIDs["b,s"] = 0
	ckvfmtSignatureIDs["b,s,s,s"] = 1
	ckvfmtSignatureIDs["i,b,s,b,s"] = 2
	ckvfmtSignatureIDs["i,i,s,s,s"] = 3
	ckvfmtSignatureIDs["i,s,b,b,s,s,s,s,s,s,i"] = 4
	ckvfmtSignatureIDs["i,s,b,b,s,s,s,s,s,s,s,i"] = 5
	ckvfmtSignatureIDs["i,s,b,f"] = 6
	ckvfmtSignatureIDs["i,s,b,i,i,i,s"] = 7
	ckvfmtSignatureIDs["i,s,b,s,s,b,s"] = 8
	ckvfmtSignatureIDs["i,s,i,s"] = 9
	ckvfmtSignatureIDs["i,s,s,b,i,s,i,i,s"] = 10
	ckvfmtSignatureIDs["i,s,s,i"] = 11
	ckvfmtSignatureIDs["i,s,s,i,s,i"] = 12
	ckvfmtSignatureIDs["i,s,s,i,s,s"] = 13
	ckvfmtSignatureIDs["i,s,s,s,b,s,s,s"] = 14
	ckvfmtSignatureIDs["i,s,s,s,i"] = 15
	ckvfmtSignatureIDs["i,s,s,s,i,s,i"] = 16
	ckvfmtSignatureIDs["i,s,s,s,s"] = 17
	ckvfmtSignatureIDs["i,s,s,s,s,b,s,s,s"] = 18
	ckvfmtSignatureIDs["i,s,s,s,s,s,i,s,s,s,s"] = 19
	ckvfmtSignatureIDs["i,s,s,s,s,s,s"] = 20
	ckvfmtSignatureIDs["i,s,s,s,s,s,s,s"] = 21
	ckvfmtSignatureIDs["i,s,s,s,s,s,s,s,i"] = 22
	ckvfmtSignatureIDs["i,s,s,s,s,s,s,s,s"] = 23
	ckvfmtSignatureIDs["i,s,s,s,s,s,s,s,s,s"] = 24
	ckvfmtSignatureIDs["i,s,s,s,s,s,s,s,s,s,s"] = 25
	ckvfmtSignatureIDs["s"] = 26
	ckvfmtSignatureIDs["s,b,s,s"] = 27
	ckvfmtSignatureIDs["s,b,s,s,s,s"] = 28
	ckvfmtSignatureIDs["s,i,i,s,s,s,s"] = 29
	ckvfmtSignatureIDs["s,i,i,s,s,s,s,s,s,s,i"] = 30
	ckvfmtSignatureIDs["s,i,i,s,s,s,s,s,s,s,s,i"] = 31
	ckvfmtSignatureIDs["s,i,s,b,b,s,s,s,s,s,s,s,i"] = 32
	ckvfmtSignatureIDs["s,i,s,b,b,s,s,s,s,s,s,s,s,i"] = 33
	ckvfmtSignatureIDs["s,i,s,i,s,s,s,s,s,s,s,i"] = 34
	ckvfmtSignatureIDs["s,i,s,s,s,i"] = 35
	ckvfmtSignatureIDs["s,i,s,s,s,i,s,i"] = 36
	ckvfmtSignatureIDs["s,i,s,s,s,s,b,s,s,s"] = 37
	ckvfmtSignatureIDs["s,i,s,s,s,s,i"] = 38
	ckvfmtSignatureIDs["s,i,s,s,s,s,i,s,i"] = 39
	ckvfmtSignatureIDs["s,i,s,s,s,s,s"] = 40
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,b,s,s,s"] = 41
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,i,s,s,s"] = 42
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,i,s,s,s,s"] = 43
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,s"] = 44
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,s,s"] = 45
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,s,s,i"] = 46
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,s,s,s"] = 47
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,s,s,s,s"] = 48
	ckvfmtSignatureIDs["s,i,s,s,s,s,s,s,s,s,s,s,s"] = 49
	ckvfmtSignatureIDs["s,s"] = 50
	ckvfmtSignatureIDs["s,s,i,i,s,s,s,s,s,s,s,s,i"] = 51
	ckvfmtSignatureIDs["s,s,i,i,s,s,s,s,s,s,s,s,s,i"] = 52
	ckvfmtSignatureIDs["s,s,i,s,i,s,s,s,s,s,s,s,i"] = 53
	ckvfmtSignatureIDs["s,s,i,s,i,s,s,s,s,s,s,s,s,i"] = 54
	ckvfmtSignatureIDs["s,s,i,s,s,s,s,s"] = 55
	ckvfmtSignatureIDs["s,s,i,s,s,s,s,s,s"] = 56
	ckvfmtSignatureIDs["s,s,i,s,s,s,s,s,s,s,i,s,s,s"] = 57
	ckvfmtSignatureIDs["s,s,i,s,s,s,s,s,s,s,s"] = 58
	ckvfmtSignatureIDs["s,s,s"] = 59
	ckvfmtSignatureIDs["s,s,s,b,s,s"] = 60
	ckvfmtSignatureIDs["s,s,s,i,s"] = 61
	ckvfmtSignatureIDs["s,s,s,i,s,i,s,s,s,s,s,s,s,s,i"] = 62
	ckvfmtSignatureIDs["s,s,s,i,s,i,s,s,s,s,s,s,s,s,s,i"] = 63
	ckvfmtSignatureIDs["s,s,s,i,s,s,s,s,s,s"] = 64
	ckvfmtSignatureIDs["s,s,s,i,s,s,s,s,s,s,s,s,i"] = 65
	ckvfmtSignatureIDs["s,s,s,i,s,s,s,s,s,s,s,s,s,i"] = 66
	ckvfmtSignatureIDs["s,s,s,s"] = 67
	ckvfmtSignatureIDs["s,s,s,s,i,s,s,s,s,s,s,s,s,s,i"] = 68
	ckvfmtSignatureIDs["s,s,s,s,i,s,s,s,s,s,s,s,s,s,s,i"] = 69
	ckvfmtSignatureIDs["s,s,s,s,s"] = 70
	ckvfmtSignatureIDs["s,s,s,s,s,i,s"] = 71
	ckvfmtSignatureIDs["s,s,s,s,s,s"] = 72
	ckvfmtSignatureIDs["s,s,s,s,s,s,i,s"] = 73
	ckvfmtSignatureIDs["s,s,s,s,s,s,i,s,s,s,i"] = 74
	ckvfmtSignatureIDs["s,s,s,s,s,s,s"] = 75
	ckvfmtSignatureIDs["s,s,s,s,s,s,s,s"] = 76
	ckvfmtSignatureIDs["s,s,s,s,s,s,s,s,i,s"] = 77
	ckvfmtSignatureIDs["s,s,s,s,s,s,s,s,i,s,s,s,i"] = 78
	ckvfmtSignatureIDs["s,s,s,s,s,s,s,s,s"] = 79
	ckvfmtSignatureIDs["s,s,s,s,s,s,s,s,s,s"] = 80
	ckvfmtSignatureIDs["s,s,s,s,s,s,s,s,s,s,s"] = 81
	ckvfmtSignatureIDs["s,s,s,s,s,s,s,s,s,s,s,s"] = 82
}
