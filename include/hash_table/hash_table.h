/*
    KrabbaScript project
    File: include/hash_table/hash_table.h
    Purpose: A hash table implementation for the symbol table
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
        - Visser13
*/

#pragma once

#include <common/common.h>
#include <stdint.h>
#include <string.h>
#include <types/types.h>

/*
 * Generates an FNV-1a hash from a null-terminated string.
 *  @param key  String to hash. Will generate an error with the ``err`` macro if
 * it's NULL
 *  @return     Unsigned 64-bit FNV-1a hash value.
 */
uint64_t hash(char* key);

/*
 * Creates a new hash table and automatically zeroes out the bucket
 *  @return     Zeroed out hash table. If failed to allocate memory, returns
 * NULL
 */
hash_table_t* newHashTable();
/*
 * Frees a hash table alongside its contents. Will return if the table is NULL
 *  @param htable  Hash table to free
 */
void freeHashTable(hash_table_t* htable);

/*
 * Inserts the entry specified into the hash table
 *  @param htable  Hash table to insert to, will generate an error with err
 * macro if NULL
 *  @param entry   Entry to insert, will generate an error with err macro if
 * NULL
 */
void htableInsert(hash_table_t* htable, htable_entry_t* entry);

/*
 * Deletes an entry from a hash table from a key
 *  @param htable  Hash table to delete from, will generate an error with err
 * macro if NULL
 *  @param key     Entry's key, does not check if it's NULL
 */
void htableDelete(hash_table_t* htable, char* key);

/*
 * Look up an entry in a hash table from a key
 *  @param htable  Hash table to look in, will generate an error with err macro
 * if NULL
 *  @param key     Entry's key, does not check if it's NULL
 *
 *  @return Found entry, if not found, returns NULL
 */
htable_entry_t* htableLookUp(hash_table_t* htable, char* key);