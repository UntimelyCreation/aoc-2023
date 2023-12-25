package main

import (
	"fmt"
	"log"
	rand "math/rand"
	"os"
	"slices"
	"strings"
)

type Edge struct {
	start string
	end   string
}

func countCuts(edges []Edge, subsets [][]string) int {
	cuts := 0

	for i := 0; i < len(edges); i++ {
		first := 0
		for first < len(subsets) && !slices.Contains(subsets[first], edges[i].start) {
			first++
		}
		second := 0
		for second < len(subsets) && !slices.Contains(subsets[second], edges[i].end) {
			second++
		}

		if first != second {
			cuts++
		}
	}

	return cuts
}

func getGroupSizes(vertices []string, edges []Edge) []int {
	subsets := [][]string{}

	for countCuts(edges, subsets) != 3 {
		subsets = [][]string{}

		for _, vertex := range vertices {
			subsets = append(subsets, []string{vertex})
		}

		for len(subsets) > 2 {
			i := rand.Int() % len(edges)

			first := 0
			for first < len(subsets) && !slices.Contains(subsets[first], edges[i].start) {
				first++
			}
			second := 0
			for second < len(subsets) && !slices.Contains(subsets[second], edges[i].end) {
				second++
			}

			if first == second {
				continue
			}

			subsets[first] = append(subsets[first], subsets[second]...)
			subsets = append(subsets[:second], subsets[second+1:]...)
		}
	}

	subsetSizes := []int{}
	for _, s := range subsets {
		subsetSizes = append(subsetSizes, len(s))
	}

	return subsetSizes
}

func divideConnectedComponents(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	componentsRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	vertices := []string{}
	edges := []Edge{}

	for _, line := range componentsRaw {
		split := strings.Split(line, ": ")

		start := split[0]
		ends := strings.Fields(split[1])

		if !slices.Contains(vertices, start) {
			vertices = append(vertices, start)
		}

		for _, end := range ends {
			if !slices.Contains(vertices, end) {
				vertices = append(vertices, end)
			}
			if !slices.Contains(edges, Edge{start: start, end: end}) &&
				!slices.Contains(edges, Edge{start: end, end: start}) {
				edges = append(edges, Edge{start: start, end: end})
			}
		}
	}

	subsetSizes := getGroupSizes(vertices, edges)
	subsetSizeProduct := 1
	for _, size := range subsetSizes {
		subsetSizeProduct *= size
	}

	return subsetSizeProduct
}

func main() {
	subsetSizeProduct := divideConnectedComponents("25/input.txt")
	fmt.Print("Part 1 solution: ", subsetSizeProduct, "\n")
}
