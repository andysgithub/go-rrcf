package rrcf

// MapLeaves traverses the tree recursively, calling Accumulate on leaves
func (rcTree RCTree) MapLeaves(node *Node, accumulator *int) {
	if node.isBranch() {
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
