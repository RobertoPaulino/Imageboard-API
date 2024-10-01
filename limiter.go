package main

import (
	"net/http"
	"sync"
	"time"
)

var (
	rateLimitWindow = 1 * time.Minute
	maxRequests     = 5
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
}

type Visitor struct {
	lastSeen time.Time
	requests int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*Visitor),
	}
}

func (r *RateLimiter) Limit(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	visitor, exists := r.visitors[ip]
	if !exists || time.Since(visitor.lastSeen) > rateLimitWindow {
		r.visitors[ip] = &Visitor{lastSeen: time.Now(), requests: 1}
		return false
	}

	visitor.lastSeen = time.Now()
	if visitor.requests >= maxRequests {
		return true
	}

	visitor.requests++
	return false
}

func rateLimitMiddleware(handler http.HandlerFunc, limiter *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if limiter.Limit(ip) {
			http.Error(w, "429 Too Many Requests - Please slow down!", http.StatusTooManyRequests)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
