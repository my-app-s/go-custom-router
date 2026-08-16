// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	r := NewRouterHandle()

	// Регистрируем тестовые эндпоинты перед проверкой
	r.AddRoute("/", func(w http.ResponseWriter, req *http.Request) {
		s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s\n", req.Method, req.Host, req.URL.Path)
		w.Write([]byte(s))
	})
	r.AddRoute("/time", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("16.08.2026 12:00:00"))
	})

	// Табличное тестирование (Table-driven tests)
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Home Page",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Method: GET",
		},
		{
			name:           "Time Page",
			method:         "GET",
			path:           "/time",
			expectedStatus: http.StatusOK,
			expectedBody:   ".", // Проверяем наличие точек в дате/времени
		},
		{
			name:           "404 Not Found",
			method:         "GET",
			path:           "/unknown",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	r := NewRouterHandle()

	// Регистрируем эндпоинт, который намеренно падает
	r.AddRoute("/crashtest", func(w http.ResponseWriter, req *http.Request) {
		var list []int
		fmt.Println(list[99]) // паника из-за выхода за границы слайса
	})

	t.Run("Recovery from CrashTest", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/crashtest", nil)
		rr := httptest.NewRecorder()

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("The router did not recover from panic!")
			}
		}()

		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500 after panic, got %d", rr.Code)
		}

		expected := "Something broke on the server."
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("expected error message %q, got %q", expected, rr.Body.String())
		}
	})
}

func TestAddRoute(t *testing.T) {
	r := NewRouterHandle()
	path := "/custom"

	r.AddRoute(path, func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("custom handler"))
	})

	req := httptest.NewRequest("GET", path, nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	if rr.Body.String() != "custom handler" {
		t.Errorf("expected 'custom handler', got %q", rr.Body.String())
	}
}

// BenchmarkRouter измеряет скорость обработки запросов роутером и количество аллокаций памяти.
func BenchmarkRouter(b *testing.B) {
	r := NewRouterHandle()

	r.AddRoute("/bench", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("benchmark ok"))
	})

	req := httptest.NewRequest("GET", "/bench", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
	}
}
