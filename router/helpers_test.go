// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestSendJSON проверяет корректность отправки JSON-ответа и заголовков.
func TestSendJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	payload := map[string]string{"status": "ok"}

	SendJSON(rec, http.StatusCreated, payload)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if res["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", res["status"])
	}
}

// TestMakeCustomHandler проверяет фабрику хэндлеров.
func TestMakeCustomHandler(t *testing.T) {
	handler := MakeCustomHandler("service", "router-api")

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	var res map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &res)

	if res["service"] != "router-api" {
		t.Errorf("expected service='router-api', got %q", res["service"])
	}
}
