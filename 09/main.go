package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func is_zeroes(slice []int) bool {
	for _, val := range slice {
		if val != 0 {
			return false
		}
	}
	return true
}

func last(slice []int) int {
	return slice[len(slice)-1]
}

func reverse(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func calc_prediction_sum(path string, backwards bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	hists_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	hists := [][]int{}
	for _, raws := range hists_raw {
		temp_hist := []int{}

		fields := strings.Fields(raws)
		if backwards {
			fields = reverse(fields)
		}
		for _, raw := range fields {
			val, _ := strconv.Atoi(raw)
			temp_hist = append(temp_hist, val)
		}
		hists = append(hists, temp_hist)
	}

	prediction_sum := 0
	for _, hist := range hists {
		derivs := [][]int{}
		derivs = append(derivs, hist)

		i := 0
		for !is_zeroes(derivs[i]) {
			next_deriv := []int{}
			for j := 1; j < len(derivs[i]); j++ {
				next_deriv = append(next_deriv, derivs[i][j]-derivs[i][j-1])
			}
			derivs = append(derivs, next_deriv)
			i += 1
		}

		n := len(derivs)
		derivs[n-1] = append(derivs[n-1], 0)

		for i := len(derivs) - 2; i >= 0; i-- {
			prediction := last(derivs[i]) + last(derivs[i+1])
			derivs[i] = append(derivs[i], prediction)
			if i == 0 {
				prediction_sum += prediction
			}
		}
	}

	return prediction_sum
}

func main() {
	prediction_sum := calc_prediction_sum("09/input.txt", false)
	prediction_sum_reverse := calc_prediction_sum("09/input.txt", true)
	fmt.Print("Part 1 solution: ", prediction_sum, "\nPart 2 solution: ", prediction_sum_reverse, "\n")
}
