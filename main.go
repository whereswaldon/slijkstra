package main

import (
	"bufio"
	"fmt"
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

func main() {
	defer handleExit()
	if len(os.Args) < 2 {
		fmt.Println("Please supply a filename as an argument")
		panic(exit{1})
	}

	fmt.Println("Opening", os.Args[1])
	file, err := os.Open(os.Args[1])
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
	fmt.Println(graph)
}
