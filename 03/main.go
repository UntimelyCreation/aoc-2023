package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"unicode"
)

type Symbol struct {
	col       int
	neighbors []int
}

func analyzeEngineSchematic(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	schematic := [][]rune{}
	symbols := [][]*Symbol{}
	partNumbersSum := 0
	gearRatiosSum := 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := []rune(scanner.Text())

		schematic = append(schematic, line)

		temp := []*Symbol{}
		for i := range line {
			if !unicode.IsDigit(line[i]) && !unicode.IsLetter(line[i]) && line[i] != '.' {
				temp = append(temp, &Symbol{
					col:       i,
					neighbors: []int{},
				})
			}
		}
		symbols = append(symbols, temp)
	}

	idxs := [][]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	n := len(schematic)

	for i := range schematic {
		j := 0
		for j < n {
			if unicode.IsDigit(schematic[i][j]) {
				start := j
				for j < n && unicode.IsDigit(schematic[i][j]) {
					j += 1
				}
				end := j

				partNumber, _ := strconv.Atoi(string(schematic[i][start:end]))
				numberIsValid := false

				for k := start; k < end; k++ {
					for _, coords := range idxs {

						l, m := i+coords[0], k+coords[1]
						if l >= 0 && l < n && m >= 0 && m < n {
							for _, sym := range symbols[l] {
								y := sym.col
								if y == m {
									numberIsValid = true
									if !slices.Contains(sym.neighbors, partNumber) {
										sym.neighbors = append(sym.neighbors, partNumber)
									}
								}
							}
						}
					}
				}

				if numberIsValid {
					partNumbersSum += partNumber
				}
			}
			j += 1
		}
	}

	for i := range symbols {
		for _, sym := range symbols[i] {
			if schematic[i][sym.col] == '*' && len(sym.neighbors) == 2 {
				ratio := 1
				for _, num := range sym.neighbors {
					ratio *= num
				}
				gearRatiosSum += ratio
			}
		}
	}

	return partNumbersSum, gearRatiosSum
}

func main() {
	enginePartsSum, gearRatiosSum := analyzeEngineSchematic("03/input.txt")
	fmt.Print("Part 1 solution: ", enginePartsSum, "\nPart 2 solution: ", gearRatiosSum, "\n")
}
