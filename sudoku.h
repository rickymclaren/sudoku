typedef struct cell {
  int row;
  int col;
  int box;
  char* values;
} Cell;

typedef struct cell_list {
  Cell* cell[9];
  int count;
} CellList;

typedef struct board {
  Cell cells[81];
  Cell* rows[9][9];
  Cell* cols[9][9];
  Cell* boxes[9][9];
} Board;
