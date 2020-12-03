package rrcf

import (
	"math/rand"
	"testing"

	"github.com/andysgithub/go-rrcf/num"
)

func TestRRCF(t *testing.T) {
	rand.Seed(0)
	n := 100
	d := 3
	X := num.Randn(n, d)
	Z := num.ArrayCopy(X)
	Z[10][0] = 1

	tree := RCTree()
	tree.Init(X, nil, 9, 0)
	duplicateTree := RCTree()
	duplicateTree.Init(Z, nil, 9, 0)
}
