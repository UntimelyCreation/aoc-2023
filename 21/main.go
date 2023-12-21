package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/grid"
	"github.com/UntimelyCreation/aoc-2023-go/pkg/utils"
	"github.com/emirpasic/gods/sets/hashset"
)

var directions []grid.Direction = []grid.Direction{grid.Up, grid.Right, grid.Down, grid.Left}

type QueueEntry struct {
	position grid.Position
	steps    int
}

func bfs(tiles grid.Grid[rune], maxSteps int) int {
	queue := hashset.New()
	for k := range tiles {
		if tiles[k] == 'S' {
			queue.Add(k)
			break
		}
	}

	for i := 1; i <= maxSteps; i++ {
		newQueue := hashset.New()

		for _, el := range queue.Values() {
			qe := el.(grid.Position)

			if _, exists := tiles[qe]; !exists {
				continue
			}

			for _, dir := range directions {
				newPos := qe.Move(dir)
				if tiles[newPos] != '#' {
					newQueue.Add(newPos)
				}
			}
		}

		queue = newQueue
	}

	return queue.Size()
}

func bfsInfinite(tiles grid.Grid[rune], maxSteps int) []int {
	rows, cols := tiles.Dimensions()

	base := utils.Mod(maxSteps, rows)
	values := []int{}

	queue := hashset.New()
	for k := range tiles {
		if tiles[k] == 'S' {
			queue.Add(k)
			break
		}
	}

	i := 1
	for len(values) < 3 {
		newQueue := hashset.New()

		for _, el := range queue.Values() {
			qe := el.(grid.Position)

			for _, dir := range directions {
				newPos := qe.Move(dir)
				newGridPos := grid.Position{Row: utils.Mod(newPos.Row, rows), Col: utils.Mod(newPos.Col, cols)}
				if tiles[newGridPos] != '#' {
					newQueue.Add(newPos)
				}
			}
		}

		queue = newQueue
		if utils.Mod(i, rows) == base {
			values = append(values, queue.Size())
		}

		i++
	}

	return values
}

// It is assumed here that x0, x1, x2 := 0, 1, 2
func applyNewtonPolynomial(values [3]int, n int) int {
	y0, y1, y2 := values[0], values[1], values[2]

	d0, d1, d2 := y0, y1-y0, y2-y1

	return d0 + d1*n + (d2-d1)*n*(n-1)/2
}

func getGardenPlotCounts(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	tilesRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	tiles := grid.Grid[rune]{}
	tilesInfinite := grid.Grid[rune]{}
	for row, line := range tilesRaw {
		for col, r := range line {
			tiles[grid.Position{Row: row, Col: col}] = r
			tilesInfinite[grid.Position{Row: row, Col: col}] = r
		}
	}

	rows, _ := tiles.Dimensions()

	gardenPlotsCount := bfs(tiles, 64)

	interpolationValues := bfsInfinite(tilesInfinite, 26501365)
	gardenPlotsCountInfinite := applyNewtonPolynomial([3]int{interpolationValues[0], interpolationValues[1], interpolationValues[2]}, 26501365/rows)

	return gardenPlotsCount, gardenPlotsCountInfinite
}

func main() {
	gardentPlotsCount, gardenPlotsCountInfinite := getGardenPlotCounts("21/input.txt")
	fmt.Print("Part 1 solution: ", gardentPlotsCount, "\nPart 2 solution: ", gardenPlotsCountInfinite, "\n")
}
