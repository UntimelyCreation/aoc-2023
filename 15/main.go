package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Lens struct {
	label        string
	focal_length int
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

func restore_focusing_lenses(path string) (int, int) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	sequence_raw := strings.Trim(string(file), "\n")
	sequence := strings.Split(sequence_raw, ",")

	sequence_hash_sum := 0
	for _, step := range sequence {
		sequence_hash_sum += hash(step)
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
			focal_length := int(step[j+1] - '0')
			if k < len(*box) {
				(*box)[k].focal_length = focal_length
			}
			if k == len(*box) {
				*box = append(*box, Lens{label, focal_length})
			}
		case '-':
			if k < len(*box) {
				*box = append((*box)[:k], (*box)[k+1:]...)
			}
		}
	}

	lenses_focusing_power := 0
	for i, box := range boxes {
		for j := range box {
			lenses_focusing_power += (i + 1) * (j + 1) * (box[j].focal_length)
		}
	}

	return sequence_hash_sum, lenses_focusing_power
}

func main() {
	sequence_hash_sum, lenses_focusing_power := restore_focusing_lenses("15/input.txt")
	fmt.Print("Part 1 solution: ", sequence_hash_sum, "\nPart 2 solution: ", lenses_focusing_power, "\n")
}
