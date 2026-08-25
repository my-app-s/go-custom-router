// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// RouterHandle содержит конфигурацию и обработчики роутера.
type RouterHandle struct {
	// ... ваши существующие поля ...
	CORS *CORS // Поле для хранения конфигурации CORS
}

// RecoveryMiddleware оборачивает http.Handler и предотвращает падение сервера при возникновении паники (panic).
// При обнаружении паники функция логирует ошибку со стеком вызовов (stack trace)
// и автоматически возвращает клиенту ответ 500 Internal Server Error.
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

// LoggerMiddleware оборачивает http.Handler и логирует сведения о входящих HTTP-запросах.
// На этапе получения запроса [START] записывается HTTP-метод и URL-путь.
// По завершении обработки [DONE] дополнительно фиксируется общее время выполнения (Duration).
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("📥 [START] %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("📤 [DONE]  %s %s | Duration: %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// ContentTypeJSONMiddleware оборачивает http.Handler и автоматически устанавливает
// HTTP-заголовок "Content-Type: application/json; charset=utf-8" для всех исходящих ответов.
// Применяется в цепочках middleware для REST/JSON API, гарантируя единообразный формат ответа.
func ContentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// Handler возвращает полностью обернутый http.Handler с базовым набором middleware.
//
// Цепочка выполнения при входящем HTTP-запросе (снаружи внутрь):
//  1. CORS — фильтрует Origin или пропускает все запросы (если r.CORS == nil, используется CorsMiddlewareOpen).
//  2. LoggerMiddleware — логирует факт запроса и время выполнения.
//  3. RecoveryMiddleware — перехватывает паники для предотвращения падения сервера.
//  4. RouterHandle — передает запрос целевому обработчику.
func (r *RouterHandle) Handler() http.Handler {
	var handler http.Handler = r

	handler = RecoveryMiddleware(handler)
	handler = LoggerMiddleware(handler)

	if r.CORS != nil {
		handler = r.CORS.CorsMiddleware(handler)
	} else {
		handler = CorsMiddlewareOpen(handler)
	}

	return handler
}

// HandlerAPI подготавливает роутер для работы в режиме REST / JSON API.
//
// В отличие от Handler, метод автоматически добавляет ContentTypeJSONMiddleware,
// который устанавливает заголовок "Content-Type: application/json; charset=utf-8".
//
// Цепочка выполнения при входящем HTTP-запросе (снаружи внутрь):
//  1. CORS — фильтрация Origin и обработка OPTIONS preflight.
//  2. LoggerMiddleware — логирование входящего запроса.
//  3. RecoveryMiddleware — перехват паник (возвращает 500 Internal Server Error).
//  4. ContentTypeJSONMiddleware — установка заголовка ответа JSON.
//  5. RouterHandle — передает запрос целевому обработчику.
func (r *RouterHandle) HandlerAPI() http.Handler {
	var handler http.Handler = r

	handler = ContentTypeJSONMiddleware(handler)
	handler = RecoveryMiddleware(handler)
	handler = LoggerMiddleware(handler)

	if r.CORS != nil {
		handler = r.CORS.CorsMiddleware(handler)
	} else {
		handler = CorsMiddlewareOpen(handler)
	}

	return handler
}
