package graph

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
	"gopkg.in/yaml.v3"
)

// GraphNode is a node in an implicit graph.
type GraphNode struct {
	id        int64
	neighbors []graph.Node
	roots     []*GraphNode
	node      *yaml.Node
}

// NewGraphNode returns a new GraphNode.
func NewGraphNode(id int64) *GraphNode {
	return &GraphNode{id: id}
}

// Node allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) Node(id int64) graph.Node {
	if id == g.id {
		return g
	}

	seen := map[int64]struct{}{g.id: {}}
	for _, root := range g.roots {
		if root.ID() == id {
			return root
		}

		if root.has(seen, id) {
			return root
		}
	}

	for _, n := range g.neighbors {
		if n.ID() == id {
			return n
		}

		if gn, ok := n.(*GraphNode); ok {
			if gn.has(seen, id) {
				return gn
			}
		}
	}

	return nil
}

func (g *GraphNode) has(seen map[int64]struct{}, id int64) bool {
	for _, root := range g.roots {
		if _, ok := seen[root.ID()]; ok {
			continue
		}

		seen[root.ID()] = struct{}{}
		if root.ID() == id {
			return true
		}

		if root.has(seen, id) {
			return true
		}

	}

	for _, n := range g.neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}

		seen[n.ID()] = struct{}{}
		if n.ID() == id {
			return true
		}

		if gn, ok := n.(*GraphNode); ok {
			if gn.has(seen, id) {
				return true
			}
		}
	}

	return false
}

// Nodes allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) Nodes() graph.Nodes {
	nodes := []graph.Node{g}
	seen := map[int64]struct{}{g.id: {}}

	for _, root := range g.roots {
		nodes = append(nodes, root)
		seen[root.ID()] = struct{}{}

		nodes = root.nodes(nodes, seen)
	}

	for _, n := range g.neighbors {
		nodes = append(nodes, n)
		seen[n.ID()] = struct{}{}

		if gn, ok := n.(*GraphNode); ok {
			nodes = gn.nodes(nodes, seen)
		}
	}

	return iterator.NewOrderedNodes(nodes)
}

func (g *GraphNode) nodes(dst []graph.Node, seen map[int64]struct{}) []graph.Node {
	for _, root := range g.roots {
		if _, ok := seen[root.ID()]; ok {
			continue
		}
		seen[root.ID()] = struct{}{}
		dst = append(dst, graph.Node(root))

		dst = root.nodes(dst, seen)
	}

	for _, n := range g.neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}

		dst = append(dst, n)
		if gn, ok := n.(*GraphNode); ok {
			dst = gn.nodes(dst, seen)
		}
	}

	return dst
}

// From allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) From(id int64) graph.Nodes {
	if id == g.ID() {
		return iterator.NewOrderedNodes(g.neighbors)
	}

	seen := map[int64]struct{}{g.id: {}}
	for _, root := range g.roots {
		seen[root.ID()] = struct{}{}

		if result := root.findNeighbors(id, seen); result != nil {
			return iterator.NewOrderedNodes(result)
		}
	}

	for _, n := range g.neighbors {
		seen[n.ID()] = struct{}{}

		if gn, ok := n.(*GraphNode); ok {
			if result := gn.findNeighbors(id, seen); result != nil {
				return iterator.NewOrderedNodes(result)
			}
		}
	}

	return nil
}

func (g *GraphNode) findNeighbors(id int64, seen map[int64]struct{}) []graph.Node {
	if id == g.ID() {
		return g.neighbors
	}

	for _, root := range g.roots {
		if _, ok := seen[root.ID()]; ok {
			continue
		}
		seen[root.ID()] = struct{}{}

		if result := root.findNeighbors(id, seen); result != nil {
			return result
		}
	}

	for _, n := range g.neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}
		seen[n.ID()] = struct{}{}

		if gn, ok := n.(*GraphNode); ok {
			if result := gn.findNeighbors(id, seen); result != nil {
				return result
			}
		}
	}

	return nil
}

// HasEdgeBetween allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) HasEdgeBetween(uid, vid int64) bool {
	return g.EdgeBetween(uid, vid) != nil
}

// Edge allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) Edge(uid, vid int64) graph.Edge {
	return g.EdgeBetween(uid, vid)
}

// EdgeBetween allows GraphNode to satisfy the graph.Graph interface.
func (g *GraphNode) EdgeBetween(uid, vid int64) graph.Edge {
	if uid == g.id || vid == g.id {
		for _, n := range g.neighbors {
			if n.ID() == uid || n.ID() == vid {
				return simple.Edge{F: g, T: n}
			}
		}
		return nil
	}

	seen := map[int64]struct{}{g.id: {}}
	for _, root := range g.roots {
		seen[root.ID()] = struct{}{}
		if result := root.edgeBetween(uid, vid, seen); result != nil {
			return result
		}
	}

	for _, n := range g.neighbors {
		seen[n.ID()] = struct{}{}
		if gn, ok := n.(*GraphNode); ok {
			if result := gn.edgeBetween(uid, vid, seen); result != nil {
				return result
			}
		}
	}

	return nil
}

func (g *GraphNode) edgeBetween(uid, vid int64, seen map[int64]struct{}) graph.Edge {
	if uid == g.id || vid == g.id {
		for _, n := range g.neighbors {
			if n.ID() == uid || n.ID() == vid {
				return simple.Edge{F: g, T: n}
			}
		}
		return nil
	}

	for _, root := range g.roots {
		if _, ok := seen[root.ID()]; ok {
			continue
		}
		seen[root.ID()] = struct{}{}
		if result := root.edgeBetween(uid, vid, seen); result != nil {
			return result
		}
	}

	for _, n := range g.neighbors {
		if _, ok := seen[n.ID()]; ok {
			continue
		}

		seen[n.ID()] = struct{}{}
		if gn, ok := n.(*GraphNode); ok {
			if result := gn.edgeBetween(uid, vid, seen); result != nil {
				return result
			}
		}
	}

	return nil
}

// ID allows GraphNode to satisfy the graph.Node interface.
func (g *GraphNode) ID() int64 {
	return g.id
}

// AddMeighbor adds an edge between g and n.
func (g *GraphNode) AddNeighbor(n *GraphNode) {
	g.neighbors = append(g.neighbors, graph.Node(n))
}

// AddRoot adds provides an entrance into the graph g from n.
func (g *GraphNode) AddRoot(n *GraphNode) {
	g.roots = append(g.roots, n)
}
