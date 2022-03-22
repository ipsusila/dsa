package dsa

import (
	"fmt"
	"io"
	"sort"

	"golang.org/x/exp/constraints"
)

// SortOrder types
type SortOrder int

// FilterFn filter node function
type FilterFn[T any] func(v T) bool

// TransformFn node function
type TransformFn[T any] func(v T) T

// SearchFn search node function
type SearchFn[T NodeType] func(nodes Nodes[T], value T) (*Node[T], int)

// List of sort order
const (
	Unsorted SortOrder = iota
	Ascending
	Descending
)

// NodeType constraint having comparable and Ordered constraint
type NodeType interface {
	comparable
	constraints.Ordered
}

// Nodes list of tree nodes
type Nodes[T NodeType] []*Node[T]

// Node stores tree node information
type Node[T NodeType] struct {
	leaf  bool
	Value T
	Nodes Nodes[T]
}

// Tree structure
type Tree[T NodeType] struct {
	so    SortOrder
	root  *Node[T]
	name  string
	fnSrc SearchFn[T]
}

// NewTree create tree structure with given name
func NewTree[T NodeType](name string) *Tree[T] {
	t := Tree[T]{
		so:   Unsorted,
		root: new(Node[T]),
		name: name,
	}
	t.fnSrc = t.searchSeq
	return &t
}

// Leaf return true if this node doesn't have any children
func (n Node[T]) Leaf() (bool, int) {
	return n.leaf, len(n.Nodes)
}

// Sort nodes in ASCENDING / DESCENDING order.
// If parameter is not specified, nodes will be sorted in ASCENDING order.
func (n Nodes[T]) Sort(order ...SortOrder) {
	if len(order) > 0 && order[0] == Descending {
		sort.Slice(n, func(i, j int) bool {
			return n[i].Value > n[j].Value
		})
	} else {
		sort.Slice(n, func(i, j int) bool {
			return n[i].Value < n[j].Value
		})
	}
}

// sequentially find node with value given as arg
func (n Nodes[T]) find(value T) (*Node[T], int) {
	for idx, v := range n {
		if v.Value == value {
			return v, idx
		}
	}
	return nil, -1
}

// debug purpose
func (n Nodes[T]) printTo(w io.Writer) {
	if len(n) > 0 {
		for i, v := range n {
			if i == 0 {
				fmt.Fprint(w, "[")
			} else {
				fmt.Fprint(w, " ")
			}
			fmt.Fprint(w, v.Value)
			if v.leaf {
				if len(v.Nodes) == 0 {
					fmt.Fprint(w, "#")
				} else {
					fmt.Fprint(w, "*")
				}
			}
		}
		fmt.Fprint(w, "] ")
	}
}

// Sorted return true if nodes are sorted either Ascending/Descending
func (t *Tree[T]) Sorted() bool {
	return t.so == Ascending || t.so == Descending
}

// SortOrder return nodes sorting order
func (t *Tree[T]) SortOrder() SortOrder {
	return t.so
}

// Empty return true if tree has no node
func (t *Tree[T]) Empty() bool {
	return t.root == nil || len(t.root.Nodes) == 0
}

// Name of the tree
func (t *Tree[T]) Name() string {
	return t.name
}

// Sort nodes in ASCENDING / DESCENDING order.
// If parameter is not specified, nodes will be sorted in ASCENDING order.
func (t *Tree[T]) Sort(order ...SortOrder) {
	if len(order) > 0 && order[0] == Descending {
		t.so = Descending
		t.fnSrc = t.searchDesc
		t.sortNodesDesc(t.root.Nodes)
	} else {
		t.so = Ascending
		t.fnSrc = t.searchAsc
		t.sortNodesAsc(t.root.Nodes)
	}
}
func (t *Tree[T]) sortNodesAsc(nodes Nodes[T]) {
	nodes.Sort(Ascending)
	for _, nd := range nodes {
		t.sortNodesAsc(nd.Nodes)
	}
}
func (t *Tree[T]) sortNodesDesc(nodes Nodes[T]) {
	nodes.Sort(Descending)
	for _, nd := range nodes {
		t.sortNodesDesc(nd.Nodes)
	}
}

func (t *Tree[T]) searchSeq(nodes Nodes[T], value T) (*Node[T], int) {
	for idx, v := range nodes {
		if v.Value == value {
			return v, idx
		}
	}
	return nil, -1
}

// Search specific value from list of nodes.
// Nodes must be *sorted* in ASCENDING order
func (t *Tree[T]) searchAsc(nodes Nodes[T], value T) (*Node[T], int) {
	idx := sort.Search(len(nodes), func(i int) bool { return nodes[i].Value >= value })
	if idx < len(nodes) && nodes[idx].Value == value {
		return nodes[idx], idx
	}
	return nil, -1
}

// Search specific value from list of nodes.
// Nodes must be *sorted* in DESCENDING order
func (t *Tree[T]) searchDesc(nodes Nodes[T], value T) (*Node[T], int) {
	idx := sort.Search(len(nodes), func(i int) bool { return nodes[i].Value <= value })
	if idx < len(nodes) && nodes[idx].Value == value {
		return nodes[idx], idx
	}
	return nil, -1
}

// Visit all values and return last visited nodes
// or nil if value not found in tree.
func (t *Tree[T]) Visit(values []T) *Node[T] {
	if t.Empty() {
		return nil
	}

	// traverse node
	var node *Node[T]
	targetNodes := t.root.Nodes
	for _, v := range values {
		if len(targetNodes) == 0 {
			return nil
		}

		// search within target nodes
		node, _ = t.fnSrc(targetNodes, v)
		if node == nil {
			return nil
		}
		targetNodes = node.Nodes
	}
	return node
}

// Match partially or all values
func (t *Tree[T]) Match(values []T) (*Node[T], bool) {
	n := t.Visit(values)
	return n, n != nil
}

// Match all values
func (t *Tree[T]) ExactMatch(values []T) (*Node[T], bool) {
	n := t.Visit(values)
	return n, n != nil && n.leaf
}

func (t *Tree[T]) insertNode(parent *Node[T], values []T) {
	// insert at the rest of nodes
	for _, vi := range values {
		nd, _ := parent.Nodes.find(vi)
		if nd == nil {
			nd = &Node[T]{
				Value: vi,
			}
			parent.Nodes = append(parent.Nodes, nd)
		}
		parent = nd
	}
	parent.leaf = true
}

// Insert values to tree
func (t *Tree[T]) Insert(values []T) {
	t.insertNode(t.root, values)
}

func (t *Tree[T]) printTo(w io.Writer, nodes Nodes[T], indent string) {
	if len(nodes) != 0 {
		fmt.Fprint(w, indent)
		nodes.printTo(w)
		fmt.Fprintln(w)
		for _, n := range nodes {
			t.printTo(w, n.Nodes, indent+"+")
		}
	}
}

// PrintTo for debugging purpose
func (t *Tree[T]) PrintTo(w io.Writer) {
	if !t.Empty() {
		t.printTo(w, t.root.Nodes, "")
	}
}
