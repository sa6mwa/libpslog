#include "bench/bench_quill.h"

#include <stdio.h>
#include <string.h>
#include <time.h>

#include <chrono>
#include <cmath>
#include <memory>
#include <mutex>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

#include "quill/Backend.h"
#include "quill/Frontend.h"
#include "quill/LogFunctions.h"
#include "quill/Logger.h"
#include "quill/sinks/Sink.h"

namespace {

struct pslog_bench_quill_frontend_options {
  static constexpr quill::QueueType queue_type =
      quill::QueueType::UnboundedBlocking;
  static constexpr size_t initial_queue_capacity = 32u * 1024u * 1024u;
  static constexpr uint32_t blocking_queue_retry_interval_ns = 800u;
  static constexpr size_t unbounded_queue_max_capacity = 32u * 1024u * 1024u;
  static constexpr quill::HugePagesPolicy huge_pages_policy =
      quill::HugePagesPolicy::Never;
};

using pslog_bench_quill_frontend =
    quill::FrontendImpl<pslog_bench_quill_frontend_options>;
using pslog_bench_quill_native_logger =
    quill::LoggerImpl<pslog_bench_quill_frontend_options>;

struct pslog_bench_quill_logger {
  pslog_bench_quill_native_logger *logger;
};

static unsigned long long pslog_bench_quill_now_ns(void) {
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

static quill::LogLevel
pslog_bench_quill_map_level(pslog_bench_quill_level level) {
  switch (level) {
  case PSLOG_BENCH_QUILL_TRACE:
    return quill::LogLevel::TraceL3;
  case PSLOG_BENCH_QUILL_DEBUG:
    return quill::LogLevel::Debug;
  case PSLOG_BENCH_QUILL_INFO:
    return quill::LogLevel::Info;
  case PSLOG_BENCH_QUILL_WARN:
    return quill::LogLevel::Warning;
  case PSLOG_BENCH_QUILL_ERROR:
    return quill::LogLevel::Error;
  case PSLOG_BENCH_QUILL_FATAL:
  case PSLOG_BENCH_QUILL_PANIC:
    return quill::LogLevel::Critical;
  case PSLOG_BENCH_QUILL_NOLEVEL:
  default:
    return quill::LogLevel::Info;
  }
}

static const char *pslog_bench_quill_level_text(pslog_bench_quill_level level) {
  switch (level) {
  case PSLOG_BENCH_QUILL_TRACE:
    return "trace";
  case PSLOG_BENCH_QUILL_DEBUG:
    return "debug";
  case PSLOG_BENCH_QUILL_INFO:
    return "info";
  case PSLOG_BENCH_QUILL_WARN:
    return "warn";
  case PSLOG_BENCH_QUILL_ERROR:
    return "error";
  case PSLOG_BENCH_QUILL_FATAL:
    return "fatal";
  case PSLOG_BENCH_QUILL_PANIC:
    return "panic";
  case PSLOG_BENCH_QUILL_NOLEVEL:
    return "nolevel";
  default:
    return "info";
  }
}

template <typename... Args>
void pslog_bench_quill_log_runtime(pslog_bench_quill_logger *logger,
                                   pslog_bench_quill_level level,
                                   char const *fmt, Args &&...args) {
  if (logger == NULL || logger->logger == NULL || fmt == NULL) {
    return;
  }
  quill::log(logger->logger, "", pslog_bench_quill_map_level(level), fmt,
             quill::SourceLocation::current(), std::forward<Args>(args)...);
}

static char pslog_bench_quill_hex_digit(unsigned int value) {
  if (value < 10u) {
    return (char)('0' + value);
  }
  return (char)('a' + (value - 10u));
}

static unsigned int pslog_bench_quill_hex_value(char ch) {
  if (ch >= '0' && ch <= '9') {
    return (unsigned int)(ch - '0');
  }
  if (ch >= 'a' && ch <= 'f') {
    return (unsigned int)(10 + (ch - 'a'));
  }
  if (ch >= 'A' && ch <= 'F') {
    return (unsigned int)(10 + (ch - 'A'));
  }
  return 16u;
}

static void pslog_bench_quill_append_hex(std::string &out,
                                         std::string_view value) {
  for (char ch : value) {
    unsigned char byte = (unsigned char)ch;

    out.push_back(pslog_bench_quill_hex_digit((unsigned int)(byte >> 4u)));
    out.push_back(pslog_bench_quill_hex_digit((unsigned int)(byte & 0x0fu)));
  }
}

static int pslog_bench_quill_append_decoded_hex(std::string &out,
                                                std::string_view value) {
  size_t i;

  if ((value.size() % 2u) != 0u) {
    return -1;
  }
  for (i = 0u; i < value.size(); i += 2u) {
    unsigned int hi = pslog_bench_quill_hex_value(value[i]);
    unsigned int lo = pslog_bench_quill_hex_value(value[i + 1u]);

    if (hi > 15u || lo > 15u) {
      return -1;
    }
    out.push_back((char)((hi << 4u) | lo));
  }
  return 0;
}

static char
pslog_bench_quill_field_kind_code(pslog_bench_quill_field_kind kind) {
  switch (kind) {
  case PSLOG_BENCH_QUILL_FIELD_STRING:
    return 's';
  case PSLOG_BENCH_QUILL_FIELD_BOOL:
    return 'b';
  case PSLOG_BENCH_QUILL_FIELD_INT64:
    return 'i';
  case PSLOG_BENCH_QUILL_FIELD_UINT64:
    return 'u';
  case PSLOG_BENCH_QUILL_FIELD_FLOAT64:
    return 'f';
  case PSLOG_BENCH_QUILL_FIELD_NULL:
  default:
    return 'n';
  }
}

static int
pslog_bench_quill_field_is_valid(const pslog_bench_quill_field_spec *field) {
  if (field == NULL) {
    return 0;
  }
  if (field->kind == PSLOG_BENCH_QUILL_FIELD_FLOAT64 &&
      !std::isfinite(field->double_value)) {
    return 0;
  }
  return 1;
}

static int pslog_bench_quill_append_field_payload(
    std::string &out, const pslog_bench_quill_field_spec *field) {
  char value_buf[64];
  int len;
  std::string_view key;
  std::string_view value;

  if (field == NULL) {
    return -1;
  }
  if (!pslog_bench_quill_field_is_valid(field)) {
    return -1;
  }
  if (!out.empty()) {
    out.push_back(';');
  }
  out.push_back(pslog_bench_quill_field_kind_code(field->kind));
  out.push_back('|');
  key =
      (field->key != NULL) ? std::string_view(field->key) : std::string_view{};
  pslog_bench_quill_append_hex(out, key);
  out.push_back('|');

  switch (field->kind) {
  case PSLOG_BENCH_QUILL_FIELD_STRING:
    value = (field->string_value != NULL)
                ? std::string_view(field->string_value)
                : std::string_view{};
    pslog_bench_quill_append_hex(out, value);
    break;
  case PSLOG_BENCH_QUILL_FIELD_BOOL:
    out.push_back(field->bool_value ? '1' : '0');
    break;
  case PSLOG_BENCH_QUILL_FIELD_INT64:
    out.append(std::to_string(field->signed_value));
    break;
  case PSLOG_BENCH_QUILL_FIELD_UINT64:
    out.append(std::to_string(field->unsigned_value));
    break;
  case PSLOG_BENCH_QUILL_FIELD_FLOAT64:
    len = snprintf(value_buf, sizeof(value_buf), "%.17g", field->double_value);
    if (len > 0) {
      out.append(value_buf, (size_t)len);
    } else {
      out.push_back('0');
    }
    break;
  case PSLOG_BENCH_QUILL_FIELD_NULL:
  default:
    break;
  }
  return 0;
}

static int
pslog_bench_quill_dispatch(pslog_bench_quill_logger *logger,
                           const pslog_bench_quill_entry_spec *entry) {
  std::string payload;
  size_t i;

  if (logger == NULL || entry == NULL) {
    return -1;
  }

  payload.reserve(entry->field_count * 32u);
  for (i = 0u; i < entry->field_count; ++i) {
    if (pslog_bench_quill_append_field_payload(payload, &entry->fields[i]) !=
        0) {
      return -1;
    }
  }
  pslog_bench_quill_log_runtime(
      logger, entry->level, "{m_lvl} {m_msg} {m_fields}",
      pslog_bench_quill_level_text(entry->level),
      (entry->message != NULL) ? entry->message : "", payload.c_str());
  return 0;
}

class pslog_bench_quill_sink final : public quill::Sink {
public:
  pslog_bench_quill_sink() = default;

  void reset(size_t capture_limit) {
    std::lock_guard<std::mutex> lock(_mu);
    _bytes = 0ull;
    _capture_limit = capture_limit;
    _capture.clear();
  }

  unsigned long long bytes() const {
    std::lock_guard<std::mutex> lock(_mu);
    return _bytes;
  }

  std::string capture() const {
    std::lock_guard<std::mutex> lock(_mu);
    return _capture;
  }

  void write_log(
      quill::MacroMetadata const * /* log_metadata */, uint64_t log_timestamp,
      std::string_view /* thread_id */, std::string_view /* thread_name */,
      std::string const & /* process_id */, std::string_view /* logger_name */,
      quill::LogLevel log_level, std::string_view /* log_level_description */,
      std::string_view /* log_level_short_code */,
      std::vector<std::pair<std::string, std::string>> const *named_args,
      std::string_view /* log_message */,
      std::string_view /* log_statement */) override {
    std::string line;
    std::string level_text_storage;
    std::string fields_payload_storage;
    std::string_view level_text;
    std::string_view message;
    std::string_view fields_payload;

    level_text = pslog_bench_quill_level_text(
        pslog_bench_quill_level_from_quill(log_level));
    if (named_args != NULL) {
      for (size_t i = 0u; i < named_args->size(); ++i) {
        auto const &named = (*named_args)[i];

        if (named.first == "m_lvl") {
          level_text_storage = named.second;
          level_text = level_text_storage;
          continue;
        }
        if (named.first == "m_msg") {
          message = named.second;
          continue;
        }
        if (named.first == "m_fields") {
          fields_payload_storage = named.second;
          fields_payload = fields_payload_storage;
        }
      }
    }

    line.reserve(256u);
    line.append("{\"ts\":\"");
    line.append(std::to_string(log_timestamp));
    line.append("\",\"lvl\":\"");
    line.append(level_text.data(), level_text.size());
    line.push_back('"');
    line.append(",\"msg\":");
    append_json_string(line, message);
    append_fields_from_payload(line, fields_payload);
    line.append("}\n");

    std::lock_guard<std::mutex> lock(_mu);
    _bytes += (unsigned long long)line.size();
    if (_capture_limit > 0u && _capture.size() < _capture_limit - 1u) {
      size_t remaining;
      size_t copy_len;

      remaining = _capture_limit - _capture.size() - 1u;
      copy_len = line.size();
      if (copy_len > remaining) {
        copy_len = remaining;
      }
      _capture.append(line.data(), copy_len);
    }
  }

  void flush_sink() noexcept override {}
  void run_periodic_tasks() noexcept override {}

private:
  static pslog_bench_quill_level
  pslog_bench_quill_level_from_quill(quill::LogLevel level) {
    switch (level) {
    case quill::LogLevel::TraceL3:
    case quill::LogLevel::TraceL2:
    case quill::LogLevel::TraceL1:
      return PSLOG_BENCH_QUILL_TRACE;
    case quill::LogLevel::Debug:
      return PSLOG_BENCH_QUILL_DEBUG;
    case quill::LogLevel::Info:
      return PSLOG_BENCH_QUILL_INFO;
    case quill::LogLevel::Warning:
      return PSLOG_BENCH_QUILL_WARN;
    case quill::LogLevel::Error:
      return PSLOG_BENCH_QUILL_ERROR;
    case quill::LogLevel::Critical:
      return PSLOG_BENCH_QUILL_FATAL;
    default:
      return PSLOG_BENCH_QUILL_INFO;
    }
  }

  static void append_fields_from_payload(std::string &line,
                                         std::string_view payload) {
    size_t offset = 0u;

    while (offset < payload.size()) {
      size_t next = payload.find(';', offset);
      std::string_view record;
      size_t first_sep;
      size_t second_sep;
      std::string key;
      char kind;

      if (next == std::string_view::npos) {
        next = payload.size();
      }
      record = payload.substr(offset, next - offset);
      offset = next + 1u;
      if (record.size() < 3u) {
        continue;
      }

      kind = record[0];
      first_sep = record.find('|');
      if (first_sep == std::string_view::npos) {
        continue;
      }
      second_sep = record.find('|', first_sep + 1u);
      if (second_sep == std::string_view::npos) {
        continue;
      }
      if (pslog_bench_quill_append_decoded_hex(
              key, record.substr(first_sep + 1u,
                                 second_sep - first_sep - 1u)) != 0) {
        continue;
      }

      line.push_back(',');
      append_json_string(line, key);
      line.push_back(':');
      if (kind == 's') {
        std::string value;

        if (pslog_bench_quill_append_decoded_hex(
                value, record.substr(second_sep + 1u)) != 0) {
          append_json_string(line, "");
        } else {
          append_json_string(line, value);
        }
        continue;
      }
      append_json_value(line, kind, record.substr(second_sep + 1u));
    }
  }

  static void append_json_string(std::string &out, std::string_view value) {
    size_t i;

    out.push_back('"');
    for (i = 0u; i < value.size(); ++i) {
      unsigned char ch;
      char buf[7];

      ch = (unsigned char)value[i];
      switch (ch) {
      case '"':
        out.append("\\\"");
        break;
      case '\\':
        out.append("\\\\");
        break;
      case '\b':
        out.append("\\b");
        break;
      case '\f':
        out.append("\\f");
        break;
      case '\n':
        out.append("\\n");
        break;
      case '\r':
        out.append("\\r");
        break;
      case '\t':
        out.append("\\t");
        break;
      default:
        if (ch < 0x20u) {
          snprintf(buf, sizeof(buf), "\\u%04x", (unsigned int)ch);
          out.append(buf);
        } else {
          out.push_back((char)ch);
        }
        break;
      }
    }
    out.push_back('"');
  }

  static void append_json_value(std::string &out, char kind,
                                std::string_view value) {
    switch (kind) {
    case 's':
    case 'm':
      append_json_string(out, value);
      break;
    case 'b':
      if (value == "true" || value == "1") {
        out.append("true");
      } else {
        out.append("false");
      }
      break;
    case 'i':
    case 'u':
    case 'f':
      if (value.empty()) {
        out.push_back('0');
      } else {
        out.append(value.data(), value.size());
      }
      break;
    case 'n':
      out.append("null");
      break;
    default:
      append_json_string(out, value);
      break;
    }
  }

  mutable std::mutex _mu;
  unsigned long long _bytes{0ull};
  size_t _capture_limit{0u};
  std::string _capture;
};

struct pslog_bench_quill_runtime {
  std::shared_ptr<pslog_bench_quill_sink> sink;
  pslog_bench_quill_logger logger;
  int ready;
};

static pslog_bench_quill_runtime &pslog_bench_quill_runtime_instance(void) {
  static pslog_bench_quill_runtime runtime{};
  static std::once_flag once;

  std::call_once(once, []() {
    quill::BackendOptions backend_options;
    quill::PatternFormatterOptions formatter{"", "%H:%M:%S.%Qns",
                                             quill::Timezone::GmtTime};
    std::shared_ptr<pslog_bench_quill_sink> sink;

    backend_options.sink_min_flush_interval = std::chrono::milliseconds{0};
    quill::Backend::start<pslog_bench_quill_frontend_options>(
        backend_options, quill::SignalHandlerOptions{});

    sink = std::dynamic_pointer_cast<pslog_bench_quill_sink>(
        pslog_bench_quill_frontend::create_or_get_sink<pslog_bench_quill_sink>(
            "pslog_bench_quill_sink"));
    runtime.sink = sink;
    runtime.logger.logger = pslog_bench_quill_frontend::create_or_get_logger(
        "pslog_bench_quill_logger", std::move(sink), formatter);
    if (runtime.logger.logger != NULL) {
      runtime.logger.logger->set_log_level(quill::LogLevel::TraceL3);
      runtime.ready = 1;
    }
  });

  return runtime;
}

} // namespace

extern "C" int pslog_bench_quill_available(void) {
  return pslog_bench_quill_runtime_instance().ready;
}

extern "C" int
pslog_bench_quill_run(const pslog_bench_quill_entry_spec *entries,
                      size_t entry_count, unsigned long long iterations,
                      pslog_bench_quill_result *result) {
  pslog_bench_quill_runtime &runtime = pslog_bench_quill_runtime_instance();
  unsigned long long start_ns;
  unsigned long long end_ns;
  unsigned long long i;
  size_t index;

  if (!runtime.ready || entries == NULL || entry_count == 0u ||
      result == NULL) {
    return -1;
  }

  runtime.sink->reset(0u);
  start_ns = pslog_bench_quill_now_ns();
  index = 0u;
  for (i = 0ull; i < iterations; ++i) {
    if (pslog_bench_quill_dispatch(&runtime.logger, &entries[index]) != 0) {
      return -1;
    }
    ++index;
    if (index == entry_count) {
      index = 0u;
    }
  }
  runtime.logger.logger->flush_log();
  end_ns = pslog_bench_quill_now_ns();

  result->elapsed_ns = end_ns - start_ns;
  result->bytes = runtime.sink->bytes();
  result->ops = iterations;
  return 0;
}

extern "C" int
pslog_bench_quill_render(const pslog_bench_quill_entry_spec *entry,
                         char *buffer, size_t capacity, size_t *written) {
  pslog_bench_quill_runtime &runtime = pslog_bench_quill_runtime_instance();
  std::string capture;
  size_t copy_len;

  if (!runtime.ready || entry == NULL || buffer == NULL || capacity == 0u) {
    return -1;
  }

  runtime.sink->reset(capacity);
  if (pslog_bench_quill_dispatch(&runtime.logger, entry) != 0) {
    return -1;
  }
  runtime.logger.logger->flush_log();

  capture = runtime.sink->capture();
  if (runtime.sink->bytes() != (unsigned long long)capture.size()) {
    return -1;
  }
  copy_len = capture.size();
  if (copy_len >= capacity) {
    copy_len = capacity - 1u;
  }
  if (copy_len > 0u) {
    memcpy(buffer, capture.data(), copy_len);
  }
  buffer[copy_len] = '\0';
  if (written != NULL) {
    *written = copy_len;
  }
  return 0;
}
