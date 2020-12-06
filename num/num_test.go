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

	unique, removedIndices, duplicateTotals := Unique(duplicated)

	assert.Equal(t, len(unique[0]), 4, "There should be 4 unique rows.")
	assert.Equal(t, unique[0], duplicated[0], "Unique row 0 should equal duplicated row 0")
	assert.Equal(t, unique[1], duplicated[1], "Unique row 1 should equal duplicated row 1")
	assert.Equal(t, unique[2], duplicated[3], "Unique row 2 should equal duplicated row 3")
	assert.Equal(t, unique[3], duplicated[7], "Unique row 3 should equal duplicated row 7")

	expectedRemoved := []int{2, 4, 5, 6}
	assert.Equal(t, removedIndices, expectedRemoved, "Incorrect removed indices")

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

	unique, removedIndices, duplicateTotals := Unique(duplicated)

	//fmt.Printf("%v %v %v\n", unique, removedIndices, duplicateTotals)

	assert.Equal(t, len(unique[0]), 4, "There should be 4 unique rows.")
	assert.Equal(t, unique[0], duplicated[0], "Unique row 0 should equal duplicated row 0")
	assert.Equal(t, unique[1], duplicated[1], "Unique row 1 should equal duplicated row 1")
	assert.Equal(t, unique[2], duplicated[2], "Unique row 2 should equal duplicated row 3")
	assert.Equal(t, unique[3], duplicated[3], "Unique row 3 should equal duplicated row 7")

	var expectedRemoved []int
	assert.Equal(t, removedIndices, expectedRemoved, "Removed indices should be empty")

	expectedDuplicates := []int{1, 1, 1, 1}
	assert.Equal(t, duplicateTotals, expectedDuplicates, "Duplicate totals should have value 1")
}
