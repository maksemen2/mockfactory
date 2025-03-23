package generator

import (
	"math"
	"testing"

	"github.com/maksemen2/mockfactory/internal/testutils"
)

func TestUnsignedGenerator_Evaluate(t *testing.T) {
	tags := map[string]string{
		"min": "10",
		"max": "20",
	}
	g := NewUnsignedGenerator[uint](tags, testRand(), testutils.TestLogger()).(*UnsignedGenerator[uint])

	for i := 0; i < 100; i++ {
		val, _ := g.Evaluate()
		if val < 10 || val > 20 {
			t.Fatalf("Generated value %d outside range [10,20]", val)
		}
	}
}

func TestUnsignedGenerator_MaxBounds(t *testing.T) {
	g := NewUnsignedGenerator[uint8](map[string]string{}, testRand(), testutils.TestLogger()).(*UnsignedGenerator[uint8])

	if g.max != math.MaxUint8 {
		t.Errorf("Expected max %d, got %d", math.MaxUint8, g.max)
	}
}
