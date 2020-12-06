package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// WriteFile saves a 2D array of floats to a csv file
func WriteFile(rndData [][]float64, fileName string) {
	csvFile, _ := os.Create(fileName)
	csvwriter := csv.NewWriter(csvFile)

	for _, dataRow := range rndData {

		var stringRow []string
		for _, value := range dataRow {
			strVal := fmt.Sprintf("%f", value)
			stringRow = append(stringRow, strVal)
		}

		_ = csvwriter.Write(stringRow)
	}

	csvwriter.Flush()
}

func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
