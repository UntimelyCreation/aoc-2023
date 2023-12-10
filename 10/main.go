package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Queue [][]int

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) Push(slice []int) {
	*q = append(*q, slice)
}

func (q *Queue) Pop() []int {
	last := (*q)[0]
	*q = (*q)[1:]

	return last
}

var directions map[rune][][]int = map[rune][][]int{
	'|': {{1, 0}, {-1, 0}},
	'-': {{0, 1}, {0, -1}},
	'F': {{0, 1}, {1, 0}},
	'7': {{0, -1}, {1, 0}},
	'J': {{0, -1}, {-1, 0}},
	'L': {{0, 1}, {-1, 0}},
}

var verticals []rune = []rune{'|', 'F', 'J', '7', 'L'}

func process_pipe_map(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	pipes_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	pipes := [][]rune{}
	for _, line := range pipes_raw {
		pipes = append(pipes, []rune(line))
	}
	start := []int{}
start:
	for i := range pipes {
		for j := range pipes[i] {
			if pipes[i][j] == 'S' {
				start = []int{i, j}
				// HARDCODED: Replace S with corresponding pipe symbol in your personal input
				pipes[i][j] = 'J'
				// pipes[i][j] = '7'
				break start
			}
		}
	}

	visited := [][]bool{}
	for range pipes {
		temp := []bool{}
		for range pipes[0] {
			temp = append(temp, false)
		}
		visited = append(visited, temp)
	}
	queue := Queue{}
	queue.Push([]int{start[0], start[1], 0})
	visited[start[0]][start[1]] = true
	furthest_dist := 0

	for !queue.IsEmpty() {
		curr := queue.Pop()

		i, j, dist := curr[0], curr[1], curr[2]
		furthest_dist = max(furthest_dist, dist)

		next_dirs := directions[pipes[i][j]]

		for _, coords := range next_dirs {
			x, y := i+coords[0], j+coords[1]
			if !visited[x][y] {
				visited[x][y] = true
				queue.Push([]int{x, y, dist + 1})
			}
		}
	}

	for i := range pipes {
		for j := range pipes[i] {
			if !visited[i][j] {
				pipes[i][j] = '.'
			}
		}
	}

	enclosed := 0
	for _, row := range pipes {
		inside := 0

		i := 0
		for i < len(row) {
			r := row[i]

			switch r {
			case '|':
				inside++
			case 'J':
				inside++
			case '7':
				inside++
			case 'F':
				i++
				for row[i] == '-' && row[i] != 'J' && row[i] != '7' {
					i++
				}
				if row[i] == 'J' {
					inside++
				}
			case 'L':
				i++
				for row[i] == '-' && row[i] != 'J' && row[i] != '7' {
					i++
				}
				if row[i] == '7' {
					inside++
				}
			case '.':
				if inside%2 == 1 {
					enclosed++
				}
			}

			i++
		}
	}

	return furthest_dist, enclosed
}

func main() {
	furthest_dist, enclosed := process_pipe_map("10/input.txt")
	fmt.Print("Part 1 solution: ", furthest_dist, "\nPart 2 solution: ", enclosed, "\n")
}
