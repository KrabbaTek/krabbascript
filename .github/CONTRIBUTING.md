
# Contributing to Krabbascript

## Table of contents
- [File header](#file-header)
- [Function naming](#function-naming)
- [Before commiting...](#before-commiting)
- [Header file structure](#header-file-structure)
- [C file structure](#c-file-structure)
- [Writing tests for your code](#writing-tests-for-your-code)
- [Commit messages](#commit-messages)
- [Writing changelog](#writing-changelog)

## File header
Every file should start with the following header:
```
/*
    KrabbaScript project
    File: src/path/to/file.c
    Purpose: Lorem Ipsum...
    Contributors:
        - person1 (person1@example.com)
*/
```

Where the purpose explains what the code in the file does. Write short, for example "A hash table implementation for the symbol table". Files with obvious purpose like main.c should only contain the contributors and the path to the file.

## Function naming
Every function should look like `addTwoNumbers` and have a comment on top that documents what it takes in, what it returns, and some additional info like does it handle NULLs. Example:
```
/*
 * Generates an FNV-1a hash from a null-terminated string.
 *  @param key  String to hash. Will generate an error with the ``err`` macro if it's NULL
 *  @return     Unsigned 64-bit FNV-1a hash value.
 */
uint64_t        hash(char* key);
...
```

## Before commiting...
Before you commit your changes, make sure that:

    - [ ] You have a test suite for your code and you tested it
    - [ ] You ran `make format`
    - [ ] You have no errors or warnings popping up when running `make`

`make format` uses clang-format to automatically format your code with the `.clang-format` in the project's root directory, worth checking out.

## Header file structure
Every header file has this structure:
```
/*
    KrabbaScript project
    File: include/example/example.h
    Contributors:
        - person1 (person1@example.com)
*/

#pragma once

#include <library.h>

#define DEFINITIONS

extern bool cooking;

void    makeBreakfast();
void    cookDinner();
```

Where functions start at the same column.
If the definition is a tweakable value, like `MAX_ENTRIES`, it should be wrapped around an ifndef and start with the `KSCRIPT_` prefix.
```
/*
    KrabbaScript project
    File: include/example/example.h
    Contributors:
        - person1 (person1@example.com)
*/

#pragma once

#include <library.h>
#include <stdio.h>

#ifndef KSCRIPT_MAX_CHEFS
#   define KSCRIPT_MAX_CHEFS 3
#endif

extern bool cooking;

void    makeBreakfast();
void    cookDinner();
```

## C file structure
Every C file follows this structure:
```
/*
    KrabbaScript project
    File: src/example/example.c
    Contributors:
        - person1 (person1@example.com)
*/

#include <example/example.h>

bool cooking;

void makeBreakfast() {
    printf("Making breakfast!\n");
}

void cookDinner() {
    printf("Cooking dinner!\n");
}
```

## Writing tests for your code
KrabbaScript uses [Unity](https://github.com/ThrowTheSwitch/Unity) for testing the project's code. All of the tests are located in `tests/` and are named `test_xyz.c`. The C file structure for the test looks like this:

```
/*
    KrabbaScript project
    File: tests/test_xyz.c
    Purpose: Tests xyz
    Contributors:
        - person1 (person1@example.com)
*/

#include <unity/unity.h>
#include <library/library.h>

void testXyz() {
    ...
}

void runXyzTests() {
    RUN_TEST(testXyz);
}
```

Where `testXyz` is a single test (like adding 2+2) and `runXyzTests` where it runs every single test.
After writing your test, make sure to register it in `tests/tests_main.c` with the extern like this:

```
/*
    KrabbaScript project
    File: tests/tests_main.c
    Purpose: Runs every test for the compiler
    Contributors:
        - person1 (person1@example.com)
*/

#include <unity/unity.h>

// Unity boilerplate that should be left blank
void setUp() {}
void tearDown() {}

...
extern void runXyz(); // Add your test here

int main() {
    UNITY_BEGIN();

    ...
    runXyz();

    return UNITY_END();
}
```

After you wrote your tests and added them to `tests_main.c`, you can now finally run the tests with `make test`. Check if the output of Unity doesn't have any failures before actually commiting.

## Commit messages
KrabbaScript's commit messages follow this style:
- `Update version from x.y.z to x.y.z`
- `Change the logo size from xyz to xyz`
- `Add xyz implementation for xyz`
- `Remove deprecated xyz macro`
- `Fix xyz memory leak`
- `Fix double free in xyz (closes #99)`

Commit messages must follow this format: `<Verb> <what> <optional detail>`. If you fixed an issue, add `(closes #xyz)` at the end where xyz is the issue number. This automatically closes the issue when pushed to the main branch.

## Writing changelog
Every time you modify something, you should document it in `CHANGELOG.md` under the current version.