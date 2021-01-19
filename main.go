package main

import (
	"crypto/rand"
	"fmt"
	"sort"
	"time"

	"github.com/andysgithub/go-rrcf/array"
	"github.com/andysgithub/go-rrcf/random"
	"github.com/andysgithub/go-rrcf/rrcf"
)

// UserMap is a map of token/user pairs
var UserMap map[string]*User

// User struct records the RRCF details for one user
type User struct {
	Forest      []rrcf.RCTree
	NumTrees    int
	TreeSize    int
	DataPoints  int
	ShingleSize int
	Shingle     []float64
}

func main() {
}

// InitForest initialises a forest from the given source data
// Returns a token to reference the forest for use in subsequent calls
func InitForest(numTrees int, treeSize int, data [][]float64, shingleSize int) string {
	if UserMap == nil {
		UserMap = make(map[string]*User)
	}

	// Generate a key token
	b := make([]byte, 4)
	rand.Read(b)
	token := fmt.Sprintf("%x", b)

	dataPoints := 0
	if data != nil {
		dataPoints = len(data)
	}

	// Add key token to user map
	UserMap[token] = &User{
		NumTrees:    numTrees,
		TreeSize:    treeSize,
		DataPoints:  dataPoints,
		ShingleSize: shingleSize,
	}

	if dataPoints == 0 {
		NewEmptyForest(token)
	} else {
		sampleSizeRange := []int{int(dataPoints / treeSize), treeSize}
		rnd := random.NewRandomState(time.Now().UTC().UnixNano())

		for i := 0; GetTotalTrees(token) < numTrees; i++ {
			// Select random subsets of points uniformly
			rows := sampleSizeRange[0]
			cols := sampleSizeRange[1]
			ixs := rnd.Array(dataPoints, rows, cols)
			for _, ix := range ixs[0 : rows-1] {
				// Produce a new array as sampled rows from source data
				sampledX := array.Sample(data, ix)
				NewRCTree(token, sampledX, ix, 9, nil)
			}
		}
	}

	// Return the token
	return token
}

// UpdateForest maintains a shingle internally by retaining previous data points
func UpdateForest(token string, sampleIndex int, point []float64) float64 {
	data := point

	if len(point) == 1 {
		// Only one data point, so use shingles
		shingleSize := UserMap[token].ShingleSize
		data = UserMap[token].Shingle

		data = append(data, point[0])
		if len(data) > shingleSize {
			data = data[1:]
		}
		UserMap[token].Shingle = data

		if len(data) < shingleSize {
			return 0
		}
	}

	return UpdatePoint(token, sampleIndex, data)
}

// ScoreForest calculates the average score at each leaf across all trees
func ScoreForest(token string) map[int]float64 {
	// Create a map to store average scores at each leaf
	avgScores := make(map[int]float64)
	// Create a map to store the total occurences of each leaf index in the forest
	leafTotals := make(map[int]float64)

	for _, tree := range UserMap[token].Forest {
		keys := []int{}
		for k := range tree.Leaves {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, key := range keys {
			codisp, _ := tree.CoDisp(key)
			avgScores[key] += codisp
			leafTotals[key]++
		}
	}
	for key := range avgScores {
		avgScores[key] /= leafTotals[key]
	}

	return avgScores
}

// NewEmptyForest creates a forest of empty trees
func NewEmptyForest(token string) {
	numTrees := UserMap[token].NumTrees
	for treeIndex := 0; treeIndex < numTrees; treeIndex++ {
		NewRCTree(token, nil, nil, 0, nil)
	}
}

// NewRCTree creates a new tree and appends it to the forest
func NewRCTree(token string, X [][]float64, indexLabels []int, precision int, randomState interface{}) {
	tree := rrcf.NewRCTree(X, indexLabels, precision, randomState)
	UserMap[token].Forest = append(UserMap[token].Forest, tree)
}

// UpdatePoint inserts a new point into each tree and updates the score
func UpdatePoint(token string, sampleIndex int, point []float64) float64 {
	treeSize := UserMap[token].TreeSize
	numTrees := UserMap[token].NumTrees
	var avgScore float64

	// For each tree in the forest
	for treeIndex := 0; treeIndex < numTrees; treeIndex++ {
		// If tree is above permitted size
		if GetTotalLeaves(token, treeIndex) > treeSize {
			// Drop the oldest point (FIFO)
			ForgetPoint(token, treeIndex, sampleIndex-treeSize)
		}
		// Insert the new point into the tree
		InsertPoint(token, treeIndex, point, sampleIndex, 0)

		// Compute codisp on the new point
		newScore, _ := GetScore(token, treeIndex, sampleIndex)
		// Take the average over all trees
		avgScore += newScore / float64(numTrees)
	}
	return avgScore
}

// InsertPoint inserts a point into a tree, creating a new leaf
func InsertPoint(token string, treeIndex int, point []float64, index int, tolerance float64) error {
	_, err := UserMap[token].Forest[treeIndex].InsertPoint(point, index, 0)
	if err == nil {
		UserMap[token].DataPoints++
	}
	return err
}

// ForgetPoint deletes a leaf from the specified tree
func ForgetPoint(token string, treeIndex int, index int) {
	UserMap[token].Forest[treeIndex].ForgetPoint(index)
}

// GetTotalTrees returns the total number of trees in the forest
func GetTotalTrees(token string) int {
	return len(UserMap[token].Forest)
}

// GetTotalLeaves returns the number of leaves in the specified tree
func GetTotalLeaves(token string, treeIndex int) int {
	return len(UserMap[token].Forest[treeIndex].Leaves)
}

// GetScore returns the collusive displacement for a leaf in the specified tree
func GetScore(token string, treeIndex int, sampleIndex int) (float64, error) {
	return UserMap[token].Forest[treeIndex].CoDisp(sampleIndex)
}
