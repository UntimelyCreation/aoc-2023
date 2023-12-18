package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type Record struct {
	springs []rune
	groups  []int
}

type CacheEntry struct {
	record       Record
	arrangements int
}

func searchCacheByRecord(cache []CacheEntry, record Record) (arrangements int, ok bool) {
	for _, ce := range cache {
		if reflect.DeepEqual(ce.record, record) {
			arrangements, ok = ce.arrangements, true
		}
	}
	return arrangements, ok
}

func (r *Record) countArrangements() int {
	cache := []CacheEntry{}
	return recArrangements(r.springs, r.groups, &cache)
}

func recArrangements(springs []rune, groups []int, cache *[]CacheEntry) int {
	record := Record{springs, groups}
	arrangements, cached := searchCacheByRecord(*cache, record)
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
		result = recArrangements(springs[1:], groups, cache)
	case '?':
		springsA := []rune{'#'}
		springsA = append(springsA, springs[1:]...)
		springsB := []rune{'.'}
		springsB = append(springsB, springs[1:]...)
		result = recArrangements(springsA, groups, cache) + recArrangements(springsB, groups, cache)
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
			result = recArrangements(springs[groups[0]+1:], groups[1:], cache)
		} else {
			result = recArrangements(springs[groups[0]:], groups[1:], cache)
		}
	}

	*cache = append(*cache, CacheEntry{record: record, arrangements: result})
	return result
}

func processSpringRecords(path string, folds int) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	recordsRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	records := []Record{}
	for _, row := range recordsRaw {
		split := strings.Fields(row)

		groupsRaw := strings.Split(split[1], ",")
		groupsFold := []int{}
		for _, numRaw := range groupsRaw {
			num, _ := strconv.Atoi(numRaw)
			groupsFold = append(groupsFold, num)
		}

		springs := []rune{}
		groups := []int{}
		for i := 0; i < folds; i++ {
			springs = append(springs, []rune(split[0])...)
			if i != folds-1 {
				springs = append(springs, '?')
			}
			groups = append(groups, groupsFold...)
		}
		record := Record{
			springs,
			groups,
		}
		records = append(records, record)
	}

	totalArrangements := 0
	for _, r := range records {
		totalArrangements += r.countArrangements()
	}

	return totalArrangements
}

func main() {
	arragementsSum1 := processSpringRecords("12/input.txt", 1)
	arragementsSum2 := processSpringRecords("12/input.txt", 5)
	fmt.Print("Part 1 solution: ", arragementsSum1, "\nPart 2 solution: ", arragementsSum2, "\n")
}
