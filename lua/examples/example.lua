local pslog = require("pslog")

local function section(title)
  io.write(("== %s ==\n"):format(title))
end

local function newline()
  io.write("\n")
end

local function example_console_and_json()
  local console
  local json

  section("console and json")

  console = pslog.new(io.stdout, {
    min_level = "trace",
  }):with("adapter", "lua-pslog", "mode", "console")
  console = console:with_level("trace")
  console:debug("hello from console mode")
  console:info("plain info line")
  console:warn("warning line")
  console:error("error line")

  json = pslog.new_json(io.stdout, {
    min_level = "trace",
  }):with("adapter", "lua-pslog", "mode", "json")
  json = json:with_level_field()
  json:info("hello from json mode")
  json:warn("warning line")

  newline()
end

local function example_log_and_levels()
  local log

  section("log() and level controls")

  log = pslog.new_json(io.stdout, {
    color = "never",
    min_level = "trace",
  })
  log = log:with_level_field()
  log = log:with_level("debug")

  log:log("info", "using log() directly", "path", "/v1/messages", "attempt", pslog.u64(3))
  log:trace("trace should stay hidden")
  log:debug("debug is visible")
  log:warn("warn is visible")

  newline()
end

local function example_typed_fields()
  local now = os.time()
  local console
  local json

  section("typed fields")

  console = pslog.new(io.stdout, {
    min_level = "trace",
  })
  console:info("typed fields in console mode",
    "ts_iso", pslog.time { sec = now, nsec = 0, offset = 0 },
    "user", "alice",
    "attempts", 3,
    "latency_ms", 12.34,
    "duration", pslog.duration_ms(123),
    "ok", true,
    "status", nil,
    "payload", pslog.bytes(string.char(0xde, 0xad, 0xbe)),
    "errno", pslog.errno(42),
    "trusted_service", pslog.trusted("checkout.worker"),
    "req_id", pslog.u64("18446744073709551615")
  )

  json = pslog.new_json(io.stdout, {
    min_level = "trace",
  })
  json:info("typed fields in json mode",
    "ts_iso", pslog.time { sec = now, nsec = 0, offset = 0 },
    "user", "alice",
    "attempts", 3,
    "latency_ms", 12.34,
    "duration", pslog.duration_ms(123),
    "ok", true,
    "status", nil,
    "payload", pslog.bytes(string.char(0xde, 0xad, 0xbe)),
    "errno", pslog.errno(42),
    "trusted_service", pslog.trusted("checkout.worker"),
    "req_id", pslog.u64("18446744073709551615")
  )

  newline()
end

local function example_palettes()
  local palettes = pslog.available_palettes()

  section("palette sweep")

  for _, name in ipairs(palettes) do
    local console = pslog.new(io.stdout, {
      min_level = "trace",
      palette = name,
    }):with_level_field()
    local json = pslog.new_json(io.stdout, {
      min_level = "trace",
      palette = name,
    }):with_level_field()

    console:trace("console palette", "palette", name, "cool", true, "err", pslog.errno(42))
    console:info("console palette", "palette", name, "cool", true, "err", pslog.errno(42))
    console:warn("console palette", "palette", name, "cool", true, "err", pslog.errno(42))
    console:error("console palette", "palette", name, "cool", true, "err", pslog.errno(42))

    json:trace("json palette", "palette", name, "cool", true, "err", pslog.errno(42))
    json:info("json palette", "palette", name, "cool", true, "err", pslog.errno(42))
    json:warn("json palette", "palette", name, "cool", true, "err", pslog.errno(42))
    json:error("json palette", "palette", name, "cool", true, "err", pslog.errno(42))

    newline()
  end
end

io.write("lua-pslog lua/examples/example.lua\n\n")
example_console_and_json()
example_log_and_levels()
example_typed_fields()
example_palettes()
