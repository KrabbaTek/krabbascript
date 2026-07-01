/*
    KrabbaScript project
    File: src/parser/parser.c
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#include <parser/parser.h>

binding_power_t getBindingPower(token_type type) {
	switch (type) {
		case KSCRIPT_TOKEN_TYPE_SBRACKET_OPEN:
		case KSCRIPT_TOKEN_TYPE_SBRACKET_CLOSED:
		case KSCRIPT_TOKEN_TYPE_DOT:
		case KSCRIPT_TOKEN_TYPE_ARROW:
		case KSCRIPT_TOKEN_TYPE_PLUS_PLUS:
		case KSCRIPT_TOKEN_TYPE_MINUS_MINUS:
			return (binding_power_t){1.0f, 1.1f};
			break;

		case KSCRIPT_TOKEN_TYPE_MORE_THAN:
		case KSCRIPT_TOKEN_TYPE_LESS_THAN:
		case KSCRIPT_TOKEN_TYPE_LESS_EQU:
		case KSCRIPT_TOKEN_TYPE_GREATER_EQU:
			return (binding_power_t){10.0f, 10.1f};
			break;

		case KSCRIPT_TOKEN_TYPE_EQUALS_EQUALS:
		case KSCRIPT_TOKEN_TYPE_NOT_EQU:
			return (binding_power_t){9.0f, 9.1f};
			break;

		case KSCRIPT_TOKEN_TYPE_LOGICAL_AND:
		case KSCRIPT_TOKEN_TYPE_LOGICAL_OR:
			return (binding_power_t){5.0f, 5.1f};
			break;

		case KSCRIPT_TOKEN_TYPE_PLUS:
		case KSCRIPT_TOKEN_TYPE_MINUS:
			return (binding_power_t){20.0f, 20.1f};
			break;
		case KSCRIPT_TOKEN_TYPE_STAR:
		case KSCRIPT_TOKEN_TYPE_SLASH:
			return (binding_power_t){30.0f, 30.1f};
			break;
		default:
			return (binding_power_t){0.0f, 0.0f};
			break;
	}
}

ast_node_t* tokenToNode(token_t token) {
	ast_node_t* node = newNode();

	switch (token.type) {
		case KSCRIPT_TOKEN_TYPE_LITERAL:
			node->type         = KSCRIPT_AST_NODE_TYPE_LITERAL;
			node->lexeme       = strdup(token.s);
			node->lexeme_owned = true;

			break;
		case KSCRIPT_TOKEN_TYPE_BOOL_LITERAL:
			node->type = KSCRIPT_AST_NODE_TYPE_BOOL_LITERAL;

			node->lexeme       = token.b == true ? "true" : "false";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_FLOAT_LITERAL:
			node->type = KSCRIPT_AST_NODE_TYPE_FLOAT_LITERAL;

			int len2     = snprintf(NULL, 0, "%f", token.d);
			node->lexeme = malloc(len2 + 1);

			snprintf(node->lexeme, len2 + 1, "%f", token.d);
			node->lexeme_owned = true;

			break;
		case KSCRIPT_TOKEN_TYPE_STR_LITERAL:
			node->type = KSCRIPT_AST_NODE_TYPE_STR_LITERAL;

			node->lexeme       = strdup(token.s);
			node->lexeme_owned = true;

			break;
		case KSCRIPT_TOKEN_TYPE_INT_LITERAL:
			node->type = KSCRIPT_AST_NODE_TYPE_INT_LITERAL;

			char buffer[32];
			int  len = snprintf(buffer, sizeof(buffer), "%d", token.i);

			node->lexeme       = malloc(len + 1);
			node->lexeme_owned = true;

			memcpy(node->lexeme, buffer, len + 1);

			break;
		case KSCRIPT_TOKEN_TYPE_EOF:
			node->type = KSCRIPT_AST_NODE_TYPE_EOF;

			node->lexeme       = "EOF";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_PLUS:
			node->type = KSCRIPT_AST_NODE_TYPE_PLUS;

			node->lexeme       = "+";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_MINUS:
			node->type = KSCRIPT_AST_NODE_TYPE_MINUS;

			node->lexeme       = "-";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_STAR:
			node->type = KSCRIPT_AST_NODE_TYPE_MUL;

			node->lexeme       = "*";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_SLASH:
			node->type = KSCRIPT_AST_NODE_TYPE_DIV;

			node->lexeme       = "/";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_LOGICAL_AND:
			node->type = KSCRIPT_AST_NODE_TYPE_LOGICAL_AND;

			node->lexeme       = "&&";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_LOGICAL_OR:
			node->type = KSCRIPT_AST_NODE_TYPE_LOGICAL_OR;

			node->lexeme       = "||";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_BWISE_AND:
			node->type = KSCRIPT_AST_NODE_TYPE_BWISE_AND;

			node->lexeme       = "&";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_BWISE_OR:
			node->type = KSCRIPT_AST_NODE_TYPE_BWISE_OR;

			node->lexeme       = "|";
			node->lexeme_owned = false;

			break;
		case KSCRIPT_TOKEN_TYPE_BWISE_XOR:
			node->type = KSCRIPT_AST_NODE_TYPE_BWISE_XOR;

			node->lexeme       = "^";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_STR:
			node->type = KSCRIPT_AST_NODE_TYPE_STR_TYPE;

			node->lexeme       = "str";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_I64:
			node->type = KSCRIPT_AST_NODE_TYPE_I64_TYPE;

			node->lexeme       = "i64";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_I32:
			node->type = KSCRIPT_AST_NODE_TYPE_I32_TYPE;

			node->lexeme       = "i32";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_I16:
			node->type = KSCRIPT_AST_NODE_TYPE_I16_TYPE;

			node->lexeme       = "i16";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_I8:
			node->type = KSCRIPT_AST_NODE_TYPE_I8_TYPE;

			node->lexeme       = "i8";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_U64:
			node->type = KSCRIPT_AST_NODE_TYPE_U64_TYPE;

			node->lexeme       = "u64";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_U32:
			node->type = KSCRIPT_AST_NODE_TYPE_U32_TYPE;

			node->lexeme       = "u32";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_U16:
			node->type = KSCRIPT_AST_NODE_TYPE_U16_TYPE;

			node->lexeme       = "u16";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_U8:
			node->type = KSCRIPT_AST_NODE_TYPE_U8_TYPE;

			node->lexeme       = "u8";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_F64:
			node->type = KSCRIPT_AST_NODE_TYPE_F64_TYPE;

			node->lexeme       = "f64";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_F32:
			node->type = KSCRIPT_AST_NODE_TYPE_F32_TYPE;

			node->lexeme       = "f32";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_BOOL:
			node->type = KSCRIPT_AST_NODE_TYPE_BOOL_TYPE;

			node->lexeme       = "bool";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_PTR:
			node->type = KSCRIPT_AST_NODE_TYPE_PTR_TYPE;

			node->lexeme       = "ptr";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_OPEN_PAREN:
			node->type = KSCRIPT_AST_NODE_TYPE_OPEN_PAREN;

			node->lexeme       = "(";
			node->lexeme_owned = false;

			break;

		case KSCRIPT_TOKEN_TYPE_CLOSED_PAREN:
			node->type = KSCRIPT_AST_NODE_TYPE_CLOSED_PAREN;

			node->lexeme       = ")";
			node->lexeme_owned = false;

			break;

		default:
			err("Unknown token ");
			deTokenizeTokenKeyword(token);
			printf("\n");

			node->type = KSCRIPT_AST_NODE_TYPE_NONE;

			break;
	}

	node->line = token.line;
	node->col  = token.col;
	return node;
}

ast_parent_t*
astParseExpression(token_vector_t* tokens, size_t* index, float min_bp) {
	ast_node_t* lhs;

	if (tokenVectorPeek(tokens, *index).type == KSCRIPT_TOKEN_TYPE_OPEN_PAREN) {
		// Consume (
		(*index)++;
		lhs = astParseExpression(tokens, index, 0.0f);

		// Consume ) if it even exists
		if (tokenVectorPeek(tokens, *index).type ==
		    KSCRIPT_TOKEN_TYPE_CLOSED_PAREN) {
			(*index)++;
		}
	} else {
		lhs = tokenToNode(tokenVectorPeek(tokens, *index));
		(*index)++;
	}

	while (true) {
		token_t op = tokenVectorPeek(tokens, *index);

		if (op.type == KSCRIPT_TOKEN_TYPE_EOF) break;
		if (op.type == KSCRIPT_TOKEN_TYPE_CLOSED_PAREN) break;

		binding_power_t bp =
		        getBindingPower(tokenVectorPeek(tokens, *index).type);

		if (bp.left < min_bp) break;

		(*index)++;
		ast_node_t* op_node = tokenToNode(op);

		ast_node_t* rhs = astParseExpression(tokens, index, bp.right);

		op_node->left  = lhs;
		op_node->right = rhs;

		lhs = op_node;
	}
	return lhs;
}

ast_parent_t* astParseTokens(token_vector_t* tokens, char* source_file) {
	ast_parent_t* root = newNode();
	newNodeBlock(root);

	size_t index = 0;

	root->type = KSCRIPT_AST_NODE_TYPE_PROGRAM;

	if (tokenVectorEof(tokens, 0)) {
		return root;
	}

	while (!tokenVectorEof(tokens, index)) {
		// Parse val statement
		if (!tokenVectorEof(tokens, index) &&
		    tokenVectorPeek(tokens, index).type == KSCRIPT_TOKEN_TYPE_VAL) {

			index++;
			ast_parent_t* parent = parserHandleVal(tokens, &index, source_file);

			nodeBlockPush(root, parent);
			index++;
		}

		// Parse var statement
		else if (!tokenVectorEof(tokens, index) &&
		         tokenVectorPeek(tokens, index).type ==
		                 KSCRIPT_TOKEN_TYPE_VAR) {

			index++;
			ast_parent_t* parent = parserHandleVar(tokens, &index, source_file);
			nodeBlockPush(root, parent);

			index++;
		}

		// Parse from statement
		else if (!tokenVectorEof(tokens, index) &&
		         tokenVectorPeek(tokens, index).type ==
		                 KSCRIPT_TOKEN_TYPE_FROM) {

			index++;
			ast_parent_t* parent =
			        parserHandleFrom(tokens, &index, source_file);
			nodeBlockPush(root, parent);

			index++;
		} else if (!tokenVectorEof(tokens, index) &&
		           tokenVectorPeek(tokens, index).type ==
		                   KSCRIPT_TOKEN_TYPE_FUNCTION) {
			index++;
			ast_parent_t* parent =
			        parserHandleFunction(tokens, &index, source_file);
			nodeBlockPush(root, parent);

			index++;
		} else
			index++;
	}

	// parserDumpNode(root, 0);
	return root;
}

ast_parent_t*
parserHandleVal(token_vector_t* tokens, size_t* index, char* source_file) {
	if (tokenVectorEof(tokens, *index)) {
		perr("Expected a literal, got <EOF>\n");

		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_LITERAL) {
		perr("Expected a literal, got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	token_t ident = tokenVectorPeek(tokens, *index);
	(*index)++;

	if (tokenVectorEof(tokens, *index)) {
		perr("Expected '=', ':', got <EOF>\n");

		return newNode();
	}

	if (!tokenVectorEof(tokens, *index) &&
	    tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_EQUALS &&
	    tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_COLON) {

		perr("Expected '=', ':', got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	// Handle the equals sign
	if (!tokenVectorEof(tokens, *index) &&
	    tokenVectorPeek(tokens, *index).type == KSCRIPT_TOKEN_TYPE_EQUALS) {
		(*index)++;
		token_vector_t* buffer = newTokenVector();

		if (tokenVectorEof(tokens, *index)) {
			perr("Expected expression, got <EOF>\n") freeTokenVector(buffer);

			return newNode();
		}

		// Collect everything into a buffer
		while (tokenVectorPeek(tokens, *index).type !=
		               KSCRIPT_TOKEN_TYPE_SEMICOLON &&
		       !tokenVectorEof(tokens, *index)) {
			tokenVectorPush(buffer, tokenVectorPeek(tokens, *index));
			(*index)++;
		}

		size_t        index_exp = 0;
		ast_parent_t* exp       = astParseExpression(buffer, &index_exp, 0.0f);
		freeTokenVector(buffer);

		ast_parent_t* node = newNode();
		node->type         = KSCRIPT_AST_NODE_TYPE_VALUE_DEF;
		node->lexeme       = strdup(ident.s);
		node->lexeme_owned = true;

		node->left = exp;
		return node;
	}

	// Handle the colon sign
	if (!tokenVectorEof(tokens, *index) &&
	    tokenVectorPeek(tokens, *index).type == KSCRIPT_TOKEN_TYPE_COLON) {
		(*index)++;
		if (tokenVectorEof(tokens, *index)) {
			perr("Expected a type, got <EOF>\n");

			return newNode();
		}
		token_t type = tokenVectorPeek(tokens, *index);

		(*index)++;

		if (tokenVectorEof(tokens, *index)) {
			perr("Expected '=' or ';', got <EOF>\n");

			return newNode();
		}

		if (!tokenVectorEof(tokens, *index) &&
		    tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_EQUALS &&
		    tokenVectorPeek(tokens, *index).type !=
		            KSCRIPT_TOKEN_TYPE_SEMICOLON) {
			perr("Expected '=' or ';', got ");

			deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
			printf("\n");

			return newNode();
		}

		// Handle the semicolon
		if (tokenVectorPeek(tokens, *index).type ==
		    KSCRIPT_TOKEN_TYPE_SEMICOLON) {
			ast_parent_t* node = newNode();
			// Type
			node->right = tokenToNode(type);

			node->type         = KSCRIPT_AST_NODE_TYPE_VALUE_DEC;
			node->lexeme       = strdup(ident.s);
			node->lexeme_owned = true;

			return node;
		}

		// Handle the equals sign
		(*index)++;

		token_vector_t* buffer = newTokenVector();

		while (tokenVectorPeek(tokens, *index).type !=
		               KSCRIPT_TOKEN_TYPE_SEMICOLON &&
		       !tokenVectorEof(tokens, *index)) {
			tokenVectorPush(buffer, tokenVectorPeek(tokens, *index));

			(*index)++;
		}

		size_t        index_exp = 0;
		ast_parent_t* exp       = astParseExpression(buffer, &index_exp, 0.0f);

		ast_parent_t* node = newNode();
		// Type
		node->right = tokenToNode(type);

		node->type         = KSCRIPT_AST_NODE_TYPE_VALUE_DEF;
		node->lexeme       = strdup(ident.s);
		node->lexeme_owned = true;

		node->left = exp;
		return node;
	}

	return newNode();
}

ast_parent_t*
parserHandleVar(token_vector_t* tokens, size_t* index, char* source_file) {
	if (tokenVectorEof(tokens, *index)) {
		perr("Expected a literal, got <EOF>\n");

		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_LITERAL) {
		perr("Expected a literal, got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	token_t ident = tokenVectorPeek(tokens, *index);
	(*index)++;

	if (tokenVectorEof(tokens, *index)) {
		perr("Expected '=', ':', got <EOF>\n");

		return newNode();
	}

	if (!tokenVectorEof(tokens, *index) &&
	    tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_EQUALS &&
	    tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_COLON) {

		perr("Expected '=', ':', got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	// Handle the equals sign
	if (!tokenVectorEof(tokens, *index) &&
	    tokenVectorPeek(tokens, *index).type == KSCRIPT_TOKEN_TYPE_EQUALS) {
		(*index)++;
		token_vector_t* buffer = newTokenVector();

		if (tokenVectorEof(tokens, *index)) {
			perr("Expected expression, got <EOF>\n") freeTokenVector(buffer);

			return newNode();
		}

		// Collect everything into a buffer
		while (tokenVectorPeek(tokens, *index).type !=
		               KSCRIPT_TOKEN_TYPE_SEMICOLON &&
		       !tokenVectorEof(tokens, *index)) {
			tokenVectorPush(buffer, tokenVectorPeek(tokens, *index));
			(*index)++;
		}

		size_t        index_exp = 0;
		ast_parent_t* exp       = astParseExpression(buffer, &index_exp, 0.0f);
		freeTokenVector(buffer);

		ast_parent_t* node = newNode();
		node->type         = KSCRIPT_AST_NODE_TYPE_VARIABLE_DEF;
		node->lexeme       = strdup(ident.s);
		node->lexeme_owned = true;

		node->left = exp;
		return node;
	}

	// Handle the colon sign
	if (!tokenVectorEof(tokens, *index) &&
	    tokenVectorPeek(tokens, *index).type == KSCRIPT_TOKEN_TYPE_COLON) {
		(*index)++;
		if (tokenVectorEof(tokens, *index)) {
			perr("Expected a type, got <EOF>\n");

			return newNode();
		}
		token_t type = tokenVectorPeek(tokens, *index);

		(*index)++;

		if (tokenVectorEof(tokens, *index)) {
			perr("Expected '=' or ';', got <EOF>\n");

			return newNode();
		}

		if (!tokenVectorEof(tokens, *index) &&
		    tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_EQUALS &&
		    tokenVectorPeek(tokens, *index).type !=
		            KSCRIPT_TOKEN_TYPE_SEMICOLON) {
			perr("Expected '=' or ';', got ");

			deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
			printf("\n");

			return newNode();
		}

		// Handle the semicolon
		if (tokenVectorPeek(tokens, *index).type ==
		    KSCRIPT_TOKEN_TYPE_SEMICOLON) {
			ast_parent_t* node = newNode();
			// Type
			node->right = tokenToNode(type);

			node->type         = KSCRIPT_AST_NODE_TYPE_VARIABLE_DEC;
			node->lexeme       = strdup(ident.s);
			node->lexeme_owned = true;

			return node;
		}

		// Handle the equals sign
		(*index)++;

		token_vector_t* buffer = newTokenVector();

		while (tokenVectorPeek(tokens, *index).type !=
		               KSCRIPT_TOKEN_TYPE_SEMICOLON &&
		       !tokenVectorEof(tokens, *index)) {
			tokenVectorPush(buffer, tokenVectorPeek(tokens, *index));

			(*index)++;
		}

		size_t        index_exp = 0;
		ast_parent_t* exp       = astParseExpression(buffer, &index_exp, 0.0f);

		ast_parent_t* node = newNode();
		// Type
		node->right = tokenToNode(type);

		node->type         = KSCRIPT_AST_NODE_TYPE_VARIABLE_DEF;
		node->lexeme       = strdup(ident.s);
		node->lexeme_owned = true;

		node->left = exp;
		return node;
	}

	return newNode();
}

ast_parent_t*
parserHandleFrom(token_vector_t* tokens, size_t* index, char* source_file) {
	if (tokenVectorEof(tokens, *index)) {
		perr("Expected a path, got <EOF>\n");

		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type !=
	    KSCRIPT_TOKEN_TYPE_STR_LITERAL) {
		perr("Expected a path, got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	token_t file = tokenVectorPeek(tokens, *index);
	(*index)++;

	if (tokenVectorEof(tokens, *index)) {
		perr("Expected import, got <EOF>\n");

		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_IMPORT) {
		perr("Expected import, got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	(*index)++;

	if (tokenVectorEof(tokens, *index)) {
		perr("Expected a literal or '*', got <EOF>\n");

		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_LITERAL &&
	    tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_STAR) {
		perr("Expected a literal or '*', got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	ast_node_t* right = tokenToNode(tokenVectorPeek(tokens, *index));

	(*index)++;

	if (tokenVectorEof(tokens, *index)) {
		perr("Expected a ';', got <EOF>\n");

		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_SEMICOLON) {
		perr("Expected a ';', got ");

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));
		printf("\n");

		return newNode();
	}

	ast_parent_t* parent = newNode();

	parent->type = KSCRIPT_AST_NODE_TYPE_IMPORT;

	parent->left  = tokenToNode(file);
	parent->right = right;

	(*index)++;
	return parent;
}

ast_parent_t*
parserHandleFunction(token_vector_t* tokens, size_t* index, char* source_file) {
	// TODO
	(void)tokens;
	(void)index;
	(void)source_file;

	return newNode();
}

void parserPrintNodeType(ast_node_t* node) {
	if (!node) return;

	switch (node->type) {
		case KSCRIPT_AST_NODE_TYPE_DIV:
			printf("Div");
			break;
		case KSCRIPT_AST_NODE_TYPE_MUL:
			printf("Mul");
			break;
		case KSCRIPT_AST_NODE_TYPE_PLUS:
			printf("Plus");
			break;
		case KSCRIPT_AST_NODE_TYPE_MINUS:
			printf("Minus");
			break;

		case KSCRIPT_AST_NODE_TYPE_LOGICAL_AND:
			printf("LogicalAnd");
			break;
		case KSCRIPT_AST_NODE_TYPE_LOGICAL_OR:
			printf("LogicalOr");
			break;
		case KSCRIPT_AST_NODE_TYPE_LOGICAL_NOT:
			printf("LogicalNot");
			break;

		case KSCRIPT_AST_NODE_TYPE_BWISE_AND:
			printf("BitwiseAnd");
			break;
		case KSCRIPT_AST_NODE_TYPE_BWISE_OR:
			printf("BitwiseOr");
			break;
		case KSCRIPT_AST_NODE_TYPE_BWISE_XOR:
			printf("BitwiseXor");
			break;
		case KSCRIPT_AST_NODE_TYPE_BWISE_NOT:
			printf("BitwiseNot");
			break;

		case KSCRIPT_AST_NODE_TYPE_VARIABLE_DEC:
			printf("VariableDec");
			break;
		case KSCRIPT_AST_NODE_TYPE_VALUE_DEC:
			printf("ValueDec");
			break;
		case KSCRIPT_AST_NODE_TYPE_VARIABLE_DEF:
			printf("VariableDef");
			break;
		case KSCRIPT_AST_NODE_TYPE_VALUE_DEF:
			printf("ValueDef");
			break;
		case KSCRIPT_AST_NODE_TYPE_FUNCTION_DEC:
			printf("FunctionDec");
			break;
		case KSCRIPT_AST_NODE_TYPE_FUNCTION_DEF:
			printf("FunctionDef");
			break;

		case KSCRIPT_AST_NODE_TYPE_LITERAL:
			printf("Literal");
			break;
		case KSCRIPT_AST_NODE_TYPE_BOOL_LITERAL:
			printf("BoolLiteral");
			break;
		case KSCRIPT_AST_NODE_TYPE_STR_LITERAL:
			printf("StrLiteral");
			break;
		case KSCRIPT_AST_NODE_TYPE_INT_LITERAL:
			printf("IntLiteral");
			break;

		case KSCRIPT_AST_NODE_TYPE_STR_TYPE:
			printf("StrType");
			break;
		case KSCRIPT_AST_NODE_TYPE_PTR_TYPE:
			printf("PtrType");
			break;
		case KSCRIPT_AST_NODE_TYPE_BOOL_TYPE:
			printf("BoolType");
			break;

		case KSCRIPT_AST_NODE_TYPE_I64_TYPE:
			printf("i64");
			break;
		case KSCRIPT_AST_NODE_TYPE_I32_TYPE:
			printf("i32");
			break;
		case KSCRIPT_AST_NODE_TYPE_I16_TYPE:
			printf("i16");
			break;
		case KSCRIPT_AST_NODE_TYPE_I8_TYPE:
			printf("i8");
			break;

		case KSCRIPT_AST_NODE_TYPE_U64_TYPE:
			printf("u64");
			break;
		case KSCRIPT_AST_NODE_TYPE_U32_TYPE:
			printf("u32");
			break;
		case KSCRIPT_AST_NODE_TYPE_U16_TYPE:
			printf("u16");
			break;
		case KSCRIPT_AST_NODE_TYPE_U8_TYPE:
			printf("u8");
			break;

		case KSCRIPT_AST_NODE_TYPE_F64_TYPE:
			printf("f64");
			break;
		case KSCRIPT_AST_NODE_TYPE_F32_TYPE:
			printf("f32");
			break;
		case KSCRIPT_AST_NODE_TYPE_OPEN_PAREN:
			printf("OpenParentheses");
			break;
		case KSCRIPT_AST_NODE_TYPE_CLOSED_PAREN:
			printf("ClosedParentheses");
			break;
		case KSCRIPT_AST_NODE_TYPE_NONE:
			printf("None");
			break;
		case KSCRIPT_AST_NODE_TYPE_EOF:
			printf("EOF");
			break;
		case KSCRIPT_AST_NODE_TYPE_PROGRAM:
			printf("Program");
			break;

		case KSCRIPT_AST_NODE_TYPE_IMPORT:
			printf("Import");
			break;

		case KSCRIPT_AST_NODE_TYPE_FLOAT_LITERAL:
			printf("FloatLiteral");
			break;

		case KSCRIPT_AST_NODE_TYPE_ROOT:
			printf("Root");
			break;

		default:
			printf("???");
			break;
	}
}

void parserDumpNode(ast_parent_t* parent, size_t depth) {
	if (parent == NULL) return;

	for (size_t i = 0; i < depth; i++)
		printf("    ");

	printf("Node: ");
	parserPrintNodeType(parent);
	printf("\n");

	if (parent->line != 0) {
		for (size_t i = 0; i < depth; i++)
			printf("    ");
		printf("    - Line: %d\n", parent->line);
	}

	if (parent->line != 0) {
		for (size_t i = 0; i < depth; i++)
			printf("    ");
		printf("    - Column: %d\n", parent->col);
	}

	if (parent->lexeme && strcmp(parent->lexeme, "") != 0) {
		for (size_t i = 0; i < depth; i++)
			printf("    ");
		printf("    - Lexeme: \"%s\"\n", parent->lexeme);
	}

	parserDumpNode(parent->left, depth + 1);
	parserDumpNode(parent->right, depth + 1);

	if (!parent->block.body) return;

	for (size_t i = 0; i < depth; i++)
		printf("    ");
	printf("    - Block:\n");

	for (size_t i = 0; i < parent->block.size; i++) {
		parserDumpNode(parent->block.body[i], depth + 2);
	}

	if (!parent->function_dec.args) return;

	for (size_t i = 0; i < depth; i++)
		printf("    ");
	printf("    - Args:\n");

	for (size_t i = 0; i < parent->function_dec.args_size; i++) {
		parserDumpNode(parent->function_dec.args[i], depth + 2);
	}
}

ast_parent_t* parserParseFunctionParam(token_vector_t* tokens,
                                       size_t*         index,
                                       char*           source_file) {
	if (tokenVectorEof(tokens, *index)) {
		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_LITERAL) {
		errors_generated++;
		printf("%s:%d:%d: \033[1;31mSYNTAX ERROR\033[0m: Incomplete function "
		       "declaration. Expected a literal, got ",
		       source_file,
		       tokenVectorPeek(tokens, *index).line,
		       tokenVectorPeek(tokens, *index).col);

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));

		return newNode();
	}

	token_t literal = tokenVectorPeek(tokens, *index);

	(*index)++;

	if (tokenVectorEof(tokens, *index)) {
		errors_generated++;
		printf("%s:%d:%d: \033[1;31mSYNTAX ERROR\033[0m: Incomplete function "
		       "declaration. Expected ':', got EOF\n",
		       source_file,
		       tokenVectorPeek(tokens, *index).line,
		       tokenVectorPeek(tokens, *index).col);

		return newNode();
	}

	if (tokenVectorPeek(tokens, *index).type != KSCRIPT_TOKEN_TYPE_COLON) {
		errors_generated++;
		printf("%s:%d:%d: \033[1;31mSYNTAX ERROR\033[0m: Incomplete function "
		       "declaration. Expected ':', got ",
		       source_file,
		       tokenVectorPeek(tokens, *index).line,
		       tokenVectorPeek(tokens, *index).col);

		deTokenizeTokenKeyword(tokenVectorPeek(tokens, *index));

		return newNode();
	}

	(*index)++;

	if (tokenVectorEof(tokens, *index)) {
		errors_generated++;
		printf("%s:%d:%d: \033[1;31mSYNTAX ERROR\033[0m: Incomplete function "
		       "declaration. Expected a type, got EOF\n",
		       source_file,
		       tokenVectorPeek(tokens, *index).line,
		       tokenVectorPeek(tokens, *index).col);

		return newNode();
	}

	ast_node_t*   type = tokenToNode(tokenVectorPeek(tokens, *index));
	ast_parent_t* node = newNode();

	node->lexeme       = strdup(literal.s);
	node->lexeme_owned = true;

	node->right = type;

	return node;
}