// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"encoding/json"
	"net/http"
)

// SendJSON отправляет структурированный JSON-ответ с указанным HTTP-статусом.
func SendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SendError отправляет стандартную JSON-ошибку с сообщением.
func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, map[string]string{"error": message})
}

// MakeCustomHandler создает HandlerFunc, возвращающий ключ-значение (key: value) в формате JSON.
// Полезно для статических ответов, заглушек (mock) и тестов.
func MakeCustomHandler(key, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{key: value}
		SendJSON(w, http.StatusOK, response)
	}
}
