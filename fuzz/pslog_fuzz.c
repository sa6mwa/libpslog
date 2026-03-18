#include "pslog.h"

#include <stddef.h>
#include <stdlib.h>
#include <string.h>

struct fuzz_sink {
  char data[16384];
  size_t len;
};

static char fuzz_key_char(unsigned char ch) {
  static const char alphabet[] = "abcdefghijklmnopqrstuvwxyz0123456789_";

  return alphabet[ch % (sizeof(alphabet) - 1u)];
}

static int fuzz_sink_write(void *userdata, const char *data, size_t len,
                           size_t *written) {
  struct fuzz_sink *sink;
  size_t copy_len;

  sink = (struct fuzz_sink *)userdata;
  copy_len = len;
  if (copy_len > sizeof(sink->data) - sink->len - 1u) {
    copy_len = sizeof(sink->data) - sink->len - 1u;
  }
  memcpy(sink->data + sink->len, data, copy_len);
  sink->len += copy_len;
  sink->data[sink->len] = '\0';
  if (written != NULL) {
    *written = copy_len;
  }
  return 0;
}

static void strip_ansi(char *dst, size_t dst_size, const char *src) {
  size_t out;

  out = 0u;
  while (*src != '\0' && out + 1u < dst_size) {
    if ((unsigned char)*src == 0x1bu && src[1] == '[') {
      src += 2;
      while (*src != '\0' && *src != 'm') {
        ++src;
      }
      if (*src == 'm') {
        ++src;
      }
      continue;
    }
    dst[out] = *src;
    ++out;
    ++src;
  }
  dst[out] = '\0';
}

int LLVMFuzzerTestOneInput(const unsigned char *data, size_t size) {
  struct fuzz_sink plain_sink;
  struct fuzz_sink color_sink;
  pslog_config config;
  pslog_logger *plain_log;
  pslog_logger *color_log;
  pslog_field fields[3];
  char msg[64];
  char key[32];
  char value[64];
  char stripped[16384];
  size_t msg_len;
  size_t key_len;
  size_t value_len;

  if (size == 0u) {
    return 0;
  }

  memset(&plain_sink, 0, sizeof(plain_sink));
  memset(&color_sink, 0, sizeof(color_sink));
  pslog_default_config(&config);
  config.mode = (data[0] & 1u) ? PSLOG_MODE_JSON : PSLOG_MODE_CONSOLE;
  config.timestamps = 0;
  config.output.write = fuzz_sink_write;
  config.output.close = NULL;
  config.output.isatty = NULL;
  config.output.userdata = &plain_sink;
  config.color = PSLOG_COLOR_NEVER;

  plain_log = pslog_new(&config);
  config.output.userdata = &color_sink;
  config.color = PSLOG_COLOR_ALWAYS;
  color_log = pslog_new(&config);
  if (plain_log == NULL || color_log == NULL) {
    if (plain_log != NULL) {
      plain_log->destroy(plain_log);
    }
    if (color_log != NULL) {
      color_log->destroy(color_log);
    }
    return 0;
  }

  msg_len = size > sizeof(msg) - 1u ? sizeof(msg) - 1u : size;
  memcpy(msg, data, msg_len);
  msg[msg_len] = '\0';

  key_len = size > 8u ? 8u : size;
  if (key_len == 0u) {
    key[0] = 'k';
    key[1] = 'e';
    key[2] = 'y';
    key[3] = '\0';
    key_len = 3u;
  } else {
    size_t i;

    for (i = 0u; i < key_len; ++i) {
      key[i] = fuzz_key_char(data[i]);
    }
    key[key_len] = '\0';
  }

  value_len = size > sizeof(value) - 1u ? sizeof(value) - 1u : size;
  memcpy(value, data, value_len);
  value[value_len] = '\0';

  fields[0] = pslog_str(key_len > 0u ? key : "key", value);
  fields[1] = pslog_i64("size", (long)size);
  fields[2] = pslog_bool("odd", (size & 1u) != 0u);

  plain_log->info(plain_log, msg, fields, 3u);
  plain_log->infof(plain_log, msg, "size=%lu odd=%b", (unsigned long)size,
                   (size & 1u) != 0u);
  color_log->info(color_log, msg, fields, 3u);
  color_log->infof(color_log, msg, "size=%lu odd=%b", (unsigned long)size,
                   (size & 1u) != 0u);
  strip_ansi(stripped, sizeof(stripped), color_sink.data);
  if (strcmp(plain_sink.data, stripped) != 0) {
    abort();
  }
  plain_log->destroy(plain_log);
  color_log->destroy(color_log);
  return 0;
}
