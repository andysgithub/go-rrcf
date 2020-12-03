package num

import (
	"fmt"
	"math"
	"math/rand"
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

func Ones_bool(rows int) []bool {
	newArray := make([]bool, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = true
	}
	return newArray
}

// AsScalar converts an array of size 1 to its scalar equivalent
func AsScalar(element []int) int {
	return element[0]
}

// FlatNonZero returns indices that are non-zero in the flattened version of the array
func FlatNonZero(array []bool) []int {
	var nonZero []int
	index := 0

	for _, element := range array {
		if element == true {
			nonZero = append(nonZero, index)
		}
		index++
	}
	return nonZero
}

// Randn generates a 2D array of normally distributed random floats
func Randn(rows, cols int) [][]float64 {
	newArray := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			newArray[i][j] = rand.NormFloat64()
		}
	}
	return newArray
}

// CopyArray copies a 2D array of floats by value and returns the copy
func CopyArray(array [][]float64) [][]float64 {
	rows := len(array)
	cols := len(array[0])

	duplicate := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		duplicate[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			duplicate[i][j] = array[i][j]
		}
	}
	return duplicate
}

// ArrayIsEqual compares an array of ints to a given value
// Returned array elements are true if equal to value
func ArrayIsEqual(array []int, value int) []bool {
	var isEqual []bool

	for _, element := range array {
		isEqual = append(isEqual, (element == value))
	}
	return isEqual
}

// ArrayIndices_int returns a slice containing array integers for the specified indices
func ArrayIndices_int(array []int, indices []int) []int {
	var returnSlice []int

	for _, index := range indices {
		returnSlice = append(returnSlice, array[index])
	}
	return returnSlice
}

// ArrayBool_float64 returns a slice containing array rows where the corresponding bool is true
func ArrayBool_float64(array [][]float64, indices []bool) [][]float64 {
	var returnSlice [][]float64

	for index, isTrue := range indices {
		if isTrue {
			returnSlice = append(returnSlice, array[index])
		}
	}
	return returnSlice
}

// MaxColValues returns the maximum value for each column in the array
func MaxColValues(array [][]float64) []float64 {
	rows := len(array)
	cols := len(array[0])
	minVal := math.SmallestNonzeroFloat64

	maxValues := []float64{minVal, minVal, minVal}

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			maxValues[col] = math.Max(array[row][col], maxValues[col])
		}
	}
	return maxValues
}

// MinColValues returns the minimum value for each column in the array
func MinColValues(array [][]float64) []float64 {
	rows := len(array)
	cols := len(array[0])
	maxVal := math.MaxFloat64

	minValues := []float64{maxVal, maxVal, maxVal}

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			minValues[col] = math.Min(array[row][col], minValues[col])
		}
	}
	return minValues
}

// ArraySum returns the total of the elements in a list
func ArraySum(array []float64) float64 {
	total := float64(0)

	for _, element := range array {
		total += element
	}
	return total
}

// ArraySub returns the difference between the elements of two lists
func ArraySub(array1 []float64, array2 []float64) []float64 {
	var returnSlice []float64

	for i := range array1 {
		returnSlice = append(returnSlice, array1[i]-array2[i])
	}
	return returnSlice
}

// ArrayDiv returns the elements of a list divided by the specified value
func ArrayDiv(array []float64, divisor float64) []float64 {
	var returnSlice []float64

	for i := range array {
		returnSlice = append(returnSlice, array[i]/divisor)
	}
	return returnSlice
}
