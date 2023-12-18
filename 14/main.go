package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/grid"
)

type CacheEntry struct {
	grid  grid.Grid[rune]
	start int
}

func searchCacheByGrid(cache []CacheEntry, platform grid.Grid[rune]) (start int, ok bool) {
	for _, ce := range cache {
		if reflect.DeepEqual(ce.grid, platform) {
			start, ok = ce.start, true
		}
	}
	return start, ok
}

func searchCacheByStart(cache []CacheEntry, start int) (grid grid.Grid[rune], ok bool) {
	for _, ce := range cache {
		if ce.start == start {
			grid, ok = ce.grid, true
		}
	}
	return grid, ok
}

func shift(platform grid.Grid[rune]) grid.Grid[rune] {
	shiftedPlatform := platform.Transpose()

	xMin, xMax := shiftedPlatform.XRange()
	yMin, yMax := shiftedPlatform.YRange()

	for i := xMin; i <= xMax; i++ {
		j := yMin
		for j <= yMax {
			for j <= yMax && shiftedPlatform[grid.Position{Row: i, Col: j}] != '.' {
				j++
			}
			k, temp := j+1, j+1
			for k <= yMax && shiftedPlatform[grid.Position{Row: i, Col: k}] == '.' {
				k++
			}
			if k <= yMax {
				switch shiftedPlatform[grid.Position{Row: i, Col: k}] {
				case 'O':
					shiftedPlatform[grid.Position{Row: i, Col: j}], shiftedPlatform[grid.Position{Row: i, Col: k}] = 'O', '.'
					j = temp
				case '#':
					j = k
				}
			} else {
				break
			}
		}
	}

	return shiftedPlatform.Transpose()
}

func cycle(platform grid.Grid[rune]) (int, int, []CacheEntry) {
	cache := []CacheEntry{}

	i := 1
	for {
		for j := 0; j < 4; j++ {
			platform = shift(platform)
			platform = platform.Rotate(complex(0, -1))
		}
		if start, cached := searchCacheByGrid(cache, platform); cached {
			return start, i - start, cache
		}
		cache = append(cache, CacheEntry{grid: platform, start: i})
		i++
	}
}

func getLoad(platform grid.Grid[rune]) int {
	load := 0
	rows, _ := platform.Dimensions()

	for k, v := range platform {
		if v == 'O' {
			load += rows - k.Row
		}
	}

	return load
}

func calcSupportBeamLoad(path string, cycles int) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	platformRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	platform := grid.Grid[rune]{}
	for row, line := range platformRaw {
		for col, r := range line {
			platform[grid.Position{Row: row, Col: col}] = r
		}
	}

	shiftedPlatform := shift(platform)

	start, period, cache := cycle(platform)
	cycledPlatform, _ := searchCacheByStart(cache, start+(cycles-start)%period)

	return getLoad(shiftedPlatform), getLoad(cycledPlatform)
}

func main() {
	totalShiftedLoad, totalCycledLoad := calcSupportBeamLoad("14/input.txt", 1000000000)
	fmt.Print("Part 1 solution: ", totalShiftedLoad, "\nPart 2 solution: ", totalCycledLoad, "\n")
}
