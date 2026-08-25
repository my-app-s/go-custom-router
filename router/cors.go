package router

import (
	"log"
	"net/http"
)

// CORS хранит конфигурацию для обработки междоменных запросов (Cross-Origin Resource Sharing).
type CORS struct {
	// AllowedOrigins содержит множество разрешенных доменов.
	AllowedOrigins map[string]struct{} `json:"allowedOrigins"`
	// AllowedMethods содержит строку с разрешенными HTTP-методами через запятую.
	AllowedMethods string              `json:"allowedMethods"`
	// AllowedHeaders содержит строку с разрешенными HTTP-заголовками через запятую.
	AllowedHeaders string              `json:"allowedHeaders"`
}

// NewCORS создает и инициализирует новый экземпляр структуры CORS
// с готовой к работе картой разрешенных доменов.
func NewCORS() *CORS {
	return &CORS{
		AllowedOrigins: make(map[string]struct{}),
	}
}

// AddOrigin регистрирует разрешенный домен (origin) с соответствующими HTTP-методами и заголовками.
// Возвращает true в случае успешного добавления или false, если один из параметров является пустой строкой.
func (c *CORS) AddOrigin(origin, methods, headers string) bool {
	if origin == "" {
		log.Println("Error input: Empty origin")
		return false
	}

	if methods == "" {
		log.Println("Error input: Empty methods")
		return false
	}

	if headers == "" {
		log.Println("Error input: Empty headers")
		return false
	}

	c.AllowedOrigins[origin] = struct{}{}
	c.AllowedMethods = methods
	c.AllowedHeaders = headers

	return true
}

// IsAllowed проверяет, находится ли указанный origin в списке разрешенных доменов.
// Безопасно возвращает false, если origin не разрешен, передан пустой строку,
// либо если объект CORS или карта AllowedOrigins не инициализированы.
func (c *CORS) IsAllowed(origin string) bool {
	if c == nil || c.AllowedOrigins == nil {
        return false
    }

	_, ok := c.AllowedOrigins[origin]
	return ok
}

// CorsMiddleware возвращает HTTP-middleware, обрабатывающий междоменные запросы (CORS).
// Если Origin входящего запроса разрешен (IsAllowed), метод устанавливает соответствующие
// заголовки Access-Control-Allow-*. Для preflight-запросов (HTTP OPTIONS) автоматически
// отправляется статус 200 OK без передачи управления следующему обработчику.
func (c *CORS) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if c.IsAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", c.AllowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", c.AllowedHeaders)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CorsMiddlewareOpen оборачивает http.Handler и предоставляет открытый доступ для любых источников (Origin: *).
// Метод автоматически устанавливает базовые заголовки CORS (GET, POST, OPTIONS), отсекает
// preflight-запросы (HTTP OPTIONS) со статусом 200 OK и подходит для публичных API или режима разработки.
func CorsMiddlewareOpen(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
