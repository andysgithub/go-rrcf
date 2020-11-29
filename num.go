package num

import (
	"math"
)

// Around evenly rounds to the given number of decimals
func Around(X [][]int64, decimals int) [][]int64 {
	multiplier := math.Pow10(decimals)
	for i, array := range X {
		for j, val := range array {
			rounded := math.Round(float64(val) * multiplier)
			X[i][j] = int64(rounded / multiplier)
		}
	}
	return X
}

// Arange returns evenly spaced integers within a given interval
func Arange(interval int) []int {
	a := make([]int, interval+1)
	for i := range a {
		a[i] = i
	}
	return a
}
