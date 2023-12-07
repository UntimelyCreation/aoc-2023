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
	high_card = iota
	pair
	two_pair
	three_kind
	full_house
	four_kind
	five_kind
)

var card_values_jack map[rune]int = map[rune]int{
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

var card_values_joker map[rune]int = map[rune]int{
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

func (hand *Hand) get_hand_kind(jokers bool) {
	card_map := map[rune]int{}

	for _, card := range hand.cards {
		card_map[rune(card)] += 1
	}
	if jokers {
		joker_count := card_map['J']

		if joker_count != 5 {
			delete(card_map, 'J')

			max_k, max_v := 'X', 0
			for i := range card_map {
				if card_map[i] > max_v {
					max_k = i
					max_v = card_map[i]
				}
			}
			card_map[max_k] += joker_count
		}
	}

	counts := []int{}
	for _, count := range card_map {
		counts = append(counts, count)
	}
	sort.Ints(counts)

	switch len(card_map) {
	case 5:
		hand.kind = high_card
	case 4:
		hand.kind = pair
	case 3:
		if reflect.DeepEqual(counts, []int{1, 2, 2}) {
			hand.kind = two_pair
		}
		if reflect.DeepEqual(counts, []int{1, 1, 3}) {
			hand.kind = three_kind
		}
	case 2:
		if reflect.DeepEqual(counts, []int{2, 3}) {
			hand.kind = full_house
		}
		if reflect.DeepEqual(counts, []int{1, 4}) {
			hand.kind = four_kind
		}
	case 1:
		hand.kind = five_kind
	}
}

func (lhs *Hand) is_lower(rhs *Hand, jokers bool) bool {
	if lhs.kind != rhs.kind {
		return lhs.kind < rhs.kind
	} else {
		for i := 0; i < 5; i++ {
			if lhs.cards[i] != rhs.cards[i] {
				if jokers {
					return card_values_joker[rune(lhs.cards[i])] < card_values_joker[rune(rhs.cards[i])]
				} else {
					return card_values_jack[rune(lhs.cards[i])] < card_values_jack[rune(rhs.cards[i])]
				}
			}
		}
	}
	return true
}

func calc_hand_set_scores(path string, jokers bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	hand_data_raw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	hands := []Hand{}

	for _, line := range hand_data_raw {
		split := strings.Fields(line)

		bid, _ := strconv.Atoi(split[1])

		hand := Hand{
			cards: split[0],
			bid:   bid,
		}
		hand.get_hand_kind(jokers)
		hands = append(hands, hand)
	}

	sort.Slice(hands, func(a, b int) bool {
		return hands[a].is_lower(&hands[b], jokers)
	})

	hand_set_score := 0
	for i, hand := range hands {
		hand_set_score += (i + 1) * hand.bid
	}

	return hand_set_score
}

func main() {
	hand_set_score_jacks := calc_hand_set_scores("07/input.txt", false)
	hand_set_score_jokers := calc_hand_set_scores("07/input.txt", true)
	fmt.Print("Part 1 solution: ", hand_set_score_jacks, "\nPart 2 solution: ", hand_set_score_jokers, "\n")
}
