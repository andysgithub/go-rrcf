package main

import (
	"C"
	"crypto/rand"
	"fmt"

	"github.com/andysgithub/go-rrcf/rrcf"
)
import "sort"

// UserMap is a map of token/user pairs
var UserMap map[string]*User

// User struct re3cords the RRCF details for one user
type User struct {
	Forest     []rrcf.RCTree
	NumTrees   int
	TreeSize   int
	DataPoints int
}

func main() {
}

//export InitRRCF
func InitRRCF(numTrees int, treeSize int, dataPoints int) string {
	if UserMap == nil {
		UserMap = make(map[string]*User)
	}

	// Generate a key token
	b := make([]byte, 4)
	rand.Read(b)
	token := fmt.Sprintf("%x", b)

	// Add key token to user map
	UserMap[token] = &User{
		NumTrees:   numTrees,
		TreeSize:   treeSize,
		DataPoints: dataPoints,
	}

	// Return the token
	return token
}

//export NewEmptyForest
func NewEmptyForest(token string) {
	numTrees := UserMap[token].NumTrees
	for treeIndex := 0; treeIndex < numTrees; treeIndex++ {
		NewRCTree(token, nil, nil, 0, nil)
	}
}

//export NewRCTree
func NewRCTree(token string, X [][]float64, indexLabels []int, precision int, randomState interface{}) {
	tree := rrcf.NewRCTree(X, indexLabels, precision, randomState)
	UserMap[token].Forest = append(UserMap[token].Forest, tree)
}

//export GetTotalTrees
func GetTotalTrees(token string) int {
	return len(UserMap[token].Forest)
}

//export GetTotalLeaves
func GetTotalLeaves(token string, treeIndex int) int {
	return len(UserMap[token].Forest[treeIndex].Leaves)
}

//export InsertPoint
func InsertPoint(token string, treeIndex int, point []float64, index int, tolerance float64) error {
	_, err := UserMap[token].Forest[treeIndex].InsertPoint(point, index, 0)
	if err == nil {
		UserMap[token].DataPoints++
	}
	return err
}

//export ForgetPoint
func ForgetPoint(token string, treeIndex int, index int) {
	UserMap[token].Forest[treeIndex].ForgetPoint(index)
}

//export GetScore
func GetScore(token string, treeIndex int, sampleIndex int) (float64, error) {
	return UserMap[token].Forest[treeIndex].CoDisp(sampleIndex)
}

//export GetAverageScore
func GetAverageScore(token string) map[int]float64 {
	// Create a map to store anomaly score of each point
	avgScore := make(map[int]float64)
	for i := 0; i < UserMap[token].DataPoints; i++ {
		avgScore[i] = 0.0
	}

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
		if index[key] == 0 {
			delete(avgScore, key)
		} else {
			avgScore[key] /= index[key]
		}
	}

	return avgScore
}

//export UpdatePoint
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

//export UpdateBatch
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
