#!/usr/bin/env bash
# Provision native AFL++ GCC-plugin instrumentation for the pkt.systems lifecycle.
set -euo pipefail

version=5.02c
revision=1
archive_name="AFLplusplus-${version}.tar.gz"
archive_sha256=118415843e5d289d63bd6d8f2252c18212978f15ac9e86acbbc75766cd45acde
repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
bootlin="$repo_root/scripts/cpkt-toolchains.sh"
die() { printf 'cpkt-aflpp: %s\n' "$*" >&2; exit 1; }
cache() {
  if [[ -n "${CPKT_TOOLCHAIN_CACHE:-}" ]]; then printf '%s\n' "$CPKT_TOOLCHAIN_CACHE"
  elif [[ -n "${XDG_CACHE_HOME:-}" ]]; then printf '%s/c.pkt.systems/toolchains\n' "$XDG_CACHE_HOME"
  elif [[ -n "${HOME:-}" ]]; then printf '%s/.cache/c.pkt.systems/toolchains\n' "$HOME"
  else die 'HOME, XDG_CACHE_HOME, or CPKT_TOOLCHAIN_CACHE is required'; fi
}
root() { printf '%s/roots/aflplusplus-%s-x86_64-linux-gnu\n' "$(cache)" "$version"; }
value() { sed -n "s/^$1=//p" <<<"$2" | tail -1; }
ready() {
  local resolved_root=$1
  [[ -x "$resolved_root/bin/afl-fuzz" && -x "$resolved_root/bin/cpkt-afl-gcc" && -x "$resolved_root/bin/cpkt-afl-g++" &&
     -f "$resolved_root/lib/afl/afl-gcc-pass.so" && -f "$resolved_root/lib/afl/afl-compiler-rt.o" &&
     -f "$resolved_root/.cpkt-aflpp-revision-$revision" ]]
}

ensure() {
  [[ "$(uname -s)" = Linux ]] || die 'AFL++ GCC-plugin fuzzing is native Linux-only'
  case "$(uname -m)" in x86_64|amd64) ;; *) die "native x86_64 Linux is required; no cross, emulator, or QEMU runner is supported";; esac
  local resolved_root cache_root archive description cc cxx bootlin_root tmp source download lock_path lock_fd
  resolved_root=$(root); cache_root=$(cache); archive="$cache_root/archives/$archive_name"
  ready "$resolved_root" && return
  [[ -x "$bootlin" ]] || die "Bootlin resolver missing: $bootlin"
  command -v flock >/dev/null 2>&1 || die 'flock is required to provision shared AFL++ safely'
  lock_path="$cache_root/locks/aflplusplus-${version}-x86_64-linux-gnu.lock"
  mkdir -p "$cache_root/archives" "$cache_root/roots" "$(dirname "$lock_path")"
  exec {lock_fd}>"$lock_path"
  flock "$lock_fd"
  if ready "$resolved_root"; then
    flock -u "$lock_fd"
    exec {lock_fd}>&-
    return
  fi
  "$bootlin" ensure x86_64-linux-gnu >/dev/null
  description=$("$bootlin" discover x86_64-linux-gnu)
  cc=$(value cc "$description"); cxx=$(value cxx "$description"); bootlin_root=$(value root "$description")
  [[ -x "$cc" && -x "$cxx" && -f "$bootlin_root/include/gmp.h" ]] || die 'Bootlin GCC plugin headers are incomplete'
  if ! [[ -f "$archive" ]] || ! printf '%s  %s\n' "$archive_sha256" "$archive" | sha256sum -c - >/dev/null 2>&1; then
    rm -f "$archive"; download="$archive.tmp.$$"
    if command -v curl >/dev/null; then
      curl -fL --retry 3 --connect-timeout 20 -o "$download" "https://github.com/AFLplusplus/AFLplusplus/archive/refs/tags/v${version}.tar.gz" || { rm -f "$download"; die 'AFL++ download failed'; }
    elif command -v wget >/dev/null; then
      wget -O "$download" "https://github.com/AFLplusplus/AFLplusplus/archive/refs/tags/v${version}.tar.gz" || { rm -f "$download"; die 'AFL++ download failed'; }
    else die 'curl or wget is required to download AFL++'; fi
    printf '%s  %s\n' "$archive_sha256" "$download" | sha256sum -c - >/dev/null || { rm -f "$download"; die 'AFL++ checksum mismatch'; }
    mv "$download" "$archive"
  fi
  tmp="$cache_root/.aflplusplus.$$"; trap 'rm -rf "${tmp:-}"' EXIT HUP INT TERM
  mkdir -p "$tmp/extract" "$tmp/root/bin" "$tmp/root/lib/afl"
  tar -xzf "$archive" -C "$tmp/extract"; source="$tmp/extract/AFLplusplus-$version"
  [[ -d "$source" ]] || die "unexpected archive layout: $archive_name"
  (
    cd "$source"; local helper="$resolved_root/lib/afl"
    make -j1 NO_PYTHON=1 CC="$cc" CXX="$cxx" PREFIX="$tmp/root" HELPER_PATH="$helper" BIN_PATH="$tmp/root/bin" afl-fuzz afl-showmap afl-tmin afl-gotcpu afl-analyze afl-cmin
    "$cc" -O3 -funroll-loops -fPIC -Wall -g -Iinclude -Iinstrumentation -DAFL_PATH=\"$helper\" -DBIN_PATH=\"$resolved_root/bin\" -DLLVM_BINDIR=\"\" -DVERSION=\"++$version\" -DLLVM_LIBDIR=\"\" -DLLVM_VERSION=\"\" -DAFL_CLANG_FLTO=\"\" -DAFL_REAL_LD=\"\" -DAFL_CLANG_LDPATH=\"\" -DAFL_CLANG_FUSELD=\"\" -DCLANG_BIN=\"$cc\" -DCLANGPP_BIN=\"$cxx\" -DUSE_BINDIR=1 -Wno-unused-function -Wno-deprecated -c src/afl-common.c -o instrumentation/afl-common.o
    "$cc" -O3 -funroll-loops -fPIC -Wall -g -Iinclude -Iinstrumentation -DAFL_PATH=\"$helper\" -DBIN_PATH=\"$resolved_root/bin\" -DLLVM_BINDIR=\"\" -DVERSION=\"++$version\" -DLLVM_LIBDIR=\"\" -DLLVM_VERSION=\"\" -DAFL_CLANG_FLTO=\"\" -DAFL_REAL_LD=\"\" -DAFL_CLANG_LDPATH=\"\" -DAFL_CLANG_FUSELD=\"\" -DCLANG_BIN=\"$cc\" -DCLANGPP_BIN=\"$cxx\" -DUSE_BINDIR=1 -Wno-unused-function -Wno-deprecated -DAFL_INCLUDE_PATH=\"$resolved_root/include/afl\" src/afl-cc.c instrumentation/afl-common.o -o afl-cc -DLLVM_MINOR=0 -DLLVM_MAJOR=0 -DCFLAGS_OPT=\"\" -lm
    ln -sf afl-cc afl-gcc-fast; ln -sf afl-cc afl-g++-fast
    make -j1 -f GNUmakefile.gcc_plugin CC="$cc" CXX="$cxx" PREFIX="$tmp/root" HELPER_PATH="$helper" BIN_PATH="$tmp/root/bin" CXXFLAGS="-O3 -g -funroll-loops -I$bootlin_root/include" LDFLAGS="-L$bootlin_root/lib -Wl,-rpath,$bootlin_root/lib"
    install -m755 afl-fuzz afl-showmap afl-tmin afl-gotcpu afl-analyze afl-cmin afl-cc "$tmp/root/bin/"
    ln -sf afl-cc "$tmp/root/bin/afl-gcc-fast"; ln -sf afl-cc "$tmp/root/bin/afl-g++-fast"
    install -m755 afl-gcc-pass.so afl-gcc-cmplog-pass.so afl-gcc-cmptrs-pass.so "$tmp/root/lib/afl/"
    install -m644 afl-compiler-rt.o dynamic_list.txt "$tmp/root/lib/afl/"
  )
  printf '#!/usr/bin/env bash\nexport AFL_PATH=%q\nexport AFL_CC=%q\nexec %q "$@"\n' "$resolved_root/lib/afl" "$cc" "$resolved_root/bin/afl-gcc-fast" > "$tmp/root/bin/cpkt-afl-gcc"
  printf '#!/usr/bin/env bash\nexport AFL_PATH=%q\nexport AFL_CC=%q\nexport AFL_CXX=%q\nexec %q "$@"\n' "$resolved_root/lib/afl" "$cc" "$cxx" "$resolved_root/bin/afl-g++-fast" > "$tmp/root/bin/cpkt-afl-g++"
  chmod +x "$tmp/root/bin/cpkt-afl-gcc" "$tmp/root/bin/cpkt-afl-g++"; touch "$tmp/root/.cpkt-aflpp-revision-$revision"
  ready "$tmp/root" || die 'incomplete AFL++ build'; rm -rf "$resolved_root"; mv "$tmp/root" "$resolved_root"; trap - EXIT HUP INT TERM; rm -rf "$tmp"
  flock -u "$lock_fd"
  exec {lock_fd}>&-
}

report() { ensure; local resolved_root=$(root); printf 'version=%s\ncache=%s\nsource=aflplusplus\nroot=%s\nafl_fuzz=%s\nafl_showmap=%s\ncc=%s\ncxx=%s\nhelper=%s\n' "$version" "$(cache)" "$resolved_root" "$resolved_root/bin/afl-fuzz" "$resolved_root/bin/afl-showmap" "$resolved_root/bin/cpkt-afl-gcc" "$resolved_root/bin/cpkt-afl-g++" "$resolved_root/lib/afl"; }
env_out() { local description cc cxx resolved_root; ensure; description=$("$bootlin" discover x86_64-linux-gnu); cc=$(value cc "$description"); cxx=$(value cxx "$description"); resolved_root=$(root); printf 'export CPKT_AFLPP_ROOT=%q\nexport AFL_PATH=%q\nexport AFL_CC=%q\nexport AFL_CXX=%q\nexport CC=%q\nexport CXX=%q\n' "$resolved_root" "$resolved_root/lib/afl" "$cc" "$cxx" "$resolved_root/bin/cpkt-afl-gcc" "$resolved_root/bin/cpkt-afl-g++"; }
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  case "${1:-}" in
    ensure) [[ $# -eq 1 ]] || die 'usage: cpkt-aflpp.sh ensure'; ensure;;
    discover) [[ $# -eq 1 ]] || die 'usage: cpkt-aflpp.sh discover'; report;;
    env) [[ $# -eq 1 ]] || die 'usage: cpkt-aflpp.sh env'; env_out;;
    *) die 'usage: cpkt-aflpp.sh {ensure|discover|env}';;
  esac
fi
