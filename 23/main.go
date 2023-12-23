package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/grid"
)

type Graph map[grid.Position][]Edge

type Edge struct {
	position grid.Position
	distance int
}

func findNonJunctionNode(graph Graph) (position grid.Position, ok bool) {
	for pos, edges := range graph {
		if len(edges) == 2 {
			position = pos
			ok = true
			break
		}
	}
	return position, ok
}

func getTrailsGraph(trails grid.Grid[rune], slippery bool) Graph {
	graph := Graph{}

	for pos := range trails {
		if trails[pos] == '#' {
			continue
		}

		edges := []Edge{}

		neighbors := getNeighbors(trails, pos, slippery)

		for _, nb := range neighbors {
			if tile, exists := trails[nb]; exists && tile != '#' {
				edges = append(edges, Edge{position: nb, distance: 1})
			}
		}

		graph[pos] = edges
	}

	pos, exists := findNonJunctionNode(graph)
	for exists {
		neighbor1, dist1 := graph[pos][0].position, graph[pos][0].distance
		neighbor2, dist2 := graph[pos][1].position, graph[pos][1].distance

		delete(graph, pos)

		for i, edge := range graph[neighbor1] {
			if edge.position == pos {
				graph[neighbor1][i] = Edge{position: neighbor2, distance: dist1 + dist2}
			}
		}
		for i, edge := range graph[neighbor2] {
			if edge.position == pos {
				graph[neighbor2][i] = Edge{position: neighbor1, distance: dist1 + dist2}
			}
		}

		pos, exists = findNonJunctionNode(graph)
	}

	return graph
}

func getNeighbors(trails grid.Grid[rune], position grid.Position, slippery bool) []grid.Position {
	neighbors := []grid.Position{}

	switch trails[position] {
	case '.':
		neighbors = append(neighbors, position.Move(grid.Up))
		neighbors = append(neighbors, position.Move(grid.Right))
		neighbors = append(neighbors, position.Move(grid.Down))
		neighbors = append(neighbors, position.Move(grid.Left))
	case '>':
		neighbors = append(neighbors, position.Move(grid.Right))
		if !slippery {
			neighbors = append(neighbors, position.Move(grid.Up))
			neighbors = append(neighbors, position.Move(grid.Down))
			neighbors = append(neighbors, position.Move(grid.Left))
		}
	case 'v':
		neighbors = append(neighbors, position.Move(grid.Down))
		if !slippery {
			neighbors = append(neighbors, position.Move(grid.Up))
			neighbors = append(neighbors, position.Move(grid.Right))
			neighbors = append(neighbors, position.Move(grid.Left))
		}
	}

	return neighbors
}

func dfs(graph Graph, visited *grid.Grid[bool], position grid.Position, target grid.Position, slippery bool) int {
	if position == target {
		return 0
	}

	maxSteps := 0

	for _, edge := range graph[position] {
		neighbor, dist := edge.position, edge.distance
		if (*visited)[neighbor] {
			continue
		}
		(*visited)[neighbor] = true

		d := dfs(graph, visited, neighbor, target, slippery)

		maxSteps = max(maxSteps, dist+d)

		(*visited)[neighbor] = false
	}

	return maxSteps
}

func getLongestWalks(path string, slippery bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	trailsRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	trails := grid.Grid[rune]{}
	visited := grid.Grid[bool]{}
	for row, line := range trailsRaw {
		for col, r := range line {
			position := grid.Position{Row: row, Col: col}
			trails[position] = r
			visited[position] = false
		}
	}

	graph := getTrailsGraph(trails, slippery)

	rows, cols := trails.Dimensions()
	start := grid.Position{Row: 0, Col: 1}
	end := grid.Position{Row: rows - 1, Col: cols - 2}

	maxSteps := dfs(graph, &visited, start, end, slippery)

	return maxSteps
}

func main() {
	maxStepsSlippery := getLongestWalks("23/input.txt", true)
	maxStepsNotSlippery := getLongestWalks("23/input.txt", false)
	fmt.Print("Part 1 solution: ", maxStepsSlippery, "\nPart 2 solution: ", maxStepsNotSlippery, "\n")
}
