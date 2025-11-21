package main

import (
    "encoding/json"
    "log"
    "net/http"
    "regexp"
)

type NotificationRequest struct {
    Email   string `json:"email"`
    Message string `json:"message"`
}

func ValidateEmail(email string) bool {
    regex := regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
    return regex.MatchString(email)
}

func NotifyHandler(w http.ResponseWriter, r *http.Request) {
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
        http.Error(w, "message cannot be empty", http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "sent"}`))
}

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"status": "ok", "service": "notification-service"}`))
    })

    http.HandleFunc("/notify", NotifyHandler)

    log.Println("Starting notification-service on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
