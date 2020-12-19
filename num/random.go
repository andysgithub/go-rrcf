package num

import (
	"math/rand"
)

// Randn1 generates a 1D array of normally distributed random floats
func Randn1(rows int) []float64 {
	newArray := make([]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = rand.NormFloat64()
	}
	return newArray
}

// Randn2 generates a 2D array of normally distributed random floats
func Randn2(rows, cols int) [][]float64 {
	newArray := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			newArray[i][j] = rand.NormFloat64()
		}
	}
	return newArray
}

// Randl2 generates a 2D array of linearly distributed random floats
func Randl2(rows, cols int) [][]float64 {
	newArray := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			newArray[i][j] = rand.Float64()
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

// RndArray generates a 2D array of random integers from a collection in a given range without replacement
func RndArray(rangeVal int, rows int, cols int) [][]int {
	s := Arange(rangeVal)

	// rand.Seed(0)

	rndArray := make([][]int, rows)
	for i := 0; i < rows; i++ {
		rndArray[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			total := len(s)
			index := rand.Intn(total)
			rndArray[i][j] = s[index]

			s[index] = s[total-1]
			s = s[:total-1]
		}
	}
	return rndArray
}

// RndUniform produces a value between min and max from a uniform distribution
func RndUniform(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RndUniformArray produces an array of values between min and max from a uniform distribution
func RndUniformArray(min float64, max float64, rows int, cols int) [][]float64 {
	rndArray := ArrayEmpty(rows, cols)
	for i, row := range rndArray {
		for j := range row {
			rndArray[i][j] = min + rand.Float64()*(max-min)
		}
	}
	return rndArray
}

// RndShuffle -
func RndShuffle(deck []int) []int {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
