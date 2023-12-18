package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/grid"
)

func getExpandedRows(image grid.Grid[rune]) []int {
	indices := []int{}

	xMin, xMax := image.XRange()
	yMin, yMax := image.YRange()

	for i := xMin; i <= xMax; i++ {
		isEmpty := true
		for j := yMin; j <= yMax; j++ {
			if image[grid.Position{Row: i, Col: j}] == '#' {
				isEmpty = false
			}
		}
		if isEmpty {
			indices = append(indices, i)
		}
	}

	return indices
}

func processGalaxyImage(path string, expansionFactor int) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	imageRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	image := grid.Grid[rune]{}
	for row, line := range imageRaw {
		for col, r := range line {
			image[grid.Position{Row: row, Col: col}] = r
		}
	}

	expandedRows := getExpandedRows(image)
	expandedCols := getExpandedRows(image.Transpose())

	galaxies := [][]int{}
	xMin, xMax := image.XRange()
	yMin, yMax := image.YRange()

	for i := xMin; i <= xMax; i++ {
		for j := yMin; j <= yMax; j++ {
			if image[grid.Position{Row: i, Col: j}] == '#' {
				k, x := xMin, xMin
				for k < i {
					if slices.Contains(expandedRows, k) {
						x += expansionFactor
					} else {
						x++
					}
					k++
				}
				l, y := yMin, yMin
				for l < j {
					if slices.Contains(expandedCols, l) {
						y += expansionFactor
					} else {
						y++
					}
					l++
				}
				galaxies = append(galaxies, []int{x, y})
			}
		}
	}

	shortestPathsSum := 0

	for i := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			a, b := galaxies[i][0], galaxies[i][1]
			c, d := galaxies[j][0], galaxies[j][1]
			shortestPathsSum += int(math.Abs(float64(c-a)) + math.Abs(float64(d-b)))
		}
	}

	return shortestPathsSum
}

func main() {
	shortestPathsSum1 := processGalaxyImage("11/input.txt", 2)
	shortestPathsSum2 := processGalaxyImage("11/input.txt", 1000000)
	fmt.Print("Part 1 solution: ", shortestPathsSum1, "\nPart 2 solution: ", shortestPathsSum2, "\n")
}
