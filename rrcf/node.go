package rrcf

import "github.com/andysgithub/go-rrcf/num"

// Node of RCTree consisting of a leaf or branch and containing at most one parent
type Node struct {
	Leaf   *Leaf
	Branch *Branch
	b      [][]float64 // Bounding box of single point or points under branch
	u      *Node       // Pointer to parent
	n      int         // Number of leaves under branch or points in leaf
}

// Branch of RCTree containing two children
type Branch struct {
	q int     // Dimension of cut
	p float64 // Value of cut
	l *Node   // Pointer to left child
	r *Node   // Pointer to right child

}

// Leaf of RCTree containing zero children
type Leaf struct {
	I int       // Index of leaf (user-specified)
	d int       // Depth of leaf
	x []float64 // Original point
}

// NodeObject stores a leaf or branch along with the node type
type NodeObject struct {
	nodeType string      // Type of node - 'Leaf' or 'Branch'
	q        int         // Dimension of cut
	p        float64     // Value of cut
	l        *Node       // Pointer to left child
	r        *Node       // Pointer to right child
	b        [][]float64 // Bounding box of single point or points under branch
	i        int         // Index of leaf (user-specified)
	d        int         // Depth of leaf
	x        []float64   // Original point
	n        int         // Number of leaves under branch or points in leaf
	ixs      int
}

// NewBranch defines a new branch of a tree
func NewBranch(q int, p float64, l *Node, r *Node, u *Node, n int, b [][]float64) *Node {
	node := Node{
		nil,
		&Branch{q, p, l, r},
		b,
		u,
		n,
	}
	return &node
}

// NewLeaf defines a new leaf of a branch
func NewLeaf(i int, d int, u *Node, x []float64, n int) *Node {
	node := Node{
		&Leaf{i, d, x},
		nil,
		num.ArrayReshapeRow(x),
		u,
		n,
	}
	return &node
}

// NewNodeObject defines a new node object for a leaf or branch
func NewNodeObject() *NodeObject {
	nodeObject := NodeObject{}
	return &nodeObject
}

func (node *Node) isLeaf() bool {
	return node.Leaf != nil
}

func (node *Node) isBranch() bool {
	return node.Branch != nil
}

func (node *Node) isRoot() bool {
	return node.u == nil
}
