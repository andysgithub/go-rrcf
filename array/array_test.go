package array

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

func TestEqualInt(t *testing.T) {
	array := []int{12, 23, 45, 12, 78, 45, 12}

	result := EqualInt(array, 12)
	assert.Equal(t, result, []bool{true, false, false, true, false, false, true}, "ArrayEq result incorrect")
}

func TestEqualFloat(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 12}
	result := EqualFloat(array, 12)
	assert.Equal(t, result, []bool{true, false, false, true, false, false, true}, "ArrayEq result incorrect")
}

func TestCompareBool(t *testing.T) {
	array1 := []bool{false, true, true, false, true}
	array2 := []bool{false, true, false, false, true}

	result := CompareBool(array1, array1)
	assert.True(t, result, "CompareBool should return true")

	result = CompareBool(array1, array2)
	assert.False(t, result, "CompareBool should return false")
}

func TestCompareFloat(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 12}
	array2 := []float64{12, 23, 54, 12, 87, 45, 12}

	result := CompareFloat(array1, array1)
	assert.True(t, result, "CompareFloat should return true")

	result = CompareFloat(array1, array2)
	assert.False(t, result, "CompareFloat should return false")
}

func TestLtFloat(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 12}
	array2 := []float64{10, 23, 54, 12, 87, 44, 12}

	result := LtFloat(array1, array2)
	assert.Equal(t, result, []bool{false, false, true, false, true, false, false}, "LtFloat result incorrect")
}

func TestGtFloat(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 12}
	array2 := []float64{10, 23, 54, 12, 87, 44, 12}

	result := GtFloat(array1, array2)
	assert.Equal(t, result, []bool{true, false, false, false, false, true, false}, "GtFloat result incorrect")
}

func TestLeqFloat(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 12}

	result := LeqFloat(array, 23.)
	assert.Equal(t, result, []bool{true, true, false, true, false, false, true}, "LeqFloat result incorrect")
}

func TestAndBool(t *testing.T) {
	array1 := []bool{false, true, true, false, true}
	array2 := []bool{false, true, false, false, true}

	result := AndBool(array1, array2)
	assert.Equal(t, result, []bool{false, true, false, false, true}, "AndBool result incorrect")
}

func TestNotBool(t *testing.T) {
	array := []bool{false, true, true, false, true}

	result := NotBool(array)
	assert.Equal(t, result, []bool{true, false, false, true, false}, "NotBool result incorrect")
}

func TestDuplicateFloat(t *testing.T) {
	array := [][]float64{{12, 23, 45, 12}, {78, 45, 12, 16}, {23, 45, 12, 82}}

	result := DuplicateFloat(array)
	assert.Equal(t, result, array, "DuplicateFloat result incorrect")
}

func TestCopyFloat(t *testing.T) {
	array1 := []float64{91, 92, 93, 94, 95, 96}
	array2 := []float64{12, 78, 45, 12, 16}

	CopyFloat(array1, array2)
	assert.Equal(t, array1, []float64{12, 78, 45, 12, 16, 96}, "CopyFloat result incorrect")
}

func TestCopyFloatWhenTrue(t *testing.T) {
	array1 := []float64{91, 92, 93, 94, 95, 96, 97, 98, 99}
	array2 := []float64{12, 23, 45, 12, 78, 45, 12, 16}
	bools := []bool{false, true, true, false, true, false, false, true}

	CopyFloatWhenTrue(array1, array2, bools)
	assert.Equal(t, array1, []float64{91, 23, 45, 94, 78, 96, 97, 16, 99}, "CopyFloatWhenTrue result incorrect")
}

func TestIndicesInt(t *testing.T) {
	array := []int{12, 23, 45, 12, 78, 45, 12}
	indices := []int{1, 2, 4, 6}

	result := IndicesInt(array, indices)
	assert.Equal(t, result, []int{23, 45, 78, 12}, "IndicesInt result incorrect")
}

func TestWhereTrueFloat(t *testing.T) {
	array := [][]float64{{12, 23, 45}, {78, 45, 12}, {23, 45, 12}, {19, 57, 24}}
	bools := []bool{true, false, true, false}
	result := WhereTrueFloat(array, bools)
	assert.Equal(t, result, [][]float64{{12, 23, 45}, {23, 45, 12}}, "WhereTrueFloat result incorrect")
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

func TestMaxValue(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 12, 16}

	result := MaxValue(array)
	assert.Equal(t, result, 78., "MaxValue result incorrect")
}

func TestMaximum(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := Maximum(array1, array2)
	assert.Equal(t, result, []float64{12, 23, 54, 13, 87, 45, 20}, "Maximum result incorrect")
}

func TestMinimum(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := Minimum(array1, array2)
	assert.Equal(t, result, []float64{10, 23, 45, 12, 78, 44, 12}, "Minimum result incorrect")
}

func TestSumFloat(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := SumFloat(array)
	assert.Equal(t, result, 235., "SumFloat result incorrect")
}

func TestSubtract(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := Subtract1D(array1, array2)
	assert.Equal(t, result, []float64{2, 0, -9, -1, -9, 1, 8}, "Subtract1D result incorrect")
}

func TestSubtractVal1D(t *testing.T) {
	array := []float64{12, 23, 45, 15, 78, 45, 20}

	result := SubtractVal1D(array, 15.)
	assert.Equal(t, result, []float64{-3, 8, 30, 0, 63, 30, 5}, "SubtractVal1D result incorrect")
}

func TestMultiply1D(t *testing.T) {
	array1 := []float64{12, 23, 45, 12, 78, 45, 20}
	array2 := []float64{10, 23, 54, 13, 87, 44, 12}

	result := Multiply1D(array1, array2)
	assert.Equal(t, result, []float64{120, 529, 2430, 156, 6786, 1980, 240}, "Multiply1D result incorrect")
}

func TestMultiplyVal1D(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := MultiplyVal1D(array, 25.)
	assert.Equal(t, result, []float64{300, 575, 1125, 300, 1950, 1125, 500}, "MultiplyVal1D result incorrect")
}

func TestMultiplyVal2D(t *testing.T) {
	array := [][]float64{{12, 23, 45}, {78, 45, 12}, {23, 24, 12}, {19, 57, 24}}
	result := MultiplyVal2D(array, 25)
	assert.Equal(t, result, [][]float64{{300, 575, 1125}, {1950, 1125, 300}, {575, 600, 300}, {475, 1425, 600}}, "MultiplyVal2D result incorrect")
}

func TestMultiplyVal1DInt(t *testing.T) {
	array := []int{12, 23, 45, 12, 78, 45, 20}

	result := MultiplyVal1DInt(array, 25.)
	assert.Equal(t, result, []float64{300, 575, 1125, 300, 1950, 1125, 500}, "MultiplyVal1DInt result incorrect")
}

func TestAddFloat2DVal(t *testing.T) {
	array1 := [][]float64{{12, 26, 45}, {78, 45, 42}, {23, 24, 37}, {19, 57, 29}}
	array2 := [][]float64{{78, 45, 42}, {23, 24, 37}, {19, 57, 29}, {12, 26, 45}}

	result := AddFloat2D(array1, array2)
	assert.Equal(t, result, [][]float64{{90, 71, 87}, {101, 69, 79}, {42, 81, 66}, {31, 83, 74}}, "AddFloat2D result incorrect")
}

func TestAddVal1D(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := AddVal1D(array, 25.)
	assert.Equal(t, result, []float64{37, 48, 70, 37, 103, 70, 45}, "AddVal1D result incorrect")
}

func TestDivVal1D(t *testing.T) {
	array := []float64{12, 23, 45, 12, 78, 45, 20}

	result := DivVal1D(array, 10.)
	assert.Equal(t, result, []float64{1.2, 2.3, 4.5, 1.2, 7.8, 4.5, 2.0}, "DivVal1D result incorrect")
}
