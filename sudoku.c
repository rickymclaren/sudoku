#include "sudoku.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

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
  if (solved(cell)) {
    return FALSE;
  }
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

void removeChar(char *str, char charToRemove) {
  int i, j;
  int len = strlen(str);

  // Iterate through the string
  for (i = 0, j = 0; i < len; i++) {
    // Copy only if current char is not the one to remove
    if (str[i] != charToRemove) {
      str[j] = str[i];
      j++;
    }
  }
  // Add null terminator at the end
  str[j] = '\0';
}

void initialize_board(struct board *b, char *line) {
  for (int r = 0; r < 9; r++) {
    for (int c = 0; c < 9; c++) {
      b->cells[r * 9 + c].row = r;
      b->cells[r * 9 + c].col = c;
      b->cells[r * 9 + c].box = (r / 3) * 3 + (c / 3);
      b->cells[r * 9 + c].values = copy_string(POSSIBLES);
      char ch = line[r * 9 + c];
      if (ch >= '1' && ch <= '9') {
        solve(&b->cells[r * 9 + c], ch);
      }
      b->rows[r][c] = &b->cells[r * 9 + c];
      b->cols[c][r] = &b->cells[r * 9 + c];
      b->boxes[r][c] = &b->cells[r * 9 + c];
    }
  }
  for (int box = 0; box < 9; box++) {
    for (int cell = 0; cell < 9; cell++) {
      int r = (box / 3) * 3 + (cell / 3);
      int c = (box % 3) * 3 + (cell % 3);
      b->boxes[box][cell] = &b->cells[r * 9 + c];
    }
  }
}

void remove_solved(Board *b) {
  for (int r = 0; r < 9; r++) {
    for (int c = 0; c < 9; c++) {
      Cell *cell = &b->cells[r * 9 + c];
      if (solved(cell)) {
        char solved_value = cell->values[0];
        for (int i = 0; i < 9; i++) {
          Cell *row_cell = b->rows[r][i];
          if (can_see(cell, row_cell)) {
            removeChar(row_cell->values, solved_value);
          }
          Cell *col_cell = b->cols[c][i];
          if (can_see(cell, col_cell)) {
            removeChar(col_cell->values, solved_value);
          }
          Cell *box_cell = b->boxes[cell->box][i];
          if (can_see(cell, box_cell)) {
            removeChar(box_cell->values, solved_value);
          }
        }
      }
    }
  }
}

int singles(Board *b) {
  for (char c = '1'; c <= '9'; c++) {
    for (int row = 0; row < 9; row++) {
      ArrayList *list = createArrayList(9);
      for (int i = 0; i < 9; i++) {
        Cell **cell = &b->rows[row][i];
        if (has(*cell, c)) {
          add(list, cell);
        }
      }
      if (list->size == 1) {
        printf("Found single at row %d for %c\n", row + 1, c);
        Cell **only = get(list, 0);
        solve(*only, c);
        freeArrayList(list);
        return TRUE;
      }
      freeArrayList(list);
    }

    for (int col = 0; col < 9; col++) {
      ArrayList *list = createArrayList(9);
      for (int i = 0; i < 9; i++) {
        Cell **cell = &b->cols[col][i];
        if (has(*cell, c)) {
          add(list, cell);
        }
      }
      if (list->size == 1) {
        printf("Found single at col %d for %c\n", col + 1, c);
        Cell *only = get(list, 0);
        solve(only, c);
        freeArrayList(list);
        return TRUE;
      }
      freeArrayList(list);
    }

    for (int box = 0; box < 9; box++) {
      ArrayList *list = createArrayList(9);
      for (int i = 0; i < 9; i++) {
        Cell **cell = &b->boxes[box][i];
        if (has(*cell, c)) {
          add(list, cell);
        }
      }
      if (list->size == 1) {
        printf("Found single at box %d for %c\n", box + 1, c);
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

int pairs(Board *board) {
  for (char a = '1'; a <= '9'; a++) {
    for (char b = '1'; b <= '9'; b++) {
      for (int row = 0; row < 9; row++) {
        ArrayList *list = createArrayList(9);
        for (int i = 0; i < 9; i++) {
          Cell *cell = board->rows[row][i];
          if (has(cell, a) || has(cell, b)) {
            add(list, cell);
          }
        }
        if (list->size == 2) {
          printf("Found pair at row %d for %c%c\n", row + 1, a, b);
          for (int i = 0; i < list->size; i++) {
            Cell *cell = get(list, i);
            cell->values[0] = a;
            cell->values[1] = b;
            cell->values[2] = '\0';
          }
          freeArrayList(list);
          return TRUE;
        }
        freeArrayList(list);
      }

      for (int col = 0; col < 9; col++) {
        ArrayList *list = createArrayList(9);
        for (int i = 0; i < 9; i++) {
          Cell *cell = board->cols[col][i];
          if (has(cell, a) || has(cell, b)) {
            add(list, cell);
          }
        }
        if (list->size == 2) {
          printf("Found pair at col %d for %c%c\n", col + 1, a, b);
          for (int i = 0; i < list->size; i++) {
            Cell *cell = get(list, i);
            cell->values[0] = a;
            cell->values[1] = b;
            cell->values[2] = '\0';
          }
          freeArrayList(list);
          return TRUE;
        }
        freeArrayList(list);
      }

      for (int box = 0; box < 9; box++) {
        ArrayList *list = createArrayList(9);
        for (int i = 0; i < 9; i++) {
          Cell *cell = board->boxes[box][i];
          if (has(cell, a) || has(cell, b)) {
            add(list, cell);
          }
        }
        if (list->size == 2) {
          printf("Found pair at box %d for %c%c\n", box + 1, a, b);
          for (int i = 0; i < list->size; i++) {
            Cell *cell = get(list, i);
            cell->values[0] = a;
            cell->values[1] = b;
            cell->values[2] = '\0';
          }
          freeArrayList(list);
          return TRUE;
        }
        freeArrayList(list);
      }
    }
  }
  return FALSE;
}

int triples(Board *board) {
  for (char a = '1'; a <= '9'; a++) {
    for (char b = '1'; b <= '9'; b++) {
      for (char c = '1'; b <= '9'; c++) {
        for (int row = 0; row < 9; row++) {
          ArrayList *list = createArrayList(9);
          for (int i = 0; i < 9; i++) {
            Cell *cell = board->rows[row][i];
            if (has(cell, a) || has(cell, b) || has(cell, c)) {
              add(list, cell);
            }
          }
          if (list->size == 3) {
            printf("Found triple at row %d for %c%c%c\n", row + 1, a, b, c);
            for (int i = 0; i < list->size; i++) {
              Cell *cell = get(list, i);
              cell->values[0] = a;
              cell->values[1] = b;
              cell->values[2] = c;
              cell->values[3] = '\0';
            }
            freeArrayList(list);
            return TRUE;
          }
          freeArrayList(list);
        }

        for (int col = 0; col < 9; col++) {
          ArrayList *list = createArrayList(9);
          for (int i = 0; i < 9; i++) {
            Cell *cell = board->cols[col][i];
            if (has(cell, a) || has(cell, b) || has(cell, c)) {
              add(list, cell);
            }
          }
          if (list->size == 3) {
            printf("Found triple at col %d for %c%c%c\n", col + 1, a, b, c);
            for (int i = 0; i < list->size; i++) {
              Cell *cell = get(list, i);
              cell->values[0] = a;
              cell->values[1] = b;
              cell->values[2] = c;
              cell->values[3] = '\0';
            }
            freeArrayList(list);
            return TRUE;
          }
          freeArrayList(list);
        }

        for (int box = 0; box < 9; box++) {
          ArrayList *list = createArrayList(9);
          for (int i = 0; i < 9; i++) {
            Cell *cell = board->boxes[box][i];
            if (has(cell, a) || has(cell, b) || has(cell, c)) {
              add(list, cell);
            }
          }
          if (list->size == 3) {
            printf("Found triple at box %d for %c%c%c\n", box + 1, a, b, c);
            for (int i = 0; i < list->size; i++) {
              Cell *cell = get(list, i);
              cell->values[0] = a;
              cell->values[1] = b;
              cell->values[2] = c;
              cell->values[3] = '\0';
            }
            freeArrayList(list);
            return TRUE;
          }
          freeArrayList(list);
        }
      }
    }
  }

  return FALSE;
}

void print_board(Board *b) {
  printf("\n");
  for (int i = 0; i < 81; i++) {
    if (solved(&b->cells[i])) {
      printf("    %s     ", b->cells[i].values);
    } else {
      printf("%10s", b->cells[i].values);
    }

    if (i > 0) {
      int j = i + 1;
      if (j % 9 == 0) {
        printf("\n");
        if (j == 27 || j == 54) {
          printf(
              "==============================|==============================|"
              "=="
              "============================\n");
        }
      } else if (j % 3 == 0) {
        printf("|");
      }
    }
  }
  printf("\n");
}

int main(void) {
  Board b;
  printf("Sudoku program started.\n");

  char buffer[100];
  FILE *f = fopen("testeasy.txt", "r");
  if (f == NULL) {
    perror("Failed to open top95.txt");
    return 1;
  }

  char *line = fgets(buffer, sizeof(buffer), f);
  while (line != NULL) {
    printf("Read line: %s\n", line);
    if (strlen(line) != 82) {
      printf("Line not 81.\n");
      line = fgets(buffer, sizeof(buffer), f);
      continue;
    }
    initialize_board(&b, line);

    for (;;) {
      sleep(2);
      remove_solved(&b);
      print_board(&b);
      if (singles(&b)) {
        continue;
      } else if (pairs(&b)) {
        continue;
      } else if (triples(&b)) {
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