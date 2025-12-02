package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type idRange struct {
	lo       int
	hi       int
	loPrefix int
	hiPrefix int
}

func (r idRange) inRange(val int) bool {
	return val >= r.lo && val <= r.hi
}

func (r idRange) String() string {
	return fmt.Sprintf("%d-%d (prefix %d-%d), size %d\n", r.lo, r.hi, r.loPrefix, r.hiPrefix, r.hi-r.lo+1)
}

func isSequence(val int) bool {
	s := strconv.Itoa(val)
	for sequenceLength := 1; sequenceLength <= len(s)/2; sequenceLength++ {
		sequence := s[:sequenceLength]
		matched := true
		for i := 0; i < len(s); i += sequenceLength {
			end := i + sequenceLength
			if end > len(s) {
				end = len(s)
			}
			if s[i:end] != sequence {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

func part1(ranges []idRange) int {
	total := 0
	fmt.Println(ranges)
	for _, r := range ranges {
		for v := r.loPrefix; v <= r.hiPrefix; v++ {
			vs := strconv.Itoa(v)
			vv, _ := strconv.Atoi(vs + vs)
			if r.inRange(vv) {
				// fmt.Println("Found:", vv)
				total += vv
			}
		}
	}

	return total
}

func part2(ranges []idRange) int {
	total := 0
	for _, r := range ranges {
		for v := r.lo; v <= r.hi; v++ {
			if isSequence(v) {
				// fmt.Println("Found:", v)
				total += v
			}
		}
	}

	return total
}

func isEven(n int) bool {
	return n%2 == 0
}

func readRanges(filename string) []idRange {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	pairs := strings.Split(string(b), ",")
	var ranges []idRange
	for _, pair := range pairs {
		var r idRange
		pat := regexp.MustCompile(`(\d+)-(\d+)`)
		matches := pat.FindStringSubmatch(pair)
		lo, _ := strconv.Atoi(matches[1])
		hi, _ := strconv.Atoi(matches[2])
		r.lo = lo
		r.hi = hi
		loLen := len(matches[1])
		hiLen := len(matches[2])
		var loPrefix, hiPrefix int
		switch {
		case !isEven(loLen) && isEven(hiLen):
			loPrefix, _ = strconv.Atoi(matches[1][:loLen/2])
			hiPrefix, _ = strconv.Atoi(matches[2][:hiLen/2])
		default:
			loPrefix, _ = strconv.Atoi(matches[1][:loLen-loLen/2])
			hiPrefix, _ = strconv.Atoi(matches[2][:hiLen-hiLen/2])
		}
		r.loPrefix = loPrefix
		r.hiPrefix = hiPrefix
		ranges = append(ranges, r)
	}
	return ranges
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
	ranges := readRanges(filename)
	fmt.Println(part1(ranges))
	fmt.Println(part2(ranges))
}
