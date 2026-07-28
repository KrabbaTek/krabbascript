/*
    KrabbaScript project
    File: tests/tests_main.c
    Purpose: Runs every test for the compiler
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#include <unity/unity.h>

void setUp() {}
void tearDown() {}

extern void runTokenizerTests();
extern void runHashTableTests();

int main() {
    UNITY_BEGIN();

    runTokenizerTests();
    runHashTableTests();

    return UNITY_END();
}