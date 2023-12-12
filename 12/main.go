package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Record struct {
	springs []rune
	groups  []int
}

func (r *Record) count_possible_arrangements() int {
	cache := map[string]int{}
	return rec_arrangements(r.springs, r.groups, &cache)
}

// Create unique memoization key from spring AND group state
func get_cache_key(springs []rune, groups []int) string {
	springs_str := string(springs)
	groups_str := fmt.Sprint(groups)
	return springs_str + groups_str
}

func rec_arrangements(springs []rune, groups []int, cache *map[string]int) int {
	key := get_cache_key(springs, groups)

	arrangements, cached := (*cache)[key]
	if cached {
		return arrangements
	}

	if len(springs) == 0 {
		if len(groups) == 0 {
			return 1
		} else {
			return 0
		}
	}

	result := -1

	switch springs[0] {
	case '.':
		result = rec_arrangements(springs[1:], groups, cache)
	case '?':
		springs_a := []rune{'#'}
		springs_a = append(springs_a, springs[1:]...)
		springs_b := []rune{'.'}
		springs_b = append(springs_b, springs[1:]...)
		result = rec_arrangements(springs_a, groups, cache) + rec_arrangements(springs_b, groups, cache)
	case '#':
		if len(groups) == 0 {
			result = 0
			break
		}
		if len(springs) < groups[0] {
			result = 0
			break
		}
		if slices.Contains(springs[:groups[0]], '.') {
			result = 0
			break
		}
		if len(groups) > 1 {
			if len(springs) < groups[0]+1 || springs[groups[0]] == '#' {
				result = 0
				break
			}
			result = rec_arrangements(springs[groups[0]+1:], groups[1:], cache)
		} else {
			result = rec_arrangements(springs[groups[0]:], groups[1:], cache)
		}
	}

	(*cache)[key] = result
	return result
}

func process_spring_records(path string, folds int) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	records_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	records := []Record{}
	for _, row := range records_raw {
		split := strings.Fields(row)

		groups_raw := strings.Split(split[1], ",")
		groups_fold := []int{}
		for _, num_raw := range groups_raw {
			num, _ := strconv.Atoi(num_raw)
			groups_fold = append(groups_fold, num)
		}

		springs := []rune{}
		groups := []int{}
		for i := 0; i < folds; i++ {
			springs = append(springs, []rune(split[0])...)
			if i != folds-1 {
				springs = append(springs, '?')
			}
			groups = append(groups, groups_fold...)
		}
		record := Record{
			springs,
			groups,
		}
		records = append(records, record)
	}

	total_possible_arrangements := 0
	for _, r := range records {
		total_possible_arrangements += r.count_possible_arrangements()
	}

	return total_possible_arrangements
}

func main() {
	arragements_sum_1 := process_spring_records("12/input.txt", 1)
	arragements_sum_2 := process_spring_records("12/input.txt", 5)
	fmt.Print("Part 1 solution: ", arragements_sum_1, "\nPart 2 solution: ", arragements_sum_2, "\n")
}
