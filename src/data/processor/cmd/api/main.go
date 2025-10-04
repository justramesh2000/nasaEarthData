package main

import (
    "log"
    "net/http"
    . "nasa-hdf-processor/internal/api"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/process", ProcessHandler)
    mux.HandleFunc("/api/data", DataHandler)

    log.Println("Starting REST API server on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

/*
y = mx + c 
*/