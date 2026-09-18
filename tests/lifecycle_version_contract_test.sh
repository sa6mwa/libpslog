#!/usr/bin/env bash
set -euo pipefail

repo_root=${1:?repository root required}
fixture_root=""

cleanup() {
  if [ -n "$fixture_root" ]; then
    rm -rf -- "$fixture_root"
  fi
}
trap cleanup EXIT HUP INT TERM

mkdir -p "$repo_root/build"
fixture_root="$(mktemp -d "$repo_root/build/lifecycle-version-contract-test.XXXXXX")"
mkdir -p "$fixture_root/scripts" "$fixture_root/lua/scripts" \
  "$fixture_root/cmake" "$fixture_root/tests"
cp "$repo_root/scripts/lifecycle_version_contract.sh" "$fixture_root/scripts/"
cp "$repo_root/lua/scripts/release_version.sh" "$fixture_root/lua/scripts/"
cp "$repo_root/cmake/pslog_version.cmake" "$fixture_root/cmake/"
cp "$repo_root/tests/version_resolution_probe.cmake" "$fixture_root/tests/"
cat >"$fixture_root/Makefile" <<'EOF'
print-release-version:
	@./lua/scripts/release_version.sh
EOF

git -C "$fixture_root" -c init.defaultBranch=trunk init >/dev/null
git -C "$fixture_root" config user.email test@example.com
git -C "$fixture_root" config user.name "libpslog lifecycle test"
git -C "$fixture_root" add Makefile scripts lua cmake tests
git -C "$fixture_root" commit -m "version fixture" >/dev/null

untagged_output=$("$fixture_root/scripts/lifecycle_version_contract.sh")
[[ "$untagged_output" == *"version 99.99.99 verified"* ]] || {
  printf 'untagged lifecycle contract did not use its reserved test tag\n' >&2
  exit 1
}
if [ ! -d "$fixture_root/build" ]; then
  printf 'lifecycle contract did not create a clean fixture build directory\n' >&2
  exit 1
fi
if git -C "$fixture_root" rev-parse --verify --quiet refs/tags/v99.99.99 >/dev/null; then
  printf 'untagged lifecycle contract left its reserved test tag behind\n' >&2
  exit 1
fi

git -C "$fixture_root" -c tag.gpgSign=false tag v98.98.98
if [ "$(git -C "$fixture_root" cat-file -t refs/tags/v98.98.98)" != commit ]; then
  printf 'test release tag was not lightweight\n' >&2
  exit 1
fi

tagged_output=$("$fixture_root/scripts/lifecycle_version_contract.sh")
[[ "$tagged_output" == *"version 98.98.98 verified"* ]] || {
  printf 'tagged lifecycle contract did not select the exact lightweight tag\n' >&2
  exit 1
}
if git -C "$fixture_root" rev-parse --verify --quiet refs/tags/v99.99.99 >/dev/null; then
  printf 'tagged lifecycle contract created its reserved test tag\n' >&2
  exit 1
fi

git -C "$fixture_root" tag -d v98.98.98 >/dev/null
git -C "$fixture_root" -c tag.gpgSign=false tag -a v97.97.97 -m "annotated lifecycle version test"
if annotated_output=$("$fixture_root/scripts/lifecycle_version_contract.sh" 2>&1); then
  printf 'lifecycle contract accepted an annotated release tag\n' >&2
  exit 1
fi
[[ "$annotated_output" == *"must be lightweight"* ]] || {
  printf 'annotated release tag failure was not actionable: %s\n' "$annotated_output" >&2
  exit 1
}
