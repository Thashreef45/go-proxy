package cache

import (
	"net/http"
	"sync"
	"time"
)

// type Request struct {
// 	Header    http.Header
// 	Body      []byte
// 	CachedAt  time.Time
// 	ExpiredIn time.Duration
// 	Key       string
// 	Prev      *Request
// 	Next      *Request
// }

type Request struct {
	Header    http.Header
	Body      []byte
	CachedAt  time.Time
	ExpiredIn time.Duration
}

type Node struct {
	Key   string
	Value Request
	Prev  *Node
	Next  *Node
}

type cache struct {
	mu       sync.Mutex
	cache    map[string]*Node
	Head     *Node
	Tail     *Node
	Capacity int
}
