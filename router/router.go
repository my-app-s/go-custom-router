// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

// Package router provides a simple custom HTTP router
// for handling requests in Go applications.
package router

import (
	"log"
	"net/http"
)

// RouterHandle хранит карту маршрутов (URL-путь -> HTTP-метод -> HandlerFunc)
// и конфигурацию CORS. Реализует интерфейс http.Handler.
type RouterHandle struct {
	// Routes: карта [Path][Method]http.HandlerFunc
	Routes map[string]map[string]http.HandlerFunc `json:"routes"`

	// CORS содержит настройки междоменного доступа (CORS).
	// Если равна nil, цепочка middleware автоматически использует CorsMiddlewareOpen.
	CORS *CORS `json:"cors,omitempty"`
}

// NewRouterHandle инициализирует новый экземпляр RouterHandle с пустой картой маршрутов.
func NewRouterHandle() *RouterHandle {
	return &RouterHandle{
		Routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// ServeHTTP выполняет диспетчеризацию входящих HTTP-запросов к соответствующим обработчикам.
// Если путь не найден — возвращает статус 404 Not Found.
// Если путь найден, но метод не зарегистрирован — возвращает статус 405 Method Not Allowed.
func (r *RouterHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("[%s] %s", req.Method, req.URL.Path)

	// 1. Проверяем, существует ли сам путь
	methods, pathExists := r.Routes[req.URL.Path]
	if !pathExists {
		http.NotFound(w, req)
		return
	}

	// 2. Проверяем, зарегистрирован ли HTTP-метод для этого пути
	handler, methodExists := methods[req.Method]
	if !methodExists {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 3. Вызываем целевой обработчик
	handler(w, req)
}
