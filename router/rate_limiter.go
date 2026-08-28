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

type clientEntry struct {
	count       int
	prevCount   int
	windowStart time.Time
	lastSeen    time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	clients     map[string]*clientEntry
	MaxRequests int
	Window      time.Duration
	CleanupTTL  time.Duration
	MaxClients  int
	stopChan    chan struct{}
}

// NewLimiter создает и инициализирует новый Rate Limiter.
func NewLimiter(maxRequests int, window time.Duration, cleanupTTL time.Duration, maxClients int) *RateLimiter {
	limiter := &RateLimiter{
		clients:     make(map[string]*clientEntry),
		MaxRequests: maxRequests,
		Window:      window,
		CleanupTTL:  cleanupTTL,
		MaxClients:  maxClients,
		stopChan:    make(chan struct{}),
	}

	go limiter.startCleanupWorker()

	return limiter
}

// Stop останавливает фоновый воркер очистки.
func (r *RateLimiter) Stop() {
	if r.stopChan != nil {
		select {
		case <-r.stopChan:
		default:
			close(r.stopChan)
		}
	}
}

func (r *RateLimiter) startCleanupWorker() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			now := time.Now()
			for ip, entry := range r.clients {
				if now.Sub(entry.lastSeen) > r.CleanupTTL {
					delete(r.clients, ip)
				}
			}
			r.mu.Unlock()
		case <-r.stopChan:
			return
		}
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	entry, exists := r.clients[ip]

	if !exists {
		if r.MaxClients > 0 && len(r.clients) >= r.MaxClients {
			return false
		}

		r.clients[ip] = &clientEntry{
			count:       1,
			windowStart: now,
			lastSeen:    now,
		}
		return true
	}

	entry.lastSeen = now

	elapsed := now.Sub(entry.windowStart)
	if elapsed >= r.Window*2 {
		entry.prevCount = 0
		entry.count = 1
		entry.windowStart = now
		return true
	} else if elapsed >= r.Window {
		entry.prevCount = entry.count
		entry.count = 1
		entry.windowStart = entry.windowStart.Add(r.Window)
		return true
	}

	timeIntoCurrentWindow := elapsed
	weight := float64(r.Window-timeIntoCurrentWindow) / float64(r.Window)
	estimatedRequests := float64(entry.prevCount)*weight + float64(entry.count)

	if estimatedRequests >= float64(r.MaxRequests) {
		return false
	}

	entry.count++
	return true
}

func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if limiter == nil {
				next.ServeHTTP(w, req)
				return
			}

			ip := getClientIP(req)

			if !limiter.Allow(ip) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "429 Too Many Requests. Try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

func getClientIP(req *http.Request) string {
	if cfIP := req.Header.Get("CF-Connecting-IP"); cfIP != "" {
		if net.ParseIP(strings.TrimSpace(cfIP)) != nil {
			return strings.TrimSpace(cfIP)
		}
	}

	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		if net.ParseIP(req.RemoteAddr) != nil {
			return req.RemoteAddr
		}
		return "0.0.0.0"
	}

	if net.ParseIP(host) != nil {
		return host
	}
	return "0.0.0.0"
}
