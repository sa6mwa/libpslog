local pslog = require("pslog")

local log = pslog.new_json(io.stdout, {
  disable_timestamp = true,
}):with("service", "checkout", "component", "worker")

log:info("ready", "port", 8080, "ok", true)
log:warn("retrying", {
  attempt = 2,
  backoff_ms = 250,
  last_error = "upstream timeout",
})
log:info("typed fields",
  "payload", pslog.bytes(string.char(0, 1, 2, 255)),
  "started_at", pslog.time { sec = 1711737600, nsec = 0, offset = 0 },
  "elapsed", pslog.duration_ms(125),
  "errno", pslog.errno(2),
  "trusted_service", pslog.trusted("checkout.worker")
)
