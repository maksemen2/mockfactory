package writer

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/maksemen2/mockfactory/internal/config"
	"github.com/maksemen2/mockfactory/internal/parser"
)

// JsonWriter writes parsed structs in a json format
type JsonWriter struct {
	BaseWriter
}

func NewJsonWriter(structs map[string][]parser.StructField, config *config.Config, logger *slog.Logger) Writer {
	return &JsonWriter{BaseWriter: BaseWriter{
		structs: structs,
		config:  config,
		logger:  logger,
	}}
}

// Write writes the parsed structs to a JSON file.
func (w *JsonWriter) Write() error {
	var file *os.File
	var err error

	if w.config.Output.OutputStrategy == config.SingleFile {
		file, err = os.Create(w.config.Output.Path)
		if err != nil {
			w.logger.Error("Failed to create single output file", "path", w.config.Output.Path, "error", err)
			return err
		}
		w.logger.Info("Single output file created", "path", w.config.Output.Path)
		defer file.Close()
	}

	for structName, fields := range w.structs {
		if w.config.Output.OutputStrategy == config.FilePerStruct {
			file, err = os.Create(filepath.Join(w.config.Output.Path, w.GetFileName(structName)+".json"))
			if err != nil {
				w.logger.Error("Failed to create output file for struct", "structName", structName, "error", err)
				return err
			}
			w.logger.Info("Output file created for struct", "structName", structName)
			defer file.Close()
		}

		// writing each struct w.config.Generation.Count times
		count := w.config.Generation.Count
		structs := make([]map[string]any, count)
		for i := 0; i < count; i++ {
			// write the struct to the file
			toWrite := make(map[string]any)
			w.logger.Debug("Writing struct instance", "structName", structName, "instance", i+1)
			for _, field := range fields {
				generator, err := field.ToGenerator(w.config.Generation.RandSeed, w.logger)
				if err != nil {
					w.logger.Error("Failed to get generator for field", "fieldName", field.Name, "error", err)
					return err
				}
				fieldValue, err := generator.EvaluateAny()
				if err != nil {
					w.logger.Error("Failed to evaluate generator for field", "fieldName", field.Name, "error", err)
					return err
				}
				toWrite[field.Name] = fieldValue
			}
			structs[i] = toWrite
		}
		jsonStruct, err := json.MarshalIndent(structs, "", " ")
		if err != nil {
			w.logger.Error("Failed to marshal JSON", "structName", structName, "error", err)
			return err
		}
		w.logger.Debug("Writing JSON data to file", "structName", structName)
		if _, err := file.Write(jsonStruct); err != nil {
			w.logger.Error("Failed to write JSON to file", "structName", structName, "error", err, "json", jsonStruct)
			return err
		}
		w.logger.Info("Successfully wrote JSON output for struct", "structName", structName)
	}
	return nil
}
