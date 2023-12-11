package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strings"
)

func get_expanded_rows(image [][]rune) []int {
	indices := []int{}

	for i, row := range image {
		if !slices.Contains(row, '#') {
			indices = append(indices, i)
		}
	}

	return indices
}

func transpose(image [][]rune) [][]rune {
	x, y := len(image[0]), len(image)
	transposed := make([][]rune, x)
	for i := range transposed {
		transposed[i] = make([]rune, y)
	}
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			transposed[i][j] = image[j][i]
		}
	}
	return transposed
}

func process_galaxy_image(path string, expansion_factor int) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	image_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	image := [][]rune{}
	for _, row := range image_raw {
		image = append(image, []rune(row))
	}

	expanded_rows := get_expanded_rows(image)
	expanded_cols := get_expanded_rows(transpose(image))

	galaxies := [][]int{}
	for i, row := range image {
		for j := range row {
			if image[i][j] == '#' {
				k, x := 0, 0
				for k < i {
					if slices.Contains(expanded_rows, k) {
						x += expansion_factor
					} else {
						x++
					}
					k++
				}
				l, y := 0, 0
				for l < j {
					if slices.Contains(expanded_cols, l) {
						y += expansion_factor
					} else {
						y++
					}
					l++
				}
				galaxies = append(galaxies, []int{x, y})
			}
		}
	}

	shortest_paths_sum := 0

	for i := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			a, b := galaxies[i][0], galaxies[i][1]
			c, d := galaxies[j][0], galaxies[j][1]
			shortest_paths_sum += int(math.Abs(float64(c-a)) + math.Abs(float64(d-b)))
		}
	}

	return shortest_paths_sum
}

func main() {
	shortest_paths_sum_1 := process_galaxy_image("11/input.txt", 2)
	shortest_paths_sum_2 := process_galaxy_image("11/input.txt", 1000000)
	fmt.Print("Part 1 solution: ", shortest_paths_sum_1, "\nPart 2 solution: ", shortest_paths_sum_2, "\n")
}
