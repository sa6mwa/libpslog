#include "pslog.h"

#include <ctype.h>
#include <string.h>

const pslog_palette pslog_builtin_palette_default = {
    "\x1b[36m",   "\x1b[1;34m", "\x1b[35m",   "\x1b[33m",
    "\x1b[90m",   "\x1b[34m",   "\x1b[32m",   "\x1b[1;32m",
    "\x1b[1;33m", "\x1b[1;31m", "\x1b[1;31m", "\x1b[1;31m",
    "\x1b[90m",   "\x1b[1m",    "\x1b[36m",   "\x1b[0m"};

const pslog_palette pslog_builtin_palette_outrun_electric = {
    "\x1b[38;5;201m", "\x1b[38;5;81m",    "\x1b[38;5;99m",  "\x1b[38;5;69m",
    "\x1b[38;5;60m",  "\x1b[38;5;33m",    "\x1b[38;5;39m",  "\x1b[38;5;45m",
    "\x1b[38;5;129m", "\x1b[38;5;205m",   "\x1b[38;5;206m", "\x1b[38;5;213m",
    "\x1b[38;5;117m", "\x1b[1;38;5;219m", "\x1b[38;5;33m",  "\x1b[0m"};

const pslog_palette pslog_builtin_palette_iosvkem = {
    "\x1b[38;5;222m", "\x1b[38;5;216m",   "\x1b[38;5;109m", "\x1b[38;5;151m",
    "\x1b[38;5;244m", "\x1b[38;5;66m",    "\x1b[38;5;72m",  "\x1b[38;5;114m",
    "\x1b[38;5;208m", "\x1b[38;5;203m",   "\x1b[38;5;197m", "\x1b[38;5;199m",
    "\x1b[38;5;242m", "\x1b[1;38;5;223m", "\x1b[38;5;67m",  "\x1b[0m"};

const pslog_palette pslog_builtin_palette_gruvbox = {
    "\x1b[38;5;214m", "\x1b[38;5;178m",   "\x1b[38;5;108m", "\x1b[38;5;142m",
    "\x1b[38;5;101m", "\x1b[38;5;66m",    "\x1b[38;5;72m",  "\x1b[38;5;107m",
    "\x1b[38;5;208m", "\x1b[38;5;167m",   "\x1b[38;5;160m", "\x1b[38;5;161m",
    "\x1b[38;5;137m", "\x1b[1;38;5;221m", "\x1b[38;5;172m", "\x1b[0m"};

const pslog_palette pslog_builtin_palette_dracula = {
    "\x1b[38;5;219m", "\x1b[38;5;141m",   "\x1b[38;5;111m", "\x1b[38;5;81m",
    "\x1b[38;5;240m", "\x1b[38;5;60m",    "\x1b[38;5;98m",  "\x1b[38;5;117m",
    "\x1b[38;5;219m", "\x1b[38;5;204m",   "\x1b[38;5;198m", "\x1b[38;5;199m",
    "\x1b[38;5;95m",  "\x1b[1;38;5;225m", "\x1b[38;5;147m", "\x1b[0m"};

const pslog_palette pslog_builtin_palette_nord = {
    "\x1b[38;5;153m", "\x1b[38;5;152m",   "\x1b[38;5;109m", "\x1b[38;5;115m",
    "\x1b[38;5;245m", "\x1b[38;5;67m",    "\x1b[38;5;74m",  "\x1b[38;5;117m",
    "\x1b[38;5;179m", "\x1b[38;5;210m",   "\x1b[38;5;204m", "\x1b[38;5;205m",
    "\x1b[38;5;109m", "\x1b[1;38;5;195m", "\x1b[38;5;110m", "\x1b[0m"};

const pslog_palette pslog_builtin_palette_tokyo_night = {
    "\x1b[38;5;69m",  "\x1b[38;5;110m",   "\x1b[38;5;176m", "\x1b[38;5;117m",
    "\x1b[38;5;244m", "\x1b[38;5;63m",    "\x1b[38;5;67m",  "\x1b[38;5;111m",
    "\x1b[38;5;173m", "\x1b[38;5;210m",   "\x1b[38;5;205m", "\x1b[38;5;219m",
    "\x1b[38;5;109m", "\x1b[1;38;5;218m", "\x1b[38;5;74m",  "\x1b[0m"};

const pslog_palette pslog_builtin_palette_solarized_nightfall = {
    "\x1b[38;5;37m",  "\x1b[38;5;86m",    "\x1b[38;5;61m",  "\x1b[38;5;65m",
    "\x1b[38;5;239m", "\x1b[38;5;24m",    "\x1b[38;5;30m",  "\x1b[38;5;36m",
    "\x1b[38;5;136m", "\x1b[38;5;160m",   "\x1b[38;5;166m", "\x1b[38;5;161m",
    "\x1b[38;5;244m", "\x1b[1;38;5;230m", "\x1b[38;5;33m",  "\x1b[0m"};

const pslog_palette pslog_builtin_palette_catppuccin_mocha = {
    "\x1b[38;5;217m", "\x1b[38;5;183m",   "\x1b[38;5;147m", "\x1b[38;5;152m",
    "\x1b[38;5;244m", "\x1b[38;5;104m",   "\x1b[38;5;109m", "\x1b[38;5;150m",
    "\x1b[38;5;216m", "\x1b[38;5;211m",   "\x1b[38;5;205m", "\x1b[38;5;204m",
    "\x1b[38;5;110m", "\x1b[1;38;5;223m", "\x1b[38;5;182m", "\x1b[0m"};

const pslog_palette pslog_builtin_palette_gruvbox_light = {
    "\x1b[38;5;130m", "\x1b[38;5;108m",   "\x1b[38;5;66m",  "\x1b[38;5;142m",
    "\x1b[38;5;180m", "\x1b[38;5;109m",   "\x1b[38;5;114m", "\x1b[38;5;73m",
    "\x1b[38;5;173m", "\x1b[38;5;167m",   "\x1b[38;5;161m", "\x1b[38;5;125m",
    "\x1b[38;5;180m", "\x1b[1;38;5;223m", "\x1b[38;5;136m", "\x1b[0m"};

const pslog_palette pslog_builtin_palette_monokai_vibrant = {
    "\x1b[38;5;229m", "\x1b[38;5;121m",   "\x1b[38;5;198m", "\x1b[38;5;118m",
    "\x1b[38;5;59m",  "\x1b[38;5;104m",   "\x1b[38;5;114m", "\x1b[38;5;121m",
    "\x1b[38;5;215m", "\x1b[38;5;197m",   "\x1b[38;5;161m", "\x1b[38;5;201m",
    "\x1b[38;5;103m", "\x1b[1;38;5;229m", "\x1b[38;5;141m", "\x1b[0m"};

const pslog_palette pslog_builtin_palette_one_dark_aurora = {
    "\x1b[38;5;110m", "\x1b[38;5;147m",   "\x1b[38;5;141m", "\x1b[38;5;115m",
    "\x1b[38;5;59m",  "\x1b[38;5;24m",    "\x1b[38;5;31m",  "\x1b[38;5;38m",
    "\x1b[38;5;178m", "\x1b[38;5;203m",   "\x1b[38;5;197m", "\x1b[38;5;199m",
    "\x1b[38;5;109m", "\x1b[1;38;5;189m", "\x1b[38;5;75m",  "\x1b[0m"};

const pslog_palette pslog_builtin_palette_synthwave_84 = {
    "\x1b[38;5;198m", "\x1b[38;5;51m",    "\x1b[38;5;207m", "\x1b[38;5;219m",
    "\x1b[38;5;102m", "\x1b[38;5;63m",    "\x1b[38;5;69m",  "\x1b[38;5;81m",
    "\x1b[38;5;220m", "\x1b[38;5;205m",   "\x1b[38;5;200m", "\x1b[38;5;201m",
    "\x1b[38;5;69m",  "\x1b[1;38;5;219m", "\x1b[38;5;45m",  "\x1b[0m"};

const pslog_palette pslog_builtin_palette_kanagawa = {
    "\x1b[1;38;5;215m", "\x1b[1;38;5;179m", "\x1b[38;5;150m",
    "\x1b[38;5;109m",   "\x1b[38;5;240m",   "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;216m", "\x1b[38;5;110m",
    "\x1b[1;38;5;181m", "\x1b[1;38;5;181m", "\x1b[1;38;5;181m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;223m", "\x1b[38;5;173m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_rose_pine = {
    "\x1b[1;38;5;217m", "\x1b[1;38;5;180m", "\x1b[38;5;146m",
    "\x1b[38;5;110m",   "\x1b[38;5;244m",   "\x1b[38;5;104m",
    "\x1b[38;5;110m",   "\x1b[1;38;5;217m", "\x1b[38;5;117m",
    "\x1b[1;38;5;210m", "\x1b[1;38;5;210m", "\x1b[1;38;5;210m",
    "\x1b[38;5;110m",   "\x1b[1;38;5;223m", "\x1b[38;5;146m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_rose_pine_dawn = {
    "\x1b[1;38;5;173m", "\x1b[1;38;5;137m", "\x1b[38;5;108m",
    "\x1b[38;5;66m",    "\x1b[38;5;244m",   "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;173m", "\x1b[38;5;109m",
    "\x1b[1;38;5;131m", "\x1b[1;38;5;131m", "\x1b[1;38;5;131m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;180m", "\x1b[38;5;137m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_everforest = {
    "\x1b[1;38;5;143m", "\x1b[1;38;5;108m", "\x1b[38;5;108m",
    "\x1b[38;5;142m",   "\x1b[38;5;239m",   "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;143m", "\x1b[38;5;109m",
    "\x1b[1;38;5;179m", "\x1b[1;38;5;179m", "\x1b[1;38;5;179m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;223m", "\x1b[38;5;108m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_everforest_light = {
    "\x1b[1;38;5;136m", "\x1b[1;38;5;108m", "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[38;5;243m",   "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;136m", "\x1b[38;5;137m",
    "\x1b[1;38;5;173m", "\x1b[1;38;5;173m", "\x1b[1;38;5;173m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;130m", "\x1b[38;5;108m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_night_owl = {
    "\x1b[1;38;5;111m", "\x1b[1;38;5;75m",  "\x1b[38;5;110m",
    "\x1b[38;5;81m",    "\x1b[38;5;244m",   "\x1b[38;5;67m",
    "\x1b[38;5;74m",    "\x1b[1;38;5;117m", "\x1b[38;5;179m",
    "\x1b[1;38;5;204m", "\x1b[1;38;5;204m", "\x1b[1;38;5;204m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;117m", "\x1b[38;5;75m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_ayu_mirage = {
    "\x1b[1;38;5;214m", "\x1b[1;38;5;179m", "\x1b[38;5;110m",
    "\x1b[38;5;109m",   "\x1b[38;5;239m",   "\x1b[38;5;66m",
    "\x1b[38;5;110m",   "\x1b[1;38;5;214m", "\x1b[38;5;173m",
    "\x1b[1;38;5;173m", "\x1b[1;38;5;173m", "\x1b[1;38;5;173m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;215m", "\x1b[38;5;179m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_ayu_light = {
    "\x1b[1;38;5;136m", "\x1b[1;38;5;109m", "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[38;5;243m",   "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;130m", "\x1b[38;5;173m",
    "\x1b[1;38;5;173m", "\x1b[1;38;5;173m", "\x1b[1;38;5;173m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;130m", "\x1b[38;5;109m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_one_light = {
    "\x1b[1;38;5;33m",  "\x1b[1;38;5;97m",  "\x1b[38;5;67m",
    "\x1b[38;5;109m",   "\x1b[38;5;243m",   "\x1b[38;5;66m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;130m", "\x1b[38;5;137m",
    "\x1b[1;38;5;167m", "\x1b[1;38;5;167m", "\x1b[1;38;5;167m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;130m", "\x1b[38;5;67m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_one_dark = {
    "\x1b[1;38;5;110m", "\x1b[1;38;5;147m", "\x1b[38;5;141m",
    "\x1b[38;5;115m",   "\x1b[38;5;239m",   "\x1b[38;5;67m",
    "\x1b[38;5;75m",    "\x1b[1;38;5;75m",  "\x1b[38;5;178m",
    "\x1b[1;38;5;203m", "\x1b[1;38;5;203m", "\x1b[1;38;5;203m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;189m", "\x1b[38;5;75m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_solarized_light = {
    "\x1b[1;38;5;64m",  "\x1b[1;38;5;37m",  "\x1b[38;5;33m",
    "\x1b[38;5;130m",   "\x1b[38;5;244m",   "\x1b[38;5;33m",
    "\x1b[38;5;37m",    "\x1b[1;38;5;136m", "\x1b[38;5;166m",
    "\x1b[1;38;5;166m", "\x1b[1;38;5;166m", "\x1b[1;38;5;166m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;136m", "\x1b[38;5;64m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_solarized_dark = {
    "\x1b[1;38;5;64m",  "\x1b[1;38;5;37m",  "\x1b[38;5;33m",
    "\x1b[38;5;130m",   "\x1b[38;5;240m",   "\x1b[38;5;33m",
    "\x1b[38;5;37m",    "\x1b[1;38;5;136m", "\x1b[38;5;166m",
    "\x1b[1;38;5;166m", "\x1b[1;38;5;166m", "\x1b[1;38;5;166m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;230m", "\x1b[38;5;64m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_github_light = {
    "\x1b[1;38;5;27m",  "\x1b[1;38;5;61m",  "\x1b[38;5;67m",
    "\x1b[38;5;109m",   "\x1b[38;5;244m",   "\x1b[38;5;67m",
    "\x1b[38;5;61m",    "\x1b[1;38;5;24m",  "\x1b[38;5;130m",
    "\x1b[1;38;5;130m", "\x1b[1;38;5;130m", "\x1b[1;38;5;130m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;24m",  "\x1b[38;5;61m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_github_dark = {
    "\x1b[1;38;5;111m", "\x1b[1;38;5;109m", "\x1b[38;5;110m",
    "\x1b[38;5;109m",   "\x1b[38;5;239m",   "\x1b[38;5;67m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;111m", "\x1b[38;5;179m",
    "\x1b[1;38;5;203m", "\x1b[1;38;5;203m", "\x1b[1;38;5;203m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;189m", "\x1b[38;5;109m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_papercolor_light = {
    "\x1b[1;38;5;33m",  "\x1b[1;38;5;64m",  "\x1b[38;5;67m",
    "\x1b[38;5;109m",   "\x1b[38;5;244m",   "\x1b[38;5;67m",
    "\x1b[38;5;61m",    "\x1b[1;38;5;130m", "\x1b[38;5;136m",
    "\x1b[1;38;5;166m", "\x1b[1;38;5;166m", "\x1b[1;38;5;166m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;130m", "\x1b[38;5;64m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_papercolor_dark = {
    "\x1b[1;38;5;109m", "\x1b[1;38;5;110m", "\x1b[38;5;110m",
    "\x1b[38;5;109m",   "\x1b[38;5;239m",   "\x1b[38;5;67m",
    "\x1b[38;5;110m",   "\x1b[1;38;5;109m", "\x1b[38;5;179m",
    "\x1b[1;38;5;203m", "\x1b[1;38;5;203m", "\x1b[1;38;5;203m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;223m", "\x1b[38;5;110m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_oceanic_next = {
    "\x1b[1;38;5;110m", "\x1b[1;38;5;75m",  "\x1b[38;5;109m",
    "\x1b[38;5;110m",   "\x1b[38;5;239m",   "\x1b[38;5;67m",
    "\x1b[38;5;75m",    "\x1b[1;38;5;110m", "\x1b[38;5;179m",
    "\x1b[1;38;5;203m", "\x1b[1;38;5;203m", "\x1b[1;38;5;203m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;189m", "\x1b[38;5;75m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_horizon = {
    "\x1b[1;38;5;214m", "\x1b[1;38;5;203m", "\x1b[38;5;209m",
    "\x1b[38;5;141m",   "\x1b[38;5;239m",   "\x1b[38;5;167m",
    "\x1b[38;5;209m",   "\x1b[1;38;5;214m", "\x1b[38;5;110m",
    "\x1b[1;38;5;203m", "\x1b[1;38;5;203m", "\x1b[1;38;5;203m",
    "\x1b[38;5;110m",   "\x1b[1;38;5;222m", "\x1b[38;5;203m",
    "\x1b[0m"};

const pslog_palette pslog_builtin_palette_palenight = {
    "\x1b[1;38;5;147m", "\x1b[1;38;5;141m", "\x1b[38;5;110m",
    "\x1b[38;5;109m",   "\x1b[38;5;239m",   "\x1b[38;5;67m",
    "\x1b[38;5;110m",   "\x1b[1;38;5;147m", "\x1b[38;5;179m",
    "\x1b[1;38;5;203m", "\x1b[1;38;5;203m", "\x1b[1;38;5;203m",
    "\x1b[38;5;109m",   "\x1b[1;38;5;189m", "\x1b[38;5;141m",
    "\x1b[0m"};

struct palette_entry {
  const char *name;
  const pslog_palette *palette;
};

static const struct palette_entry palette_entries[] = {
    {"ayu-light", &pslog_builtin_palette_ayu_light},
    {"ayu-mirage", &pslog_builtin_palette_ayu_mirage},
    {"catppuccin-mocha", &pslog_builtin_palette_catppuccin_mocha},
    {"default", &pslog_builtin_palette_default},
    {"dracula", &pslog_builtin_palette_dracula},
    {"everforest", &pslog_builtin_palette_everforest},
    {"everforest-light", &pslog_builtin_palette_everforest_light},
    {"github-dark", &pslog_builtin_palette_github_dark},
    {"github-light", &pslog_builtin_palette_github_light},
    {"gruvbox", &pslog_builtin_palette_gruvbox},
    {"gruvbox-light", &pslog_builtin_palette_gruvbox_light},
    {"horizon", &pslog_builtin_palette_horizon},
    {"iosvkem", &pslog_builtin_palette_iosvkem},
    {"kanagawa", &pslog_builtin_palette_kanagawa},
    {"monokai-vibrant", &pslog_builtin_palette_monokai_vibrant},
    {"night-owl", &pslog_builtin_palette_night_owl},
    {"nord", &pslog_builtin_palette_nord},
    {"oceanic-next", &pslog_builtin_palette_oceanic_next},
    {"one-dark", &pslog_builtin_palette_one_dark},
    {"one-dark-aurora", &pslog_builtin_palette_one_dark_aurora},
    {"one-light", &pslog_builtin_palette_one_light},
    {"outrun-electric", &pslog_builtin_palette_outrun_electric},
    {"palenight", &pslog_builtin_palette_palenight},
    {"papercolor-dark", &pslog_builtin_palette_papercolor_dark},
    {"papercolor-light", &pslog_builtin_palette_papercolor_light},
    {"rose-pine", &pslog_builtin_palette_rose_pine},
    {"rose-pine-dawn", &pslog_builtin_palette_rose_pine_dawn},
    {"solarized-dark", &pslog_builtin_palette_solarized_dark},
    {"solarized-light", &pslog_builtin_palette_solarized_light},
    {"solarized-nightfall", &pslog_builtin_palette_solarized_nightfall},
    {"synthwave-84", &pslog_builtin_palette_synthwave_84},
    {"tokyo-night", &pslog_builtin_palette_tokyo_night}};

struct alias_entry {
  const char *alias;
  const char *canonical;
};

static const struct alias_entry palette_aliases[] = {
    {"ayulight", "ayu-light"},
    {"ayumirage", "ayu-mirage"},
    {"catppuccinmocha", "catppuccin-mocha"},
    {"doom-dracula", "dracula"},
    {"doom-gruvbox", "gruvbox"},
    {"doom-iosvkem", "iosvkem"},
    {"doom-nord", "nord"},
    {"doomdracula", "dracula"},
    {"doomgruvbox", "gruvbox"},
    {"doomiosvkem", "iosvkem"},
    {"doomnord", "nord"},
    {"everforestlight", "everforest-light"},
    {"githubdark", "github-dark"},
    {"githublight", "github-light"},
    {"gruvboxlight", "gruvbox-light"},
    {"monokaivibrant", "monokai-vibrant"},
    {"nightowl", "night-owl"},
    {"oceanicnext", "oceanic-next"},
    {"onedark", "one-dark"},
    {"onedarkaurora", "one-dark-aurora"},
    {"onelight", "one-light"},
    {"outrunelectric", "outrun-electric"},
    {"papercolordark", "papercolor-dark"},
    {"papercolorlight", "papercolor-light"},
    {"rosepine", "rose-pine"},
    {"rosepinedawn", "rose-pine-dawn"},
    {"solarizeddark", "solarized-dark"},
    {"solarizedlight", "solarized-light"},
    {"solarizednightfall", "solarized-nightfall"},
    {"synthwave84", "synthwave-84"},
    {"tokyonight", "tokyo-night"}};

static void normalize_name(const char *src, char *dst, size_t dst_size) {
  size_t out;
  char ch;

  out = 0u;
  while (src != NULL && *src != '\0' && out + 1u < dst_size) {
    ch = (char)tolower((unsigned char)*src);
    if (ch == '_' || ch == ' ') {
      ch = '-';
    }
    if (ch == '-' && out > 0u && dst[out - 1u] == '-') {
      ++src;
      continue;
    }
    dst[out] = ch;
    ++out;
    ++src;
  }
  dst[out] = '\0';
  while (out > 0u && dst[out - 1u] == '-') {
    dst[out - 1u] = '\0';
    --out;
  }
  if (strncmp(dst, "palette-", 8u) == 0) {
    memmove(dst, dst + 8, strlen(dst + 8) + 1u);
  } else if (strncmp(dst, "palette", 7u) == 0) {
    memmove(dst, dst + 7, strlen(dst + 7) + 1u);
    while (*dst == '-') {
      memmove(dst, dst + 1, strlen(dst));
    }
  }
}

const pslog_palette *pslog_palette_default(void) {
  return &pslog_builtin_palette_default;
}

const pslog_palette *pslog_palette_by_name(const char *name) {
  char normalized[64];
  size_t i;

  normalize_name(name, normalized, sizeof(normalized));
  if (normalized[0] == '\0') {
    return &pslog_builtin_palette_default;
  }
  for (i = 0u; i < sizeof(palette_aliases) / sizeof(palette_aliases[0]); ++i) {
    if (strcmp(normalized, palette_aliases[i].alias) == 0) {
      strcpy(normalized, palette_aliases[i].canonical);
      break;
    }
  }
  for (i = 0u; i < sizeof(palette_entries) / sizeof(palette_entries[0]); ++i) {
    if (strcmp(normalized, palette_entries[i].name) == 0) {
      return palette_entries[i].palette;
    }
  }
  return &pslog_builtin_palette_default;
}

size_t pslog_palette_count(void) {
  return sizeof(palette_entries) / sizeof(palette_entries[0]);
}

const char *pslog_palette_name(size_t index) {
  if (index >= pslog_palette_count()) {
    return NULL;
  }
  return palette_entries[index].name;
}
