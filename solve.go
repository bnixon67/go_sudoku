/*
Copyright 2017 Bill Nixon

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func isValid(num byte, puzzle [9][9]byte, row byte, col byte) bool {

	// check if num already exists in row or column
	for i := byte(0); i < 9; i++ {
		// does num already exist in row?
		if puzzle[row][i] == num {
			return false
		}

		// does num already exist in column?
		if puzzle[i][col] == num {
			return false
		}
	}

	// check 9x9 grid
	topRow := row / 3 * 3
	topCol := col / 3 * 3

	for i := topRow; i < topRow+3; i++ {
		for j := topCol; j < topCol+3; j++ {
			if !(i == row && j == col) {
				if puzzle[i][j] == num {
					return false
				}
			}
		}
	}

	return true
}

func solve(puzzle *[9][9]byte, row byte, col byte) bool {

	// past last row?
	if row == 9 {
		return true
	}

	// puzzle cell already has a value, move to next cell
	if puzzle[row][col] != 0 {
		if col == 8 {
			if solve(puzzle, row+1, 0) {
				return true
			}
		} else {
			if solve(puzzle, row, col+1) {
				return true
			}
		}
	} else {
		for num := byte(1); num < 10; num++ {
			if isValid(num, *puzzle, row, col) {
				puzzle[row][col] = num
				if col == 8 {
					if solve(puzzle, row+1, 0) {
						return true
					}
				} else {
					if solve(puzzle, row, col+1) {
						return true
					}
				}
				puzzle[row][col] = 0
			}
		}
	}

	return false
}

func print_puzzle(puzzle [9][9]byte) {

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			fmt.Print(puzzle[row][col])
		}
		fmt.Println()
	}

}

func main() {

	var puzzle [9][9]byte

	if len(os.Args) != 2 {
		_, file := filepath.Split(os.Args[0])
		fmt.Println("usage: ", file, "filename")
		os.Exit(1)
	}

	filename := os.Args[1]
	fmt.Println("File is", filename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	//reader := bufio.NewReader(os.Stdin)

	for row := 0; row < 9; row++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		for col, rune := range line[:9] {
			puzzle[row][col] = byte(rune - '0')
		}
	}

	fmt.Println("Puzzle to solve:")
	print_puzzle(puzzle)
	fmt.Println()

	if solve(&puzzle, 0, 0) {
		fmt.Println("Puzzle solved:")
		print_puzzle(puzzle)
	} else {
		fmt.Println("Cannot solve this puzzle")
	}

}
