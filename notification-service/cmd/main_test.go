package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestValidateEmail(t *testing.T) {
    if !ValidateEmail("test@example.com") {
        t.Errorf("expected valid email")
    }
    if ValidateEmail("bad-email") {
        t.Errorf("expected invalid email")
    }
}

func TestNotifyHandler_Success(t *testing.T) {
    body := []byte(`{"email":"test@example.com","message":"hi"}`)
    req := httptest.NewRequest("POST", "/notify", bytes.NewBuffer(body))
    w := httptest.NewRecorder()

    NotifyHandler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}
