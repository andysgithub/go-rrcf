package num

import (
	"math/rand"
)

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

// RndChoice generates a random integer from the given range according to a probability density function
func RndChoice(rangeVal int, pdf []float64) int {

	// Calculate the Cumulative Distribution Function
	cdf := make([]float64, len(pdf))
	cdf[0] = pdf[0]
	for i := 1; i < len(pdf); i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}

	randVal := rand.Float64()

	// Select the index from the random value according to the cdf
	index := 0
	for randVal > cdf[index] {
		index++
	}

	rangeArray := Arange(rangeVal)
	return rangeArray[index]
}

// RndUniform produces a value between min and max from a uniform distribution
func RndUniform(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RndShuffle -
func RndShuffle(deck []int) []int {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
