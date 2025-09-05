#include "sudoku.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "arraylist.h"
#include "utils.h"

#define FALSE 0
#define TRUE (!FALSE)

char *POSSIBLES = "123456789";

void solve(Cell *cell, char c) {
  cell->values[0] = c;
  cell->values[1] = '\0';
}

int solved(Cell *cell) { return strlen(cell->values) == 1; }

int board_solved(Board *b) {
  for (int i = 0; i < 81; i++) {
    if (!solved(&b->cells[i])) {
      return FALSE;
    }
  }
  return TRUE;
}

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

int has_others(Cell *cell, char *chars) {
  if (solved(cell)) {
    return FALSE;
  }
  char *values = copy_string(cell->values);
  for (int i = 0; chars[i] != '\0'; i++) {
    removeChar(values, chars[i]);
  }
  int others = strlen(values) > 0;
  free(values);
  return others;
}

int can_see(Cell *a, Cell *b) {
  if (a == b) {
    // A cell cannot see itself
    return FALSE;
  }
  return (a->row == b->row) || (a->col == b->col) || (a->box == b->box);
}

int same_row(ArrayList *cells) {
  if (cells->size == 0) {
    return FALSE;
  }
  Cell *first = get(cells, 0);
  for (int i = 1; i < cells->size; i++) {
    Cell *other = get(cells, i);
    if (other->row != first->row) {
      return FALSE;
    }
  }
  return TRUE;
}

int same_col(ArrayList *cells) {
  if (cells->size == 0) {
    return FALSE;
  }
  Cell *first = get(cells, 0);
  for (int i = 1; i < cells->size; i++) {
    Cell *other = get(cells, i);
    if (other->col != first->col) {
      return FALSE;
    }
  }
  return TRUE;
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
  int removed = TRUE;
  while (removed) {
    removed = FALSE;
    for (int i = 0; i < 81; i++) {
      Cell *cell = &b->cells[i];
      if (solved(cell)) {
        char ch = cell->values[0];
        for (int j = 0; j < 81; j++) {
          Cell *other = &b->cells[j];
          if (can_see(cell, other) && !solved(other)) {
            if (removeChar(other->values, ch)) {
              removed = TRUE;
            }
          }
        }
      }
    }
  }
}

int remove_from_cells_outside_box(Cell *cells[9], int box, char ch) {
  int removed = FALSE;
  for (int i = 0; i < 9; i++) {
    Cell *cell = cells[i];
    if (cell->box != box && !solved(cell)) {
      removed |= removeChar(cell->values, ch);
    }
  }
  return removed;
}

/*

*/
int singles(Board *b) {
  for (char c = '1'; c <= '9'; c++) {
    for (int row = 0; row < 9; row++) {
      ArrayList *list = createArrayList(9);
      for (int i = 0; i < 9; i++) {
        Cell *cell = b->rows[row][i];
        if (has(cell, c)) {
          add(list, cell);
        }
      }
      if (list->size == 1) {
        printf("Found single at row %d for %c\n", row + 1, c);
        Cell *only = get(list, 0);
        solve(only, c);
        freeArrayList(list);
        return TRUE;
      }
      freeArrayList(list);
    }

    for (int col = 0; col < 9; col++) {
      ArrayList *list = createArrayList(9);
      for (int i = 0; i < 9; i++) {
        Cell *cell = b->cols[col][i];
        if (has(cell, c)) {
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
        Cell *cell = b->boxes[box][i];
        if (has(cell, c)) {
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

/*
 * Look for pairs of cells in a row, column, or box that only contain the
 * same two values. If found, these values can be removed from all other cells
 */
int naked_pairs(Board *board) {
  ArrayList *naked;
  ArrayList *others;

  for (char a = '1'; a <= '9'; a++) {
    for (char b = a + 1; b <= '9'; b++) {
      char *pair = malloc(3);
      pair[0] = a;
      pair[1] = b;
      pair[2] = '\0';

      // Check rows for naked pairs
      for (int row = 0; row < 9; row++) {
        if (strcmp(pair, "12") == 0 && row == 3) {
          int debug = 1;
        }
        naked = createArrayList(9);
        others = createArrayList(9);
        for (int i = 0; i < 9; i++) {
          if (solved(board->rows[row][i])) {
            continue;
          }
          if (has_others(board->rows[row][i], pair)) {
            add(others, board->rows[row][i]);
          } else {
            add(naked, board->rows[row][i]);
          }
        }
        if (naked->size == 2) {
          int removed = FALSE;
          for (int i = 0; i < others->size; i++) {
            Cell *cell = get(others, i);
            removed |= removeChar(cell->values, a);
            removed |= removeChar(cell->values, b);
          }
          freeArrayList(naked);
          freeArrayList(others);
          if (removed) {
            printf("Found naked pair at row %d for %c%c\n", row + 1, a, b);
            return TRUE;
          }
        }
      }

      // Check columns for naked pairs
      for (int col = 0; col < 9; col++) {
        if (strcmp(pair, "12") == 0 && col == 3) {
          int debug = 1;
        }
        naked = createArrayList(9);
        others = createArrayList(9);
        for (int i = 0; i < 9; i++) {
          if (solved(board->cols[col][i])) {
            continue;
          }
          if (has_others(board->cols[col][i], pair)) {
            add(others, board->cols[col][i]);
          } else {
            add(naked, board->cols[col][i]);
          }
        }
        if (naked->size == 2) {
          int removed = FALSE;
          for (int i = 0; i < others->size; i++) {
            Cell *cell = get(others, i);
            removed |= removeChar(cell->values, a);
            removed |= removeChar(cell->values, b);
          }
          freeArrayList(naked);
          freeArrayList(others);
          if (removed) {
            printf("Found naked pair at col %d for %c%c\n", col + 1, a, b);
            return TRUE;
          }
        }
      }

      // Check boxes for naked pairs
      for (int box = 0; box < 9; box++) {
        naked = createArrayList(9);
        others = createArrayList(9);
        for (int i = 0; i < 9; i++) {
          if (solved(board->boxes[box][i])) {
            continue;
          }
          if (has_others(board->boxes[box][i], pair)) {
            add(others, board->boxes[box][i]);
          } else {
            add(naked, board->boxes[box][i]);
          }
        }
        if (naked->size == 2) {
          int removed = FALSE;
          for (int i = 0; i < others->size; i++) {
            Cell *cell = get(others, i);
            removed |= removeChar(cell->values, a);
            removed |= removeChar(cell->values, b);
          }
          freeArrayList(naked);
          freeArrayList(others);
          if (removed) {
            printf("Found naked pair at box %d for %c%c\n", box + 1, a, b);
            return TRUE;
          }
        }
      }
    }
  }
  return FALSE;
}

/*
 * Look for triples of cells in a row, column, or box that only contain the
 * same three values. If found, these values can be removed from all other cells
 */
int naked_triples(Board *board) {
  ArrayList *naked;
  ArrayList *others;
  int debug = 0;

  for (char a = '1'; a <= '9'; a++) {
    for (char b = a + 1; b <= '9'; b++) {
      for (char c = b + 1; c <= '9'; c++) {
        char *triple = malloc(4);
        triple[0] = a;
        triple[1] = b;
        triple[2] = c;
        triple[3] = '\0';

        // Check rows for naked triples
        for (int row = 0; row < 9; row++) {
          naked = createArrayList(9);
          others = createArrayList(9);
          if ((strcmp(triple, "125") == 0) && (row == 5)) {
            debug = 1;
          }
          for (int i = 0; i < 9; i++) {
            if (solved(board->rows[row][i])) {
              continue;
            }
            if (has_others(board->rows[row][i], triple)) {
              add(others, board->rows[row][i]);
            } else {
              add(naked, board->rows[row][i]);
            }
          }
          if (naked->size == 3) {
            int removed = FALSE;
            for (int i = 0; i < others->size; i++) {
              Cell *cell = get(others, i);
              removed |= removeChar(cell->values, a);
              removed |= removeChar(cell->values, b);
              removed |= removeChar(cell->values, c);
            }
            freeArrayList(naked);
            freeArrayList(others);
            if (removed) {
              printf("Found naked triple at row %d for %c%c%c\n", row + 1, a, b,
                     c);
              return TRUE;
            }
          }
        }

        // Check columns for naked triples
        for (int col = 0; col < 9; col++) {
          naked = createArrayList(9);
          others = createArrayList(9);
          if ((strcmp(triple, "125") == 0) && (col == 4)) {
            debug = 1;
          }
          for (int i = 0; i < 9; i++) {
            if (solved(board->cols[col][i])) {
              continue;
            }
            if (has_others(board->cols[col][i], triple)) {
              add(others, board->cols[col][i]);
            } else {
              add(naked, board->cols[col][i]);
            }
          }
          if (naked->size == 3) {
            int removed = FALSE;
            for (int i = 0; i < others->size; i++) {
              Cell *cell = get(others, i);
              removed |= removeChar(cell->values, a);
              removed |= removeChar(cell->values, b);
              removed |= removeChar(cell->values, c);
            }
            freeArrayList(naked);
            freeArrayList(others);
            if (removed) {
              printf("Found naked triple at col %d for %c%c%c\n", col + 1, a, b,
                     c);
              return TRUE;
            }
          }
        }

        // Check boxes for naked triples
        for (int box = 0; box < 9; box++) {
          naked = createArrayList(9);
          others = createArrayList(9);
          for (int i = 0; i < 9; i++) {
            if (solved(board->boxes[box][i])) {
              continue;
            }
            if (has_others(board->boxes[box][i], triple)) {
              add(others, board->boxes[box][i]);
            } else {
              add(naked, board->boxes[box][i]);
            }
          }
          if (naked->size == 3) {
            int removed = FALSE;
            for (int i = 0; i < others->size; i++) {
              Cell *cell = get(others, i);
              removed |= removeChar(cell->values, a);
              removed |= removeChar(cell->values, b);
              removed |= removeChar(cell->values, c);
            }
            freeArrayList(naked);
            freeArrayList(others);
            if (removed) {
              printf("Found naked triple at box %d for %c%c%c\n", box + 1, a, b,
                     c);
              return TRUE;
            }
          }
        }
      }
    }
  }
  return FALSE;
}

/*
 * Look for triples of cells in a row, column, or box that only contain the
 * same three values. If found, these values can be removed from all other cells
 */
int pointing_pairs(Board *board) {
  int debug = 0;

  for (char ch = '1'; ch <= '9'; ch++) {
    for (int box = 0; box < 9; box++) {
      if ((ch == '3') && (box == 5)) {
        debug = 1;
      }
      ArrayList *pair = createArrayList(9);
      for (int i = 0; i < 9; i++) {
        if (has(board->boxes[box][i], ch)) {
          add(pair, board->boxes[box][i]);
        }
      }
      if (pair->size == 2 || pair->size == 3) {
        if (same_row(pair)) {
          Cell *first = get(pair, 0);
          // Both in same row, remove from other cells in that row outside box
          if (remove_from_cells_outside_box(board->rows[first->row], box, ch)) {
            printf("Found pointing pair at box %d for %c in row %d\n", box + 1,
                   ch, first->row + 1);
            freeArrayList(pair);
            return TRUE;
          }
        }
        if (same_col(pair)) {
          Cell *first = get(pair, 0);
          // Both in same row, remove from other cells in that row outside box
          if (remove_from_cells_outside_box(board->cols[first->col], box, ch)) {
            printf("Found pointing pair at box %d for %c in col %d\n", box + 1,
                   ch, first->col + 1);
            freeArrayList(pair);
            return TRUE;
          }
        }
      }
      freeArrayList(pair);
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
      sleep(1);
      remove_solved(&b);
      print_board(&b);
      if (board_solved(&b)) {
        printf("Board solved!\n");
        break;
      }
      if (singles(&b)) {
        continue;
      } else if (naked_pairs(&b)) {
        continue;
      } else if (naked_triples(&b)) {
        continue;
      } else if (pointing_pairs(&b)) {
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