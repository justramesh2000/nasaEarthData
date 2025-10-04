package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// Configuration holds application configuration
type Configuration struct {
	InputDir     string
	OutputDir    string
	OldDir       string
	LogLevel     string
	Interval     time.Duration
	Help         bool
	RunOnce      bool
	Cleanup      bool
	OutputFormat string
}

// ParseFlags parses CLI flags and returns a Configuration
func ParseFlags() *Configuration {
	cfg := &Configuration{}
	flag.StringVar(&cfg.InputDir, "input-dir", "./input", "Directory containing HDF files to process")
	flag.StringVar(&cfg.OutputDir, "output-dir", "./output", "Directory to save converted JSON files")
	flag.StringVar(&cfg.OldDir, "old-dir", "", "Old directory name (relative to output-dir)")
	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.DurationVar(&cfg.Interval, "interval", 5*time.Minute, "Processing interval")
	flag.BoolVar(&cfg.Help, "help", false, "Show help")
	flag.BoolVar(&cfg.RunOnce, "once", false, "Run once and exit")
	flag.BoolVar(&cfg.Cleanup, "cleanup", false, "Cleanup old files and exit")
	flag.StringVar(&cfg.OutputFormat, "format", "json", "Output format")

	flag.Parse()
	return cfg
}

func Validate(cfg *Configuration) error {
	if err := os.MkdirAll(cfg.InputDir, 0755); err != nil {
		return fmt.Errorf("failed to create input dir: %v", err)
	}
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %v", err)
	}
	if cfg.Interval <= 0 {
		return fmt.Errorf("interval must be > 0")
	}
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, lvl := range validLevels {
		if cfg.LogLevel == lvl {
			return nil
		}
	}
	return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
}
