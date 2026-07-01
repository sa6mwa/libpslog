# Performance Logs

Committed benchmark artifacts live here instead of under `build/` so perf
history stays reviewable in git.

Conventions:

- keep pure C and Lua benchmark summaries and raw runs here
- treat `pure-c-baseline.txt` as the current pure C comparison baseline
- treat `lua-baseline.txt` as the current Lua binding comparison baseline
- include enough context in filenames to explain what changed and how the run
  was measured
- prefer pinned runs when comparing deltas, for example `taskset -c 0`
- start every committed artifact with a `# pslog_perf_artifact:` metadata
  header before benchmark rows
- record provenance in that header: commit or baseline/current source files,
  measured commit, command, iterations, run count, pinning, CMake
  preset/options, optional benchmark dependencies, host, CPU model, OS/kernel,
  compiler, date, and summary method
- when backfilling historical artifacts, say `not recorded` for provenance
  that cannot be recovered instead of implying precision from filenames alone
