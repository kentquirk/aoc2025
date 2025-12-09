package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"sort"
	"strings"
)

type point3 struct {
	x, y, z int
	circuit *circuit // pointer to the circuit this point belongs to
}

func (p *point3) linearDist2(o *point3) int {
	dx := p.x - o.x
	dy := p.y - o.y
	dz := p.z - o.z
	return dx*dx + dy*dy + dz*dz
}

type closePair struct {
	a, b  *point3
	dist2 int
}

func (cp *closePair) String() string {
	return fmt.Sprintf("(%d,%d,%d) <-> (%d,%d,%d) = %d",
		cp.a.x, cp.a.y, cp.a.z,
		cp.b.x, cp.b.y, cp.b.z,
		cp.dist2)
}

type closestPairs struct {
	pairs   []*closePair
	longest int
}

func NewClosestPairs() *closestPairs {
	return &closestPairs{
		pairs:   make([]*closePair, 0),
		longest: math.MaxInt,
	}
}

func (cp *closestPairs) addPair(a, b *point3) {
	dist2 := a.linearDist2(b)
	if dist2 < cp.longest {
		cp.pairs = append(cp.pairs, &closePair{a: a, b: b, dist2: dist2})
	}
}

type circuit []*point3

func (n *circuit) size() int {
	return len(*n)
}

func (n *circuit) has(p *point3) bool {
	return p.circuit != nil && p.circuit == n
}

func (n *circuit) add(p *point3) {
	if !n.has(p) {
		*n = append(*n, p)
		p.circuit = n
	}
}

func (n *circuit) join(other *circuit) {
	for _, p := range *other {
		if !n.has(p) {
			n.add(p)
		}
	}
}

func (n *circuit) String() string {
	s := "[ "
	for _, p := range *n {
		s += fmt.Sprintf("(%d,%d,%d) ", p.x, p.y, p.z)
	}
	s += "]"
	return s
}

func part1(points []*point3, numConnections int) int {
	cpairs := NewClosestPairs()
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			cpairs.addPair(points[i], points[j])
		}
		// at the end of each i loop, we can sort and trim cpairs.pairs to only
		// keep the closest numConnections pairs
		// Sort pairs by distance
		sort.Slice(cpairs.pairs, func(a, b int) bool {
			return cpairs.pairs[a].dist2 < cpairs.pairs[b].dist2
		})
		if len(cpairs.pairs) > numConnections {
			// Keep only the closest numConnections pairs
			if len(cpairs.pairs) > numConnections {
				cpairs.pairs = cpairs.pairs[:numConnections]
			}
		}
		// Update the longest distance if necessary (it's sorted
		// so the last element is the longest)
		if len(cpairs.pairs) > 0 {
			cpairs.longest = cpairs.pairs[len(cpairs.pairs)-1].dist2
		}
	}

	// now we have the closest pairs, we can add them all to circuits
	circuits := make([]*circuit, 0)
	for _, pair := range cpairs.pairs {
		// keep track of the most recently updated circuit
		var lastUpdated *circuit
		switch {
		case pair.a.circuit != nil && pair.b.circuit != nil:
			// if the're already in the same circuit, do nothing
			if pair.a.circuit == pair.b.circuit {
				continue
			}
			// both points are already in circuits, we need to join them
			if pair.a.circuit != pair.b.circuit {
				// find the index of pair.b.circuit
				ix := slices.Index(circuits, pair.b.circuit)
				if ix == -1 {
					fmt.Printf("circuits: %v\n", circuits)
					fmt.Printf("pair.a.circuit: %v\n", pair.a.circuit)
					fmt.Printf("pair.b.circuit: %v\n", pair.b.circuit)
					log.Fatalf("circuit %v not found in circuits slice %v", pair.b.circuit, circuits)
				}
				// join the two circuits
				pair.a.circuit.join(pair.b.circuit)
				// remove pair.b.circuit from circuits slice
				circuits = append(circuits[:ix], circuits[ix+1:]...)
				lastUpdated = pair.a.circuit
			}
		case pair.a.circuit != nil:
			pair.a.circuit.add(pair.b)
			lastUpdated = pair.a.circuit
		case pair.b.circuit != nil:
			pair.b.circuit.add(pair.a)
			lastUpdated = pair.b.circuit
		case pair.a.circuit == nil && pair.b.circuit == nil:
			// neither point is in a circuit yet
			newCircuit := &circuit{pair.a, pair.b}
			pair.a.circuit = newCircuit
			pair.b.circuit = newCircuit
			circuits = append(circuits, newCircuit)
			lastUpdated = newCircuit
		}

		if lastUpdated != nil && lastUpdated.size() == len(points) {
			// all points are now connected, we can stop early
			// return value is the product of the x values of the pair we just added
			product := pair.a.x * pair.b.x
			return product
		}
	}

	// now we have circuits, we can sort them by size
	sort.Slice(circuits, func(a, b int) bool {
		return circuits[a].size() > circuits[b].size()
	})

	// the product of the top 3 is the answer
	product := 1
	for i := 0; i < 3 && i < len(circuits); i++ {
		product *= circuits[i].size()
	}
	return product
}

func part2(points []point3) int {
	return 0
}

func readlines(filename string) []*point3 {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	var points []*point3
	for _, line := range lines {
		if line == "" {
			continue
		}
		var p point3
		_, err := fmt.Sscanf(line, "%d,%d,%d", &p.x, &p.y, &p.z)
		if err != nil {
			log.Fatalf("Failed to parse line '%s': %v", line, err)
		}
		points = append(points, &p)
	}
	return points
}

func main() {
	args := os.Args[1:]
	filename := "sample"
	numConnections := 10
	if len(args) > 0 {
		switch args[0] {
		case "sample", "input":
			filename = args[0]
		case "-s":
			filename = "sample"
			numConnections = 10
		case "-i":
			filename = "input"
			numConnections = 1000
		default:
			log.Fatalf("Unknown filename: %s", args[0])
		}
	}
	lines := readlines(filename)
	fmt.Println(part1(lines, numConnections))
	lines = readlines(filename)
	fmt.Println(part1(lines, 10000))
}
