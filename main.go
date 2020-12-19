package main

import (
	"math"
	"sort"

	"github.com/andysgithub/go-rrcf/num"
	"github.com/andysgithub/go-rrcf/random"
	"github.com/andysgithub/go-rrcf/rrcf"
	"github.com/andysgithub/go-rrcf/utils"
)

func main() {
	plotPoints := StreamingTrial()
	utils.WriteArray(plotPoints, "results/streaming/plot_points.csv")

	// plotPoints = BatchTrial()
	// utils.WriteArray(plotPoints, "results/batch/plot_points.csv")
}

// StreamingTrial shows how the algorithm can be used to detect anomalies in streaming time series data
func StreamingTrial() [][]float64 {
	// Generate data
	n := 730
	A := 50
	center := 100
	phi := 30
	T := 2.0 * math.Pi / 100
	t := num.Arange(n)

	diff := num.ArraySubVal(num.ArrayMulValInt(t, T), float64(phi)*T)
	mul := num.ArrayMulVal(num.ArraySin(diff), float64(A))
	sin := num.ArrayAddVal(mul, float64(center))

	num.ArrayFillElements(sin, 235, 255, float64(80))

	// Construct forest of empty RCTrees

	// Set tree parameters
	numTrees := 40
	shingleSize := 4
	treeSize := 256

	// Create a forest of empty trees
	var forest []rrcf.RCTree
	for i := range make([]int, numTrees) {
		tree := rrcf.NewRCTree(nil, nil, 9, nil, i)
		forest = append(forest, tree)
	}

	// rrcf.SaveForest(forest, "data/forest/stream.json")

	// Insert streaming points into tree and compute anomaly score
	// Use the "shingle" generator to create a rolling window
	shingle := rrcf.NewShingleList(sin, shingleSize)
	totalPoints := len(sin) - shingleSize

	// Create a map to store anomaly score of each point
	avgCodisp := make(map[int]float64)

	// For each shingle
	for index := 0; index < totalPoints; index++ {
		point := shingle.NextInList()
		// For each tree in the forest
		for i := range forest {
			// If tree is above permitted size
			if len(forest[i].Leaves) > treeSize {
				// Drop the oldest point (FIFO)
				forest[i].ForgetPoint(index - treeSize)
			}
			// Insert the new point into the tree
			forest[i].InsertPoint(point, index, 0)
			// Compute codisp on the new point
			newCodisp, _ := forest[i].CoDisp(index)
			// Take the average over all trees
			if _, ok := avgCodisp[index]; !ok {
				avgCodisp[index] = 0
			}
			avgCodisp[index] += newCodisp / float64(numTrees)
		}
	}

	// Close logs
	for _, tree := range forest {
		tree.Log.Write()
		tree.Log.Close()
	}

	// Compile points for plotting
	plotPoints := num.ArrayEmpty(totalPoints, 2)

	for i := 0; i < totalPoints; i++ {
		plotPoints[i][0] = sin[i]
		plotPoints[i][1] = avgCodisp[i]
	}

	return plotPoints
}

// BatchTrial shows how the algorithm can be used to detect outliers in a batch setting
func BatchTrial() [][]float64 {
	rnd := random.NewRandom()

	// Set sample parameters
	rnd.Seed(int64(0))
	n := 2010
	d := 3

	// Generate data
	X := num.ArrayEmpty(n, d)
	num.ArrayFillColumn(X, 0, 0, 999, 5.)
	num.ArrayFillColumn(X, 0, 1000, 1999, -5.)

	randArray := num.Randl2(len(X), len(X[0]))
	//randArray := rnd.Randn2(len(X), len(X[0]))
	mulArray := num.Array2DMulVal(randArray, 0.01)
	X = num.Array2DAdd(X, mulArray)

	utils.WriteArray(X, "data/batch.csv")

	// Construct a random forest

	// Set forest parameters
	numTrees := 100
	treeSize := 256
	sampleSizeRange := []int{int(n / treeSize), treeSize}

	// Construct forest
	var forest []rrcf.RCTree

	for i := 0; len(forest) < numTrees; i++ {
		// Select random subsets of points uniformly
		rows := sampleSizeRange[0]
		cols := sampleSizeRange[1]
		ixs := num.RndArray(n, rows, cols)
		//ixs := rnd.RndArray(n, rows, cols)
		for _, ix := range ixs[0 : rows-1] {
			// Produce a new array as sampled rows from X
			sampledX := num.ArraySample(X, ix)
			tree := rrcf.NewRCTree(sampledX, ix, 9, nil, i)
			forest = append(forest, tree)
		}
	}

	// rrcf.SaveForest(forest, "data/forest/batch.json")

	// Create a map to store anomaly score of each point
	avgCodisp := make(map[int]float64)
	for i := 0; i < n; i++ {
		avgCodisp[i] = 0.0
	}

	// Compute average CoDisp
	index := make([]float64, n)
	for _, tree := range forest {

		keys := []int{}
		for k := range tree.Leaves {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, key := range keys {
			codisp, _ := tree.CoDisp(key)
			avgCodisp[key] += codisp
			index[key]++
		}
	}
	for key := range avgCodisp {
		if index[key] == 0 {
			delete(avgCodisp, key)
		} else {
			avgCodisp[key] /= index[key]
		}
	}

	// Compile points for plotting
	values := sortMap(avgCodisp)
	threshold := values[len(values)-10]

	totalPoints := len(X)
	plotPoints := num.ArrayEmpty(totalPoints, 5)

	for i := 0; i < totalPoints; i++ {
		plotPoints[i][0] = X[i][0]
		plotPoints[i][1] = X[i][1]
		plotPoints[i][2] = X[i][2]
		plotPoints[i][3] = avgCodisp[i]
		plotPoints[i][4] = boolToFloat(avgCodisp[i] >= threshold)
	}

	return plotPoints
}

func sortMap(m map[int]float64) []float64 {
	// Convert map to slice of values.
	values := []float64{}
	for _, value := range m {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.
	}
	return 0.
}
