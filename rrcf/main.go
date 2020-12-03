package rrcf

import (
	"math/rand"
	"time"

	"github.com/andysgithub/go-rrcf/num"
)

// RRCF - Robust Random Cut Forest
type RRCF struct {
	rng         float64
	leaves      map[int]*Leaf // Map containing pointers to all leaves in tree
	root        *Branch       // Pointer to root of tree
	ndim        int           // Dimension of points in the tree
	indexLabels []int         // Index labels
	u           *Node         // Parent of the current node
}

type Node struct {
	u *Node
}

// NewRRCF - Returns a new random cut forest
func RCTree() RRCF {
	rand.Seed(time.Now().UTC().UnixNano())
	rrcf := RRCF{
		rand.Float64(),
		make(map[int]*Leaf),
		nil, 0, nil, nil,
	}

	return rrcf
}

// Init - Initialises the random cut forest
func (rrcf RRCF) Init(X [][]float64, indexLabels []int, precision int, random_state int64) {
	if random_state != 0 {
		// Random number generation with provided seed
		rand.Seed(random_state)
	}
	rrcf.rng = rand.Float64()

	// Round data to avoid sorting errors
	X = num.Around(X, precision)
	if indexLabels == nil {
		indexLabels = num.Arange(len(X))
	}
	rrcf.indexLabels = indexLabels

	// Remove duplicated rows
	U, I, N := num.Unique(X)

	//fmt.Printf("%v %v %v\n", U, I, N)

	dataRows := len(U)
	dataCols := len(U[0])

	//fmt.Printf("%v %v\n", dataRows, dataCols)

	// Store dimension of dataset
	rrcf.ndim = dataCols

	// Set node above to nil in case of bottom-up search
	rrcf.u = nil

	// Create RRC Tree
	S := num.Ones_bool(dataRows)
	rrcf.MakeTree(X, S, N, I, nil, "root", 0)

	// Remove parent of root
	rrcf.root.u = nil
	// Count all leaves under each branch
	rrcf.CountAllTopDown(rrcf.root)
	// Set bboxes of all branches
	rrcf.GetBboxTopDown(rrcf.root)
}

func (rrcf RRCF) MakeTree(X [][]float64, S []bool, N []int, I []int, parent *Branch, side string, depth int) {
	// Increment depth as we traverse down
	depth++
	// Create a cut according to definition 1
	S1, S2, branch := rrcf.Cut(X, S, parent, side)
	// If S1 does not contain an isolated point
	if contains(S1, true) {
		// Recursively construct tree on S1
		rrcf.MakeTree(X, S1, N, I, branch, "l", depth)
	} else {
		// Create a leaf node from the isolated point
		i := int(num.AsScalar(num.FlatNonZero(S1)))
		leaf := NewLeaf(i, depth, branch, X[i][:], N[i])
		// Link leaf node to parent
		branch.l = leaf
		// If duplicates exist
		if i > 0 {
			// Add a key in the leaves map pointing to leaf for all duplicate indices
			J := num.FlatNonZero(num.ArrayIsEqual(I, i))
			// Get index label
			J = num.ArrayIndices_int(rrcf.indexLabels, J)
			for _, j := range J {
				rrcf.leaves[j] = leaf
			}
		} else {
			i = rrcf.indexLabels[i]
			rrcf.leaves[i] = leaf
		}
	}
	// If S2 does not contain an isolated point
	if contains(S2, true) {
		// Recursively construct tree on S2
		rrcf.MakeTree(X, S2, N, I, branch, "r", depth)
	} else {
		// Create a leaf node from isolated point
	}
}

func contains(array []bool, value bool) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}
	return false
}

// Cut -
func (rrcf RRCF) Cut(X [][]float64, S []bool, parent *Branch, side string) ([]bool, []bool, *Branch) {
	subset := num.ArrayBool_float64(X, S)
	// Find max and min over all d dimensions
	xmax := num.MaxColValues(subset)
	xmin := num.MinColValues(subset)

	// Compute l
	l := num.ArraySub(xmax, xmin)
	l = num.ArrayDiv(l, num.ArraySum(l))

	// Determine dimension to cut
	q := num.RngChoice(rrcf.ndim, l)
	// Determine value for split
	p := num.RngUniform(xmin[q], xmax[q])

	// Determine subset of points to left
	arrayLeq := num.ArrayLeq(num.GetColumn(X, q), p)
	S1 := num.ArrayAnd(arrayLeq, S)
	// Determine subset of points to right
	arrayNot := num.ArrayNot(S1)
	S2 := num.ArrayAnd(arrayNot, S)

	// Create new child node
	child := NewBranch(q, p, nil, nil, parent, 0, nil)

	// Link child node to parent
	if parent != nil {
		switch side {
		case "l":
			parent.l = child
		case "r":
			parent.r = child
		}
	}

	return S1, S2, child
}

// CountAllTopDown -
func (rrcf RRCF) CountAllTopDown(branch *Branch) {

}

// GetBboxTopDown -
func (rrcf RRCF) GetBboxTopDown(branch *Branch) {

}
