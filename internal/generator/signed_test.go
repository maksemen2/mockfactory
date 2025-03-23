package generator

import (
	"math"
	"testing"

	"github.com/maksemen2/mockfactory/internal/testutils"
)

func TestNewSignedGenerator_Int64FullRange(t *testing.T) {
	tags := map[string]string{
		"min": "-9223372036854775808",
		"max": "9223372036854775807",
	}
	g := NewSignedGenerator[int64](tags, testRand(), testutils.TestLogger()).(*SignedGenerator[int64])

	if g.min != math.MinInt64 || g.max != math.MaxInt64 {
		t.Fatalf("Expected full int64 range, got min=%d, max=%d", g.min, g.max)
	}

	val, _ := g.Evaluate()
	if val < math.MinInt64 || val > math.MaxInt64 {
		t.Errorf("Value %d out of range", val)
	}
}

func TestSignedGenerator_Clamping(t *testing.T) {
	tags := map[string]string{
		"min": "-200",
		"max": "500",
	}
	g := NewSignedGenerator[int8](tags, testRand(), testutils.TestLogger()).(*SignedGenerator[int8])

	if g.min != -128 || g.max != 127 {
		t.Errorf("Expected clamping to int8 bounds, got min=%d, max=%d", g.min, g.max)
	}
}

func TestSignedGenerator_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		typeGen  func() *SignedGenerator[int]
		expected func(int) bool
	}{
		{
			name: "default range",
			typeGen: func() *SignedGenerator[int] {
				return NewSignedGenerator[int](map[string]string{}, testRand(), testutils.TestLogger()).(*SignedGenerator[int])
			},
			expected: func(v int) bool { return v >= math.MinInt && v <= math.MaxInt },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.typeGen()
			val, _ := g.Evaluate()
			if !tt.expected(val) {
				t.Errorf("Unexpected value: %d", val)
			}
		})
	}
}
