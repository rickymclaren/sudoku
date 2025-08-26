#ifndef ARRAYLIST_H
#define ARRAYLIST_H
typedef struct {
  void** data;   // Pointer to the array
  int size;      // Current number of elements
  int capacity;  // Total allocated space
} ArrayList;

ArrayList* createArrayList(int initialCapacity);
int add(ArrayList* list, void* element);
void* get(ArrayList* list, int index);
int removeAt(ArrayList* list, int index);
void freeArrayList(ArrayList* list);

#endif  // ARRAYLIST_H