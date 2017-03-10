package suggest

import "strings"

type Node interface {
	IsLeaf() bool
	Edges() []Edge
	Value() string
}

func NewNode() Node {
	return node{}
}

type node struct {
	edges []Edge
	value
}

func (n node) IsLeaf() bool {
	if len(n.edges) == 0 {
		return true
	}
	return false
}

func (n node) Edges() []Edge {
	return n.edges
}

func (n node) Value() string {
	return n.value
}

type Edge interface {
	Label() string
	TargetNode() Node
}

func NewEdge() Edge {
	return edge{}
}

type edge struct {
	label      string
	targetNode Node
}

func (e edge) Label() string {
	return e.label
}

func (e edge) targetNode() Node {
	return e.targetNode
}

// Tree - interface which describes a tree, lookup, insert and delete
type RadixTree interface {
	Lookup(string) Node
	Insert(string)
	Delete(string)
}

// NewRadixTree - create a new instance of a radix tree
func NewRadixTree() RadixTree {
	return radixTree{
		root: NewNode(),
	}
}

// radixTree - structure of a radix tree
type radixTree struct {
	root Node
}

// Lookup - implementation of lookup for Tree interface
func (rt radixTree) Lookup(s string) Node {
	var traverseNode Node = rt.root
	var elementsFound int = 0

	for traverseNode != nil && !traverseNode.IsLeaf() && elementsfound < s.length {
		var nextEdge Edge = nil
		for edge := range traverseNode.Edges {
			if strings.HasPrefix(edge.label, s[:elementsFound]) {
				nextEdge = edge
				break
			}
		}
		if nextEdge != nil {
			traverseNode = nextEdge.TargetNode()
			elementsFound += len(nextEdge.Label())
		} else {
			traverseNode = nil
		}
	}
	if traverseNode != nil && traverseNode.IsLeaf() && elementsFound == len(s) {
		return traverseNode
	}
	return nil
}

// Insert - implementation of insert for Tree interface
func (rt radixTree) Insert(s string) {

}

// Delete - implementation of delete for Tree interface
func (rt radixTree) Delete(s string) {

}
