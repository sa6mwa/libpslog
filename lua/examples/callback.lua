local pslog = require("pslog")

local chunks = {}
local log = pslog.new_json({
  output = function(chunk)
    chunks[#chunks + 1] = chunk
  end,
  no_color = true,
  disable_timestamp = true,
})

log:info("callback sink", "service", "edge", "ok", true)
log:close()

io.write(table.concat(chunks))
