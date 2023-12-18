package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func processGameDraws(path string) (int, int) {
	pattern := `Game (\d+): (.*)`
	regex := regexp.MustCompile(pattern)

	maximums := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	possibleIdsSum := 0
	gamePossible := map[int]bool{}

	powerSum := 0
	gameMinimums := map[int][]int{}

	for scanner.Scan() {
		line := scanner.Text()

		matches := regex.FindAllStringSubmatch(line, -1)
		id, _ := strconv.Atoi(matches[0][1])

		gamePossible[id] = true
		gameMinimums[id] = []int{0, 0, 0}

		draws := strings.Split(matches[0][2], "; ")
		for i := range draws {
			colors := strings.Split(draws[i], ", ")

			for j := range colors {
				split := strings.Split(colors[j], " ")
				count, _ := strconv.Atoi(split[0])
				color := split[1]
				colorId := 0

				if count > maximums[color] {
					gamePossible[id] = false
				}

				switch color {
				case "red":
					colorId = 0
				case "blue":
					colorId = 1
				case "green":
					colorId = 2
				}
				gameMinimums[id][colorId] = max(gameMinimums[id][colorId], count)
			}
		}

	}

	for k, v := range gamePossible {
		if v {
			possibleIdsSum += k
		}
	}
	for _, v := range gameMinimums {
		powerSum += v[0] * v[1] * v[2]
	}

	return possibleIdsSum, powerSum
}

func main() {
	possibleIdsSum, powerSum := processGameDraws("02/input.txt")
	fmt.Print("Part 1 solution: ", possibleIdsSum, "\nPart 2 solution: ", powerSum, "\n")
}
