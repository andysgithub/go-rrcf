package rrcf

// Node of RCTree containing zero or two children and at most one parent
type Node struct {
	// For leaf:
	i  int       // Index of leaf (user-specified)
	d  int       // Depth of leaf
	x  []float64 // Original point (1D)
	b1 []float64 // Bounding box of point (1D)

	// For branch:
	q  int         // Dimension of cut
	p  float64     // Value of cut
	l  *Node       // Pointer to left child
	r  *Node       // Pointer to right child
	b2 [][]float64 // Bounding box of points under branch (2D)

	// Common:
	u *Node // Pointer to parent
	n int   // Number of leaves under branch or points in leaf
}

// NewBranch defines a new branch of a tree
func NewBranch(q int, p float64, l *Node, r *Node, u *Node, n int, b [][]float64) *Node {
	branch := Node{
		0, 0, nil, nil, q, p, l, r, b, u, n,
	}
	return &branch
}

// NewLeaf defines a new leaf of a branch
func NewLeaf(i int, d int, u *Node, x []float64, n int) *Node {
	leaf := Node{
		i, d, x, x, 0, 0, nil, nil, nil, u, n,
	}
	return &leaf
}
