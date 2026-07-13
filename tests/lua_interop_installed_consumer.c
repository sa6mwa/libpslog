#include <pslog_lua.h>

int main(void) { return pslog_lua_is_logger(NULL, 0) == 0 ? 0 : 1; }
