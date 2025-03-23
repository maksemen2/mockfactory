package generator

import (
	"math/rand"
)

func testRand() *rand.Rand {
	return rand.New(rand.NewSource(1))
}
