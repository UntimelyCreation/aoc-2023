package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/space"
	"github.com/UntimelyCreation/aoc-2023-go/pkg/utils"
)

const (
	// rangeMin = 7
	// rangeMax = 27
	rangeMin = 2e14
	rangeMax = 4e14
)

func different(a, b float64) bool {
	// Compensates for floating point precision errors
	// Arbitrary threshold, chosen to make the solution work
	return math.Abs(a-b) > 1
}

type Intersection struct {
	position space.Position[float64]
	time     float64
}

type Hailstone struct {
	position, velocity space.Position[float64]
}

func (h Hailstone) get2DLineCoefficients() (float64, float64) {
	slope := h.velocity.Y / h.velocity.X
	intercept := h.position.Y - slope*h.position.X

	return slope, intercept
}

func (h Hailstone) checkFuture(intersection Intersection) bool {
	futureX := utils.Sign(intersection.position.X-h.position.X) == utils.Sign(h.velocity.X)
	futureY := utils.Sign(intersection.position.Y-h.position.Y) == utils.Sign(h.velocity.Y)

	return futureX && futureY
}

func (h Hailstone) intersects2DInFuture(other Hailstone) (intersection Intersection, ok bool) {
	hSlope, hIntercept := h.get2DLineCoefficients()
	otherSlope, otherIntercept := other.get2DLineCoefficients()

	if hSlope == otherSlope {
		return intersection, ok
	}

	interX := (otherIntercept - hIntercept) / (hSlope - otherSlope)
	interY := hSlope*interX + hIntercept
	interTime := (interX - h.position.X) / h.velocity.X

	interPos := space.Position[float64]{X: interX, Y: interY, Z: 0}
	intersection = Intersection{position: interPos, time: interTime}

	if h.checkFuture(intersection) && other.checkFuture(intersection) {
		ok = true
	}

	return intersection, ok
}

func getHailstoneIntersections(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	hailstonesRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	hailstones := []Hailstone{}
	for _, line := range hailstonesRaw {
		split := strings.Split(line, " @ ")

		posRaw := strings.Split(split[0], ", ")
		velRaw := strings.Split(split[1], ", ")

		posX, _ := strconv.Atoi(strings.TrimSpace(posRaw[0]))
		posY, _ := strconv.Atoi(strings.TrimSpace(posRaw[1]))
		posZ, _ := strconv.Atoi(strings.TrimSpace(posRaw[2]))

		velX, _ := strconv.Atoi(strings.TrimSpace(velRaw[0]))
		velY, _ := strconv.Atoi(strings.TrimSpace(velRaw[1]))
		velZ, _ := strconv.Atoi(strings.TrimSpace(velRaw[2]))

		hailstones = append(hailstones, Hailstone{
			position: space.Position[float64]{X: float64(posX), Y: float64(posY), Z: float64(posZ)},
			velocity: space.Position[float64]{X: float64(velX), Y: float64(velY), Z: float64(velZ)},
		})
	}

	velXMin, velXMax := math.MaxFloat64, float64(0)
	velYMin, velYMax := math.MaxFloat64, float64(0)
	velZMin, velZMax := math.MaxFloat64, float64(0)
	validIntersections := 0
	for i, h := range hailstones {

		velXMin = min(velXMin, h.velocity.X)
		velXMax = max(velXMax, h.velocity.X)
		velYMin = min(velYMin, h.velocity.Y)
		velYMax = max(velYMax, h.velocity.Y)
		velZMin = min(velZMin, h.velocity.Z)
		velZMax = max(velZMax, h.velocity.Z)

		for _, other := range hailstones[i+1:] {
			if intersection, ok := h.intersects2DInFuture(other); ok {
				if intersection.position.X >= rangeMin &&
					intersection.position.X <= rangeMax &&
					intersection.position.Y >= rangeMin &&
					intersection.position.Y <= rangeMax {
					validIntersections++
				}
			}
		}
	}

	rockPosSum := 0
	// -2 and +2 are for the test input
findRockPosition:
	for vX := velXMin - 2; vX <= velXMax+2; vX++ {
	calcInter2D:
		for vY := velYMin - 2; vY <= velYMax+2; vY++ {
		calcInter3D:
			for vZ := velZMin - 2; vZ <= velZMax+2; vZ++ {
				adjustedHailstones := []Hailstone{}
				for _, h := range hailstones[:5] {
					adjustedHailstones = append(adjustedHailstones, Hailstone{
						position: h.position,
						velocity: space.Position[float64]{
							X: h.velocity.X - vX,
							Y: h.velocity.Y - vY,
							Z: h.velocity.Z - vZ,
						},
					})
				}

				intersections := []Intersection{}
				for _, h := range adjustedHailstones[1:] {
					intersection, ok := h.intersects2DInFuture(adjustedHailstones[0])
					if !ok {
						continue calcInter2D
					}
					intersections = append(intersections, intersection)
				}
				for _, intersection := range intersections[1:] {
					if different(intersection.position.X, intersections[0].position.X) ||
						different(intersection.position.Y, intersections[0].position.Y) {
						continue calcInter2D
					}
				}

				interZs := []float64{}
				for i, h := range adjustedHailstones[1:] {
					interZs = append(interZs, h.position.Z+intersections[i].time*h.velocity.Z)
				}
				for _, interZ := range interZs[1:] {
					if different(interZ, interZs[0]) {
						continue calcInter3D
					}
				}

				rockPosSum = int(intersections[0].position.X + intersections[0].position.Y + interZs[0])
				break findRockPosition
			}
		}
	}

	return validIntersections, rockPosSum
}

func main() {
	validIntersections, rockPosSum := getHailstoneIntersections("24/input.txt")
	fmt.Print("Part 1 solution: ", validIntersections, "\nPart 2 solution: ", rockPosSum, "\n")
}
