package main

import "fmt"

type Cell struct {
	row 	int
	column int
	value     string
	possibles string
	solved    bool
}

var board [81]Cell
var rows [9][9]*Cell
var columns [9][9]*Cell
var boxes [9][9]*Cell
var all [27][9]*Cell

func setup() {
	for i := 0; i < len(board); i++ {
		row := i / 9
		column := i % 9
		board[i] = Cell{possibles: "123456789", row: row, column: column }
	}
	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			rows[row][column] = &board[row * 9 + column]
			columns[column][row] = &board[column * 9 + row]
		}
	}
}

func parse(s string) {
	for i := 0; i < len(s); i++ {
		c := string(s[i])
		if c != "." {
			board[i].value = c
			board[i].solved = true
		}
	}
}

func printBoard() {
	for i := 0; i < len(board); i++ {
		if board[i].solved {
			fmt.Print(board[i].value)
		} else {
			fmt.Print(".")
		}

		if i > 0 {
			j := i + 1
			if j%9 == 0 {
				fmt.Println()
				if j == 27 || j == 54 {
					fmt.Println("===|===|===")
				}
			} else if j%3 == 0 {
				fmt.Print("|")
			}
		}
	}
	fmt.Println()
}

func singles() {
	for _, row := range rows {
		for _, cell := range row {
			fmt.Println(cell.row, cell.column)

		}
	}
}

func main() {
	setup()
	parse("..6..7..8..1.3....25......9..7.58...9.......1...14.7..8......16....9.4..4..5..8..")
	printBoard()
	singles()
}
