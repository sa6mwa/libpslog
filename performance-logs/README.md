# Performance Logs

Committed benchmark artifacts live here instead of under `build/` so perf
history stays reviewable in git.

Conventions:

- keep pure C benchmark summaries and raw runs here
- treat `pure-c-baseline.txt` as the current local comparison baseline
- include enough context in filenames to explain what changed and how the run
  was measured
- prefer pinned runs when comparing deltas, for example `taskset -c 0`
