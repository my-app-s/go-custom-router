// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type inputCORS struct {
	origin  string
	methods string
	headers string
}

type wantCORS struct {
	allowed bool
	origin  string
	methods string
	headers string
}

var defaultInput = inputCORS{
	origin:  "https://my-app.com",
	methods: "GET, POST, PUT, DELETE, OPTIONS",
	headers: "Content-Type, Authorization",
}

func TestNewCORS(t *testing.T) {
	cors := NewCORS()

	if cors == nil {
		t.Fatal("NewCORS() returned nil")
	}

	if cors.AllowedOrigins == nil {
		t.Error("AllowedOrigins map is nil")
	}
}

func TestAddOrigin(t *testing.T) {
	testCases := []struct {
		name        string
		input       inputCORS
		checkOrigin string
		wantSuccess bool
		want        wantCORS
	}{
		{
			name:        "Valid Origin Default",
			input:       defaultInput,
			checkOrigin: "https://my-app.com",
			wantSuccess: true,
			want: wantCORS{
				allowed: true,
				origin:  "https://my-app.com",
				methods: "GET, POST, PUT, DELETE, OPTIONS",
				headers: "Content-Type, Authorization",
			},
		},
		{
			name:        "Empty Origin",
			input:       inputCORS{origin: "", methods: "GET", headers: "Content-Type"},
			checkOrigin: "",
			wantSuccess: false,
			want:        wantCORS{allowed: false},
		},
		{
			name:        "Empty Methods",
			input:       inputCORS{origin: "https://my-app.com", methods: "", headers: "Content-Type"},
			checkOrigin: "https://my-app.com",
			wantSuccess: false,
			want:        wantCORS{allowed: false},
		},
		{
			name:        "Empty Headers",
			input:       inputCORS{origin: "https://my-app.com", methods: "GET", headers: ""},
			checkOrigin: "https://my-app.com",
			wantSuccess: false,
			want:        wantCORS{allowed: false},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cors := NewCORS()

			gotSuccess := cors.AddOrigin(tt.input.origin, tt.input.methods, tt.input.headers)
			if gotSuccess != tt.wantSuccess {
				t.Errorf("AddOrigin() = %v; want %v", gotSuccess, tt.wantSuccess)
			}

			gotAllowed := cors.IsAllowed(tt.checkOrigin)
			if gotAllowed != tt.want.allowed {
				t.Errorf("IsAllowed(%q) = %v; want %v", tt.checkOrigin, gotAllowed, tt.want.allowed)
			}

			if tt.wantSuccess {
				opt, exists := cors.AllowedOrigins[tt.want.origin]
				if !exists {
					t.Fatalf("origin %q not found in AllowedOrigins map", tt.want.origin)
				}
				if opt.AllowedMethods != tt.want.methods {
					t.Errorf("AllowedMethods = %q; want %q", opt.AllowedMethods, tt.want.methods)
				}
				if opt.AllowedHeaders != tt.want.headers {
					t.Errorf("AllowedHeaders = %q; want %q", opt.AllowedHeaders, tt.want.headers)
				}
			}
		})
	}
}

func TestMultiAddOrigin(t *testing.T) {
	cors := NewCORS()
	cors.AddOrigin(defaultInput.origin, defaultInput.methods, defaultInput.headers)
	cors.AddOrigin(defaultInput.origin, defaultInput.methods, defaultInput.headers)

	count := len(cors.AllowedOrigins)
	if count != 1 {
		t.Errorf("Add origin duplicate count = %v; want: %d", count, 1)
	}
}

func TestCorsMiddleware(t *testing.T) {
	cors := NewCORS()
	cors.AddOrigin("https://my-app.com", "GET, POST", "Content-Type")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := cors.CorsMiddleware(nextHandler)

	t.Run("Allowed Origin Preflight", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
		req.Header.Set("Origin", "https://my-app.com")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://my-app.com" {
			t.Errorf("Access-Control-Allow-Origin = %q; want %q", got, "https://my-app.com")
		}
		if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "GET, POST" {
			t.Errorf("Access-Control-Allow-Methods = %q; want %q", got, "GET, POST")
		}
	})

	t.Run("Forbidden Origin Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "https://untrusted.com")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("Access-Control-Allow-Origin = %q; want empty", got)
		}
	})
}

func TestCorsMiddlewareOpen(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CorsMiddlewareOpen(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/public", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q; want %q", got, "*")
	}
}
