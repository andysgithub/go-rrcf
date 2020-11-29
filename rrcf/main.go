package rrcf

import (
  "math/rand"
  "github.com/andysgithub/go-rrcf/num"
	"time"
)

// RRCF - Robust Random Cut Forest
type RRCF struct {
	rng    float64
	leaves map[string]int64 // Map containing pointers to all leaves in tree
	root   int64            // Pointer to root of tree
  ndim   int              // Dimension of points in the tree
  index_labels []int      // Index labels
}

// NewRRCF - Returns a new forest
func NewRRCF() RRCF {
	rand.Seed(time.Now().UTC().UnixNano())
	rrcf := RRCF{
		rand.Float64(),
		make(map[string]int64),
		0,
		0,
	}

	return rrcf
}

// Init - Initialises the random cut forest
func (rrcf RRCF) Init(X [][]int64, index_labels int, precision int, random_state int64) {
	if random_state != 0 {
		// Random number generation with provided seed
		rand.Seed(random_state)
		rrcf.rng = rand.Float64()
	}
  // Round data to avoid sorting errors
  X = num.Around(X, precision)

  // Initialize index labels
  if len(index_labels) == 0 {
    rrcf.index_labels = num.Arange(len(X[0]))
  } else {
    rrcf.index_labels = index_labels
  }

  // Check for duplicates
  U, I, N = num.Unique(X, return_inverse=True, return_counts=True,
    axis=0)

}
