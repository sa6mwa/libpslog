#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

/* Check actual mappings, not just the linker's intended search order. */
int main(int argc, char **argv) {
  FILE *maps;
  char line[8192];
  int saw_libc = 0;
  int status;
  pid_t child;
  const char *root = PSLOG_RUNTIME_ROOT;
  (void)argv;
  maps = fopen("/proc/self/maps", "r");
  if (maps == NULL)
    return 1;
  while (fgets(line, sizeof(line), maps) != NULL) {
    char *path = strchr(line, '/');
    if (path == NULL || strstr(line, "r-xp") == NULL ||
        strstr(path, ".so") == NULL)
      continue;
    if (strncmp(path, root, strlen(root)) != 0 || path[strlen(root)] != '/') {
      fprintf(stderr, "Unexpected runtime mapping: %s", path);
      fclose(maps);
      return 2;
    }
    if (strstr(path, "libc.so") != NULL)
      saw_libc = 1;
  }
  fclose(maps);
  if (!saw_libc)
    return 3;
  if (argc > 1)
    return 0;
  child = fork();
  if (child == 0) {
    execl("/proc/self/exe", "runtime_probe", "child", (char *)NULL);
    _exit(4);
  }
  if (child < 0 || waitpid(child, &status, 0) != child || !WIFEXITED(status) ||
      WEXITSTATUS(status) != 0)
    return 5;
  /* Host children keep their own interpreter and no injected library paths. */
  return system("/bin/sh -c 'test -z \"${LD_LIBRARY_PATH:-}\" && test -z "
                "\"${LD_PRELOAD:-}\"'") == 0
             ? 0
             : 6;
}
