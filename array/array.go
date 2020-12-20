package array

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
	indices := make(map[string]int)
	indexU := 0

	for _, values := range X {
		rowKey := SliceToString(values, ",")
		// If this row of values has not been recorded yet
		if _, value := keys[rowKey]; !value {
			// Set the total occurrences for this row to 1
			keys[rowKey] = 1
			rowKeys = append(rowKeys, rowKey)
			U = append(U, values)
			// Create a new index in U of this row
			indices[rowKey] = indexU
			indexU++
		} else {
			// Increment the total occurrences for this row
			keys[rowKey]++
		}
		// Record the index in U for this row value
		I = append(I, indices[rowKey])
	}
	// Compile the list of unique row totals
	for _, rowKey := range rowKeys {
		N = append(N, keys[rowKey])
	}
	return U, I, N
}

// OnesBool returns a list of bools set to true
func OnesBool(rows int) []bool {
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

// EqualInt compares an array of ints to a given value
// Returned array elements are true if equal to value
func EqualInt(array []int, value int) []bool {
	var isEqual []bool

	for _, element := range array {
		isEqual = append(isEqual, (element == value))
	}
	return isEqual
}

// EqualFloat compares an array of floats to a given value
// Returned array elements are true if equal to value
func EqualFloat(array []float64, value float64) []bool {
	var isEqual []bool

	for _, element := range array {
		isEqual = append(isEqual, (element == value))
	}
	return isEqual
}

// CompareBool compares two boolean arrays
// Returns true if all array elements are equal
func CompareBool(array1 []bool, array2 []bool) bool {
	for i := range array1 {
		if array1[i] != array2[i] {
			return false
		}
	}
	return true
}

// CompareFloat compares two float arrays
// Returns true if all array elements are equal
func CompareFloat(array1 []float64, array2 []float64) bool {
	for i := range array1 {
		if array1[i] != array2[i] {
			return false
		}
	}
	return true
}

// LtFloat compares two arrays of floats
// Returned array elements are true if array1 less than array2
func LtFloat(array1 []float64, array2 []float64) []bool {
	var isLt []bool

	for i := range array1 {
		isLt = append(isLt, (array1[i] < array2[i]))
	}
	return isLt
}

// GtFloat compares two arrays of floats
// Returned array elements are true if array1 greater than array2
func GtFloat(array1 []float64, array2 []float64) []bool {
	var isGt []bool

	for i := range array1 {
		isGt = append(isGt, (array1[i] > array2[i]))
	}
	return isGt
}

// LeqFloat compares an array of floats to a given value
// Returned array elements are true if less than or equal to value
func LeqFloat(array []float64, value float64) []bool {
	var isLeq []bool

	for _, element := range array {
		isLeq = append(isLeq, (element <= value))
	}
	return isLeq
}

// AndBool returns the logical and of two boolean arrays
func AndBool(array1 []bool, array2 []bool) []bool {
	var andArray []bool

	for i := range array1 {
		andArray = append(andArray, array1[i] && array2[i])
	}
	return andArray
}

// NotBool returns the inverse of a boolean array
func NotBool(array []bool) []bool {
	var inverseArray []bool

	for _, element := range array {
		inverseArray = append(inverseArray, !element)
	}
	return inverseArray
}

// DuplicateFloat copies a 2D array of floats by value and returns the copy
func DuplicateFloat(array [][]float64) [][]float64 {
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

// CopyFloat copies array2 elements into array1
func CopyFloat(array1 []float64, array2 []float64) {
	for i, element := range array2 {
		array1[i] = element
	}
}

// CopyFloatWhenTrue copies array2 elements into array1 where the bools array element is true
func CopyFloatWhenTrue(array1 []float64, array2 []float64, bools []bool) {
	for i, element := range array2 {
		if bools[i] {
			array1[i] = element
		}
	}
}

// IndicesInt returns a slice containing array integers for the specified indices
func IndicesInt(array []int, indices []int) []int {
	var returnSlice []int

	for _, index := range indices {
		returnSlice = append(returnSlice, array[index])
	}
	return returnSlice
}

// WhereTrueFloat returns a slice containing array rows where the corresponding bool is true
func WhereTrueFloat(array [][]float64, indices []bool) [][]float64 {
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
	minVal := -math.MaxFloat64

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

// MaxValue returns the maximum value in a list
func MaxValue(array []float64) float64 {
	maxVal := -math.MaxFloat64

	for _, value := range array {
		maxVal = math.Max(value, maxVal)
	}
	return maxVal
}

// Maximum compares two arrays and returns a new array containing the element-wise maxima
func Maximum(array1 []float64, array2 []float64) []float64 {
	var maxValues []float64
	cols := len(array1)

	for col := 0; col < cols; col++ {
		maxValues = append(maxValues, math.Max(array1[col], array2[col]))
	}
	return maxValues
}

// Minimum compares two arrays and returns a new array containing the element-wise minima
func Minimum(array1 []float64, array2 []float64) []float64 {
	var minValues []float64
	cols := len(array1)

	for col := 0; col < cols; col++ {
		minValues = append(minValues, math.Min(array1[col], array2[col]))
	}
	return minValues
}

// SumFloat returns the total of the elements in a list
func SumFloat(array []float64) float64 {
	total := float64(0)

	for _, element := range array {
		total += element
	}
	return total
}

// Subtract1D returns the difference between the elements of two lists
func Subtract1D(array1 []float64, array2 []float64) []float64 {
	var returnSlice []float64

	for i := range array1 {
		returnSlice = append(returnSlice, array1[i]-array2[i])
	}
	return returnSlice
}

// SubtractVal1D subtracts the specified value from the elements of a list
func SubtractVal1D(array []float64, value float64) []float64 {
	var returnSlice []float64

	for _, element := range array {
		returnSlice = append(returnSlice, element-value)
	}
	return returnSlice
}

// Multiply1D returns the product of the elements of two lists
func Multiply1D(array1 []float64, array2 []float64) []float64 {
	var returnSlice []float64

	for i := range array1 {
		returnSlice = append(returnSlice, array1[i]*array2[i])
	}
	return returnSlice
}

// MultiplyVal1D returns the elements of a list multiplied by the specified value
func MultiplyVal1D(array []float64, value float64) []float64 {
	var returnSlice []float64

	for _, element := range array {
		returnSlice = append(returnSlice, element*value)
	}
	return returnSlice
}

// MultiplyVal2D returns the elements of a 2D array multiplied by the specified value
func MultiplyVal2D(array [][]float64, value float64) [][]float64 {
	rows := len(array)
	cols := len(array[0])

	result := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = array[i][j] * value
		}
	}
	return result
}

// MultiplyVal1DInt returns the elements of a integer list multiplied by the specified float value
func MultiplyVal1DInt(array []int, value float64) []float64 {
	var returnSlice []float64

	for _, element := range array {
		returnSlice = append(returnSlice, float64(element)*value)
	}
	return returnSlice
}

// AddFloat2D returns the sum of two 2D arrays
func AddFloat2D(array1 [][]float64, array2 [][]float64) [][]float64 {
	rows := len(array1)
	cols := len(array1[0])

	result := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			result[i][j] = array1[i][j] + array2[i][j]
		}
	}
	return result
}

// AddVal1D returns the elements of a list multiplied by the specified value
func AddVal1D(array []float64, value float64) []float64 {
	var returnSlice []float64

	for _, element := range array {
		returnSlice = append(returnSlice, element+value)
	}
	return returnSlice
}

// DivVal1D returns the elements of a list divided by the specified value
func DivVal1D(array []float64, value float64) []float64 {
	var returnSlice []float64

	for _, element := range array {
		returnSlice = append(returnSlice, element/value)
	}
	return returnSlice
}

// GetColumn returns the column with the specified index from an array
func GetColumn(array [][]float64, columnIndex int) []float64 {
	column := make([]float64, 0)
	for _, row := range array {
		column = append(column, row[columnIndex])
	}
	return column
}

// FillElements fills the elements of a list with the specified value
func FillElements(array []float64, rowStart int, rowEnd int, value float64) {
	for i := rowStart; i <= rowEnd; i++ {
		array[i] = value
	}
}

// FillRows fills array rows with the specified value
func FillRows(array [][]float64, rowStart int, rowEnd int, value float64) {
	for i := rowStart; i <= rowEnd; i++ {
		for j := 0; j < len(array[0]); j++ {
			array[i][j] = value
		}
	}
}

// FillColumn fills array column with the specified value
func FillColumn(array [][]float64, col int, rowStart int, rowEnd int, value float64) {
	for i := rowStart; i <= rowEnd; i++ {
		array[i][col] = value
	}
}

// SumTrue totals the true values in the list
func SumTrue(array []bool) int {
	sum := 0
	for _, isTrue := range array {
		if isTrue {
			sum++
		}
	}
	return sum
}

// VStack stacks the sequence of input arrays vertically to make a single array
func VStack(array1 []float64, array2 []float64) [][]float64 {
	var stacked [][]float64

	stacked = append(stacked, array1)
	stacked = append(stacked, array2)
	return stacked
}

// Full returns a 1D array of given length, filled with fillValue
func Full(length int, fillValue float64) []float64 {
	var filledArray []float64

	for i := 0; i < length; i++ {
		filledArray = append(filledArray, fillValue)
	}
	return filledArray
}

// Zero2D returns a zeroed 2D array
func Zero2D(rows int, cols int) [][]float64 {
	array := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		array[i] = make([]float64, cols)
	}
	return array
}

// IsClose returns a boolean array where two arrays are element-wise equal within a tolerance
func IsClose(array1 []float64, array2 []float64, tolerance float64) []bool {
	var result []bool
	for i, value := range array1 {
		result = append(result, math.Abs(array2[i]-value) <= tolerance)
	}
	return result
}

// AllClose returns true if two arrays are element-wise equal within a tolerance
func AllClose(array1 [][]float64, array2 [][]float64, tolerance float64) bool {
	for i, rows := range array1 {
		for j, value := range rows {
			if math.Abs(array2[i][j]-value) > tolerance {
				return false
			}
		}
	}
	return true
}

// AnyEqFloat returns true if any item in array1 equals corresponding item in array2
func AnyEqFloat(array1 []float64, array2 []float64) bool {
	for i, value := range array1 {
		if array2[i] == value {
			return true
		}
	}
	return false
}

// ContainsBool returns true if the boolean array contains the boolean value
func ContainsBool(array []bool, value bool) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}
	return false
}

// AnyTrueBool returns true if any item in array is true
func AnyTrueBool(array []bool) bool {
	for _, value := range array {
		if value {
			return true
		}
	}
	return false
}

// AllTrueBool returns true if all items in array are true
func AllTrueBool(array []bool) bool {
	for _, value := range array {
		if !value {
			return false
		}
	}
	return true
}

// CumSum returns the cumulative sum of the elements in a list
func CumSum(array []float64) []float64 {
	accumulator := float64(0)
	var result []float64

	for _, value := range array {
		accumulator += value
		result = append(result, accumulator)
	}
	return result
}

// Sin returns the sine of the elements in a list
func Sin(array []float64) []float64 {
	var result []float64

	for _, value := range array {
		result = append(result, math.Sin(value))
	}
	return result
}

// ReshapeRow converts a list into a 2D array with one row
func ReshapeRow(array []float64) [][]float64 {
	result := Zero2D(1, len(array))

	for i, value := range array {
		result[0][i] = value
	}
	return result
}

// ReshapeCol converts a list into a 2D array with one column
func ReshapeCol(array []float64) [][]float64 {
	result := Zero2D(len(array), 1)

	for i, value := range array {
		result[i][0] = value
	}
	return result
}

// Sample produces a 2D array consisting of sampled rows from the original array
func Sample(array [][]float64, samples []int) [][]float64 {
	rows := len(samples)
	cols := len(array[0])

	sampled := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		sampled[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			row := samples[i]
			sampled[i][j] = array[row][j]
		}
	}
	return sampled
}

// SliceToString returns a string of array values separated by the specified character
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
