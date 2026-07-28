/*
    KrabbaScript project
    File: src/tokenizer/tokenizer.c
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#include <tokenizer/tokenizer.h>

char_vector_t* fileToCharVector(const char* path) {
	FILE* file_ptr = fopen(path, "r");

	if (file_ptr == NULL) {
		perror("\033[1;31mERROR\033[0m: Failed to open file");
		exit(1);
	}

	char_vector_t* vector = newCharVector();

	int c;
	while ((c = fgetc(file_ptr)) != EOF) {
		charVectorPush(vector, (char)c);
	}

	fclose(file_ptr);
	return vector;
}

void deTokenize(token_vector_t* vector) {
	for (size_t i = 0; i < vector->size; i++) {
		deTokenizeToken(vector->data[i]);
	}
}

void deTokenizeToken(token_t token) {
	switch (token.type) {
		case KSCRIPT_TOKEN_TYPE_I64:
			printf("i64 ");
			break;
		case KSCRIPT_TOKEN_TYPE_I32:
			printf("i32 ");
			break;
		case KSCRIPT_TOKEN_TYPE_I16:
			printf("i16 ");
			break;
		case KSCRIPT_TOKEN_TYPE_I8:
			printf("i8 ");
			break;
		case KSCRIPT_TOKEN_TYPE_U64:
			printf("u64 ");
			break;
		case KSCRIPT_TOKEN_TYPE_U32:
			printf("u32 ");
			break;
		case KSCRIPT_TOKEN_TYPE_U16:
			printf("u16 ");
			break;
		case KSCRIPT_TOKEN_TYPE_U8:
			printf("u8 ");
			break;
		case KSCRIPT_TOKEN_TYPE_F64:
			printf("f64 ");
			break;
		case KSCRIPT_TOKEN_TYPE_F32:
			printf("f32 ");
			break;
		case KSCRIPT_TOKEN_TYPE_PTR:
			printf("ptr ");
			break;
		case KSCRIPT_TOKEN_TYPE_STR:
			printf("str ");
			break;
		case KSCRIPT_TOKEN_TYPE_BOOL:
			printf("bool ");
			break;

		case KSCRIPT_TOKEN_TYPE_LITERAL:
			printf("Literal: %s ", token.s);
			break;
		case KSCRIPT_TOKEN_TYPE_INT_LITERAL:
			printf("Int literal: %d ", token.i);
			break;
		case KSCRIPT_TOKEN_TYPE_FLOAT_LITERAL:
			printf("Float literal: %f", token.d);
			break;
		case KSCRIPT_TOKEN_TYPE_STR_LITERAL:
			printf("String literal: \"%s\" ", token.s);
			break;
		case KSCRIPT_TOKEN_TYPE_BOOL_LITERAL:
			printf("Bool literal: %s", token.b == true ? "true" : "false");
			break;

		case KSCRIPT_TOKEN_TYPE_FUNCTION:
			printf("function ");
			break;
		case KSCRIPT_TOKEN_TYPE_VAR:
			printf("var ");
			break;
		case KSCRIPT_TOKEN_TYPE_VAL:
			printf("val ");
			break;
		case KSCRIPT_TOKEN_TYPE_IMPORT:
			printf("import ");
			break;
		case KSCRIPT_TOKEN_TYPE_FROM:
			printf("from ");
			break;
		case KSCRIPT_TOKEN_TYPE_MODULE:
			printf("module ");
			break;
		case KSCRIPT_TOKEN_TYPE_WHILE:
			printf("while ");
			break;
		case KSCRIPT_TOKEN_TYPE_TRUE:
			printf("true ");
			break;
		case KSCRIPT_TOKEN_TYPE_FALSE:
			printf("false ");
			break;
		case KSCRIPT_TOKEN_TYPE_FOR:
			printf("for ");
			break;

		case KSCRIPT_TOKEN_TYPE_OPEN_PAREN:
			printf("( ");
			break;
		case KSCRIPT_TOKEN_TYPE_CLOSED_PAREN:
			printf(") ");
			break;
		case KSCRIPT_TOKEN_TYPE_SEMICOLON:
			printf("; ");
			break;
		case KSCRIPT_TOKEN_TYPE_EQUALS:
			printf("= ");
			break;
		case KSCRIPT_TOKEN_TYPE_EXMARK:
			printf("! ");
			break;
		case KSCRIPT_TOKEN_TYPE_QMARK:
			printf("? ");
			break;
		case KSCRIPT_TOKEN_TYPE_DOT:
			printf(". ");
			break;
		case KSCRIPT_TOKEN_TYPE_COMMA:
			printf(", ");
			break;
		case KSCRIPT_TOKEN_TYPE_CBRACKET_OPEN:
			printf("{ ");
			break;
		case KSCRIPT_TOKEN_TYPE_CBRACKET_CLOSED:
			printf("} ");
			break;
		case KSCRIPT_TOKEN_TYPE_LESS_THAN:
			printf("< ");
			break;
		case KSCRIPT_TOKEN_TYPE_MORE_THAN:
			printf("> ");
			break;
		case KSCRIPT_TOKEN_TYPE_STAR:
			printf("* ");
			break;
		case KSCRIPT_TOKEN_TYPE_MINUS:
			printf("- ");
			break;
		case KSCRIPT_TOKEN_TYPE_COLON:
			printf(": ");
			break;
		case KSCRIPT_TOKEN_TYPE_SBRACKET_OPEN:
			printf("[ ");
			break;
		case KSCRIPT_TOKEN_TYPE_SBRACKET_CLOSED:
			printf("] ");
			break;
		case KSCRIPT_TOKEN_TYPE_SLASH:
			printf("/ ");
			break;
		case KSCRIPT_TOKEN_TYPE_BSLASH:
			printf("\\ ");
			break;
		case KSCRIPT_TOKEN_TYPE_IF:
			printf("if ");
			break;
		case KSCRIPT_TOKEN_TYPE_ELIF:
			printf("elif ");
			break;
		case KSCRIPT_TOKEN_TYPE_ELSE:
			printf("else ");
			break;
		case KSCRIPT_TOKEN_TYPE_ARROW:
			printf("-> ");
			break;
		case KSCRIPT_TOKEN_TYPE_EQUALS_EQUALS:
			printf("== ");
			break;

		case KSCRIPT_TOKEN_TYPE_ARRAY:
			printf("array ");
			break;
		case KSCRIPT_TOKEN_TYPE_PLUS:
			printf("+ ");
			break;

		case KSCRIPT_TOKEN_TYPE_PLUS_EQ:
			printf("+= ");
			break;
		case KSCRIPT_TOKEN_TYPE_MINUS_EQ:
			printf("-= ");
			break;
		case KSCRIPT_TOKEN_TYPE_MUL_EQ:
			printf("*= ");
			break;
		case KSCRIPT_TOKEN_TYPE_DIV_EQ:
			printf("/= ");
			break;

		case KSCRIPT_TOKEN_TYPE_PLUS_PLUS:
			printf("++ ");
			break;
		case KSCRIPT_TOKEN_TYPE_MINUS_MINUS:
			printf("-- ");
			break;

		case KSCRIPT_TOKEN_TYPE_EOF:
			printf("<EOF> ");
			break;

		case KSCRIPT_TOKEN_TYPE_LOGICAL_AND:
			printf("&& ");
			break;
		case KSCRIPT_TOKEN_TYPE_LOGICAL_OR:
			printf("|| ");
			break;

		case KSCRIPT_TOKEN_TYPE_BWISE_XOR:
			printf("^ ");
			break;
		case KSCRIPT_TOKEN_TYPE_BWISE_AND:
			printf("& ");
			break;

		case KSCRIPT_TOKEN_TYPE_BWISE_OR:
			printf("| ");
			break;
		case KSCRIPT_TOKEN_TYPE_POWER_OF:
			printf("** ");
			break;

		case KSCRIPT_TOKEN_TYPE_GREATER_EQU:
			printf(">=");
			break;

		case KSCRIPT_TOKEN_TYPE_LESS_EQU:
			printf("<= ");
			break;

		case KSCRIPT_TOKEN_TYPE_NOT_EQU:
			printf("!= ");
			break;

		case KSCRIPT_TOKEN_TYPE_PUB:
			printf("pub ");
			break;

		case KSCRIPT_TOKEN_TYPE_RETURN:
			printf("return ");
			break;

		default:
			printf("Unknown token: %d", token.type);
			break;
	}
	printf("\n");
}

void deTokenizeTokenKeyword(token_t token) {
	switch (token.type) {
		case KSCRIPT_TOKEN_TYPE_I64:
			printf("i64");
			break;
		case KSCRIPT_TOKEN_TYPE_I32:
			printf("i32");
			break;
		case KSCRIPT_TOKEN_TYPE_I16:
			printf("i16");
			break;
		case KSCRIPT_TOKEN_TYPE_I8:
			printf("i8");
			break;
		case KSCRIPT_TOKEN_TYPE_U64:
			printf("u64");
			break;
		case KSCRIPT_TOKEN_TYPE_U32:
			printf("u32");
			break;
		case KSCRIPT_TOKEN_TYPE_U16:
			printf("u16");
			break;
		case KSCRIPT_TOKEN_TYPE_U8:
			printf("u8");
			break;
		case KSCRIPT_TOKEN_TYPE_F64:
			printf("f64");
			break;
		case KSCRIPT_TOKEN_TYPE_F32:
			printf("f32");
			break;
		case KSCRIPT_TOKEN_TYPE_PTR:
			printf("ptr");
			break;
		case KSCRIPT_TOKEN_TYPE_STR:
			printf("str");
			break;
		case KSCRIPT_TOKEN_TYPE_BOOL:
			printf("bool");
			break;

		case KSCRIPT_TOKEN_TYPE_LITERAL:
			printf("%s", token.s);
			break;
		case KSCRIPT_TOKEN_TYPE_INT_LITERAL:
			printf("%d", token.i);
			break;
		case KSCRIPT_TOKEN_TYPE_STR_LITERAL:
			printf("\"%s\"", token.s);
			break;
		case KSCRIPT_TOKEN_TYPE_BOOL_LITERAL:
			printf("%s", token.b == true ? "true" : "false");
			break;
		case KSCRIPT_TOKEN_TYPE_FLOAT_LITERAL:
			printf("%f", token.d);
			break;

		case KSCRIPT_TOKEN_TYPE_FUNCTION:
			printf("func");
			break;
		case KSCRIPT_TOKEN_TYPE_VAR:
			printf("var");
			break;
		case KSCRIPT_TOKEN_TYPE_VAL:
			printf("val");
			break;
		case KSCRIPT_TOKEN_TYPE_IMPORT:
			printf("import");
			break;
		case KSCRIPT_TOKEN_TYPE_FROM:
			printf("from");
			break;
		case KSCRIPT_TOKEN_TYPE_MODULE:
			printf("module");
			break;
		case KSCRIPT_TOKEN_TYPE_WHILE:
			printf("while");
			break;
		case KSCRIPT_TOKEN_TYPE_TRUE:
			printf("true");
			break;
		case KSCRIPT_TOKEN_TYPE_FALSE:
			printf("false");
			break;
		case KSCRIPT_TOKEN_TYPE_FOR:
			printf("for");
			break;

		case KSCRIPT_TOKEN_TYPE_OPEN_PAREN:
			printf("(");
			break;
		case KSCRIPT_TOKEN_TYPE_CLOSED_PAREN:
			printf(")");
			break;
		case KSCRIPT_TOKEN_TYPE_SEMICOLON:
			printf(";");
			break;
		case KSCRIPT_TOKEN_TYPE_EQUALS:
			printf("=");
			break;
		case KSCRIPT_TOKEN_TYPE_EXMARK:
			printf("!");
			break;
		case KSCRIPT_TOKEN_TYPE_QMARK:
			printf("?");
			break;
		case KSCRIPT_TOKEN_TYPE_DOT:
			printf(".");
			break;
		case KSCRIPT_TOKEN_TYPE_COMMA:
			printf(",");
			break;
		case KSCRIPT_TOKEN_TYPE_CBRACKET_OPEN:
			printf("{");
			break;
		case KSCRIPT_TOKEN_TYPE_CBRACKET_CLOSED:
			printf("}");
			break;
		case KSCRIPT_TOKEN_TYPE_LESS_THAN:
			printf("<");
			break;
		case KSCRIPT_TOKEN_TYPE_MORE_THAN:
			printf(">");
			break;
		case KSCRIPT_TOKEN_TYPE_STAR:
			printf("*");
			break;
		case KSCRIPT_TOKEN_TYPE_MINUS:
			printf("-");
			break;
		case KSCRIPT_TOKEN_TYPE_COLON:
			printf(":");
			break;
		case KSCRIPT_TOKEN_TYPE_SBRACKET_OPEN:
			printf("[");
			break;
		case KSCRIPT_TOKEN_TYPE_SBRACKET_CLOSED:
			printf("]");
			break;
		case KSCRIPT_TOKEN_TYPE_SLASH:
			printf("/");
			break;
		case KSCRIPT_TOKEN_TYPE_BSLASH:
			printf("\\");
			break;
		case KSCRIPT_TOKEN_TYPE_IF:
			printf("if");
			break;
		case KSCRIPT_TOKEN_TYPE_ELIF:
			printf("elif");
			break;
		case KSCRIPT_TOKEN_TYPE_ELSE:
			printf("else");
			break;
		case KSCRIPT_TOKEN_TYPE_ARROW:
			printf("->");
			break;
		case KSCRIPT_TOKEN_TYPE_EQUALS_EQUALS:
			printf("==");
			break;

		case KSCRIPT_TOKEN_TYPE_ARRAY:
			printf("array");
			break;
		case KSCRIPT_TOKEN_TYPE_PLUS:
			printf("+");
			break;

		case KSCRIPT_TOKEN_TYPE_PLUS_EQ:
			printf("+=");
			break;
		case KSCRIPT_TOKEN_TYPE_MINUS_EQ:
			printf("-=");
			break;
		case KSCRIPT_TOKEN_TYPE_MUL_EQ:
			printf("*=");
			break;
		case KSCRIPT_TOKEN_TYPE_DIV_EQ:
			printf("/=");
			break;

		case KSCRIPT_TOKEN_TYPE_PLUS_PLUS:
			printf("++");
			break;
		case KSCRIPT_TOKEN_TYPE_MINUS_MINUS:
			printf("--");
			break;

		case KSCRIPT_TOKEN_TYPE_EOF:
			printf("EOF");
			break;

		case KSCRIPT_TOKEN_TYPE_LOGICAL_AND:
			printf("&&");
			break;
		case KSCRIPT_TOKEN_TYPE_LOGICAL_OR:
			printf("||");
			break;

		case KSCRIPT_TOKEN_TYPE_BWISE_XOR:
			printf("^");
			break;
		case KSCRIPT_TOKEN_TYPE_BWISE_AND:
			printf("&");
			break;

		case KSCRIPT_TOKEN_TYPE_BWISE_OR:
			printf("|");
			break;
		case KSCRIPT_TOKEN_TYPE_POWER_OF:
			printf("**");
			break;
		case KSCRIPT_TOKEN_TYPE_PUB:
			printf("pub");
			break;

		case KSCRIPT_TOKEN_TYPE_RETURN:
			printf("return");
			break;

		default:
			printf("??");
			break;
	}
	printf("\n");
}

token_vector_t* tokenize(char_vector_t* vector, const char* debug_path) {

	char_vector_t*  buffer = newCharVector();
	token_vector_t* tokens = newTokenVector();

	int line = 1;
	int col  = 1;
	for (size_t i = 0; i < vector->size; i++) {
		char current = vector->data[i];

		if (isdigit(current)) {
			while (i < vector->size && isdigit(current)) {
				charVectorPush(buffer, current);

				i++;
				col++;
				current = vector->data[i];
			}

			if (i < vector->size && current == '.') {
				charVectorPush(buffer, '.');
				i++;
				col++;
				current = vector->data[i];

				while (i < vector->size && isdigit(current)) {
					charVectorPush(buffer, current);

					i++;
					col++;
					current = (i < vector->size) ? vector->data[i] : '\0';
				}

				charVectorPush(buffer, '\0');
				i--;

				errno        = 0;
				double value = strtod(buffer->data, NULL);
				if (errno == ERANGE) {
					errors_generated++;
					printf("%s:%d:%d: \033[1;31mOVERFLOW ERROR\033[0m: Float "
					       "literal is "
					       "too large\n",
					       debug_path,
					       line,
					       col);
				}

				resetCharVector(buffer);
				tokenVectorPush(
				        tokens,
				        (token_t){.type = KSCRIPT_TOKEN_TYPE_FLOAT_LITERAL,
				                  .d    = value,
				                  .line = line,
				                  .col  = col});
			} else {
				// Not a float
				charVectorPush(buffer, '\0');
				i--;

				errno  = 0;
				long v = strtol(buffer->data, NULL, 10);
				if (errno == ERANGE || v > INT_MAX || v < INT_MIN) {
					errors_generated++;
					printf("%s:%d:%d: \033[1;31mOVERFLOW ERROR\033[0m: Integer "
					       "literal is "
					       "too large\n",
					       debug_path,
					       line,
					       col);
				}

				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_INT_LITERAL,
				                          .i    = (int)v,
				                          .line = line,
				                          .col  = col});
			}
		} else if (current == '(') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_OPEN_PAREN,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == ')') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_CLOSED_PAREN,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == ';') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_SEMICOLON,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '=') {
			if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(
				        tokens,
				        (token_t){.type = KSCRIPT_TOKEN_TYPE_EQUALS_EQUALS,
				                  .line = line,
				                  .col  = col});
				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_EQUALS,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '!') {
			if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_NOT_EQU,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_EXMARK,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '?') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_QMARK,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '.') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_DOT,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == ',') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_COMMA,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '{') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_CBRACKET_OPEN,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '}') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_CBRACKET_CLOSED,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '<') {
			if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_LESS_EQU,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			} else if (tokens->size > 0 &&
			           tokenVectorPeek(tokens, tokens->size - 1).type ==
			                   KSCRIPT_TOKEN_TYPE_FROM) {
				resetCharVector(buffer);
				i++;
				col++;

				while (i < vector->size && vector->data[i] != '>') {
					charVectorPush(buffer, vector->data[i]);
					i++;
					col++;
				}

				col++;
				charVectorPush(buffer, '\0');

				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_STR_LITERAL,
				                          .s       = strdup(buffer->data),
				                          .line    = line,
				                          .col     = col,
				                          .s_owned = true});
				resetCharVector(buffer);
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_LESS_THAN,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '>') {
			if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_GREATER_EQU,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_MORE_THAN,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '"') {
			resetCharVector(buffer);
			i++;
			col++;

			while (i < vector->size) {
				current = vector->data[i];

				if (current == '\\' && i + 1 < vector->size) {
					char next = vector->data[i + 1];
					if (next == '"') {
						charVectorPush(buffer, '"');
						i += 2;
						col += 2;
						continue;
					}
				}

				if (current == '"') break;

				charVectorPush(buffer, current);
				if (current == '\n') {
					line++;
					col = 1;
				} else
					col++;
				i++;
			}

			charVectorPush(buffer, '\0');
			tokenVectorPush(tokens,
			                (token_t){.type    = KSCRIPT_TOKEN_TYPE_STR_LITERAL,
			                          .s       = strdup(buffer->data),
			                          .line    = line,
			                          .col     = col,
			                          .s_owned = true});

			resetCharVector(buffer);
		} else if (current == '*') {
			if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_MUL_EQ,
				                          .line = line,
				                          .col  = col});

				i++;
				col += 2;
			}
			if (i + 1 < vector->size && vector->data[i + 1] == '*') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_POWER_OF,
				                          .line = line,
				                          .col  = col});

				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_STAR,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '-') {
			// Handle ->
			if (i + 1 < vector->size && vector->data[i + 1] == '>') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_ARROW,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			}
			// Handle -=
			else if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_MINUS_EQ,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			}
			// Handle --
			else if (i + 1 < vector->size && vector->data[i + 1] == '-') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_MINUS_MINUS,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_MINUS,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '^') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_BWISE_XOR,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '|') {
			if (i + 1 < vector->size && vector->data[i + 1] == '|') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_LOGICAL_OR,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_BWISE_OR,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '&') {
			if (i + 1 < vector->size && vector->data[i + 1] == '&') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_LOGICAL_AND,
				                          .line = line,
				                          .col  = col});
				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_BWISE_AND,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == ':') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_COLON,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '[') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_SBRACKET_OPEN,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == ']') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_SBRACKET_CLOSED,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '/') {
			// Handle /=
			if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_DIV_EQ,
				                          .line = line,
				                          .col  = col});
				i++;

				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_SLASH,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (current == '\\') {
			tokenVectorPush(tokens,
			                (token_t){.type = KSCRIPT_TOKEN_TYPE_BSLASH,
			                          .line = line,
			                          .col  = col});
			col++;
		} else if (current == '#') {
			while (i < vector->size && current != '\n') {
				i++;
				col++;

				current = (i < vector->size) ? vector->data[i] : '\0';
			}

			if (i < vector->size && vector->data[i] == '\n') {
				line++;
				col = 1;
			}
		} else if (current == '+') {
			// Handle +=
			if (i + 1 < vector->size && vector->data[i + 1] == '=') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_PLUS_EQ,
				                          .line = line,
				                          .col  = col});

				i++;
				col += 2;
			}
			// Handle ++
			else if (i + 1 < vector->size && vector->data[i + 1] == '+') {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_PLUS_PLUS,
				                          .line = line,
				                          .col  = col});

				i++;
				col += 2;
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_PLUS,
				                          .line = line,
				                          .col  = col});
				col++;
			}
		} else if (isalnum(current) || current == '_') {
			while (i < vector->size && (isalnum(current) || current == '_')) {
				charVectorPush(buffer, current);

				i++;
				col++;
				current = (i < vector->size) ? vector->data[i] : '\0';
			}
			charVectorPush(buffer, '\0');
			i--;

			if (strcmp(buffer->data, "func") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_FUNCTION,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "var") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_VAR,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "val") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_VAL,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "import") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_IMPORT,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "from") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_FROM,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "module") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_MODULE,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "while") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_WHILE,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "true") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(
				        tokens,
				        (token_t){.type = KSCRIPT_TOKEN_TYPE_BOOL_LITERAL,
				                  .b    = true,
				                  .line = line,
				                  .col  = col});
			} else if (strcmp(buffer->data, "false") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(
				        tokens,
				        (token_t){.type = KSCRIPT_TOKEN_TYPE_BOOL_LITERAL,
				                  .b    = false,
				                  .line = line,
				                  .col  = col});
			} else if (strcmp(buffer->data, "for") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_FOR,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "i64") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_I64,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "i32") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_I32,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "i16") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_I16,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "i8") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_I8,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "u64") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_U64,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "u32") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_U32,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "u16") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_U16,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "u8") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_U8,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "f64") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_F64,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "f32") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_F32,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "ptr") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_PTR,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "str") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_STR,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "bool") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_BOOL,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "array") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_ARRAY,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "if") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_IF,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "elif") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_ELIF,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "else") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_ELSE,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "pub") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_PUB,
				                          .line = line,
				                          .col  = col});
			} else if (strcmp(buffer->data, "return") == 0) {
				resetCharVector(buffer);
				tokenVectorPush(tokens,
				                (token_t){.type = KSCRIPT_TOKEN_TYPE_RETURN,
				                          .line = line,
				                          .col  = col});
			} else {
				tokenVectorPush(tokens,
				                (token_t){.type    = KSCRIPT_TOKEN_TYPE_LITERAL,
				                          .s       = strdup(buffer->data),
				                          .line    = line,
				                          .col     = col,
				                          .s_owned = true});
				resetCharVector(buffer);
			}
		} else if (isspace(current)) {
			if (current == '\n') {
				line++;
				col = 1;
			}
			continue;
		} else {
			errors_generated++;

			printf("%s:%d:%d: \033[1;31mLEXICAL ERROR\033[0m: Unexpected "
			       "character '%c'\n",
			       debug_path,
			       line,
			       col,
			       current);

			unsigned char b = (unsigned char)current;
			int           len;

			if ((b & 0x80) == 0x0)
				len = 1;
			else if ((b & 0xE0) == 0xC0)
				len = 2;
			else if ((b & 0xF0) == 0xE0)
				len = 3;
			else if ((b & 0xF8) == 0xF0)
				len = 4;
			else
				len = 1;

			// Skip the utf-8 characters
			i += len - 1;
			col++;
		}
	}
	tokenVectorPush(tokens,
	                (token_t){.type = KSCRIPT_TOKEN_TYPE_EOF,
	                          .line = line,
	                          .col  = col});

	freeCharVector(buffer);
	freeCharVector(vector);

	return tokens;
}
