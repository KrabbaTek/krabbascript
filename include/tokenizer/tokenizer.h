/*
    KrabbaScript project
    File: include/tokenizer/tokenizer.h
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#pragma once

#define _GNU_SOURCE
#include <common/common.h>
#include <ctype.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <vectors/vectors.h>

char_vector_t* fileToCharVector(const char* path);

void deTokenize(token_vector_t* vector);
void deTokenizeToken(token_t token);

void deTokenizeTokenKeyword(token_t token);

// WARNING: Calls free on the char vector provided!
token_vector_t* tokenize(char_vector_t* vector, const char* debug_path);