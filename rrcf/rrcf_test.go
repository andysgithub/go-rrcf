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

	tree = RCTree()
	tree.Init(X, nil, 9, 0)

	duplicateTree = RCTree()
	duplicateTree.Init(Z, nil, 9, 0)

	deck := num.Arange(n)
	deck = num.RndShuffle(deck)
	indexes = deck[:5]
}

func TestBatch(t *testing.T) {
	TestInit(t)

	var branches []Node

	// Check stored bounding boxes and leaf counts after instantiating from batch
	branches = tree.MapBranches(tree.root, branches)
	leafcount := tree.CountLeaves(tree.root)
	assert.Equal(t, leafcount, 100, "Wrong number of total leaves")

	for _, branch := range branches {
		leafcount := tree.CountLeaves(&branch)
		assert.Equal(t, leafcount, branch.n, "Wrong number of leaves on branch")
		bbox := tree.GetBbox(&branch)

		for i := 0; i < len(bbox); i++ {
			for j := 0; j < len(bbox[0]); j++ {
				assert.Equal(t, bbox[i][j], branch.b2[i][j], "Wrong bounding box value for branch")
			}
		}

		// fmt.Printf("%v \n", bbox)
	}
}

func TestCoDisp(t *testing.T) {
	TestInit(t)

	for i := range [100]int{} {
		codisp, _ := tree.CoDisp(i)
		assert.Greater(t, codisp, float64(0), fmt.Sprintf("Codisp value %f not greater than 0\n", codisp))
	}
}

func TestDisp(t *testing.T) {
	TestInit(t)

	for i := range [100]int{} {
		disp, _ := tree.Disp(i)
		assert.Greater(t, disp, int(0), fmt.Sprintf("Codisp value %d not greater than 0\n", disp))
	}
}

func TestForgetBatch(t *testing.T) {
}

func TestInsertBatch(t *testing.T) {
}

func TestBatchWithDuplicates(t *testing.T) {
}

func TestInsertDuplicate(t *testing.T) {
}

func TestFindDuplicate(t *testing.T) {
}

func TestForgetDuplicate(t *testing.T) {
}

func TestShingle(t *testing.T) {
}

func TestRandomState(t *testing.T) {
}

func TestInsertDepth(t *testing.T) {
}

func TestToDict(t *testing.T) {
}

func TestFromDict(t *testing.T) {
}

func TestPrint(t *testing.T) {
}
