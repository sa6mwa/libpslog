#include "pslog.h"

#include <stdio.h>
#include <string.h>

struct memory_sink {
  char data[4096];
  size_t len;
};

static int memory_sink_write(void *userdata, const char *data, size_t len,
                             size_t *written) {
  struct memory_sink *sink;
  size_t available;

  sink = (struct memory_sink *)userdata;
  available = sizeof(sink->data) - sink->len - 1u;
  if (len > available) {
    if (written != NULL) {
      *written = 0u;
    }
    return -1;
  }
  memcpy(sink->data + sink->len, data, len);
  sink->len += len;
  sink->data[sink->len] = '\0';
  if (written != NULL) {
    *written = len;
  }
  return 0;
}

int main(void) {
  struct memory_sink sink;
  pslog_config config;
  pslog_field root_field;
  pslog_field event_field;
  pslog_logger *root;
  pslog_logger *child;

  memset(&sink, 0, sizeof(sink));
  pslog_default_config(&config);
  config.mode = PSLOG_MODE_JSON;
  config.color = PSLOG_COLOR_NEVER;
  config.timestamps = 0;
  config.output.write = memory_sink_write;
  config.output.userdata = &sink;
  root = pslog_new(&config);
  if (root == NULL) {
    fputs("failed to create logger\n", stderr);
    return 1;
  }

  root_field = pslog_str("service", "valgrind");
  child = root->with(root, &root_field, 1u);
  if (child == NULL) {
    root->destroy(root);
    fputs("failed to derive logger\n", stderr);
    return 1;
  }
  event_field = pslog_i64("request", 42);
  child->info(child, "facade-memory-check", &event_field, 1u);
  child->destroy(child);
  root->destroy(root);

  if (strstr(sink.data, "\"msg\":\"facade-memory-check\"") == NULL ||
      strstr(sink.data, "\"service\":\"valgrind\"") == NULL ||
      strstr(sink.data, "\"request\":42") == NULL) {
    fputs("unexpected facade output\n", stderr);
    return 1;
  }
  return 0;
}
