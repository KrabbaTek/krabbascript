#include <stdio.h>
#include <stdbool.h>
#include <string.h>

/*
- Visser
This is a hash table to start making the symbol table. This has hasn't been fully tested yet and I need to look up more
about has tables and symbol tables. Granted it's a start (This may get moved to a separate file for a hash table).
*/

#define MAX_KEY_LENGTH 256
#define TABLE_SIZE 10000

typedef struct entry
{
    char key[MAX_KEY_LENGTH];
    void *value;
    struct entry *next_entry;
} entry;

typedef struct hash_table
{
    int entry_count;
    entry *entries[TABLE_SIZE];
} hash_table;

unsigned int hash(char *key)
{
    int length = strnlen(key, MAX_KEY_LENGTH);
    unsigned int hash_value = 0;
    for (int i = 0; i < length; i++)
    {
        hash_value += key[i];
        hash_value = (hash_value * key[i]) % TABLE_SIZE;
    }
    return hash_value;
}

hash_table init_hash_table()
{   
    hash_table tmp_table;

    for (int i = 0; i < TABLE_SIZE; i++)
    {
        tmp_table.entries[i] = NULL;
    }
    
    return tmp_table;

}

void print_hash_table(hash_table *hash_t)
{
    printf("Start\n");
    for (int i = 0; i < TABLE_SIZE; i++)
    {
        if (hash_t->entries[i] == NULL)
        {
            printf("\t%i\t---\n", i);
        }
        else
        {
            printf("\t%i\t%s\n", i, hash_t->entries[i]);
            entry *tmp = hash_t->entries[i];
            while (tmp != NULL)
            {
                printf("%s - ", tmp->key);
                tmp = tmp->next_entry;
            }
            printf("\n");
        }
    }
    printf("End\n");
}

bool hash_table_insert(hash_table *hash_t, entry *k)
{
    if (k == NULL)
        return false;

    int index = hash(k->key);

    k->next_entry = hash_t->entries[index];
    hash_t->entries[index] = k;

    return true;
}

entry *hash_table_lookup(hash_table *hash_t, entry *key)
{
    int index = hash(key);

    entry *tmp = hash_t->entries[index];
    while (tmp != NULL && strncmp(tmp->key, key, MAX_KEY_LENGTH) != 0)
    {
        tmp = tmp->next_entry;
    }
    return tmp;
}

entry *hash_table_delete(hash_table *hash_t, entry *key)
{
    int index = hash(key);
    entry *tmp = hash_t->entries[index];
    entry *prev = NULL;

    while (tmp != NULL && strncmp(tmp->key, key, MAX_KEY_LENGTH) != 0)
    {
        tmp = tmp->next_entry;
    }
    if (tmp != NULL)
        return NULL;
    if (prev == NULL)
    {
        // deleting the head of the linked list
        hash_t->entries[index] = tmp->next_entry;
    }
    else
    {
        prev->next_entry = tmp->next_entry;
    }

    return tmp;
}






