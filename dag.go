// Package dag provides dag
package dag

import (
	"sort"
)

// Node is an integer that denotes Node.
type Node int

// Edge is a pair of depender node and dependee node.
type Edge struct {
	Depender Node
	Dependee Node
}

// Graph is SUPPOSED to be a DAG, but there is no guarantee.
// Graph MAY be disconnected graph.
type Graph struct {
	// Nodes SHOULD be sorted but not needed
	Nodes []Node
	// Edges SHOULD be sorted but not needed
	Edges []Edge
}

// VerifyAcyclicity returns true if the graph is DAG
func (g *Graph) VerifyAcyclicity() bool {
	panic("unimplemented")
	return true
}

// HasNode returns true if the graph has the node
func (g *Graph) HasNode(n Node) bool {
	for _, x := range g.Nodes {
		if x == n {
			return true
		}
	}
	return false
}

// AddNode adds a node if not existed
func (g *Graph) AddNode(n Node) {
	if !g.HasNode(n) {
		g.Nodes = append(g.Nodes, n)
		sort.Sort(NodesSorter(g.Nodes))
	}
}

// HasEdge returns true if the graph has the node
func (g *Graph) HasEdge(e Edge) bool {
	for _, x := range g.Edges {
		if x.Depender == e.Depender && x.Dependee == e.Dependee {
			return true
		}
	}
	return false
}

// AddEdge adds an edge if not existed
func (g *Graph) AddEdge(e Edge) {
	if !g.HasEdge(e) {
		g.Edges = append(g.Edges, e)
		sort.Sort(EdgesSorter(g.Edges))
		g.AddNode(e.Depender)
		g.AddNode(e.Dependee)
	}
}

// NodesSorter sorts nodes
type NodesSorter []Node

// Len implements sorter
func (x NodesSorter) Len() int { return len(x) }

// Swap implements sorter
func (x NodesSorter) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// Less implements sorter
func (x NodesSorter) Less(i, j int) bool { return x[i] < x[j] }

// EdgesSorter sorts edges
type EdgesSorter []Edge

// Len implements sorter
func (x EdgesSorter) Len() int { return len(x) }

// Swap implements sorter
func (x EdgesSorter) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// Less implements sorter
func (x EdgesSorter) Less(i, j int) bool {
	return x[i].Depender < x[j].Depender ||
		(x[i].Depender == x[j].Depender && x[i].Dependee < x[j].Dependee)
}

// ConnectedComponentRoots returns nodes with no dependee
// See also https://en.wikipedia.org/wiki/Connected_component_(graph_theory)
func (g *Graph) ConnectedComponentRoots() []Node {
	nonRoot := make(map[Node]struct{}, 0)
	for _, edge := range g.Edges {
		nonRoot[edge.Depender] = struct{}{}
	}
	var roots []Node
	for _, n := range g.Nodes {
		_, ok := nonRoot[n]
		if !ok {
			roots = append(roots, n)
		}
	}
	return roots
}

// DirectDependers returns direct dependers, not indirect ones
func (g *Graph) DirectDependers(dependee Node) []Node {
	var dependers []Node
	for _, e := range g.Edges {
		if e.Dependee == dependee {
			dependers = appendIfUnique(dependers, e.Depender)
		}
	}
	sort.Sort(NodesSorter(dependers))
	return dependers
}

// DirectDependees returns direct dependees, not indirect ones
func (g *Graph) DirectDependees(depender Node) []Node {
	var dependees []Node
	for _, e := range g.Edges {
		if e.Depender == depender {
			dependees = appendIfUnique(dependees, e.Dependee)
		}
	}
	sort.Sort(NodesSorter(dependees))
	return dependees
}

func appendIfUnique(l []Node, e Node) []Node {
	for _, f := range l {
		if f == e {
			return l
		}
	}
	return append(l, e)
}
