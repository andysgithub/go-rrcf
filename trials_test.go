package main

import (
	"testing"

	"github.com/andysgithub/go-rrcf/array"
	"github.com/andysgithub/go-rrcf/random"
	"github.com/andysgithub/go-rrcf/utils"
)

func TestTrials(t *testing.T) {
	plotPoints := StreamingTrial()
	utils.WriteArray(plotPoints, "results/streaming/plot_points.csv")

	// plotPoints := BatchTrial()
	// utils.WriteArray(plotPoints, "results/batch/plot_points.csv")
}

// StreamingTrial shows how the algorithm can be used to detect anomalies in streaming time series data
func StreamingTrial() [][]float64 {
	points, err := utils.ReadFromCsv("data/sin.csv")
	if err != nil {
		return nil
	}

	// Construct a forest of empty trees
	token := InitRRCF(40, 256, 0, 3)
	NewEmptyForest(token)

	// Create a map to store anomaly score of each point
	avgScore := make(map[int]float64)

	// For each point
	for sampleIndex, point := range points {
		// Update the forest with this point
		avgScore[sampleIndex] = UpdateForest(token, sampleIndex, point)
	}

	// Compile points for plotting
	plotPoints := array.Zero2D(len(points), 2)
	for sampleIndex, point := range points {
		plotPoints[sampleIndex][0] = point[0]
		plotPoints[sampleIndex][1] = avgScore[sampleIndex]
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
	token := InitRRCF(numTrees, treeSize, n, 0)

	for i := 0; GetTotalTrees(token) < numTrees; i++ {
		// Select random subsets of points uniformly
		rows := sampleSizeRange[0]
		cols := sampleSizeRange[1]
		ixs := rnd.Array(n, rows, cols)
		for _, ix := range ixs[0 : rows-1] {
			// Produce a new array as sampled rows from X
			sampledX := array.Sample(X, ix)
			NewRCTree(token, sampledX, ix, 9, nil)
		}
	}

	// Compute average anomaly score
	avgScore := GetAverageScore(token)

	// Calculate the threshold for the 99.5th percentile
	values := utils.SortMap(avgScore)
	threshold := values[len(values)-10]

	// Compile points for plotting
	totalPoints := len(X)
	plotPoints := array.Zero2D(totalPoints, 5)

	for i := 0; i < totalPoints; i++ {
		plotPoints[i][0] = X[i][0]
		plotPoints[i][1] = X[i][1]
		plotPoints[i][2] = X[i][2]
		plotPoints[i][3] = avgScore[i]
		plotPoints[i][4] = utils.BoolToFloat(avgScore[i] >= threshold)
	}

	return plotPoints
}
