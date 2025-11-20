package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        if _, err := w.Write([]byte(`{"status": "ok", "service": "notification-service"}`)); err != nil {
            log.Printf("Error writing response: %v", err)
        }
    })

    log.Println("Starting notification-service on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}