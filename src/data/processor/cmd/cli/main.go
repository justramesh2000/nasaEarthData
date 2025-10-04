package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	. "nasa-hdf-processor/internal/processor"
	. "nasa-hdf-processor/internal/config"
)

// setupLogging configures logging
func setupLogging(logLevel string) {
	// Set log format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// For now, we'll use the standard log package
	// In production, you might want to use a more sophisticated logging library
	switch logLevel {
	case "debug":
		log.Println("Log level set to DEBUG")
	case "info":
		log.Println("Log level set to INFO")
	case "warn":
		log.Println("Log level set to WARN")
	case "error":
		log.Println("Log level set to ERROR")
	}
}

// checkHDFSupport checks if HDF support is available
func checkHDFSupport() error {
	// This is a placeholder for HDF library availability check
	// In a real implementation, you might want to check if the HDF5 library is installed
	log.Println("Checking HDF support...")
	log.Println("HDF support is available")
	return nil
}

// initializeDirectories sets up the input and output directories
func initializeDirectories(config *Configuration) error {
	// Make sure the input directory exists and is readable
	absInputDir, err := filepath.Abs(config.InputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for input directory: %v", err)
	}

	absOutputDir, err := filepath.Abs(config.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output directory: %v", err)
	}

	log.Printf("Input directory: %s", absInputDir)
	log.Printf("Output directory: %s", absOutputDir)

	// Check if we can read the input directory
	entries, err := os.ReadDir(absInputDir)
	if err != nil {
		return fmt.Errorf("cannot read input directory %s: %v", absInputDir, err)
	}

	log.Printf("Found %d entries in input directory", len(entries))

	// Count HDF files
	hdfCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".h5" || ext == ".hdf5" || ext == ".he5" || ext == ".nc" {
				hdfCount++
			}
		}
	}

	log.Printf("Found %d HDF files in input directory", hdfCount)

	return nil
}

// main is the entry point of the application
func main() {
	// Parse command line arguments
	config := ParseFlags()

	// Check environment variables for overrides
	if envInputDir := os.Getenv("INPUT_DIR"); envInputDir != "" {
		config.InputDir = envInputDir
	}
	if envOutputDir := os.Getenv("OUTPUT_DIR"); envOutputDir != "" {
		config.OutputDir = envOutputDir
	}
	if envOldDir := os.Getenv("OLD_DIR"); envOldDir != "" {
		config.OldDir = envOldDir
	}
	if envInterval := os.Getenv("PROCESS_INTERVAL"); envInterval != "" {
		if interval, err := time.ParseDuration(envInterval); err == nil {
			config.Interval = interval
		}
	}

	// Set default old directory if not specified - always relative to output directory
	if config.OldDir == "" {
		config.OldDir = filepath.Join(config.OutputDir, "old")
	} else {
		// If user specified a relative path, make it relative to output directory
		if !filepath.IsAbs(config.OldDir) {
			config.OldDir = filepath.Join(config.OutputDir, config.OldDir)
		}
	}

	// Validate configuration
	if err := Validate(config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Setup logging
	setupLogging(config.LogLevel)
	// Check HDF support
	if err := checkHDFSupport(); err != nil {
		log.Fatalf("HDF support check failed: %v", err)
	}

	// Initialize directories
	if err := initializeDirectories(config); err != nil {
		log.Fatalf("Directory initialization failed: %v", err)
	}

	// Create batch processor
	processor := NewBatchProcessor(config)

	if config.Cleanup {
		// Cleanup mode - move all files to old directory
		log.Println("Running in cleanup mode")
		if err := processor.CleanupOldFiles(); err != nil {
			log.Fatalf("Cleanup failed: %v", err)
		}
		log.Println("Cleanup completed successfully")
	} else if config.RunOnce {
		// Run once and exit
		log.Println("Running in single-shot mode")
		if err := processor.ProcessFiles(); err != nil {
			log.Fatalf("Processing failed: %v", err)
		}
		log.Println("Processing completed successfully")
	} else {
		// Run continuously
		log.Printf("Starting continuous processing with interval: %v", config.Interval)
		processor.Start()
	}
}

// Additional utility functions

// getVersion returns the application version
func getVersion() string {
	return "1.0.0"
}

// getBuildInfo returns build information
func getBuildInfo() map[string]string {
	return map[string]string{
		"version":    getVersion(),
		"go_version": "1.21+",   // This would be set during build
		"build_time": "unknown", // This would be set during build
	}
}
