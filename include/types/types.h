/*
    KrabbaScript project
    File: include/types/types.h
    Purpose: Contains types for the tokenizer like char vectors, tokens...
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
        - Visser13
*/

#pragma once

#include <stdbool.h>
#include <stddef.h>

#ifndef KSCRIPT_HTABLE_BUCKET_SIZE
// Hash table bucket size
#	define KSCRIPT_HTABLE_BUCKET_SIZE 1000
#endif

typedef enum {
	// Types
	KSCRIPT_TOKEN_TYPE_I64,
	KSCRIPT_TOKEN_TYPE_I32,
	KSCRIPT_TOKEN_TYPE_I16,
	KSCRIPT_TOKEN_TYPE_I8,

	KSCRIPT_TOKEN_TYPE_U64,
	KSCRIPT_TOKEN_TYPE_U32,
	KSCRIPT_TOKEN_TYPE_U16,
	KSCRIPT_TOKEN_TYPE_U8,

	KSCRIPT_TOKEN_TYPE_F64,
	KSCRIPT_TOKEN_TYPE_F32,

	KSCRIPT_TOKEN_TYPE_PTR,
	KSCRIPT_TOKEN_TYPE_STR,
	KSCRIPT_TOKEN_TYPE_BOOL,

	KSCRIPT_TOKEN_TYPE_ARRAY,

	// Compiler stuff,
	KSCRIPT_TOKEN_TYPE_LITERAL,
	KSCRIPT_TOKEN_TYPE_INT_LITERAL,
	KSCRIPT_TOKEN_TYPE_FLOAT_LITERAL,
	KSCRIPT_TOKEN_TYPE_STR_LITERAL,
	KSCRIPT_TOKEN_TYPE_BOOL_LITERAL,
	KSCRIPT_TOKEN_TYPE_NONE,

	KSCRIPT_TOKEN_TYPE_EOF,
	KSCRIPT_TOKEN_TYPE_ARROW,
	KSCRIPT_TOKEN_TYPE_EQUALS_EQUALS,

	// Keywords
	KSCRIPT_TOKEN_TYPE_FUNCTION,
	KSCRIPT_TOKEN_TYPE_VAR,
	KSCRIPT_TOKEN_TYPE_VAL,
	KSCRIPT_TOKEN_TYPE_IMPORT,
	KSCRIPT_TOKEN_TYPE_FROM,
	KSCRIPT_TOKEN_TYPE_MODULE,
	KSCRIPT_TOKEN_TYPE_WHILE,
	KSCRIPT_TOKEN_TYPE_TRUE,
	KSCRIPT_TOKEN_TYPE_FALSE,
	KSCRIPT_TOKEN_TYPE_FOR,
	KSCRIPT_TOKEN_TYPE_PUB,

	// Symbols
	KSCRIPT_TOKEN_TYPE_OPEN_PAREN,
	KSCRIPT_TOKEN_TYPE_CLOSED_PAREN,

	KSCRIPT_TOKEN_TYPE_SEMICOLON,
	KSCRIPT_TOKEN_TYPE_COLON,
	KSCRIPT_TOKEN_TYPE_EQUALS,

	KSCRIPT_TOKEN_TYPE_EXMARK,
	KSCRIPT_TOKEN_TYPE_QMARK,

	KSCRIPT_TOKEN_TYPE_DOT,
	KSCRIPT_TOKEN_TYPE_COMMA,

	KSCRIPT_TOKEN_TYPE_CBRACKET_OPEN,
	KSCRIPT_TOKEN_TYPE_CBRACKET_CLOSED,

	KSCRIPT_TOKEN_TYPE_SBRACKET_OPEN,
	KSCRIPT_TOKEN_TYPE_SBRACKET_CLOSED,

	KSCRIPT_TOKEN_TYPE_LESS_THAN,
	KSCRIPT_TOKEN_TYPE_MORE_THAN,
	KSCRIPT_TOKEN_TYPE_STAR,

	KSCRIPT_TOKEN_TYPE_PLUS,
	KSCRIPT_TOKEN_TYPE_MINUS,

	KSCRIPT_TOKEN_TYPE_SLASH,
	KSCRIPT_TOKEN_TYPE_BSLASH,

	KSCRIPT_TOKEN_TYPE_IF,
	KSCRIPT_TOKEN_TYPE_ELIF,
	KSCRIPT_TOKEN_TYPE_ELSE,

	KSCRIPT_TOKEN_TYPE_RETURN,

	KSCRIPT_TOKEN_TYPE_PLUS_EQ,
	KSCRIPT_TOKEN_TYPE_MINUS_EQ,
	KSCRIPT_TOKEN_TYPE_MUL_EQ,
	KSCRIPT_TOKEN_TYPE_DIV_EQ,

	KSCRIPT_TOKEN_TYPE_NOT_EQU,
	KSCRIPT_TOKEN_TYPE_GREATER_EQU,
	KSCRIPT_TOKEN_TYPE_LESS_EQU,

	KSCRIPT_TOKEN_TYPE_PLUS_PLUS,
	KSCRIPT_TOKEN_TYPE_MINUS_MINUS,

	KSCRIPT_TOKEN_TYPE_LOGICAL_OR,
	KSCRIPT_TOKEN_TYPE_LOGICAL_AND,

	KSCRIPT_TOKEN_TYPE_BWISE_AND,
	KSCRIPT_TOKEN_TYPE_BWISE_OR,

	KSCRIPT_TOKEN_TYPE_BWISE_XOR,

	KSCRIPT_TOKEN_TYPE_POWER_OF

} token_type;

typedef struct {
	token_type type; /**< Token type */

	union {
		int    i;
		float  f;
		double d;

		char* s;
		bool  b;
	};

	bool s_owned; /**< If the s value was mallocated. If it was, s gets cleaned
	                 up */

	int line;
	int col;

} token_t;

typedef struct {
	char* data;

	size_t size;
	size_t capacity;

} char_vector_t;

typedef struct {
	token_t* data;

	size_t size;
	size_t capacity;

} token_vector_t;

typedef enum {
	// Arithmetic operators
	KSCRIPT_AST_NODE_TYPE_DIV,
	KSCRIPT_AST_NODE_TYPE_MUL,
	KSCRIPT_AST_NODE_TYPE_PLUS,
	KSCRIPT_AST_NODE_TYPE_MINUS,

	// Logical & bitwise operators
	KSCRIPT_AST_NODE_TYPE_LOGICAL_AND,
	KSCRIPT_AST_NODE_TYPE_LOGICAL_OR,
	KSCRIPT_AST_NODE_TYPE_LOGICAL_NOT,

	KSCRIPT_AST_NODE_TYPE_BWISE_AND,
	KSCRIPT_AST_NODE_TYPE_BWISE_OR,
	KSCRIPT_AST_NODE_TYPE_BWISE_XOR,
	KSCRIPT_AST_NODE_TYPE_BWISE_NOT,

	KSCRIPT_AST_NODE_TYPE_VARIABLE_DEC,
	KSCRIPT_AST_NODE_TYPE_VALUE_DEC,

	KSCRIPT_AST_NODE_TYPE_VARIABLE_DEF,
	KSCRIPT_AST_NODE_TYPE_VALUE_DEF,

	KSCRIPT_AST_NODE_TYPE_FUNCTION_DEC,
	KSCRIPT_AST_NODE_TYPE_FUNCTION_DEF,

	KSCRIPT_AST_NODE_TYPE_LITERAL,
	KSCRIPT_AST_NODE_TYPE_BOOL_LITERAL,
	KSCRIPT_AST_NODE_TYPE_STR_LITERAL,
	KSCRIPT_AST_NODE_TYPE_INT_LITERAL,

	KSCRIPT_AST_NODE_TYPE_STR_TYPE,
	KSCRIPT_AST_NODE_TYPE_PTR_TYPE,
	KSCRIPT_AST_NODE_TYPE_BOOL_TYPE,

	KSCRIPT_AST_NODE_TYPE_I64_TYPE,
	KSCRIPT_AST_NODE_TYPE_I32_TYPE,
	KSCRIPT_AST_NODE_TYPE_I16_TYPE,
	KSCRIPT_AST_NODE_TYPE_I8_TYPE,

	KSCRIPT_AST_NODE_TYPE_U64_TYPE,
	KSCRIPT_AST_NODE_TYPE_U32_TYPE,
	KSCRIPT_AST_NODE_TYPE_U16_TYPE,
	KSCRIPT_AST_NODE_TYPE_U8_TYPE,

	KSCRIPT_AST_NODE_TYPE_F64_TYPE,
	KSCRIPT_AST_NODE_TYPE_F32_TYPE,

	KSCRIPT_AST_NODE_TYPE_OPEN_PAREN,
	KSCRIPT_AST_NODE_TYPE_CLOSED_PAREN,

	KSCRIPT_AST_NODE_TYPE_IMPORT,

	KSCRIPT_AST_NODE_TYPE_FLOAT_LITERAL,

	KSCRIPT_AST_NODE_TYPE_NONE,
	KSCRIPT_AST_NODE_TYPE_PROGRAM,
	KSCRIPT_AST_NODE_TYPE_ROOT,
	KSCRIPT_AST_NODE_TYPE_EOF
} ast_node_type;

typedef struct ast_node_t {
	ast_node_type type;

	struct ast_node_t* left;
	struct ast_node_t* right;

	char* lexeme;
	bool  lexeme_owned; /**< If the lexeme value was mallocated. If it was,
	                       lexeme gets cleaned up */

	int line;
	int col;

	union {
		struct {
			struct ast_node_t** body;

			size_t size;
			size_t capacity;
		} block;

		struct {
			struct ast_node_t** args;

			size_t args_size;
			size_t args_capacity;

			struct ast_node_t* ret_type;
			struct ast_node_t* body;
		} function_dec;
	};

} ast_node_t;

typedef ast_node_t ast_parent_t;

typedef struct {
	float left;
	float right;

} binding_power_t;

typedef enum { KSCRIPT_STABLE_ENUM_SCOPE_GLOBAL } stable_enum_scope_t;

typedef struct {
	char* lexeme;
	bool  lexeme_owned; /**< If the lexeme value was mallocated. If it was,
	                       lexeme gets cleaned up */

	union {
		stable_enum_scope_t enum_scope;
		ast_node_t*         node_scope;
	} scope;
} stable_entry_t;

typedef enum {
	KSCRIPT_PATH_TYPE_FILE,
	KSCRIPT_PATH_TYPE_DIR,
	KSCRIPT_PATH_TYPE_NONE
} path_type;

typedef struct {
	char* key;
	bool  key_owned;

	void* value;
	bool  value_owned;

	struct htable_entry_t* next_entry;
} htable_entry_t;

typedef struct {
	size_t           count;
	htable_entry_t** bucket;
} hash_table_t;

typedef struct {
	char*	bin_name;

	// Will be used for VERSION_MAJOR, VERSION_MINOR and VERSION_PATCH macros
	int		ver_major;
	int		ver_minor;
	int		ver_patch;
} compile_flags_t;