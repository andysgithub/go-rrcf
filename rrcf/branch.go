package rrcf

// Branch of RCTree containing two children and at most one parent
type Branch struct {
	q int     // Dimension of cut
	p float64 // Value of cut
	l *Branch // Pointer to left child
	r *Branch // Pointer to right child
	u *Branch // Pointer to parent
	n int     // Number of leaves under branch
	b [][]int // Bounding box of points under branch (2D)
}

// NewBranch defines a new branch of a tree
func NewBranch(q int, p float64, l *Branch, r *Branch, u *Branch, n int, b [][]int) *Branch {
	branch := Branch{
		q, p, l, r, u, n, b,
	}
	return &branch
}
