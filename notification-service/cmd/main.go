package main

import (
    "encoding/json"
    "log"
    "net/http"
    "regexp"
    "time"
)

var emailRegex = regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)

type NotificationRequest struct {
    Email   string `json:"email"`
    Message string `json:"message"`
}

func ValidateEmail(email string) bool {
    return emailRegex.MatchString(email)
}

func NotifyHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
    defer func() {
        if err := r.Body.Close(); err != nil {
            log.Printf("error closing request body: %v", err)
        }
    }()

    var req NotificationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    if !ValidateEmail(req.Email) {
        http.Error(w, "invalid email", http.StatusBadRequest)
        return
    }

    if req.Message == "" {
        http.Error(w, "mmmmessage cannot be empty", http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    if _, err := w.Write([]byte(`{"status": "sent"}`)); err != nil {
        log.Printf("write error: %v", err)
    }
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        if _, err := w.Write([]byte(`{"status": "ok", "service": "notification-service"}`)); err != nil {
            log.Printf("write error: %v", err)
        }
    })

    mux.HandleFunc("/notify", NotifyHandler)

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    log.Println("Starting notification-service on :8080")
    log.Fatal(srv.ListenAndServe())
}
