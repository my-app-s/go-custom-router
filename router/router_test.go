package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	r := NewRouterHandle()

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
			// Создаем имитацию запроса
			req := httptest.NewRequest(tt.method, tt.path, nil)
			// Создаем инструмент для записи ответа
			rr := httptest.NewRecorder()

			// Вызываем ServeHTTP напрямую
			r.ServeHTTP(rr, req)

			// Проверяем статус-код
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Проверяем, содержит ли тело ответа ожидаемую строку
			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	r := NewRouterHandle()

	t.Run("Recovery from CrashTest", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/crashtest", nil)
		rr := httptest.NewRecorder()

		// Оборачиваем в проверку, чтобы тест не завершился аварийно сам по себе
		// (хотя ServeHTTP уже содержит recover, это двойная проверка)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The router did not recover from panic!")
			}
		}()

		r.ServeHTTP(rr, req)

		// Проверяем, что роутер перехватил панику и вернул 500
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
