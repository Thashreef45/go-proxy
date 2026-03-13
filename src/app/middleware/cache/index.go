package cache

import (
	"net/http"
	"time"

	"github.com/Thashreef45/proxy-server/src/internal/model"
)

func CacheHandler(c *cache, config model.CacheConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// cache is only for GET requests
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		path := r.URL.Path
		url := r.URL.String()

		pathExist := false
		var ttl time.Duration

		for _, cfg := range config.Routes {
			if cfg.Path == path {
				pathExist = true
				ttl = time.Duration(cfg.Ttl) * time.Second
				break
			}
		}
		if !pathExist {
			next.ServeHTTP(w, r)
			return
		}

		// if cache hits (data found on cache)
		data, valueExist := c.GetCachedData(url)
		if valueExist {
			for key, val := range data.Value.Header {
				for _, v := range val {
					w.Header().Add(key, v)
				}
			}
			w.Write(data.Value.Body)
			return
		}

		// if cache miss (data not exist in cache)
		resRecorder := &responseRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(resRecorder, r)

		// add new response to cache
		req := Request{
			Header:    resRecorder.Header(),
			Body:      resRecorder.body,
			ExpiredIn: ttl,
		}
		c.write(url, req)
	})
}
