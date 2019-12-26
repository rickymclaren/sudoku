package main

import (
	"fmt"
)

type Cell struct {
	row       int
	column    int
	possibles string
}

// --- Methods of Cell ---
func (c Cell) solved() bool {
	return len(c.possibles) == 1
}

func (cell Cell) inRow(row int) bool {
	return cell.row == row
}

func (cell Cell) inCol(col int) bool {
	return cell.column == col
}

// -----------------------

func filter(cells []Cell, v int, include func(*Cell, int) bool) []*Cell {
	result := []*Cell{}
	for _, c := range cells {
		if include(&c, v) {
			result = append(result, &c)
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
var columns = [][]*Cell{
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
var numbers string = "12345689"

func init() {
	for i, _ := range b {
		row := i / 9
		column := i % 9
		b[i] = Cell{possibles: numbers, row: row, column: column}
	}
}

func parse(s string) {
	for i := 0; i < len(s); i++ {
		c := string(s[i])
		if c != "." {
			b[i].possibles = c
		}
	}
}

func printb() {
	for i, _ := range b {
		if b[i].solved() {
			fmt.Printf("%-9s", "    "+b[i].possibles)
		} else {
			fmt.Printf("%-9s", b[i].possibles)
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

func main() {
	parse("..6..7..8..1.3....25......9..7.58...9.......1...14.7..8......16....9.4..4..5..8..")
	printb()
}
