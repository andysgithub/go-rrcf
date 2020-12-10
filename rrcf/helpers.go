package rrcf

import "github.com/andysgithub/go-rrcf/num"

// IncrementDepth increments the depth attribute of a leaf
func (rrcf RRCF) IncrementDepth(node *Node, increment int) {
	node.leaf.d += increment
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
	lt := num.ArrayLt(x.leaf.x, mins)
	gt := num.ArrayGt(x.leaf.x, maxes)

	num.ArrayCopyWhenTrue(mins, x.leaf.x, lt)
	num.ArrayCopyWhenTrue(maxes, x.leaf.x, gt)
}

// RemoveIndex removes the element at index and move all later values up
// Returns the element removed
func RemoveIndex(s map[int]*Node, index int) *Node {
	element := s[index]
	delete(s, index)
	return element
}
