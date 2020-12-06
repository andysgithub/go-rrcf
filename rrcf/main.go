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

// MakeTree -
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

// Cut -
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

	// ALERT: If S1 = S, S2 will be all false
	// arrayLeq must be false in a position where S is true

	if num.ArrayCompare(S, S1) {
		bufferMain.Reset()
		bufferMain.WriteString("Warning: S = S1\n")
		bufferMain.WriteString(fmt.Sprintf("Total points in S: %d\n", num.ArraySumBool(S)))
		bufferMain.WriteString(fmt.Sprintf("Total points in S1: %d\n", num.ArraySumBool(S1)))
		bufferMain.WriteString(fmt.Sprintf("Total points in S2: %d\n", num.ArraySumBool(S2)))
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

// GetBbox computes the bounding box of all points underneath a given branch
func (rrcf RRCF) GetBbox(branch *Node) [][]float64 {
	if branch == nil {
		branch = rrcf.root
	}

	mins := num.Full(rrcf.ndim, math.Inf(1))
	maxes := num.Full(rrcf.ndim, math.Inf(-1))
	//rrcf.MapLeaves(branch, GetBbox, mins, maxes)
	bbox := num.ArrayVStack(mins, maxes)

	return bbox
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

// MapLeaves traverses the tree recursively, calling function f on leaves
func (rrcf RRCF) MapLeaves(node *Node, f func(*Node, *int), accumulator *int) {
	if isBranch(node) {
		if node.l != nil {
			rrcf.MapLeaves(node.l, f, accumulator)
		}
		if node.r != nil {
			rrcf.MapLeaves(node.r, f, accumulator)
		}
	} else {
		f(node, accumulator)
	}
}

// GetNodes accumulates a list of all leaves in a subtree
func (rrcf RRCF) GetNodes(node *Node, stack []Node) []Node {
	stack = append(stack, *node)
	return stack
}

// Accumulate counts the number of points in a subtree
func (rrcf RRCF) Accumulate(node *Node, accumulator *int) {
	*accumulator += node.n
}

// CountLeaves counts the total leaves underneath a single node
func (rrcf RRCF) CountLeaves(branch *Node) int {
	var numLeaves int

	rrcf.MapLeaves(branch, rrcf.Accumulate, &numLeaves)
	return numLeaves
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

func isBranch(node *Node) bool {
	return node.p != 0
}
