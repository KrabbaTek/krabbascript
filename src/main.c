/*
    KrabbaScript project
    File: src/main.c
    Contributors (To this file):
        - khytryy (yehor.khytryy@gmail.com)
*/

#include <common/common.h>
#include <linux/limits.h>
#include <parser/parser.h>
#include <stdio.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <tokenizer/tokenizer.h>
#include <unistd.h>

int main(int argc, char** argv) {
	if (argc < 2) {
		printUsage();

		return 0;
	}

	if (strcmp(argv[1], "version") == 0) {
		printf("Krabbascript (kscript) %d.%d.%d-%s (%s)\n",
		       KSCRIPT_VERSION_MAJOR,
		       KSCRIPT_VERSION_MINOR,
		       KSCRIPT_VERSION_PATCH,
		       KSCRIPT_VERSION_IDENT,
		       KSCRIPT_VENDOR);

		printf("Copyright © 2026 KrabbaTek Group\n");

		return 0;
	} else if (strcmp(argv[1], "build") == 0) {
		if (argc < 3) {
			err("No input file|directory\n");
			usageBuild();

			return 1;
		}

		int result = isFile(argv[2]);
		if (result == KSCRIPT_PATH_TYPE_FILE) {
			char_vector_t* source = fileToCharVector(argv[2]);

			token_vector_t* tokens = tokenize(source, argv[2]);
			deTokenize(tokens);

			ast_parent_t* program = astParseTokens(tokens, argv[2]);
			freeTokenVector(tokens);

			parserDumpNode(program, 0);

			freeNodeBlock(program);
			freeNode(program);
		}

		else if (result == KSCRIPT_PATH_TYPE_DIR) {
			compile_flags_t* flags = buildDirectory(argv[2]);
			printErrorsGenerated();
			
			err("Directory buidling is not yet implemented, only the manifest parsing\n");

			free(flags->bin_name);
			free(flags);
			return 0;
		} else {
			err("No such file or directory \"");
			printf("%s\"\n", argv[2]);

			return 1;
		}

		return 0;
	} else if (strcmp(argv[1], "init") == 0) {
		if (argc < 3) {
			err("No directory name\n");
			usageBuild();

			return 1;
		}

		struct stat st = {0};

		if (stat(argv[2], &st) == -1) {
			// Create the directories for the project
			if (mkdir(argv[2], 0700) == -1) {
				err("Failed to create project directory\n");
				return 1;
			}

			char src_path[PATH_MAX];

			if (snprintf(src_path, sizeof(src_path), "%s/src", argv[2]) >=
			    (int)sizeof(src_path)) {
				err("Project path too long\n");
				return 1;
			}

			if (mkdir(src_path, 0700) == -1) {
				err("Failed to create src directory\n");
				return 1;
			}

			if (snprintf(src_path,
			             sizeof(src_path),
			             "%s/src/main.krs",
			             argv[2]) >= (int)sizeof(src_path)) {
				err("Project path too long\n");
				return 1;
			}

			// Create a new hello world source file
			FILE* hello_world = fopen(src_path, "w");

			if (!hello_world) {
				err("Failed to create main.krs\n");
				return 1;
			}

			fprintf(hello_world,
			        "from <stdlib.krs> import *;\nfunc main() -> i32 {\n    "
			        "std.print(\"Hello, world!\");\n\n    return 0;\n}\n");
			fclose(hello_world);

			printf("Initialized %s project\n", argv[2]);
			return 0;
		} else {
			err("Directory already exists\n");
			return 1;
		}

	} else if (strcmp(argv[1], "krabba") == 0) {
		printKrabba();
		printf("Krabba!\n");
		return 0;
	} else if (strcmp(argv[1], "help") == 0) {
		printUsage();
		return 0;
	}

	err("Unknown command \"");
	printf("%s\"\n", argv[1]);

	printUsage();

	return 1;
}
