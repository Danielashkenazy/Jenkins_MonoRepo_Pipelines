package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status": "ok", "service": "notification-service"}`))
    })

    fmt.Println("notification-service running on port 8080")
    fmt.Println("Welcome to the Notification Serviceeeeeeeeeeeeeeeawwewewweweweeweeawewewewewewewe")
    http.ListenAndServe(":8080", nil)
}
