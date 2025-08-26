#include "arraylist.h"

#include <stdio.h>
#include <stdlib.h>

#include "sudoku.h"

// Initialize the ArrayList
ArrayList* createArrayList(int initialCapacity) {
  ArrayList* list = (ArrayList*)malloc(sizeof(ArrayList));
  if (!list) return NULL;

  list->data = (void**)malloc(initialCapacity * sizeof(void*));
  if (!list->data) {
    free(list);
    return NULL;
  }

  list->size = 0;
  list->capacity = initialCapacity;
  return list;
}

// Add element to the end
int add(ArrayList* list, void* element) {
  if (!list) return 0;

  // Resize if needed
  if (list->size >= list->capacity) {
    int newCapacity = list->capacity * 2;
    void** newData = (void**)realloc(list->data, newCapacity * sizeof(Cell));
    if (!newData) return 0;

    list->data = newData;
    list->capacity = newCapacity;
  }

  list->data[list->size++] = element;
  return 1;
}

// Get element at index
void* get(ArrayList* list, int index) {
  if (!list || index < 0 || index >= list->size) {
    printf("Index out of bounds\n");
    exit(1);
  };
  return list->data[index];
}

// Remove element at index
int removeAt(ArrayList* list, int index) {
  if (!list || index < 0 || index >= list->size) return 0;

  for (int i = index; i < list->size - 1; i++) {
    list->data[i] = list->data[i + 1];
  }
  list->size--;
  return 1;
}

// Free the ArrayList
void freeArrayList(ArrayList* list) {
  if (list) {
    free(list->data);
    free(list);
  }
}
