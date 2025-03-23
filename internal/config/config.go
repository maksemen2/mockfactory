package config

type OutputStrategy int

const (
	FilePerStruct OutputStrategy = iota // Generates a file per struct
	SingleFile                          // Writes all structs to single file
)

type FieldIgnoreStrategy int

const (
	IgnoreUntagged FieldIgnoreStrategy = iota // Ignores all fields with no "mock" tag
	IgnoreWithTag                             // Ignores all fields with "mock:ignore" tag
	IgnoreAll                                 // IgnoreUntagged && IgnoreWithTag (ignores fields with no "mock" tag and with "mock:ignore" tag) (default)
	IncludeAll                                // Include all fields
)

type Config struct {
	InputPath  string           `validate:"required,input_file"` // Path to input file.
	Generation GenerationConfig `validate:"required"`
	Output     OutputConfig     `validate:"required"`
	Fields     FieldsConfig     `validate:"required"`
	Logging    LoggingConfig
}

type GenerationConfig struct {
	StructNames []string // Names of structs to generate. If empty, all structs will be generated
	Count       int      `validate:"min=1"` // Count of mocks to generate per struct
	RandSeed    int64    // Seed for random values
	Format      string   `validate:"oneof=json"` // Format of output file. Currently only JSON is supported
}

type OutputConfig struct {
	Path             string         `validate:"required"`               // Path to write output files. Can be a file path or folder
	OutputStrategy   OutputStrategy `validate:"required,file_strategy"` // Strategy for writing output files
	FileNameTemplate string         `validate:"file_name_template"`     // Template for file names. Can contain the following placeholders: {struct} - struct name, {count} - count of mocks
}

type FieldsConfig struct {
	IgnoreStrategy FieldIgnoreStrategy `validate:"required,ignore_strategy"`
}

type LoggingConfig struct {
	Level string
}
