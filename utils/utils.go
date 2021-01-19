package utils

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/andysgithub/go-rrcf/array"
)

// ReadFromCsv will read the csv file at filePath as a 2d array of floats
func ReadFromCsv(filePath string) ([][]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	stringValues, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	values, err := StringsToFloats(stringValues)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// WriteToCsv saves a 2D array of floats to a csv file
func WriteToCsv(data [][]float64, fileName string) error {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	csvwriter := csv.NewWriter(csvFile)

	for _, dataRow := range data {

		var stringRow []string
		for _, value := range dataRow {
			strVal := fmt.Sprintf("%f", value)
			stringRow = append(stringRow, strVal)
		}

		csvwriter.Write(stringRow)
	}

	csvwriter.Flush()
	return nil
}

// SortMap converts a map to a slice of values and returns the sorted slice
func SortMap(m map[int]float64) []float64 {
	values := []float64{}
	for _, value := range m {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}

// BoolToFloat converts a boolean to 1.0 (true) or 0.0 (false)
func BoolToFloat(b bool) float64 {
	if b {
		return 1.
	}
	return 0.
}

// StringsToFloats converts a 2d array of strings into a 2d array of floats
func StringsToFloats(stringValues [][]string) ([][]float64, error) {
	values := array.Zero2D(len(stringValues), len(stringValues[0]))

	for rowIndex := range values {
		for colIndex := range values[rowIndex] {
			var err error = nil

			trimString :=
				strings.TrimSpace(stringValues[rowIndex][colIndex])

			values[rowIndex][colIndex], err =
				strconv.ParseFloat(trimString, 64)

			if err != nil {
				fmt.Println(err)
				return values, err
			}
		}
	}

	return values, nil
}

// GetDataPoints compiles data points and score values into a 2D array of floats
func GetDataPoints(points [][]float64, score map[int]float64, threshold float64) [][]float64 {
	dataCols := len(points[0])
	scoreCols := 1
	if threshold > 0 {
		scoreCols = 2
	}

	plotPoints := array.Zero2D(len(points), dataCols+scoreCols)

	for row := range points {
		// Store the data points for this row
		for col, value := range points[row] {
			plotPoints[row][col] = value
		}
		// Next column is the score for this row
		plotPoints[row][dataCols] = score[row]
		if threshold > 0 {
			// Last column is above/below threshold (1/0)
			plotPoints[row][dataCols+1] = BoolToFloat(score[row] >= threshold)
		}
	}
	return plotPoints
}

// GetThreshold calculates the threshold for the given percentile
func GetThreshold(scores map[int]float64, percentile float64) float64 {
	// Sort the scores into numerical order
	values := SortMap(scores)

	thresholdIndex := math.Round(float64(len(values)) * percentile / 100)

	return values[int(thresholdIndex)]
}
