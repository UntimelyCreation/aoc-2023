package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func calc_racing_sum(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	race_data := strings.Split(string(file), "\n")
	times_raw := strings.Fields(race_data[0])[1:]
	dists_raw := strings.Fields(race_data[1])[1:]

	time_raw := ""
	dist_raw := ""

	times := []int{}
	dists := []int{}

	for _, raw := range times_raw {
		time_raw += raw
		time, _ := strconv.Atoi(raw)
		times = append(times, time)
	}
	for _, raw := range dists_raw {
		dist_raw += raw
		dist, _ := strconv.Atoi(raw)
		dists = append(dists, dist)
	}

	time, _ := strconv.Atoi(time_raw)
	dist, _ := strconv.Atoi(dist_raw)

	ways_to_beat_product_1 := 1

	for i := 0; i < len(times); i++ {
		ways_to_beat_race := 0
		time := times[i]
		dist := dists[i]
		for k := 0; k < time; k++ {
			if k*(time-k) > dist {
				ways_to_beat_race++
			}
		}
		ways_to_beat_product_1 *= ways_to_beat_race
	}

	ways_to_beat_race_2 := 0

	for k := 0; k < time; k++ {
		if k*(time-k) > dist {
			ways_to_beat_race_2++
		}
	}

	return ways_to_beat_product_1, ways_to_beat_race_2
}

func main() {
	ways_to_beat_product_1, ways_to_beat_race_2 := calc_racing_sum("06/input.txt")
	fmt.Print("Part 1 solution: ", ways_to_beat_product_1, "\nPart 2 solution: ", ways_to_beat_race_2, "\n")
}
