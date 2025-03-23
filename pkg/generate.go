package pkg

import (
	"log/slog"

	"github.com/maksemen2/mockfactory/internal/config"
	"github.com/maksemen2/mockfactory/internal/parser"
	"github.com/maksemen2/mockfactory/internal/writer"
)

func GenerateFromFile(cfg *config.Config, logger *slog.Logger) error {

	p := parser.NewParser(cfg, logger)
	fields, err := p.ParseFile()
	if err != nil {
		panic(err)
	}

	w := writer.NewJsonWriter(fields, cfg, logger)
	return w.Write()
}
