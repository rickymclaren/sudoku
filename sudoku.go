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
	"fmt"
	"os"
	"sort"
	"strings"
)

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

// --- Methods of Cell ---
func (cell Cell) solved() bool {
	return len(cell.possibles) == 1
}

func (cell Cell) inRow(row int) bool {
	return cell.row == row
}

func (cell Cell) inCol(col int) bool {
	return cell.col == col
}

func (cell Cell) inBox(box int) bool {
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
	newPossibles := []string{}
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

func (cell *Cell) hasAllOf(possibles []string) bool {
	if cell.solved() {
		return false
	}
	for _, possible := range possibles {
		if !cell.hasPossible(possible) {
			return false
		}
	}
	return true
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

// --- Methods of Cells ---

func (cells Cells) remove(values []string) bool {
	found := false
	for _, cell := range cells {
		if cell.removePossibles(values) {
			found = true
		}
	}
	return found
}

func (cells Cells) filterInclude(include func(*Cell) bool) Cells {
	result := []*Cell{}
	for _, c := range cells {
		if include(c) {
			result = append(result, c)
		}
	}
	return result
}

func (cells Cells) filterExclude(include func(*Cell) bool) Cells {
	result := []*Cell{}
	for _, c := range cells {
		if !include(c) {
			result = append(result, c)
		}
	}
	return result
}

func (cells Cells) filterHasPossible(s string) Cells {
	result := []*Cell{}
	for _, c := range cells {
		if c.hasPossible(s) {
			result = append(result, c)
		}
	}
	return result
}

func (cells Cells) possibles() []string {
	result := []string{}
	for _, cell := range cells {
		if !cell.solved() {
			for _, possible := range cell.possibles {
				result = append(result, possible)
			}
		}
	}
	return unique(result)
}

//------- General functions ----------

func unique(values []string) []string {
	result := []string{}
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

func solution() string {
	solution := ""
	for _, cell := range b {
		solution += strings.Join(cell.possibles, "")
	}
	return solution
}

//-------------------------------------------

var b [81]Cell
var rows = []Cells{
	{&b[0], &b[1], &b[2], &b[3], &b[4], &b[5], &b[6], &b[7], &b[8]},
	{&b[9], &b[10], &b[11], &b[12], &b[13], &b[14], &b[15], &b[16], &b[17]},
	{&b[18], &b[19], &b[20], &b[21], &b[22], &b[23], &b[24], &b[25], &b[26]},
	{&b[27], &b[28], &b[29], &b[30], &b[31], &b[32], &b[33], &b[34], &b[35]},
	{&b[36], &b[37], &b[38], &b[39], &b[40], &b[41], &b[42], &b[43], &b[44]},
	{&b[45], &b[46], &b[47], &b[48], &b[49], &b[50], &b[51], &b[52], &b[53]},
	{&b[54], &b[55], &b[56], &b[57], &b[58], &b[59], &b[60], &b[61], &b[62]},
	{&b[63], &b[64], &b[65], &b[66], &b[67], &b[68], &b[69], &b[70], &b[71]},
	{&b[72], &b[73], &b[74], &b[75], &b[76], &b[77], &b[78], &b[79], &b[80]},
}
var cols = []Cells{
	{&b[0], &b[9], &b[18], &b[27], &b[36], &b[45], &b[54], &b[63], &b[72]},
	{&b[1], &b[10], &b[19], &b[28], &b[37], &b[46], &b[55], &b[64], &b[73]},
	{&b[2], &b[11], &b[20], &b[29], &b[38], &b[47], &b[56], &b[65], &b[74]},
	{&b[3], &b[12], &b[21], &b[30], &b[39], &b[48], &b[57], &b[66], &b[75]},
	{&b[4], &b[13], &b[22], &b[31], &b[40], &b[49], &b[58], &b[67], &b[76]},
	{&b[5], &b[14], &b[23], &b[32], &b[41], &b[50], &b[59], &b[68], &b[77]},
	{&b[6], &b[15], &b[24], &b[33], &b[42], &b[51], &b[60], &b[69], &b[78]},
	{&b[7], &b[16], &b[25], &b[34], &b[43], &b[52], &b[61], &b[70], &b[79]},
	{&b[8], &b[17], &b[26], &b[35], &b[44], &b[53], &b[62], &b[71], &b[80]},
}
var boxes = []Cells{
	{&b[0], &b[1], &b[2], &b[9], &b[10], &b[11], &b[18], &b[19], &b[20]},
	{&b[3], &b[4], &b[5], &b[12], &b[13], &b[14], &b[21], &b[22], &b[23]},
	{&b[6], &b[7], &b[8], &b[15], &b[16], &b[17], &b[24], &b[25], &b[26]},
	{&b[27], &b[28], &b[29], &b[36], &b[37], &b[38], &b[45], &b[46], &b[47]},
	{&b[30], &b[31], &b[32], &b[39], &b[40], &b[41], &b[48], &b[49], &b[50]},
	{&b[33], &b[34], &b[35], &b[42], &b[43], &b[44], &b[51], &b[52], &b[53]},
	{&b[54], &b[55], &b[56], &b[63], &b[64], &b[65], &b[72], &b[73], &b[74]},
	{&b[57], &b[58], &b[59], &b[66], &b[67], &b[68], &b[75], &b[76], &b[77]},
	{&b[60], &b[61], &b[62], &b[69], &b[70], &b[71], &b[78], &b[79], &b[80]},
}
var blocks = []Cells{}
var numbers []string = strings.Split("123456789", "")
var indexes = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
var combinations = [][]string{}

func init() {
	for i := range b {
		row := i / 9
		col := i % 9
		box := (row / 3 * 3) + col/3
		b[i] = Cell{possibles: numbers, row: row, col: col, box: box}
	}
	for _, row := range rows {
		blocks = append(blocks, row)
	}
	for _, col := range cols {
		blocks = append(blocks, col)
	}
	for _, box := range boxes {
		blocks = append(blocks, box)
	}
	combinations = combosOfString(numbers, 2)
}

// Found on topcoder
// Imagine all numbers from 0 to 2^len-1
// The bit patterns of these numbers are the combinations
func combosOfString(elems []string, min int) [][]string {
	result := [][]string{}
	n := len(elems)
	for num := 0; num < (1 << uint(n)); num++ {
		combination := []string{}
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
		combination := []int{}
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

func parse(s string) {
	if len(s) != 81 {
		fmt.Printf("!!! Parse expected length 81 but got %v\n", len(s))
	}
	for i := 0; i < len(s); i++ {
		c := string(s[i])
		if c == "." {
			b[i].possibles = numbers
		} else {
			b[i].possibles = []string{c}
		}
	}
}

func boardSolved() bool {
	solved := true
	for _, cell := range b {
		solved = solved && cell.solved()
	}
	return solved
}

func printb() {
	fmt.Println()
	for i := range b {
		if b[i].solved() {
			fmt.Printf("%-10s", "    "+b[i].possibles[0])
		} else {
			fmt.Printf("%-10s", strings.Join(b[i].possibles, ""))
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

func nameOfBlock(block int) string {
	if block < 9 {
		return fmt.Sprintf("Row %v", block+1)
	} else if block < 18 {
		return fmt.Sprintf("Col %v", block-9+1)
	} else {
		return fmt.Sprintf("Box %v", block-18+1)
	}
}

func signatureOfBlock(i int) string {
	block := blocks[i]
	result := nameOfBlock(i) + ":"
	for _, cell := range block {
		result += strings.Join(cell.possibles, "") + "|"
	}
	return result
}

func removeSolved() {
	removed := false
	found := true
	for found {
		found = false
		for i, block := range blocks {
			oldSignature := signatureOfBlock(i)
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
			newSignature := signatureOfBlock(i)
			if newSignature != oldSignature {
				// fmt.Printf("%s -> %s\n", oldSignature, newSignature)
			}
		}
	}
	if removed {
		fmt.Println("Removed Solved")
		printb()
	}
}

// If a cell is the only one to contain a possible then it is the solution.
func singles() bool {
	for index, cells := range blocks {
		for _, possible := range numbers {
			matches := cells.filterHasPossible(possible)
			if len(matches) == 1 {
				fmt.Printf("Single %s in %v\n", possible, nameOfBlock(index))
				matches[0].solve(possible)
				return true
			}
		}
	}
	return false
}

// If two or more cells are the only ones to contain a combo then the combo can be removed from other cells.
func nakeds() bool {
	for _, combo := range combinations {
		// Need to match combo is [1,2,3] and value is [1] or [2,3] etc.
		comboFlavours := combosOfString(combo, 1)
		hasCombo := func(cell *Cell) bool {
			possibles := strings.Join(cell.possibles, "")
			for _, c := range comboFlavours {
				if possibles == strings.Join(c, "") {
					return true
				}
			}
			return false
		}
		for index, block := range blocks {
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
					fmt.Printf("Naked %v found in %s\n", combo, nameOfBlock(index))
					return true
				}
			}
		}
	}
	return false
}

// If two or more cells are the only ones to contain a combo then any other possibles can be removed from those cells.
func hiddens() bool {
	for index, block := range blocks {
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
					fmt.Printf("Hidden %v found in %s\n", combo, nameOfBlock(index))
					return true
				}
			}

		}
	}
	return false
}

// If 2 or 3 cells in a box have a possible only in the same row/col then it can be removed from the rest of that row/col.
func pointingPairs() bool {
	for i, box := range boxes {
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
				others := rows[row].filterExclude(cellInBox)
				if others.remove([]string{number}) {
					fmt.Printf("Pointing pair: %v in box %v row %v\n", number, i+1, row+1)
					return true
				}
			}

			if scanCol {
				col := matches[0].col
				others := cols[col].filterExclude(cellInBox)
				if others.remove([]string{number}) {
					fmt.Printf("Pointing pair: %v in box %v col %v\n", number, i+1, col+1)
					return true
				}
			}

		}
	}
	return false
}

// If 2 or 3 cells in a row/col have a possible only in the same box then it can be removed from the rest of that box.
func boxLineReduction() bool {
	for i, row := range rows {
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
				others := boxes[box].filterExclude(cellInRow)
				if others.remove([]string{number}) {
					fmt.Printf("Box Line Reduction: %v in box %v row %v\n", number, box+1, i+1)
					return true
				}

			}
		}
	}

	for i, col := range cols {
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
				others := boxes[box].filterExclude(cellInCol)
				if others.remove([]string{number}) {
					fmt.Printf("Box Line Reduction: %v in box %v col %v\n", number, box+1, i+1)
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
func xwing() bool {
	combos := combosOfInt(indexes, 2)
	for _, possible := range numbers {
		for _, combo := range combos {
			matchedRows := []int{}
			for i, row := range rows {
				matchedCells := row.filterHasPossible(possible)
				if len(matchedCells) == 2 &&
					matchedCells[0].col == combo[0] &&
					matchedCells[1].col == combo[1] {
					matchedRows = append(matchedRows, i)
				}
			}
			if len(matchedRows) == 2 {
				others := Cells{}
				for _, cell := range cols[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range cols[combo[1]] {
					others = append(others, cell)
				}
				inRow := func(cell *Cell) bool {
					return cell.inRow(matchedRows[0]) || cell.inRow(matchedRows[1])
				}
				others = others.filterExclude(inRow)
				if others.remove([]string{possible}) {
					fmt.Printf("XWing %v in rows %v,%v cols %v,%v\n", possible, matchedRows[0]+1, matchedRows[1]+1, combo[0]+1, combo[1]+1)
					return true
				}
			}

			matchedCols := []int{}
			for i, col := range cols {
				matchedCells := col.filterHasPossible(possible)
				if len(matchedCells) == 2 &&
					matchedCells[0].row == combo[0] &&
					matchedCells[1].row == combo[1] {
					matchedCols = append(matchedCols, i)
				}
			}
			if len(matchedCols) == 2 {
				others := Cells{}
				for _, cell := range rows[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range rows[combo[1]] {
					others = append(others, cell)
				}
				inCol := func(cell *Cell) bool {
					return cell.inCol(matchedCols[0]) || cell.inCol(matchedCols[1])
				}
				others = others.filterExclude(inCol)
				if others.remove([]string{possible}) {
					fmt.Printf("XWing %v in cols %v,%v rows %v,%v\n", possible, matchedCols[0]+1, matchedCols[1]+1, combo[0]+1, combo[1]+1)
					return true
				}
			}
		}
	}
	return false
}

// Swordfish is the 3 row/col variant of XWing.
func swordfish() bool {
	combos := combosOfInt(indexes, 3)
	for _, possible := range numbers {
		for _, combo := range combos {
			matchedRows := []int{}
			for i, row := range rows {
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
				for _, cell := range cols[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range cols[combo[1]] {
					others = append(others, cell)
				}
				for _, cell := range cols[combo[2]] {
					others = append(others, cell)
				}
				inRow := func(cell *Cell) bool {
					return cell.inRow(matchedRows[0]) || cell.inRow(matchedRows[1]) || cell.inRow(matchedRows[2])
				}
				others = others.filterExclude(inRow)
				if others.remove([]string{possible}) {
					fmt.Printf("Swordfish %v in rows %v,%v,%v cols %v,%v,%v\n", possible, matchedRows[0]+1, matchedRows[1]+1, matchedRows[2]+1,
						combo[0]+1, combo[1]+1, combo[2]+1)
					return true
				}
			}

			matchedCols := []int{}
			for i, col := range cols {
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
				for _, cell := range rows[combo[0]] {
					others = append(others, cell)
				}
				for _, cell := range rows[combo[1]] {
					others = append(others, cell)
				}
				for _, cell := range rows[combo[2]] {
					others = append(others, cell)
				}
				inCol := func(cell *Cell) bool {
					return cell.inCol(matchedCols[0]) || cell.inCol(matchedCols[1]) || cell.inCol(matchedCols[2])
				}
				others = others.filterExclude(inCol)
				if others.remove([]string{possible}) {
					fmt.Printf("Swordfish %v in cols %v,%v,%v rows %v,%v,%v\n", possible, matchedCols[0]+1, matchedCols[1]+1, matchedCols[2]+1,
						combo[0]+1, combo[1]+1, combo[2]+1)
					return true
				}
			}
		}
	}
	return false
}

func solvePuzzle(puzzle string) (bool, string) {
	strategies := []func() bool{
		singles,
		nakeds,
		hiddens,
		pointingPairs,
		boxLineReduction,
		xwing,
		swordfish,
	}
	parse(puzzle)
	printb()
	removeSolved()
	if boardSolved() {
		fmt.Println("Done !!!")
		return true, solution()
	}
	for {
		found := false
		for _, strategy := range strategies {
			if strategy() {
				found = true
				if boardSolved() {
					fmt.Println("Done !!!")
					return true, solution()
				}
			}
			if found {
				printb()
				removeSolved()
				if boardSolved() {
					fmt.Println("Done !!!")
					return true, solution()
				}
				break
			}
		}

		if !found {
			fmt.Println("Beats me !!!")
			return false, solution()
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
		fmt.Printf(puzzle)
		solvedIt, solution := solvePuzzle(puzzle)
		if solvedIt {
			status += "S"
			solved++
			if solution != expected[index] {
				fmt.Println("Incorrect solution")
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
