# examples

`examples/example.c` is the main usage walkthrough for `libpslog`.

It is intentionally direct: each section creates and uses the public API without wrapping logger creation behind helper functions.

Build it in normal library mode:

```sh
cmake --preset host
cmake --build --preset host
cd examples
"$(sed -n 's/^CMAKE_C_COMPILER:FILEPATH=//p' ../build/host/CMakeCache.txt)" -I../build/host/generated/include -I../include \
  -o example example.c ../build/host/libpslog.a -pthread
./example
```

The compiler path above comes from the configured host CMake cache. Do not
substitute an ambient `cc`, `gcc`, or `clang`: Linux builds and consumer smoke
tests use the pinned Bootlin collection selected for that target.

Build it in single-header mode:

```sh
cmake --preset host
cmake --build ../build/host --target package-single-header
cd examples
"$(sed -n 's/^CMAKE_C_COMPILER:FILEPATH=//p' ../build/host/CMakeCache.txt)" -DPSLOG_EXAMPLE_SINGLE_HEADER=1 \
  -I../build/host/generated/include \
  -o example example.c -pthread
./example
```

It demonstrates:

- console and JSON loggers
- `log()` and level-specific methods
- `with()`
- `withf()`
- `with_level()` and `with_level_field()`
- typed fields
- `infof` / `kvfmt`
- `pslog_new_from_env()`
- palette iteration
- adaptive color by default: ANSI is emitted only when stdout is a tty unless forced
