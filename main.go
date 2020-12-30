package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"time"

	"sort"

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

// InitForest -
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
				// Produce a new array as sampled rows from X
				sampledX := array.Sample(data, ix)
				NewRCTree(token, sampledX, ix, 9, nil)
			}
		}
	}

	// Return the token
	return token
}

// NewEmptyForest -
func NewEmptyForest(token string) {
	numTrees := UserMap[token].NumTrees
	for treeIndex := 0; treeIndex < numTrees; treeIndex++ {
		NewRCTree(token, nil, nil, 0, nil)
	}
}

// NewRCTree -
func NewRCTree(token string, X [][]float64, indexLabels []int, precision int, randomState interface{}) {
	tree := rrcf.NewRCTree(X, indexLabels, precision, randomState)
	UserMap[token].Forest = append(UserMap[token].Forest, tree)
}

// GetTotalTrees -
func GetTotalTrees(token string) int {
	return len(UserMap[token].Forest)
}

// GetTotalLeaves -
func GetTotalLeaves(token string, treeIndex int) int {
	return len(UserMap[token].Forest[treeIndex].Leaves)
}

// InsertPoint -
func InsertPoint(token string, treeIndex int, point []float64, index int, tolerance float64) error {
	_, err := UserMap[token].Forest[treeIndex].InsertPoint(point, index, 0)
	if err == nil {
		UserMap[token].DataPoints++
	}
	return err
}

// ForgetPoint -
func ForgetPoint(token string, treeIndex int, index int) {
	UserMap[token].Forest[treeIndex].ForgetPoint(index)
}

// GetScore -
func GetScore(token string, treeIndex int, sampleIndex int) (float64, error) {
	return UserMap[token].Forest[treeIndex].CoDisp(sampleIndex)
}

// ScoreForest -
func ScoreForest(token string) []float64 {
	// Create a map to store anomaly score of each point
	avgScore := make([]float64, UserMap[token].DataPoints)

	index := make([]float64, UserMap[token].DataPoints)
	for _, tree := range UserMap[token].Forest {

		keys := []int{}
		for k := range tree.Leaves {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, key := range keys {
			codisp, _ := tree.CoDisp(key)
			avgScore[key] += codisp
			index[key]++
		}
	}
	for key := range avgScore {
		avgScore[key] /= index[key]
	}

	return avgScore
}

// GetThreshold calculates the threshold for the given percentile
func GetThreshold(token string, percentile float64) float64 {
	scores := ScoreForest(token)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i] < scores[j]
	})
	thresholdIndex := math.Round(float64(len(scores)) * percentile / 100)

	return scores[int(thresholdIndex)]
}

// UpdatePoint -
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

// UpdateBatch returns the average scores for each point, and the next sample index
func UpdateBatch(token string, sampleIndex int, points [][]float64) ([]float64, int) {
	index := sampleIndex
	var avgScore []float64

	for _, point := range points {
		avgScore = append(avgScore, UpdatePoint(token, index, point))
		index++
	}
	return avgScore, index
}
