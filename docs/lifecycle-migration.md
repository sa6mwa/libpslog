# Lifecycle migration ledger

This ledger is temporary and records the hard cutover to the pkt.systems CMake
lifecycle. Remove it only after the repository has converged.

| Old behavior | Lifecycle surface | Preserved behavior | Verification |
| --- | --- | --- | --- |
| Distro and `$HOME/.local/cross` Linux compiler paths | Pinned toolchains | Same six Linux target IDs and Darwin osxcross support | Resolver syntax and resolver contract tests |
| `linux-gnu-release` / `linux-musl-release` preset names | Standard CMake release presets | x86_64 GNU and musl release artifacts | CMake preset listing and release matrix |
| Host `cc` in Go dataset generation | Lifecycle-owned host compiler | Generated benchmark dataset | Host build then gobencher tests |
| Package build presets per target | Direct CMake target invocation from configured release preset | Archive, source, single-header, checksum, and privacy gates | `make package-verify` and `make release-matrix` |

Remaining work is limited to verification and any toolchain-specific package
failures surfaced by the new matrix; no legacy Linux toolchain path remains.
