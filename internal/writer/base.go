package writer

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/maksemen2/mockfactory/internal/config"
	"github.com/maksemen2/mockfactory/internal/parser"
)

// BaseWriter implements base writer functionality
type BaseWriter struct {
	structs map[string][]parser.StructField
	config  *config.Config
	logger  *slog.Logger
}

// GetFileName returns the file name for the struct
// depending on is there a template in the config
// and replaces the placeholders with the actual values
func (w *BaseWriter) GetFileName(structName string) string {
	if w.config.Output.FileNameTemplate != "" {
		result := w.config.Output.FileNameTemplate
		result = strings.ReplaceAll(result, "{struct}", structName)
		result = strings.ReplaceAll(result, "{count}", strconv.Itoa(w.config.Generation.Count))
		return result
	}
	return structName
}
