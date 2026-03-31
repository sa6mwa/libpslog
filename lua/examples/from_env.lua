local pslog = require("pslog")

local log = pslog.from_env("APP_LOG_", {
  output = io.stdout,
  mode = "console",
  no_color = true,
}):with("service", "payments")

log:info("before overrides")
log:debug("set APP_LOG_MODE=json or APP_LOG_LEVEL=debug to see env-driven changes")
