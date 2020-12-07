package rrcf

import "github.com/andysgithub/go-rrcf/num"

// IncrementDepth increments the depth attribute of a leaf
func (rrcf RRCF) IncrementDepth() {
}

// Accumulate counts the number of points in a subtree
func (rrcf RRCF) Accumulate(node *Node, accumulator *int) {
	*accumulator += node.n
}

// GetNodes accumulates a list of all leaves in a subtree
func (rrcf RRCF) GetNodes(node *Node, stack []Node) []Node {
	stack = append(stack, *node)
	return stack
}

// ComputeBbox computes the bbox of a point
func (rrcf RRCF) ComputeBbox(x *Node, mins []float64, maxes []float64) {
	lt := num.ArrayLt(x.x, mins)
	gt := num.ArrayGt(x.x, maxes)

	mins = num.ArrayCopyWhenTrue(mins, x.x, lt)
	maxes = num.ArrayCopyWhenTrue(maxes, x.x, gt)
}

func isBranch(node *Node) bool {
	return node.p != 0
}

func isRoot(node *Node) bool {
	return node.u == nil
}
