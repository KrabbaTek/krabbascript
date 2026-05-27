#pragma once

#define _POSIX_C_SOURCE 200809L
#include <common/common.h>
#include <stdio.h>
#include <string.h>
#include <tokenizer/tokenizer.h>
#include <types/types.h>

#define perr(msg)                                                              \
    sourceError(tokenVectorPeek(tokens, *index).line,                          \
                tokenVectorPeek(tokens, *index).col,                           \
                source_file,                                                   \
                msg);

binding_power_t getBindingPower(token_type type);
ast_node_type   tokenToExpression(token_t token);

ast_parent_t*
astParseExpression(token_vector_t* tokens, size_t* index, float min_bp);
ast_parent_t* astParseTokens(token_vector_t* tokens, char* source_file);

symbol_table_t* createSymbolTableFromTokens(token_vector_t* tokens);

ast_parent_t*
parserHandleVal(token_vector_t* tokens, size_t* index, char* source_file);

ast_parent_t*
parserHandleVar(token_vector_t* tokens, size_t* index, char* source_file);

ast_parent_t*
parserHandleFrom(token_vector_t* tokens, size_t* index, char* source_file);

ast_parent_t*
parserHandleFunction(token_vector_t* tokens, size_t* index, char* source_file);

void parserPrintNodeType(ast_node_t* node);
void parserDumpNode(ast_parent_t* parent, size_t depth);

ast_parent_t* parserParseFunctionParam(token_vector_t* tokens,
                                       size_t*         index,
                                       char*           source_file);