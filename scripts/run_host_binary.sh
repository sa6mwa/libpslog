#!/usr/bin/env bash
# Run a native x86_64 Bootlin-built program through the configured host sysroot.
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
exec "$repo_root/scripts/run_sysroot_binary.sh" \
  --build-dir "$repo_root/build/host" \
  "$@"
