package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	start int
	end   int
}

type MapRange struct {
	start int
	end   int
	shift int
}

type Map struct {
	ranges []MapRange
}

func parse_map(input string) Map {
	ranges := []MapRange{}

	rows := strings.Split(input, "\n")[1:]

	for _, row := range rows {
		nums := strings.Fields(row)
		dest_range_start, _ := strconv.Atoi(nums[0])
		src_range_start, _ := strconv.Atoi(nums[1])
		range_len, _ := strconv.Atoi(nums[2])
		ranges = append(ranges, MapRange{
			start: src_range_start,
			end:   src_range_start + range_len,
			shift: dest_range_start - src_range_start,
		})
	}

	sort.Slice(ranges, func(a, b int) bool {
		return ranges[a].start < ranges[b].start
	})

	return Map{ranges}
}

func (m *Map) map_seed(seed int) int {
	for _, r := range m.ranges {
		if r.start <= seed && seed < r.end {
			return seed + r.shift
		}
	}
	return seed
}

func (m *Map) map_seed_range(seeds Range) []Range {
	result := []Range{}

	start := seeds.start
	cur := start

	i := 0
	for i < len(m.ranges) && m.ranges[i].end < start {
		i += 1
	}
	for _, r := range m.ranges[i:] {
		if r.start > cur {
			result = append(result, Range{start: cur, end: min(seeds.end, r.start)})
			cur = r.start
		}
		if cur >= r.end {
			break
		}
		result = append(result, Range{start: cur + r.shift, end: min(seeds.end, r.end) + r.shift})
		cur = r.end
		if cur >= seeds.end {
			break
		}
	}
	if cur < seeds.end {
		result = append(result, Range{start: cur, end: seeds.end})
	}

	return result
}

func process_almanac_1(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	almanac := strings.Trim(string(file), "\n")
	maps_raw := strings.Split(almanac, "\n\n")
	maps := []Map{}
	for _, input := range maps_raw[1:] {
		maps = append(maps, parse_map(input))
	}

	seeds_raw := strings.Fields(strings.Split(strings.Split(maps_raw[0], "\n")[0], ": ")[1])
	seeds := []int{}
	for _, seed_raw := range seeds_raw {
		seed, _ := strconv.Atoi(seed_raw)
		seeds = append(seeds, seed)
	}
	min_seed := int(^uint(0) >> 1)

	for _, seed := range seeds {
		for _, m := range maps {
			seed = m.map_seed(seed)
		}
		min_seed = min(min_seed, seed)
	}

	return min_seed
}

func process_almanac_2_brute_force(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	almanac := strings.Trim(string(file), "\n")
	maps_raw := strings.Split(almanac, "\n\n")
	maps := []Map{}
	for _, input := range maps_raw[1:] {
		maps = append(maps, parse_map(input))
	}

	seeds_raw := strings.Fields(strings.Split(strings.Split(maps_raw[0], "\n")[0], ": ")[1])
	seeds := []int{}
	for _, seed_raw := range seeds_raw {
		seed, _ := strconv.Atoi(seed_raw)
		seeds = append(seeds, seed)
	}
	min_seed := int(^uint(0) >> 1)

	for i := 0; i < len(seeds); i += 2 {
		seed_range_start := seeds[i]
		seed_range_len := seeds[i+1]

		min_seed_tmp := int(^uint(0) >> 1)

		for j := seed_range_start; j < seed_range_start+seed_range_len; j++ {
			seed := j

			for _, m := range maps {
				seed = m.map_seed(seed)
			}
			min_seed_tmp = min(min_seed_tmp, seed)
		}

		min_seed = min(min_seed, min_seed_tmp)
	}

	return min_seed
}

func process_almanac_2_intervals(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	almanac := strings.Trim(string(file), "\n")
	maps_raw := strings.Split(almanac, "\n\n")
	maps := []Map{}
	for _, input := range maps_raw[1:] {
		maps = append(maps, parse_map(input))
	}

	seeds_raw := strings.Fields(strings.Split(strings.Split(maps_raw[0], "\n")[0], ": ")[1])
	seeds := []int{}
	for _, seed_raw := range seeds_raw {
		seed, _ := strconv.Atoi(seed_raw)
		seeds = append(seeds, seed)
	}
	min_seed := int(^uint(0) >> 1)

	for i := 0; i < len(seeds); i += 2 {
		cur := []Range{{start: seeds[i], end: seeds[i] + seeds[i+1]}}

		for _, m := range maps {
			mapped := []Range{}
			for _, r := range cur {
				result := m.map_seed_range(r)
				mapped = append(mapped, result...)
			}
			cur = mapped
		}
		for _, r := range cur {
			min_seed = min(min_seed, r.start)
		}
	}

	return min_seed
}

func main() {
	min_location_number_1 := process_almanac_1("05/input.txt")
	// min_location_number_2 := process_almanac_2_brute_force("05/input_test.txt")
	min_location_number_2 := process_almanac_2_intervals("05/input.txt")
	fmt.Print("Part 1 solution: ", min_location_number_1, "\nPart 2 solution: ", min_location_number_2, "\n")
}
