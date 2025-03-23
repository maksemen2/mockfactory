package root

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/maksemen2/mockfactory/internal/config"
	"github.com/maksemen2/mockfactory/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mockfactory",
	Short: "Generate mock data from Go structs",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, logger, _ := ExtractConfig(cmd)
		logger.Debug("Initializing process", "config", cfg)
		pkg.GenerateFromFile(cfg, logger)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("input", "i", "", "Path to input Go file (required)")
	rootCmd.MarkPersistentFlagRequired("input")
	rootCmd.PersistentFlags().StringSlice("structs", []string{}, "Comma-separated list of struct names")
	rootCmd.PersistentFlags().Int("count", 1, "Number of objects to generate per struct")
	rootCmd.PersistentFlags().Int64("seed", 0, "Random seed (default time.Now().UnixNano())")
	rootCmd.PersistentFlags().String("format", "json", "Output format")
	rootCmd.PersistentFlags().StringP("output", "o", ".", "Output path")
	rootCmd.PersistentFlags().String("strategy", "per-struct", "Output strategy: per-struct|single-file")
	rootCmd.PersistentFlags().String("template", "", "File name template (e.g. {struct}_{count}.json)")
	rootCmd.PersistentFlags().String("ignore", "all", "Ignore strategy: untagged|with-tag|all|none")
	rootCmd.PersistentFlags().String("log-level", "", "Log level: debug|info|warn|error")
}

func ExtractConfig(cmd *cobra.Command) (*config.Config, *slog.Logger, error) {
	var err error
	cfg := &config.Config{}

	cfg.InputPath, err = cmd.Flags().GetString("input")
	if err != nil {
		return nil, nil, err
	}

	cfg.Generation.StructNames, err = cmd.Flags().GetStringSlice("structs")
	if err != nil {
		return nil, nil, err
	}

	cfg.Generation.Count, err = cmd.Flags().GetInt("count")
	if err != nil {
		return nil, nil, err
	}

	cfg.Generation.RandSeed, err = cmd.Flags().GetInt64("seed")
	if err != nil {
		return nil, nil, err
	}

	cfg.Generation.Format, err = cmd.Flags().GetString("format")
	if err != nil {
		return nil, nil, err
	}

	cfg.Output.Path, err = cmd.Flags().GetString("output")
	if err != nil {
		return nil, nil, err
	}

	strategy, err := cmd.Flags().GetString("strategy")
	if err != nil {
		return nil, nil, err
	}
	switch strategy {
	case "per-struct":
		cfg.Output.OutputStrategy = config.FilePerStruct
	case "single-file":
		cfg.Output.OutputStrategy = config.SingleFile
	default:
		return nil, nil, fmt.Errorf("invalid file strategy: %s", strategy)
	}

	cfg.Output.FileNameTemplate, err = cmd.Flags().GetString("template")
	if err != nil {
		return nil, nil, err
	}

	ignore, err := cmd.Flags().GetString("ignore")
	if err != nil {
		return nil, nil, err
	}
	switch ignore {
	case "untagged":
		cfg.Fields.IgnoreStrategy = config.IgnoreUntagged
	case "with-tag":
		cfg.Fields.IgnoreStrategy = config.IgnoreWithTag
	case "all":
		cfg.Fields.IgnoreStrategy = config.IgnoreAll
	case "none":
		cfg.Fields.IgnoreStrategy = config.IncludeAll
	default:
		return nil, nil, fmt.Errorf("invalid ignore strategy: %s", ignore)
	}

	logLevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return nil, nil, err
	}

	slogLevel := getLogLevel(logLevel)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))

	return cfg, logger, nil
}

func getLogLevel(str string) slog.Level {
	switch str {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
