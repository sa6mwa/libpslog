#!/usr/bin/env bash
set -eu

tag="$(git describe --tags --exact-match HEAD 2>/dev/null || true)"

case "$tag" in
  v*)
    printf '%s\n' "${tag#v}"
    ;;
  *)
    printf '0.0.0\n'
    ;;
esac
