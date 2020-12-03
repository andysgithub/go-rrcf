package rrcf

// Leaf of RCTree containing no children and at most one parent
type Leaf struct {
	i int       // Index of leaf (user-specified)
	d int       // Depth of leaf
	u *Branch   // Pointer to parent
	x []float64 // Original point (1D)
	n int       // Number of points in leaf (1 if no duplicates)
	b []float64 // Bounding box of point (1D)
}

// NewLeaf defines a new leaf of a branch
func NewLeaf(i int, d int, u *Branch, x []float64, n int) *Leaf {
	leaf := Leaf{
		i, d, u, x, n, x,
	}
	return &leaf
}
