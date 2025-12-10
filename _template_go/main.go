package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func part1(data []string) int {
	return 0
}

func part2(data []string) int {
	return 0
}

func parse(filename string) []string {
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
