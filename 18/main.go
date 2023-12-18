package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/grid"
)

var letterToDir map[string]grid.Direction = map[string]grid.Direction{
	"U": grid.Up,
	"R": grid.Right,
	"D": grid.Down,
	"L": grid.Left,
}

type Instruction struct {
	direction grid.Direction
	distance  int
}

func getLagoonCapacity(digInstructions []Instruction) int {
	inside := 0
	boundary := 0
	vertices := []grid.Position{}

	start := grid.Position{Row: 0, Col: 0}
	curr := start
	vertices = append(vertices, curr)

	for _, inst := range digInstructions {
		dir, dist := inst.direction, inst.distance

		for i := 0; i < dist; i++ {
			boundary++
			next := curr.Move(dir)
			curr = next
		}
		vertices = append(vertices, curr)

	}

	// Shoelace formula
	for i := 0; i < len(vertices)-1; i++ {
		a, b := vertices[i], vertices[i+1]
		inside += (a.Col * b.Row) - (a.Row * b.Col)
	}

	// Pick's theorem with a +2 to compensate for both divisions
	return boundary/2 + inside/2 + 1
}

func processDigInstructions(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	digPlanRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	digInstructions := []Instruction{}
	digInstructionsHex := []Instruction{}
	for _, line := range digPlanRaw {
		elements := strings.Fields(line)

		direction := letterToDir[elements[0]]
		distance, _ := strconv.Atoi(elements[1])
		digInstructions = append(digInstructions, Instruction{direction, distance})

		hexDist, _ := strconv.ParseInt(elements[2][2:7], 16, 64)
		hexDir, _ := strconv.Atoi(string(elements[2][7]))
		digInstructionsHex = append(digInstructionsHex, Instruction{grid.Direction(hexDir), int(hexDist)})
	}

	lagoonCapacity := getLagoonCapacity(digInstructions)
	lagoonCapacityHex := getLagoonCapacity(digInstructionsHex)

	return lagoonCapacity, lagoonCapacityHex
}

func main() {
	lagoonCapacity, lagoonCapacityHex := processDigInstructions("18/input.txt")
	fmt.Print("Part 1 solution: ", lagoonCapacity, "\nPart 2 solution: ", lagoonCapacityHex, "\n")
}
