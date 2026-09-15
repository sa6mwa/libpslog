# Performance Logs

`make perf-gate` automatically selects a C/Lua baseline pair from
`baselines/<id>/`. `make bench-freeze-baseline` explicitly captures or refreshes
the current host's detailed baseline. Unknown hosts fail the gate until a
baseline is deliberately captured; normal runs never rewrite baselines.

Each directory has `pure-c-baseline.txt`, `lua-baseline.txt`, and an `identity`
file containing one or more aliases, one per line:

```text
fingerprint-md5 <32 lowercase hex digits>
hostname-md5 <32 lowercase hex digits>
```

Selection searches all detailed aliases first, then all hostname aliases.
Multiple matching directories at the winning level are an error. Multiple
aliases within one directory are allowed. A matching but incomplete baseline
is an error, not a reason to fall back. Comments and blank lines are allowed.

The detailed hash is MD5 of this exact UTF-8 text with LF line endings and a
final newline (values come from the running host):

```text
v1
host=<hostname -s>
os=<uname -s>
arch=<uname -m>
cpu=<CPU model>
logical_cpus=<online logical CPU count>
```

Linux reads the CPU model from `/proc/cpuinfo` and the count from `getconf`;
macOS uses `sysctl`. The hostname-only hash is MD5 of `hostname -s` without its
trailing newline. Kernel and compiler versions are excluded from identity, so
updates do not silently bypass an established baseline. New baselines register
only the detailed hash. A refresh follows an existing detailed alias, preserving
its other aliases. A hostname-only match is retained as history: capture creates
a separate detailed baseline for the current hardware. Add hostname aliases deliberately when broader matching
is wanted. The historical pair retains its original timings and is registered
using only the supplied historical host's hostname hash.

Names are hashed before storage. These identifiers obscure names; they are not
secrets or a security boundary. A hostname fallback intentionally tolerates
hardware differences under that name. Use unique short names for unrelated hosts.

Baseline headers record the measured commit/worktree state, UTC capture time,
compiler, preset, run count, and C/Lua measurement settings. A fresh baseline
establishes a reference for future runs, not evidence against past regressions.
Other files here remain historical artifacts; retain their recorded provenance
and mark unavailable historical details as `not recorded`.
