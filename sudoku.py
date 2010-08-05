#! /usr/bin/env python
import sys
import datetime
import glob
from xml.etree import ElementTree
from PyQt4.QtCore import *
from PyQt4.QtGui import *

box_size = 70

def get_data():
    return '700200040400030907015040000000700000057000210000009000000080360503070002020006004'

class Cell (object):
    def __init__(self, board, row, col):
        self.board = board
        self.row = row
        self.col = col
        self.possibles = '123456789'
            
    def box_number(self):
        return (self.row / 3) * 3 + (self.col / 3) 

    def has_possible(self, value):
        return len(self.possibles) > 1 and str(value) in self.possibles
                
    def remove_possibles(self, values, msg = ""):
        if len(self.possibles) > 1:
            for value in str(values):
                if value in self.possibles:
                    self.possibles = self.possibles.replace(value, '')
                    print "{0}removed {1} from {2}:{3} => {4}".format(msg, value, self.row, self.col, self.possibles)
                    if len(self.possibles) == 1:
                        self.board.solve(self, self.possibles)
                    return True
        return False
    
    
class Board (QWidget):
    def __init__(self, data, parent=None):
        QWidget.__init__(self, parent)
        self.setGeometry(QRect(400, 100, box_size * 9, box_size * 9))
        self.cells = []
        for row in range(0,9):
            for col in range(0,9):
                self.cells.append(Cell(self, row, col))
        for row in range(0,9):
            for col in range(0,9):
                value = data[row * 9 + col]
                if value != '0':
                    self.solve(self.get_cell(row, col), value)
        self.print_out()
        
                
    def print_out(self):
        data = []
        for c in self.cells:
            data.append(c.possibles)
            
        print '=' * 60
        for row in range(0, 9):
            print repr(data[row * 9: row * 9 + 9])
        print '=' * 60

    def paintEvent(self, event):
        painter = QPainter()
        painter.begin(self)
        
        solvedFont = QFont('Helvetica', 14, QFont.Bold)
        possiblesFont = QFont('Helvetica', 10, QFont.Normal)

        # draw the grid
        for y in range(0, 10):
            if y % 3 == 0:
                painter.setPen(QPen(Qt.black, 3))
            else:
                painter.setPen(QPen(Qt.black, 1))
            painter.drawLine(0, y * box_size, box_size * 9, y * box_size)        
        for x in range(0, 10):
            if x % 3 == 0:
                painter.setPen(QPen(Qt.black, 3))
            else:
                painter.setPen(QPen(Qt.black, 1))
            painter.drawLine(x * box_size, 0, x * box_size, box_size * 9)        
            
        # draw the solved squares
        painter.setFont(solvedFont)
        for row in range(0, 9):
            x = 0
            for cell in self.get_row_text(row):
                if cell != '0': painter.drawText((x + 0.5) * box_size, (row + 0.5) * box_size, cell)
                x = x + 1
            
        # draw the possibles
        painter.setFont(possiblesFont)
        painter.setPen(QPen(Qt.gray))
        for row in range(0, 9):
            x = 0
            for cell in self.get_row_text(row):
                if cell == '0': 
                    poss = self.get_cell(row, x).possibles
                    painter.drawText(x * box_size + 4, row * box_size + 15, poss)
                x = x + 1
            
        painter.end()
        
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
        
    def get_other_possibles_by_row(self, row, col):
        data = []
        for cell in self.cells:
            if cell.row == row:
                if cell.col != col:
                    if len(cell.possibles) > 1:
                        data.append(cell.possibles)
        return "".join(data)
        
    def get_other_possibles_by_col(self, row, col):
        data = []
        for cell in self.cells:
            if cell.col == col:
                if cell.row != row:
                    if len(cell.possibles) > 1:
                        data.append(cell.possibles)
        return "".join(data)
        
    def get_other_possibles_by_box(self, row, col):
        data = []
        box = self.get_cell(row, col).box_number()
        for cell in self.cells:
            if cell.box_number() == box:
                if cell.row != row or cell.col != col:
                    if len(cell.possibles) > 1:
                        data.append(cell.possibles)
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
        print "========= Singles ============"
        for row in range(0, 9):
            for col in range(0, 9):
                cell = self.get_cell(row, col)
                # Unique by row
                if len(cell.possibles) > 1:
                    for value in cell.possibles:
                        if not value in self.get_other_possibles_by_row(row, col):
                            print "{0}:{1} unique by row => {2}".format(row, col, value)
                            self.solve(cell, value)
                            return True

                # Unique by col
                if len(cell.possibles) > 1:
                    for value in cell.possibles:
                        if not value in self.get_other_possibles_by_col(row, col):
                            print "{0}:{1} unique by col => {2}".format(row, col, value)
                            self.solve(cell, value)
                            return True

                # Unique by box
                if len(cell.possibles) > 1:
                    for value in cell.possibles:
                        if not value in self.get_other_possibles_by_box(row, col):
                            print "{0}:{1} unique by box => {2}".format(row, col, value)
                            self.solve(cell, value)
                            return True
                            
        return False

    def find_naked_pairs(self):
        """If two cells have only the same two possibles then those possibles can be removed from the rest of the row, col, and box"""
        print "========= Naked Pairs ============"
        for row in range(0, 9):
            cells = self.get_cells_by_row(row)
            twos = []
            for cell in cells:
                if len(cell.possibles) == 2:
                    twos.append(cell)
            if len(twos) == 2 and twos[0].possibles == twos[1].possibles:
                msg = "row {0} has naked pair {1}: ".format(row, twos[0].possibles)
                for cell in cells:
                    if cell.possibles != twos[0].possibles:
                        if cell.remove_possibles(twos[0].possibles, msg):
                            return True
                            
        for col in range(0, 9):
            cells = self.get_cells_by_col(col)
            twos = []
            for cell in cells:
                if len(cell.possibles) == 2:
                    twos.append(cell)
            if len(twos) == 2 and twos[0].possibles == twos[1].possibles:
                msg = "col {0} has naked pair {1}: ".format(col, twos[0].possibles)
                for cell in cells:
                    if cell.possibles != twos[0].possibles:
                        if cell.remove_possibles(twos[0].possibles, msg):
                            return True
                            
        for box in range(0, 9):
            cells = self.get_cells_by_box_number(box)
            twos = []
            for cell in cells:
                if len(cell.possibles) == 2:
                    twos.append(cell)
            if len(twos) == 2 and twos[0].possibles == twos[1].possibles:
                msg = "box {0} has naked pair {1}: ".format(box, twos[0].possibles)
                for cell in cells:
                    if cell.possibles != twos[0].possibles:
                        if cell.remove_possibles(twos[0].possibles, msg):
                            return True
                            
        return False

    def find_pointing_pairs(self):
        """If a possible occurs only twice in a box and these are on a row/col then the possible can be removed from the rest of the row/col"""
        print "========= Pointing pairs ============"
        found = False
        for box in range(0, 9):
            cells = self.get_cells_by_box_number(box)
            for possible in range(0,9):
                pair = []
                for cell in cells:
                    if cell.has_possible(possible):
                        pair.append(cell)
                if len(pair) == 2:
                    if pair[0].row == pair[1].row:
                        row = pair[0].row
                        for cell in self.get_cells_by_row(row):
                            if cell != pair[0] and cell != pair[1]:
                                msg = "box {0} has pointing pair {1} in row {2}: ".format(box, possible, row)
                                if cell.remove_possibles(possible):
                                    found = True

                    if pair[0].col == pair[1].col:
                        col = pair[0].col
                        for cell in self.get_cells_by_col(col):
                            if cell != pair[0] and cell != pair[1]:
                                msg = "box {0} has pointing pair {1} in col {2}: ".format(box, possible, col)
                                if cell.remove_possibles(possible, msg):
                                    found = True
                            
                if len(pair) == 3:
                    if pair[0].row == pair[1].row and pair[1].row == pair[2].row:
                        row = pair[0].row
                        for cell in self.get_cells_by_row(row):
                            if cell != pair[0] and cell != pair[1] and cell != pair[2]:
                                msg = "box {0} has pointing triple {1} in row {2}: ".format(box, possible, row)
                                if cell.remove_possibles(possible, msg):
                                    found = True

                    if pair[0].col == pair[1].col and pair[1].col == pair[2].col:
                        col = pair[0].col
                        for cell in self.get_cells_by_col(col):
                            if cell != pair[0] and cell != pair[1] and cell != pair[2]:
                                msg = "box {0} has pointing triple {1} in col {2}: ".format(box, possible, col)
                                if cell.remove_possibles(possible, msg):
                                    found = True
                            
        return found

    def box_line_reduction(self):
        """If a possible occurs in only one box of a row/col then it can be removed from the rest of the box.
        It is the reverse of pointing pairs
        """
        print "========= Box Line Reduction ============"
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
        print "========= X-Wing ============"
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
        print "========= Swordfish ============"
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
            

app = QApplication(sys.argv)
board = Board(get_data())

running = True
while running:
    if board.find_singles():
        pass
    elif board.find_naked_pairs():
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
    print "=== SOLVED ==="
else:            
    board.show()
    app.exec_()

