package ratelimiter

import (
	"net"
	"net/http"
)

func RateLimiter(l *LRUCache, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		if !l.AllowRequest(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests. Please try again later."))
			return
		}
		next.ServeHTTP(w, r)

	})
}
