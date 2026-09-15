# Performance Logs

`make perf-gate` automatically selects a C/Lua baseline pair from
`baselines/<id>/`. `make bench-freeze-baseline` explicitly captures or refreshes
the current host's detailed baseline. Unknown hosts fail the gate until a
baseline is deliberately captured; normal runs never rewrite baselines.

Each directory has `pure-c-baseline.txt`, `lua-baseline.txt`, and an `identity`
file containing one or more aliases, one per line:

```text
fingerprint-md5 <32 lowercase hex digits>
nodename-md5 <32 lowercase hex digits>
```

Selection searches all detailed aliases first, then all node-name aliases.
Multiple matching directories at the winning level are an error. Multiple
aliases within one directory are allowed. A matching but incomplete baseline
is an error, not a reason to fall back. Comments and blank lines are allowed.

The detailed hash is MD5 of this exact UTF-8 text with LF line endings and a
final newline (values come from the running host):

```text
v1
host=<uname -n>
os=<uname -s>
arch=<uname -m>
cpu=<CPU model>
logical_cpus=<online logical CPU count>
```

Linux reads the CPU model from `/proc/cpuinfo` and the count from `getconf`;
macOS uses `sysctl`. The node-name hash is MD5 of `uname -n` without its
trailing newline. Kernel and compiler versions are excluded from identity, so
updates do not silently bypass an established baseline. New baselines register
only the detailed hash. A refresh follows an existing detailed alias, preserving
its other aliases. A node-name-only match is retained as history: capture creates
a separate detailed baseline for the current hardware. Add node-name aliases deliberately when broader matching
is wanted. The historical pair retains its original timings and is registered
using only the supplied historical host's node-name hash.

Names are hashed before storage. These identifiers obscure names; they are not
secrets or a security boundary. A node-name fallback intentionally tolerates
hardware differences under that node name. Use distinct node names for unrelated hosts.

Baseline headers record the measured commit/worktree state, UTC capture time,
compiler, preset, run count, and C/Lua measurement settings. A fresh baseline
establishes a reference for future runs, not evidence against past regressions.
Other files here remain historical artifacts; retain their recorded provenance
and mark unavailable historical details as `not recorded`.
