package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/utils"
)

func calcMinSteps1(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	directionsRaw := strings.Split(strings.Trim(string(file), "\n"), "\n\n")
	instructions := directionsRaw[0]

	nodesRaw := strings.Split(directionsRaw[1], "\n")
	nodes := map[string][]string{}

	pattern := `([A-Z0-9]{3}) = \(([A-Z0-9]{3}), ([A-Z0-9]{3})\)`
	regex := regexp.MustCompile(pattern)
	for _, node := range nodesRaw {
		matches := regex.FindAllStringSubmatch(node, -1)
		nodes[matches[0][1]] = []string{matches[0][2], matches[0][3]}
	}

	minSteps := 0
	curr := "AAA"
	for curr != "ZZZ" {
		for _, ch := range instructions {
			switch ch {
			case 'L':
				curr = nodes[curr][0]
			case 'R':
				curr = nodes[curr][1]
			}
			minSteps += 1
		}
	}

	return minSteps
}

func calcMinSteps2(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	directionsRaw := strings.Split(strings.Trim(string(file), "\n"), "\n\n")
	instructions := directionsRaw[0]

	nodesRaw := strings.Split(directionsRaw[1], "\n")
	nodes := map[string][]string{}

	pattern := `([A-Z0-9]{3}) = \(([A-Z0-9]{3}), ([A-Z0-9]{3})\)`
	regex := regexp.MustCompile(pattern)
	for _, node := range nodesRaw {
		matches := regex.FindAllStringSubmatch(node, -1)
		nodes[matches[0][1]] = []string{matches[0][2], matches[0][3]}
	}

	startSteps := []int{}
	for node := range nodes {
		if node[2] == 'A' {
			curr := node
			minSteps := 0
			for curr[2] != 'Z' {
				for _, ch := range instructions {
					switch ch {
					case 'L':
						curr = nodes[curr][0]
					case 'R':
						curr = nodes[curr][1]
					}
					minSteps += 1
				}
			}
			startSteps = append(startSteps, minSteps)
		}
	}

	minSteps := 1
	for _, count := range startSteps {
		minSteps = utils.Lcm(minSteps, count)
	}

	return minSteps
}

func main() {
	minSteps1 := calcMinSteps1("08/input.txt")
	minSteps2 := calcMinSteps2("08/input.txt")
	fmt.Print("Part 1 solution: ", minSteps1, "\nPart 2 solution: ", minSteps2, "\n")
}
