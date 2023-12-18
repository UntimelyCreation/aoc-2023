package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
)

func transpose(pattern []string) []string {
	cols := len(pattern[0])
	transposed := make([]string, cols)
	for i := 0; i < cols; i++ {
		for j := range pattern {
			transposed[i] += string(pattern[len(pattern)-1-j][i])
		}
	}
	return transposed
}

func hammingDistance(a []string, b []string) int {
	distance := 0

	for i := range a {
		for j := range a[0] {
			if a[i][j] != b[i][j] {
				distance++
			}
		}
	}

	return distance
}

func getReflectionLine(pattern []string, smudged bool) int {
	for r := 1; r < len(pattern); r++ {
		copied := make([]string, len(pattern))
		copy(copied, pattern)

		left := copied[:r]
		right := copied[r:]
		min := min(len(left), len(right))

		left = left[len(left)-min:]
		right = right[:min]
		slices.Reverse(right)

		if (!smudged && reflect.DeepEqual(left, right)) || (smudged && hammingDistance(left, right) == 1) {
			return r
		}
	}
	return -1
}

func processGroundPattern(path string, smudged bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	patternsRaw := strings.Split(strings.Trim(string(file), "\n"), "\n\n")
	patterns := [][]string{}
	for _, raw := range patternsRaw {
		pattern := strings.Split(raw, "\n")
		patterns = append(patterns, pattern)
	}

	summary := 0
	for _, pattern := range patterns {
		rows := getReflectionLine(pattern, smudged)
		cols := getReflectionLine(transpose(pattern), smudged)
		if rows != -1 {
			summary += 100 * rows
		}
		if cols != -1 {
			summary += cols
		}
	}

	return summary
}

func main() {
	summaryClean := processGroundPattern("13/input.txt", false)
	summarySmudged := processGroundPattern("13/input.txt", true)
	fmt.Print("Part 1 solution: ", summaryClean, "\nPart 2 solution: ", summarySmudged, "\n")
}
