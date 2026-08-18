// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

// Package router provides a simple custom HTTP router
// for handling requests in Go applications.

package router

import (
	"log"
	"net/http"
)

// RouterHandle stores a mapping of URL paths to handler functions.
// It implements the http.Handler interface.
type RouterHandle struct {
	Routes map[string]http.HandlerFunc
}

// NewRouterHandle initializes a new, clean RouterHandle.
func NewRouterHandle() *RouterHandle {
	return &RouterHandle{
		Routes: make(map[string]http.HandlerFunc),
	}
}

// ServeHTTP dispatches incoming requests to the appropriate handler
// based on the request path. It includes panic recovery to ensure
// the server remains stable if a handler crashes, returning a 500 status.
// If no handler is found, it returns 404 Not Found.
func (r *RouterHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Critical error: %v", err)
			http.Error(w, "Something broke on the server.", 500)
		}
	}()

	// Realization CORS request
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPRIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if handler, ok := r.Routes[req.URL.Path]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// AddRoute registers a handler for the given path and returns the router handle.
func (r *RouterHandle) AddRoute(path string, handler http.HandlerFunc) *RouterHandle {
	r.Routes[path] = handler
	return r
}
