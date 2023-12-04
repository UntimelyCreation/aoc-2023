package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func process_scratchcard_game(path string) (int, int) {
	pattern := `Card\s+\d+: (.*)`
	regex := regexp.MustCompile(pattern)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scratchcards := [189]int{}
	for i := range scratchcards {
		scratchcards[i] = 1
	}

	total_scratchcard_points := 0
	scratchcard_count := 0

	line_num := 0
	for scanner.Scan() {
		line := scanner.Text()

		matches := regex.FindAllStringSubmatch(line, -1)

		numbers := strings.Split(matches[0][1], " | ")
		winning_nums_raw := strings.Fields(numbers[0])
		winning_nums := []int{}
		chosen_nums_raw := strings.Fields(numbers[1])

		for i := range winning_nums_raw {
			num, _ := strconv.Atoi(winning_nums_raw[i])
			winning_nums = append(winning_nums, num)
		}

		winners := 0
		for i := range chosen_nums_raw {
			num, _ := strconv.Atoi(chosen_nums_raw[i])
			if slices.Contains(winning_nums, num) {
				winners += 1
			}
		}

		copies := scratchcards[line_num]
		scratchcard_count += copies

		if winners > 0 {
			total_scratchcard_points += int(math.Pow(2, float64(winners-1)))

			for j := line_num + 1; j < line_num+winners+1; j++ {
				scratchcards[j] += copies
			}
		}

		line_num += 1
	}

	return total_scratchcard_points, scratchcard_count
}

func main() {
	total_scratchcard_points, scratchcard_count := process_scratchcard_game("04/input.txt")
	fmt.Print("Part 1 solution: ", total_scratchcard_points, "\nPart 2 solution: ", scratchcard_count, "\n")
}
