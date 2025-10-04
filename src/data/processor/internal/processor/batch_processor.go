package processor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	. "nasa-hdf-processor/internal/config"
)

// BatchProcessor handles batch processing of HDF files
type BatchProcessor struct {
	config    *Configuration
	hdfReader *HDFReader
	processed map[string]time.Time
	mutex     sync.RWMutex
	stopCh    chan struct{}
}

// ProcessingResult represents the result of processing a file
type ProcessingResult struct {
	FileName     string    `json:"file_name"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message,omitempty"`
	ProcessedAt  time.Time `json:"processed_at"`
	OutputFile   string    `json:"output_file,omitempty"`
	FileSize     int64     `json:"file_size,omitempty"`
	Duration     string    `json:"duration,omitempty"`
}

// NewBatchProcessor creates a new batch processor instance
func NewBatchProcessor(config *Configuration) *BatchProcessor {
	return &BatchProcessor{
		config:    config,
		hdfReader: NewHDFReader(),
		processed: make(map[string]time.Time),
		stopCh:    make(chan struct{}),
	}
}

// Start begins continuous processing
func (bp *BatchProcessor) Start() {
	log.Println("Starting batch processor...")

	// Process files immediately on startup
	if err := bp.ProcessFiles(); err != nil {
		log.Printf("Initial processing failed: %v", err)
	}

	// Set up ticker for periodic processing
	ticker := time.NewTicker(bp.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := bp.ProcessFiles(); err != nil {
				log.Printf("Processing failed: %v", err)
			}
		case <-bp.stopCh:
			log.Println("Batch processor stopped")
			return
		}
	}
}

// Stop stops the batch processor
func (bp *BatchProcessor) Stop() {
	close(bp.stopCh)
}

// ProcessFiles processes all HDF files in the input directory
func (bp *BatchProcessor) ProcessFiles() error {
	log.Println("Starting file processing cycle...")

	// Get list of HDF files
	files, err := ListHDFFiles(bp.config.InputDir)
	if err != nil {
		return fmt.Errorf("failed to list HDF files: %v", err)
	}

	if len(files) == 0 {
		log.Println("No HDF files found to process")
		return nil
	}

	log.Printf("Found %d HDF files to process", len(files))

	var results []ProcessingResult
	var wg sync.WaitGroup
	resultCh := make(chan ProcessingResult, len(files))

	// Process files concurrently (limit to avoid overwhelming the system)
	maxWorkers := 3
	semaphore := make(chan struct{}, maxWorkers)

	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := bp.processFile(filePath)
			resultCh <- result
		}(file)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	for result := range resultCh {
		results = append(results, result)
	}

	// Generate summary report
	//bp.generateReport(results)

	return nil
}

// processFile processes a single HDF file
func (bp *BatchProcessor) processFile(filePath string) ProcessingResult {
	startTime := time.Now()

	result := ProcessingResult{
		FileName:    filepath.Base(filePath),
		ProcessedAt: startTime,
	}

	// Check if file was already processed recently (within 1 hour)
	bp.mutex.RLock()
	if lastProcessed, exists := bp.processed[filePath]; exists {
		if time.Since(lastProcessed) < time.Hour {
			bp.mutex.RUnlock()
			log.Printf("Skipping %s (already processed recently)", result.FileName)
			result.Success = true
			result.Duration = time.Since(startTime).String()
			return result
		}
	}
	bp.mutex.RUnlock()

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("failed to get file info: %v", err)
		return result
	}
	result.FileSize = info.Size()

	// Check file size limit (100MB)
	if result.FileSize > 100*1024*1024 {
		result.ErrorMessage = "file too large (>100MB)"
		return result
	}

	// Convert HDF to JSON
	jsonData, err := bp.hdfReader.ProcessFile(filePath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("conversion failed: %v", err)
		return result
	}

	// Generate output filename
	outputFile := bp.generateOutputPath(filePath)

	// Write JSON to output file
	if err := bp.writeJSONFile(outputFile, jsonData); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed to write output file: %v", err)
		return result
	}

	result.Success = true
	result.OutputFile = outputFile
	result.Duration = time.Since(startTime).String()

	// Mark as processed
	bp.mutex.Lock()
	bp.processed[filePath] = startTime
	bp.mutex.Unlock()

	log.Printf("Successfully processed %s -> %s (%s)",
		result.FileName, filepath.Base(result.OutputFile), result.Duration)

	return result
}

// generateOutputPath generates the output file path for an HDF file
func (bp *BatchProcessor) generateOutputPath(inputFile string) string {
	// Get relative path from input directory
	relPath, err := filepath.Rel(bp.config.InputDir, inputFile)
	if err != nil {
		relPath = filepath.Base(inputFile)
	}

	// Change extension to .json
	ext := filepath.Ext(relPath)
	baseName := strings.TrimSuffix(relPath, ext)

	// Add timestamp to avoid conflicts
	timestamp := time.Now().Format("20060102_150405")
	outputPath := filepath.Join(bp.config.OutputDir, fmt.Sprintf("%s_%s.json", baseName, timestamp))

	return outputPath
}

// writeJSONFile writes JSON data to a file
func (bp *BatchProcessor) writeJSONFile(filePath string, jsonData []byte) error {
	// Ensure output directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	//cleanup old files
	bp.CleanupOldFiles()
	// Write file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// GetProcessedFiles returns the list of processed files
func (bp *BatchProcessor) GetProcessedFiles() map[string]time.Time {
	bp.mutex.RLock()
	defer bp.mutex.RUnlock()

	result := make(map[string]time.Time)
	for file, timestamp := range bp.processed {
		result[file] = timestamp
	}
	return result
}

// CleanupOldFiles moves all processed files to an "old" directory
func (bp *BatchProcessor) CleanupOldFiles() error {
	log.Printf("Moving all processed files to old directory")

	moved := 0

	// Create old directory if it doesn't exist
	if err := os.MkdirAll(bp.config.OldDir, 0755); err != nil {
		return fmt.Errorf("failed to create old directory: %v", err)
	}

	err := filepath.Walk(bp.config.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the old directory itself
		if strings.Contains(path, "/old/") {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			// Create subdirectory in old folder based on date
			dateDir := info.ModTime().Format("2006-01-02")
			targetDir := filepath.Join(bp.config.OldDir, dateDir)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				log.Printf("Failed to create date directory %s: %v", targetDir, err)
				return nil
			}

			// Move file to old directory
			targetPath := filepath.Join(targetDir, filepath.Base(path))
			if err := os.Rename(path, targetPath); err != nil {
				log.Printf("Failed to move file %s to %s: %v", path, targetPath, err)
			} else {
				moved++
				log.Printf("Moved file: %s -> %s", filepath.Base(path), targetPath)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("cleanup failed: %v", err)
	}

	log.Printf("Cleanup completed: moved %d files to %s", moved, bp.config.OldDir)
	return nil
}

// GetStatus returns the current status of the batch processor
func (bp *BatchProcessor) GetStatus() map[string]interface{} {
	bp.mutex.RLock()
	defer bp.mutex.RUnlock()

	return map[string]interface{}{
		"processed_files":     len(bp.processed),
		"input_directory":     bp.config.InputDir,
		"output_directory":    bp.config.OutputDir,
		"processing_interval": bp.config.Interval.String(),
		"last_check":          time.Now(),
	}
}

// MoveToOldDirectory manually moves a specific file to the old directory
func (bp *BatchProcessor) MoveToOldDirectory(filePath string) error {
	// Create old directory if it doesn't exist
	if err := os.MkdirAll(bp.config.OldDir, 0755); err != nil {
		return fmt.Errorf("failed to create old directory: %v", err)
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Create date-based subdirectory
	dateDir := info.ModTime().Format("2006-01-02")
	targetDir := filepath.Join(bp.config.OldDir, dateDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create date directory: %v", err)
	}

	// Move file
	targetPath := filepath.Join(targetDir, filepath.Base(filePath))
	if err := os.Rename(filePath, targetPath); err != nil {
		return fmt.Errorf("failed to move file: %v", err)
	}

	log.Printf("Moved file: %s -> %s", filepath.Base(filePath), targetPath)
	return nil
}
