package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	aq "github.com/emirpasic/gods/queues/arrayqueue"
)

type Part map[rune]int

func parsePart(input string) Part {
	part := Part{}

	trimmed := strings.Trim(input, "{}")
	split := strings.Split(trimmed, ",")
	for _, rating := range split {
		splitRating := strings.Split(rating, "=")
		category := rune(splitRating[0][0])
		value, _ := strconv.Atoi(splitRating[1])

		part[category] = value
	}

	return part
}

type Agglomerate struct {
	// Disgusting but hey, it works
	ranges       map[rune]*[2]int
	nextWorkflow string
}

func (a Agglomerate) getCombinations() int64 {
	var combinations int64 = 1
	for _, v := range a.ranges {
		combinations *= int64(v[1] - v[0] + 1)
	}
	return combinations
}

func (a Agglomerate) deepCopy() Agglomerate {
	newRanges := make(map[rune]*[2]int)
	for k, v := range a.ranges {
		newVal := [2]int{(*v)[0], (*v)[1]}
		newRanges[k] = &newVal
	}
	return Agglomerate{newRanges, a.nextWorkflow}
}

type Rule struct {
	category    rune
	operator    rune
	threshold   int
	destination string
}

func parseRule(input string) Rule {
	rule := Rule{}

	split := strings.Split(input, ":")

	if len(split) == 1 {
		destination := split[0]
		rule.destination = destination
	} else {
		category := rune(split[0][0])
		operator := rune(split[0][1])
		threshold, _ := strconv.Atoi(split[0][2:])

		destination := split[1]

		rule.category = category
		rule.operator = operator
		rule.threshold = threshold
		rule.destination = destination
	}

	return rule
}

func (r Rule) check(part Part) bool {
	switch r.operator {
	case '>':
		return part[r.category] > r.threshold
	case '<':
		return part[r.category] < r.threshold
	}

	return true
}

func (r Rule) process(agglo Agglomerate) [2]Agglomerate {
	succeeded, failed := agglo.deepCopy(), agglo.deepCopy()

	switch r.operator {
	case '>':
		for k, v := range succeeded.ranges {
			if k == r.category {
				(*v)[0] = max((*v)[0], r.threshold+1)
			}
		}
		for k, v := range failed.ranges {
			if k == r.category {
				(*v)[1] = min((*v)[1], r.threshold)
			}
		}

	case '<':
		for k, v := range succeeded.ranges {
			if k == r.category {
				(*v)[1] = min((*v)[1], r.threshold-1)
			}
		}
		for k, v := range failed.ranges {
			if k == r.category {
				(*v)[0] = max((*v)[0], r.threshold)
			}
		}
	}

	return [2]Agglomerate{succeeded, failed}
}

type Workflow struct {
	rules []Rule
}

func parseWorkflow(input string) (name string, workflow Workflow) {
	split := strings.Split(input, "{")

	name = split[0]

	trimmed := strings.Trim(split[1], "{}")
	rulesRaw := strings.Split(trimmed, ",")
	rules := []Rule{}
	for _, raw := range rulesRaw {
		rules = append(rules, parseRule(raw))
	}

	workflow.rules = rules

	return name, workflow
}

func (w Workflow) check(part Part) string {
	i := 0
	for i < len(w.rules) && !w.rules[i].check(part) {
		i++
	}
	return w.rules[i].destination
}

func (w Workflow) process(agglo Agglomerate) []Agglomerate {
	agglos := []Agglomerate{}

	curr := agglo
	i := 0
	for i < len(w.rules) {
		processed := w.rules[i].process(curr)
		succeeded, failed := processed[0], processed[1]

		succeeded.nextWorkflow = w.rules[i].destination
		agglos = append(agglos, succeeded)

		curr = failed
		i++
	}

	return agglos
}

func processPartsAndWorkflows(path string, minValue, maxValue int) (int, int64) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	inputRaw := strings.Split(strings.Trim(string(file), "\n"), "\n\n")
	workflowsRaw := strings.Split(inputRaw[0], "\n")
	partsRaw := strings.Split(inputRaw[1], "\n")

	workflows := map[string]Workflow{}
	for _, line := range workflowsRaw {
		name, workflow := parseWorkflow(line)
		workflows[name] = workflow
	}

	parts := []Part{}
	for _, line := range partsRaw {
		parts = append(parts, parsePart(line))
	}

	totalRating := 0
	for _, part := range parts {
		curr := workflows["in"].check(part)

		for curr != "R" && curr != "A" {
			curr = workflows[curr].check(part)
		}

		if curr == "A" {
			for _, v := range part {
				totalRating += v
			}
		}
	}

	agglos := aq.New()
	agglos.Enqueue(Agglomerate{
		ranges: map[rune]*[2]int{
			'x': {minValue, maxValue},
			'm': {minValue, maxValue},
			'a': {minValue, maxValue},
			's': {minValue, maxValue},
		},
		nextWorkflow: "in",
	})
	var totalCombinations int64 = 0

	for !agglos.Empty() {
		qe, _ := agglos.Dequeue()
		agglo := qe.(Agglomerate)

		switch agglo.nextWorkflow {
		case "R":
			continue
		case "A":
			totalCombinations += agglo.getCombinations()
		default:
			newAgglos := workflows[agglo.nextWorkflow].process(agglo)
			for _, na := range newAgglos {
				agglos.Enqueue(na)
			}
		}
	}

	return totalRating, totalCombinations
}

func main() {
	totalRating, totalCombinations := processPartsAndWorkflows("19/input.txt", 1, 4000)
	fmt.Print("Part 1 solution: ", totalRating, "\nPart 2 solution: ", totalCombinations, "\n")
}
