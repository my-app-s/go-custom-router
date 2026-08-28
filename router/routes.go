// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import "net/http"

func (r *Router) GET(path string, handler http.HandlerFunc) *Router {
	return r.Handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) *Router {
	return r.Handle(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) *Router {
	return r.Handle(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) *Router {
	return r.Handle(http.MethodDelete, path, handler)
}

func (r *Router) PATCH(path string, handler http.HandlerFunc) *Router {
	return r.Handle(http.MethodPatch, path, handler)
}
