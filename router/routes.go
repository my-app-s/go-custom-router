// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import "net/http"

// AddRoute регистрирует обработчик (http.HandlerFunc) для указанного HTTP-метода и пути.
// Инициализирует внутреннюю карту маршрутов, если она пуста, и возвращает *RouterHandle для поддержки Fluent API.
func (r *RouterHandle) AddRoute(method, path string, handler http.HandlerFunc) *RouterHandle {
	if r.Routes[path] == nil {
		r.Routes[path] = make(map[string]http.HandlerFunc)
	}
	r.Routes[path][method] = handler
	return r
}

// GET — удобный метод-помощник для регистрации обработчика GET-запроса.
func (r *RouterHandle) GET(path string, handler http.HandlerFunc) *RouterHandle {
	return r.AddRoute(http.MethodGet, path, handler)
}

// POST — удобный метод-помощник для регистрации обработчика POST-запроса.
func (r *RouterHandle) POST(path string, handler http.HandlerFunc) *RouterHandle {
	return r.AddRoute(http.MethodPost, path, handler)
}

// PUT — удобный метод-помощник для регистрации обработчика PUT-запроса.
func (r *RouterHandle) PUT(path string, handler http.HandlerFunc) *RouterHandle {
	return r.AddRoute(http.MethodPut, path, handler)
}

// DELETE — удобный метод-помощник для регистрации обработчика DELETE-запроса.
func (r *RouterHandle) DELETE(path string, handler http.HandlerFunc) *RouterHandle {
	return r.AddRoute(http.MethodDelete, path, handler)
}
