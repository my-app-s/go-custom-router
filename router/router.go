// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net/http"
	"path"
)

type Router struct {
	routes      map[string]map[string]http.Handler
	CORS        *CORS
	RateLimiter *RateLimiter
}

// NewRouter инициализирует новый экземпляр Router.
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.Handler),
	}
}

// Handle регистрирует новый маршрут.
func (r *Router) Handle(method, pathStr string, handler http.Handler) *Router {
	cleanPath := path.Clean(pathStr)
	if r.routes[cleanPath] == nil {
		r.routes[cleanPath] = make(map[string]http.Handler)
	}
	r.routes[cleanPath][method] = handler
	return r
}

// ServeHTTP реализует интерфейс http.Handler для диспетчеризации запросов.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cleanPath := path.Clean(req.URL.Path)

	methods, pathExists := r.routes[cleanPath]
	if !pathExists {
		http.NotFound(w, req)
		return
	}

	handler, methodExists := methods[req.Method]
	if !methodExists {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	handler.ServeHTTP(w, req)
}
