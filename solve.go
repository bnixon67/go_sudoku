/*
Copyright 2021 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Puzzle struct {
	values [9][9]byte
}

// Stringer interafce for Puzzle type
func (puzzle Puzzle) String() string {

	var buffer bytes.Buffer

	for row := 0; row < 9; row++ {
		if row%3 == 0 {
			buffer.WriteString("+---+---+---+\n")
		}
		for col := 0; col < 9; col++ {
			if col%3 == 0 {
				buffer.WriteString("|")
			}
			buffer.WriteString(fmt.Sprint(puzzle.values[row][col]))
		}
		buffer.WriteString("|\n")
	}
	buffer.WriteString("+---+---+---+")

	return buffer.String()
}

func isValid(num byte, puzzle Puzzle, row byte, col byte) bool {

	// check if num already exists in row or column
	for i := byte(0); i < 9; i++ {
		// does num already exist in row?
		if puzzle.values[row][i] == num {
			return false
		}

		// does num already exist in column?
		if puzzle.values[i][col] == num {
			return false
		}
	}

	// check 9x9 grid
	topRow := row / 3 * 3
	topCol := col / 3 * 3

	for i := topRow; i < topRow+3; i++ {
		for j := topCol; j < topCol+3; j++ {
			if !(i == row && j == col) {
				if puzzle.values[i][j] == num {
					return false
				}
			}
		}
	}

	return true
}

func solve(puzzle *Puzzle, row byte, col byte) bool {

	// past last row?
	if row == 9 {
		return true
	}

	// puzzle cell already has a value, move to next cell
	if puzzle.values[row][col] != 0 {
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
				puzzle.values[row][col] = num
				if col == 8 {
					if solve(puzzle, row+1, 0) {
						return true
					}
				} else {
					if solve(puzzle, row, col+1) {
						return true
					}
				}
				puzzle.values[row][col] = 0
			}
		}
	}

	return false
}

func main() {

	var puzzle Puzzle

	var file *os.File
	var err error

	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help {
		_, exec := filepath.Split(os.Args[0])
		fmt.Println("usage: ", exec, "filename")
		os.Exit(1)
	}

	switch len(flag.Args()) {

	case 0:
		file = os.Stdin

	case 1:
		filename := flag.Arg(0)
		file, err = os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

	default:
		_, exec := filepath.Split(os.Args[0])
		fmt.Println("usage: ", exec, "filename")
		os.Exit(1)
	}

	reader := bufio.NewReader(file)

	for row := 0; row < 9; row++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		for col, rune := range line[:9] {
			puzzle.values[row][col] = byte(rune - '0')
		}
	}

	fmt.Println("Puzzle to solve:")
	fmt.Println(puzzle)
	fmt.Println()

	if solve(&puzzle, 0, 0) {
		fmt.Println("Puzzle solved:")
		fmt.Println(puzzle)
	} else {
		fmt.Println("Cannot solve this puzzle")
	}

}
