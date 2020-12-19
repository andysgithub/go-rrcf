package utils

import (
	"encoding/csv"
	"fmt"
	"os"
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

// Min -
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max -
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
