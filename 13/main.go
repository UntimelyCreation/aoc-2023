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

func hamming_distance(a []string, b []string) int {
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

func get_reflection_line(pattern []string, smudged bool) int {
	for r := 1; r < len(pattern); r++ {
		copied := make([]string, len(pattern))
		copy(copied, pattern)

		left := copied[:r]
		right := copied[r:]
		min := min(len(left), len(right))

		left = left[len(left)-min:]
		right = right[:min]
		slices.Reverse(right)

		if (!smudged && reflect.DeepEqual(left, right)) || (smudged && hamming_distance(left, right) == 1) {
			return r
		}
	}
	return -1
}

func process_ground_pattern(path string, smudged bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	patterns_raw := strings.Split(strings.Trim(string(file), "\n"), "\n\n")
	patterns := [][]string{}
	for _, raw := range patterns_raw {
		pattern := strings.Split(raw, "\n")
		patterns = append(patterns, pattern)
	}

	summary := 0
	for _, pattern := range patterns {
		rows := get_reflection_line(pattern, smudged)
		cols := get_reflection_line(transpose(pattern), smudged)
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
	summary_clean := process_ground_pattern("13/input.txt", false)
	summary_smudged := process_ground_pattern("13/input.txt", true)
	fmt.Print("Part 1 solution: ", summary_clean, "\nPart 2 solution: ", summary_smudged, "\n")
}
