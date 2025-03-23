package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/maksemen2/mockfactory/internal/config"
	"github.com/maksemen2/mockfactory/internal/testutils"
)

func createTempFile(content string) (string, func()) {
	tmpFile, err := os.CreateTemp("", "testfile-*.go")
	if err != nil {
		panic(err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		panic(err)
	}
	tmpFile.Close()

	return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
}

func compareMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

func TestParser_ParseFile(t *testing.T) {
	testContent := `
package testdata

type User struct {
    ID    int    ` + "`mock:\"min=10;max=20\"`" + `
    Name  string ` + "`mock:\"ignore\"`" + `
    Email string
}

type Account struct {
    ID     int
    Number string
}
`
	filePath, cleanup := createTempFile(testContent)
	defer cleanup()

	tests := []struct {
		name        string
		config      *config.Config
		wantStructs map[string][]StructField
	}{
		{
			name:   "parse all structs",
			config: &config.Config{InputPath: filePath, Fields: config.FieldsConfig{config.IncludeAll}},
			wantStructs: map[string][]StructField{
				"User": {
					{Name: "ID", Type: "int", MockTags: map[string]string{"min": "10", "max": "20"}},
					{Name: "Name", Type: "string", MockTags: map[string]string{"ignore": ""}},
					{Name: "Email", Type: "string", MockTags: map[string]string{}},
				},
				"Account": {
					{Name: "ID", Type: "int", MockTags: map[string]string{}},
					{Name: "Number", Type: "string", MockTags: map[string]string{}},
				},
			},
		},
		{
			name: "filter structs by name",
			config: &config.Config{
				InputPath: filePath,
				Generation: config.GenerationConfig{
					StructNames: []string{"User"},
				},
				Fields: config.FieldsConfig{config.IncludeAll},
			},
			wantStructs: map[string][]StructField{
				"User": {
					{Name: "ID", Type: "int", MockTags: map[string]string{"min": "10", "max": "20"}},
					{Name: "Name", Type: "string", MockTags: map[string]string{"ignore": ""}},
					{Name: "Email", Type: "string", MockTags: map[string]string{}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.config, testutils.TestLogger())
			got, err := p.ParseFile()

			if err != nil {
				t.Errorf("ParseFile() error = %v", err)
				return
			}

			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.wantStructs) {
				t.Errorf("ParseFile() got = %v, want = %v", got, tt.wantStructs)
			}
		})
	}
}
