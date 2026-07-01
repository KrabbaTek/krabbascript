/*
    KrabbaScript project
    File: src/hash_table/hash_table.c
    Purpose: A hash table implementation for the symbol table
    Contributors:
        - khytryy (yehor.khytryy@gmail.com)
        - Visser13
*/

#define _POSIX_C_SOURCE 200809L
#include <hash_table/hash_table.h>

/*
- Visser
This is a hash table to start making the symbol table. This has hasn't been
fully tested yet and I need to look up more about hash tables and symbol tables.
Granted it's a start (This may get moved to a separate file for a hash table).
*/

uint64_t hash(char* key) {
	if (!key) {
		err("Key fed into hash is NULL\n");
		return 0;
	}

	uint64_t h = 0xcbf29ce484222325;
	while (*key) {
		h ^= (unsigned int)*key;
		h *= 0x100000001b3;

		key++;
	}

	return h;
}

hash_table_t* newHashTable() {
	hash_table_t* table = malloc(sizeof(hash_table_t));
	if (!table) {
		err("Failed to malloc a hash table\n");
		return NULL;
	}

	table->count  = 0;
	table->bucket = (htable_entry_t**)calloc(KSCRIPT_HTABLE_BUCKET_SIZE,
	                                         sizeof(htable_entry_t*));

	if (!table->bucket) {
		err("Failed to calloc hash table entries\n");

		free(table);
		return NULL;
	}

	return table;
}

void freeHashTable(hash_table_t* htable) {
	if (!htable) return;

	for (size_t i = 0; i < KSCRIPT_HTABLE_BUCKET_SIZE; i++) {
		htable_entry_t* entry = htable->bucket[i];
		while (entry) {
			htable_entry_t* next = (htable_entry_t*)entry->next_entry;

			if (entry->value_owned) free(entry->value);
			if (entry->key_owned) free(entry->key);

			free(entry);
			entry = next;
		}
	}

	free(htable->bucket);
	free(htable);
}

void htableInsert(hash_table_t* htable, htable_entry_t* entry) {
	if (!htable) {
		err("Hash table fed into htableInsert is NULL\n");
		return;
	}

	if (!entry) {
		err("Hash table entry fed into htableInsert is NULL\n");
		return;
	}

	uint64_t h = hash(entry->key) % KSCRIPT_HTABLE_BUCKET_SIZE;

	entry->next_entry = (struct htable_entry_t*)htable->bucket[h];

	htable->bucket[h] = entry;
	htable->count++;
}

htable_entry_t* htableLookUp(hash_table_t* htable, char* key) {
	if (!htable) {
		err("Hash table fed into htableLookUp is NULL\n");
		return NULL;
	}

	uint64_t        h     = hash(key) % KSCRIPT_HTABLE_BUCKET_SIZE;
	htable_entry_t* entry = htable->bucket[h];

	while (entry) {
		if (strcmp(entry->key, key) == 0) return entry;
		entry = (htable_entry_t*)entry->next_entry;
	}
	return NULL;
}

void htableDelete(hash_table_t* htable, char* key) {
	if (!htable) {
		err("Hash table fed into htableDelete is NULL\n");
		return;
	}

	uint64_t        index = hash(key) % KSCRIPT_HTABLE_BUCKET_SIZE;
	htable_entry_t* entry = htable->bucket[index];
	htable_entry_t* prev  = NULL;

	while (entry) {
		if (strcmp(entry->key, key) == 0) {
			if (prev)
				prev->next_entry = entry->next_entry;
			else
				htable->bucket[index] = (htable_entry_t*)entry->next_entry;

			if (entry->value_owned) free(entry->value);
			if (entry->key_owned) free(entry->key);

			free(entry);
			htable->count--;
			return;
		}

		prev  = entry;
		entry = (htable_entry_t*)entry->next_entry;
	}
}