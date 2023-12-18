package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/grid"
	aq "github.com/emirpasic/gods/queues/arrayqueue"
)

type QueueEntry struct {
	position grid.Position
	distance int
}

var pipeDirs map[rune][]grid.Direction = map[rune][]grid.Direction{
	'|': {grid.Up, grid.Down},
	'-': {grid.Left, grid.Right},
	'F': {grid.Right, grid.Down},
	'7': {grid.Left, grid.Down},
	'J': {grid.Left, grid.Up},
	'L': {grid.Right, grid.Up},
}

func processPipeMap(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	pipesRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	pipes := grid.Grid[rune]{}
	for row, line := range pipesRaw {
		for col, r := range line {
			pipes[grid.Position{Row: row, Col: col}] = r
		}
	}

	start := grid.Position{}
	for k, v := range pipes {
		if v == 'S' {
			start = k
			// HARDCODED: Replace S with corresponding pipe symbol in personal input
			pipes[k] = 'J'
			// pipes[k] = '7'
		}
	}

	visited := grid.Grid[bool]{}
	for k := range pipes {
		visited[k] = false
	}
	queue := aq.New()
	queue.Enqueue(QueueEntry{position: start, distance: 0})
	visited[start] = true
	furthestDist := 0

	for !queue.Empty() {
		qe, _ := queue.Dequeue()

		position, dist := qe.(QueueEntry).position, qe.(QueueEntry).distance
		furthestDist = max(furthestDist, dist)

		nextDirs := pipeDirs[pipes[position]]

		for _, dir := range nextDirs {
			newPosition := position.Move(dir)
			if !visited[newPosition] {
				visited[newPosition] = true
				queue.Enqueue(QueueEntry{position: newPosition, distance: dist + 1})
			}
		}
	}

	for k := range pipes {
		if !visited[k] {
			pipes[k] = '.'
		}
	}

	xMin, xMax := pipes.XRange()
	yMin, yMax := pipes.YRange()

	enclosed := 0
	for i := xMin; i <= xMax; i++ {
		inside := 0

		j := yMin
		for j <= yMax {
			r := pipes[grid.Position{Row: i, Col: j}]

			switch r {
			case '|':
				inside++
			case 'J':
				inside++
			case '7':
				inside++
			case 'F':
				j++
				for pipes[grid.Position{Row: i, Col: j}] == '-' && pipes[grid.Position{Row: i, Col: j}] != 'J' && pipes[grid.Position{Row: i, Col: j}] != '7' {
					j++
				}
				if pipes[grid.Position{Row: i, Col: j}] == 'J' {
					inside++
				}
			case 'L':
				j++
				for pipes[grid.Position{Row: i, Col: j}] == '-' && pipes[grid.Position{Row: i, Col: j}] != 'J' && pipes[grid.Position{Row: i, Col: j}] != '7' {
					j++
				}
				if pipes[grid.Position{Row: i, Col: j}] == '7' {
					inside++
				}
			case '.':
				if inside%2 == 1 {
					enclosed++
				}
			}

			j++
		}
	}

	return furthestDist, enclosed
}

func main() {
	furthestDist, enclosed := processPipeMap("10/input.txt")
	fmt.Print("Part 1 solution: ", furthestDist, "\nPart 2 solution: ", enclosed, "\n")
}
