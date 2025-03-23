package generator

import (
	"math"
	"testing"

	"github.com/maksemen2/mockfactory/internal/testutils"
)

func TestFloatGenerator_Range(t *testing.T) {
	tags := map[string]string{
		"min": "5.5",
		"max": "10.5",
	}
	g := NewFloatGenerator[float64](tags, testRand(), testutils.TestLogger()).(*FloatGenerator[float64])

	for i := 0; i < 100; i++ {
		val, _ := g.Evaluate()
		if val < 5.5 || val > 10.5 {
			t.Fatalf("Value %f out of range [5.5,10.5]", val)
		}
	}
}

func TestFloatGenerator_Defaults(t *testing.T) {
	g := NewFloatGenerator[float32](map[string]string{}, testRand(), testutils.TestLogger()).(*FloatGenerator[float32])

	if g.min != -math.MaxFloat32 || g.max != math.MaxFloat32 {
		t.Errorf("Unexpected default range for float32")
	}
}
