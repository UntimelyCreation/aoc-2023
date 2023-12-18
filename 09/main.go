package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/utils"
)

func isZeroes(slice []int) bool {
	for _, val := range slice {
		if val != 0 {
			return false
		}
	}
	return true
}

func calcPredictionSum(path string, backwards bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	histsRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	hists := [][]int{}
	for _, raws := range histsRaw {
		tempHist := []int{}

		fields := strings.Fields(raws)
		if backwards {
			fields = utils.Reverse(fields)
		}
		for _, raw := range fields {
			val, _ := strconv.Atoi(raw)
			tempHist = append(tempHist, val)
		}
		hists = append(hists, tempHist)
	}

	predictionSum := 0
	for _, hist := range hists {
		derivs := [][]int{}
		derivs = append(derivs, hist)

		i := 0
		for !isZeroes(derivs[i]) {
			nextDeriv := []int{}
			for j := 1; j < len(derivs[i]); j++ {
				nextDeriv = append(nextDeriv, derivs[i][j]-derivs[i][j-1])
			}
			derivs = append(derivs, nextDeriv)
			i += 1
		}

		n := len(derivs)
		derivs[n-1] = append(derivs[n-1], 0)

		for i := len(derivs) - 2; i >= 0; i-- {
			prediction := utils.Last(derivs[i]) + utils.Last(derivs[i+1])
			derivs[i] = append(derivs[i], prediction)
			if i == 0 {
				predictionSum += prediction
			}
		}
	}

	return predictionSum
}

func main() {
	predictionSum := calcPredictionSum("09/input.txt", false)
	predictionSumReverse := calcPredictionSum("09/input.txt", true)
	fmt.Print("Part 1 solution: ", predictionSum, "\nPart 2 solution: ", predictionSumReverse, "\n")
}
