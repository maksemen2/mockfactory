package generator

import (
	"strings"
	"testing"

	"github.com/maksemen2/mockfactory/internal/testutils"
)

func TestStringGenerator_Length(t *testing.T) {
	tags := map[string]string{
		"len":    "10",
		"prefix": "test_",
		"suffix": "_end",
	}
	g := NewStringGenerator(tags, testRand(), testutils.TestLogger()).(*StringGenerator)

	val, _ := g.Evaluate()
	if len(val) != 10+len("test_")+len("_end") {
		t.Errorf("Unexpected length: %d", len(val))
	}
	if !strings.HasPrefix(val, "test_") || !strings.HasSuffix(val, "_end") {
		t.Errorf("Invalid prefix/suffix")
	}
}

func TestStringGenerator_DefaultLength(t *testing.T) {
	g := NewStringGenerator(map[string]string{}, testRand(), testutils.TestLogger()).(*StringGenerator)
	val, _ := g.Evaluate()
	if len(val) != defaultLength {
		t.Errorf("Expected default length %d, got %d", defaultLength, len(val))
	}
}
