package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func calculate_calibration_values(path string, pattern string) int {
	regex := regexp.MustCompile(pattern)

	digits := map[string]int{
		"one":   1,
		"1":     1,
		"two":   2,
		"2":     2,
		"three": 3,
		"3":     3,
		"four":  4,
		"4":     4,
		"five":  5,
		"5":     5,
		"six":   6,
		"6":     6,
		"seven": 7,
		"7":     7,
		"eight": 8,
		"8":     8,
		"nine":  9,
		"9":     9,
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	result := 0
	for scanner.Scan() {
		line := scanner.Text()
		matches := []int{}

		for i := range line {
			sub := line[i:]
			if match := regex.FindStringIndex(sub); match != nil {
				l, r := match[0], match[1]
				matches = append(matches, digits[sub[l:r]])
			}
		}

		result += 10*matches[0] + matches[len(matches)-1]
	}

	return result
}

func main() {
	calibration_sum_1 := calculate_calibration_values("01/input.txt", `^\d`)
	fmt.Print("Part 1 solution: ", calibration_sum_1, "\n")

	calibration_sum_2 := calculate_calibration_values("01/input.txt", `^one|^two|^three|^four|^five|^six|^seven|^eight|^nine|^\d`)
	fmt.Print("Part 2 solution: ", calibration_sum_2, "\n")
}
