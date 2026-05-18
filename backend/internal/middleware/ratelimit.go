package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter implementa rate limiting por IP
type RateLimiter struct {
	visitors          map[string]*rate.Limiter
	mu                sync.RWMutex
	maxVisitors       int
	requestsPerSecond float64
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(requestsPerSecond float64, maxVisitors int) *RateLimiter {
	rl := &RateLimiter{
		visitors:          make(map[string]*rate.Limiter),
		maxVisitors:       maxVisitors,
		requestsPerSecond: requestsPerSecond,
	}

	// Limpiar visitors inactivos cada 10 minutos
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			rl.mu.Lock()
			for ip, limiter := range rl.visitors {
				if !limiter.Allow() {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

// getVisitor obtiene o crea un limiter para una IP
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		// Limitar número máximo de visitantes
		if len(rl.visitors) >= rl.maxVisitors {
			// Usar un limiter muy restrictivo
			return rate.NewLimiter(rate.Limit(0), 0)
		}

		limiter := rate.NewLimiter(rate.Limit(rl.requestsPerSecond), 1)
		rl.visitors[ip] = limiter
		return limiter
	}

	return v
}

// RateLimit middleware para rate limiting
func RateLimit(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener IP del cliente
			ip := getClientIP(r)

			limiter := rl.getVisitor(ip)
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP obtiene la IP real del cliente
func getClientIP(r *http.Request) string {
	// Verificar X-Forwarded-For (detrás de proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Verificar X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Usar RemoteAddr como último recurso
	return r.RemoteAddr
}
