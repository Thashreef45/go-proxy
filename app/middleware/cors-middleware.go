package middleware

import (
	"net/http"

	"github.com/Thashreef45/proxy-server/internal/model"
)

func CORSMiddleware(corsCfg model.CorsConfig, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")

		// check origin
		for _, o := range corsCfg.AllowedOrigins {
			if o == "*" || o == origin {
				w.Header().Set("Access-Control-Allow-Origin", o)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", joinOrDefault(corsCfg.AllowedMethods, "GET,POST,PUT,DELETE,OPTIONS"))
		w.Header().Set("Access-Control-Allow-Headers", joinOrDefault(corsCfg.AllowedHeaders, "Content-Type,Authorization"))

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func joinOrDefault(arr []string, defaultVal string) string {
	if len(arr) == 0 {
		return defaultVal
	}
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}
