package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type bits int

func (b bits) String() string {
	return fmt.Sprintf("%03x", int(b))
}

func (b bits) asBits(n int) string {
	var result strings.Builder
	for i := 0; i < n; i++ {
		if b&(1<<i) != 0 {
			result.WriteByte('#')
		} else {
			result.WriteByte('.')
		}
	}
	return result.String()
}

// we call them switches in part 1 and buttons in part 2
type machine struct {
	nbits    int
	lamps    bits
	buttons  [][]int // each button is a list of indices
	switches []bits  // the bit patterns for each switch
	joltages []int
	state    bits
	children []*machine
	jstate   []int
	presses  []int
}

func (m machine) String() string {
	return fmt.Sprintf("lamps: %s, switches: %v, joltages: %v", m.lamps.asBits(m.nbits), m.switches, m.joltages)
}

func (m machine) JString() string {
	return fmt.Sprintf("jstate: %v, presses: %v", m.jstate, m.presses)
}

func (m machine) print() {
	fmt.Printf("Lamps - %s\n", m.lamps.asBits(m.nbits))
	for i, s := range m.switches {
		fmt.Printf(" Sw%02d   %s\n", i, s.asBits(m.nbits))
	}
	fmt.Printf(" J      {%v}\n", m.joltages)
}

func (m machine) pressSwitch(i int) (*machine, bool) {
	newState := m.state ^ m.switches[i]
	if newState == m.lamps {
		return nil, true
	}
	child := &machine{
		nbits:    m.nbits,
		lamps:    m.lamps,
		joltages: m.joltages,
		state:    newState,
	}
	// remove the switch we just used from the list
	for j := range m.switches {
		if j == i {
			continue
		}
		child.switches = append(child.switches, m.switches[j])
	}
	// sort the switches so we can hash them
	sort.Slice(child.switches, func(i, j int) bool {
		return child.switches[i] < child.switches[j]
	})
	return child, false
}

// This is a breadth-first search of the solution space.
// we push each of the switches and clone children, but
// stop if pushing any of them solves the machine.
// If we've tried all and we're not done, we go deeper by
// visiting each child.
func (m machine) search(depth int) int {
	// fmt.Printf("%*s at depth %d: %s\n", depth*2, "", depth, m)
	for i := range m.switches {
		n, done := m.pressSwitch(i)
		if done {
			return depth
		}
		m.children = append(m.children, n)
	}
	dmin := 1000
	for _, child := range m.children {
		if d := child.search(depth + 1); d != -1 {
			if d < dmin {
				dmin = d
			}
		}
	}
	if dmin != 1000 {
		return dmin
	}
	return -1
}

func (m *machine) pressButton() bool {
	// calculate how far jstate is from joltages
	deltas := make([]int, len(m.joltages))
	for i, st := range m.jstate {
		deltas[i] = m.joltages[i] - st
	}
	// find the best candidate button based on deltas
	// value of a button is sum of deltas corresponding
	// to the joltages that the button effects -- but if
	// a button has 0 for a delta, skip it.
	best := -1
	bestValue := -1
outer:
	for i, b := range m.buttons {
		value := 0
		for _, ix := range b {
			if deltas[ix] == 0 {
				continue outer
			}
			value += deltas[ix]
		}
		if value > bestValue {
			bestValue = value
			best = i
		}
	}
	if best == -1 {
		return false
	}

	m.presses = append(m.presses, best)
	for _, i := range m.buttons[best] {
		m.jstate[i]++
	}
	return true
}

func (m *machine) unpressButton() {
	ix := rand.IntN(len(m.presses))
	button := m.presses[ix]
	m.presses = append(m.presses[:ix], m.presses[ix+1:]...)
	for _, i := range m.buttons[button] {
		m.jstate[i]--
	}
}

type joltstate int

const (
	TooHigh joltstate = 1
	TooLow  joltstate = -1
	Equal   joltstate = 0
)

func (m *machine) check() joltstate {
	isEqual := true
	for i := range m.joltages {
		if m.jstate[i] > m.joltages[i] {
			return TooHigh
		}
		if m.jstate[i] < m.joltages[i] {
			isEqual = false
		}
	}
	if isEqual {
		return Equal
	}
	return TooLow
}

func (m *machine) reset() {
	m.jstate = make([]int, len(m.joltages))
	m.presses = nil
}

func parseLamps(s string) (int, bits) {
	// string is . (0) and # (1), and low order bit is first
	// so ".#.#" should be A and "#..###" should be 56
	var result int
	for i, c := range s {
		b := 0
		if c == '#' {
			b = 1
		}
		result = result | (b << i)
	}
	return len(s), bits(result)
}

func parseNumberList(s string) []int {
	// takes a comma-separated list of numbers and returns a slice of ints
	parts := strings.Split(s, ",")
	var result []int
	for _, p := range parts {
		num, err := strconv.Atoi(p)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, num)
	}
	return result
}

func parseSwitch(s string) bits {
	// takes a comma-separated list of positions and sets them to true in the result
	var result int
	values := parseNumberList(s)
	for _, pos := range values {
		result |= 1 << pos
	}
	return bits(result)
}

func parseButtons(s string) []int {
	// takes a comma-separated list of positions and sets them to true in the result
	values := parseNumberList(s)
	return values
}

func parseJolts(s string) []int {
	return parseNumberList(s)
}

var lampPat = regexp.MustCompile(`\[([^\]]+)\]`)
var switchPat = regexp.MustCompile(`\(([^\)]+)\)`)
var joltsPat = regexp.MustCompile(`\{([^\}]+)\}`)

func parseMachine(line string) machine {
	machine := machine{}
	ml := lampPat.FindStringSubmatch(line)
	if ml == nil {
		log.Fatalf("invalid format %s", line)
	}
	machine.nbits, machine.lamps = parseLamps(ml[1])

	ms := switchPat.FindAllStringSubmatch(line, -1)
	for _, sub := range ms {
		machine.switches = append(machine.switches, parseSwitch(sub[1]))
		machine.buttons = append(machine.buttons, parseButtons(sub[1]))
	}

	mj := joltsPat.FindStringSubmatch(line)
	if mj != nil {
		machine.joltages = parseJolts(mj[1])
	}
	machine.jstate = make([]int, len(machine.joltages))
	return machine
}

func parse(filename string) []machine {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")

	var machines []machine
	for _, line := range lines {
		if line == "" {
			continue
		}
		machine := parseMachine(line)
		// machine.print()
		machines = append(machines, machine)
	}

	return machines
}

func part1(data []machine) int {
	sum := 0
	for _, m := range data {
		if d := m.search(1); d != -1 {
			sum += d
			fmt.Printf("solved in %d steps at %s: %s\n", d, time.Now().Format(time.RFC3339), m)
		} else {
			fmt.Printf("failed to solve machine at %s: %s\n", time.Now().Format(time.RFC3339), m)
		}
	}
	return sum
}

func part2(data []machine) int {
	total := 0
	for _, m := range data {
		attempts := 0
		best := 10000
		score := 10000
		for range 1 {
			m.reset()
			done := false
			attempts++
			for !done {
				switch m.check() {
				case TooHigh:
					m.unpressButton()
				case TooLow:
					if !m.pressButton() {
						m.unpressButton()
					}
				case Equal:
					score = len(m.presses)
					done = true
				}
			}
			if score < best {
				best = score
				fmt.Printf("%s\n", m)
				fmt.Printf("Solved with %d presses in %d attempts: %s\n", score, attempts, m.JString())
			}
		}
		total += best
	}
	return total
}

func main() {
	args := os.Args[1:]
	filename := "sample"
	if len(args) > 0 {
		switch args[0] {
		case "-s":
			filename = "sample"
		case "-i":
			filename = "input"
		default:
			filename = args[0]
		}
	}
	data := parse(filename)
	// fmt.Println(part1(data))
	fmt.Println(part2(data))
}
