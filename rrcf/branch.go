package rrcf

// Branch of RCTree containing two children and at most one parent
type Branch struct {
	node *Node
	q    int
	p    int
	l    *Branch
	r    *Branch
	u    *Branch
	n    int
	b    [][]int
}

func NewBranch(q int, p int, l *Branch, r *Branch, u *Branch, n int, b [][]int) *Branch {
	branch := Branch{
		q, p, l, r, u, n, b,
	}
	return *branch
}
