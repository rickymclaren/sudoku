package main

/*
 * Terminology
 * ===========
 *
 * The sudoku is made up 81 cells organised in 9 rows and 9 cols
 * These are referred to as the board or b for short.
 * A 3x3 cells group is called a box. There are 9 boxes.
 *
 * For strategies that can use row/col/box interchangably these can
 * be accessed as 27 groups i.e. 9 rows, 9 cols, and 9 boxes.
 */

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strings"
)

//-------------------------------------------

var numbers []string = strings.Split("123456789", "")
var indexes = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
var combinations = combosOfString(numbers, 2)

//-------------------------------------------

// Cell is the basic building block of a sudoku.
// If it has only one possible then it is solved.
type Cell struct {
	row       int
	col       int
	box       int
	possibles []string
}

// Cells is a group of cells. This could be all the cells in a row, column, or box.
type Cells []*Cell

// Chain is a sequence of cell pairs that are connected.
// At most there can be 3 links. One per row, col, and box.
// Links in a chain alternate between two colours.
type Chain struct {
	cell   *Cell
	colour int
	links  []Chain
}

type Board struct {
	cells   [81]Cell
	rows    [9]Cells
	cols    [9]Cells
	boxes   [9]Cells
	blocks  []Cells
	verbose bool
}

// --- Methods of Cell ---
func (cell *Cell) String() string {
	return fmt.Sprintf("C:%v,%v", cell.row+1, cell.col+1)
}

func (cell *Cell) solved() bool {
	return len(cell.possibles) == 1
}

func (cell *Cell) inRow(row int) bool {
	return cell.row == row
}

func (cell *Cell) inCol(col int) bool {
	return cell.col == col
}

func (cell *Cell) inBox(box int) bool {
	return cell.box == box
}

func (cell *Cell) solve(value string) {
	cell.possibles = []string{value}
}

func (cell *Cell) hasPossible(s string) bool {
	if cell.solved() {
		return false
	}
	for _, possible := range cell.possibles {
		if possible == s {
			return true
		}
	}
	return false
}

func (cell *Cell) removePossible(value string) bool {
	if cell.solved() {
		return false
	}
	newPossibles := make([]string, 0, 9)
	result := false
	for _, possible := range cell.possibles {
		if possible == value {
			result = true
		} else {
			newPossibles = append(newPossibles, possible)
		}
	}
	if len(newPossibles) != len(cell.possibles) {
		// fmt.Printf("Setting cell %v %v from %v to %v\n", cell.row, cell.col, cell.possibles, newPossibles)
		cell.possibles = newPossibles
	}
	return result
}

func (cell *Cell) removePossibles(values []string) bool {
	result := false
	for _, value := range values {
		if cell.removePossible(value) {
			result = true
		}
	}
	return result
}

func (cell *Cell) removePossiblesApartFrom(values []string) bool {
	result := false
	for _, possible := range cell.possibles {
		isValue := false
		for _, value := range values {
			if possible == value {
				isValue = true
			}
		}
		if !isValue && cell.removePossible(possible) {
			result = true
		}
	}
	return result
}

func (cell *Cell) hasAnyOf(possibles []string) bool {
	if cell.solved() {
		return false
	}
	for _, possible := range possibles {
		if cell.hasPossible(possible) {
			return true
		}
	}
	return false
}

func (cell *Cell) canSee(chain Chain) int {
	result := 0
	if cell != chain.cell {
		if cell.col == chain.cell.col || cell.row == chain.cell.row || cell.box == chain.cell.box {
			result = result | chain.colour
		}
	}
	for _, link := range chain.links {
		result = result | cell.canSee(link)
	}
	return result
}

// --- Methods of Cells ---

func (cells *Cells) remove(values []string) bool {
	found := false
	for _, cell := range *cells {
		if cell.removePossibles(values) {
			found = true
		}
	}
	return found
}

func (cells *Cells) filterInclude(include func(*Cell) bool) Cells {
	result := []*Cell{}
	for _, c := range *cells {
		if include(c) {
			result = append(result, c)
		}
	}
	return result
}

func (cells *Cells) filterExclude(include func(*Cell) bool) Cells {
	result := make([]*Cell, 0, 9)
	for _, c := range *cells {
		if !include(c) {
			result = append(result, c)
		}
	}
	return result
}

func (cells *Cells) filterHasPossible(s string) Cells {
	result := make([]*Cell, 0, 9)
	for _, c := range *cells {
		if c.hasPossible(s) {
			result = append(result, c)
		}
	}
	return result
}

func (cells *Cells) possibles() []string {
	result := make([]string, 0, 9)
	for _, cell := range *cells {
		if !cell.solved() {
			for _, possible := range cell.possibles {
				result = append(result, possible)
			}
		}
	}
	return unique(result)
}

//------- Methods of Chain -----------

func (chain *Chain) String() string {
	return fmt.Sprintf("%v,c%v,->%v", chain.cell, chain.colour, chain.links)
}

func (chain *Chain) hasCell(cell *Cell) bool {
	return chain.findCell(cell) != nil
}

func (chain *Chain) findCell(cell *Cell) *Chain {
	if chain.cell == cell {
		return chain
	}
	for linkIndex := range chain.links {
		c := chain.links[linkIndex].findCell(cell)
		if c != nil {
			return c
		}
	}
	return nil
}

// --- Methods of Board ---

func (b *Board) init() {
	for i := range b.cells {
		row := i / 9
		col := i % 9
		box := (row / 3 * 3) + col/3
		b.cells[i] = Cell{possibles: numbers, row: row, col: col, box: box}
	}

	b.rows = [9]Cells{
		{&b.cells[0], &b.cells[1], &b.cells[2], &b.cells[3], &b.cells[4], &b.cells[5], &b.cells[6], &b.cells[7], &b.cells[8]},
		{&b.cells[9], &b.cells[10], &b.cells[11], &b.cells[12], &b.cells[13], &b.cells[14], &b.cells[15], &b.cells[16], &b.cells[17]},
		{&b.cells[18], &b.cells[19], &b.cells[20], &b.cells[21], &b.cells[22], &b.cells[23], &b.cells[24], &b.cells[25], &b.cells[26]},
		{&b.cells[27], &b.cells[28], &b.cells[29], &b.cells[30], &b.cells[31], &b.cells[32], &b.cells[33], &b.cells[34], &b.cells[35]},
		{&b.cells[36], &b.cells[37], &b.cells[38], &b.cells[39], &b.cells[40], &b.cells[41], &b.cells[42], &b.cells[43], &b.cells[44]},
		{&b.cells[45], &b.cells[46], &b.cells[47], &b.cells[48], &b.cells[49], &b.cells[50], &b.cells[51], &b.cells[52], &b.cells[53]},
		{&b.cells[54], &b.cells[55], &b.cells[56], &b.cells[57], &b.cells[58], &b.cells[59], &b.cells[60], &b.cells[61], &b.cells[62]},
		{&b.cells[63], &b.cells[64], &b.cells[65], &b.cells[66], &b.cells[67], &b.cells[68], &b.cells[69], &b.cells[70], &b.cells[71]},
		{&b.cells[72], &b.cells[73], &b.cells[74], &b.cells[75], &b.cells[76], &b.cells[77], &b.cells[78], &b.cells[79], &b.cells[80]},
	}

	b.cols = [9]Cells{
		{&b.cells[0], &b.cells[9], &b.cells[18], &b.cells[27], &b.cells[36], &b.cells[45], &b.cells[54], &b.cells[63], &b.cells[72]},
		{&b.cells[1], &b.cells[10], &b.cells[19], &b.cells[28], &b.cells[37], &b.cells[46], &b.cells[55], &b.cells[64], &b.cells[73]},
		{&b.cells[2], &b.cells[11], &b.cells[20], &b.cells[29], &b.cells[38], &b.cells[47], &b.cells[56], &b.cells[65], &b.cells[74]},
		{&b.cells[3], &b.cells[12], &b.cells[21], &b.cells[30], &b.cells[39], &b.cells[48], &b.cells[57], &b.cells[66], &b.cells[75]},
		{&b.cells[4], &b.cells[13], &b.cells[22], &b.cells[31], &b.cells[40], &b.cells[49], &b.cells[58], &b.cells[67], &b.cells[76]},
		{&b.cells[5], &b.cells[14], &b.cells[23], &b.cells[32], &b.cells[41], &b.cells[50], &b.cells[59], &b.cells[68], &b.cells[77]},
		{&b.cells[6], &b.cells[15], &b.cells[24], &b.cells[33], &b.cells[42], &b.cells[51], &b.cells[60], &b.cells[69], &b.cells[78]},
		{&b.cells[7], &b.cells[16], &b.cells[25], &b.cells[34], &b.cells[43], &b.cells[52], &b.cells[61], &b.cells[70], &b.cells[79]},
		{&b.cells[8], &b.cells[17], &b.cells[26], &b.cells[35], &b.cells[44], &b.cells[53], &b.cells[62], &b.cells[71], &b.cells[80]},
	}
	b.boxes = [9]Cells{
		{&b.cells[0], &b.cells[1], &b.cells[2], &b.cells[9], &b.cells[10], &b.cells[11], &b.cells[18], &b.cells[19], &b.cells[20]},
		{&b.cells[3], &b.cells[4], &b.cells[5], &b.cells[12], &b.cells[13], &b.cells[14], &b.cells[21], &b.cells[22], &b.cells[23]},
		{&b.cells[6], &b.cells[7], &b.cells[8], &b.cells[15], &b.cells[16], &b.cells[17], &b.cells[24], &b.cells[25], &b.cells[26]},
		{&b.cells[27], &b.cells[28], &b.cells[29], &b.cells[36], &b.cells[37], &b.cells[38], &b.cells[45], &b.cells[46], &b.cells[47]},
		{&b.cells[30], &b.cells[31], &b.cells[32], &b.cells[39], &b.cells[40], &b.cells[41], &b.cells[48], &b.cells[49], &b.cells[50]},
		{&b.cells[33], &b.cells[34], &b.cells[35], &b.cells[42], &b.cells[43], &b.cells[44], &b.cells[51], &b.cells[52], &b.cells[53]},
		{&b.cells[54], &b.cells[55], &b.cells[56], &b.cells[63], &b.cells[64], &b.cells[65], &b.cells[72], &b.cells[73], &b.cells[74]},
		{&b.cells[57], &b.cells[58], &b.cells[59], &b.cells[66], &b.cells[67], &b.cells[68], &b.cells[75], &b.cells[76], &b.cells[77]},
		{&b.cells[60], &b.cells[61], &b.cells[62], &b.cells[69], &b.cells[70], &b.cells[71], &b.cells[78], &b.cells[79], &b.cells[80]},
	}

	for _, row := range b.rows {
		b.blocks = append(b.blocks, row)
	}
	for _, col := range b.cols {
		b.blocks = append(b.blocks, col)
	}
	for _, box := range b.boxes {
		b.blocks = append(b.blocks, box)
	}
}

func (b *Board) solution() string {
	solution := ""
	for _, cell := range b.cells {
		if len(cell.possibles) == 1 {
			solution += strings.Join(cell.possibles, "")
		} else {
			solution = solution + "."
		}

	}
	return solution
}

func (b *Board) solved() bool {
	solved := true
	for _, cell := range b.cells {
		solved = solved && cell.solved()
	}
	return solved
}

func (b *Board) removeSolved() {
	removed := false
	found := true
	for found {
		found = false
		for _, block := range b.blocks {
			solved := []string{}
			for _, cell := range block {
				if cell.solved() {
					solved = append(solved, cell.possibles[0])
				}
			}
			if block.remove(solved) {
				found = true
				removed = true
			}
		}
	}
	if removed {
		b.print()
	}
}

func (b *Board) parse(s string) {
	if len(s) != 81 {
		fmt.Printf("!!! Parse expected length 81 but got %v\n", len(s))
	}
	for i := 0; i < len(s); i++ {
		c := string(s[i])
		if c == "." {
			b.cells[i].possibles = numbers
		} else {
			b.cells[i].possibles = []string{c}
		}
	}
}

func (b *Board) print() {
	if !b.verbose {
		return
	}

	fmt.Println()
	for i := range b.cells {
		if b.cells[i].solved() {
			fmt.Printf("%-10s", "    "+b.cells[i].possibles[0])
		} else {
			fmt.Printf("%-10s", strings.Join(b.cells[i].possibles, ""))
		}

		if i > 0 {
			j := i + 1
			if j%9 == 0 {
				fmt.Println()
				if j == 27 || j == 54 {
					fmt.Println("==============================|==============================|==============================")
				}
			} else if j%3 == 0 {
				fmt.Print("|")
			}
		}
	}
	fmt.Println()
}

// If a cell is the only one to contain a possible then it is the solution.
func (b *Board) singles() bool {
	for index, cells := range b.blocks {
		for _, possible := range numbers {
			matches := cells.filterHasPossible(possible)
			if len(matches) == 1 {
				if b.verbose {
					fmt.Printf("Single %s in %v\n", possible, nameOfBlock(index))
				}
				matches[0].solve(possible)
				return true
			}
		}
	}
	return false
}

// If two or more cells are the only ones to contain a combo then the combo can be removed from other cells.
func (b *Board) nakeds() bool {
	for index, block := range b.blocks {
		possibles := block.possibles()
		combinations := combosOfString(possibles, 2)
		for _, combo := range combinations {
			// Need to match combo is [1,2,3] and value is [1] or [2,3] etc.
			comboFlavours := combosOfString(combo, 1)
			hasCombo := func(cell *Cell) bool {
				for _, c := range comboFlavours {
					if sliceOfStringsEqual(c, cell.possibles) {
						return true
					}
				}
				return false
			}
			matches := block.filterInclude(hasCombo)
			if len(matches) == len(combo) {
				inMatches := func(cell *Cell) bool {
					for _, match := range matches {
						if cell == match {
							return true
						}
					}
					return false
				}
				others := block.filterExclude(inMatches)
				found := others.remove(combo)
				if found {
					if b.verbose {
						fmt.Printf("Naked %v found in %s\n", combo, nameOfBlock(index))
					}
					return true
				}
			}
		}
	}
	return false
}

// If two or more cells are the only ones to contain a combo then any other possibles can be removed from those cells.
func (b *Board) hiddens() bool {
	for index, block := range b.blocks {
		possibles := block.possibles()
		combinations := combosOfString(possibles, 2)
		for _, combo := range combinations {
			hasCombo := func(cell *Cell) bool {
				return cell.hasAnyOf(combo)
			}
			matches := block.filterInclude(hasCombo)
			if len(matches) == len(combo) {
				found := false
				for _, match := range matches {
					if match.removePossiblesApartFrom(combo) {
						found = true
					}
				}
				if found {
					if b.verbose {
						fmt.Printf("Hidden %v found in %s\n", combo, nameOfBlock(index))
					}
					return true
				}
			}

		}
	}
	return false
}

// If 2 or 3 cells in a box have a possible only in the same row/col then it can be removed from the rest of that row/col.
func (b *Board) pointingPairs() bool {
	for i, box := range b.boxes {
		cellInBox := func(cell *Cell) bool {
			return cell.box == i
		}
		for _, number := range numbers {
			matches := box.filterHasPossible(number)
			scanRow := false
			scanCol := false
			if len(matches) == 2 {
				if matches[0].row == matches[1].row {
					scanRow = true
				} else if matches[0].col == matches[1].col {
					scanCol = true
				}
			} else if len(matches) == 3 {
				if matches[0].row == matches[1].row && matches[0].row == matches[2].row {
					scanRow = true
				} else if matches[0].col == matches[1].col && matches[0].col == matches[2].col {
					scanCol = true
				}
			}

			if scanRow {
				row := matches[0].row
				others := b.rows[row].filterExclude(cellInBox)
				if others.remove([]string{number}) {
					if b.verbose {
						fmt.Printf("Pointing pair: %v in box %v row %v\n", number, i+1, row+1)
					}
					return true
				}
			}

			if scanCol {
				col := matches[0].col
				others := b.cols[col].filterExclude(cellInBox)
				if others.remove([]string{number}) {
					if b.verbose {
						fmt.Printf("Pointing pair: %v in box %v col %v\n", number, i+1, col+1)
					}
					return true
				}
			}

		}
	}
	return false
}

// If 2 or 3 cells in a row/col have a possible only in the same box then it can be removed from the rest of that box.
func (b *Board) boxLineReduction() bool {
	for i, row := range b.rows {
		cellInRow := func(cell *Cell) bool {
			return cell.row == i
		}
		for _, number := range numbers {
			matches := row.filterHasPossible(number)
			if len(matches) == 0 {
				continue
			}
			box := matches[0].box
			if (len(matches) == 2 && matches[1].box == box) ||
				(len(matches) == 3 && matches[1].box == box && matches[2].box == box) {
				others := b.boxes[box].filterExclude(cellInRow)
				if others.remove([]string{number}) {
					if b.verbose {
						fmt.Printf("Box Line Reduction: %v in box %v row %v\n", number, box+1, i+1)
					}
					return true
				}

			}
		}
	}

	for i, col := range b.cols {
		cellInCol := func(cell *Cell) bool {
			return cell.col == i
		}
		for _, number := range numbers {
			matches := col.filterHasPossible(number)
			if len(matches) == 0 {
				continue
			}
			box := matches[0].box
			if (len(matches) == 2 && matches[1].box == box) ||
				(len(matches) == 3 && matches[1].box == box && matches[2].box == box) {
				others := b.boxes[box].filterExclude(cellInCol)
				if others.remove([]string{number}) {
					if b.verbose {
						fmt.Printf("Box Line Reduction: %v in box %v col %v\n", number, box+1, i+1)
					}
					return true
				}

			}
		}
	}
	return false
}

// If there are only 2 possibles in the same two columns of two rows, i.e. forming an X, then the possible
// can be removed from the rest of the cells in those two columns.
//
// Repeat swapping rows and colums.
func (b *Board) xwing() bool {
	combos := combosOfInt(indexes, 2)
	for _, possible := range numbers {
		for _, combo := range combos {
			matchedRows := []int{}
			for i, row := range b.rows {
				matchedCells := row.filterHasPossible(possible)
				if len(matchedCells) == 2 &&
					matchedCells[0].col == combo[0] &&
					matchedCells[1].col == combo[1] {
					matchedRows = append(matchedRows, i)
				}
			}
			if len(matchedRows) == 2 {
				others := Cells{}
				for _, cell := range b.cols[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range b.cols[combo[1]] {
					others = append(others, cell)
				}
				inRow := func(cell *Cell) bool {
					return cell.inRow(matchedRows[0]) || cell.inRow(matchedRows[1])
				}
				others = others.filterExclude(inRow)
				if others.remove([]string{possible}) {
					if b.verbose {
						fmt.Printf("XWing %v in rows %v,%v cols %v,%v\n", possible, matchedRows[0]+1, matchedRows[1]+1, combo[0]+1, combo[1]+1)
					}
					return true
				}
			}

			matchedCols := []int{}
			for i, col := range b.cols {
				matchedCells := col.filterHasPossible(possible)
				if len(matchedCells) == 2 &&
					matchedCells[0].row == combo[0] &&
					matchedCells[1].row == combo[1] {
					matchedCols = append(matchedCols, i)
				}
			}
			if len(matchedCols) == 2 {
				others := Cells{}
				for _, cell := range b.rows[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range b.rows[combo[1]] {
					others = append(others, cell)
				}
				inCol := func(cell *Cell) bool {
					return cell.inCol(matchedCols[0]) || cell.inCol(matchedCols[1])
				}
				others = others.filterExclude(inCol)
				if others.remove([]string{possible}) {
					if b.verbose {
						fmt.Printf("XWing %v in cols %v,%v rows %v,%v\n", possible, matchedCols[0]+1, matchedCols[1]+1, combo[0]+1, combo[1]+1)
					}
					return true
				}
			}
		}
	}
	return false
}

// Swordfish is the 3 row/col variant of XWing.
// Note: It has to handle 3 and 2 row/col combinations.
func (b *Board) swordfish() bool {
	combos := combosOfInt(indexes, 3)
	for _, possible := range numbers {
		for _, combo := range combos {
			matchedRows := []int{}
			for i, row := range b.rows {
				matchedCells := row.filterHasPossible(possible)
				if len(matchedCells) == 3 &&
					matchedCells[0].col == combo[0] &&
					matchedCells[1].col == combo[1] &&
					matchedCells[2].col == combo[2] {
					matchedRows = append(matchedRows, i)
				} else if len(matchedCells) == 2 &&
					matchedCells[0].col == combo[0] &&
					matchedCells[1].col == combo[1] {
					matchedRows = append(matchedRows, i)
				} else if len(matchedCells) == 2 &&
					matchedCells[0].col == combo[1] &&
					matchedCells[1].col == combo[2] {
					matchedRows = append(matchedRows, i)
				} else if len(matchedCells) == 2 &&
					matchedCells[0].col == combo[0] &&
					matchedCells[1].col == combo[2] {
					matchedRows = append(matchedRows, i)
				}
			}
			if len(matchedRows) == 3 {
				others := Cells{}
				for _, cell := range b.cols[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range b.cols[combo[1]] {
					others = append(others, cell)
				}
				for _, cell := range b.cols[combo[2]] {
					others = append(others, cell)
				}
				inRow := func(cell *Cell) bool {
					return cell.inRow(matchedRows[0]) || cell.inRow(matchedRows[1]) || cell.inRow(matchedRows[2])
				}
				others = others.filterExclude(inRow)
				if others.remove([]string{possible}) {
					if b.verbose {
						fmt.Printf("Swordfish %v in rows %v,%v,%v cols %v,%v,%v\n", possible, matchedRows[0]+1, matchedRows[1]+1, matchedRows[2]+1,
							combo[0]+1, combo[1]+1, combo[2]+1)
					}
					return true
				}
			}

			matchedCols := []int{}
			for i, col := range b.cols {
				matchedCells := col.filterHasPossible(possible)
				if len(matchedCells) == 3 &&
					matchedCells[0].row == combo[0] &&
					matchedCells[1].row == combo[1] &&
					matchedCells[2].row == combo[2] {
					matchedCols = append(matchedCols, i)
				} else if len(matchedCells) == 2 &&
					matchedCells[0].row == combo[0] &&
					matchedCells[1].row == combo[1] {
					matchedCols = append(matchedCols, i)
				} else if len(matchedCells) == 2 &&
					matchedCells[0].row == combo[1] &&
					matchedCells[1].row == combo[2] {
					matchedCols = append(matchedCols, i)
				} else if len(matchedCells) == 2 &&
					matchedCells[0].row == combo[0] &&
					matchedCells[1].row == combo[2] {
					matchedCols = append(matchedCols, i)
				}
			}
			if len(matchedCols) == 3 {
				others := Cells{}
				for _, cell := range b.rows[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range b.rows[combo[1]] {
					others = append(others, cell)
				}
				for _, cell := range b.rows[combo[2]] {
					others = append(others, cell)
				}
				inCol := func(cell *Cell) bool {
					return cell.inCol(matchedCols[0]) || cell.inCol(matchedCols[1]) || cell.inCol(matchedCols[2])
				}
				others = others.filterExclude(inCol)
				if others.remove([]string{possible}) {
					if b.verbose {
						fmt.Printf("Swordfish %v in cols %v,%v,%v rows %v,%v,%v\n", possible, matchedCols[0]+1, matchedCols[1]+1, matchedCols[2]+1,
							combo[0]+1, combo[1]+1, combo[2]+1)
					}
					return true
				}
			}
		}
	}
	return false
}

// All pairs of possibles are grouped into chains. These are then coloured in alternate colours for odd/even.
// Any cell that can see two different colours can be removed.
// This includes cells in the chain but a cell cannot see itself.
func (b *Board) simplecolouring() bool {
	for _, possible := range numbers {
		pairs := []Cells{}
		for _, block := range b.blocks {
			matches := block.filterHasPossible(possible)
			if len(matches) == 2 {
				pairs = append(pairs, matches)
			}
		}
		if len(pairs) < 2 {
			return false
		}
		chains := createChainsFrom(pairs)
		for _, chain := range chains {
			for _, row := range b.rows {
				cells := row.filterHasPossible(possible)
				for _, cell := range cells {
					if cell.canSee(chain) == 3 {
						if b.verbose {
							fmt.Printf("Simple Colouring %v: %v can see two colours in %v\n", possible, cell, chain)
						}
						cell.removePossible(possible)
						return true
					}
				}
			}
		}
	}
	return false
}

//------- General functions ----------

func unique(values []string) []string {
	result := make([]string, 0, 9)
	m := make(map[string]bool)
	for _, value := range values {
		m[value] = true
	}
	for k := range m {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

// Found on topcoder
// Imagine all numbers from 0 to 2^len-1
// The bit patterns of these numbers are the combinations
func combosOfString(elems []string, min int) [][]string {
	result := [][]string{}
	n := len(elems)
	for num := 0; num < (1 << uint(n)); num++ {
		combination := make([]string, 0, 9)
		for ndx := 0; ndx < n; ndx++ {
			// (is the bit "on" in this number?)
			if num&(1<<uint(ndx)) != 0 {
				// (then add it to the combination)
				combination = append(combination, elems[ndx])
			}
		}
		if len(combination) >= min {
			result = append(result, combination)
		}
	}
	return result
}

func combosOfInt(elems []int, size int) [][]int {
	result := [][]int{}
	n := len(elems)
	for num := 0; num < (1 << uint(n)); num++ {
		combination := make([]int, 0, 9)
		for ndx := 0; ndx < n; ndx++ {
			// (is the bit "on" in this number?)
			if num&(1<<uint(ndx)) != 0 {
				// (then add it to the combination)
				combination = append(combination, elems[ndx])
			}
		}
		if len(combination) == size {
			result = append(result, combination)
		}
	}
	return result
}

func nameOfBlock(block int) string {
	if block < 9 {
		return fmt.Sprintf("Row %v", block+1)
	} else if block < 18 {
		return fmt.Sprintf("Col %v", block-9+1)
	} else {
		return fmt.Sprintf("Box %v", block-18+1)
	}
}

func sliceOfStringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func createChainsFrom(pairs []Cells) []Chain {
	// assign the pairs to chains
	chains := []Chain{}
	for len(pairs) > 0 {
		matchedChain := false
		for chainIndex := range chains {
			for i, pair := range pairs {
				c1 := chains[chainIndex].findCell(pair[0])
				c2 := chains[chainIndex].findCell(pair[1])
				if c1 != nil {
					if c1.colour == 1 {
						c1.links = append(c1.links, Chain{cell: pair[1], colour: 2, links: []Chain{}})
					} else {
						c1.links = append(c1.links, Chain{cell: pair[1], colour: 1, links: []Chain{}})
					}
				} else if c2 != nil {
					if c2.colour == 1 {
						c2.links = append(c2.links, Chain{cell: pair[0], colour: 2, links: []Chain{}})
					} else {
						c2.links = append(c2.links, Chain{cell: pair[0], colour: 1, links: []Chain{}})
					}
				}
				if c1 != nil || c2 != nil {
					matchedChain = true
					pairs = append(pairs[:i], pairs[i+1:]...)
					break
				}
			}
		}
		if !matchedChain {
			pair := pairs[0]
			pairs = pairs[1:]
			c1 := Chain{cell: pair[0], colour: 1, links: []Chain{}}
			c2 := Chain{cell: pair[1], colour: 2, links: []Chain{}}
			c1.links = append(c1.links, c2)
			chains = append(chains, c1)
		}
	}
	return chains
}

// ---------------------------------------

func solvePuzzle(puzzle string) (bool, string) {
	var b = new(Board)
	strategies := []func() bool{
		b.singles,
		b.nakeds,
		b.hiddens,
		b.pointingPairs,
		b.boxLineReduction,
		b.xwing,
		b.swordfish,
		b.simplecolouring,
	}
	b.init()
	b.parse(puzzle)
	b.print()
	b.removeSolved()
	if b.solved() {
		return true, b.solution()
	}
	for {
		found := false
		for _, strategy := range strategies {
			if strategy() {
				found = true
				if b.solved() {
					return true, b.solution()
				}
			}
			if found {
				b.print()
				b.removeSolved()
				if b.solved() {
					return true, b.solution()
				}
				break
			}
		}

		if !found {
			if b.verbose {
				fmt.Println("Beats me !!!")
			}
			return false, b.solution()
		}
	}
}

func loadFile(name string) ([]string, bool) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return nil, true
	}
	defer file.Close()

	var puzzles = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		puzzles = append(puzzles, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, true
	}
	return puzzles, false
}

func main() {

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	status := ""

	puzzles, done := loadFile("top95.txt")
	if done {
		return
	}

	expected, done := loadFile("top95expected.txt")
	if done {
		return
	}

	var solved, total int
	for index, puzzle := range puzzles {
		total++
		fmt.Printf("### Puzzle %v ###\n", total)
		fmt.Printf("%s\n", puzzle)
		solvedIt, solution := solvePuzzle(puzzle)
		fmt.Printf("%s\n", solution)
		if solvedIt {
			status += "S"
			solved++
			if solution != expected[index] {
				fmt.Printf("Incorrect solution in puzzle %v\n", index+1)
				fmt.Println(solution)
				fmt.Println(expected[index])
				return
			}
		} else {
			status += "."
		}
		if total%10 == 0 {
			status += "\n"
		}
	}

	fmt.Printf("Solved %v out of %v\n", solved, total)
	fmt.Println(status)
}
