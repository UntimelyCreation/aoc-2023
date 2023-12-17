package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/utils"
	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

type QueueEntry struct {
	position  utils.Position
	direction utils.Direction
	streak    int
	heatLoss  int
}

type CacheEntry struct {
	position  utils.Position
	direction utils.Direction
	streak    int
}

func compareHeatLoss(a, b any) int {
	return a.(QueueEntry).heatLoss - b.(QueueEntry).heatLoss
}

var dirToDirs map[int][]int = map[int][]int{
	0: {0, 1, 3},
	1: {1, 2, 0},
	2: {2, 3, 1},
	3: {3, 0, 2},
}

var dirToCoords map[int][]int = map[int][]int{
	0: {-1, 0},
	1: {0, 1},
	2: {1, 0},
	3: {0, -1},
}

func dijkstra(blocks utils.Grid[int], start utils.Position, target utils.Position, minStreak, maxStreak int) int {
	queue := pq.NewWith(compareHeatLoss)
	queue.Enqueue(QueueEntry{
		position:  start,
		direction: utils.Right,
		streak:    1,
		heatLoss:  0,
	})

	cache := make(map[CacheEntry]int)

	for !queue.Empty() {
		next, _ := queue.Dequeue()
		qe := next.(QueueEntry)

		if _, exists := blocks[qe.position]; !exists {
			continue
		}

		heatLoss := qe.heatLoss + *blocks[qe.position]

		// TODO: Algorithm works for real input but not for second test input
		// Should add check for qe.streak >= minStreak as well
		if qe.position == target {
			return heatLoss
		}

		ce := CacheEntry{
			position:  qe.position,
			direction: qe.direction,
			streak:    qe.streak,
		}
		if minHeatLoss, cached := cache[ce]; cached {
			if minHeatLoss <= heatLoss {
				continue
			}
		}
		cache[ce] = heatLoss

		if qe.streak >= minStreak {
			left := utils.TurnLeft(qe.direction)
			queue.Enqueue(QueueEntry{
				position:  qe.position.Move(left),
				direction: left,
				streak:    1,
				heatLoss:  heatLoss,
			})

			right := utils.TurnRight(qe.direction)
			queue.Enqueue(QueueEntry{
				position:  qe.position.Move(right),
				direction: right,
				streak:    1,
				heatLoss:  heatLoss,
			})
		}

		if qe.streak < maxStreak {
			queue.Enqueue(QueueEntry{
				position:  qe.position.Move(qe.direction),
				direction: qe.direction,
				streak:    qe.streak + 1,
				heatLoss:  heatLoss,
			})
		}
	}
	return -1
}

func navigateCrucible(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	blocksRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	blocks := utils.Grid[int]{}
	for row, line := range blocksRaw {
		for col, r := range line {
			val := int(r - '0')
			blocks[utils.Position{Row: row, Col: col}] = &val
		}
	}
	rows, cols := blocks.Dimensions()

	minHeatLoss1 := dijkstra(blocks, utils.Position{Row: 0, Col: 1}, utils.Position{Row: rows - 1, Col: cols - 1}, 0, 3)
	minHeatLoss2 := dijkstra(blocks, utils.Position{Row: 0, Col: 1}, utils.Position{Row: rows - 1, Col: cols - 1}, 4, 10)

	return minHeatLoss1, minHeatLoss2
}

func main() {
	minHeatLoss1, minHeatLoss2 := navigateCrucible("17/input.txt")
	fmt.Print("Part 1 solution: ", minHeatLoss1, "\nPart 2 solution: ", minHeatLoss2, "\n")
}
