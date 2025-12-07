package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func part1(lines [][]byte) int {
	splitCount := 0
	startPos := bytes.Index(lines[0], []byte("S"))
	beams := make([]map[int]struct{}, len(lines))
	beams[0] = map[int]struct{}{startPos: {}}
	for i := 1; i < len(lines)-1; i++ {
		// copy the previous row's beams or split them
		beams[i] = make(map[int]struct{})
		for b := range beams[i-1] {
			if lines[i][b] == '^' {
				if b > 0 {
					beams[i][b-1] = struct{}{}
				}
				if b < len(lines[i])-1 {
					beams[i][b+1] = struct{}{}
				}
				splitCount++
			} else {
				beams[i][b] = struct{}{}
			}
		}
	}
	return splitCount
}

// memoize the recursive calls
var memo map[[2]int]int

func doBeam(lines [][]byte, row, col int) int {
	if val, ok := memo[[2]int{row, col}]; ok {
		return val
	}
	// fmt.Printf("At row %d col %d char %c\n", row, col, lines[row][col])
	if row >= len(lines)-1 {
		return 1 // each time we reach the bottom, count 1 path
	}
	pathCount := 0
	if lines[row][col] == '^' {
		if col > 0 {
			pathCount += doBeam(lines, row+1, col-1)
		}
		if col < len(lines[row])-1 {
			pathCount += doBeam(lines, row+1, col+1)
		}
	} else {
		pathCount += doBeam(lines, row+1, col)
	}
	memo[[2]int{row, col}] = pathCount
	return pathCount
}

func part2(lines [][]byte) int {
	memo = make(map[[2]int]int)
	splitCount := 0
	startPos := bytes.Index(lines[0], []byte("S"))
	splitCount += doBeam(lines, 1, startPos)
	return splitCount
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
	return bytes.Split(b, []byte("\n"))
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
