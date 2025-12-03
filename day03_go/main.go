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

// algorithm: find first max character in a string not including the last
// character -- and its position, which must be at least n digits from the end.
// Then find the next max character after that position. If we ever run out of
// string, just take whatever is left. Repeat until we have n digits.
func solve(lines []string, numdigits int) int {
	total := 0
	// BUG can't truncate to less than n digits
	for _, line := range lines {
		digits := ""
		for n := numdigits - 1; n >= 0; n-- {
			if len(line) <= n {
				digits += line
				break
			}
			sorted := strings.Split(line[:len(line)-n], "")
			sort.Strings(sorted)
			maxdigit := sorted[len(sorted)-1]
			p := strings.Index(line, maxdigit)
			line = line[p+1:]
			digits += maxdigit
		}
		value, _ := strconv.Atoi(digits)
		total += value
	}
	return total
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
	fmt.Println(solve(lines, 2))
	fmt.Println(solve(lines, 12))
}
