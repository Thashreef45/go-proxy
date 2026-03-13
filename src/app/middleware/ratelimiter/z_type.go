package ratelimiter

import (
	"sync"
	"time"
)

type Bucket struct {
	IP        string
	Tokens    int
	LastVisit time.Time
	Next      *Bucket
	Prev      *Bucket
}

type LRUCache struct {
	mu         sync.Mutex
	Cache      map[string]*Bucket
	Head       *Bucket
	Tail       *Bucket
	Capacity   int
	RefillRate int
	MaxTokens  int
}
