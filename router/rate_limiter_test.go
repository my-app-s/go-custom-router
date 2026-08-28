// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCustomRateLimiter_Allow проверяет логику блокировки и сброса окна.
func TestCustomRateLimiter_Allow(t *testing.T) {
	// Лимит: 2 запроса в 100мс, максимальное число уникальных IP: 1000
	limiter := NewLimiter(2, 100*time.Millisecond, 1*time.Minute, 1000)
	defer limiter.Stop() // Останавливаем воркер после теста

	ip := "192.168.1.1"

	if !limiter.Allow(ip) {
		t.Error("1st request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Error("2nd request should be allowed")
	}
	if limiter.Allow(ip) {
		t.Error("3rd request should be blocked (limit exceeded)")
	}

	// Ждем окончания окна
	time.Sleep(120 * time.Millisecond)

	if !limiter.Allow(ip) {
		t.Error("request after window expiration should be allowed")
	}
}

// TestRateLimitMiddleware проверяет работу middleware и статусы HTTP 429.
func TestRateLimitMiddleware(t *testing.T) {
	limiter := NewLimiter(1, 1*time.Second, 1*time.Minute, 1000)
	defer limiter.Stop() // Останавливаем воркер после теста

	middleware := RateLimitMiddleware(limiter)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(nextHandler)

	// Запрос 1: Разрешен (200 OK)
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "10.0.0.1:12345"
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec1.Code)
	}

	// Запрос 2 с того же IP (разные порты): Блокировка (429 Too Many Requests)
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "10.0.0.1:54321"
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too Many Requests, got %d", rec2.Code)
	}
}
