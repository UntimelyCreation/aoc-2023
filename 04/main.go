package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func processScratchcardGame(path string) (int, int) {
	pattern := `Card\s+\d+: (.*)`
	regex := regexp.MustCompile(pattern)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scratchcards := [189]int{}
	for i := range scratchcards {
		scratchcards[i] = 1
	}

	totalPoints := 0
	scratchcardCount := 0

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()

		matches := regex.FindAllStringSubmatch(line, -1)

		numbers := strings.Split(matches[0][1], " | ")
		winningNumsRaw := strings.Fields(numbers[0])
		winningNums := []int{}
		chosenNumsRaw := strings.Fields(numbers[1])

		for i := range winningNumsRaw {
			num, _ := strconv.Atoi(winningNumsRaw[i])
			winningNums = append(winningNums, num)
		}

		winners := 0
		for i := range chosenNumsRaw {
			num, _ := strconv.Atoi(chosenNumsRaw[i])
			if slices.Contains(winningNums, num) {
				winners += 1
			}
		}

		copies := scratchcards[lineNum]
		scratchcardCount += copies

		if winners > 0 {
			totalPoints += int(math.Pow(2, float64(winners-1)))

			for j := lineNum + 1; j < lineNum+winners+1; j++ {
				scratchcards[j] += copies
			}
		}

		lineNum += 1
	}

	return totalPoints, scratchcardCount
}

func main() {
	totalPoints, scratchcardCount := processScratchcardGame("04/input.txt")
	fmt.Print("Part 1 solution: ", totalPoints, "\nPart 2 solution: ", scratchcardCount, "\n")
}
