// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Тестовая реализация ServeHTTP для RouterHandle, чтобы проверять работу цепочки.
func (r *RouterHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/panic" {
		panic("test panic error")
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// TestRecoveryMiddleware проверяет перехват паники и штатное прохождение запроса.
func TestRecoveryMiddleware(t *testing.T) {
	t.Run("Panicking Handler", func(t *testing.T) {
		panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("something crashed!")
		})

		handler := RecoveryMiddleware(panickingHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}

		if !strings.Contains(rec.Body.String(), "Internal Server Error") {
			t.Errorf("expected body to contain 'Internal Server Error', got %q", rec.Body.String())
		}
	})

	t.Run("Normal Handler", func(t *testing.T) {
		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})

		handler := RecoveryMiddleware(normalHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
		if rec.Body.String() != "OK" {
			t.Errorf("expected body 'OK', got %q", rec.Body.String())
		}
	})
}

// TestLoggerMiddleware проверяет, что логгер не нарушает прохождение HTTP-запроса.
func TestLoggerMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	handler := LoggerMiddleware(nextHandler)
	req := httptest.NewRequest(http.MethodPost, "/users/create", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}
}

// TestContentTypeJSONMiddleware проверяет автоматическую установку HTTP-заголовка JSON.
func TestContentTypeJSONMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := ContentTypeJSONMiddleware(nextHandler)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expectedHeader := "application/json; charset=utf-8"
	gotHeader := rec.Header().Get("Content-Type")

	if gotHeader != expectedHeader {
		t.Errorf("expected Content-Type %q, got %q", expectedHeader, gotHeader)
	}
}

// TestRouterHandle_Handler проверяет базовую цепочку middleware без JSON заголовка.
func TestRouterHandle_Handler(t *testing.T) {
	t.Run("Default CORS (Nil)", func(t *testing.T) {
		router := &RouterHandle{CORS: nil}
		handler := router.Handler()

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		// При CORS == nil используется CorsMiddlewareOpen -> Access-Control-Allow-Origin: *
		if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("expected CORS header '*', got %q", rec.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("Custom CORS Configured", func(t *testing.T) {
		cors := NewCORS()
		_ = cors.AddOrigin("https://example.com", "GET, POST, PUT, DELETE, OPTIONS", "Content-Type, Authorization")

		router := &RouterHandle{CORS: cors}
		handler := router.Handler()

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
			t.Errorf("expected CORS header 'https://example.com', got %q", rec.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("Panic Handling in Full Chain", func(t *testing.T) {
		router := &RouterHandle{}
		handler := router.Handler()

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d on panic, got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}

// TestRouterHandle_HandlerAPI проверяет полную цепочку middleware для JSON API.
func TestRouterHandle_HandlerAPI(t *testing.T) {
	router := &RouterHandle{CORS: nil}
	handler := router.HandlerAPI()

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Проверяем наличие заголовка Content-Type
	expectedContentType := "application/json; charset=utf-8"
	if rec.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("expected Content-Type %q, got %q", expectedContentType, rec.Header().Get("Content-Type"))
	}
}
