package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
)

// WriteArray saves a 2D array of floats to a csv file
func WriteArray(data [][]float64, fileName string) {
	csvFile, _ := os.Create(fileName)
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
