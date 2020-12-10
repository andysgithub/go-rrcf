package rrcf

import (
	"bytes"
	"errors"
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
		make(map[int]*Node),
		nil, 0, nil, nil,
	}

	return rrcf
}

// Init - Initialises the random cut forest
func (rrcf *RRCF) Init(X [][]float64, indexLabels []int, precision int, randomState interface{}) {
	os.Remove("logs/rctree.log")
	file, err := os.OpenFile("logs/rctree.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		defer logMain.Close()
		logMain = file
	}

	switch randomState.(type) {
	case int:
		// Random number generation with provided seed
		rand.Seed((int64)(randomState.(int)))
	}

	if X != nil {
		// Round data to avoid sorting errors
		X = num.Around(X, precision)
		if indexLabels == nil {
			indexLabels = num.Arange(len(X))
		}
		rrcf.indexLabels = indexLabels

		// Remove duplicated rows
		X, I, N := num.Unique(X)

		dataRows := len(X)
		dataCols := len(X[0])

		bufferMain.Reset()
		bufferMain.WriteString("\n-------- Tree Build Start --------\n")
		bufferMain.WriteString(fmt.Sprintf("Rows: %d  Cols: %d\n", dataRows, dataCols))
		fmt.Fprintln(logMain, bufferMain.String())

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
}

// MakeTree generates a random cut tree
func (rrcf *RRCF) MakeTree(X [][]float64, S []bool, N []int, I []int, parent *Node, side string, depth int) {

	bufferMain.Reset()
	bufferMain.WriteString("\nCreate Branch\n")
	bufferMain.WriteString(fmt.Sprintf("Side: %s  Depth: %d\n", side, depth))
	fmt.Fprintln(logMain, bufferMain.String())

	// Increment depth as we traverse down
	depth++
	// Create a cut according to definition 1
	S1, S2, node := rrcf.Cut(X, S, parent, side)
	// If S1 does not contain an isolated point
	if num.ArraySumBool(S1) > 1 {
		// Recursively construct tree on S1
		rrcf.MakeTree(X, S1, N, I, node, "l", depth)
	} else {
		// Create a leaf node from the isolated point
		i := int(num.AsScalar(num.FlatNonZero(S1)))
		leaf := NewLeaf(i, depth, node, X[i][:], N[i])
		// Link leaf node to parent
		node.branch.l = leaf
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
		rrcf.MakeTree(X, S2, N, I, node, "r", depth)
	} else {
		// Create a leaf node from isolated point
		i := num.AsScalar(num.FlatNonZero(S2))
		leaf := NewLeaf(i, depth, node, X[i][:], N[i])
		// Link leaf node to parent
		node.branch.r = leaf
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

// Cut creates a child node to the left or right of the parent
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

	if num.ArrayCompare(S, S1) {
		bufferMain.Reset()
		bufferMain.WriteString("Warning: S equals S1\n")
		fmt.Fprintln(logMain, bufferMain.String())
	}

	// Create new child node
	child := NewBranch(q, p, nil, nil, parent, 0, nil)

	// Link child node to parent
	if parent != nil || side == "root" {
		switch side {
		case "l":
			parent.branch.l = child
		case "r":
			parent.branch.r = child
		case "root":
			rrcf.root = child
		}
	}

	return S1, S2, child
}

// ForgetPoint deletes a leaf from the tree
func (rrcf RRCF) ForgetPoint(index int) *Node {
	// Get leaf from the leaves array
	node := rrcf.leaves[index]
	// If duplicate points exist
	if node.n > 1 {
		// Decrement the number of points in the leaf and for all branches above
		rrcf.UpdateLeafCountUpwards(node, -1)
		return RemoveIndex(rrcf.leaves, index)
	}

	// If node is the root
	if node.isRoot() {
		rrcf.root = nil
		rrcf.ndim = 0
		return RemoveIndex(rrcf.leaves, index)
	}

	// Find parent
	parent := node.u
	// Find sibling
	sibling := parent.branch.l
	if node == parent.branch.l {
		sibling = parent.branch.r
	}
	// If parent is the root
	if parent.isRoot() {
		// Set sibling as new root
		sibling.u = nil
		rrcf.root = sibling
		// Update depths
		if sibling.isLeaf() {
			sibling.leaf.d = 0
		} else {
			rrcf.MapDepths(sibling, -1)
		}
		return RemoveIndex(rrcf.leaves, index)
	}
	// Find grandparent
	grandparent := parent.u
	// Set parent of sibling to grandparent
	sibling.u = grandparent
	// Short-circuit grandparent to sibling
	if parent == grandparent.branch.l {
		grandparent.branch.l = sibling
	} else {
		grandparent.branch.r = sibling
	}
	// Update depths
	parent = grandparent
	rrcf.MapDepths(sibling, -1)
	// Update leaf counts under each branch
	rrcf.UpdateLeafCountUpwards(parent, -1)
	// Update bounding boxes
	point := node.leaf.x
	rrcf.RelaxBboxUpwards(parent, point)
	return RemoveIndex(rrcf.leaves, index)
}

// UpdateLeafCountUpwards updates the stored count of leaves beneath each branch (branch.n)
func (rrcf RRCF) UpdateLeafCountUpwards(node *Node, inc int) {
	for node != nil {
		node.n += inc
		node = node.u
	}
}

// InsertPoint inserts a point into the tree, creating a new leaf
func (rrcf RRCF) InsertPoint(point []float64, index int, tolerance float64) (*Node, error) {
	if rrcf.root == nil {
		leafNode := NewLeaf(index, 0, nil, point, 1)
		rrcf.root = leafNode
		rrcf.ndim = len(point)
		rrcf.leaves[index] = leafNode
		return leafNode, nil
	}
	// If leaves already exist in tree, check dimensions of point
	if len(point) != rrcf.ndim {
		err := fmt.Errorf("Point dimension (%d) not equal to existing points in tree (%d)", len(point), rrcf.ndim)
		return nil, err
	}
	// Check for existing index in leaves map
	if _, exists := rrcf.leaves[index]; exists {
		err := fmt.Errorf("Index %d already exists in leaves map", index)
		return nil, err
	}
	// Check for duplicate points
	duplicate := rrcf.FindDuplicate(point, tolerance)
	if duplicate != nil {
		rrcf.UpdateLeafCountUpwards(duplicate, 1)
		rrcf.leaves[index] = duplicate
		return duplicate, nil
	}
	// Tree has points and point is not a duplicate, so continue
	maxDepth := math.MinInt64
	for _, node := range rrcf.leaves {
		if node.leaf.d > maxDepth {
			maxDepth = node.leaf.d
		}
	}

	depth := 0
	var branchNode *Node
	var leafNode *Node
	var side string

	currentNode := rrcf.root
	parent := currentNode.u

	for range make([]int, maxDepth+1) {
		if currentNode.isBranch() {
			bbox := currentNode.branch.b
			cutDimension, cut, _ := rrcf.InsertPointCut(point, bbox)
			if cut <= bbox[0][cutDimension] {
				leafNode = NewLeaf(index, depth, nil, point, 1)
				branchNode = NewBranch(cutDimension, cut, leafNode, currentNode, nil, leafNode.n+currentNode.n, nil)
				break
			} else if cut >= bbox[len(bbox)-1][cutDimension] {
				leafNode = NewLeaf(index, depth, nil, point, 1)
				branchNode = NewBranch(cutDimension, cut, leafNode, currentNode, nil, leafNode.n+currentNode.n, nil)
				break
			} else {
				depth++
				parent = currentNode
				if point[currentNode.branch.q] <= currentNode.branch.p {
					currentNode = currentNode.branch.l
					side = "l"
				} else {
					currentNode = currentNode.branch.r
					side = "r"
				}
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
			parent.branch.l = branchNode
		case "r":
			parent.branch.r = branchNode
		}
	} else {
		// If a new root was created, assign the attribute
		rrcf.root = branchNode
	}
	// Increment depths below branch
	rrcf.MapDepths(branchNode, 1)
	// Increment leaf count above branch
	rrcf.UpdateLeafCountUpwards(parent, 1)
	// Update bounding boxes
	rrcf.TightenBboxUpwards(branchNode)
	// Add leaf to leaves dict
	rrcf.leaves[index] = leafNode
	// Return inserted leaf for convenience
	return leafNode, nil
}

// Query searches for leaf nearest to point
func (rrcf RRCF) Query(point []float64, node *Node) *Node {
	if node == nil {
		node = rrcf.root
	}
	return queryPoint(point, node)
}

// queryPoint recursively searches for the nearest leaf to a given point
func queryPoint(point []float64, node *Node) *Node {
	if node.isLeaf() {
		return node
	}
	if point[node.branch.q] <= node.branch.p {
		return queryPoint(point, node.branch.l)
	}
	return queryPoint(point, node.branch.r)
}

// Disp computes displacement at leaf
func (rrcf RRCF) Disp(param interface{}) (int, error) {
	leaf, ok := param.(*Node)
	if !ok {
		index, ok := param.(int)
		if !ok {
			return 0, fmt.Errorf("CoDisp parameter not recognised: %v", leaf)
		}
		leaf = rrcf.leaves[index]
	}

	// Handle case where leaf is root
	if leaf.isRoot() {
		return 0, nil
	}

	parent := leaf.u
	// Find sibling
	sibling := parent.branch.l
	if leaf == parent.branch.l {
		sibling = parent.branch.r
	}
	// Count number of nodes in sibling subtree
	displacement := sibling.n
	return displacement, nil
}

// CoDisp computes collusive displacement at leaf
func (rrcf RRCF) CoDisp(param interface{}) (float64, error) {
	leaf, ok := param.(*Node)
	if !ok {
		index, ok := param.(int)
		if !ok {
			return 0, fmt.Errorf("CoDisp parameter not recognised: %v", leaf)
		}
		leaf = rrcf.leaves[index]
	}

	// Handle case where leaf is root
	if leaf.isRoot() {
		return 0, nil
	}
	node := leaf
	leafDepth := node.leaf.d
	var results []float64

	for i := 0; i < leafDepth; i++ {
		parent := node.u
		if parent == nil {
			break
		}
		sibling := parent.branch.l
		if node == parent.branch.l {
			sibling = parent.branch.r
		}
		numDeleted := node.n
		displacement := sibling.n
		result := float64(displacement / numDeleted)
		results = append(results, result)
		node = parent
	}
	coDisplacement := num.ArrayMaxValue(results)
	return coDisplacement, nil
}

// GetBbox computes the bounding box of all points underneath a given branch
func (rrcf RRCF) GetBbox(branch *Node) [][]float64 {
	if branch == nil {
		branch = rrcf.root
	}

	mins := num.Full(rrcf.ndim, math.Inf(1))
	maxes := num.Full(rrcf.ndim, math.Inf(-1))
	rrcf.MapBboxes(branch, mins, maxes)
	bbox := num.ArrayVStack(mins, maxes)

	return bbox
}

// FindDuplicate returns the leaf containing the duplicate of an existing point in the tree
// Returns nil if no duplicate found
func (rrcf RRCF) FindDuplicate(point []float64, tolerance float64) *Node {
	nearest := rrcf.Query(point, nil)
	if tolerance == 0 {
		if num.ArrayCompareFloat(nearest.leaf.x, point) {
			return nearest
		}
	} else {
		result := num.IsClose(nearest.leaf.x, point, tolerance)
		if num.AllTrue(result) {
			return nearest
		}
	}
	return nil
}

// ToDict -
func (rrcf RRCF) ToDict() {

}

// Serialise -
func (rrcf RRCF) Serialise() {

}

// LoadDict -
func (rrcf RRCF) LoadDict() {

}

// Deserialise -
func (rrcf RRCF) Deserialise() {

}

// FromDict -
func (rrcf RRCF) FromDict() {

}

// lrBranchBbox computes the bbox of a node based on bboxes of the node's children
func lrBranchBbox(branchNode *Node) [][]float64 {
	var bbLeft, bbRight, bbLastLeft, bbLastRight []float64

	node := branchNode.branch.l
	if node.isBranch() {
		lastLeft := len(node.branch.b) - 1
		bbLeft = node.branch.b[0][:]
		bbLastLeft = node.branch.b[lastLeft][:]
	} else {
		bbLeft = node.leaf.b[:]
		bbLastLeft = bbLeft
	}

	node = branchNode.branch.r
	if node.isBranch() {
		lastRight := len(node.branch.b) - 1
		bbRight = node.branch.b[0][:]
		bbLastRight = node.branch.b[lastRight][:]
	} else {
		bbRight = node.leaf.b[:]
		bbLastRight = bbRight
	}

	bbox := num.ArrayVStack(
		num.ArrayMinimum(bbLeft, bbRight),
		num.ArrayMaximum(bbLastLeft, bbLastRight))

	return bbox
}

// GetBboxTopDown recursively computes bboxes of all branches from root to leaves
func (rrcf RRCF) GetBboxTopDown(node *Node) {
	if node.isBranch() {
		if node.branch.l != nil {
			rrcf.GetBboxTopDown(node.branch.l)
		}
		if node.branch.r != nil {
			rrcf.GetBboxTopDown(node.branch.r)
		}
		bbox := lrBranchBbox(node)
		node.branch.b = bbox
	}
}

// CountAllTopDown recursively computes the number of leaves below each branch from root to leaves
func (rrcf RRCF) CountAllTopDown(node *Node) {
	if node.isBranch() {
		if node.branch.l != nil {
			rrcf.CountAllTopDown(node.branch.l)
		}
		if node.branch.r != nil {
			rrcf.CountAllTopDown(node.branch.r)
		}
		node.n = node.branch.l.n + node.branch.r.n
	}
}

// CountLeaves counts the total leaves underneath a single node
func (rrcf RRCF) CountLeaves(branch *Node) int {
	var numLeaves int

	rrcf.MapLeaves(branch, &numLeaves)
	return numLeaves
}

// SearchForLeaf -
func (rrcf RRCF) SearchForLeaf() {

}

// TightenBboxUpwards expands bbox of all nodes above new point if point is outside the existing bbox
func (rrcf RRCF) TightenBboxUpwards(node *Node) {
	bbox := lrBranchBbox(node)
	node.branch.b = bbox
	node = node.u
	for node != nil {
		lastNode := len(node.branch.b) - 1
		lastBbox := len(bbox) - 1
		lt := num.ArrayLt(bbox[0][:], node.branch.b[0][:])
		gt := num.ArrayGt(bbox[lastBbox][:], node.branch.b[lastNode][:])
		ltAny := num.AnyTrueBool(lt)
		gtAny := num.AnyTrueBool(gt)
		if ltAny || gtAny {
			if ltAny {
				num.ArrayCopyWhenTrue(node.branch.b[0][:], bbox[0][:], lt)
			}
			if gtAny {
				num.ArrayCopyWhenTrue(node.branch.b[lastNode][:], bbox[lastBbox][:], gt)
			}
		} else {
			break
		}
		node = node.u
	}
}

// RelaxBboxUpwards contracts bbox of all nodes above a deleted point
// if the deleted point defined the boundary of the bbox
func (rrcf RRCF) RelaxBboxUpwards(node *Node, point []float64) {
	for node != nil {
		bbox := lrBranchBbox(node)
		lastIndex := len(node.branch.b) - 1
		if !(num.AnyTrue(node.branch.b[0][:], point) || num.AnyTrue(node.branch.b[lastIndex][:], point)) {
			break
		}
		num.ArrayCopy(node.branch.b[0][:], bbox[0][:])
		lastIndex = len(node.branch.b) - 1
		lastBbox := len(bbox) - 1
		num.ArrayCopy(node.branch.b[lastIndex][:], bbox[lastBbox][:])
		node = node.u
	}
}

// InsertPointCut generates the cut dimension and cut value based on InsertPoint()
func (rrcf RRCF) InsertPointCut(point []float64, bbox [][]float64) (int, float64, error) {
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
