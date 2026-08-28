// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRoutes проверяет регистрацию маршрутов через Handle и вызов через ServeHTTP.
func TestRoutes(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Test Handle"))
	})

	r := NewRouter()
	r.Handle(http.MethodGet, "/", testHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, res.StatusCode)
	}
}
