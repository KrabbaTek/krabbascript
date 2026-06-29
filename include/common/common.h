/*
    KrabbaScript project
    File: include/common/common.h
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#pragma once

#include <dirent.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>

#define KSCRIPT_DIRECTORY 0
#define KSCRIPT_FILE 1

#define KSCRIPT_VERSION_MAJOR 2
#define KSCRIPT_VERSION_MINOR 6
#define KSCRIPT_VERSION_PATCH 6

#define KSCRIPT_VERSION_IDENT "pre-alpha"

#ifndef KSCRIPT_VENDOR

#	define KSCRIPT_VENDOR "None"

#endif

extern int errors_generated;

void printUsage();

void usageInit();
void usageBuild();

void printKrabba();

int isFile(const char* name);

void printErrorsGenerated();

void sourceError(int line, int col, char* source, char* msg);
void kscriptError(char* msg);

#define err(msg) kscriptError(msg);