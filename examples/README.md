# examples

`examples/example.c` is the main usage walkthrough for `libpslog`.

It is intentionally direct: each section creates and uses the public API without wrapping logger creation behind helper functions.

Build it after the release library exists:

```sh
cmake --preset host
cmake --build --preset host
cd examples
cc -I../build/host/generated/include -I../include \
  -o example example.c ../build/host/libpslog.a
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
