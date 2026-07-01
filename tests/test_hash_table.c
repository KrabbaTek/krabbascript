/*
    KrabbaScript project
    File: tests/test_tokenizer.c
    Purpose: Tests the hash table implementation
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
*/

#include <unity/unity.h>
#include <hash_table/hash_table.h>

void testHashTableLookUp() {
    hash_table_t* t = newHashTable();

    htable_entry_t* e = malloc(sizeof(htable_entry_t));

    e->key = "Krabba!";
    e->key_owned = false;
    e->value_owned = false;

    htableInsert(t, e);

    htable_entry_t* res = htableLookUp(t, "Krabba!");
    TEST_ASSERT_NOT_EQUAL(NULL, res);
    TEST_ASSERT_EQUAL_STRING(e->key, res->key);

    freeHashTable(t);
}

void testHashTableEntryNoLongerExists() {
    hash_table_t* t = newHashTable();

    htable_entry_t* e = malloc(sizeof(htable_entry_t));

    e->key = "Krabba!";
    e->key_owned = false;
    e->value_owned = false;

    htableInsert(t, e);

    htable_entry_t* res = htableLookUp(t, "Krabba!");
    TEST_ASSERT_NOT_EQUAL(NULL, res);
    TEST_ASSERT_EQUAL_STRING("Krabba!", res->key);

    htableDelete(t, "Krabba!");
    res = htableLookUp(t, "Krabba!");
    TEST_ASSERT_EQUAL(NULL, res);

    freeHashTable(t);
}

void runHashTableTests() {
    RUN_TEST(testHashTableLookUp);
    RUN_TEST(testHashTableEntryNoLongerExists);
}