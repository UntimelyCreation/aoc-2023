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

func parseMap(input string) Map {
	ranges := []MapRange{}

	rows := strings.Split(input, "\n")[1:]

	for _, row := range rows {
		nums := strings.Fields(row)
		destRangeStart, _ := strconv.Atoi(nums[0])
		srcRangeStart, _ := strconv.Atoi(nums[1])
		rangeLen, _ := strconv.Atoi(nums[2])
		ranges = append(ranges, MapRange{
			start: srcRangeStart,
			end:   srcRangeStart + rangeLen,
			shift: destRangeStart - srcRangeStart,
		})
	}

	sort.Slice(ranges, func(a, b int) bool {
		return ranges[a].start < ranges[b].start
	})

	return Map{ranges}
}

func (m *Map) mapSeed(seed int) int {
	for _, r := range m.ranges {
		if r.start <= seed && seed < r.end {
			return seed + r.shift
		}
	}
	return seed
}

func (m *Map) mapSeedRange(seeds Range) []Range {
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

func processAlmanac1(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	almanac := strings.Trim(string(file), "\n")
	mapsRaw := strings.Split(almanac, "\n\n")
	maps := []Map{}
	for _, input := range mapsRaw[1:] {
		maps = append(maps, parseMap(input))
	}

	seedsRaw := strings.Fields(strings.Split(strings.Split(mapsRaw[0], "\n")[0], ": ")[1])
	seeds := []int{}
	for _, seedRaw := range seedsRaw {
		seed, _ := strconv.Atoi(seedRaw)
		seeds = append(seeds, seed)
	}
	minSeed := int(^uint(0) >> 1)

	for _, seed := range seeds {
		for _, m := range maps {
			seed = m.mapSeed(seed)
		}
		minSeed = min(minSeed, seed)
	}

	return minSeed
}

func processAlmanac2BruteForce(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	almanac := strings.Trim(string(file), "\n")
	mapsRaw := strings.Split(almanac, "\n\n")
	maps := []Map{}
	for _, input := range mapsRaw[1:] {
		maps = append(maps, parseMap(input))
	}

	seedsRaw := strings.Fields(strings.Split(strings.Split(mapsRaw[0], "\n")[0], ": ")[1])
	seeds := []int{}
	for _, seedRaw := range seedsRaw {
		seed, _ := strconv.Atoi(seedRaw)
		seeds = append(seeds, seed)
	}
	minSeed := int(^uint(0) >> 1)

	for i := 0; i < len(seeds); i += 2 {
		seedRangeStart := seeds[i]
		seedRangeLen := seeds[i+1]

		minSeedTmp := int(^uint(0) >> 1)

		for j := seedRangeStart; j < seedRangeStart+seedRangeLen; j++ {
			seed := j

			for _, m := range maps {
				seed = m.mapSeed(seed)
			}
			minSeedTmp = min(minSeedTmp, seed)
		}

		minSeed = min(minSeed, minSeedTmp)
	}

	return minSeed
}

func processAlmanac2Intervals(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	almanac := strings.Trim(string(file), "\n")
	mapsRaw := strings.Split(almanac, "\n\n")
	maps := []Map{}
	for _, input := range mapsRaw[1:] {
		maps = append(maps, parseMap(input))
	}

	seedsRaw := strings.Fields(strings.Split(strings.Split(mapsRaw[0], "\n")[0], ": ")[1])
	seeds := []int{}
	for _, seedRaw := range seedsRaw {
		seed, _ := strconv.Atoi(seedRaw)
		seeds = append(seeds, seed)
	}
	minSeed := int(^uint(0) >> 1)

	for i := 0; i < len(seeds); i += 2 {
		cur := []Range{{start: seeds[i], end: seeds[i] + seeds[i+1]}}

		for _, m := range maps {
			mapped := []Range{}
			for _, r := range cur {
				result := m.mapSeedRange(r)
				mapped = append(mapped, result...)
			}
			cur = mapped
		}
		for _, r := range cur {
			minSeed = min(minSeed, r.start)
		}
	}

	return minSeed
}

func main() {
	minLocationNumber1 := processAlmanac1("05/input.txt")
	// minLocationNumber2 := processAlmanac2BruteForce("05/input_test.txt")
	minLocationNumber2 := processAlmanac2Intervals("05/input.txt")
	fmt.Print("Part 1 solution: ", minLocationNumber1, "\nPart 2 solution: ", minLocationNumber2, "\n")
}
