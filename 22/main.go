package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/space"
)

type Brick struct {
	start space.Position
	end   space.Position
}

func (b Brick) overlaps(o Brick) bool {
	xOverlaps := (b.start.X <= o.end.X) && (o.start.X <= b.end.X)
	yOverlaps := (b.start.Y <= o.end.Y) && (o.start.Y <= b.end.Y)

	return xOverlaps && yOverlaps
}

func settleBrickStack(bricks []Brick) ([]Brick, int) {
	settledBricks := make([]Brick, len(bricks))
	copy(settledBricks, bricks)

	slices.SortFunc(settledBricks, func(a, b Brick) int {
		return a.start.Z - b.start.Z
	})

	movedBricks := 0

	for i := range settledBricks {
		moved := false
	fall:
		for settledBricks[i].start.Z > 1 {
			for _, other := range settledBricks[:i] {
				if settledBricks[i].overlaps(other) && (settledBricks[i].start.Z <= other.end.Z+1) {
					break fall
				}
			}
			settledBricks[i].start.Z--
			settledBricks[i].end.Z--
			moved = true
		}

		if moved {
			movedBricks++
		}
	}

	return settledBricks, movedBricks
}

func simulateFallingBricks(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	bricksRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	bricks := []Brick{}
	for _, line := range bricksRaw {
		split := strings.Split(line, "~")
		start := strings.Split(split[0], ",")
		end := strings.Split(split[1], ",")

		sX, _ := strconv.Atoi(start[0])
		sY, _ := strconv.Atoi(start[1])
		sZ, _ := strconv.Atoi(start[2])
		startPos := space.Position{X: sX, Y: sY, Z: sZ}

		eX, _ := strconv.Atoi(end[0])
		eY, _ := strconv.Atoi(end[1])
		eZ, _ := strconv.Atoi(end[2])
		endPos := space.Position{X: eX, Y: eY, Z: eZ}

		bricks = append(bricks, Brick{startPos, endPos})
	}

	safeBricks := 0
	movedBricksCount := 0

	settledBricks, _ := settleBrickStack(bricks)

	for i := range settledBricks {
		copiedBricks := make([]Brick, len(settledBricks))
		copy(copiedBricks, settledBricks)

		otherBricks := append(copiedBricks[:i], copiedBricks[i+1:]...)
		_, movedBricks := settleBrickStack(otherBricks)

		if movedBricks == 0 {
			safeBricks++
		}
		movedBricksCount += movedBricks
	}

	return safeBricks, movedBricksCount
}

func main() {
	safeBricks, movedBricks := simulateFallingBricks("22/input.txt")
	fmt.Print("Part 1 solution: ", safeBricks, "\nPart 2 solution: ", movedBricks, "\n")
}
