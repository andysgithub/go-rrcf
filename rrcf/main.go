package rrcf

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/andysgithub/go-rrcf/num"
)

var (
	logMain    *os.File
	bufferMain bytes.Buffer
)

// RRCF - Robust Random Cut Forest
type RRCF struct {
	rng         float64
	leaves      map[int]*Node // Map containing pointers to all leaves in tree
	root        *Node         // Pointer to root of tree
	ndim        int           // Dimension of points in the tree
	indexLabels []int         // Index labels
	u           *Node         // Parent of the current node
}

// RCTree returns a new random cut forest
func RCTree() RRCF {
	rand.Seed(time.Now().UTC().UnixNano())
	rrcf := RRCF{
		rand.Float64(),
		make(map[int]*Node),
		nil, 0, nil, nil,
	}

	return rrcf
}

// Init - Initialises the random cut forest
func (rrcf *RRCF) Init(X [][]float64, indexLabels []int, precision int, randomState int64) {
	os.Remove("logs/rctree.log")
	file, err := os.OpenFile("logs/rctree.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		defer logMain.Close()
		logMain = file
	}

	if randomState != 0 {
		// Random number generation with provided seed
		rand.Seed(randomState)
	}
	rrcf.rng = rand.Float64()

	// Round data to avoid sorting errors
	X = num.Around(X, precision)
	if indexLabels == nil {
		indexLabels = num.Arange(len(X))
	}
	rrcf.indexLabels = indexLabels

	// Remove duplicated rows
	X, I, N := num.Unique(X)

	//fmt.Printf("%v %v %v\n", X, I, N)

	dataRows := len(X)
	dataCols := len(X[0])

	bufferMain.Reset()
	bufferMain.WriteString("\n-------- Tree Build Start --------\n")
	bufferMain.WriteString(fmt.Sprintf("Rows: %d  Cols: %d\n", dataRows, dataCols))
	fmt.Fprintln(logMain, bufferMain.String())

	//fmt.Printf("%v %v\n", dataRows, dataCols)

	// Store dimension of dataset
	rrcf.ndim = dataCols

	// Set node above to nil in case of bottom-up search
	rrcf.u = nil

	// Create RRC Tree
	S := num.OnesBool(dataRows)
	rrcf.MakeTree(X, S, N, I, rrcf.root, "root", 0)

	// Remove parent of root
	rrcf.root.u = nil
	// Count all leaves under each branch
	rrcf.CountAllTopDown(rrcf.root)
	// Set bboxes of all branches
	rrcf.GetBboxTopDown(rrcf.root)
}

// MakeTree generates a random cut tree
func (rrcf *RRCF) MakeTree(X [][]float64, S []bool, N []int, I []int, parent *Node, side string, depth int) {

	bufferMain.Reset()
	bufferMain.WriteString("\nCreate Branch\n")
	bufferMain.WriteString(fmt.Sprintf("Side: %s  Depth: %d\n", side, depth))
	fmt.Fprintln(logMain, bufferMain.String())

	// Increment depth as we traverse down
	depth++
	// Create a cut according to definition 1
	S1, S2, branch := rrcf.Cut(X, S, parent, side)
	// If S1 does not contain an isolated point
	if num.ArraySumBool(S1) > 1 {
		// Recursively construct tree on S1
		rrcf.MakeTree(X, S1, N, I, branch, "l", depth)
	} else {
		// Create a leaf node from the isolated point
		i := int(num.AsScalar(num.FlatNonZero(S1)))
		leaf := NewLeaf(i, depth, branch, X[i][:], N[i])
		// Link leaf node to parent
		branch.l = leaf
		// If duplicates exist
		if I != nil {
			// Add a key in the leaves map pointing to leaf for all duplicate indices
			J := num.FlatNonZero(num.ArrayEqInt(I, i))
			// Get index label
			J = num.ArrayIndicesInt(rrcf.indexLabels, J)
			for _, j := range J {
				rrcf.leaves[j] = leaf
			}
		} else {
			i = rrcf.indexLabels[i]
			rrcf.leaves[i] = leaf
		}
	}
	// If S2 does not contain an isolated point
	if num.ArraySumBool(S2) > 1 {
		// Recursively construct tree on S2
		rrcf.MakeTree(X, S2, N, I, branch, "r", depth)
	} else {
		// Create a leaf node from isolated point
		i := num.AsScalar(num.FlatNonZero(S2))
		leaf := NewLeaf(i, depth, branch, X[i][:], N[i])
		// Link leaf node to parent
		branch.r = leaf
		// If duplicates exist
		if I != nil {
			// Add a key in the leaves map pointing to leaf for all duplicate indices
			J := num.FlatNonZero(num.ArrayEqInt(I, i))
			// Get index label
			J = num.ArrayIndicesInt(rrcf.indexLabels, J)
			for _, j := range J {
				rrcf.leaves[j] = leaf
			}
		} else {
			i = rrcf.indexLabels[i]
			rrcf.leaves[i] = leaf
		}
	}
	depth--
}

// Cut creates a child node to the left or right of the parent
func (rrcf *RRCF) Cut(X [][]float64, S []bool, parent *Node, side string) ([]bool, []bool, *Node) {
	subset := num.ArrayBoolFloat64(X, S)
	// Find max and min over all d dimensions
	xmax := num.MaxColValues(subset)
	xmin := num.MinColValues(subset)

	// Compute l
	l := num.ArraySub(xmax, xmin)
	l = num.ArrayDiv(l, num.ArraySumFloat(l))

	// Determine dimension to cut
	q := num.RndChoice(rrcf.ndim, l)
	// Determine value for split
	p := num.RndUniform(xmin[q], xmax[q])

	bufferMain.Reset()
	bufferMain.WriteString("\nCut Tree\n")
	bufferMain.WriteString(fmt.Sprintf("l: %v\nq: %d  p: %f\n", l, q, p))
	fmt.Fprintln(logMain, bufferMain.String())

	// Determine subset of points to left
	arrayLeq := num.ArrayLeq(num.GetColumn(X, q), p)
	S1 := num.ArrayAnd(arrayLeq, S) // S1 is all points in S with random dimension < split value
	// Determine subset of points to right
	arrayNot := num.ArrayNot(S1)
	S2 := num.ArrayAnd(arrayNot, S) // S2 is all the points in S not in S1

	if num.ArrayCompare(S, S1) {
		bufferMain.Reset()
		bufferMain.WriteString("Warning: S equals S1\n")
		fmt.Fprintln(logMain, bufferMain.String())
	}

	// Create new child node
	child := NewBranch(q, p, nil, nil, parent, 0, nil)

	// Link child node to parent
	if parent != nil || side == "root" {
		switch side {
		case "l":
			parent.l = child
		case "r":
			parent.r = child
		case "root":
			rrcf.root = child
		}
	}

	return S1, S2, child
}

// MapLeaves traverses the tree recursively, calling Accumulate on leaves
func (rrcf RRCF) MapLeaves(node *Node, accumulator *int) {
	if isBranch(node) {
		if node.l != nil {
			rrcf.MapLeaves(node.l, accumulator)
		}
		if node.r != nil {
			rrcf.MapLeaves(node.r, accumulator)
		}
	} else {
		rrcf.Accumulate(node, accumulator)
	}
}

// MapBranches traverses the tree recursively, calling GetNodes on branches
func (rrcf RRCF) MapBranches(node *Node, branches []Node) []Node {
	if isBranch(node) {
		if node.l != nil {
			branches = rrcf.MapBranches(node.l, branches)
		}
		if node.r != nil {
			branches = rrcf.MapBranches(node.r, branches)
		}
		branches = rrcf.GetNodes(node, branches)
	}
	return branches
}

// MapBboxes traverses the tree recursively, calling GetBbox on leaves
func (rrcf RRCF) MapBboxes(node *Node, mins []float64, maxes []float64) {
	if isBranch(node) {
		if node.l != nil {
			rrcf.MapBboxes(node.l, mins, maxes)
		}
		if node.r != nil {
			rrcf.MapBboxes(node.r, mins, maxes)
		}
	} else {
		rrcf.ComputeBbox(node, mins, maxes)
	}
}

// ForgetPoint -
func (rrcf RRCF) ForgetPoint() {

}

// UpdateLeafCountUpwards -
func (rrcf RRCF) UpdateLeafCountUpwards() {

}

// InsertPoint -
func (rrcf RRCF) InsertPoint() {

}

// Query -
func (rrcf RRCF) Query() {

}

// Disp computes displacement at leaf
func (rrcf RRCF) Disp(param interface{}) (int, error) {
	leaf, ok := param.(*Node)
	if !ok {
		index, ok := param.(int)
		if !ok {
			return 0, fmt.Errorf("CoDisp parameter not recognised: %v", leaf)
		}
		leaf = rrcf.leaves[index]
	}

	// Handle case where leaf is root
	if isRoot(leaf) {
		return 0, nil
	}

	parent := leaf.u
	// Find sibling
	sibling := parent.l
	if leaf == parent.l {
		sibling = parent.r
	}
	// Count number of nodes in sibling subtree
	displacement := sibling.n
	return displacement, nil
}

// CoDisp computes collusive displacement at leaf
func (rrcf RRCF) CoDisp(param interface{}) (float64, error) {
	leaf, ok := param.(*Node)
	if !ok {
		index, ok := param.(int)
		if !ok {
			return 0, fmt.Errorf("CoDisp parameter not recognised: %v", leaf)
		}
		leaf = rrcf.leaves[index]
	}

	// Handle case where leaf is root
	if isRoot(leaf) {
		return 0, nil
	}
	node := leaf
	var results []float64

	for i := 0; i < node.d; i++ {
		parent := node.u
		if parent == nil {
			break
		}
		sibling := parent.l
		if node == parent.l {
			sibling = parent.r
		}
		numDeleted := node.n
		displacement := sibling.n
		result := float64(displacement / numDeleted)
		results = append(results, result)
		node = parent
	}
	coDisplacement := num.ArrayMaxValue(results)
	return coDisplacement, nil
}

// GetBbox computes the bounding box of all points underneath a given branch
func (rrcf RRCF) GetBbox(branch *Node) [][]float64 {
	if branch == nil {
		branch = rrcf.root
	}

	mins := num.Full(rrcf.ndim, math.Inf(1))
	maxes := num.Full(rrcf.ndim, math.Inf(-1))
	rrcf.MapBboxes(branch, mins, maxes)
	bbox := num.ArrayVStack(mins, maxes)

	return bbox
}

// FindDuplicate -
func (rrcf RRCF) FindDuplicate() {

}

// ToDict -
func (rrcf RRCF) ToDict() {

}

// Serialise -
func (rrcf RRCF) Serialise() {

}

// LoadDict -
func (rrcf RRCF) LoadDict() {

}

// Deserialise -
func (rrcf RRCF) Deserialise() {

}

// FromDict -
func (rrcf RRCF) FromDict() {

}

// lrBranchBbox computes the bbox of a node based on bboxes of the node's children
func lrBranchBbox(node *Node) [][]float64 {
	var bbLeft, bbRight, bbLastLeft, bbLastRight []float64

	if isBranch(node.l) {
		lastLeft := len(node.l.b2) - 1
		bbLeft = node.l.b2[0][:]
		bbLastLeft = node.l.b2[lastLeft][:]
	} else {
		bbLeft = node.l.b1[:]
		bbLastLeft = bbLeft
	}

	if isBranch(node.r) {
		lastRight := len(node.r.b2) - 1
		bbRight = node.r.b2[0][:]
		bbLastRight = node.r.b2[lastRight][:]
	} else {
		bbRight = node.r.b1[:]
		bbLastRight = bbRight
	}

	bbox := num.ArrayVStack(
		num.ArrayMinimum(bbLeft, bbRight),
		num.ArrayMaximum(bbLastLeft, bbLastRight))

	return bbox
}

// GetBboxTopDown recursively computes bboxes of all branches from root to leaves
func (rrcf RRCF) GetBboxTopDown(node *Node) {
	if isBranch(node) {
		if node.l != nil {
			rrcf.GetBboxTopDown(node.l)
		}
		if node.r != nil {
			rrcf.GetBboxTopDown(node.r)
		}
		bbox := lrBranchBbox(node)
		node.b2 = bbox
	}
}

// CountAllTopDown recursively computes the number of leaves below each branch from root to leaves
func (rrcf RRCF) CountAllTopDown(node *Node) {
	if isBranch(node) {
		if node.l != nil {
			rrcf.CountAllTopDown(node.l)
		}
		if node.r != nil {
			rrcf.CountAllTopDown(node.r)
		}
		node.n = node.l.n + node.r.n
	}
}

// CountLeaves counts the total leaves underneath a single node
func (rrcf RRCF) CountLeaves(branch *Node) int {
	var numLeaves int

	rrcf.MapLeaves(branch, &numLeaves)
	return numLeaves
}

// SearchForLeaf -
func (rrcf RRCF) SearchForLeaf() {

}

// TightenBboxUpwards -
func (rrcf RRCF) TightenBboxUpwards() {

}

// RelaxBboxUpwards -
func (rrcf RRCF) RelaxBboxUpwards() {

}

// InsertPointCut -
func (rrcf RRCF) InsertPointCut() {

}
