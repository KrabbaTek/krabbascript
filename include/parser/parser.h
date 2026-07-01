/*
    KrabbaScript project
    File: include/parser/parser.h
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

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

/*
 * Returns a binding power of a token
 *  @param type     Token
 *  @return         Binding power of the token
 */
binding_power_t getBindingPower(token_type type);

/*
 * Converts a token to a node
 *  @param type     Token
 *  @return         Pointer to a node. If the token is unknown, returns a
 * newNode()
 */
ast_node_t* tokenToNode(token_t token);

/*
 * Converts a sequence of tokens to an expression
 *  @param tokens   Tokens, does not handle if its NULL
 *  @param index    Pointer to an index. Will be used by the function, so
 * recommended to leave it 0
 *  @param min_bp   Minimum binding power, leave 0.0f
 *  @return         Pointer to a parent node.
 */
ast_parent_t*
astParseExpression(token_vector_t* tokens, size_t* index, float min_bp);

/*
 * Converts a sequence of tokens to an AST
 *  @param tokens       Tokens, does not handle if its NULL
 *  @param source_file  Path to the source file, used for printing errors
 *  @return             Pointer to a parent node.
 */
ast_parent_t* astParseTokens(token_vector_t* tokens, char* source_file);

/*
 * Handles the val declaration/definition
 *  @param tokens       Tokens, does not handle if its NULL
 *  @param index        Pointer to an index. Will be used by the function
 *  @param source_file  Path to the source file, used for printing errors
 *  @return             Pointer to a parent node.
 */
ast_parent_t*
parserHandleVal(token_vector_t* tokens, size_t* index, char* source_file);

/*
 * Handles the var declaration/definition
 *  @param tokens       Tokens, does not handle if its NULL
 *  @param index        Pointer to an index. Will be used by the function
 *  @param source_file  Path to the source file, used for printing errors
 *  @return             Pointer to a parent node.
 */
ast_parent_t*
parserHandleVar(token_vector_t* tokens, size_t* index, char* source_file);

/*
 * Handles the from declaration
 *  @param tokens       Tokens, does not handle if its NULL
 *  @param index        Pointer to an index. Will be used by the function
 *  @param source_file  Path to the source file, used for printing errors
 *  @return             Pointer to a parent node.
 */
ast_parent_t*
parserHandleFrom(token_vector_t* tokens, size_t* index, char* source_file);

/*
 * Handles the function declaration/definition
 *  @param tokens       Tokens, does not handle if its NULL
 *  @param index        Pointer to an index. Will be used by the function
 *  @param source_file  Path to the source file, used for printing errors
 *  @return             Pointer to a parent node.
 */
ast_parent_t*
parserHandleFunction(token_vector_t* tokens, size_t* index, char* source_file);

/*
 * Prints a node type
 *  @param node     Node, returns if its NULL
 */
void parserPrintNodeType(ast_node_t* node);

/*
 * Prints the parent node and it's children
 *  @param node     Pointer to a parent node, returns if NULL
 *  @param depth    How many tabs to put. Leave 0
 */
void parserDumpNode(ast_parent_t* parent, size_t depth);

/*
 * Helper function that handles the function parameters
 *  @param tokens       Tokens, does not handle if its NULL
 *  @param index        Pointer to an index. Will be used by the function
 *  @param source_file  Path to the source file, used for printing errors
 *  @return             Pointer to a parent node.
 */
ast_parent_t* parserParseFunctionParam(token_vector_t* tokens,
                                       size_t*         index,
                                       char*           source_file);