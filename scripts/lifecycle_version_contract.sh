#!/usr/bin/env bash
set -euo pipefail

repo_root="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
reserved_tag="v99.99.99"
scratch=""

cleanup() {
  if [ -n "$scratch" ]; then
    rm -rf -- "$scratch"
  fi
  git -C "$repo_root" tag -d "$reserved_tag" >/dev/null 2>&1 || true
}
trap cleanup EXIT HUP INT TERM

if ! git -C "$repo_root" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  printf 'lifecycle-version-contract: a Git worktree is required\n' >&2
  exit 1
fi

# This is reserved test-only state. Delete an interrupted previous run before
# inspecting release tags so it cannot affect the contract.
git -C "$repo_root" tag -d "$reserved_tag" >/dev/null 2>&1 || true
head_commit="$(git -C "$repo_root" rev-parse HEAD)"
exact_tags=()

while IFS= read -r tag; do
  [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]] || continue
  tag_commit="$(git -C "$repo_root" rev-parse "refs/tags/$tag^{}" 2>/dev/null || true)"
  [ "$tag_commit" = "$head_commit" ] || continue
  tag_type="$(git -C "$repo_root" cat-file -t "refs/tags/$tag")"
  if [ "$tag_type" != commit ]; then
    printf 'lifecycle-version-contract: exact release tag %s must be lightweight\n' "$tag" >&2
    exit 1
  fi
  exact_tags+=("$tag")
done < <(git -C "$repo_root" for-each-ref --format='%(refname:strip=2)' refs/tags)

if [ "${#exact_tags[@]}" -gt 1 ]; then
  printf 'lifecycle-version-contract: multiple exact lightweight release tags point at HEAD: %s\n' "${exact_tags[*]}" >&2
  exit 1
fi

if [ "${#exact_tags[@]}" -eq 1 ]; then
  expected_version="${exact_tags[0]#v}"
else
  git -C "$repo_root" -c tag.gpgSign=false tag "$reserved_tag"
  if [ "$(git -C "$repo_root" cat-file -t "refs/tags/$reserved_tag")" != commit ]; then
    printf 'lifecycle-version-contract: reserved tag is not lightweight\n' >&2
    exit 1
  fi
  expected_version="${reserved_tag#v}"
fi

make_version="$(make -C "$repo_root" --no-print-directory print-release-version)"
if [ "$make_version" != "$expected_version" ]; then
  printf 'lifecycle-version-contract: Make resolved %s, expected %s\n' "$make_version" "$expected_version" >&2
  exit 1
fi

mkdir -p "$repo_root/build"
scratch="$(mktemp -d "$repo_root/build/lifecycle-version-contract.XXXXXX")"
cmake -DPSLOG_ROOT="$repo_root" \
  -DPSLOG_VERSION_SOURCE_DIR="$repo_root" \
  -DPSLOG_VERSION_PROBE_OUTPUT="$scratch/version.txt" \
  -P "$repo_root/tests/version_resolution_probe.cmake" >/dev/null
cmake_version="$(cat "$scratch/version.txt" | cut -d '|' -f 1)"
if [ "$cmake_version" != "$expected_version" ]; then
  printf 'lifecycle-version-contract: CMake resolved %s, expected %s\n' "$cmake_version" "$expected_version" >&2
  exit 1
fi

if [ "${#exact_tags[@]}" -eq 0 ]; then
  git -C "$repo_root" tag -d "$reserved_tag" >/dev/null
  make_version="$(make -C "$repo_root" --no-print-directory print-release-version)"
  if [ "$make_version" != 0.0.0 ]; then
    printf 'lifecycle-version-contract: untagged Make version must be 0.0.0, got %s\n' "$make_version" >&2
    exit 1
  fi
fi

printf 'lifecycle-version-contract: version %s verified\n' "$expected_version"
