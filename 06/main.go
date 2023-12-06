package main

import (
	"fmt"
	"log"
	"math"
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

func calc_racing_sum_opt(path string) (int, int) {
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
		time := times[i]
		dist := dists[i]

		disc := math.Sqrt(float64(time*time - 4*dist))
		// int(Floor(x + 1)), int(Ceil(x - 1)) deals with edge cases
		a, b := int(math.Floor((-float64(time)+disc)/-2+1)), int(math.Ceil((-float64(time)-disc)/-2-1))
		ways_to_beat_product_1 *= (b - a) + 1
	}

	disc := math.Sqrt(float64(time*time - 4*dist))
	a, b := int(math.Floor((-float64(time)+disc)/-2+1)), int(math.Ceil((-float64(time)-disc)/-2-1))
	ways_to_beat_race_2 := (b - a) + 1

	return ways_to_beat_product_1, ways_to_beat_race_2
}

func main() {
	// ways_to_beat_product_1, ways_to_beat_race_2 := calc_racing_sum("06/input.txt")
	ways_to_beat_product_1, ways_to_beat_race_2 := calc_racing_sum_opt("06/input.txt")
	fmt.Print("Part 1 solution: ", ways_to_beat_product_1, "\nPart 2 solution: ", ways_to_beat_race_2, "\n")
}
