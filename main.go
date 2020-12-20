package main

import (
	"math"
	"sort"

	"github.com/andysgithub/go-rrcf/array"
	"github.com/andysgithub/go-rrcf/random"
	"github.com/andysgithub/go-rrcf/rrcf"
	"github.com/andysgithub/go-rrcf/utils"
)

func main() {
	plotPoints := StreamingTrial()
	utils.WriteArray(plotPoints, "results/streaming/plot_points.csv")

	plotPoints = BatchTrial()
	utils.WriteArray(plotPoints, "results/batch/plot_points.csv")
}

// StreamingTrial shows how the algorithm can be used to detect anomalies in streaming time series data
func StreamingTrial() [][]float64 {
	// Generate data
	n := 730
	A := 50
	center := 100
	phi := 30
	T := 2.0 * math.Pi / 100
	t := array.Arange(n)

	diff := array.SubtractVal1D(array.MultiplyVal1DInt(t, T), float64(phi)*T)
	mul := array.MultiplyVal1D(array.Sin(diff), float64(A))
	sin := array.AddVal1D(mul, float64(center))

	array.FillElements(sin, 235, 255, float64(80))

	// Construct forest of empty RCTrees

	// Set tree parameters
	numTrees := 40
	shingleSize := 4
	treeSize := 256

	// Create a forest of empty trees
	var forest []rrcf.RCTree
	for range make([]int, numTrees) {
		tree := rrcf.NewRCTree(nil, nil, 9, nil)
		forest = append(forest, tree)
	}

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

	// Compile points for plotting
	plotPoints := array.Zero2D(totalPoints, 2)

	for i := 0; i < totalPoints; i++ {
		plotPoints[i][0] = sin[i]
		plotPoints[i][1] = avgCodisp[i]
	}

	return plotPoints
}

// BatchTrial shows how the algorithm can be used to detect outliers in a batch setting
func BatchTrial() [][]float64 {
	// Set sample parameters
	rnd := random.NewRandomState(0)
	n := 2010
	d := 3

	// Generate data
	X := array.Zero2D(n, d)
	array.FillColumn(X, 0, 0, 999, 5.)
	array.FillColumn(X, 0, 1000, 1999, -5.)

	randArray := rnd.Normal2D(len(X), len(X[0]))
	mulArray := array.MultiplyVal2D(randArray, 0.01)
	X = array.AddFloat2D(X, mulArray)

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
		ixs := rnd.Array(n, rows, cols)
		for _, ix := range ixs[0 : rows-1] {
			// Produce a new array as sampled rows from X
			sampledX := array.Sample(X, ix)
			tree := rrcf.NewRCTree(sampledX, ix, 9, nil)
			forest = append(forest, tree)
		}
	}

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

	// Calculate the threshold for the 99.5th percentile
	values := utils.SortMap(avgCodisp)
	threshold := values[len(values)-10]

	// Compile points for plotting
	totalPoints := len(X)
	plotPoints := array.Zero2D(totalPoints, 5)

	for i := 0; i < totalPoints; i++ {
		plotPoints[i][0] = X[i][0]
		plotPoints[i][1] = X[i][1]
		plotPoints[i][2] = X[i][2]
		plotPoints[i][3] = avgCodisp[i]
		plotPoints[i][4] = utils.BoolToFloat(avgCodisp[i] >= threshold)
	}

	return plotPoints
}
