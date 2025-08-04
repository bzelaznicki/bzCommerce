package main

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     time.Duration
	burst    int
}

type visitor struct {
	lastSeen time.Time
	tokens   int
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}
	
	// Clean up old visitors every 30 minutes
	go rl.cleanupVisitors()
	
	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	v, exists := rl.visitors[ip]
	now := time.Now()
	
	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: now,
			tokens:   rl.burst - 1,
		}
		return true
	}
	
	// Add tokens based on time elapsed
	elapsed := now.Sub(v.lastSeen)
	tokensToAdd := int(elapsed / rl.rate)
	
	if tokensToAdd > 0 {
		v.tokens += tokensToAdd
		if v.tokens > rl.burst {
			v.tokens = rl.burst
		}
		v.lastSeen = now
	}
	
	if v.tokens > 0 {
		v.tokens--
		return true
	}
	
	return false
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(30 * time.Minute)
		rl.mu.Lock()
		
		cutoff := time.Now().Add(-time.Hour)
		for ip, v := range rl.visitors {
			if v.lastSeen.Before(cutoff) {
				delete(rl.visitors, ip)
			}
		}
		
		rl.mu.Unlock()
	}
}

func (cfg *apiConfig) withRateLimit(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP, considering X-Forwarded-For header
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.Header.Get("X-Real-IP")
			}
			if ip == "" {
				ip = r.RemoteAddr
			}
			
			if !rl.Allow(ip) {
				respondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}