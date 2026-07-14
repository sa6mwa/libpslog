#ifndef PSLOG_LUA_H
#define PSLOG_LUA_H

#include <stddef.h>

#include <lua.h>

#include <pslog.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * ABI version for the public Lua/C logger interop structs.
 *
 * Callers initialize `size` and `abi_version` before passing a struct to an
 * interop function. libpslog rejects mismatches instead of guessing layout.
 */
#define PSLOG_LUA_INTEROP_ABI_VERSION 1u

/**
 * Borrowed C view of a live Lua-owned `pslog.logger` userdata.
 *
 * `logger` remains owned by Lua. It is valid only while the backing userdata is
 * alive and only on the owning `lua_State`. Calling `log:close()` from Lua
 * marks the userdata closed and future checks reject it, but an already
 * returned borrowed logger pointer remains allocated until the userdata is
 * collected.
 */
typedef struct pslog_lua_logger_view {
  /** Size of this struct as compiled by the caller. */
  size_t size;
  /** ABI version expected by the caller. */
  unsigned int abi_version;
  /** Borrowed logger pointer owned by the Lua userdata. */
  pslog_logger *logger;
  /** Reserved for future view flags. */
  unsigned int flags;
} pslog_lua_logger_view;

/**
 * Retained reference to a Lua-owned logger plus an optional derived C logger.
 *
 * `borrowed` remains owned by Lua. `derived`, when non-NULL, is owned by this
 * reference and is destroyed by `pslog_lua_unref_logger()`. The registry
 * reference keeps the Lua userdata and borrowed logger allocation alive even if
 * Lua code later calls `log:close()`.
 */
typedef struct pslog_lua_logger_ref {
  /** Size of this struct as compiled by the caller. */
  size_t size;
  /** ABI version expected by the caller. */
  unsigned int abi_version;
  /** Owning Lua state for `registry_ref` and `borrowed`. */
  lua_State *L;
  /** Lua registry reference keeping the logger userdata alive. */
  int registry_ref;
  /** Borrowed logger pointer owned by the Lua userdata. */
  pslog_logger *borrowed;
  /** Optional derived logger owned by this ref. */
  pslog_logger *derived;
} pslog_lua_logger_ref;

/**
 * Return non-zero when `index` is a live `pslog.logger` userdata on `L`.
 *
 * Closed logger userdata is treated as not live.
 *
 * This function does not raise Lua errors for ordinary validation failures.
 */
PSLOG_API int pslog_lua_is_logger(lua_State *L, int index);

/**
 * Validate a live `pslog.logger` userdata and fill a borrowed view.
 *
 * `out->size` and `out->abi_version` must be initialized by the caller. Returns
 * 0 on success or an errno-style status such as `EINVAL`, `EPROTO`, or
 * `ENOENT` on validation failure. This function does not transfer ownership and
 * does not raise Lua errors for ordinary validation failures. Closed logger
 * userdata is rejected.
 */
PSLOG_API int pslog_lua_check_logger(lua_State *L, int index,
                                     pslog_lua_logger_view *out);

/**
 * Retain a Lua logger and optionally derive a C-owned logger with `fields`.
 *
 * `out->size` and `out->abi_version` must be initialized by the caller. When
 * `field_count` is zero, `derived` is NULL and consumers should use `borrowed`.
 * When `field_count` is non-zero, `derived` is owned by `out` and released by
 * `pslog_lua_unref_logger()`. Closed logger userdata is rejected.
 */
PSLOG_API int pslog_lua_ref_logger(lua_State *L, int index,
                                   const pslog_field *fields,
                                   size_t field_count,
                                   pslog_lua_logger_ref *out);

/**
 * Release a retained Lua logger reference and its helper-owned derived logger.
 *
 * The borrowed Lua logger is never destroyed by this function. NULL is allowed.
 */
PSLOG_API void pslog_lua_unref_logger(pslog_lua_logger_ref *ref);

#ifdef __cplusplus
}
#endif

#endif
