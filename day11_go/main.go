package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type node struct {
	name     string
	children []*node
}

func (n *node) String() string {
	s := fmt.Sprintf("%s: ", n.name)
	for _, child := range n.children {
		s += fmt.Sprintf("%s ", child.name)
	}
	return s
}

type graph struct {
	nodes map[string]*node
}

func (g graph) String() string {
	s := ""
	for _, node := range g.nodes {
		s += fmt.Sprintf("%s\n", node)
	}
	return s
}

type path []*node

func (p path) String() string {
	s := ""
	for _, n := range p {
		s += fmt.Sprintf("%s ", n.name)
	}
	return s
}

func traverse(n *node, pathIn path, successes *[]path) {
	for _, child := range n.children {
		newpath := append(path{}, pathIn...)
		newpath = append(newpath, child)
		if child.name == "out" {
			*successes = append(*successes, newpath)
			continue
		}
		traverse(child, newpath, successes)
	}
}

func part1(data graph) int {
	// fmt.Printf("%s", data)
	var successes []path
	traverse(data.nodes["you"], path{data.nodes["you"]}, &successes)
	// for _, p := range successes {
	// 	fmt.Println(p)
	// }
	return len(successes)
}

func part2(data graph) int {
	return 0
}

func parse(filename string) graph {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	g := graph{nodes: make(map[string]*node)}
	g.nodes["out"] = &node{name: "out"}
	// we do this in two passes -- first make the nodes, then populate the children
	for _, line := range lines {
		if line == "" {
			continue
		}
		name := line[:3] // we know it's always 3 letters
		node := &node{name: name}
		g.nodes[name] = node
	}
	// second pass -- populate the children
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line[5:])
		node := g.nodes[line[:3]]
		for _, childName := range parts {
			child, ok := g.nodes[childName]
			if !ok {
				log.Fatalf("child node %s not found", childName)
			}
			node.children = append(node.children, child)
		}
	}
	return g
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
