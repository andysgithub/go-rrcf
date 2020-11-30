package rrcf

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/andysgithub/go-rrcf/num"
)

// RRCF - Robust Random Cut Forest
type RRCF struct {
	rng          float64
	leaves       map[string]*Leaf // Map containing pointers to all leaves in tree
	root         *Branch          // Pointer to root of tree
	ndim         int              // Dimension of points in the tree
	index_labels []int            // Index labels
	u            *Node            // Parent of the current node
}

type Node struct {
	u *Node
}

type Leaf struct {
}

// NewRRCF - Returns a new random cut forest
func NewRRCF() RRCF {
	rand.Seed(time.Now().UTC().UnixNano())
	rrcf := RRCF{
		rand.Float64(),
		make(map[string]*Leaf),
		nil, 0, nil, nil,
	}

	return rrcf
}

// Init - Initialises the random cut forest
func (rrcf RRCF) Init(X [][]float64, index_labels []int, precision int, random_state int64) {
	if random_state != 0 {
		// Random number generation with provided seed
		rand.Seed(random_state)
		rrcf.rng = rand.Float64()
	}
	// Round data to avoid sorting errors
	X = num.Around(X, precision)
	if index_labels == nil {
		rrcf.index_labels = num.Arange(len(X[0]))
	} else {
		rrcf.index_labels = index_labels
	}

	// Remove duplicated rows
	U, I, N := num.Unique(X)

	fmt.Printf("%v %v %v\n", U, I, N)

	dataRows := len(U)
	dataCols := len(U[0])

	fmt.Printf("%v %v\n", dataRows, dataCols)

	// Store dimension of dataset
	rrcf.ndim = dataCols

	// Set node above to nil in case of bottom-up search
	rrcf.u = nil

	// Create RRC Tree
	S := num.Ones_bool(dataRows, dataCols)
	rrcf.MakeTree(X, S, N, I, nil, "root", 0)

	// Remove parent of root
	rrcf.root.node.u = nil
	// Count all leaves under each branch
	rrcf.CountAllTopDown(rrcf.root)
	// Set bboxes of all branches
	rrcf.GetBboxTopDown(rrcf.root)
}

func (rrcf RRCF) MakeTree(X [][]float64, S [][]bool, N []int, I []int, parent *Node, side string, depth int) {
	// Increment depth as we traverse down
	depth++
	// Create a cut according to definition 1
	S1, S2, branch = rrcf.Cut(X, S, parent, side)
	// If S1 does not contain an isolated point
	if S1.sum() > 1 {
		// Recursively construct tree on S1
		rrcf.MakeTree(X, S1, N, I, branch.node, 'l', depth)
	} else {
		// Create a leaf node from the isolated point
        i = num.AsScalar(num.FlatNonZero(S1))
	}
}

func (rrcf RRCF) Cut(X [][]float64, S [][]bool, parent *Node, side string) [][]bool, [][]bool,  {

}

func (rrcf RRCF) CountAllTopDown(branch *Branch) {

}

func (rrcf RRCF) GetBboxTopDown(branch *Branch) {

}
