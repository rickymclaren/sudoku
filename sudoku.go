package main

import (
	"fmt"
	"strings"
)

type Cell struct {
	row       int
	column    int
	value     string
	possibles []string
}

// --- Methods of Cell ---
func (c Cell) solved() bool {
	return len(c.possibles) == 0
}

func (cell Cell) inRow(row int) bool {
	return cell.row == row
}

func (cell Cell) inCol(col int) bool {
	return cell.column == col
}

func (cell *Cell) solve(value string) {
	cell.value = value
	cell.possibles = make([]string, 0)
	removeSolved()
}

func (cell *Cell) has(s string) bool {
	for _, possible := range cell.possibles {
		if possible == s {
			return true
		}
	}
	return false
}

func (cell *Cell) remove(value string) bool {
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
	if len(newPossibles) == 1 {
		cell.value = newPossibles[0]
		cell.possibles = []string{}
	} else if len(newPossibles) != len(cell.possibles){
		// fmt.Printf("Setting cell %v %v from %v to %v\n", cell.row, cell.column, cell.possibles, newPossibles)
		cell.possibles = newPossibles
	}
	return result
}

func (cell *Cell) removeValues(values []string) bool {
	result := false
	for _, value := range values {
		if cell.remove(value) {
			result = true
		}
	}
	return result
}

// -----------------------

func filter(cells []*Cell, v int, include func(*Cell, int) bool) []*Cell {
	result := []*Cell{}
	for _, c := range cells {
		if include(c, v) {
			result = append(result, c)
		}
	}
	return result
}

func filterHas(cells []*Cell, s string) []*Cell {
	result := []*Cell{}
	for _, c := range cells {
		if c.has(s) {
			result = append(result, c)
		}
	}
	return result
}

func filterCells(cells []*Cell, remove []*Cell) []*Cell {
	result := []*Cell{}
	for _, c := range cells {
		doAppend := true
		for _, r := range remove {
			if r == c {
				doAppend = false
			}
		}
		if doAppend {
			result = append(result, c)
		}
	}
	return result

}

var b [81]Cell
var rows = [][]*Cell{
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
var cols = [][]*Cell{
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
var boxes = [][]*Cell{
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
var blocks = [][]*Cell{}
var numbers []string = strings.Split("123456789", "")
var combinations = [][]string{}

func init() {
	for i, _ := range b {
		row := i / 9
		column := i % 9
		b[i] = Cell{possibles: numbers, row: row, column: column}
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
	combinations = makeCombinations(numbers)
}

func makeCombinations(elems []string) [][]string {
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
		if len(combination) > 1 {
			result = append(result, combination)
		}
	}
	return result
}

func parse(s string) {
	for i := 0; i < len(s); i++ {
		c := string(s[i])
		if c != "." {
			b[i].value = c
			b[i].possibles = make([]string, 0)
		}
	}
}

func printb() {
	fmt.Println()
	for i, _ := range b {
		if b[i].solved() {
			fmt.Printf("%-9s", "    "+b[i].value)
		} else {
			fmt.Printf("%-9s", strings.Join(b[i].possibles, ""))
		}

		if i > 0 {
			j := i + 1
			if j%9 == 0 {
				fmt.Println()
				if j == 27 || j == 54 {
					fmt.Println("===========================|===========================|===========================")
				}
			} else if j%3 == 0 {
				fmt.Print("|")
			}
		}
	}
	fmt.Println()
}

func name(i int) string {
	if i < 9 {
		return fmt.Sprintf("Row %v", i+1)
	} else if i < 18 {
		return fmt.Sprintf("Col %v", i-9+1)
	} else {
		return fmt.Sprintf("Box %v", i-18+1)
	}
}

func removeSolved() bool {
	fmt.Println("=== Remove Solved")
	result := false
	for _, cells := range blocks {
		solved := []string{}
		for _, cell := range cells {
			if cell.solved() {
				solved = append(solved, cell.value)
			}
		}
		for _, cell := range cells {
			if cell.removeValues(solved) {
				result = true
			}
		}
	}
	return result
}

func singles() bool {
	fmt.Println("=== Singles")
	for index, cells := range blocks {
		for _, r := range numbers {
			matches := filterHas(cells, r)
			if len(matches) == 1 {
				fmt.Printf("Single %s in %v\n", r, name(index))
				matches[0].solve(r)
				return true
			}
		}
	}
	return false
}

func filterForCombo(cells []*Cell, combo []string) []*Cell {
	result := []*Cell{}
	for _, cell := range cells {
		if strings.Join(cell.possibles, "") == strings.Join(combo, "") {
			result = append(result, cell)
		}
	}
	return result
}

func removeCombo(block []*Cell, matches []*Cell, combo []string) bool {
	others := filterCells(block, matches)
	found := false
	for _, cell := range others {
		if cell.removeValues(combo) {
			found = true
		}
	}
	return found
}

func nakeds() bool {
	fmt.Println("=== Nakeds")
	for _, combo := range combinations {
		for index, block := range blocks {
			matches := filterForCombo(block, combo)
			if len(matches) == len(combo) {
				found := removeCombo(block, matches, combo)
				if found {
					fmt.Printf("Naked %v found in %s\n", combo, name(index))
					return true
				}
			}
		}
	}
	return false
}

func main() {
	strategies := []func() bool{
		singles,
		nakeds,
	}
	parse(".36..54.9451..23..982...5616.7...98...346....5..287.1..4..7..963.....7....59468.2")
	printb()
	removeSolved()
	printb()
	for {
		found := false
		for _, strategy := range strategies {
			if strategy() {
				found = true
			}
			if found {
				printb()
				continue
			}
		}

		if !found {
			fmt.Println("Beats me !!!")
			return
		}
	}
}
