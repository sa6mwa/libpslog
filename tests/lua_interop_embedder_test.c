#include "pslog_lua.h"

#include <lauxlib.h>
#include <lua.h>
#include <lualib.h>

#include <stdio.h>

#define TEST_ASSERT(expr)                                                      \
  do {                                                                         \
    if (!(expr)) {                                                             \
      fprintf(stderr, "assertion failed: %s (%s:%d)\n", #expr, __FILE__,       \
              __LINE__);                                                       \
      return 1;                                                                \
    }                                                                          \
  } while (0)

extern int luaopen_pslog_core(lua_State *L);

static int preload_pslog_core(lua_State *L) {
  luaL_requiref(L, "pslog.core", luaopen_pslog_core, 1);
  lua_pop(L, 1);
  return 0;
}

static int run_chunk(lua_State *L, const char *chunk) {
  if (luaL_dostring(L, chunk) != LUA_OK) {
    fprintf(stderr, "lua error: %s\n", lua_tostring(L, -1));
    lua_pop(L, 1);
    return 1;
  }
  return 0;
}

static int test_borrowed_logger_writes(lua_State *L) {
  pslog_lua_logger_view view;

  TEST_ASSERT(run_chunk(L,
                        "sink = {}\n"
                        "log = pslog.new({ output = function(chunk) "
                        "sink[#sink + 1] = chunk end, mode = 'json', "
                        "disable_timestamp = true, no_color = true })") == 0);
  lua_getglobal(L, "log");
  TEST_ASSERT(pslog_lua_is_logger(L, -1));
  view.size = sizeof(view);
  view.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &view) == 0);
  TEST_ASSERT(view.logger != NULL);
  pslog_info(view.logger, "borrowed", NULL, 0u);
  lua_pop(L, 1);
  TEST_ASSERT(
      run_chunk(L,
                "payload = table.concat(sink)\n"
                "assert(payload:find('\"msg\":\"borrowed\"', 1, true))") == 0);
  return 0;
}

static int test_borrowed_logger_streams_callback_chunks(lua_State *L) {
  pslog_lua_logger_view view;
  char message[128];
  size_t i;

  for (i = 0u; i + 1u < sizeof(message); ++i) {
    message[i] = (char)('a' + (i % 26u));
  }
  message[sizeof(message) - 1u] = '\0';

  TEST_ASSERT(run_chunk(L, "sink = {}\n"
                           "log = pslog.new({ output = function(chunk) "
                           "sink[#sink + 1] = chunk end, mode = 'json', "
                           "disable_timestamp = true, no_color = true, "
                           "line_buffer_capacity = 16 })") == 0);
  lua_getglobal(L, "log");
  view.size = sizeof(view);
  view.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &view) == 0);
  TEST_ASSERT(view.logger != NULL);
  pslog_info(view.logger, message, NULL, 0u);
  lua_pop(L, 1);
  TEST_ASSERT(
      run_chunk(
          L,
          "payload = table.concat(sink)\n"
          "assert(#sink > 1)\n"
          "assert(payload:find('abcdefghijklmnopqrstuvwxyz', 1, true))") == 0);
  return 0;
}

static int test_borrowed_logger_view_survives_lua_close(lua_State *L) {
  pslog_lua_logger_view view;

  TEST_ASSERT(run_chunk(L,
                        "sink = {}\n"
                        "log = pslog.new({ output = function(chunk) "
                        "sink[#sink + 1] = chunk end, mode = 'json', "
                        "disable_timestamp = true, no_color = true })") == 0);
  lua_getglobal(L, "log");
  view.size = sizeof(view);
  view.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &view) == 0);
  TEST_ASSERT(view.logger != NULL);
  lua_pop(L, 1);

  TEST_ASSERT(run_chunk(L, "assert(log:close())") == 0);
  lua_getglobal(L, "log");
  TEST_ASSERT(!pslog_lua_is_logger(L, -1));
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &view) != 0);
  lua_pop(L, 1);

  pslog_info(view.logger, "borrowed_after_close", NULL, 0u);
  return 0;
}

static int test_ref_logger_derives_and_retains(lua_State *L) {
  pslog_lua_logger_ref ref;
  pslog_field field = pslog_str("service", "embedder");

  TEST_ASSERT(run_chunk(L,
                        "sink = {}\n"
                        "log = pslog.new({ output = function(chunk) "
                        "sink[#sink + 1] = chunk end, mode = 'json', "
                        "disable_timestamp = true, no_color = true })") == 0);
  lua_getglobal(L, "log");
  ref.size = sizeof(ref);
  ref.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_ref_logger(L, -1, &field, 1u, &ref) == 0);
  TEST_ASSERT(ref.borrowed != NULL);
  TEST_ASSERT(ref.derived != NULL);
  lua_pop(L, 1);
  TEST_ASSERT(run_chunk(L, "log = nil; collectgarbage('collect')") == 0);
  pslog_info(ref.derived, "derived", NULL, 0u);
  pslog_lua_unref_logger(&ref);
  TEST_ASSERT(ref.borrowed == NULL);
  TEST_ASSERT(ref.derived == NULL);
  TEST_ASSERT(
      run_chunk(
          L, "payload = table.concat(sink)\n"
             "assert(payload:find('\"msg\":\"derived\"', 1, true))\n"
             "assert(payload:find('\"service\":\"embedder\"', 1, true))") == 0);
  return 0;
}

static int test_ref_logger_survives_lua_close(lua_State *L) {
  pslog_lua_logger_ref ref;

  TEST_ASSERT(run_chunk(L,
                        "sink = {}\n"
                        "log = pslog.new({ output = function(chunk) "
                        "sink[#sink + 1] = chunk end, mode = 'json', "
                        "disable_timestamp = true, no_color = true })") == 0);
  lua_getglobal(L, "log");
  ref.size = sizeof(ref);
  ref.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_ref_logger(L, -1, NULL, 0u, &ref) == 0);
  TEST_ASSERT(ref.borrowed != NULL);
  TEST_ASSERT(ref.derived == NULL);
  lua_pop(L, 1);

  TEST_ASSERT(run_chunk(L, "assert(log:close())") == 0);
  lua_getglobal(L, "log");
  TEST_ASSERT(!pslog_lua_is_logger(L, -1));
  lua_pop(L, 1);

  pslog_info(ref.borrowed, "after_close", NULL, 0u);
  pslog_lua_unref_logger(&ref);
  TEST_ASSERT(ref.borrowed == NULL);
  return 0;
}

static int test_negative_paths(lua_State *L) {
  pslog_lua_logger_view view;
  pslog_lua_logger_ref ref;
  int status;

  lua_pushinteger(L, 42);
  TEST_ASSERT(!pslog_lua_is_logger(L, -1));
  view.size = sizeof(view);
  view.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &view) != 0);
  lua_pop(L, 1);

  TEST_ASSERT(run_chunk(L, "log = pslog.new({ output = function(_) end, "
                           "mode = 'json', disable_timestamp = true })") == 0);
  lua_getglobal(L, "log");
  view.size = sizeof(view) - 1u;
  view.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  status = pslog_lua_check_logger(L, -1, &view);
  TEST_ASSERT(status != 0);
  view.size = sizeof(view);
  view.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION + 1u;
  status = pslog_lua_check_logger(L, -1, &view);
  TEST_ASSERT(status != 0);
  ref.size = sizeof(ref);
  ref.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_ref_logger(L, -1, NULL, 1u, &ref) != 0);
  lua_pop(L, 1);
  TEST_ASSERT(run_chunk(L, "log = pslog.new({ output = function(_) end, "
                           "mode = 'json', disable_timestamp = true })\n"
                           "assert(log:close())") == 0);
  lua_getglobal(L, "log");
  TEST_ASSERT(!pslog_lua_is_logger(L, -1));
  view.size = sizeof(view);
  view.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &view) != 0);
  lua_pop(L, 1);
  return 0;
}

static int test_multi_logger_views(lua_State *L) {
  pslog_lua_logger_view first;
  pslog_lua_logger_view second;

  TEST_ASSERT(run_chunk(L, "a = pslog.new({ output = function(_) end, "
                           "mode = 'json', disable_timestamp = true })\n"
                           "b = pslog.new({ output = function(_) end, "
                           "mode = 'json', disable_timestamp = true })") == 0);
  lua_getglobal(L, "a");
  first.size = sizeof(first);
  first.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &first) == 0);
  lua_pop(L, 1);
  lua_getglobal(L, "b");
  second.size = sizeof(second);
  second.abi_version = PSLOG_LUA_INTEROP_ABI_VERSION;
  TEST_ASSERT(pslog_lua_check_logger(L, -1, &second) == 0);
  lua_pop(L, 1);
  TEST_ASSERT(first.logger != NULL);
  TEST_ASSERT(second.logger != NULL);
  TEST_ASSERT(first.logger != second.logger);
  return 0;
}

int main(void) {
  lua_State *L = luaL_newstate();
  int failed = 0;

  TEST_ASSERT(L != NULL);
  luaL_openlibs(L);
  preload_pslog_core(L);
  if (run_chunk(
          L, "package.path = './lua/?.lua;./lua/?/init.lua;' .. package.path\n"
             "pslog = require('pslog')") != 0) {
    lua_close(L);
    return 1;
  }

  failed |= test_borrowed_logger_writes(L);
  failed |= test_borrowed_logger_streams_callback_chunks(L);
  failed |= test_borrowed_logger_view_survives_lua_close(L);
  failed |= test_ref_logger_derives_and_retains(L);
  failed |= test_ref_logger_survives_lua_close(L);
  failed |= test_negative_paths(L);
  failed |= test_multi_logger_views(L);
  lua_close(L);
  return failed;
}
