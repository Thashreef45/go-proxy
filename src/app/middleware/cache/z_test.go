package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Thashreef45/proxy-server/src/internal/model"
)

func TestCacheWriteAndGet(t *testing.T) {
	c := NewCache(10)

	key := "/api/products"
	req := Request{
		Header:    http.Header{"Content-Type": []string{"application/json"}},
		Body:      []byte(`{"name":"product"}`),
		ExpiredIn: 5 * time.Minute,
	}

	c.write(key, req)

	node, exist := c.GetCachedData(key)
	if !exist {
		t.Fatalf("expected cache entry to exist")
	}

	if string(node.Value.Body) != `{"name":"product"}` {
		t.Errorf("unexpected cached body")
	}
}

func TestCacheExpiration(t *testing.T) {
	c := NewCache(10)

	key := "/api/products"
	req := Request{
		Body:      []byte("test-body"),
		ExpiredIn: 1 * time.Second,
	}

	c.write(key, req)

	time.Sleep(2 * time.Second)

	_, exist := c.GetCachedData(key)
	if exist {
		t.Errorf("expected cache to expire")
	}
}

func TestCacheLRUEviction(t *testing.T) {
	c := NewCache(2)

	c.write("a", Request{Body: []byte("A"), ExpiredIn: time.Minute})
	c.write("b", Request{Body: []byte("B"), ExpiredIn: time.Minute})
	c.write("c", Request{Body: []byte("C"), ExpiredIn: time.Minute})

	if _, ok := c.cache["a"]; ok {
		t.Errorf("expected 'a' to be evicted")
	}

	if _, ok := c.cache["c"]; !ok {
		t.Errorf("expected 'c' to exist")
	}
}

func TestCacheAccessUpdatesLRU(t *testing.T) {
	c := NewCache(2)

	c.write("a", Request{Body: []byte("A"), ExpiredIn: time.Minute})
	c.write("b", Request{Body: []byte("B"), ExpiredIn: time.Minute})

	// access "a"
	c.GetCachedData("a")

	c.write("c", Request{Body: []byte("C"), ExpiredIn: time.Minute})

	if _, ok := c.cache["b"]; ok {
		t.Errorf("expected 'b' to be evicted")
	}
}

func TestCacheMiddleware(t *testing.T) {
	c := NewCache(10)

	cfg := model.CacheConfig{
		Routes: []model.CacheRoute{
			{
				Path: "/api/products",
				Ttl:  60,
			},
		},
	}

	handler := CacheHandler(c, cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend response"))
	}))

	req := httptest.NewRequest("GET", "/api/products", nil)
	rec := httptest.NewRecorder()

	// first request (cache miss)
	handler.ServeHTTP(rec, req)

	if rec.Body.String() != "backend response" {
		t.Errorf("unexpected response")
	}

	// second request (cache hit)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)

	if rec2.Body.String() != "backend response" {
		t.Errorf("cache did not return expected response")
	}
}

func TestCacheMiddlewareSkipRoute(t *testing.T) {
	c := NewCache(10)

	cfg := model.CacheConfig{
		Routes: []model.CacheRoute{
			{Path: "/api/products", Ttl: 60},
		},
	}

	count := 0

	handler := CacheHandler(c, cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.Write([]byte("backend"))
	}))

	req := httptest.NewRequest("GET", "/api/users", nil)

	handler.ServeHTTP(httptest.NewRecorder(), req)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if count != 2 {
		t.Errorf("route should not be cached")
	}
}

func TestCacheMiddlewareCachesRoute(t *testing.T) {
	c := NewCache(10)

	cfg := model.CacheConfig{
		Routes: []model.CacheRoute{
			{Path: "/api/products", Ttl: 60},
		},
	}

	count := 0

	handler := CacheHandler(c, cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.Write([]byte("backend"))
	}))

	req := httptest.NewRequest("GET", "/api/products", nil)

	handler.ServeHTTP(httptest.NewRecorder(), req)
	handler.ServeHTTP(httptest.NewRecorder(), req)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if count != 1 {
		t.Errorf("expected backend to be called once, got %d", count)
	}
}

func TestCacheQueryKey(t *testing.T) {
	c := NewCache(10)

	cfg := model.CacheConfig{
		Routes: []model.CacheRoute{
			{Path: "/api/products", Ttl: 60},
		},
	}

	count := 0

	handler := CacheHandler(c, cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.Write([]byte("data"))
	}))

	req1 := httptest.NewRequest("GET", "/api/products?page=1", nil)
	req2 := httptest.NewRequest("GET", "/api/products?page=2", nil)

	handler.ServeHTTP(httptest.NewRecorder(), req1)
	handler.ServeHTTP(httptest.NewRecorder(), req2)

	if count != 2 {
		t.Errorf("queries should create different cache entries")
	}
}
