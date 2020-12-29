package rrcf

///// 2D VERSION /////

// Shingle generates shingles (a rolling window) of a given size from a 2D array
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

///// 1D VERSION /////

// ShingleList generates shingles (a rolling window) of a given size from a list
type ShingleList struct {
	Sequence     []float64
	Size         int
	RowStart     *int
	TotalSamples int
}

// NewShingleList returns an initialised Shingle object
func NewShingleList(sequence []float64, size int) *ShingleList {
	rowStart := 0
	shingle := ShingleList{
		Sequence:     sequence,
		Size:         size,
		RowStart:     &rowStart,
		TotalSamples: len(sequence) - size,
	}
	return &shingle
}

// NextInList returns the next collection of a given size from a sequence
func (shingle ShingleList) NextInList() []float64 {
	// Read size number of rows from sequence array starting at row start
	first := *shingle.RowStart
	last := first + shingle.Size
	// Increment row start by one
	*shingle.RowStart++
	// Return the array read
	return shingle.Sequence[first:last]
}
