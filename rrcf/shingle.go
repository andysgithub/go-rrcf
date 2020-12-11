package rrcf

// Shingle generates shingles (a rolling window) of a given size
type Shingle struct {
	Sequence [][]float64
	Size     int
	RowStart *int
}

// NewShingle returns an initialised Shingle object
func NewShingle(sequence [][]float64, size int) *Shingle {
	rowStart := 0
	shingle := Shingle{
		Sequence: sequence,
		Size:     size,
		RowStart: &rowStart,
	}
	return &shingle
}

// Next returns the next collection of a given size from a sequence
func (shingle Shingle) Next() [][]float64 {
	// Read size number of rows from sequence array starting at row start
	first := *shingle.RowStart
	last := first + shingle.Size
	// Increment row start by one
	*shingle.RowStart++
	// Return the array read
	return shingle.Sequence[first:last]
}
