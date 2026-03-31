local core = require("pslog.core")

local M = {}

local function normalize_sec_nsec(sec, nsec)
  sec = tonumber(sec or 0) or error("sec must be numeric", 3)
  nsec = tonumber(nsec or 0) or error("nsec must be numeric", 3)

  if nsec >= 1000000000 or nsec <= -1000000000 then
    local carry = math.floor(nsec / 1000000000)
    sec = sec + carry
    nsec = nsec - carry * 1000000000
  end

  if sec > 0 and nsec < 0 then
    sec = sec - 1
    nsec = nsec + 1000000000
  elseif sec < 0 and nsec > 0 then
    sec = sec + 1
    nsec = nsec - 1000000000
  end

  return sec, nsec
end

local function split_scaled(value, scale)
  local total = tonumber(value) or error("value must be numeric", 3)
  local sec
  local nsec

  total = total * scale
  if total >= 0 then
    sec = math.floor(total / 1000000000)
  else
    sec = math.ceil(total / 1000000000)
  end
  nsec = total - sec * 1000000000
  return normalize_sec_nsec(sec, nsec)
end

local function copy_table(input)
  local output = {}
  local key, value

  if input == nil then
    return output
  end

  for key, value in pairs(input) do
    output[key] = value
  end
  return output
end

local function copy_exports()
  for key, value in pairs(core) do
    M[key] = value
  end
end

local function normalize_options(output_or_opts, opts)
  local config

  if opts ~= nil then
    config = copy_table(opts)
    if output_or_opts ~= nil then
      config.output = output_or_opts
    end
    return config
  end

  if type(output_or_opts) == "table" then
    return copy_table(output_or_opts)
  end

  if output_or_opts == nil then
    return nil
  end

  return { output = output_or_opts }
end

copy_exports()

function M.new(output_or_opts, opts)
  return core.new_logger(normalize_options(output_or_opts, opts))
end

function M.new_json(output_or_opts, opts)
  local config = normalize_options(output_or_opts, opts) or {}

  config.mode = "json"
  return core.new_logger(config)
end

M.new_structured = M.new_json

function M.from_env(prefix_or_opts, maybe_opts)
  if type(prefix_or_opts) == "table" and maybe_opts == nil then
    return core.new_logger_from_env(nil, prefix_or_opts)
  end
  return core.new_logger_from_env(prefix_or_opts, maybe_opts)
end

function M.bytes(value)
  if type(value) ~= "string" then
    error("pslog.bytes expects a Lua string", 2)
  end
  return core._wrap_bytes(value)
end

function M.time_unix(sec, nsec, utc_offset_minutes)
  sec, nsec = normalize_sec_nsec(sec, nsec)
  return core._wrap_time(
    sec,
    nsec,
    tonumber(utc_offset_minutes or 0) or error("utc_offset_minutes must be numeric", 2)
  )
end

function M.time(spec)
  if type(spec) ~= "table" then
    error("pslog.time expects a table", 2)
  end
  return M.time_unix(
    spec.sec or spec.epoch_seconds or 0,
    spec.nsec or spec.nanoseconds or 0,
    spec.offset or spec.utc_offset_minutes or 0
  )
end

function M.duration(sec, nsec)
  sec, nsec = normalize_sec_nsec(sec, nsec)
  return core._wrap_duration(sec, nsec)
end

function M.duration_ns(value)
  local sec, nsec = split_scaled(value, 1)
  return M.duration(sec, nsec)
end

function M.duration_us(value)
  local sec, nsec = split_scaled(value, 1000)
  return M.duration(sec, nsec)
end

function M.duration_ms(value)
  local sec, nsec = split_scaled(value, 1000000)
  return M.duration(sec, nsec)
end

function M.duration_s(value)
  local sec, nsec = split_scaled(value, 1000000000)
  return M.duration(sec, nsec)
end

function M.errno(code)
  code = math.tointeger(code)
  if code == nil then
    error("pslog.errno expects an integer code", 2)
  end
  return core._wrap_errno(code)
end

function M.trusted(value)
  if type(value) ~= "string" then
    error("pslog.trusted expects a string", 2)
  end
  return core._wrap_trusted(value)
end

function M.u64(value)
  if type(value) == "string" then
    if not value:match("^[0-9]+$") then
      error("pslog.u64 expects a non-negative integer or decimal string", 2)
    end
    return core._wrap_u64(value)
  end

  value = math.tointeger(value)
  if value == nil or value < 0 then
    error("pslog.u64 expects a non-negative integer or decimal string", 2)
  end
  return core._wrap_u64(value)
end

return M
