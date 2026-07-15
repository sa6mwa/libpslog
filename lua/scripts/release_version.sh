#!/usr/bin/env bash
set -eu

tag="$(git describe --tags --exact-match HEAD 2>/dev/null || true)"

case "$tag" in
  v*)
    printf '%s\n' "${tag#v}"
    ;;
  *)
    if [ ! -d .git ] && [ -f VERSION ]; then
      version=$(sed -n '1p' VERSION)
      case "$version" in
        [0-9]*.[0-9]*.[0-9]*)
          printf '%s\n' "$version"
          ;;
        *)
          printf '0.0.0\n'
          ;;
      esac
    else
      printf '0.0.0\n'
    fi
    ;;
esac
