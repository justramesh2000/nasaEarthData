package api

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "path/filepath"
    "strings"
    . "nasa-hdf-processor/internal/processor"
    . "nasa-hdf-processor/internal/config"
)

// global processor instance
var processor *BatchProcessor
var config *Configuration

func main() {
    // Set fixed config for simplicity, adjust paths as needed
    config := &Configuration{
        InputDir:  "./input",
        OutputDir: "./output",
    }

    // Create processor instance
    processor = NewBatchProcessor(config)

    // Simple handlers
    http.HandleFunc("/process", ProcessHandler)
    http.HandleFunc("/api/data", DataHandler)

    log.Println("Starting API server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

// processHandler runs batch processing once
func ProcessHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Processing triggered via API")
    err := processor.ProcessFiles()
    if err != nil {
        log.Printf("Processing error: %v", err)
        http.Error(w, "Processing failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"processing completed"}`))
}

// dataHandler reads all JSON files in output directory and returns combined JSON array
func DataHandler(w http.ResponseWriter, r *http.Request) {
    files, err := ioutil.ReadDir(config.OutputDir)
    if err != nil {
        http.Error(w, "Failed to read output directory: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var combined []map[string]interface{}

    for _, file := range files {
        if file.IsDir() {
            continue
        }
        if strings.HasSuffix(file.Name(), ".json") {
            path := filepath.Join(config.OutputDir, file.Name())
            data, err := ioutil.ReadFile(path)
            if err != nil {
                log.Printf("Failed to read file %s: %v", path, err)
                continue
            }
            var obj map[string]interface{}
            if err := json.Unmarshal(data, &obj); err != nil {
                log.Printf("Failed to parse JSON file %s: %v", path, err)
                continue
            }
            combined = append(combined, obj)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(combined)
}
