package cache

import "net/http"

func CacheHandler(c *cache, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check value exist in cache
		// if exist , return value
		// else call next

		data, valueExist := c.GetCachedData("")
		if valueExist {
			data = data
			// response
			return
		}
		next.ServeHTTP(w, r)
	})
}
