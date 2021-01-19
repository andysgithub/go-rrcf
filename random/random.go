package random

import (
	"math/rand"

	"github.com/andysgithub/go-rrcf/array"
)

// RandomState holds a reference to the random number generator for a specific instance
type RandomState struct {
	rnd *rand.Rand
}

// NewRandomState sets the seed value for the random number generator
func NewRandomState(seed int64) *RandomState {
	randomState := RandomState{
		rand.New(rand.NewSource(seed)),
	}

	return &randomState
}

// Normal1D generates a 1D array of normally distributed random floats
func (rng *RandomState) Normal1D(rows int) []float64 {
	newArray := make([]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = rng.rnd.NormFloat64()
	}
	return newArray
}

// Normal2D generates a 2D array of normally distributed random floats
func (rng *RandomState) Normal2D(rows, cols int) [][]float64 {
	newArray := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			newArray[i][j] = rng.rnd.NormFloat64()
		}
	}
	return newArray
}

// Linear2D generates a 2D array of linearly distributed random floats
func (rng *RandomState) Linear2D(rows, cols int) [][]float64 {
	newArray := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			newArray[i][j] = rng.rnd.Float64()
		}
	}
	return newArray
}

// Choice generates a random integer from the given range according to a probability density function
func (rng *RandomState) Choice(rangeVal int, pdf []float64) int {
	// Calculate the Cumulative Distribution Function
	cdf := make([]float64, len(pdf))
	cdf[0] = pdf[0]
	for i := 1; i < len(pdf); i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}

	randVal := rng.rnd.Float64()

	// Select the index from the random value according to the cdf
	index := 0
	for randVal > cdf[index] {
		index++
	}

	rangeArray := array.Arange(rangeVal)
	return rangeArray[index]
}

// Array generates a 2D array of random integers from a collection in a given range without replacement
func (rng *RandomState) Array(rangeVal int, rows int, cols int) [][]int {
	s := array.Arange(rangeVal)

	rndArray := make([][]int, rows)
	for i := 0; i < rows; i++ {
		rndArray[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			total := len(s)
			index := rng.rnd.Intn(total)
			rndArray[i][j] = s[index]

			s[index] = s[total-1]
			s = s[:total-1]
		}
	}
	return rndArray
}

// Uniform produces a value between min and max from a uniform distribution
func (rng *RandomState) Uniform(min float64, max float64) float64 {
	return min + rng.rnd.Float64()*(max-min)
}

// UniformArray produces an array of values between min and max from a uniform distribution
func (rng *RandomState) UniformArray(min float64, max float64, rows int, cols int) [][]float64 {
	rndArray := array.Zero2D(rows, cols)
	for i, row := range rndArray {
		for j := range row {
			rndArray[i][j] = min + rng.rnd.Float64()*(max-min)
		}
	}
	return rndArray
}

// Shuffle randomises the ordering of rows in a slice of ints
func (rng *RandomState) Shuffle(deck []int) []int {
	rng.rnd.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
