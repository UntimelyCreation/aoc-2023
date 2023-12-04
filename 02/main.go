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

func process_game_draws(path string) (int, int) {
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

	possible_ids_sum := 0
	game_possible := map[int]bool{}

	power_sum := 0
	game_minimums := map[int][]int{}

	for scanner.Scan() {
		line := scanner.Text()

		matches := regex.FindAllStringSubmatch(line, -1)
		id, _ := strconv.Atoi(matches[0][1])

		game_possible[id] = true
		game_minimums[id] = []int{0, 0, 0}

		draws := strings.Split(matches[0][2], "; ")
		for i := range draws {
			colors := strings.Split(draws[i], ", ")

			for j := range colors {
				split := strings.Split(colors[j], " ")
				count, _ := strconv.Atoi(split[0])
				color := split[1]
				color_id := 0

				if count > maximums[color] {
					game_possible[id] = false
				}

				switch color {
				case "red":
					color_id = 0
				case "blue":
					color_id = 1
				case "green":
					color_id = 2
				}
				game_minimums[id][color_id] = max(game_minimums[id][color_id], count)
			}
		}

	}

	for k, v := range game_possible {
		if v {
			possible_ids_sum += k
		}
	}
	for _, v := range game_minimums {
		power_sum += v[0] * v[1] * v[2]
	}

	return possible_ids_sum, power_sum
}

func main() {
	possible_ids_sum, power_sum := process_game_draws("02/input.txt")
	fmt.Print("Part 1 solution: ", possible_ids_sum, "\nPart 2 solution: ", power_sum, "\n")
}
