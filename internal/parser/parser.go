package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"slices"
	"strings"

	"log/slog"

	"github.com/maksemen2/mockfactory/internal/config"
)

type Parser struct {
	config *config.Config
	logger *slog.Logger
}

func NewParser(cfg *config.Config, logger *slog.Logger) *Parser {
	return &Parser{config: cfg, logger: logger}
}

// ParseFile parses the file at the configured InputPath and returns a map of struct names to their fields.
func (p *Parser) ParseFile() (map[string][]StructField, error) {
	p.logger.Debug("Starting file parsing", "inputPath", p.config.InputPath)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, p.config.InputPath, nil, parser.ParseComments)
	if err != nil {
		p.logger.Error("Failed to parse file", "inputPath", p.config.InputPath, "error", err)
		return nil, err
	}

	structs := make(map[string][]StructField) // structName -> fields

	// Iterate over all declarations in the file.
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue // skip non-type declarations
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structName := typeSpec.Name.Name
			// Only parse required structs based on configuration.
			if !p.shouldParseStruct(structName) {
				p.logger.Debug("Skipping struct based on configuration", "structName", structName)
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				p.logger.Warn("TypeSpec is not a struct; skipping", "structName", structName)
				continue // skip if not a struct
			}

			fields := p.extractFields(structType)
			structs[structName] = fields
			p.logger.Debug("Parsed struct", "structName", structName, "fieldCount", len(fields))
		}
	}
	p.logger.Info("File parsed successfully", "structCount", len(structs))
	return structs, nil
}

func (p *Parser) extractFields(structType *ast.StructType) []StructField {
	var fields []StructField
	for _, field := range structType.Fields.List {
		fieldType := fmt.Sprintf("%s", field.Type)

		var mockTags map[string]string
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			mockTags = parseMockTags(tag.Get("mock"))
		}

		if p.shouldAddField(mockTags) {
			for _, name := range field.Names {
				fields = append(fields, StructField{
					Name:     name.Name,
					Type:     fieldType,
					MockTags: mockTags,
				})
			}
		}
	}
	return fields
}

func (p *Parser) shouldParseStruct(structName string) bool {
	if p.config.Generation.StructNames == nil || len(p.config.Generation.StructNames) == 0 {
		return true
	}
	return slices.Contains(p.config.Generation.StructNames, structName)
}

func (p *Parser) shouldAddField(tags map[string]string) bool {
	switch p.config.Fields.IgnoreStrategy {
	case config.IncludeAll:
		return true
	case config.IgnoreUntagged:
		return len(tags) != 0
	case config.IgnoreWithTag:
		_, ok := tags["ignore"]
		return !ok
	case config.IgnoreAll:
		_, ok := tags["ignore"]
		return len(tags) != 0 || !ok
	default:
		return true
	}
}

func parseMockTags(tag string) map[string]string {
	result := make(map[string]string)
	if tag == "" {
		return result
	}

	pairs := strings.Split(tag, ";")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else if len(kv) == 1 {
			result[kv[0]] = ""
		}
	}

	return result
}
