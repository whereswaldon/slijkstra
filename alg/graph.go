package alg

import (
	"container/list"
	"fmt"
	"github.com/akonneker/golib/gopqueue"
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

// HasEdge returns true if the two given nodes are connected
func (g *Graph) HasEdge(u, v int) bool {
	list := g.AdjacencyLists[u]
	for cursor := list.Front(); cursor != nil; cursor = cursor.Next() {
		edge := cursor.Value.(*Edge)
		if (edge.V == v && edge.U == u) || (edge.V == u && edge.U == v) {
			return true
		}
	}
	return false
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

func (g *Graph) FindShortestPathTree(rootVertex int) *Table {
	// create a new empty table
	table := NewTable(g.Order, rootVertex)
	// create a new empty priority queue
	q := pqueue.New(g.Order)
	// Set the root vertex to have no parent and have distance 0
	table.Set(rootVertex, 0, -1)
	q.Enqueue(WeightedVertex{Vertex: rootVertex, Weight: 0})
	for !q.IsEmpty() {
		// Pull off top of queue and cast to concrete type
		current := q.Dequeue().(WeightedVertex).Vertex
		if table.Visited(current) {
			continue // skip evaluation
		}
		for edge := g.AdjacencyLists[current].Front(); edge != nil; edge = edge.Next() {
			concrete := edge.Value.(*Edge)
			// ensure that we don't assume which is U and which is V
			other := concrete.U
			if other == current {
				other = concrete.V
			}
			// skip this edge if we've been to the other end
			if table.Visited(other) {
				continue
			}
			// if it is unseen, give it the current distance
			if table.Distance(other) < 0 ||
				table.Distance(other) > table.Distance(current)+concrete.Weight {
				table.Set(other, table.Distance(current)+concrete.Weight, current)
			}
			q.Enqueue(WeightedVertex{Vertex: other, Weight: table.Distance(other)})
		}
		table.Visit(current)
	}
	return table
}

// Find diameter returns the starting point, ending point, and distance of the
// shortest-path that defines the diameter of the graph
func (g *Graph) FindDiameter() (startNode, endNode, distance int) {
	currentTable := g.FindShortestPathTree(0)
	startNode = currentTable.Root
	endNode = currentTable.FurthestNode
	distance = currentTable.MaxDistance
	for i := 1; i < g.Order; i++ {
		currentTable = g.FindShortestPathTree(i)
		if currentTable.MaxDistance > distance {
			startNode = currentTable.Root
			endNode = currentTable.FurthestNode
			distance = currentTable.MaxDistance
		}
	}
	return
}

type WeightedVertex struct {
	Vertex int
	Weight int
}

func (w WeightedVertex) Less(o interface{}) bool {
	return w.Weight < o.(WeightedVertex).Weight
}

// Table represents the progress of a run of the shortest path
// tree algorithm
type Table struct {
	MaxDistance, Root, FurthestNode int
	visited                         []bool
	distance                        []int
	parent                          []int
}

// NewTable creates a new empty table
func NewTable(order, root int) *Table {
	distance := make([]int, order)
	for i := range distance {
		distance[i] = -1
	}
	return &Table{
		Root:         root,
		FurthestNode: root,
		MaxDistance:  0,
		visited:      make([]bool, order),
		distance:     distance,
		parent:       make([]int, order),
	}
}

// Visited returns whether or not the given vertex has been marked
// as visited
func (t *Table) Visited(vertex int) bool {
	return t.visited[vertex]
}

// Visit marks a vertex as visited
func (t *Table) Visit(vertex int) {
	if t.distance[vertex] > t.MaxDistance {
		t.MaxDistance = t.distance[vertex]
		t.FurthestNode = vertex
	}
	t.visited[vertex] = true
}

// Distance returns the shortest distance known to the given vertex, -1
// if no shortest distance is known.
func (t *Table) Distance(vertex int) int {
	return t.distance[vertex]
}

// Parent returns the parent of the given vertex, -1 if the vertex is the
// root of the tree
func (t *Table) Parent(vertex int) int {
	return t.parent[vertex]
}

// Set updates the table entry for a vertex with a new shortest distance
// and parent node
func (t *Table) Set(vertex, distance, parent int) {
	t.distance[vertex] = distance
	t.parent[vertex] = parent
}

func (t *Table) String() string {
	res := fmt.Sprintf("%8s %8s %8s %8s\n", "Vertex", "Visited", "Distance", "Parent")
	for i := range t.visited {
		res += fmt.Sprintf("%8d %8t %8d %8d\n", i, t.visited[i], t.distance[i], t.parent[i])
	}
	return res
}
