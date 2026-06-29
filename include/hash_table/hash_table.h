/*
    KrabbaScript project
    File: include/hash_table/hash_table.h
    Purpose: A hash table implementation for the symbol table
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
        - Visser13
*/
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

hash_table_t* newHashTable();
void          freeHashTable(hash_table_t* htable);

void htableInsert(hash_table_t* htable, htable_entry_t* entry);
void htableDelete(hash_table_t* htable, char* key);

htable_entry_t* htableLookUp(hash_table_t* htable, char* key);