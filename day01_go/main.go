package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func part1(lines []string) int {
	dial := 50
	count := 0

	for _, l := range lines {
		if l == "" {
			continue
		}
		s := l[0]
		v, _ := strconv.Atoi(l[1:])
		switch s {
		case 'L':
			dial = ((dial - v) % 100)
			for dial < 0 {
				dial += 100
			}
		case 'R':
			dial = (dial + v) % 100
		}

		if dial == 0 {
			count++
		}

	}

	return count
}

// cases:
//
//	Rv:
//	  all values: (dial+v)/100 clicks
//	Lv:
//	  dial-v > 0: 0 clicks
//	  dial-v = 0: 1 click
//	  dial-v < 0: -((dial-v-99)/100) clicks
func part2(lines []string) int {
	dial := 50
	zeroClicks := 0
	clicks := 0

	for _, l := range lines {
		if l == "" {
			continue
		}
		direction := l[0]
		value, _ := strconv.Atoi(l[1:])
		switch direction {
		case 'L':
			newdial := dial - value
			switch {
			case newdial > 0:
				// still positive, no clicks
				clicks = 0
				dial = newdial
			case newdial <= 0:
				clicks = -((newdial - 100) / 100)
				if dial == 0 {
					// if we start on zero, we overcounted by 1
					clicks--
				}
				// in go, modulus of negative numbers is negative, so adjust
				dial = newdial % 100
				if dial < 0 {
					dial += 100
				}
			}
		case 'R':
			newdial := dial + value
			// we just divide by 100 to get the number of clicks
			clicks = newdial / 100
			dial = newdial % 100
		}
		zeroClicks += clicks
		// fmt.Printf("%c%d %d %d %d\n", direction, value, dial, clicks, zeroClicks)
	}

	return zeroClicks
}

func readlines(filename string) []string {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(b), "\n")
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
	lines := readlines(filename)
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}
