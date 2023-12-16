package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var dirs map[int][]int = map[int][]int{
	0: {-1, 0},
	1: {0, 1},
	2: {1, 0},
	3: {0, -1},
}

func simulate_beam(layout *[][]rune, mark *[][][]int, i int, j int, dir int) {
	rows, cols := len(*layout), len((*layout)[0])

	if slices.Contains((*mark)[i][j], dir) {
		return
	}
	(*mark)[i][j] = append((*mark)[i][j], dir)

	new_dirs := []int{}

	switch (*layout)[i][j] {
	case '/':
		switch dir {
		case 0:
			new_dirs = append(new_dirs, 1)
		case 1:
			new_dirs = append(new_dirs, 0)
		case 2:
			new_dirs = append(new_dirs, 3)
		case 3:
			new_dirs = append(new_dirs, 2)
		}
	case '\\':
		switch dir {
		case 0:
			new_dirs = append(new_dirs, 3)
		case 1:
			new_dirs = append(new_dirs, 2)
		case 2:
			new_dirs = append(new_dirs, 1)
		case 3:
			new_dirs = append(new_dirs, 0)
		}
	case '-':
		switch dir {
		case 0:
			new_dirs = append(new_dirs, 1)
			new_dirs = append(new_dirs, 3)
		case 1:
			new_dirs = append(new_dirs, 1)
		case 2:
			new_dirs = append(new_dirs, 1)
			new_dirs = append(new_dirs, 3)
		case 3:
			new_dirs = append(new_dirs, 3)
		}
	case '|':
		switch dir {
		case 0:
			new_dirs = append(new_dirs, 0)
		case 1:
			new_dirs = append(new_dirs, 0)
			new_dirs = append(new_dirs, 2)
		case 2:
			new_dirs = append(new_dirs, 2)
		case 3:
			new_dirs = append(new_dirs, 0)
			new_dirs = append(new_dirs, 2)
		}
	case '.':
		new_dirs = append(new_dirs, dir)
	}

	for _, d := range new_dirs {
		x, y := dirs[d][0], dirs[d][1]

		if (i+x) >= 0 && (i+x) < rows && (j+y) >= 0 && (j+y) < cols {
			simulate_beam(layout, mark, i+x, j+y, d)
		}
	}
}

func calculate_energized_tiles(layout [][]rune, i int, j int, dir int) int {
	mark := [][][]int{}
	for _, row := range layout {
		temp := [][]int{}
		for range row {
			temp = append(temp, []int{})
		}
		mark = append(mark, temp)
	}
	simulate_beam(&layout, &mark, i, j, dir)

	energized := 0
	for _, row := range mark {
		for _, dirs := range row {
			if len(dirs) > 0 {
				energized++
			}
		}
	}
	return energized
}

func simulate_beams_and_mirrors(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	layout_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	layout := [][]rune{}
	for _, row := range layout_raw {
		layout = append(layout, []rune(row))
	}

	energized := calculate_energized_tiles(layout, 0, 0, 1)

	rows, cols := len(layout), len(layout[0])
	max_energized := 0
	for i := 0; i < rows; i++ {
		max_energized = max(max_energized, calculate_energized_tiles(layout, i, 0, 1))
		max_energized = max(max_energized, calculate_energized_tiles(layout, i, cols-1, 3))
	}
	for j := 0; j < cols; j++ {
		max_energized = max(max_energized, calculate_energized_tiles(layout, 0, j, 2))
		max_energized = max(max_energized, calculate_energized_tiles(layout, rows-1, j, 0))
	}

	return energized, max_energized
}

func main() {
	energized, max_energized := simulate_beams_and_mirrors("16/input.txt")
	fmt.Print("Part 1 solution: ", energized, "\nPart 2 solution: ", max_energized, "\n")
}
