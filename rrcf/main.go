package rrcf

import (
	"math/rand"
	"time"

	"github.com/andysgithub/go-rrcf/num"
)

// RRCF - Robust Random Cut Forest
type RRCF struct {
	rng          float64
	leaves       map[string]Leaf // Map containing pointers to all leaves in tree
	root         Branch          // Pointer to root of tree
	ndim         int             // Dimension of points in the tree
	index_labels []int           // Index labels
	u            Node            // Parent of the current node
}

type Tree struct {
}

type Node struct {
}

type Leaf struct {
}

// NewRRCF - Returns a new forest
func NewRRCF() RRCF {
	rand.Seed(time.Now().UTC().UnixNano())
	rrcf := RRCF{
		rand.Float64(),
		make(map[string]int),
		0,
		0,
	}

	return rrcf
}

// Init - Initialises the random cut forest
func (rrcf RRCF) Init(X [][]float64, index_labels int, precision int, random_state int) {
	if random_state != 0 {
		// Random number generation with provided seed
		rand.Seed(random_state)
		rrcf.rng = rand.Float64()
	}
	// Round data to avoid sorting errors
	X = num.Around(X, precision)
	if index_labels == 0 {
		rrcf.index_labels = num.Arange(len(X[0]))
	} else {
		rrcf.index_labels = index_labels
	}

	// Remove duplicated rows
	U, I, N := num.Unique(X)

	dataRows := len(U)
	dataCols := len(U[0])

	// Store dimension of dataset
	rrcf.ndim = dataCols
	// Set node above to None in case of bottom-up search

}
