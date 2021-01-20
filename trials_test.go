package main

import (
	"testing"

	"github.com/andysgithub/go-rrcf/utils"
)

func TestTrials(t *testing.T) {
	plotPoints := BatchTrial()
	utils.WriteToCsv(plotPoints, "results/batch/plot_points.csv")

	plotPoints = StreamingTrial()
	utils.WriteToCsv(plotPoints, "results/streaming/plot_points.csv")

	plotPoints = TrainingTrial()
	utils.WriteToCsv(plotPoints, "results/training/plot_points.csv")
}

// BatchTrial shows how the algorithm can be used to detect outliers in a batch setting
func BatchTrial() [][]float64 {
	// Get random 3D data with anomalies
	points, _ := utils.ReadFromCsv("data/random3D.csv")

	// Construct a random forest
	token := InitForest(100, 256, points, 0)

	// Compute average anomaly score
	scores := ScoreForest(token)

	// Calculate the threshold for the 99.5th percentile
	threshold := utils.GetThreshold(scores, 99.5)

	// Return points for plotting
	return utils.GetDataPoints(points, scores, threshold)
}

// StreamingTrial shows how the algorithm can be used to detect anomalies in streaming time series data
func StreamingTrial() [][]float64 {
	// Get sine function data with anomalies
	points, _ := utils.ReadFromCsv("data/sine.csv")

	// Construct a forest of empty trees
	token := InitForest(40, 256, nil, 3)

	// Create a map to store the anomaly score of each point
	scores := make(map[int]float64)

	// For each streamed data point
	for sampleIndex, point := range points {
		// Update the forest with this point and record the average score
		scores[sampleIndex] = UpdateForest(token, sampleIndex, point)
	}

	// Return points for plotting
	return utils.GetDataPoints(points, scores, 0)
}

// TrainingTrial shows how the forest can be pre-trained with anomaly-free data then continue streaming
func TrainingTrial() [][]float64 {
	// Get sine function data for training
	points, _ := utils.ReadFromCsv("data/training.csv")

	// Construct a forest of empty trees
	token := InitForest(40, 256, nil, 3)

	// For each training data point
	for sampleIndex, point := range points {
		// Update the forest with this point
		UpdateForest(token, sampleIndex, point)
	}
	lastIndex := len(points)

	// Get sine function data with anomalies
	points, _ = utils.ReadFromCsv("data/sine.csv")

	// Create a map to store the anomaly score of each point
	scores := make(map[int]float64)

	// For each streamed data point
	for sampleIndex, point := range points {
		// Update the forest with this point and record the average score
		scores[sampleIndex] = UpdateForest(token, lastIndex+sampleIndex, point)
	}

	// Return points for plotting
	return utils.GetDataPoints(points, scores, 0)
}
