package rrcf

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/andysgithub/go-rrcf/num"
	"github.com/stretchr/testify/assert"
)

var (
	tree          RRCF
	duplicateTree RRCF
	indexes       []int
)

func TestInit(t *testing.T) {
	rand.Seed(0)
	n := 100
	d := 3
	X := num.Randn(n, d)
	Z := num.ArrayCopy(X)
	Z = num.ArrayFillRows(Z, 90, 99, float64(1))

	tree := RCTree()
	tree.Init(X, nil, 9, 0)

	duplicateTree := RCTree()
	duplicateTree.Init(Z, nil, 9, 0)

	deck := num.Arange(n)
	deck = num.RndShuffle(deck)
	indexes = deck[:5]
}

func TestBatch(t *testing.T) {
	rand.Seed(0)
	n := 100
	d := 3
	X := num.Randn(n, d)
	Z := num.ArrayCopy(X)
	Z = num.ArrayFillRows(Z, 90, 99, float64(1))

	tree := RCTree()
	tree.Init(X, nil, 9, 0)

	duplicateTree := RCTree()
	duplicateTree.Init(Z, nil, 9, 0)

	deck := num.Arange(n)
	deck = num.RndShuffle(deck)
	indexes = deck[:5]

	//////////////////////////////////////////

	var branches []Node

	// Check stored bounding boxes and leaf counts after instantiating from batch
	branches = tree.MapBranches(tree.root, branches)
	leafcount := tree.CountLeaves(tree.root)
	assert.Equal(t, leafcount, 100, "Wrong number of total leaves")

	for _, branch := range branches {
		leafcount := tree.CountLeaves(&branch)
		assert.Equal(t, leafcount, branch.n, "Wrong number of leaves on branch")
		bbox := tree.GetBbox(&branch)
		fmt.Printf("%v \n", bbox)
	}

}
