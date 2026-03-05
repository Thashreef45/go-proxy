package ratelimiter

import (
	"net/http"
)

func RateLimiter(l *LRUCache, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr

		if !l.AllowRequest(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests. Please try again later."))
			return
		}
		next.ServeHTTP(w, r)

	})
}
