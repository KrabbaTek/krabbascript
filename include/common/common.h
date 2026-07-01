/*
    KrabbaScript project
    File: include/common/common.h
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#pragma once

#define _POSIX_C_SOURCE 200809L
#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <types/types.h>
#include <string.h>

#include <linux/limits.h>

#include <tomlc17/tomlc17.h>

#define KSCRIPT_VERSION_MAJOR 2
#define KSCRIPT_VERSION_MINOR 6
#define KSCRIPT_VERSION_PATCH 6

#define KSCRIPT_VERSION_IDENT "pre-alpha"

#ifndef KSCRIPT_VENDOR

#	define KSCRIPT_VENDOR "None"

#endif

extern int errors_generated;

/*
 * Prints usage of the compiler
 */
void printUsage();

/*
 * Prints usage of the init commant
 */
void usageInit();

/*
 * Prints usage of the build commant
 */
void usageBuild();

/*
 * Prints an epic crab
 */
void printKrabba();

/*
 * Checks if the path is a file or a directory
 *  @param name     Path
 *  @return         KSCRIPT_PATH_TYPE_FILE, KSCRIPT_PATH_TYPE_DIR or
 * KSCRIPT_PATH_TYPE_NONE if the path doesn't exist
 */
path_type isFile(const char* name);

/*
 * Prints the amount of errors generated from err macros
 */
void printErrorsGenerated();

/*
 * Prints formatted error where an error happened in a file. Will automatically
 * increment the errors_generated number
 *  @param line     Line number
 *  @param col      Column number
 *  @param source   Source file
 *  @param msg      Error message
 */
void sourceError(int line, int col, char* source, char* msg);

/*
 * Prints a general error, like memory failing to allocate
 *  @param msg      Error message
 */
void kscriptError(char* msg);
void kscriptTomlError(char* msg);

#define err(msg) kscriptError(msg);
#define terr(msg) kscriptTomlError(msg);

/*
 * Builds a directory and parses the TOML file in the project alongside with the compiler flags in it
 *  @param dir      Directory
 *  @return         Compiler flags
 */
compile_flags_t* buildDirectory(char* dir);