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

/*
 * Converts a file to a char vector
 *  @param path         Path to file. If it doesn't exist, calls exit(1)
 *  @return             Pointer to a char vector
 */
char_vector_t* fileToCharVector(const char* path);

/*
 * Detokenizes contents of a token vector and prints them
 *  @param vector       Vector, returns if NULL
 */
void deTokenize(token_vector_t* vector);

/*
 * Detokenizes contents of a token and prints them
 *  @param token        Token
 */
void deTokenizeToken(token_t token);

/*
 * Detokenizes token's keyword and prints it
 *  @param token        Token
 */
void deTokenizeTokenKeyword(token_t token);

/*
 * Tokenizes a char vector to a token vector
 *  @param vector       Char vector, returns an empty vector if NULL. The
 * function calls freeCharVector internally on it
 *  @param debug_path   Source file path, used for printing errors
 *  @return             Pointer to a token vector
 */
token_vector_t* tokenize(char_vector_t* vector, const char* debug_path);