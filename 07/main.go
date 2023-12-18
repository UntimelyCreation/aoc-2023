package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Hand struct {
	cards string
	kind  int
	bid   int
}

const (
	highCard = iota
	pair
	twoPair
	threeKind
	fullHouse
	fourKind
	fiveKind
)

var cardValuesJack map[rune]int = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var cardValuesJoker map[rune]int = map[rune]int{
	'J': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'Q': 12,
	'K': 13,
	'A': 14,
}

func (hand *Hand) getHandKind(jokers bool) {
	cardMap := map[rune]int{}

	for _, card := range hand.cards {
		cardMap[rune(card)] += 1
	}
	if jokers {
		jokerCount := cardMap['J']

		if jokerCount != 5 {
			delete(cardMap, 'J')

			maxK, maxV := 'X', 0
			for i := range cardMap {
				if cardMap[i] > maxV {
					maxK = i
					maxV = cardMap[i]
				}
			}
			cardMap[maxK] += jokerCount
		}
	}

	counts := []int{}
	for _, count := range cardMap {
		counts = append(counts, count)
	}
	sort.Ints(counts)

	switch len(cardMap) {
	case 5:
		hand.kind = highCard
	case 4:
		hand.kind = pair
	case 3:
		if reflect.DeepEqual(counts, []int{1, 2, 2}) {
			hand.kind = twoPair
		}
		if reflect.DeepEqual(counts, []int{1, 1, 3}) {
			hand.kind = threeKind
		}
	case 2:
		if reflect.DeepEqual(counts, []int{2, 3}) {
			hand.kind = fullHouse
		}
		if reflect.DeepEqual(counts, []int{1, 4}) {
			hand.kind = fourKind
		}
	case 1:
		hand.kind = fiveKind
	}
}

func (lhs *Hand) isLower(rhs *Hand, jokers bool) bool {
	if lhs.kind != rhs.kind {
		return lhs.kind < rhs.kind
	} else {
		for i := 0; i < 5; i++ {
			if lhs.cards[i] != rhs.cards[i] {
				if jokers {
					return cardValuesJoker[rune(lhs.cards[i])] < cardValuesJoker[rune(rhs.cards[i])]
				} else {
					return cardValuesJack[rune(lhs.cards[i])] < cardValuesJack[rune(rhs.cards[i])]
				}
			}
		}
	}
	return true
}

func getHandSetScores(path string, jokers bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	handsDataRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	hands := []Hand{}

	for _, line := range handsDataRaw {
		split := strings.Fields(line)

		bid, _ := strconv.Atoi(split[1])

		hand := Hand{
			cards: split[0],
			bid:   bid,
		}
		hand.getHandKind(jokers)
		hands = append(hands, hand)
	}

	sort.Slice(hands, func(a, b int) bool {
		return hands[a].isLower(&hands[b], jokers)
	})

	handSetScore := 0
	for i, hand := range hands {
		handSetScore += (i + 1) * hand.bid
	}

	return handSetScore
}

func main() {
	handSetScoreJacks := getHandSetScores("07/input.txt", false)
	handSetScoreJokers := getHandSetScores("07/input.txt", true)
	fmt.Print("Part 1 solution: ", handSetScoreJacks, "\nPart 2 solution: ", handSetScoreJokers, "\n")
}
