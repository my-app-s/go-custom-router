// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 PANIC RECOVERED: %v\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("📥 [START] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("📤 [DONE]  %s %s | Duration: %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func ContentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func (r *Router) HandlerAPI() http.Handler {
	var handler http.Handler = r

	if r.RateLimiter != nil {
		handler = RateLimitMiddleware(r.RateLimiter)(handler)
	}

	if r.CORS != nil {
		handler = r.CORS.CorsMiddleware(handler)
	} else {
		handler = CorsMiddlewareOpen(handler)
	}

	handler = ContentTypeJSONMiddleware(handler)
	handler = RecoveryMiddleware(handler)
	handler = LoggerMiddleware(handler)

	return handler
}
