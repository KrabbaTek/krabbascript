/*
    KrabbaScript project
    File: src/common/common.c
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

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

path_type isFile(const char* name) {
	DIR* directory = opendir(name);

	if (directory != NULL) {
		closedir(directory);
		return KSCRIPT_PATH_TYPE_DIR;
	}

	if (errno == ENOTDIR) return KSCRIPT_PATH_TYPE_FILE;

	return KSCRIPT_PATH_TYPE_NONE;
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

void kscriptTomlError(char* msg) {
	errors_generated++;
	printf("\033[1;31mERROR\033[0m: TOML: %s\n", msg);
}

compile_flags_t* buildDirectory(char* dir) {
	compile_flags_t* compile_flags = malloc(sizeof(compile_flags_t));

	char toml_file_path[PATH_MAX];

	snprintf(toml_file_path, sizeof(toml_file_path), "%s/manifest.toml", dir);
	toml_result_t result = toml_parse_file_ex(toml_file_path);

	if (!result.ok) {
		terr(result.errmsg);
		toml_free(result);

		return NULL;
	}

	// Get the project fields
	// Project name
	toml_datum_t project_name = toml_seek(result.toptab, "project.name");

	if (project_name.type == TOML_UNKNOWN) {
		terr("project.name doesn't exist");
		toml_free(result);

		return NULL;
	}


	if (project_name.type != TOML_STRING) {
		terr("project.name must be a string");
		toml_free(result);

		return NULL;
	}

	compile_flags->bin_name = strdup(project_name.u.s);

	// Project version
	toml_datum_t project_ver = toml_seek(result.toptab, "project.version");

	if (project_ver.type == TOML_UNKNOWN) {
		terr("project.version doesn't exist");
		toml_free(result);

		return NULL;
	}

	if (project_ver.type != TOML_STRING) {
		terr("project.version must be a string");
		toml_free(result);

		return NULL;
	}

	toml_free(result);
	return compile_flags;
}