package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"math/rand/v2"
	"os"
	"strings"
)

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type edge struct {
	start, end point
	horizontal bool
}

func newEdge(start, end point) edge {
	// Use closed intervals [start, end]
	return edge{start: start, end: end, horizontal: start.y == end.y}
}

func (e edge) isHorizontal() bool {
	return e.horizontal
}

func (e edge) isVertical() bool {
	return !e.horizontal
}

// on returns true if point p lies on this edge
func (e edge) on(p point) bool {
	if e.isVertical() {
		// Point is on vertical edge if x matches and y is in range [start.y, end.y]
		return p.x == e.start.x && between(p.y, e.start.y, e.end.y)
	} else {
		// Point is on horizontal edge if y matches and x is in range [start.x, end.x]
		return p.y == e.start.y && between(p.x, e.start.x, e.end.x)
	}
}

func between(v, start, end int) bool {
	if start > end {
		start, end = end, start
	}
	return (v >= start && v <= end)
}

// This function checks if two edges cross each other
// Edges are always either horizontal or vertical
// Other will always be horizontal (the ray)
func (e edge) crosses(ray edge) bool {
	if e.isVertical() {
		// Check if vertical edge crosses horizontal ray
		rayY := ray.start.y
		// Ray must be within the vertical edge's Y range (handle non-normalized edges)
		if between(rayY, e.start.y, e.end.y) {
			// Vertical edge must be to the left of ray start (since ray goes left)
			if e.start.x < ray.start.x {
				// Special handling for vertex intersections to avoid double counting
				if rayY == e.start.y || rayY == e.end.y {
					// Use "start vertex rule" - only count if ray hits the lower vertex
					minY := e.start.y
					if e.end.y < minY {
						minY = e.end.y
					}
					return rayY == minY
				}
				// Regular edge intersection
				return true
			}
		}
	} else {
		// Both horizontal - parallel rays don't cross horizontal edges
		return false
	}
	return false
}

type rectangle struct {
	min, max point
}

func (r rectangle) String() string {
	return fmt.Sprintf("[%v-%v]", r.min, r.max)
}

func newRectangle(p1, p2 point) rectangle {
	x1 := p1.x
	y1 := p1.y
	x2 := p2.x
	y2 := p2.y
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	return rectangle{min: point{x: x1, y: y1}, max: point{x: x2, y: y2}}
}

func (r rectangle) area() int {
	width := r.max.x - r.min.x + 1
	height := r.max.y - r.min.y + 1
	return width * height
}

func (r rectangle) randomPoints() iter.Seq[point] {
	// Sample up to 500 points within the rectangle or up to half the area.
	// This gives a good chance of finding an outside point quickly
	// without spending too much time on large rectangles.
	n := min(500, r.area()/2)
	return func(yield func(point) bool) {
		for i := 0; i < n; i++ {
			x := r.min.x + rand.IntN(r.max.x-r.min.x+1)
			y := r.min.y + rand.IntN(r.max.y-r.min.y+1)
			if !yield(point{x: x, y: y}) {
				return
			}
		}
	}
}

func (r rectangle) edgePoints() iter.Seq[point] {
	// If we test all the points in the edges of the rectangle, we can be sure
	// that if all edge points are inside the shape, then the entire rectangle
	// is inside the shape because the shape can't have holes.
	return func(yield func(point) bool) {
		// Top and bottom edges
		for x := r.min.x; x <= r.max.x; x++ {
			if !yield(point{x: x, y: r.min.y}) {
				return
			}
			if !yield(point{x: x, y: r.max.y}) {
				return
			}
		}
		// Left and right edges (skip corners to avoid duplicates)
		for y := r.min.y + 1; y < r.max.y; y++ {
			if !yield(point{x: r.min.x, y: y}) {
				return
			}
			if !yield(point{x: r.max.x, y: y}) {
				return
			}
		}
	}
}

type shape struct {
	edges []edge
}

func (s *shape) addEdge(start, end point) {
	s.edges = append(s.edges, newEdge(start, end))
}

func (s shape) isInside(p point) bool {
	// First check if point lies on any edge (boundary points are considered inside)
	for _, e := range s.edges {
		if e.on(p) {
			// fmt.Printf("Point %v is on edge %v\n", p, e)
			return true
		}
	}

	// Cast a ray to the left and count intersections with vertical edges
	// An odd count means we're inside
	count := 0
	ray := edge{start: p, end: point{x: -1, y: p.y}} // ray going left
	for _, e := range s.edges {
		if e.crosses(ray) {
			count++
			// fmt.Printf("Ray from %v crosses edge %v (start:%v end:%v)\n", p, e, e.start, e.end)
		}
	}
	// fmt.Printf("Point %v: %d crossings, inside: %v\n", p, count, count%2 == 1)
	return count%2 == 1
}

func part1(points []point) int {
	// we'll try brute force for now
	maxarea := 0
	for i, p1 := range points {
		for j, p2 := range points {
			if i == j {
				continue
			}
			x1 := p1.x
			y1 := p1.y
			x2 := p2.x
			y2 := p2.y
			if x2 < x1 {
				x1, x2 = x2, x1
			}
			if y2 < y1 {
				y1, y2 = y2, y1
			}
			area := (x2 - x1 + 1) * (y2 - y1 + 1)
			if area > maxarea {
				maxarea = area
			}
		}
	}
	return maxarea
}

func part2(points []point) int {
	shape := &shape{}
	for i := 1; i < len(points); i++ {
		p1 := points[i-1]
		p2 := points[i]
		shape.addEdge(p1, p2)
	}
	// Also add edge from last point back to first to close the shape
	if len(points) > 0 {
		shape.addEdge(points[len(points)-1], points[0])
	}

	// fmt.Printf("Shape has %d edges:\n", len(shape.edges))
	// for i, e := range shape.edges {
	// 	fmt.Printf("  Edge %d: %v to %v (horizontal: %v)\n", i, e.start, e.end, e.horizontal)
	// }
	largestArea := 0
	// now create each rectangle between each pair of points
	// and start generating random points within that rectangle
	// and see if they're inside the shape
	for i, p1 := range points {
	inner:
		for j := i + 1; j < len(points); j++ {
			p2 := points[j]
			rect := newRectangle(p1, p2)
			// look for an early out by testing a random sample of points
			for p := range rect.randomPoints() {
				if !shape.isInside(p) {
					// fmt.Printf("Skipping %v because point %v is outside\n", rect, p)
					continue inner
				}
			}
			// if any point is not inside the shape, skip this rectangle
			for p := range rect.edgePoints() {
				if !shape.isInside(p) {
					// fmt.Printf("Skipping %v because edge point %v is outside\n", rect, p)
					continue inner
				}
			}
			area := rect.area()
			if area > largestArea {
				largestArea = area
				fmt.Printf("New largest area: %d between points %v and %v\n", largestArea, p1, p2)
			}
		}
	}
	return largestArea
}

func readlines(filename string) []point {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	var points []point
	for _, line := range lines {
		if line == "" {
			continue
		}
		var p point
		_, err := fmt.Sscanf(line, "%d,%d", &p.x, &p.y)
		if err != nil {
			log.Fatal(err)
		}
		points = append(points, p)
	}
	return points
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
	points := readlines(filename)
	fmt.Println(part1(points))
	fmt.Println(part2(points))
}
