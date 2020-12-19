package rrcf

import (
	"fmt"
	"math"
	"testing"

	"github.com/andysgithub/go-rrcf/num"
	"github.com/stretchr/testify/assert"
)

var (
	tree                RCTree
	duplicateTree       RCTree
	treeSeeded          RCTree
	duplicateTreeSeeded RCTree
	indexes             []int
	n                   int
	d                   int
	X                   [][]float64
)

func TestInit(t *testing.T) {
	n = 100
	d = 3

	X = num.Randn2(n, d)
	Z := num.ArrayDuplicate(X)
	num.ArrayFillRows(Z, 90, 99, float64(1))

	tree = NewRCTree(X, nil, 9, 0, 1)

	duplicateTree = NewRCTree(Z, nil, 9, 0, 2)

	deck := num.Arange(n)
	deck = num.RndShuffle(deck)
	indexes = deck[:5]
}

func TestBatch(t *testing.T) {
	TestInit(t)

	var branches []Node

	// Check stored bounding boxes and leaf counts after instantiating from batch
	branches = tree.MapBranches(tree.Root, branches)
	leafcount := tree.CountLeaves(tree.Root)
	assert.Equal(t, leafcount, 100, "Wrong number of total leaves")

	for _, node := range branches {
		leafcount := tree.CountLeaves(&node)
		assert.Equal(t, leafcount, node.n, "Wrong number of leaves on branch")
		bbox := tree.GetBbox(&node)

		for i := 0; i < len(bbox); i++ {
			for j := 0; j < len(bbox[0]); j++ {
				assert.Equal(t, bbox[i][j], node.b[i][j], "Wrong bounding box value for branch")
			}
		}
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
	TestInit(t)

	// Check stored bounding boxes and leaf counts after forgetting points
	for _, index := range indexes {
		forgotten := tree.ForgetPoint(index)
		var branches []Node
		branches = tree.MapBranches(tree.Root, branches)
		for _, node := range branches {
			leafcount := tree.CountLeaves(&node)
			assert.Equal(t, leafcount, node.n, fmt.Sprintf("%f - Computed: %d  Stored: %d\n", forgotten.Leaf.x, leafcount, node.n))

			bbox := tree.GetBbox(&node)
			result := num.AllClose(bbox, node.b, math.Pow(10, -8))
			assert.True(t, result, fmt.Sprintf("%f - Computed: %d  Stored: %d\n", forgotten.Leaf.x, leafcount, node.n))
		}
	}
}

func TestInsertBatch(t *testing.T) {
	TestInit(t)

	// Check stored bounding boxes and leaf counts after inserting points
	for _, index := range indexes {
		x := num.Randn1(d)
		_, err := tree.InsertPoint(x, index, 0)
		if err == nil {
			var branches []Node
			branches = tree.MapBranches(tree.Root, branches)
			for _, node := range branches {
				leafCount := tree.CountLeaves(&node)
				assert.Equal(t, leafCount, node.n, fmt.Sprintf("Computed: %d  Stored: %d\n", leafCount, node.n))

				bbox := tree.GetBbox(&node)
				result := num.AllClose(bbox, node.b, math.Pow(10, -8))
				assert.True(t, result, fmt.Sprintf("Computed: %v  Stored: %v\n", bbox, node.b))
			}
		}
	}
}

func TestBatchWithDuplicates(t *testing.T) {
	TestInit(t)
	// Instantiate tree with 10 duplicates
	leafCount := duplicateTree.CountLeaves(duplicateTree.Root)
	assert.Equal(t, leafCount, n, fmt.Sprintf("Leaf count %d not equal to samples %d\n", leafCount, n))

	for i := 90; i < 100; i++ {
		message := fmt.Sprintf("Leaf count %d in duplicate tree not equal to 10\n", duplicateTree.Leaves[i].n)
		assert.Equal(t, duplicateTree.Leaves[i].n, 10, message)
	}
}

func TestInsertDuplicate(t *testing.T) {
	TestInit(t)
	// Insert duplicate point
	point := []float64{1., 1., 1.}
	leaf, _ := duplicateTree.InsertPoint(point, 100, 0)
	assert.Equal(t, leaf.n, 11, "Leaf count %d in duplicate tree not equal to 11\n", leaf.n)
	for i := 90; i < 100; i++ {
		message := fmt.Sprintf("Leaf count %d in duplicate tree not equal to 11\n", duplicateTree.Leaves[i].n)
		assert.Equal(t, duplicateTree.Leaves[i].n, 11, message)
	}
}

func TestFindDuplicate(t *testing.T) {
	TestInsertDuplicate(t)
	// Find duplicate point
	point := []float64{1., 1., 1.}
	duplicate := duplicateTree.FindDuplicate(point, 0)
	assert.NotNil(t, duplicate, "Duplicate in duplicate tree is nil\n")
}

func TestForgetDuplicate(t *testing.T) {
	TestFindDuplicate(t)
	// Forget duplicate point
	duplicateTree.ForgetPoint(100)
	for i := 90; i < 100; i++ {
		message := fmt.Sprintf("Leaf count %d in duplicate tree not equal to 10\n", duplicateTree.Leaves[i].n)
		assert.Equal(t, duplicateTree.Leaves[i].n, 10, message)
	}
}

func TestShingle(t *testing.T) {
	TestInit(t)
	shingle := NewShingle(X, 3)
	step0 := shingle.Next()
	step1 := shingle.Next()

	message := fmt.Sprintf("Shingles misaligned: %v vs %v", step0[1], step1[0])
	assert.True(t, num.ArrayCompareFloat(step0[1], step1[0]), message)
}

func TestInsertDepth(t *testing.T) {
	tree = NewRCTree(nil, nil, 0, 0, 1)

	tree.InsertPoint([]float64{0., 0.}, 0, 0)
	tree.InsertPoint([]float64{0., 0.}, 1, 0)
	tree.InsertPoint([]float64{0., 0.}, 2, 0)
	tree.InsertPoint([]float64{0., 1.}, 3, 0)
	tree.ForgetPoint(3)

	minDepth := math.MaxInt64
	for _, node := range tree.Leaves {
		if node.Leaf.d < minDepth {
			minDepth = node.Leaf.d
		}
	}
	assert.GreaterOrEqual(t, minDepth, 0)
}
