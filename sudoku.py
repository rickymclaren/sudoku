#! /usr/bin/env python
import sys
import datetime
import glob
import re
import itertools

class Cell (object):
    def __init__(self, board, row, col):
        self.board = board
        self.row = row
        self.col = col
        self.possibles = '123456789'

    def __str__(self):
        return "[Cell:{0}:{1}:{2}]".format(self.row, self.col, self.possibles)

    def __repr__(self):
        # List uses __repr__ to print list contents
        return self.__str__()

    def box_number(self):
        return (self.row / 3) * 3 + (self.col / 3) 

    def has_possible(self, value):
        return len(self.possibles) > 1 and str(value) in self.possibles

    def is_subset_of(self, values):
        return len(set(self.possibles) - set(values)) == 0
                
    def remove_possibles(self, values, msg = None):
        found = False
        if len(self.possibles) > 1:
            for value in str(values):
                if value in self.possibles:
                    self.possibles = self.possibles.replace(value, '')
                    if msg:
                        print "{0}removed {1} from {2}:{3} => {4}".format(msg, value, self.row, self.col, self.possibles)
                    if len(self.possibles) == 1:
                        self.board.solve(self, self.possibles)
                    found = True
        return found
    
    
class Board (object):
    def __init__(self, data):
        self.cells = []
        for row in range(0,9):
            for col in range(0,9):
                self.cells.append(Cell(self, row, col))
        for row in range(0,9):
            for col in range(0,9):
                value = data[row * 9 + col]
                if re.match('[1-9]', value):
                    self.solve(self.get_cell(row, col), value)
        self.print_out()
        
                
    def print_out(self):
        data = []
        for c in self.cells:
            data.append(c.possibles)
        size = max(map(lambda x: len(x), data))

        data = map(lambda x: x.ljust(size), data)

        width = size * 9 + 9
        print '=' * width
        for row in range(0, 9):
            #print repr(data[row * 9: row * 9 + 9])
            print "|".join(data[row * 9: row * 9 + 9])
            if row in [2,5]:
                print '-' * width
        print '=' * width

    def get_cell(self, row, col):
        return self.cells[row * 9 + col]
        
    def get_row_text(self, row):
        row_data = []
        for cell in self.cells:
            if cell.row == row:
                if len(cell.possibles) == 1:
                    row_data.append(cell.possibles[0])
                else:
                    row_data.append('0')
        return "".join(row_data)
        
    def get_col_text(self, col):
        col_data = []
        for cell in self.cells:
            if cell.col == col:
                if len(cell.possibles) == 1:
                    col_data.append(cell.possibles[0])
                else:
                    col_data.append('0')
        return "".join(col_data)
        
    def get_box_text(self, row, col):
        data = []
        box = self.get_cell(row, col).box_number()
        for cell in self.cells:
            if cell.box_number() == box:
                if len(cell.possibles) == 1:
                    data.append(cell.possibles[0])
                else:
                    data.append('0')
        return "".join(data)
        
    def get_cells_by_row(self, row):
        cells = []
        for cell in self.cells:
            if cell.row == row:
                cells.append(cell)
        return cells

    def get_cells_by_col(self, col):
        cells = []
        for cell in self.cells:
            if cell.col == col:
                cells.append(cell)
        return cells

    def get_cells_by_box_number(self, box):
        cells = []
        for cell in self.cells:
            if cell.box_number() == box:
                cells.append(cell)
        return cells

    def solve(self, cell, value):
        """Solve a cell and remove any matching possibles on the same row, col, or box"""
        cell.possibles = value
        for c in self.cells:
            if c.row == cell.row or c.col == cell.col or c.box_number() == cell.box_number():
                c.remove_possibles(value) 
            
    def find_singles(self):
        """Find unique possibles by row, col, and box"""
        #print "========= Singles ============"
        for value in "123456789":
            for row in range(9):
                cells = filter(lambda x: x.has_possible(value), self.get_cells_by_row(row))
                if len(cells) == 1:
                    print "Single {0} in row {1}".format(value, row)
                    self.solve(cells[0], value)
                    return True
                            
            for col in range(9):
                cells = filter(lambda x: x.has_possible(value), self.get_cells_by_col(col))
                if len(cells) == 1:
                    print "Single {0} in col {1}".format(value, col)
                    self.solve(cells[0], value)
                    return True
                            
            for box in range(9):
                cells = filter(lambda x: x.has_possible(value), self.get_cells_by_box_number(box))
                if len(cells) == 1:
                    print "Single {0} in box {1}".format(value, box)
                    self.solve(cells[0], value)
                    return True
                            
        return False

    def possibles(self, cells):
        """ Return a list of all the unique outstanding possibles for a set of cells """
        possibles = map(lambda x: x.possibles if len(x.possibles) > 1 else "", cells)
        possibles = set("".join(possibles))
        return sorted(possibles)

    def combinations(self, possibles):
        """ all combinations of possibles from pairs to length - 1 """
        result = []
        for length in range(2,len(possibles)):
            result.extend(itertools.combinations(possibles, length))
        return result

    def naked_cells(self, cells):
        possibles = self.possibles(cells)
        combinations = self.combinations(possibles)
        for combo in combinations:
            matches = filter(lambda x: x.is_subset_of(combo), cells)
            if (len(combo) == len(matches)):
                found = False
                for cell in cells:
                    if not cell in matches:
                        if cell.remove_possibles(combo):
                            found = True
                if found:
                    print "Naked {0} in {1}".format(combo, matches)
                    return True

        return False

    def nakeds(self):
        for row in range(9):
            if self.naked_cells(self.get_cells_by_row(row)):
                return True
        for col in range(9):
            if self.naked_cells(self.get_cells_by_col(col)):
                return True
        for box in range(9):
            if self.naked_cells(self.get_cells_by_box_number(box)):
                return True
        return False

    def find_pointing_pairs(self):
        """If a possible occurs only twice in a box and these are on a row/col then the possible can be removed from the rest of the row/col"""
        #print "========= Pointing pairs ============"
        found = False
        for box in range(9):
            for value in "123456789":
                msg = "Pointing pair {0} in box {1} ".format(value, box)
                matches = filter(lambda x: x.has_possible(value), self.get_cells_by_box_number(box))
                cols = list(set(map(lambda x: x.col, matches)))
                rows = list(set(map(lambda x: x.row, matches)))
                if len(matches) in [2,3]:
                    if len(cols) == 1:
                        cells = set(self.get_cells_by_col(cols[0])) - set(matches)
                        for cell in cells:
                            if cell.remove_possibles(value, msg):
                                found = True
                        if found:
                            return True

                    if len(rows) == 1:
                        cells = set(self.get_cells_by_row(rows[0])) - set(matches)
                        for cell in cells:
                            if cell.remove_possibles(value, msg):
                                found = True
                        if found:
                            return True

        return False

    def box_line_reduction(self):
        """If a possible occurs in only one box of a row/col then it can be removed from the rest of the box.
        It is the reverse of pointing pairs
        """
        #print "========= Box Line Reduction ============"
        found = False
        for row in range(0, 9):
            for possible in range(0,9):
                boxes = []
                for cell in self.get_cells_by_row(row):
                    if cell.has_possible(possible):
                        box = cell.box_number()
                        if not box in boxes:
                            boxes.append(box)
                if len(boxes) == 1:
                    for cell in self.get_cells_by_box_number(boxes[0]):
                        if cell.row != row:
                            msg = "Box line reduction in box {0} row {1} for {2}: ".format(boxes[0], row, possible) 
                            if cell.remove_possibles(possible, msg):
                                found = True
                            
        for col in range(0, 9):
            for possible in range(0,9):
                boxes = []
                for cell in self.get_cells_by_col(row):
                    if cell.has_possible(possible):
                        box = cell.box_number()
                        if not box in boxes:
                            boxes.append(box)
                if len(boxes) == 1:
                    for cell in self.get_cells_by_box_number(boxes[0]):
                        if cell.col != col:
                            msg = "Box line reduction in box {0} col {1} for {2}: ".format(boxes[0], col, possible) 
                            if cell.remove_possibles(possible):
                                found = True
                            
        return found

    def x_wing(self):
        """If two rows have a possible in the same two columns then the possible can be removed from other rows on the same columns
        """
        #print "========= X-Wing ============"
        found = False
        for possible in range(1,10):
            rows = []
            for row in range(0, 9):
                cols = [possible, row]
                for cell in self.get_cells_by_row(row):
                    if cell.has_possible(possible):
                        cols.append(cell.col)
                if len(cols) == 4:
                    rows.append(cols)
    
            while len(rows) > 1:            
                for i in range(1, len(rows)):
                    first = str(rows[0][2]) + str(rows[0][3])
                    second = str(rows[i][2]) + str(rows[i][3])
                    if first == second:
                        col_one = rows[0][2]
                        col_two = rows[0][3]
                        row_one = rows[0][1]
                        row_two = rows[i][1]
                        msg = "X wing for {0} in rows {1} and {2}: ".format(possible, row_one, row_two) 
                        for cells in self.get_cells_by_col(col_one):
                            if cell.row != row_one and cell.row != row_two:
                                if cell.remove_possibles(possible, msg):
                                    found = True
                        for cells in self.get_cells_by_col(col_two):
                            if cell.row != row_one and cell.row != row_two:
                                if cell.remove_possibles(possible, msg):
                                    found = True
                        rows.pop(i)
                        break
                rows.pop(0)
                            
        for possible in range(1,10):
            cols = []
            for col in range(0, 9):
                rows = [possible, col]
                for cell in self.get_cells_by_col(col):
                    if cell.has_possible(possible):
                        rows.append(cell.row)
                if len(rows) == 4:
                    cols.append(rows)
    
            while len(cols) > 1:            
                for i in range(1, len(cols)):
                    first = str(cols[0][2]) + str(cols[0][3])
                    second = str(cols[i][2]) + str(cols[i][3])
                    if first == second:
                        row_one = cols[0][2]
                        row_two = cols[0][3]
                        col_one = cols[0][1]
                        col_two = cols[i][1]
                        msg = "X wing for {0} in cols {1} and {2}".format(possible, col_one, col_two) 
                        for cells in self.get_cells_by_row(row_one):
                            if cell.col != col_one and cell.col != col_two:
                                if cell.remove_possibles(possible, msg):
                                    found = True
                        for cells in self.get_cells_by_row(row_two):
                            if cell.col != col_one and cell.col != col_two:
                                if cell.remove_possibles(possible, msg):
                                    found = True
                        cols.pop(i)
                        break
                cols.pop(0)
                            
        return found

    def swordfish(self):
        """If a possible exists in the same 3 columns of 3 rows then the possible can be removed from other rows on the same columns
           Note: It doesn't need to be 3 - it could be 2.
        """
        #print "========= Swordfish ============"
        found = False
        for possible in range(1,10):
            rows = []
            for row in range(0, 9):
                cols = [possible, row]
                for cell in self.get_cells_by_row(row):
                    if cell.has_possible(possible):
                        cols.append(cell.col)
                if len(cols) == 4 or len(cols) == 5:
                    rows.append(cols)
            
            # for rows with 3 cols check if there are 2 other rows that are in those cols
            for row in range(0, len(rows)):
                if len(rows[row]) == 5:   # 3 cols
                    cols = rows[row][2:]
                    triple = [rows[row]]
                    for match in range(0, len(rows)):
                        if match != row:
                            is_match = True
                            for col in rows[match][2:]: 
                                if not col in cols:
                                    is_match = False
                            if is_match:
                                triple.append(rows[match])
                    if len(triple) == 3:
                        msg = "triple found " + repr(triple)
                        triple_rows = [triple[0][1], triple[1][1], triple[2][1]]
                        for col in cols:
                            for cell in self.get_cells_by_col(col):
                                if not cell.row in triple_rows:
                                    if cell.remove_possibles(possible, msg):
                                        found = True
                        
            
        return found

    def solved(self):
    
        for row in range(0, 9):
            answers = []
            for cell in self.get_cells_by_row(row):
                answers.append(cell.possibles)
            answers.sort()
            answer = "".join(answers)
            if answer != "123456789":
                return False
                
        for col in range(0, 9):
            answers = []
            for cell in self.get_cells_by_col(col):
                answers.append(cell.possibles)
            answers.sort()
            answer = "".join(answers)
            if answer != "123456789":
                return False
                
        for box in range(0, 9):
            answers = []
            for cell in self.get_cells_by_box_number(box):
                answers.append(cell.possibles)
            answers.sort()
            answer = "".join(answers)
            if answer != "123456789":
                return False
                
        return True

    def solution(self):
        solution = ""
        if self.solved:
            for row in range(0, 9):
                for col in range(0, 9):
                    solution += self.get_cell(row, col).possibles
            return solution
        else:
            return 'Not solved'
            
solutions = []
with open('top95expected.txt') as f:
    solutions = f.readlines()

problems = []
with open('top95.txt') as f:
    problems = f.readlines()

results = ""
index = 0
solved = 0
for problem in problems:
    index += 1
    if (len(sys.argv) == 2) and (sys.argv[1] != str(index)):
        continue

    board = Board(problem)

    running = True
    while running:
        if board.find_singles():
            pass
        elif board.nakeds():
            pass
        elif board.find_pointing_pairs():
            pass
        elif board.box_line_reduction():
            pass
        elif board.x_wing():
            pass
        elif board.swordfish():
            pass
        else:
            running = False

        board.print_out()
            
    if board.solved():
        print ">>> SOLVED >>>"
        expected = solutions[index-1].rstrip()
        solution = board.solution()
        if solution != expected:
            print "Houston we have a problem"
            print "Expected \n[{0}] \ngot\n[{1}]".format(expected, solution)
            sys.exit()
        solved += 1
        results += 'S'
    else:            
        print ">>> BEATS ME >>>"
        results += '.'
    if index % 10 == 0:
        results += '\n'

print results
print "Solved {0} of {1}".format(solved, index)
 
