#if defined(PSLOG_EXAMPLE_SINGLE_HEADER)
#if (defined(__unix__) || defined(__APPLE__) || defined(__FreeBSD__)) &&       \
    !defined(_POSIX_C_SOURCE)
#define _POSIX_C_SOURCE 200809L
#endif
#define PSLOG_IMPLEMENTATION
#include "pslog_single_header.h"
#else
#include "../include/pslog.h"
#endif

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

static void example_console_and_json(void) {
  pslog_config console_config;
  pslog_config json_config;
  pslog_logger *console;
  pslog_logger *json;
  pslog_logger *next;
  pslog_field fields[2];
#define __FIELDS_LEN sizeof(fields) / sizeof(pslog_field)

  puts("== console and json ==");

  pslog_default_config(&console_config);
  console_config.mode = PSLOG_MODE_CONSOLE;
  console_config.color = PSLOG_COLOR_AUTO;
  console_config.min_level = PSLOG_LEVEL_TRACE;
  console_config.timestamps = 1;
  console_config.output = pslog_output_from_fp(stdout, 0);

  console = pslog_new(&console_config);
  fields[0] = pslog_str("adapter", "libpslog");
  fields[1] = pslog_str("mode", "console");
  next = console->with(console, fields, 2u);
  console->destroy(console);
  console = next;
  next = console->with_level(console, PSLOG_LEVEL_TRACE);
  console->destroy(console);
  console = next;
  console->debug(console, "hello from console mode", NULL, 0u);
  console->info(console, "plain info line", NULL, 0u);
  console->warn(console, "warning line", NULL, 0u);
  console->error(console, "error line", NULL, 0u);
  console->destroy(console);

  pslog_default_config(&json_config);
  json_config.mode = PSLOG_MODE_JSON;
  json_config.color = PSLOG_COLOR_AUTO;
  json_config.min_level = PSLOG_LEVEL_TRACE;
  json_config.timestamps = 1;
  json_config.output = pslog_output_from_fp(stdout, 0);

  json = pslog_new(&json_config);
  fields[0] = pslog_str("adapter", "libpslog");
  fields[1] = pslog_str("mode", "json");
  next = json->with(json, fields, 2u);
  json->destroy(json);
  json = next;
  next = json->with_level_field(json);
  json->destroy(json);
  json = next;
  json->info(json, "hello from json mode", NULL, 0u);
  json->warn(json, "warning line", NULL, 0u);
  json->destroy(json);

  puts("");
}

static void example_log_and_levels(void) {
  pslog_config config;
  pslog_logger *log;
  pslog_logger *next;
  pslog_field fields[2];

  puts("== log() and level controls ==");

  pslog_default_config(&config);
  config.mode = PSLOG_MODE_JSON;
  config.color = PSLOG_COLOR_NEVER;
  config.min_level = PSLOG_LEVEL_TRACE;
  config.timestamps = 1;
  config.output = pslog_output_from_fp(stdout, 0);

  log = pslog_new(&config);
  next = log->with_level_field(log);
  log->destroy(log);
  log = next;
  next = log->with_level(log, PSLOG_LEVEL_DEBUG);
  log->destroy(log);
  log = next;

  fields[0] = pslog_str("path", "/v1/messages");
  fields[1] = pslog_u64("attempt", 3ul);
  log->log(log, PSLOG_LEVEL_INFO, "using log() directly", fields, 2u);
  log->trace(log, "trace should stay hidden", NULL, 0u);
  log->debug(log, "debug is visible", NULL, 0u);
  log->warn(log, "warn is visible", NULL, 0u);
  log->log(log, PSLOG_LEVEL_NOLEVEL, "nolevel event", NULL, 0u);

  log->destroy(log);
  puts("");
}

static void example_typed_fields(void) {
  pslog_config console_config;
  pslog_config json_config;
  pslog_logger *console;
  pslog_logger *json;
  pslog_field fields[8];
  pslog_time_value now_value;
  time_t now;
  unsigned char bytes_data[3];

  puts("== typed fields ==");

  now = time((time_t *)0);
  now_value.epoch_seconds = (long)now;
  now_value.nanoseconds = 0l;
  now_value.utc_offset_minutes = 0;

  bytes_data[0] = 0xdeu;
  bytes_data[1] = 0xadu;
  bytes_data[2] = 0xbeu;

  fields[0] = pslog_time_field("ts_iso", now_value.epoch_seconds, 0l, 0);
  fields[1] = pslog_str("user", "alice");
  fields[2] = pslog_i64("attempts", 3l);
  fields[3] = pslog_f64("latency_ms", 12.34);
  fields[4] = pslog_duration_field("duration", 0l, 123000l);
  fields[5] = pslog_bool("ok", 1);
  fields[6] = pslog_null("status");
  fields[7] = pslog_bytes_field("payload", bytes_data, 3u);

  pslog_default_config(&console_config);
  console_config.mode = PSLOG_MODE_CONSOLE;
  console_config.color = PSLOG_COLOR_AUTO;
  console_config.min_level = PSLOG_LEVEL_TRACE;
  console_config.timestamps = 1;
  console_config.output = pslog_output_from_fp(stdout, 0);

  console = pslog_new(&console_config);
  console->info(console, "typed fields in console mode", fields, 8u);
  console->destroy(console);

  pslog_default_config(&json_config);
  json_config.mode = PSLOG_MODE_JSON;
  json_config.color = PSLOG_COLOR_AUTO;
  json_config.min_level = PSLOG_LEVEL_TRACE;
  json_config.timestamps = 1;
  json_config.output = pslog_output_from_fp(stdout, 0);

  json = pslog_new(&json_config);
  json->info(json, "typed fields in json mode", fields, 8u);
  json->destroy(json);

  puts("");
}

static void example_kvfmt(void) {
  pslog_config config;
  pslog_logger *log;
  pslog_logger *child;
  pslog_logger *next;

  puts("== kvfmt ==");

  pslog_default_config(&config);
  config.mode = PSLOG_MODE_CONSOLE;
  config.color = PSLOG_COLOR_AUTO;
  config.min_level = PSLOG_LEVEL_TRACE;
  config.timestamps = 1;
  config.output = pslog_output_from_fp(stdout, 0);

  log = pslog_new(&config);
  next = log->with_level_field(log);
  log->destroy(log);
  log = next;
  child = log->withf(log, "service=%s subsystem=%s", "api", "worker");
  if (child != NULL) {
    child->info(child, "derived static kvfmt fields", NULL, 0u);
    child->destroy(child);
  }
  log->infof(log, "hello", "user=%s attempts=%d ok=%b ms=%f", "alice", 3, 1,
             12.34);
  log->warnf(log, "oops", "err=%m path=%s code=%d", "/v1/messages", 503);
  log->destroy(log);

  puts("");
}

static void example_env(void) {
  pslog_config seed;
  pslog_logger *log;
  pslog_logger *next;

  puts("== pslog_new_from_env() ==");

  pslog_default_config(&seed);
  seed.output = pslog_output_from_fp(stdout, 0);
  setenv("LOG_MODE", "json", 1);
  setenv("LOG_LEVEL", "trace", 1);
  setenv("LOG_FORCE_COLOR", "true", 1);
  setenv("LOG_PALETTE", "nord", 1);

  log = pslog_new_from_env("LOG_", &seed);
  if (log != NULL) {
    next = log->with_level_field(log);
    log->destroy(log);
    log = next;
    log->debugf(log, "logger from env", "source=%s ok=%b", "env", 1);
    log->destroy(log);
  }

  unsetenv("LOG_MODE");
  unsetenv("LOG_LEVEL");
  unsetenv("LOG_FORCE_COLOR");
  unsetenv("LOG_PALETTE");

  puts("");
}

static void example_palettes(void) {
  size_t i;

  puts("== palette sweep ==");

  for (i = 0u; i < pslog_palette_count(); ++i) {
    pslog_config console_config;
    pslog_config json_config;
    const char *name;
    const pslog_palette *palette;
    pslog_logger *console;
    pslog_logger *json;
    pslog_logger *next;
    pslog_field fields[3];
    int err;

    name = pslog_palette_name(i);
    palette = pslog_palette_by_name(name);

    pslog_default_config(&console_config);
    console_config.mode = PSLOG_MODE_CONSOLE;
    console_config.color = PSLOG_COLOR_AUTO;
    console_config.min_level = PSLOG_LEVEL_TRACE;
    console_config.timestamps = 1;
    console_config.output = pslog_output_from_fp(stdout, 0);
    console_config.palette = palette;

    pslog_default_config(&json_config);
    json_config.mode = PSLOG_MODE_JSON;
    json_config.color = PSLOG_COLOR_AUTO;
    json_config.min_level = PSLOG_LEVEL_TRACE;
    json_config.timestamps = 1;
    json_config.output = pslog_output_from_fp(stdout, 0);
    json_config.palette = palette;

    console = pslog_new(&console_config);
    json = pslog_new(&json_config);

    err = ENOMSG;
    fields[0] = pslog_str("palette", name);
    fields[1] = pslog_bool("cool", 1);
    fields[2] = pslog_errno("err", err);

    next = console->with_level_field(console);
    console->destroy(console);
    console = next;
    next = json->with_level_field(json);
    json->destroy(json);
    json = next;
    console->trace(console, "console palette", fields, __FIELDS_LEN);
    console->debug(console, "console palette", fields, __FIELDS_LEN);
    console->info(console, "console palette", fields, __FIELDS_LEN);
    console->warn(console, "console palette", fields, __FIELDS_LEN);
    console->error(console, "console palette", fields, __FIELDS_LEN);

    json->trace(json, "json palette", fields, __FIELDS_LEN);
    json->debug(json, "json palette", fields, __FIELDS_LEN);
    json->info(json, "json palette", fields, __FIELDS_LEN);
    json->warn(json, "json palette", fields, __FIELDS_LEN);
    json->error(json, "json palette", fields, __FIELDS_LEN);

    puts("");

    console->destroy(console);
    json->destroy(json);
  }
}

int main(void) {
  puts("libpslog examples/example.c");
  puts("library mode:");
  puts("  cc -I ../build/host/generated/include \\");
  puts("    -I ../include -o example example.c ../build/host/libpslog.a");
  puts("single-header mode:");
  puts("  cc -DPSLOG_EXAMPLE_SINGLE_HEADER=1 \\");
  puts("    -I ../build/host/generated/include -o example example.c");
  puts("");

  example_console_and_json();
  example_log_and_levels();
  example_typed_fields();
  example_kvfmt();
  example_env();
  example_palettes();
  return 0;
}
