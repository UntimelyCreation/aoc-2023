package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func get_string_key(platform [][]rune) string {
	key := ""
	for _, row := range platform {
		key += string(row)
	}
	return key
}

func get_key_by_value(m map[string]int, value int) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func transpose(platform [][]rune) [][]rune {
	x, y := len(platform[0]), len(platform)
	transposed := make([][]rune, x)
	for i := range transposed {
		transposed[i] = make([]rune, y)
	}
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			transposed[i][j] = platform[j][i]
		}
	}
	return transposed
}

func rotate_clockwise(platform *[][]rune) {
	n := len(*platform)

	for i := 0; i < n/2; i++ {
		for j := i; j < n-i-1; j++ {
			temp := (*platform)[i][j]
			(*platform)[i][j] = (*platform)[n-1-j][i]
			(*platform)[n-1-j][i] = (*platform)[n-1-i][n-1-j]
			(*platform)[n-1-i][n-1-j] = (*platform)[j][n-1-i]
			(*platform)[j][n-1-i] = temp
		}
	}
}

func shift_platform(platform *[][]rune) {
	shifted_platform := transpose(*platform)

	for _, row := range shifted_platform {
		j := 0
		for j < len(row) {
			for j < len(row) && row[j] != '.' {
				j++
			}
			k := j + 1
			save := k
			for k < len(row) && row[k] == '.' {
				k++
			}
			if k < len(row) {
				switch row[k] {
				case 'O':
					row[j], row[k] = 'O', '.'
					j = save
				case '#':
					j = k
				}
			} else {
				break
			}
		}
	}

	(*platform) = transpose(shifted_platform)
}

func cycle_platform(platform *[][]rune) (int, int, map[string]int) {
	cache := map[string]int{}

	i := 1
	for {
		for j := 0; j < 4; j++ {
			shift_platform(platform)
			rotate_clockwise(platform)
		}
		start, cached := cache[get_string_key(*platform)]
		if cached {
			return start, i - start, cache
		}
		cache[get_string_key(*platform)] = i
		i++
	}
}

func calculate_load(platform [][]rune) int {
	load := 0
	rows := len(platform)

	for i, row := range platform {
		row_load := 0
		for _, r := range row {
			if r == 'O' {
				row_load++
			}
		}
		load += row_load * (rows - i)
	}

	return load
}

func calculate_shifted_support_beam_load(path string) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	platform_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	platform := [][]rune{}
	for _, row := range platform_raw {
		platform = append(platform, []rune(row))
	}

	shift_platform(&platform)

	return calculate_load(platform)
}

func calculate_cycled_support_beam_load(path string, cycles int) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	platform_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	platform := [][]rune{}
	for _, row := range platform_raw {
		platform = append(platform, []rune(row))
	}

	start, period, cache := cycle_platform(&platform)
	shifted_platform_str, _ := get_key_by_value(cache, start+(cycles-start)%period)

	cols := len(platform[0])
	shifted_platform := [][]rune{}
	for i := cols; i <= len(shifted_platform_str); i += cols {
		shifted_platform = append(shifted_platform, []rune(shifted_platform_str[i-cols:i]))
	}

	return calculate_load(shifted_platform)
}

func main() {
	total_shifted_load := calculate_shifted_support_beam_load("14/input.txt")
	total_cycled_load := calculate_cycled_support_beam_load("14/input.txt", 1000000000)
	fmt.Print("Part 1 solution: ", total_shifted_load, "\nPart 2 solution: ", total_cycled_load, "\n")
}
