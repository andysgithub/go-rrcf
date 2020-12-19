package rrcf

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/andysgithub/go-rrcf/logging"
	"github.com/andysgithub/go-rrcf/num"
	"github.com/andysgithub/go-rrcf/random"
)

// RCTree - Robust Random Cut Forest
type RCTree struct {
	Leaves      map[int]*Node // Map containing pointers to all leaves in tree
	Root        *Node         // Pointer to root of tree
	Ndim        int           // Dimension of points in the tree
	IndexLabels []int         // Index labels
	Parent      *Node         // Parent of the current node
	Log         *logging.Logger
	rnd         random.Random
}

// NewRCTree returns a new random cut forest
func NewRCTree(X [][]float64, indexLabels []int, precision int, randomState interface{}, id int) RCTree {
	rnd := random.NewRandom()

	log := logging.NewLogger("logs/rctree-all.log")
	if id == 0 {
		log = logging.NewLogger(fmt.Sprintf("logs/rctree-%d.log", id))
	}

	rct := RCTree{
		make(map[int]*Node),
		nil, 0, nil, nil,
		log,
		rnd,
	}

	rct.Init(X, indexLabels, precision, randomState)
	return rct
}

// Init - Initialises the random cut forest
func (rct *RCTree) Init(X [][]float64, indexLabels []int, precision int, randomState interface{}) {
	switch randomState.(type) {
	case int:
		// Random number generation with provided seed
		rand.Seed((int64)(randomState.(int)))
		// rct.rnd.Seed((int64)(randomState.(int)))
	default:
		// Random number generation with random seed
		rand.Seed(time.Now().UTC().UnixNano())
		// rct.rnd.Seed(time.Now().UTC().UnixNano())
	}

	if X != nil {
		// Round data to avoid sorting errors
		X = num.Around(X, precision)
		if indexLabels == nil {
			indexLabels = num.Arange(len(X))
		}
		rct.IndexLabels = indexLabels

		// Remove duplicated rows
		X, I, N := num.Unique(X)

		dataRows := len(X)
		dataCols := len(X[0])

		rct.Log.Section("\n-------- Tree Build Start --------\n")
		rct.Log.AddLast(fmt.Sprintf("Rows: %d  Cols: %d\n", dataRows, dataCols))

		// Store dimension of dataset
		rct.Ndim = dataCols

		// Set node above to nil in case of bottom-up search
		rct.Parent = nil

		// Create RRC Tree
		S := num.OnesBool(dataRows)
		rct.MakeTree(X, S, N, I, rct.Root, "root", 0)

		// Remove parent of root
		rct.Root.u = nil
		// Count all leaves under each branch
		rct.CountAllTopDown(rct.Root)
		// Set bboxes of all branches
		rct.GetBboxTopDown(rct.Root)
	}
}

// MakeTree generates a random cut tree
func (rct *RCTree) MakeTree(X [][]float64, S []bool, N []int, I []int, parent *Node, side string, depth int) {

	rct.Log.Section("\nCreate Branch\n")
	rct.Log.AddLast(fmt.Sprintf("Side: %s  Depth: %d\n", side, depth))

	// Increment depth as we traverse down
	depth++
	// Create a cut according to definition 1
	S1, S2, node := rct.Cut(X, S, parent, side)
	// If S1 does not contain an isolated point
	if num.ArraySumBool(S1) > 1 {
		// Recursively construct tree on S1
		rct.MakeTree(X, S1, N, I, node, "l", depth)
	} else {
		// Create a leaf node from the isolated point
		i := int(num.AsScalar(num.FlatNonZero(S1)))
		leaf := NewLeaf(i, depth, node, X[i][:], N[i])
		// Link leaf node to parent
		node.Branch.l = leaf
		// If duplicates exist
		if I != nil {
			// Add a key in the leaves map pointing to leaf for all duplicate indices
			J := num.FlatNonZero(num.ArrayEqInt(I, i))
			// Get index label
			J = num.ArrayIndicesInt(rct.IndexLabels, J)
			for _, j := range J {
				rct.Leaves[j] = leaf
			}
		} else {
			i = rct.IndexLabels[i]
			rct.Leaves[i] = leaf
		}
	}
	// If S2 does not contain an isolated point
	if num.ArraySumBool(S2) > 1 {
		// Recursively construct tree on S2
		rct.MakeTree(X, S2, N, I, node, "r", depth)
	} else {
		// Create a leaf node from isolated point
		i := num.AsScalar(num.FlatNonZero(S2))
		leaf := NewLeaf(i, depth, node, X[i][:], N[i])
		// Link leaf node to parent
		node.Branch.r = leaf
		// If duplicates exist
		if I != nil {
			// Add a key in the leaves map pointing to leaf for all duplicate indices
			J := num.FlatNonZero(num.ArrayEqInt(I, i))
			// Get index label
			J = num.ArrayIndicesInt(rct.IndexLabels, J)
			for _, j := range J {
				rct.Leaves[j] = leaf
			}
		} else {
			i = rct.IndexLabels[i]
			rct.Leaves[i] = leaf
		}
	}
	depth--
}

// Cut creates a child node to the left or right of the parent
func (rct *RCTree) Cut(X [][]float64, S []bool, parent *Node, side string) ([]bool, []bool, *Node) {
	subset := num.ArrayBoolFloat(X, S)
	// Find max and min over all d dimensions
	xmax := num.MaxColValues(subset)
	xmin := num.MinColValues(subset)

	// Compute l
	l := num.ArraySub(xmax, xmin)
	l = num.ArrayDivVal(l, num.ArraySumFloat(l))

	// Determine dimension to cut
	q := num.RndChoice(rct.Ndim, l)
	// q := rct.rnd.RndChoice(rct.Ndim, l)
	// Determine value for split
	p := num.RndUniform(xmin[q], xmax[q])
	// p := rct.rnd.RndUniform(xmin[q], xmax[q])

	rct.Log.Section("\nCut Tree\n")
	rct.Log.Add(fmt.Sprintf("l: %v\nq: %d  p: %f\n", l, q, p))

	// Determine subset of points to left
	arrayLeq := num.ArrayLeq(num.GetColumn(X, q), p)
	S1 := num.ArrayAnd(arrayLeq, S) // S1 is all points in S with random dimension < split value
	// Determine subset of points to right
	arrayNot := num.ArrayNot(S1)
	S2 := num.ArrayAnd(arrayNot, S) // S2 is all the points in S not in S1

	if num.ArrayCompareBool(S, S1) {
		rct.Log.Add("Warning: S equals S1\n")
	}

	// Create new child node
	child := NewBranch(q, p, nil, nil, parent, 0, nil)

	// Link child node to parent
	if parent != nil || side == "root" {
		switch side {
		case "l":
			parent.Branch.l = child
		case "r":
			parent.Branch.r = child
		case "root":
			rct.Root = child
		}
	}

	rct.Log.Write()
	return S1, S2, child
}

// ForgetPoint deletes a leaf from the tree
func (rct *RCTree) ForgetPoint(index int) *Node {
	// Get leaf from the leaves array
	node := rct.Leaves[index]
	// If duplicate points exist
	if node.n > 1 {
		// Decrement the number of points in the leaf and for all branches above
		rct.UpdateLeafCountUpwards(node, -1)
		return RemoveIndex(rct.Leaves, index)
	}

	// If node is the root
	if node.isRoot() {
		rct.Root = nil
		rct.Ndim = 0
		return RemoveIndex(rct.Leaves, index)
	}

	// Find parent
	parent := node.u
	// Find sibling
	sibling := parent.Branch.l
	if node == parent.Branch.l {
		sibling = parent.Branch.r
	}
	// If parent is the root
	if parent.isRoot() {
		// Set sibling as new root
		sibling.u = nil
		rct.Root = sibling
		// Update depths
		if sibling.isLeaf() {
			sibling.Leaf.d = 0
		} else {
			rct.MapDepths(sibling, -1)
		}
		return RemoveIndex(rct.Leaves, index)
	}
	// Find grandparent
	grandparent := parent.u
	// Set parent of sibling to grandparent
	sibling.u = grandparent
	// Short-circuit grandparent to sibling
	if parent == grandparent.Branch.l {
		grandparent.Branch.l = sibling
	} else {
		grandparent.Branch.r = sibling
	}
	// Update depths
	parent = grandparent
	rct.MapDepths(sibling, -1)
	// Update leaf counts under each branch
	rct.UpdateLeafCountUpwards(parent, -1)
	// Update bounding boxes
	point := node.Leaf.x
	rct.RelaxBboxUpwards(parent, point)
	return RemoveIndex(rct.Leaves, index)
}

// UpdateLeafCountUpwards updates the stored count of leaves beneath each branch (branch.n)
func (rct *RCTree) UpdateLeafCountUpwards(node *Node, inc int) {
	for node != nil {
		node.n += inc
		node = node.u
	}
}

// InsertPoint inserts a point into the tree, creating a new leaf
func (rct *RCTree) InsertPoint(point []float64, index int, tolerance float64) (*Node, error) {
	if rct.Root == nil {
		leafNode := NewLeaf(index, 0, nil, point, 1)
		rct.Root = leafNode
		rct.Ndim = len(point)
		rct.Leaves[index] = leafNode
		return leafNode, nil
	}
	// If leaves already exist in tree, check dimensions of point
	if len(point) != rct.Ndim {
		err := fmt.Errorf("Point dimension (%d) not equal to existing points in tree (%d)", len(point), rct.Ndim)
		return nil, err
	}
	// Check for existing index in leaves map
	if _, exists := rct.Leaves[index]; exists {
		err := fmt.Errorf("Index %d already exists in leaves map", index)
		return nil, err
	}
	// Check for duplicate points
	duplicate := rct.FindDuplicate(point, tolerance)
	if duplicate != nil {
		rct.UpdateLeafCountUpwards(duplicate, 1)
		rct.Leaves[index] = duplicate
		return duplicate, nil
	}
	// Tree has points and point is not a duplicate, so continue
	maxDepth := math.MinInt64
	for _, node := range rct.Leaves {
		if node.Leaf.d > maxDepth {
			maxDepth = node.Leaf.d
		}
	}

	depth := 0
	var branchNode *Node
	var leafNode *Node
	var side string

	currentNode := rct.Root
	parent := currentNode.u

	for range make([]int, maxDepth+1) {
		bbox := currentNode.b
		cutDimension, cut, _ := rct.InsertPointCut(point, bbox)

		if cut <= bbox[0][cutDimension] {
			leafNode = NewLeaf(index, depth, nil, point, 1)
			branchNode = NewBranch(cutDimension, cut, leafNode, currentNode, nil, leafNode.n+currentNode.n, nil)
			break
		} else if cut >= bbox[len(bbox)-1][cutDimension] {
			leafNode = NewLeaf(index, depth, nil, point, 1)
			branchNode = NewBranch(cutDimension, cut, currentNode, leafNode, nil, leafNode.n+currentNode.n, nil)
			break
		} else {
			depth++
			parent = currentNode
			if point[currentNode.Branch.q] <= currentNode.Branch.p {
				currentNode = currentNode.Branch.l
				side = "l"
			} else {
				currentNode = currentNode.Branch.r
				side = "r"
			}
		}
	}
	if branchNode == nil {
		err := fmt.Errorf("A cut was not found for index %d", index)
		return nil, err
	}

	// Set parent of new leaf and old branch
	currentNode.u = branchNode
	leafNode.u = branchNode

	// Set parent of new branch
	branchNode.u = parent
	if parent != nil {
		// Set child of parent to new branch
		switch side {
		case "l":
			parent.Branch.l = branchNode
		case "r":
			parent.Branch.r = branchNode
		}
	} else {
		// If a new root was created, assign the attribute
		rct.Root = branchNode
	}
	// Increment depths below branch
	rct.MapDepths(branchNode, 1)
	// Increment leaf count above branch
	rct.UpdateLeafCountUpwards(parent, 1)
	// Update bounding boxes
	rct.TightenBboxUpwards(branchNode)
	// Add leaf to leaves dict
	rct.Leaves[index] = leafNode
	// Return inserted leaf for convenience
	return leafNode, nil
}

// Query searches for leaf nearest to point
func (rct RCTree) Query(point []float64, node *Node) *Node {
	if node == nil {
		node = rct.Root
	}
	return queryPoint(point, node)
}

// queryPoint recursively searches for the nearest leaf to a given point
func queryPoint(point []float64, node *Node) *Node {
	if node.isLeaf() {
		return node
	}
	if point[node.Branch.q] <= node.Branch.p {
		return queryPoint(point, node.Branch.l)
	}
	return queryPoint(point, node.Branch.r)
}

// Disp computes displacement at leaf
func (rct RCTree) Disp(param interface{}) (int, error) {
	leaf, ok := param.(*Node)
	if !ok {
		index, ok := param.(int)
		if !ok {
			return 0, fmt.Errorf("CoDisp parameter not recognised: %v", leaf)
		}
		leaf = rct.Leaves[index]
	}

	// Handle case where leaf is root
	if leaf.isRoot() {
		return 0, nil
	}

	parent := leaf.u
	// Find sibling
	sibling := parent.Branch.l
	if leaf == parent.Branch.l {
		sibling = parent.Branch.r
	}
	// Count number of nodes in sibling subtree
	displacement := sibling.n
	return displacement, nil
}

// CoDisp computes collusive displacement (anomaly score) at leaf
func (rct RCTree) CoDisp(param interface{}) (float64, error) {
	leaf, ok := param.(*Node)
	if !ok {
		index, ok := param.(int)
		if !ok {
			return 0, fmt.Errorf("CoDisp parameter not recognised: %v", leaf)
		}
		leaf = rct.Leaves[index]
	}

	// Handle case where leaf is root
	if leaf.isRoot() {
		return 0, nil
	}
	node := leaf
	leafDepth := node.Leaf.d
	var results []float64

	for i := 0; i < leafDepth; i++ {
		parent := node.u
		if parent == nil {
			break
		}
		sibling := parent.Branch.l
		if node == parent.Branch.l {
			sibling = parent.Branch.r
		}
		numDeleted := node.n
		displacement := sibling.n
		result := float64(displacement) / float64(numDeleted)
		results = append(results, result)
		node = parent
	}
	coDisplacement := num.ArrayMaxValue(results)
	return coDisplacement, nil
}

// GetBbox computes the bounding box of all points underneath a given branch
func (rct *RCTree) GetBbox(branch *Node) [][]float64 {
	if branch == nil {
		branch = rct.Root
	}

	mins := num.Full(rct.Ndim, math.Inf(1))
	maxes := num.Full(rct.Ndim, math.Inf(-1))
	rct.MapBboxes(branch, mins, maxes)
	bbox := num.ArrayVStack(mins, maxes)

	return bbox
}

// FindDuplicate returns the leaf containing the duplicate of an existing point in the tree
// Returns nil if no duplicate found
func (rct *RCTree) FindDuplicate(point []float64, tolerance float64) *Node {
	nearest := rct.Query(point, nil)
	if tolerance == 0 {
		if num.ArrayCompareFloat(nearest.Leaf.x, point) {
			return nearest
		}
	} else {
		result := num.IsClose(nearest.Leaf.x, point, tolerance)
		if num.AllTrueBool(result) {
			return nearest
		}
	}
	return nil
}

// lrBranchBbox computes the bbox of a node based on bboxes of the node's children
func lrBranchBbox(branchNode *Node) [][]float64 {
	var bbLeft, bbRight, bbLastLeft, bbLastRight []float64

	node := branchNode.Branch.l
	if node.isBranch() {
		lastLeft := len(node.b) - 1
		bbLeft = node.b[0][:]
		bbLastLeft = node.b[lastLeft][:]
	} else {
		bbLeft = node.b[0][:]
		bbLastLeft = bbLeft
	}

	node = branchNode.Branch.r
	if node.isBranch() {
		lastRight := len(node.b) - 1
		bbRight = node.b[0][:]
		bbLastRight = node.b[lastRight][:]
	} else {
		bbRight = node.b[0][:]
		bbLastRight = bbRight
	}

	bbox := num.ArrayVStack(
		num.ArrayMinimum(bbLeft, bbRight),
		num.ArrayMaximum(bbLastLeft, bbLastRight))

	return bbox
}

// GetBboxTopDown recursively computes bboxes of all branches from root to leaves
func (rct *RCTree) GetBboxTopDown(node *Node) {
	if node.isBranch() {
		if node.Branch.l != nil {
			rct.GetBboxTopDown(node.Branch.l)
		}
		if node.Branch.r != nil {
			rct.GetBboxTopDown(node.Branch.r)
		}
		bbox := lrBranchBbox(node)
		node.b = bbox
	}
}

// CountAllTopDown recursively computes the number of leaves below each branch from root to leaves
func (rct *RCTree) CountAllTopDown(node *Node) {
	if node.isBranch() {
		if node.Branch.l != nil {
			rct.CountAllTopDown(node.Branch.l)
		}
		if node.Branch.r != nil {
			rct.CountAllTopDown(node.Branch.r)
		}
		node.n = node.Branch.l.n + node.Branch.r.n
	}
}

// CountLeaves counts the total leaves underneath a single node
func (rct *RCTree) CountLeaves(branch *Node) int {
	var numLeaves int

	rct.MapLeaves(branch, &numLeaves)
	return numLeaves
}

// SearchForLeaf -
func (rct *RCTree) SearchForLeaf() {

}

// TightenBboxUpwards expands bbox of all nodes above new point if point is outside the existing bbox
func (rct *RCTree) TightenBboxUpwards(node *Node) {
	bbox := lrBranchBbox(node)
	node.b = bbox
	node = node.u
	for node != nil {
		lastNode := len(node.b) - 1
		lastBbox := len(bbox) - 1
		lt := num.ArrayLt(bbox[0][:], node.b[0][:])
		gt := num.ArrayGt(bbox[lastBbox][:], node.b[lastNode][:])
		ltAny := num.AnyTrueBool(lt)
		gtAny := num.AnyTrueBool(gt)
		if ltAny || gtAny {
			if ltAny {
				num.ArrayCopyWhenTrue(node.b[0][:], bbox[0][:], lt)
			}
			if gtAny {
				num.ArrayCopyWhenTrue(node.b[lastNode][:], bbox[lastBbox][:], gt)
			}
		} else {
			break
		}
		node = node.u
	}
}

// RelaxBboxUpwards contracts bbox of all nodes above a deleted point
// if the deleted point defined the boundary of the bbox
func (rct *RCTree) RelaxBboxUpwards(node *Node, point []float64) {
	for node != nil {
		bbox := lrBranchBbox(node)
		lastIndex := len(node.b) - 1
		if !(num.AnyEqFloat(node.b[0][:], point) || num.AnyEqFloat(node.b[lastIndex][:], point)) {
			break
		}
		num.ArrayCopy(node.b[0][:], bbox[0][:])
		lastIndex = len(node.b) - 1
		lastBbox := len(bbox) - 1
		num.ArrayCopy(node.b[lastIndex][:], bbox[lastBbox][:])
		node = node.u
	}
}

// InsertPointCut generates the cut dimension and cut value based on InsertPoint()
func (rct *RCTree) InsertPointCut(point []float64, bbox [][]float64) (int, float64, error) {
	// Generate the bounding box
	bboxHat := num.ArrayEmpty(len(bbox), len(bbox[0]))
	// Update the bounding box based on the internal point
	lastBbox := len(bbox) - 1
	lastBboxHat := len(bboxHat) - 1
	minima := num.ArrayMinimum(bbox[0][:], point)
	num.ArrayCopy(bboxHat[0][:], minima)
	maxima := num.ArrayMaximum(bbox[lastBbox][:], point)
	num.ArrayCopy(bboxHat[lastBboxHat][:], maxima)
	bSpan := num.ArraySub(bboxHat[lastBboxHat][:], bboxHat[0][:])
	bRange := num.ArraySumFloat(bSpan)
	r := num.RndUniform(0, bRange)
	// r := rct.rnd.RndUniform(0, bRange)
	spanSum := num.ArrayCumSum(bSpan)
	cutDimension := math.MaxInt64
	for j := range make([]int, len(spanSum)) {
		if spanSum[j] >= r {
			cutDimension = j
			break
		}
	}
	if cutDimension == math.MaxInt64 {
		err := errors.New("Cut dimension is too large")
		return 0, 0, err
	}
	cut := bboxHat[0][cutDimension] + spanSum[cutDimension] - r
	return cutDimension, cut, nil
}
