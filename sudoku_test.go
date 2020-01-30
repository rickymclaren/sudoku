package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSimpleChain(t *testing.T) {
	cells := []Cell{
		{row: 4, col: 1, box: 3},
		{row: 4, col: 7, box: 5},
		{row: 6, col: 7, box: 8},
		{row: 7, col: 7, box: 8},
	}
	pairs := []Cells{
		{&cells[0], &cells[1]},
		{&cells[1], &cells[2]},
		{&cells[2], &cells[3]},
	}
	fmt.Printf("pairs: %v\n", pairs)
	chains := createChainsFrom(pairs)
	assert.NotNil(t, chains)
	fmt.Printf("chains: %v\n", chains)
	assert.Equal(t, 1, len(chains), "Length of chains %v", chains)
	assert.Equal(t, &cells[0], chains[0].cell)
	assert.Equal(t, &cells[1], chains[0].links[0].cell)
	assert.Equal(t, &cells[2], chains[0].links[0].links[0].cell)

}

func TestCreateComplexChains(t *testing.T) {
	//     1     25        36        |24568     23568     24568     |    7     358           9
	// 69            4     369       |568       3568          7     |    2     358           1
	//     8     25            7     |1259      2359      1259      |    6     35            4
	// ==============================|==============================|==============================
	//     2         7     489       |589           1     589       |489           6         3
	//     3     19        48        |267       267       26        |48        19            5
	//     5         6     189       |89            4         3     |189           2         7
	// ==============================|==============================|==============================
	// 469       139       169       |1245679   25679     124569    |35        19            8
	// 469       89            5     |    3     689       14689     |19            7         2
	//     7     1389          2     |1589      589       1589      |35            4         6
	//
	// All the 1s from above puzzle
	// Note: 7,3 is part of the chain but can be removed because it can see two colours
	//       via 6,3 -> 6,7 -> 5,8 -> 7,8

	cells := []Cell{
		{row: 2, col: 3, box: 1}, // 0 - C3,4
		{row: 2, col: 5, box: 1}, // 1 - C3.6
		{row: 4, col: 1, box: 3}, // 2 - C5,2
		{row: 4, col: 7, box: 5}, // 3 - C5,8
		{row: 5, col: 2, box: 3}, // 4 - C6,3
		{row: 5, col: 6, box: 5}, // 5 - C6,7
		{row: 7, col: 5, box: 7}, // 6 - C8,6
		{row: 7, col: 6, box: 8}, // 7 - C8,7
		{row: 6, col: 2, box: 6}, // 8 - C7,3
		{row: 6, col: 7, box: 8}, // 9 - C7,8
	}
	pairs := []Cells{
		{&cells[0], &cells[1]}, // [C:3,4 C:3,6]
		{&cells[2], &cells[3]}, // [C:5,2 C:5,8]
		{&cells[4], &cells[5]}, // [C:6,3 C:6,7]
		{&cells[6], &cells[7]}, // [C:8,6 C:8,7]
		{&cells[4], &cells[8]}, // [C:6,3 C:7,3]
		{&cells[5], &cells[7]}, // [C:6,7 C:8,7]
		{&cells[3], &cells[9]}, // [C:5,8 C:7,8]
		{&cells[0], &cells[1]}, // [C:3,4 C:3,6] duplicate
		{&cells[2], &cells[4]}, // [C:5,2 C:6,3]
		{&cells[3], &cells[5]}, // [C:5,8 C:6,7]
		{&cells[9], &cells[7]}, // [C:7,8 C:8,7]
	}
	fmt.Printf("pairs: %v\n", pairs)
	chains := createChainsFrom(pairs)
	assert.NotNil(t, chains)
	fmt.Printf("chains: %v\n", chains)
	assert.Equal(t, 2, len(chains), "Length of chains %v", chains)

	// First two cells are not linked to anybody else
	// but there is a duplicate
	assert.Equal(t, &cells[0], chains[0].cell)
	assert.Equal(t, 2, len(chains[0].links))
	assert.Equal(t, &cells[1], chains[0].links[0].cell)
	assert.Equal(t, &cells[1], chains[0].links[1].cell)

	// Rest are one complex chain
	// 5,2 -> 5,8 -> 7,8 -> 8,7
	//            -> 6,7
	//     -> 6,3 -> 6,7 -> 8,7 -> 8,6
	//            -> 7,3
	assert.Equal(t, &cells[2], chains[1].cell)
	assert.Equal(t, &cells[3], chains[1].links[0].cell)
	assert.Equal(t, &cells[9], chains[1].links[0].links[0].cell)
	assert.Equal(t, &cells[7], chains[1].links[0].links[0].links[0].cell)

	assert.Equal(t, &cells[5], chains[1].links[0].links[1].cell)

	assert.Equal(t, &cells[4], chains[1].links[1].cell)
	assert.Equal(t, &cells[5], chains[1].links[1].links[0].cell)
	assert.Equal(t, &cells[7], chains[1].links[1].links[0].links[0].cell)
	assert.Equal(t, &cells[6], chains[1].links[1].links[0].links[0].links[0].cell)

	assert.Equal(t, &cells[8], chains[1].links[1].links[1].cell)
}
