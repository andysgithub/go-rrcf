package rrcf

// MapLeaves traverses the tree recursively, calling Accumulate on leaves
func (rcTree RCTree) MapLeaves(node *Node, accumulator *int) {
	if node.isBranch() {

		// Process without recursion if both children are leaves
		if node.Branch.l.isLeaf() && node.Branch.r.isLeaf() {
			rcTree.Accumulate(node.Branch.l, accumulator)
			rcTree.Accumulate(node.Branch.r, accumulator)
			return
		}

		if node.Branch.l != nil {
			rcTree.MapLeaves(node.Branch.l, accumulator)
		}
		if node.Branch.r != nil {
			rcTree.MapLeaves(node.Branch.r, accumulator)
		}
	} else {
		rcTree.Accumulate(node, accumulator)
	}
}

// MapBranches traverses the tree recursively, calling GetNodes on branches
func (rcTree RCTree) MapBranches(node *Node, branches []Node) []Node {
	if node.isBranch() {

		// Process without recursion if both children are leaves
		if node.Branch.l.isLeaf() && node.Branch.r.isLeaf() {
			branches = rcTree.GetNodes(node.Branch.l, branches)
			branches = rcTree.GetNodes(node.Branch.r, branches)
			return branches
		}

		if node.Branch.l != nil {
			branches = rcTree.MapBranches(node.Branch.l, branches)
		}
		if node.Branch.r != nil {
			branches = rcTree.MapBranches(node.Branch.r, branches)
		}
		branches = rcTree.GetNodes(node, branches)
	}
	return branches
}

// MapBboxes traverses the tree recursively, calling GetBbox on leaves
func (rcTree RCTree) MapBboxes(node *Node, mins []float64, maxes []float64) {
	if node.isBranch() {

		// Process without recursion if both children are leaves
		if node.Branch.l.isLeaf() && node.Branch.r.isLeaf() {
			rcTree.ComputeBbox(node.Branch.l, mins, maxes)
			rcTree.ComputeBbox(node.Branch.r, mins, maxes)
			return
		}

		if node.Branch.l != nil {
			rcTree.MapBboxes(node.Branch.l, mins, maxes)
		}
		if node.Branch.r != nil {
			rcTree.MapBboxes(node.Branch.r, mins, maxes)
		}
	} else {
		rcTree.ComputeBbox(node, mins, maxes)
	}
}

// MapDepths traverses the tree recursively, calling IncrementDepth on leaves
func (rcTree RCTree) MapDepths(node *Node, inc int) {
	if node.isBranch() {

		// Process without recursion if both children are leaves
		if node.Branch.l.isLeaf() && node.Branch.r.isLeaf() {
			rcTree.IncrementDepth(node.Branch.l, inc)
			rcTree.IncrementDepth(node.Branch.r, inc)
			return
		}

		if node.Branch.l != nil {
			rcTree.MapDepths(node.Branch.l, inc)
		}
		if node.Branch.r != nil {
			rcTree.MapDepths(node.Branch.r, inc)
		}
	} else {
		rcTree.IncrementDepth(node, inc)
	}
}
