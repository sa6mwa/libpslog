#define _POSIX_C_SOURCE 1993009L
#include <time.h>

#include "pslog.h"

void sleep_ms(int ms);

int main(void) {
  pslog_config root_config;
  pslog_config console_config;
  pslog_config json_config;
  pslog_logger *root;
  pslog_logger *console;
  pslog_logger *json;
  char *logger = "pkt.systems/c/libpslog";

  pslog_default_config(&console_config);
  console_config.mode = PSLOG_MODE_CONSOLE;
  console_config.min_level = PSLOG_LEVEL_TRACE;
  console_config.output = pslog_output_from_fp(stdout, 0);
  console_config.palette = &pslog_builtin_palette_outrun_electric;

  pslog_default_config(&json_config);
  json_config.mode = PSLOG_MODE_JSON;
  json_config.min_level = PSLOG_LEVEL_TRACE;
  json_config.output = pslog_output_from_fp(stdout, 0);
  json_config.palette = &pslog_builtin_palette_catppuccin_mocha;

  console = pslog_new(&console_config);
  json = pslog_new(&json_config);
  console = console->withf(console, "mode=%s", "console");
  json = json->withf(json, "mode=%s", "json");

  console->infof(console, "Hi! 😀", NULL);
  sleep_ms(2000);
  json->infof(json, "Structured is cooler 😎", NULL);
  sleep_ms(3000);

  console->warnf(console, "No 🤡", NULL);
  sleep_ms(3000);
  json->debugf(json, "🥱", NULL);
  sleep_ms(2200);
  console->tracef(console,
                  "Sleepy? This is the C port of pslog, it goes vrooom... 💨",
                  NULL);
  sleep_ms(3000);
  json->infof(json, "💯", "cool=%b", 1);
  sleep_ms(500);
  console->infof(console, "Get libpslog now!", "url=%s",
                 "https://pkt.systems/c/libpslog");
}

void sleep_ms(int ms) {
  struct timespec ts;
  ts.tv_sec = (time_t)(ms / 1000);
  ts.tv_nsec = (long)(ms % 1000) * 1000000L;
  nanosleep(&ts, NULL);
}
