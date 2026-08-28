// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"log"
	"net/http"
)

// CORS хранит конфигурацию для обработки междоменных запросов (Cross-Origin Resource Sharing).
type CORS struct {
	AllowedOrigins map[string]OriginsOptions `json:"allowedOrigins"`
}

type OriginsOptions struct {
	AllowedMethods string `json:"allowedMethods"`
	AllowedHeaders string `json:"allowedHeaders"`
}

func NewCORS() *CORS {
	return &CORS{
		AllowedOrigins: make(map[string]OriginsOptions),
	}
}

// AddOrigin регистрирует разрешенный домен.
func (c *CORS) AddOrigin(origin, methods, headers string) bool {
	if origin == "" || methods == "" || headers == "" {
		log.Println("Error input: Empty parameters for CORS origin")
		return false
	}

	if c.AllowedOrigins == nil {
		c.AllowedOrigins = make(map[string]OriginsOptions)
	}

	c.AllowedOrigins[origin] = OriginsOptions{
		AllowedMethods: methods,
		AllowedHeaders: headers,
	}

	return true
}

// IsAllowed проверяет, разрешен ли домен.
func (c *CORS) IsAllowed(origin string) bool {
	if c == nil || c.AllowedOrigins == nil {
		return false
	}
	_, ok := c.AllowedOrigins[origin]
	return ok
}

// CorsMiddleware обрабатывает CORS для настроенных доменов.
func (c *CORS) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if c != nil && c.IsAllowed(origin) {
			opt := c.AllowedOrigins[origin]
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", opt.AllowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", opt.AllowedHeaders)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// CorsMiddlewareOpen предоставляет открытый доступ для любых источников (Origin: *).
func CorsMiddlewareOpen(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
