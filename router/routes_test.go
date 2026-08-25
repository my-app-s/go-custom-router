// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net/http"
	"testing"
)

// TestAddRoute_And_FluentAPI проверяет регистрацию маршрутов через AddRoute и helpers (GET, POST, PUT, DELETE).
func TestAddRoute_And_FluentAPI(t *testing.T) {
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {}

	r := &RouterHandle{
		Routes: make(map[string]map[string]http.HandlerFunc),
	}

	// Проверяем работу цепочки вызовов (Fluent API)
	r.GET("/users", dummyHandler).
		POST("/users", dummyHandler).
		PUT("/users/1", dummyHandler).
		DELETE("/users/1", dummyHandler).
		AddRoute(http.MethodPatch, "/users/1", dummyHandler)

	tests := []struct {
		path   string
		method string
	}{
		{"/users", http.MethodGet},
		{"/users", http.MethodPost},
		{"/users/1", http.MethodPut},
		{"/users/1", http.MethodDelete},
		{"/users/1", http.MethodPatch},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			methodsMap, exists := r.Routes[tt.path]
			if !exists {
				t.Fatalf("expected path %q to exist in Routes", tt.path)
			}

			handler, methodExists := methodsMap[tt.method]
			if !methodExists {
				t.Fatalf("expected method %q to exist for path %q", tt.method, tt.path)
			}

			if handler == nil {
				t.Errorf("expected non-nil handler for %s %s", tt.method, tt.path)
			}
		})
	}
}
