package rrcf

import (
	"fmt"
	"math"
	"testing"

	"github.com/andysgithub/go-rrcf/num"
	"github.com/stretchr/testify/assert"
)

var (
	tree                RRCF
	duplicateTree       RRCF
	treeSeeded          RRCF
	duplicateTreeSeeded RRCF
	indexes             []int
	n                   int
	d                   int
)

func TestInit(t *testing.T) {
	n = 100
	d = 3
	X := num.Randn2(n, d)
	Z := num.ArrayDuplicate(X)
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

	for _, node := range branches {
		leafcount := tree.CountLeaves(&node)
		assert.Equal(t, leafcount, node.n, "Wrong number of leaves on branch")
		bbox := tree.GetBbox(&node)

		for i := 0; i < len(bbox); i++ {
			for j := 0; j < len(bbox[0]); j++ {
				assert.Equal(t, bbox[i][j], node.branch.b[i][j], "Wrong bounding box value for branch")
			}
		}
	}
}

func TestCoDisp(t *testing.T) {
	TestBatch(t)

	for i := range [100]int{} {
		codisp, _ := tree.CoDisp(i)
		assert.Greater(t, codisp, float64(0), fmt.Sprintf("Codisp value %f not greater than 0\n", codisp))
	}
}

func TestDisp(t *testing.T) {
	TestCoDisp(t)

	for i := range [100]int{} {
		disp, _ := tree.Disp(i)
		assert.Greater(t, disp, int(0), fmt.Sprintf("Codisp value %d not greater than 0\n", disp))
	}
}

func TestForgetBatch(t *testing.T) {
	TestDisp(t)

	// Check stored bounding boxes and leaf counts after forgetting points
	for _, index := range indexes {
		forgotten := tree.ForgetPoint(index)
		var branches []Node
		branches = tree.MapBranches(tree.root, branches)
		for _, node := range branches {
			leafcount := tree.CountLeaves(&node)
			assert.Equal(t, leafcount, node.n, fmt.Sprintf("%f - Computed: %d  Stored: %d\n", forgotten.leaf.x, leafcount, node.n))

			bbox := tree.GetBbox(&node)
			result := num.AllClose(bbox, node.branch.b, math.Pow(10, -8))
			assert.True(t, result, fmt.Sprintf("%f - Computed: %d  Stored: %d\n", forgotten.leaf.x, leafcount, node.n))
		}
	}
}

func TestInsertBatch(t *testing.T) {
	TestForgetBatch(t)

	// Check stored bounding boxes and leaf counts after inserting points
	for _, index := range indexes {
		x := num.Randn1(d)
		_, err := tree.InsertPoint(x, index, 0)
		if err == nil {
			var branches []Node
			branches = tree.MapBranches(tree.root, branches)
			for _, node := range branches {
				leafCount := tree.CountLeaves(&node)
				assert.Equal(t, leafCount, node.n, fmt.Sprintf("Computed: %d  Stored: %d\n", leafCount, node.n))

				bbox := tree.GetBbox(&node)
				result := num.AllClose(bbox, node.branch.b, math.Pow(10, -8))
				assert.True(t, result, fmt.Sprintf("Computed: %v  Stored: %v\n", bbox, node.branch.b))
			}
		}
	}
}

func TestBatchWithDuplicates(t *testing.T) {
	TestInsertBatch(t)

	// Instantiate tree with 10 duplicates
	leafCount := duplicateTree.CountLeaves(tree.root)
	assert.Equal(t, leafCount, n, fmt.Sprintf("Leaf count %d not equal to samples %d\n", leafCount, n))

	for i := 90; i < 100; i++ {
		message := fmt.Sprintf("Leaf count %d in duplicate tree not equal to 10\n", duplicateTree.leaves[i].n)
		assert.Equal(t, duplicateTree.leaves[i].n, 10, message)
	}
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
