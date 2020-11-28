package go-rrcf

import (
)

// RRCF - Robust Random Cut Forest
type RRCF struct {
}

// NewRRCF - Returns a new forest
func NewRRCF() RRCF {
  rrcf := RRCF{}
  return rrcf
}

// Init - Initialises the random cut forest
func (rrcf RRCF) Init(X int, index_labels int, precision int, 
                 random_state int) {
}
