#include <common/common.h>

int errors_generated = 0;

void printUsage() {
    printf("Usage: kscript [options]\n\n");
    printf("Commands: \n");

    printf("    init    <name>      Create a new Krabbascript project\n");
    printf("    build   <file|dir>  Build a project or a file\n");
    printf("    version             Display the current compiler version\n");
    printf("    help                Display this menu\n");
}

void usageInit() {
    printf("Usage: kscript init <name>\n");
}

void usageBuild() {
    printf("Usage: kscript build <file|dir>\n");
}

int isFile(const char* name) {
    DIR* directory = opendir(name);

    if (directory != NULL) {
        closedir(directory);
        return 0;
    }

    if (errno == ENOTDIR) return 1;

    return -1;
}

void printKrabba() {
    printf("              ▓▓▓▓▓▓▓▓▓▓▓▓                    \n"
           "  ▓▓▓▓      ▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓                  \n"
           "  ▓▓      ▓▓▓▓    ▓▓██    ██                  \n"
           "  ▓▓▓▓    ▓▓▓▓    ████              ▓▓▓▓      \n"
           "    ▓▓    ▓▓▓▓                        ▓▓▓▓    \n"
           "    ▓▓▓▓  ▓▓████        ▓▓        ▓▓    ██    \n"
           "      ▓▓▓▓  ██▓▓▓▓▓▓▓▓▓▓          ▓▓▓▓▓▓██    \n"
           "▓▓      ████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓    ▓▓▓▓██    \n"
           "▓▓▓▓▓▓    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓        ▓▓██    \n"
           "  ▓▓▓▓████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓      ▓▓██    \n"
           "          ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓    ████  ▓▓\n"
           "    ▓▓██████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████████    ██\n"
           "  ▓▓▓▓      ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████████    ████\n"
           "  ▓▓  ▓▓▓▓██████▓▓▓▓▓▓▓▓▓▓████████  ▓▓▓▓▓▓██  \n"
           "    ▓▓▓▓        ████▓▓▓▓██████████▓▓          \n"
           "    ▓▓▓▓            ▓▓████████    ▓▓▓▓    ██  \n"
           "      ▓▓▓▓              ▓▓  ▓▓▓▓    ▓▓▓▓▓▓██  \n"
           "                        ▓▓    ▓▓▓▓            \n"
           "                      ████      ██            \n"
           "                  ▓▓████      ████            \n");

    // Krabba!!!!
}

void printErrorsGenerated() {
    if (errors_generated > 0) {
        printf("\033[1;31m==== BUILD FAILED with %d %s ====\033[0m\n",
               errors_generated,
               errors_generated == 1 ? "error" : "errors");

        exit(1);
    }
}

void sourceError(int line, int col, char* source, char* msg) {
    errors_generated++;
    printf("%s:%d:%d: \033[1;31mERROR\033[0m: %s", source, line, col, msg);
}

void kscriptError(char* msg) {
    errors_generated++;
    printf("\033[1;31mERROR\033[0m: %s", msg);
}