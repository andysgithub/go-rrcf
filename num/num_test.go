package num

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAround(t *testing.T) {
	values := [][]float64{
		{1.2345, 2.3456, 3.45678},
		{4.54321, 5.6743, 6.78912},
	}

	rounded := Around(values, 3)

	expectedRounded := [][]float64{
		{1.235, 2.346, 3.457},
		{4.543, 5.674, 6.789},
	}
	assert.Equal(t, rounded, expectedRounded, "Decimal rounding incorrect")
}

func TestArange(t *testing.T) {
	duplicated := [][]float64{
		{4, 5, 6, 7},
		{8, 9, 10, 11},
	}

	indexLabels := Arange(len(duplicated[0]))

	expectedRange := []int{0, 1, 2, 3}
	assert.Equal(t, indexLabels, expectedRange, "Label range incorrect")
}

func TestUnique(t *testing.T) {
	duplicated := [][]float64{
		{0, 1, 2, 3},
		{4, 5, 6, 7},
		{0, 1, 2, 3},
		{8, 9, 10, 11},
		{4, 5, 6, 7},
		{0, 1, 2, 3},
		{8, 9, 10, 11},
		{12, 13, 14, 15},
	}

	unique, indicesMap, duplicateTotals := Unique(duplicated)

	assert.Equal(t, len(unique[0]), 4, "There should be 4 unique rows.")
	assert.Equal(t, unique[0], duplicated[0], "Unique row 0 should equal duplicated row 0")
	assert.Equal(t, unique[1], duplicated[1], "Unique row 1 should equal duplicated row 1")
	assert.Equal(t, unique[2], duplicated[3], "Unique row 2 should equal duplicated row 3")
	assert.Equal(t, unique[3], duplicated[7], "Unique row 3 should equal duplicated row 7")

	expectedMap := []int{0, 1, 0, 2, 1, 0, 2, 3}
	assert.Equal(t, indicesMap, expectedMap, "Indices map is incorrect")

	expectedDuplicates := []int{3, 2, 2, 1}
	assert.Equal(t, duplicateTotals, expectedDuplicates, "Incorrect duplicate totals")
}

func TestNoDuplicates(t *testing.T) {
	duplicated := [][]float64{
		{0, 1, 2, 3},
		{4, 5, 6, 7},
		{8, 9, 10, 11},
		{12, 13, 14, 15},
	}

	unique, indicesMap, duplicateTotals := Unique(duplicated)

	assert.Equal(t, len(unique[0]), 4, "There should be 4 unique rows.")
	assert.Equal(t, unique[0], duplicated[0], "Unique row 0 should equal duplicated row 0")
	assert.Equal(t, unique[1], duplicated[1], "Unique row 1 should equal duplicated row 1")
	assert.Equal(t, unique[2], duplicated[2], "Unique row 2 should equal duplicated row 2")
	assert.Equal(t, unique[3], duplicated[3], "Unique row 3 should equal duplicated row 3")

	expectedMap := []int{0, 1, 2, 3}
	assert.Equal(t, indicesMap, expectedMap, "Indices map is incorrect")

	expectedDuplicates := []int{1, 1, 1, 1}
	assert.Equal(t, duplicateTotals, expectedDuplicates, "Duplicate totals should have value 1")
}

func TestNonZero(t *testing.T) {
	array := []bool{false, true, true, false, true}

	result := FlatNonZero(array)
	assert.Equal(t, result, []int{1, 2, 4}, "NonZero array incorrect")
}

func TestArrayEqInt(t *testing.T) {
	array := []int{12, 23, 45, 12, 78, 45, 12}

	result := ArrayEqInt(array, 12)
	assert.Equal(t, result, []bool{true, false, false, true, false, false, true}, "ArrayEq result incorrect")
}

func TestArrayEqFloat(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 12}
	result := ArrayEqFloat(array, 12)
	assert.Equal(t, result, []bool{true, false, false, true, false, false, true}, "ArrayEq result incorrect")
}

func TestArrayCompareBool(t *testing.T) {
	array1 := []bool{false, true, true, false, true}
	array2 := []bool{false, true, false, false, true}

	result := ArrayCompareBool(array1, array1)
	assert.True(t, result, "ArrayCompareBool should return true")

	result = ArrayCompareBool(array1, array2)
	assert.False(t, result, "ArrayCompareBool should return false")
}

func TestArrayCompareFloat(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 12}
	array2 := []float64{12, 23, 54, 12, 87, 45, 12}

	result := ArrayCompareFloat(array1, array1)
	assert.True(t, result, "ArrayCompareFloat should return true")

	result = ArrayCompareFloat(array1, array2)
	assert.False(t, result, "ArrayCompareFloat should return false")
}

func TestArrayLt(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 12}
	array2 := []float64{10, 23, 54, 12, 87, 44, 12}

	result := ArrayLt(array1, array2)
	assert.Equal(t, result, []bool{false, false, true, false, true, false, false}, "ArrayLt result incorrect")
}

func TestArrayGt(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 12}
	array2 := []float64{10, 23, 54, 12, 87, 44, 12}

	result := ArrayGt(array1, array2)
	assert.Equal(t, result, []bool{true, false, false, false, false, true, false}, "ArrayGt result incorrect")
}

func TestArrayLeq(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 12}

	result := ArrayLeq(array, 23.)
	assert.Equal(t, result, []bool{true, true, false, true, false, false, true}, "ArrayLeq result incorrect")
}

func TestArrayAnd(t *testing.T) {
	array1 := []bool{false, true, true, false, true}
	array2 := []bool{false, true, false, false, true}

	result := ArrayAnd(array1, array2)
	assert.Equal(t, result, []bool{false, true, false, false, true}, "ArrayAnd result incorrect")
}

func TestArrayNot(t *testing.T) {
	array := []bool{false, true, true, false, true}

	result := ArrayNot(array)
	assert.Equal(t, result, []bool{true, false, false, true, false}, "ArrayNot result incorrect")
}

func TestArrayDuplicate(t *testing.T) {
	array := [][]float64{{12, 23, 45, 12}, {78, 45, 12, 16}, {23, 45, 12, 82}}

	result := ArrayDuplicate(array)
	assert.Equal(t, result, array, "ArrayDuplicate result incorrect")
}

func TestArrayCopy(t *testing.T) {
	array1 := []float64{91, 92, 93, 94, 95, 96}
	array2 := []float64{12, 78, 45, 12, 16}

	ArrayCopy(array1, array2)
	assert.Equal(t, array1, []float64{12, 78, 45, 12, 16, 96}, "ArrayCopy result incorrect")
}

func TestArrayCopyWhenTrue(t *testing.T) {
	array1 := []float64{91, 92, 93, 94, 95, 96, 97, 98, 99}
	array2 := []float64{12, 23, 45, 12, 78, 45, 12, 16}
	bools := []bool{false, true, true, false, true, false, false, true}

	ArrayCopyWhenTrue(array1, array2, bools)
	assert.Equal(t, array1, []float64{91, 23, 45, 94, 78, 96, 97, 16, 99}, "ArrayCopyWhenTrue result incorrect")
}

func TestArrayIndicesInt(t *testing.T) {
	array := []int{12, 23, 45, 12, 78, 45, 12}
	indices := []int{1, 2, 4, 6}

	result := ArrayIndicesInt(array, indices)
	assert.Equal(t, result, []int{23, 45, 78, 12}, "ArrayIndicesInt result incorrect")
}

func TestArrayBoolFloat(t *testing.T) {
	array := [][]float64{{12, 23, 45}, {78, 45, 12}, {23, 45, 12}, {19, 57, 24}}
	bools := []bool{true, false, true, false}
	result := ArrayBoolFloat(array, bools)
	assert.Equal(t, result, [][]float64{{12, 23, 45}, {23, 45, 12}}, "ArrayBoolFloat result incorrect")
}

func TestMaxColValues(t *testing.T) {
	array := [][]float64{{12, 23, 45}, {78, 45, 12}, {23, 24, 12}, {19, 57, 24}}
	result := MaxColValues(array)
	assert.Equal(t, result, []float64{78, 57, 45}, "MaxColValues result incorrect")
}

func TestMinColValues(t *testing.T) {
	array := [][]float64{{12, 26, 45}, {78, 45, 42}, {23, 24, 37}, {19, 57, 29}}
	result := MinColValues(array)
	assert.Equal(t, result, []float64{12, 24, 29}, "MinColValues result incorrect")
}

func TestArrayMaxValue(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 12, 16}

	result := ArrayMaxValue(array)
	assert.Equal(t, result, 78., "ArrayMaxValue result incorrect")
}

func TestArrayMaximum(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := ArrayMaximum(array1, array2)
	assert.Equal(t, result, []float64{12, 23, 54, 13, 87, 45, 20}, "ArrayMaximum result incorrect")
}

func TestArrayMinimum(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := ArrayMinimum(array1, array2)
	assert.Equal(t, result, []float64{10, 23, 45, 12, 78, 44, 12}, "ArrayMinimum result incorrect")
}

func TestArraySumFloat(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := ArraySumFloat(array)
	assert.Equal(t, result, 235., "ArraySumFloat result incorrect")
}

func TestArraySub(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := ArraySub(array1, array2)
	assert.Equal(t, result, []float64{2, 0, -9, -1, -9, 1, 8}, "ArraySub result incorrect")
}

func TestArraySubVal(t *testing.T) {
	array := []float64{12, 23, 45, 15, 78, 45, 20}

	result := ArraySubVal(array, 15.)
	assert.Equal(t, result, []float64{-3, 8, 30, 0, 63, 30, 5}, "ArraySubVal result incorrect")
}

func TestArrayMul(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := ArrayMul(array1, array2)
	assert.Equal(t, result, []float64{120, 529, 2430, 156, 6786, 1980, 240}, "ArrayMul result incorrect")
}

func TestArrayMulVal(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := ArrayMulVal(array, 25.)
	assert.Equal(t, result, []float64{300, 575, 1125, 300, 1950, 1125, 500}, "ArrayMulVal result incorrect")
}

func TestArray2DMulVal(t *testing.T) {
	array := [][]float64{{12, 23, 45}, {78, 45, 12}, {23, 24, 12}, {19, 57, 24}}
	result := Array2DMulVal(array, 25)
	assert.Equal(t, result, [][]float64{{300, 575, 1125}, {1950, 1125, 300}, {575, 600, 300}, {475, 1425, 600}}, "Array2DMulVal result incorrect")
}

func TestArrayMulValInt(t *testing.T) {
	array := []int{12, 23, 45, 12, 78, 45, 20}

	result := ArrayMulValInt(array, 25.)
	assert.Equal(t, result, []float64{300, 575, 1125, 300, 1950, 1125, 500}, "ArrayMulValInt result incorrect")
}

func TestArray2DAddVal(t *testing.T) {
	array1 := [][]float64{{12, 26, 45}, {78, 45, 42}, {23, 24, 37}, {19, 57, 29}}
	array2 := [][]float64{{78, 45, 42}, {23, 24, 37}, {19, 57, 29}, {12, 26, 45}}

	result := Array2DAdd(array1, array2)
	assert.Equal(t, result, [][]float64{{90, 71, 87}, {101, 69, 79}, {42, 81, 66}, {31, 83, 74}}, "Array2DAdd result incorrect")
}

func TestArrayAddVal(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := ArrayAddVal(array, 25.)
	assert.Equal(t, result, []float64{37, 48, 70, 37, 103, 70, 45}, "ArrayAddVal result incorrect")
}

func TestArrayDivVal(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := ArrayDivVal(array, 10.)
	assert.Equal(t, result, []float64{1.2, 2.3, 4.5, 1.2, 7.8, 4.5, 2.0}, "ArrayDivVal result incorrect")
}
