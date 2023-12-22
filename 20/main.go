package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/UntimelyCreation/aoc-2023-go/pkg/utils"
	aq "github.com/emirpasic/gods/queues/arrayqueue"
)

var (
	modules       map[string]Module
	moduleInputs  map[string][]string
	moduleOutputs map[string][]string
	messages      *aq.Queue
)

type Pulse int

const (
	Low Pulse = iota
	High
)

type Message struct {
	sender   string
	pulse    Pulse
	receiver string
}

type Module interface {
	Send(sender string, pulse Pulse)
}

type FlipFlop struct {
	name    string
	powered bool
}

func (ff *FlipFlop) Send(sender string, pulse Pulse) {
	if pulse == Low {
		ff.powered = !ff.powered
		newPulse := Low
		if ff.powered {
			newPulse = High
		}
		for _, dest := range moduleOutputs[ff.name] {
			messages.Enqueue(Message{ff.name, newPulse, dest})
		}
	}
}

type Conjunction struct {
	name string
	mrp  map[string]Pulse
}

func (c *Conjunction) Send(sender string, pulse Pulse) {
	c.mrp[sender] = pulse
	newPulse := Low
	for _, p := range c.mrp {
		if p == Low {
			newPulse = High
			break
		}
	}
	for _, dest := range moduleOutputs[c.name] {
		messages.Enqueue(Message{c.name, newPulse, dest})
	}
}

type Broadcaster struct {
	name string
}

func (b *Broadcaster) Send(sender string, pulse Pulse) {
	for _, dest := range moduleOutputs[b.name] {
		messages.Enqueue(Message{b.name, pulse, dest})
	}
}

type Button struct {
	name string
}

func (b *Button) Send(sender string, pulse Pulse) {
	messages.Enqueue(Message{b.name, pulse, "broadcaster"})
}

func propagatePulses(path string, waitRxLow bool) int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	modules = map[string]Module{
		"button": &Button{name: "button"},
	}
	moduleInputs = map[string][]string{}
	moduleOutputs = map[string][]string{}

	modulesRaw := strings.Split(strings.Trim(string(file), "\n"), "\n")
	for _, line := range modulesRaw {
		split := strings.Split(line, " -> ")

		name := split[0]
		dests := strings.Split(split[1], ", ")

		if name == "broadcaster" {
			modules[name] = &Broadcaster{name: split[0]}
		} else {
			prefix := name[0]
			name = name[1:]
			switch prefix {
			case '%':
				modules[name] = &FlipFlop{name: name, powered: false}
			case '&':
				modules[name] = &Conjunction{name: name, mrp: map[string]Pulse{}}
			}
		}

		for _, dest := range dests {
			moduleInputs[dest] = append(moduleInputs[dest], name)
		}
		moduleOutputs[name] = dests
	}

	for _, module := range modules {
		if cj, ok := module.(*Conjunction); ok {
			for _, name := range moduleInputs[cj.name] {
				cj.mrp[name] = Low
			}
		}
	}

	i := 0

	if !waitRxLow {
		pulseCounter := map[Pulse]int{
			0: 0,
			1: 0,
		}

		for i < 1000 {
			messages = aq.New()
			messages.Enqueue(Message{"button", Low, "broadcaster"})

			for !messages.Empty() {
				qe, _ := messages.Dequeue()
				msg := qe.(Message)

				pulseCounter[msg.pulse]++

				// Prevent crashing on non-existent dummy output node
				if destModule := modules[msg.receiver]; destModule != nil {
					destModule.Send(msg.sender, msg.pulse)
				}
			}

			i++
		}

		return pulseCounter[Low] * pulseCounter[High]
	}

	// HARDCODED: Replace module name with corresponding one connecting to rx in personal input
	rxInputName := "kz"
	var rxInput *Conjunction
	firstHighPulses := map[string]int{}

	for _, module := range modules {
		if cj, ok := module.(*Conjunction); ok {
			if cj.name == rxInputName {
				rxInput = cj
			}
		}
	}

	for len(firstHighPulses) != len(rxInput.mrp) {
		messages = aq.New()
		messages.Enqueue(Message{"button", Low, "broadcaster"})

		for !messages.Empty() {
			qe, _ := messages.Dequeue()
			msg := qe.(Message)

			if msg.receiver == rxInputName && msg.pulse == High {
				firstHighPulses[msg.sender] = i + 1
			}

			// Prevent crashing on non-existent dummy output node
			if destModule := modules[msg.receiver]; destModule != nil {
				destModule.Send(msg.sender, msg.pulse)
			}
		}

		i++
	}

	minButtonPresses := 1
	for _, count := range firstHighPulses {
		minButtonPresses = utils.Lcm(minButtonPresses, count)
	}
	return minButtonPresses
}

func main() {
	totalPulses := propagatePulses("20/input.txt", false)
	minButtonPresses := propagatePulses("20/input.txt", true)
	fmt.Print("Part 1 solution: ", totalPulses, "\nPart 2 solution: ", minButtonPresses, "\n")
}
