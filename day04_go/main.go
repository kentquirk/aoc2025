package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func countNeighbors(grid [][]byte, row, col int) int {
	count := 0
	directions := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1} /*******/, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for _, dir := range directions {
		r, c := row+dir[0], col+dir[1]
		if r >= 0 && r < len(grid) && c >= 0 && c < len(grid[row]) {
			if grid[r][c] == '@' {
				count++
			}
		}
	}
	return count
}

func countRolls(grid [][]byte) int {
	count := 0
	for r := range len(grid) {
		for c := range len(grid[r]) {
			if grid[r][c] == '@' {
				count++
			}
		}
	}
	return count
}

func part1(lines [][]byte) int {
	total := 0
	for r := range len(lines) {
		for c := range len(lines[r]) {
			neighbors := countNeighbors(lines, r, c)
			// fmt.Print(neighbors)
			if lines[r][c] == '@' && neighbors < 4 {
				total++
			}
		}
		// fmt.Println()
	}
	return total
}

func findRemoveables(lines [][]byte) [][2]int {
	var removeables [][2]int
	for r := range len(lines) {
		for c := range len(lines[r]) {
			neighbors := countNeighbors(lines, r, c)
			if lines[r][c] == '@' && neighbors < 4 {
				removeables = append(removeables, [2]int{r, c})
			}
		}
	}
	return removeables
}

func part2(lines [][]byte) int {
	rounds := 0
	rolls := countRolls(lines)
	fmt.Println("Initial rolls:", rolls)
	for {
		removeables := findRemoveables(lines)
		if len(removeables) == 0 {
			break
		}
		fmt.Println("Removing:", len(removeables))
		for _, rc := range removeables {
			lines[rc[0]][rc[1]] = '.'
		}
		rounds++
	}
	remaining := countRolls(lines)
	fmt.Println("Remaining rolls:", remaining)
	removed := rolls - remaining
	fmt.Println("Total removed:", removed)
	return removed
}

func readlines(filename string) [][]byte {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := bytes.Split(b, []byte("\n"))
	// assumes that there are no extra or funky characters in the input
	return lines
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
