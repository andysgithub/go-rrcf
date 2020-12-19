package random

import (
	"math"

	"github.com/andysgithub/go-rrcf/num"
)

const (
	a = 134775813
	c = 1
)

var (
	m float64
)

// Random -
type Random struct {
	seed int64
}

// NewRandom -
func NewRandom() Random {
	m = math.Pow(2, 32)

	rnd := Random{}
	return rnd
}

// Seed -
func (rnd *Random) Seed(seed int64) {
	rnd.seed = 12345678
}

// Float -
func (rnd *Random) Float() float64 {
	newRand := int64(math.Mod(float64(a*rnd.seed+c), m))
	rnd.seed = int64(newRand / 1000)
	return float64(newRand) / (m + 1)
}

// // NormFloat -
// func (rnd Random) NormFloat() float64 {
// rnd.Seed = int64(math.Mod(float64(a*rnd.Seed+c), m))
// return float64(rnd.Seed)
// }

// // Randn1 generates a 1D array of normally distributed random floats
// func (rnd Random) Randn1(rows int) []float64 {
// 	newArray := make([]float64, rows)

// 	for i := 0; i < rows; i++ {
// 		newArray[i] = rnd.NormFloat()
// 	}
// 	return newArray
// }

// // Randn2 generates a 2D array of normally distributed random floats
// func (rnd Random) Randn2(rows, cols int) [][]float64 {
// 	newArray := make([][]float64, rows)

// 	for i := 0; i < rows; i++ {
// 		newArray[i] = make([]float64, cols)
// 		for j := 0; j < cols; j++ {
// 			newArray[i][j] = rnd.NormFloat()
// 		}
// 	}
// 	return newArray
// }

// RndChoice generates a random integer from the given range according to a probability density function
func (rnd *Random) RndChoice(rangeVal int, pdf []float64) int {

	// Calculate the Cumulative Distribution Function
	cdf := make([]float64, len(pdf))
	cdf[0] = pdf[0]
	for i := 1; i < len(pdf); i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}

	randVal := rnd.Float()

	// Select the index from the random value according to the cdf
	index := 0
	for randVal > cdf[index] {
		index++
	}

	rangeArray := num.Arange(rangeVal)
	return rangeArray[index]
}

// RndUniform produces a value between min and max from a uniform distribution
func (rnd *Random) RndUniform(min float64, max float64) float64 {
	return min + rnd.Float()*(max-min)
}

// Randn2 generates a 2D array of normally distributed random floats
func (rnd *Random) Randn2(rows, cols int) [][]float64 {
	newArray := make([][]float64, rows)

	for i := 0; i < rows; i++ {
		newArray[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			newArray[i][j] = rnd.Float()
		}
	}
	return newArray
}

// RndArray generates a 2D array of random integers from a collection in a given range without replacement
func (rnd *Random) RndArray(rangeVal int, rows int, cols int) [][]int {
	s := num.Arange(rangeVal)

	// rand.Seed(0)
	rnd.Seed(12345678)

	rndArray := make([][]int, rows)
	for i := 0; i < rows; i++ {
		rndArray[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			total := len(s)
			index := int(rnd.Float() * float64(total))
			rndArray[i][j] = s[index]

			s[index] = s[total-1]
			s = s[:total-1]
		}
	}
	return rndArray
}
