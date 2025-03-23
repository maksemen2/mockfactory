package writer

import (
	"log/slog"

	"github.com/maksemen2/mockfactory/internal/config"
	"github.com/maksemen2/mockfactory/internal/parser"
)

// WriterFactory defines an interface
// for creating Writer instances using parsed struct fields.
type WriterFactory interface {
	Create(structs map[string][]parser.StructField, config *config.Config, logger *slog.Logger) Writer
}

var WriterFactories = map[string]WriterFactory{
	"json": &JsonWriterFactory{},
}

type JsonWriterFactory struct{}

// Create instantiates a new JSON Writer using the provided struct definitions.
func (f *JsonWriterFactory) Create(structs map[string][]parser.StructField, config *config.Config, logger *slog.Logger) Writer {
	return NewJsonWriter(structs, config, logger)
}
