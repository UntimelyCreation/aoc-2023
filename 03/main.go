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

func analyze_engine_schematic(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	schematic := [][]rune{}
	symbols := [][]*Symbol{}
	part_numbers_sum := 0
	gear_ratios_sum := 0

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

				part_number, _ := strconv.Atoi(string(schematic[i][start:end]))
				number_is_valid := false

				for k := start; k < end; k++ {
					for _, coords := range idxs {

						l, m := i+coords[0], k+coords[1]
						if l >= 0 && l < n && m >= 0 && m < n {
							for _, sym := range symbols[l] {
								y := sym.col
								if y == m {
									number_is_valid = true
									if !slices.Contains(sym.neighbors, part_number) {
										sym.neighbors = append(sym.neighbors, part_number)
									}
								}
							}
						}
					}
				}

				if number_is_valid {
					part_numbers_sum += part_number
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
				gear_ratios_sum += ratio
			}
		}
	}

	return part_numbers_sum, gear_ratios_sum
}

func main() {
	engine_parts_sum, gear_ratios_sum := analyze_engine_schematic("03/input.txt")
	fmt.Print("Part 1 solution: ", engine_parts_sum, "\nPart 2 solution: ", gear_ratios_sum, "\n")
}
