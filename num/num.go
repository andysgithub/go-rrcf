package num

import (
	"fmt"
	"math"
	"strings"
)

// Around evenly rounds to the given number of decimals
func Around(X [][]float64, decimals int) [][]float64 {
	multiplier := math.Pow10(decimals)
	for i, array := range X {
		for j, val := range array {
			rounded := math.Round(float64(val) * multiplier)
			X[i][j] = rounded / multiplier
		}
	}
	return X
}

// Arange returns evenly spaced integers within a given interval
func Arange(interval int) []int {
	a := make([]int, interval)
	for i := range a {
		a[i] = i
	}
	return a
}

// Unique returns an array with unique rows
// U: The resulting array with unique rows
// I: Row indices which were removed
// N: The number of times each of the unique rows appears in X
func Unique(X [][]float64) ([][]float64, []int, []int) {
	var U [][]float64
	var I []int
	var N []int
	var rowKeys []string

	// Make a map to record the contents of each recorded row
	keys := make(map[string]int)
	indexU := 0
	for i, values := range X {
		rowKey := SliceToString(values, ",")
		// If this row of values has not been recorded yet
		if _, value := keys[rowKey]; !value {
			// Set the total occurrences for this row to 1
			keys[rowKey] = 1
			rowKeys = append(rowKeys, rowKey)
			U = append(U, values)
			indexU++
		} else {
			// Record the index of this duplicated row
			I = append(I, i)
			// Increment the total occurrences for this row
			keys[rowKey]++
		}
	}
	// Compile the list of unique row totals
	for _, rowKey := range rowKeys {
		N = append(N, keys[rowKey])
	}
	return U, I, N
}

func SliceToString(a []float64, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = fmt.Sprintf("%f", v)
	}
	return strings.Join(b, sep)
}
