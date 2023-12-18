package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Lens struct {
	label       string
	focalLength int
}

func hash(step string) int {
	curr := 0
	for _, r := range step {
		curr += int(r)
		curr *= 17
		curr = curr % 256
	}
	return curr
}

func restoreFocusingLenses(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	sequenceRaw := strings.Trim(string(file), "\n")
	sequence := strings.Split(sequenceRaw, ",")

	sequenceHashSum := 0
	for _, step := range sequence {
		sequenceHashSum += hash(step)
	}

	boxes := [256][]Lens{}

	for _, step := range sequence {

		j := 0
		for j < len(step) && step[j] != '=' && step[j] != '-' {
			j++
		}

		label := step[:j]
		box := &boxes[hash(label)]

		k := 0
		for k < len(*box) && (*box)[k].label != label {
			k++
		}

		switch step[j] {
		case '=':
			focalLength := int(step[j+1] - '0')
			if k < len(*box) {
				(*box)[k].focalLength = focalLength
			}
			if k == len(*box) {
				*box = append(*box, Lens{label, focalLength})
			}
		case '-':
			if k < len(*box) {
				*box = append((*box)[:k], (*box)[k+1:]...)
			}
		}
	}

	lensesFocusingPower := 0
	for i, box := range boxes {
		for j := range box {
			lensesFocusingPower += (i + 1) * (j + 1) * (box[j].focalLength)
		}
	}

	return sequenceHashSum, lensesFocusingPower
}

func main() {
	sequenceHashSum, lensesFocusingPower := restoreFocusingLenses("15/input.txt")
	fmt.Print("Part 1 solution: ", sequenceHashSum, "\nPart 2 solution: ", lensesFocusingPower, "\n")
}
