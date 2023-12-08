package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func gcd(a int, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a int, b int) int {
	return int(a * b / gcd(a, b))
}

func calc_min_lr_steps_1(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	directions_raw := strings.Split(strings.Trim(string(file), "\n"), "\n\n")
	instructions := directions_raw[0]

	nodes_raw := strings.Split(directions_raw[1], "\n")
	nodes_map := map[string][]string{}

	pattern := `([A-Z0-9]{3}) = \(([A-Z0-9]{3}), ([A-Z0-9]{3})\)`
	regex := regexp.MustCompile(pattern)
	for _, node := range nodes_raw {
		matches := regex.FindAllStringSubmatch(node, -1)
		nodes_map[matches[0][1]] = []string{matches[0][2], matches[0][3]}
	}

	min_steps := 0
	curr := "AAA"
	for curr != "ZZZ" {
		for _, ch := range instructions {
			switch ch {
			case 'L':
				curr = nodes_map[curr][0]
			case 'R':
				curr = nodes_map[curr][1]
			}
			min_steps += 1
		}
	}

	return min_steps
}

func calc_min_lr_steps_2(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	directions_raw := strings.Split(strings.Trim(string(file), "\n"), "\n\n")
	instructions := directions_raw[0]

	nodes_raw := strings.Split(directions_raw[1], "\n")
	nodes_map := map[string][]string{}

	pattern := `([A-Z0-9]{3}) = \(([A-Z0-9]{3}), ([A-Z0-9]{3})\)`
	regex := regexp.MustCompile(pattern)
	for _, node := range nodes_raw {
		matches := regex.FindAllStringSubmatch(node, -1)
		nodes_map[matches[0][1]] = []string{matches[0][2], matches[0][3]}
	}

	start_steps := []int{}
	for node := range nodes_map {
		if node[2] == 'A' {
			curr := node
			min_steps := 0
			for curr[2] != 'Z' {
				for _, ch := range instructions {
					switch ch {
					case 'L':
						curr = nodes_map[curr][0]
					case 'R':
						curr = nodes_map[curr][1]
					}
					min_steps += 1
				}
			}
			start_steps = append(start_steps, min_steps)
		}
	}

	min_steps := 1
	for _, count := range start_steps {
		min_steps = lcm(min_steps, count)
	}

	return min_steps
}

func main() {
	min_steps_1 := calc_min_lr_steps_1("08/input.txt")
	min_steps_2 := calc_min_lr_steps_2("08/input.txt")
	fmt.Print("Part 1 solution: ", min_steps_1, "\nPart 2 solution: ", min_steps_2, "\n")
}
