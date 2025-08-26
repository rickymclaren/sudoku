#include "sudoku.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "arraylist.h"

#define FALSE 0
#define TRUE (!FALSE)

char *POSSIBLES = "123456789";

void solve(Cell *cell, char c) {
  cell->values[0] = c;
  cell->values[1] = '\0';
}

int solved(Cell *cell) { return strlen(cell->values) == 1; }

int has(Cell *cell, char c) {
  for (int i = 0; cell->values[i] != '\0'; i++) {
    if (cell->values[i] == c) {
      return TRUE;
    }
  }
  return FALSE;
}

int can_see(Cell *a, Cell *b) {
  if (a == b) {
    // A cell cannot see itself
    return FALSE;
  }
  return (a->row == b->row) || (a->col == b->col) || (a->box == b->box);
}

char *copy_string(const char *s) {
  int len = 0;
  while (s[len] != '\0') {
    len++;
  }
  char *new_str = (char *)malloc((len + 1) * sizeof(char));
  for (int i = 0; i <= len; i++) {
    new_str[i] = s[i];
  }
  return new_str;
}

void initialize_board(struct board *b, char *line) {
  for (int r = 0; r < 9; r++) {
    for (int c = 0; c < 9; c++) {
      b->cells[r][c].row = r;
      b->cells[r][c].col = c;
      b->cells[r][c].box = (r / 3) * 3 + (c / 3);
      b->cells[r][c].values = copy_string(POSSIBLES);
      char ch = line[r * 9 + c];
      if (ch >= '1' && ch <= '9') {
        solve(&b->cells[r][c], ch);
      }
    }
  }
}

int singles(Board *b) {
  for (int row = 0; row < 9; row++) {
    for (char c = '1'; c <= '9'; c++) {
      ArrayList *list = createArrayList(9);
      for (int col = 0; col < 9; col++) {
        Cell *cell = &b->cells[row][col];
        if (has(cell, c) && !solved(cell)) {
          add(list, cell);
        }
      }
      if (list->size == 1) {
        printf("Found single at row %d for %c\n", row, c);
        Cell *only = get(list, 0);
        solve(only, c);
        freeArrayList(list);
        return TRUE;
      }
      freeArrayList(list);
    }
  }

  for (int col = 0; col < 9; col++) {
    for (char c = '1'; c <= '9'; c++) {
      ArrayList *list = createArrayList(9);
      for (int row = 0; row < 9; row++) {
        Cell *cell = &b->cells[row][col];
        if (has(cell, c) && !solved(cell)) {
          add(list, cell);
        }
      }
      if (list->size == 1) {
        printf(">>> Found single at col %d for %c\n", col, c);
        Cell *only = get(list, 0);
        solve(only, c);
        freeArrayList(list);
        return TRUE;
      }
      freeArrayList(list);
    }
  }

  return FALSE;
}

int main(void) {
  Board b;
  printf("Sudoku program started.\n");

  char buffer[100];
  FILE *f = fopen("top95.txt", "r");
  if (f == NULL) {
    perror("Failed to open top95.txt");
    return 1;
  }

  char *line = fgets(buffer, sizeof(buffer), f);
  while (line != NULL) {
    printf("Read line: %s\n", line);
    initialize_board(&b, line);

    for (;;) {
      if (singles(&b)) {
        continue;
      } else {
        printf("Beats me.\n");
        break;
      }
    }

    line = fgets(buffer, sizeof(buffer), f);
  }

  fclose(f);

  return 0;
}