package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// --- UNIT TESTS ---

func TestDecrementToken(t *testing.T) {
	cache := NewLRUCache(10, 5, 1)
	ip := "1.2.3.4"

	cache.addBucket(ip)
	if cache.Cache[ip].Tokens != 5 {
		t.Errorf("expected 5 tokens, got %d", cache.Cache[ip].Tokens)
	}

	cache.decrementToken(ip)
	if cache.Cache[ip].Tokens != 4 {
		t.Errorf("expected 4 tokens, got %d", cache.Cache[ip].Tokens)
	}
}

func TestRefillBucket(t *testing.T) {
	cache := NewLRUCache(10, 5, 2) // refillRate = 2 tokens/sec
	ip := "1.2.3.5"

	cache.addBucket(ip)
	cache.Cache[ip].Tokens = 0

	// simulate 2 seconds elapsed
	cache.Cache[ip].LastVisit = time.Now().Add(-2 * time.Second)
	cache.refillBucket(ip)

	expected := 4 // 2 * 2 tokens
	if cache.Cache[ip].Tokens != expected {
		t.Errorf("expected %d tokens, got %d", expected, cache.Cache[ip].Tokens)
	}
}

func TestLRUEviction(t *testing.T) {
	cache := NewLRUCache(2, 5, 1)
	cache.addBucket("ip1")
	cache.addBucket("ip2")
	cache.addBucket("ip3") // should evict ip1

	if _, exists := cache.Cache["ip1"]; exists {
		t.Errorf("expected ip1 to be evicted")
	}
	if _, exists := cache.Cache["ip3"]; !exists {
		t.Errorf("expected ip3 to exist")
	}
}

func TestAllowRequest(t *testing.T) {
	cache := NewLRUCache(2, 2, 1)
	ip := "1.2.3.6"

	allowed := cache.AllowRequest(ip) // first token
	if !allowed {
		t.Errorf("first request should be allowed")
	}
	allowed = cache.AllowRequest(ip) // second token
	if !allowed {
		t.Errorf("second request should be allowed")
	}
	allowed = cache.AllowRequest(ip) // third token, should fail
	if allowed {
		t.Errorf("third request should be denied")
	}
}

// --- INTEGRATION TEST: HTTP MIDDLEWARE ---

func TestRateLimiterMiddleware(t *testing.T) {
	cache := NewLRUCache(10, 2, 1)
	handler := RateLimiter(cache, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	// First request should succeed
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	// Second request should succeed
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	// Third request should get 429
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too Many Requests, got %d", rec.Code)
	}
}

// --- CONCURRENCY TEST ---

func TestConcurrency(t *testing.T) {
	cache := NewLRUCache(10, 100, 1)
	ip := "1.2.3.7"

	var wg sync.WaitGroup
	for i := 0; i < 50000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.AllowRequest(ip)
		}()
	}
	wg.Wait()
}

// --- OPTIONAL: Test that tokens do not go negative ---

func TestTokenNeverNegative(t *testing.T) {
	cache := NewLRUCache(2, 1, 1)
	ip := "1.2.3.8"

	cache.addBucket(ip)
	cache.decrementToken(ip) // Tokens = 0
	cache.decrementToken(ip) // Should stay at 0, not negative

	if cache.Cache[ip].Tokens < 0 {
		t.Errorf("tokens went negative: %d", cache.Cache[ip].Tokens)
	}
}
