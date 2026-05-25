#pragma once

#include <assert.h>
#include <common/common.h>
#include <stdio.h>
#include <stdlib.h>
#include <types/types.h>

char_vector_t* newCharVector();
char_vector_t* charVectorFromString(const char* string);
void           resetCharVector(char_vector_t* vector);
void           freeCharVector(char_vector_t* vector);

void charVectorPush(char_vector_t* vector, char val);
char charVectorPop(char_vector_t* vector);

char charVectorPeek(char_vector_t* vector, size_t index);

token_vector_t* newTokenVector();
void            tokenVectorPush(token_vector_t* vector, token_t val);
token_t         tokenVectorPop(token_vector_t* vector);

token_t tokenVectorPeek(token_vector_t* vector, size_t index);
void    freeTokenVector(token_vector_t* vector);
void    freeToken(token_t* token);

bool tokenVectorEof(token_vector_t* vector, size_t index);

ast_node_t* newNode();
void        freeNode(ast_node_t* node);

symbol_table_t* newSymbolTable();
void            freeSymbolTable(symbol_table_t* table);

void       symbolTablePush(symbol_table_t* table, st_entry_t entry);
st_entry_t symbolTablePop(symbol_table_t* table);

st_entry_t symbolTablePeek(symbol_table_t* table, size_t index);

void newNodeBlock(ast_node_t* node);
void freeNodeBlock(ast_node_t* node);

void        nodeBlockPush(ast_node_t* node, ast_node_t* child);
ast_node_t* nodeBlockPop(ast_node_t* node);

void newFuncArgs(ast_node_t* node);
void freeFuncArgs(ast_node_t* node);

void        funcArgsPush(ast_node_t* node, ast_node_t* child);
ast_node_t* funcArgskPop(ast_node_t* node);

void resetTokenVector(token_vector_t* vector);