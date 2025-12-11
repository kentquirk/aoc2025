package main

import (
	"fmt"
	"io"
	"log"
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

type machine struct {
	nbits    int
	lamps    bits
	switches []bits
	joltages []int
	state    bits
	children []*machine
}

func (m machine) String() string {
	return fmt.Sprintf("lamps: %s, switches: %v, joltages: %v", m.lamps.asBits(m.nbits), m.switches, m.joltages)
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
	}

	mj := joltsPat.FindStringSubmatch(line)
	if mj != nil {
		machine.joltages = parseJolts(mj[1])
	}
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
	return 0
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
	fmt.Println(part1(data))
}
