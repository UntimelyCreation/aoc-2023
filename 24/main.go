package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/space"
	"gonum.org/v1/gonum/mat"
)

const (
	// rangeMin = 7
	// rangeMax = 27
	rangeMin = 2e14
	rangeMax = 4e14
)

type Hailstone struct {
	position, velocity space.Position[float64]
}

func (h Hailstone) get2DLineCoefficients() (float64, float64) {
	slope := h.velocity.Y / h.velocity.X
	intercept := h.position.Y - slope*h.position.X

	return slope, intercept
}

func (h Hailstone) checkFuture(intersection space.Position[float64]) bool {
	futureX := (intersection.X-h.position.X)/h.velocity.X > 0
	futureY := (intersection.Y-h.position.Y)/h.velocity.Y > 0

	return futureX && futureY
}

func (h Hailstone) intersects2DInFuture(other Hailstone) (intersection space.Position[float64], ok bool) {
	hSlope, hIntercept := h.get2DLineCoefficients()
	otherSlope, otherIntercept := other.get2DLineCoefficients()

	if hSlope == otherSlope {
		return intersection, ok
	}

	interX := (otherIntercept - hIntercept) / (hSlope - otherSlope)
	interY := hSlope*interX + hIntercept

	intersection = space.Position[float64]{X: interX, Y: interY, Z: 0}

	if h.checkFuture(intersection) && other.checkFuture(intersection) {
		ok = true
	}

	return intersection, ok
}

func getHailstoneIntersections(path string, iterations int) (int, int) {
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

	validIntersections := 0
	for i, h := range hailstones {
		for _, other := range hailstones[i+1:] {
			if intersection, ok := h.intersects2DInFuture(other); ok {
				if intersection.X >= rangeMin &&
					intersection.X <= rangeMax &&
					intersection.Y >= rangeMin &&
					intersection.Y <= rangeMax {
					validIntersections++
				}
			}
		}
	}

	// Linear algebra: cross-product (posHi - posRock) ^ (velHi - velRock) = 0, where Hi is the ith hailstone
	// There are 6 unknowns, so 6 equations are formed by equating for different hailstones, e.g. i with i+1 and i with i+2
	rockPosSumTotal := float64(0)
	for i := 0; i < iterations; i++ {

		h1, h2, h3 := hailstones[i], hailstones[i+1], hailstones[i+2]

		A := mat.NewDense(6, 6, []float64{
			-(h1.velocity.Y - h2.velocity.Y), h1.velocity.X - h2.velocity.X, 0, h1.position.Y - h2.position.Y, -(h1.position.X - h2.position.X), 0,
			-(h1.velocity.Y - h3.velocity.Y), h1.velocity.X - h3.velocity.X, 0, h1.position.Y - h3.position.Y, -(h1.position.X - h3.position.X), 0,

			0, -(h1.velocity.Z - h2.velocity.Z), h1.velocity.Y - h2.velocity.Y, 0, h1.position.Z - h2.position.Z, -(h1.position.Y - h2.position.Y),
			0, -(h1.velocity.Z - h3.velocity.Z), h1.velocity.Y - h3.velocity.Y, 0, h1.position.Z - h3.position.Z, -(h1.position.Y - h3.position.Y),

			-(h1.velocity.Z - h2.velocity.Z), 0, h1.velocity.X - h2.velocity.X, h1.position.Z - h2.position.Z, 0, -(h1.position.X - h2.position.X),
			-(h1.velocity.Z - h3.velocity.Z), 0, h1.velocity.X - h3.velocity.X, h1.position.Z - h3.position.Z, 0, -(h1.position.X - h3.position.X),
		})
		b := mat.NewVecDense(6, []float64{
			(h1.position.Y*h1.velocity.X - h2.position.Y*h2.velocity.X) - (h1.position.X*h1.velocity.Y - h2.position.X*h2.velocity.Y),
			(h1.position.Y*h1.velocity.X - h3.position.Y*h3.velocity.X) - (h1.position.X*h1.velocity.Y - h3.position.X*h3.velocity.Y),

			(h1.position.Z*h1.velocity.Y - h2.position.Z*h2.velocity.Y) - (h1.position.Y*h1.velocity.Z - h2.position.Y*h2.velocity.Z),
			(h1.position.Z*h1.velocity.Y - h3.position.Z*h3.velocity.Y) - (h1.position.Y*h1.velocity.Z - h3.position.Y*h3.velocity.Z),

			(h1.position.Z*h1.velocity.X - h2.position.Z*h2.velocity.X) - (h1.position.X*h1.velocity.Z - h2.position.X*h2.velocity.Z),
			(h1.position.Z*h1.velocity.X - h3.position.Z*h3.velocity.X) - (h1.position.X*h1.velocity.Z - h3.position.X*h3.velocity.Z),
		})

		var rock mat.VecDense
		if err := rock.SolveVec(A, b); err != nil {
			fmt.Println(err)
		}

		rockPosSumTotal += rock.At(0, 0) + rock.At(1, 0) + rock.At(2, 0)
	}

	return validIntersections, int(rockPosSumTotal) / iterations
}

func main() {
	// Iteration count may need fiddling to get exact value
	validIntersections, rockPosSum := getHailstoneIntersections("24/input.txt", 25)
	fmt.Print("Part 1 solution: ", validIntersections, "\nPart 2 solution: ", rockPosSum, "\n")
}
