// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

package router

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// clientEntry содержит статистику запросов для одного IP-адреса.
type clientEntry struct {
	count     int
	windowStart time.Time
	lastSeen  time.Time
}

// CustomRateLimiter отвечает за ограничение количества HTTP-запросов (Rate Limiting).
type CustomRateLimiter struct {
	mu          sync.Mutex
	clients     map[string]*clientEntry
	maxRequests int
	window      time.Duration
	cleanupTTL  time.Duration
}

// NewCustomRateLimiter создает и инициализирует новый Rate Limiter.
// Запускает фоновую горутину для очистки неактивных IP из памяти.
func NewCustomRateLimiter(maxRequests int, window time.Duration, cleanupTTL time.Duration) *CustomRateLimiter {
	limiter := &CustomRateLimiter{
		clients:     make(map[string]*clientEntry),
		maxRequests: maxRequests,
		window:      window,
		cleanupTTL:  cleanupTTL,
	}

	go limiter.startCleanupWorker()

	return limiter
}

// startCleanupWorker удаляет устаревшие записи клиентов для предотвращения утечек памяти.
func (l *CustomRateLimiter) startCleanupWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, entry := range l.clients {
			if now.Sub(entry.lastSeen) > l.cleanupTTL {
				delete(l.clients, ip)
			}
		}
		l.mu.Unlock()
	}
}

// Allow проверяет, превысил ли IP-адрес лимит запросов за текущий временной интервал.
func (l *CustomRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry, exists := l.clients[ip]

	if !exists {
		l.clients[ip] = &clientEntry{
			count:       1,
			windowStart: now,
			lastSeen:    now,
		}
		return true
	}

	entry.lastSeen = now

	// Если временное окно истекло — сбрасываем счетчик и обновляем начало окна
	if now.Sub(entry.windowStart) >= l.window {
		entry.count = 1
		entry.windowStart = now
		return true
	}

	entry.count++
	return entry.count <= l.maxRequests
}

// RateLimitMiddleware возвращает HTTP-middleware для ограничения частоты запросов.
func RateLimitMiddleware(limiter *CustomRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if limiter == nil {
				next.ServeHTTP(w, req)
				return
			}

			ip := getClientIP(req)

			if !limiter.Allow(ip) {
				http.Error(w, "429 Too Many Requests. Try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

// getClientIP извлекает реальный IP-адрес клиента из заголовков прокси (Cloudflare, X-Forwarded-For) или RemoteAddr.
func getClientIP(req *http.Request) string {
	if cfIP := req.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return strings.TrimSpace(cfIP)
	}

	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}
	return host
}
