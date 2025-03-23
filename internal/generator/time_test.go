package generator

import (
	"testing"
	"time"

	"github.com/maksemen2/mockfactory/internal/testutils"
)

func TestTimeGenerator_Ranges(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name   string
		tags   map[string]string
		assert func(time.Time) bool
	}{
		{
			name:   "past",
			tags:   map[string]string{"range": "past"},
			assert: func(tm time.Time) bool { return tm.Before(now) },
		},
		{
			name:   "future",
			tags:   map[string]string{"range": "future"},
			assert: func(tm time.Time) bool { return tm.After(now) },
		},
		{
			name:   "default",
			tags:   map[string]string{},
			assert: func(tm time.Time) bool { return tm.Sub(now) < time.Second },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewTimeGenerator(tt.tags, testRand(), testutils.TestLogger()).(*TimeGenerator)
			val, _ := g.Evaluate()
			if !tt.assert(val) {
				t.Errorf("Time validation failed for %s", tt.name)
			}
		})
	}
}
