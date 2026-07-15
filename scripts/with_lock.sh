#!/usr/bin/env sh
# Run a command while holding a portable advisory lock implemented with atomic
# directory creation. This avoids GNU flock, which is not present on stock macOS.
set -eu

if [ "$#" -lt 2 ]; then
    printf 'usage: %s LOCK_DIR COMMAND [ARG ...]\n' "$0" >&2
    exit 2
fi

lock_dir=$1
shift
timeout=${PSLOG_LOCK_TIMEOUT:-120}

case "$timeout" in
    ''|*[!0-9]*)
        printf 'with_lock: PSLOG_LOCK_TIMEOUT must be a non-negative integer\n' >&2
        exit 2
        ;;
esac

parent_dir=$(dirname -- "$lock_dir")
mkdir -p "$parent_dir"
owner_file="$lock_dir/owner"
owner_pid=$$
owner_host=$(hostname 2>/dev/null || printf 'unknown')

write_owner() {
    printf 'pid=%s\nhost=%s\n' "$owner_pid" "$owner_host" >"$owner_file"
}

owner_value() {
    key=$1
    sed -n "s/^$key=//p" "$owner_file" 2>/dev/null | tail -n 1
}

try_reclaim_stale_lock() {
    reclaim_marker="$lock_dir/reclaiming"
    reclaim_dir="${lock_dir}.reclaimed.${owner_pid}"

    if [ ! -e "$owner_file" ]; then
        return
    fi

    lock_pid=$(owner_value pid)
    lock_host=$(owner_value host)
    if [ -n "$lock_pid" ] &&
       [ "$lock_host" = "$owner_host" ] &&
       ! kill -0 "$lock_pid" 2>/dev/null; then
        if mkdir "$reclaim_marker" 2>/dev/null; then
            if [ "$(owner_value pid)" = "$lock_pid" ] &&
               [ "$(owner_value host)" = "$lock_host" ] &&
               ! kill -0 "$lock_pid" 2>/dev/null &&
               mv "$lock_dir" "$reclaim_dir" 2>/dev/null; then
                rm -rf "$reclaim_dir"
            else
                rmdir "$reclaim_marker" 2>/dev/null || true
            fi
        fi
    fi
}

try_reclaim_ownerless_lock() {
    if [ ! -e "$owner_file" ]; then
        rmdir "$lock_dir" 2>/dev/null || true
    fi
}

elapsed=0
while :; do
    if mkdir "$lock_dir" 2>/dev/null; then
        if write_owner 2>/dev/null &&
           [ "$(owner_value pid)" = "$owner_pid" ] &&
           [ "$(owner_value host)" = "$owner_host" ]; then
            break
        fi
        rm -rf "$lock_dir" 2>/dev/null || true
        elapsed=0
        continue
    fi

    try_reclaim_stale_lock
    if [ ! -e "$lock_dir" ]; then
        elapsed=0
        continue
    fi

    if [ "$elapsed" -ge "$timeout" ]; then
        try_reclaim_ownerless_lock
        if [ ! -e "$lock_dir" ]; then
            elapsed=0
            continue
        fi
        printf 'with_lock: timed out waiting for lock: %s\n' "$lock_dir" >&2
        exit 1
    fi

    sleep 1
    elapsed=$((elapsed + 1))
done

cleanup() {
    rm -f "$owner_file"
    rmdir "$lock_dir" 2>/dev/null || true
}
trap cleanup EXIT HUP INT TERM

"$@"
