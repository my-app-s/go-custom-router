// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewRouter проверяет корректность инициализации нового роутера.
func TestNewRouter(t *testing.T) {
	r := NewRouter()
	if r == nil {
		t.Fatal("expected non-nil Router")
	}
	if r.routes == nil {
		t.Error("expected initialized routes map")
	}
}

// TestServeHTTP_Routing проверяет корректность диспетчеризации запросов (200, 404, 405).
func TestServeHTTP_Routing(t *testing.T) {
	r := NewRouter()

	// Регистрируем тестовые эндпоинты
	r.GET("/api/ping", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	t.Run("200 OK - Successful Route Match", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
		if rec.Body.String() != "pong" {
			t.Errorf("expected body 'pong', got %q", rec.Body.String())
		}
	})

	t.Run("404 Not Found - Path Does Not Exist", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("405 Method Not Allowed - Unregistered HTTP Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/ping", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}
