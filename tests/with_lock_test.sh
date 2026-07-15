#!/usr/bin/env bash
set -euo pipefail

repo_root=$1
test_root=$(mktemp -d)
cleanup() { rm -rf "$test_root"; }
trap cleanup EXIT HUP INT TERM

lock_dir="$test_root/lua-rock.lock.d"
output_file="$test_root/output"

"$repo_root/scripts/with_lock.sh" "$lock_dir" sh -c 'printf locked >"$1"' sh "$output_file"
[[ "$(cat "$output_file")" == locked ]] || {
  printf 'with_lock did not run the protected command\n' >&2
  exit 1
}
[[ ! -e "$lock_dir" ]] || {
  printf 'with_lock did not remove a completed lock\n' >&2
  exit 1
}

mkdir "$lock_dir"
printf 'pid=999999\nhost=%s\n' "$(hostname 2>/dev/null || printf unknown)" >"$lock_dir/owner"
PSLOG_LOCK_TIMEOUT=2 "$repo_root/scripts/with_lock.sh" "$lock_dir" sh -c 'printf reclaimed >"$1"' sh "$output_file"
[[ "$(cat "$output_file")" == reclaimed ]] || {
  printf 'with_lock did not reclaim a stale same-host lock\n' >&2
  exit 1
}

mkdir "$lock_dir"
printf 'pid=%s\nhost=%s\n' "$$" "$(hostname 2>/dev/null || printf unknown)" >"$lock_dir/owner"
if PSLOG_LOCK_TIMEOUT=1 "$repo_root/scripts/with_lock.sh" "$lock_dir" true 2>"$test_root/error"; then
  printf 'with_lock acquired a live-owner lock\n' >&2
  exit 1
fi
grep -Fq 'timed out waiting for lock' "$test_root/error" || {
  printf 'with_lock did not report timeout for live-owner lock\n' >&2
  cat "$test_root/error" >&2
  exit 1
}
rm -rf "$lock_dir"

mkdir "$lock_dir"
PSLOG_LOCK_TIMEOUT=1 "$repo_root/scripts/with_lock.sh" "$lock_dir" sh -c 'printf ownerless-reclaimed >"$1"' sh "$output_file"
[[ "$(cat "$output_file")" == ownerless-reclaimed ]] || {
  printf 'with_lock did not reclaim an ownerless stale lock at the timeout boundary\n' >&2
  exit 1
}

grep -Fq 'reclaiming' "$repo_root/scripts/with_lock.sh" || {
  printf 'with_lock stale-owner cleanup does not claim reclamation ownership\n' >&2
  exit 1
}
grep -Fq 'mv "$lock_dir" "$reclaim_dir"' "$repo_root/scripts/with_lock.sh" || {
  printf 'with_lock stale-owner cleanup does not rename the stale lock before deletion\n' >&2
  exit 1
}
