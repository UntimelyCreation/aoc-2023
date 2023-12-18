package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func calcRacingSum(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	raceData := strings.Split(string(file), "\n")
	timesRaw := strings.Fields(raceData[0])[1:]
	distsRaw := strings.Fields(raceData[1])[1:]

	timeRaw := ""
	distRaw := ""

	times := []int{}
	dists := []int{}

	for _, raw := range timesRaw {
		timeRaw += raw
		time, _ := strconv.Atoi(raw)
		times = append(times, time)
	}
	for _, raw := range distsRaw {
		distRaw += raw
		dist, _ := strconv.Atoi(raw)
		dists = append(dists, dist)
	}

	time, _ := strconv.Atoi(timeRaw)
	dist, _ := strconv.Atoi(distRaw)

	waysToBeatProduct1 := 1

	for i := 0; i < len(times); i++ {
		waysToBeatRace := 0
		time := times[i]
		dist := dists[i]
		for k := 0; k < time; k++ {
			if k*(time-k) > dist {
				waysToBeatRace++
			}
		}
		waysToBeatProduct1 *= waysToBeatRace
	}

	waysToBeatRace2 := 0

	for k := 0; k < time; k++ {
		if k*(time-k) > dist {
			waysToBeatRace2++
		}
	}

	return waysToBeatProduct1, waysToBeatRace2
}

func calcRacingSumOpt(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	raceData := strings.Split(string(file), "\n")
	timesRaw := strings.Fields(raceData[0])[1:]
	distsRaw := strings.Fields(raceData[1])[1:]

	timeRaw := ""
	distRaw := ""

	times := []int{}
	dists := []int{}

	for _, raw := range timesRaw {
		timeRaw += raw
		time, _ := strconv.Atoi(raw)
		times = append(times, time)
	}
	for _, raw := range distsRaw {
		distRaw += raw
		dist, _ := strconv.Atoi(raw)
		dists = append(dists, dist)
	}

	time, _ := strconv.Atoi(timeRaw)
	dist, _ := strconv.Atoi(distRaw)

	waysToBeatProduct1 := 1

	for i := 0; i < len(times); i++ {
		time := times[i]
		dist := dists[i]

		disc := math.Sqrt(float64(time*time - 4*dist))
		// int(Floor(x + 1)), int(Ceil(x - 1)) deals with edge cases
		a, b := int(math.Floor((-float64(time)+disc)/-2+1)), int(math.Ceil((-float64(time)-disc)/-2-1))
		waysToBeatProduct1 *= (b - a) + 1
	}

	disc := math.Sqrt(float64(time*time - 4*dist))
	a, b := int(math.Floor((-float64(time)+disc)/-2+1)), int(math.Ceil((-float64(time)-disc)/-2-1))
	waysToBeatRace2 := (b - a) + 1

	return waysToBeatProduct1, waysToBeatRace2
}

func main() {
	// waysToBeatProduct1, waysToBeatRace2 := calcRacingSum("06/input.txt")
	waysToBeatProduct1, waysToBeatRace2 := calcRacingSumOpt("06/input.txt")
	fmt.Print("Part 1 solution: ", waysToBeatProduct1, "\nPart 2 solution: ", waysToBeatRace2, "\n")
}
