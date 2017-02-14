package main

import (
	"container/list"
	"fmt"
)

// Edge represents a weighted edge.
type Edge struct {
	U, V, Weight int
}

func (e *Edge) String() string {
	return fmt.Sprintf("Edge{U: %d, V:%d, Weight: %d,}", e.U, e.V, e.Weight)
}

// Graph represents a graph with a set
// of adjacency lists.
type Graph struct {
	Order          int
	AdjacencyLists []*list.List
}

// NewGraph creates an empty graph of the provided
// order.
func NewGraph(order int) *Graph {
	lists := make([]*list.List, order)
	for i := range lists {
		lists[i] = list.New()
	}
	return &Graph{
		Order:          order,
		AdjacencyLists: lists,
	}
}

// InsertEdge adds the provded edge to the adjacency lists of
// both of its end points.
func (g *Graph) InsertEdge(u, v, weight int) {
	edge := &Edge{U: u, V: v, Weight: weight}
	g.AdjacencyLists[u].PushBack(edge)
	g.AdjacencyLists[v].PushBack(edge)
}

func (g *Graph) String() string {
	res := fmt.Sprintf("Order: %d\n", g.Order)
	for i, l := range g.AdjacencyLists {
		res += fmt.Sprintf("\t%d: ", i)
		cursor := l.Front()
		for cursor != nil {
			res += fmt.Sprint(cursor.Value.(*Edge))
			cursor = cursor.Next()
		}
		res += "\n"
	}
	return res
}
