#define _POSIX_C_SOURCE 1993009L
#include <stdio.h>
#include <time.h>

#include "pslog.h"

void sleep_ms(int ms);

int main(void) {
  pslog_config console_config;
  pslog_config json_config;
  pslog_logger *console;
  pslog_logger *json;
  pslog_logger *console_view;
  pslog_logger *json_view;

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
  console_view = console->withf(console, "mode=%s", "console");
  json_view = json->withf(json, "mode=%s", "json");

  console_view->infof(console_view, "Hi! 😀", NULL);
  sleep_ms(2000);
  json_view->infof(json_view, "Structured is cooler 😎", NULL);
  sleep_ms(3000);

  console_view->warnf(console_view, "No 🤡", NULL);
  sleep_ms(3000);
  json_view->debugf(json_view, "🥱", NULL);
  sleep_ms(2200);
  console_view->tracef(
      console_view, "Sleepy? This is the C port of pslog, it goes vrooom... 💨",
      NULL);
  sleep_ms(3000);
  json_view->infof(json_view, "💯", "cool=%b", 1);
  sleep_ms(500);
  console_view->infof(console_view, "Get libpslog now!", "url=%s",
                      "https://pkt.systems/c/libpslog");
  console_view->destroy(console_view);
  json_view->destroy(json_view);
  console->destroy(console);
  json->destroy(json);
  return 0;
}

void sleep_ms(int ms) {
  struct timespec ts;
  ts.tv_sec = (time_t)(ms / 1000);
  ts.tv_nsec = (long)(ms % 1000) * 1000000L;
  nanosleep(&ts, NULL);
}
