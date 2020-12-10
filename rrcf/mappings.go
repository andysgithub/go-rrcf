package rrcf

// MapLeaves traverses the tree recursively, calling Accumulate on leaves
func (rrcf RRCF) MapLeaves(node *Node, accumulator *int) {
	if node.isBranch() {
		if node.branch.l != nil {
			rrcf.MapLeaves(node.branch.l, accumulator)
		}
		if node.branch.r != nil {
			rrcf.MapLeaves(node.branch.r, accumulator)
		}
	} else {
		rrcf.Accumulate(node, accumulator)
	}
}

// MapBranches traverses the tree recursively, calling GetNodes on branches
func (rrcf RRCF) MapBranches(node *Node, branches []Node) []Node {
	if node.isBranch() {
		if node.branch.l != nil {
			branches = rrcf.MapBranches(node.branch.l, branches)
		}
		if node.branch.r != nil {
			branches = rrcf.MapBranches(node.branch.r, branches)
		}
		branches = rrcf.GetNodes(node, branches)
	}
	return branches
}

// MapBboxes traverses the tree recursively, calling GetBbox on leaves
func (rrcf RRCF) MapBboxes(node *Node, mins []float64, maxes []float64) {
	if node.isBranch() {
		if node.branch.l != nil {
			rrcf.MapBboxes(node.branch.l, mins, maxes)
		}
		if node.branch.r != nil {
			rrcf.MapBboxes(node.branch.r, mins, maxes)
		}
	} else {
		rrcf.ComputeBbox(node, mins, maxes)
	}
}

// MapDepths traverses the tree recursively, calling IncrementDepth on leaves
func (rrcf RRCF) MapDepths(node *Node, inc int) {
	if node.isBranch() {
		if node.branch.l != nil {
			rrcf.MapDepths(node.branch.l, inc)
		}
		if node.branch.r != nil {
			rrcf.MapDepths(node.branch.r, inc)
		}
	} else {
		rrcf.IncrementDepth(node, inc)
	}
}
