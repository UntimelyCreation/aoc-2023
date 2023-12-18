package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/grid"
)

func simulateBeam(layout *grid.Grid[rune], cache *map[grid.Position][]grid.Direction, position grid.Position, dir grid.Direction) {
	if _, exists := (*layout)[position]; !exists {
		return
	}

	if slices.Contains((*cache)[position], dir) {
		return
	}
	(*cache)[position] = append((*cache)[position], dir)

	newDirs := []grid.Direction{}

	switch (*layout)[position] {
	case '/':
		switch dir {
		case grid.Up, grid.Down:
			newDirs = append(newDirs, grid.TurnRight(dir))
		case grid.Left, grid.Right:
			newDirs = append(newDirs, grid.TurnLeft(dir))
		}
	case '\\':
		switch dir {
		case grid.Up, grid.Down:
			newDirs = append(newDirs, grid.TurnLeft(dir))
		case grid.Left, grid.Right:
			newDirs = append(newDirs, grid.TurnRight(dir))
		}
	case '-':
		switch dir {
		case grid.Up, grid.Down:
			newDirs = append(newDirs, grid.TurnLeft(dir))
			newDirs = append(newDirs, grid.TurnRight(dir))
		case grid.Left, grid.Right:
			newDirs = append(newDirs, dir)
		}
	case '|':
		switch dir {
		case grid.Up, grid.Down:
			newDirs = append(newDirs, dir)
		case grid.Left, grid.Right:
			newDirs = append(newDirs, grid.TurnLeft(dir))
			newDirs = append(newDirs, grid.TurnRight(dir))
		}
	case '.':
		newDirs = append(newDirs, dir)
	}

	for _, d := range newDirs {
		simulateBeam(layout, cache, position.Move(d), d)
	}
}

func getEnergizedTiles(layout grid.Grid[rune], position grid.Position, dir grid.Direction) int {
	cache := map[grid.Position][]grid.Direction{}
	simulateBeam(&layout, &cache, position, dir)

	return len(cache)
}

func simulateBeamsAndMirrors(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	layoutRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	layout := grid.Grid[rune]{}
	for row, line := range layoutRaw {
		for col, r := range line {
			layout[grid.Position{Row: row, Col: col}] = r
		}
	}

	energized := getEnergizedTiles(layout, grid.Position{Row: 0, Col: 0}, grid.Right)

	rows, cols := layout.Dimensions()
	maxEnergized := 0
	for i := 0; i < rows; i++ {
		maxEnergized = max(maxEnergized, getEnergizedTiles(layout, grid.Position{Row: i, Col: 0}, 1))
		maxEnergized = max(maxEnergized, getEnergizedTiles(layout, grid.Position{Row: i, Col: cols - 1}, 3))
	}
	for j := 0; j < cols; j++ {
		maxEnergized = max(maxEnergized, getEnergizedTiles(layout, grid.Position{Row: 0, Col: j}, 2))
		maxEnergized = max(maxEnergized, getEnergizedTiles(layout, grid.Position{Row: rows - 1, Col: j}, 0))
	}

	return energized, maxEnergized
}

func main() {
	energized, max_energized := simulateBeamsAndMirrors("16/input.txt")
	fmt.Print("Part 1 solution: ", energized, "\nPart 2 solution: ", max_energized, "\n")
}
