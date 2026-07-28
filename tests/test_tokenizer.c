/*
    KrabbaScript project
    File: tests/test_tokenizer.c
    Purpose: Tests the tokenizer
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#define TOSTRING(x) #x
#define STRINGIFY(x) TOSTRING(x)

#include <unity/unity.h>
#include <tokenizer/tokenizer.h>

void testTokenizerStringLiteral() {
    char_vector_t* v = charVectorFromString("\"Krabba!\"");
    token_vector_t* t = tokenize(v, "test");

    TEST_ASSERT_EQUAL_INT(KSCRIPT_TOKEN_TYPE_STR_LITERAL, tokenVectorPeek(t, 0).type);
    TEST_ASSERT_EQUAL_STRING("Krabba!", tokenVectorPeek(t, 0).s);

    freeTokenVector(t);
}

void testTokenizerStringLiteralNewLine() {
    char_vector_t* v = charVectorFromString("\"Krabba!\n\"");
    token_vector_t* t = tokenize(v, "test");

    TEST_ASSERT_EQUAL_INT(KSCRIPT_TOKEN_TYPE_STR_LITERAL, tokenVectorPeek(t, 0).type);
    TEST_ASSERT_EQUAL_STRING("Krabba!\n", tokenVectorPeek(t, 0).s);

    freeTokenVector(t);
}

void testTokenizerIntLiteral() {
    char_vector_t* v = charVectorFromString("500");
    token_vector_t* t = tokenize(v, "test");

    TEST_ASSERT_EQUAL_INT(KSCRIPT_TOKEN_TYPE_INT_LITERAL, tokenVectorPeek(t, 0).type);
    TEST_ASSERT_EQUAL_INT(500, tokenVectorPeek(t, 0).i);

    freeTokenVector(t);
}

void testTokenizerFloatLiteral() {
    char_vector_t* v = charVectorFromString("3.14");
    token_vector_t* t = tokenize(v, "test");

    TEST_ASSERT_EQUAL_INT(KSCRIPT_TOKEN_TYPE_FLOAT_LITERAL, tokenVectorPeek(t, 0).type);
    TEST_ASSERT_DOUBLE_WITHIN(0.0001, 3.14, tokenVectorPeek(t, 0).d);

    freeTokenVector(t);
}

// Will be called by tests_main.c
void runTokenizerTests() {
    RUN_TEST(testTokenizerStringLiteral);
    RUN_TEST(testTokenizerStringLiteralNewLine);
    RUN_TEST(testTokenizerIntLiteral);
    RUN_TEST(testTokenizerFloatLiteral);
}