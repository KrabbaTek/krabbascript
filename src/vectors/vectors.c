#include <vectors/vectors.h>

char_vector_t* newCharVector() {
    char_vector_t* vec = (char_vector_t*)malloc(sizeof(char_vector_t));

    if (!vec) {
        printf("\033[1;31mALLOCATION ERRO\033[0m: Failed to allocate memory "
               "for a char_vector\n");
        exit(1);
    }

    vec->data = NULL;

    vec->size     = 0;
    vec->capacity = 0;

    return vec;
}

char_vector_t* charVectorFromString(const char* string) {
    char_vector_t* vec = newCharVector();

    while (*string != '\0') {
        char c = *string;
        charVectorPush(vec, c);
        string++;
    }

    return vec;
}

void resetCharVector(char_vector_t* vector) {
    vector->size = 0;
}

void freeCharVector(char_vector_t* vector) {
    if (!vector) return;

    if (!vector->data) {
        free(vector);
        return;
    } else {
        free(vector->data);
        free(vector);
    }
}

void charVectorPush(char_vector_t* vector, char val) {
    if (vector->data == NULL) {
        vector->data = (char*)malloc(sizeof(char));

        vector->capacity = 1;
    } else if (vector->size >= vector->capacity) {
        vector->capacity *= 2;
        vector->data =
                (char*)realloc(vector->data, sizeof(char) * vector->capacity);
    }

    vector->data[vector->size] = val;
    vector->size++;
}

char charVectorPop(char_vector_t* vector) {
    if (vector->size > 0) {
        char c = vector->data[vector->size - 1];

        vector->size--;
        return c;
    }

    return '\0';
}

char charVectorPeek(char_vector_t* vector, size_t index) {
    if (index > vector->size) return '\0';

    return vector->data[index];
}

token_vector_t* newTokenVector() {
    token_vector_t* vec = (token_vector_t*)malloc(sizeof(token_vector_t));
    if (!vec) {
        err("Failed to allocate a token vector");
        exit(1);
    }

    vec->data = NULL;

    vec->size     = 0;
    vec->capacity = 0;

    return vec;
}

void tokenVectorPush(token_vector_t* vector, token_t val) {
    if (!vector->data) {
        vector->data = (token_t*)malloc(sizeof(token_t));

        vector->capacity = 1;
    } else if (vector->size >= vector->capacity) {
        vector->capacity *= 2;
        vector->data = (token_t*)realloc(vector->data,
                                         sizeof(token_t) * vector->capacity);
    }

    vector->data[vector->size] = val;
    vector->size++;
}

token_t tokenVectorPop(token_vector_t* vector) {
    if (vector->size > 0) {
        token_t t = vector->data[vector->size - 1];

        vector->size--;
        return t;
    }

    return (token_t){.type = KSCRIPT_TOKEN_TYPE_EOF};
}

token_t tokenVectorPeek(token_vector_t* vector, size_t index) {
    if (index >= vector->size) {
        return (token_t){.type = KSCRIPT_TOKEN_TYPE_EOF};
    }

    return vector->data[index];
}

void freeTokenVector(token_vector_t* vector) {
    if (!vector) return;

    if (!vector->data) {
        free(vector);
        return;
    }
    // Free every token
    size_t i = 0;
    while (!tokenVectorEof(vector, i)) {
        freeToken(&vector->data[i]);
        i++;
    }

    free(vector->data);
    free(vector);
}

void freeToken(token_t* token) {
    if (token && token->s_owned) {
        free(token->s);
        token->s       = NULL;
        token->s_owned = false;
    }
}

bool tokenVectorEof(token_vector_t* vector, size_t index) {
    if (tokenVectorPeek(vector, index).type == KSCRIPT_TOKEN_TYPE_EOF) {
        return true;
    }

    return false;
}

ast_node_t* newNode() {
    ast_node_t* node = (ast_node_t*)calloc(1, sizeof(ast_node_t));

    return node;
}

void freeNode(ast_node_t* node) {
    if (!node) return;
    if (node->lexeme && node->lexeme_owned) free(node->lexeme);

    freeNode(node->left);
    freeNode(node->right);

    free(node);
}

symbol_table_t* newSymbolTable() {
    symbol_table_t* t = (symbol_table_t*)malloc(sizeof(symbol_table_t));

    t->entries = NULL;

    t->size     = 0;
    t->capacity = 0;

    return t;
}

void freeSymbolTable(symbol_table_t* table) {
    if (!table) return;

    if (!table->entries) {
        free(table);
        return;
    } else {
        free(table->entries);
        free(table);

        return;
    }
}

void symbolTablePush(symbol_table_t* table, st_entry_t entry) {
    if (table->entries == NULL) {
        table->entries = (st_entry_t*)malloc(sizeof(st_entry_t));

        table->capacity = 1;
    } else if (table->size >= table->capacity) {
        table->capacity *= 2;
        table->entries =
                (st_entry_t*)realloc(table->entries,
                                     sizeof(st_entry_t) * table->capacity);
    }

    table->entries[table->size] = entry;
    table->size++;
}

st_entry_t symbolTablePop(symbol_table_t* table) {
    if (table->size > 0) {
        st_entry_t t = table->entries[table->size - 1];

        table->size--;
        return t;
    }

    return (st_entry_t){.type = KSCRIPT_ST_ENTRY_TYPE_EOF};
}

st_entry_t symbolTablePeek(symbol_table_t* table, size_t index) {
    if (index > table->size) {
        return (st_entry_t){.type = KSCRIPT_ST_ENTRY_TYPE_EOF};
    }

    return table->entries[index];
}

void newNodeBlock(ast_node_t* node) {
    node->block.body = NULL;

    node->block.size     = 0;
    node->block.capacity = 0;
}

void freeNodeBlock(ast_node_t* node) {
    if (!node->block.body) return;

    for (size_t i = 0; i < node->block.size; i++) {
        freeNode(node->block.body[i]);
    }

    free(node->block.body);

    node->block.body     = NULL;
    node->block.size     = 0;
    node->block.capacity = 0;
}

void nodeBlockPush(ast_node_t* node, ast_node_t* child) {
    if (node->block.body == NULL) {
        node->block.body     = (ast_node_t**)malloc(sizeof(ast_node_t*));
        node->block.capacity = 1;
    } else if (node->block.size >= node->block.capacity) {
        node->block.capacity *= 2;
        node->block.body = (ast_node_t**)realloc(node->block.body,
                                                 sizeof(ast_node_t*) *
                                                         node->block.capacity);
    }

    node->block.body[node->block.size] = child;
    node->block.size++;
}

ast_node_t* nodeBlockPop(ast_node_t* node) {
    if (node->block.size > 0) {
        ast_node_t* child = node->block.body[node->block.size - 1];
        node->block.size--;
        return child;
    }

    return NULL;
}

void resetTokenVector(token_vector_t* vector) {
    vector->size = 0;
}

void newFuncArgs(ast_node_t* node) {
    node->function_dec.args = NULL;

    node->function_dec.args_size     = 0;
    node->function_dec.args_capacity = 0;
}

void freeFuncArgs(ast_node_t* node) {
    if (!node->function_dec.args) return;

    for (size_t i = 0; i < node->function_dec.args_size; i++) {
        freeNode(node->function_dec.args[i]);
    }

    free(node->function_dec.args);

    node->function_dec.args          = NULL;
    node->function_dec.args_size     = 0;
    node->function_dec.args_capacity = 0;
}

void funcArgsPush(ast_node_t* node, ast_node_t* child) {
    if (node->function_dec.args == NULL) {
        node->function_dec.args = (ast_node_t**)malloc(sizeof(ast_node_t*));
        node->function_dec.args_capacity = 1;
    } else if (node->function_dec.args_size >=
               node->function_dec.args_capacity) {
        node->function_dec.args_capacity *= 2;
        node->function_dec.args =
                (ast_node_t**)realloc(node->function_dec.args,
                                      sizeof(ast_node_t*) *
                                              node->function_dec.args_capacity);
    }

    node->function_dec.args[node->function_dec.args_size] = child;
    node->function_dec.args_size++;
}

ast_node_t* funcArgskPop(ast_node_t* node) {
    if (node->function_dec.args_size > 0) {
        ast_node_t* child =
                node->function_dec.args[node->function_dec.args_size - 1];
        node->function_dec.args_size--;
        return child;
    }

    return NULL;
}