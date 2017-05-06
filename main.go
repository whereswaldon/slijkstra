package main

import (
	"bufio"
	"fmt"
	. "github.com/whereswaldon/dijkstra/alg"
	"io"
	"os"
	"strconv"
	"strings"
)

// exit wraps an exit code with a type
type exit struct {
	Code int
}

// handleExit recovers from panics and checks whether the
// panic was caused by an exit request. If it was, it exits
// with that status code. Otherwise, it panics again.
func handleExit() {
	if r := recover(); r != nil {
		if e, ok := r.(exit); !ok {
			panic(r)
		} else {
			os.Exit(e.Code)
		}
	}
}

// setupGraph parses the file contents and creates a graph
// representing the file
func setupGraph(filename string) *Graph {
	fmt.Println("Opening", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to open file! ", err)
		panic(exit{1})
	}
	defer file.Close()
	buffReader := bufio.NewReader(file)
	nodes, err := buffReader.ReadString('\n')
	if err != nil {
		fmt.Println("Unable to read line! ", err)
		panic(exit{1})
	}
	nodes = strings.TrimSpace(nodes)
	order, err := strconv.Atoi(nodes)
	if err != nil {
		fmt.Println("Unable to convert nodes to integer")
		panic(exit{1})
	}
	graph := NewGraph(order)
	fmt.Println("There are", nodes, "nodes in the graph!")
	var u, v, weight int
	var line string
	for err == nil || err != io.EOF {
		line, err = buffReader.ReadString('\n')
		fmt.Sscanf(line, "%d %d %d\n", &u, &v, &weight)
		graph.InsertEdge(u, v, weight)
	}
	return graph
}

func main() {
	defer handleExit()
	var root int = 0
	var err error
	if len(os.Args) < 2 {
		fmt.Println("Please supply a filename as an argument")
		panic(exit{1})
	} else if len(os.Args) < 3 {
		fmt.Println("No root node specified, defaulting to zero")
	} else {
		root, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Unable to convert second argument to number: ", err)
			panic(exit{1})
		}
	}

	graph := setupGraph(os.Args[1])
	table := graph.FindShortestPathTree(root)

	fmt.Println(table)
	fmt.Println("Max distance: ", table.MaxDistance, " from ", table.Root, " to ", table.FurthestNode)
	s, e, d := graph.FindDiameter()
	fmt.Println("diameter: ", d, " start: ", s, " end: ", e)
}
