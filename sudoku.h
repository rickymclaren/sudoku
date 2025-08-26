typedef struct cell {
  int row;
  int col;
  int box;
  char *values;
} Cell;

typedef struct cell_list {
  Cell *cell[9];
  int count;
} CellList;

typedef struct board {
  Cell cells[9][9];
} Board;
