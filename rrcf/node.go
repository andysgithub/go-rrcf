package rrcf

// Node of RCTree consisting of a leaf or branch and containing at most one parent
type Node struct {
	leaf   *Leaf
	branch *Branch
	u      *Node // Pointer to parent
	n      int   // Number of leaves under branch or points in leaf
}

// Branch of RCTree containing two children
type Branch struct {
	q int         // Dimension of cut
	p float64     // Value of cut
	l *Node       // Pointer to left child
	r *Node       // Pointer to right child
	b [][]float64 // Bounding box of points under branch
}

// Leaf of RCTree containing zero children
type Leaf struct {
	i int       // Index of leaf (user-specified)
	d int       // Depth of leaf
	x []float64 // Original point
	b []float64 // Bounding box of single point
}

// NewBranch defines a new branch of a tree
func NewBranch(q int, p float64, l *Node, r *Node, u *Node, n int, b [][]float64) *Node {
	node := Node{
		nil,
		&Branch{q, p, l, r, b},
		u,
		n,
	}
	return &node
}

// NewLeaf defines a new leaf of a branch
func NewLeaf(i int, d int, u *Node, x []float64, n int) *Node {
	node := Node{
		&Leaf{i, d, x, x},
		nil,
		u,
		n,
	}
	return &node
}

// func (node *Node) nVal() int {
// 	if node.isBranch() {
// 		return node.branch.n
// 	}
// 	return node.leaf.n
// }

func (node *Node) isLeaf() bool {
	return node.leaf != nil
}

func (node *Node) isBranch() bool {
	return node.branch != nil
}

func (node *Node) isRoot() bool {
	return node.u == nil
}
