package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type shape struct {
	rows  [][]byte
	count int
}

type region struct {
	w      int
	h      int
	counts []int
}

func part1(shapes []shape, regions []region) int {
	possible := 0
	for i, r := range regions {
		areaOfRegion := r.w * r.h
		areaOfShapes := 0
		for j, sh := range shapes {
			areaOfShapes += sh.count * r.counts[j]
		}
		if areaOfShapes > areaOfRegion {
			fmt.Printf("Area of shapes %d exceeds area of region %d for region %d %v\n", areaOfShapes, areaOfRegion, i, r)
		} else {
			possible++
		}
	}
	return possible
}

func part2(shapes []shape, regions []region) int {
	return 0
}

func extractNumbers(s string) []int {
	numpat := regexp.MustCompile("[0-9]+")
	matches := numpat.FindAllString(s, -1)
	var numbers []int
	for _, m := range matches {
		n, err := strconv.Atoi(m)
		if err != nil {
			log.Fatal(err)
		}
		numbers = append(numbers, n)
	}
	return numbers
}

func parse(filename string) ([]shape, []region) {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	var shapes []shape
	// we know there are 6 shapes, each with 5 lines
	for i := 0; i < 6; i++ {
		sh := shape{}
		for j := 1; j <= 3; j++ {
			sh.rows = append(sh.rows, []byte(lines[i*5+j]))
			sh.count += bytes.Count(sh.rows[j-1], []byte{'#'})
		}
		shapes = append(shapes, sh)
	}

	var regions []region
	for _, l := range lines[30:] {
		nums := extractNumbers(l)
		region := region{
			w:      nums[0],
			h:      nums[1],
			counts: nums[2:],
		}
		regions = append(regions, region)
	}
	return shapes, regions
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
	shapes, regions := parse(filename)
	fmt.Println(part1(shapes, regions))
}
