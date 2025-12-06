package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type idRange struct {
	start int
	end   int
}

func size(r idRange) int {
	return r.end - r.start + 1
}

func (r idRange) contains(value int) bool {
	return r.start <= value && value <= r.end
}

func (r idRange) overlaps(other idRange) bool {
	return r.contains(other.start) || r.contains(other.end) ||
		other.contains(r.start) || other.contains(r.end)
}

func merge(a, b idRange) (idRange, bool) {
	if a.overlaps(b) {
		if b.start < a.start {
			a.start = b.start
		}
		if b.end > a.end {
			a.end = b.end
		}
		return a, true
	}
	return idRange{}, false
}

func part1(ranges []idRange, values []int) int {
	freshcount := 0

	for _, val := range values {
		for _, r := range ranges {
			if r.contains(val) {
				freshcount++
				break
			}
		}
	}
	return freshcount
}

// need to consolidate ranges, then measure the total size of
// all ranges
func part2(ranges []idRange, values []int) int {
	// sort ranges by start value
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	// now merge overlapping ranges in place by comparing sequentially
	last := len(ranges) - 1
	for i := 0; i < last; {
		merged, ok := merge(ranges[i], ranges[i+1])
		if ok {
			// replace ranges[i] with merged, and remove ranges[i+1]
			ranges[i] = merged
			ranges = append(ranges[:i+1], ranges[i+2:]...)
			last--
		} else {
			i++
		}
	}

	total := 0
	for _, r := range ranges {
		total += size(r)
	}
	return total
}

func readlines(filename string) ([]idRange, []int) {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	ranges := make([]idRange, 0)
	values := make([]int, 0)
	for _, line := range strings.Split(string(b), "\n") {
		if strings.Contains(line, "-") {
			// parse range
			vals := strings.Split(line, "-")
			if len(vals) != 2 {
				log.Fatalf("Invalid range: %s", line)
			}
			lo, _ := strconv.Atoi(vals[0])
			hi, _ := strconv.Atoi(vals[1])
			r := idRange{start: lo, end: hi}
			ranges = append(ranges, r)
		} else if line != "" {
			// parse single value
			val, _ := strconv.Atoi(line)
			values = append(values, val)
		}
	}
	return ranges, values
}

func main() {
	args := os.Args[1:]
	filename := "sample"
	if len(args) > 0 {
		switch args[0] {
		case "sample", "input":
			filename = args[0]
		case "-s":
			filename = "sample"
		case "-i":
			filename = "input"
		default:
			log.Fatalf("Unknown filename: %s", args[0])
		}
	}
	ranges, values := readlines(filename)
	fmt.Println(part1(ranges, values))
	fmt.Println(part2(ranges, values))
}
