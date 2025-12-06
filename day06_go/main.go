package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func part1(lines [][]string) int {
	grandtotal := 0
	ops := lines[len(lines)-1]
	lines = lines[:len(lines)-1]
	for i := range len(ops) {
		op := ops[i]
		switch op {
		case "+":
			total := 0
			for _, line := range lines {
				num := 0
				fmt.Sscanf(line[i], "%d", &num)
				total += num
			}
			grandtotal += total
		case "*":
			product := 1
			for _, line := range lines {
				num := 0
				fmt.Sscanf(line[i], "%d", &num)
				product *= num
			}
			grandtotal += product
		default:
			log.Fatalf("Unknown op: %s", op)
		}
	}
	return grandtotal
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
	lines := strings.Split(string(b), "\n")
	return lines
}

func splitLinesByBlanks(lines []string) [][]string {
	data := make([][]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		row := regexp.MustCompile(`\s+`).Split(line, -1)
		data = append(data, row)
	}
	return data
}

func rotateLines(lines []string) []string {
	maxlen := 0
	for _, line := range lines {
		if len(line) > maxlen {
			maxlen = len(line)
		}
	}
	// Pad lines to maxlen with spaces
	for i, line := range lines {
		if len(line) < maxlen {
			lines[i] = line + strings.Repeat(" ", maxlen-len(line))
		}
	}
	rotated := make([]string, 0)
	for i := maxlen - 1; i >= 0; i-- {
		s := ""
		for _, line := range lines {
			s += string(line[i])
		}
		rotated = append(rotated, s)
	}
	return rotated
}

func part2(lines []string) int {
	linedata := rotateLines(lines)
	grandtotal := 0
	values := []int{}
	pat := regexp.MustCompile(`\s*(\d+)\s*([+*]?)`)
	for lin := range linedata {
		matches := pat.FindStringSubmatch(linedata[lin])
		if matches == nil {
			continue
		}
		op := matches[2]
		num := 0
		if op == "" {
			fmt.Sscanf(matches[1], "%d", &num)
			values = append(values, num)
		} else {
			fmt.Sscanf(matches[1], "%d", &num)
			values = append(values, num)
			result := 0
			// fmt.Println(values, op)
			for _, v := range values {
				switch op {
				case "+":
					result += v
				case "*":
					if result == 0 {
						result = 1
					}
					result *= v
				default:
					log.Fatalf("Unknown op: '%s'", op)
				}
			}
			grandtotal += result
			values = []int{}
		}
	}

	return grandtotal
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
	linesData := splitLinesByBlanks(lines)
	fmt.Println(part1(linesData))
	fmt.Println(part2(lines))
}
