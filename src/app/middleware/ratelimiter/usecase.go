package ratelimiter

import (
	"time"
)

// LRU Cache with token bucket rate limiter
func NewLRUCache(capacity, maxTokens, refillRate int) *LRUCache {

	/*capacity - maximum number of buckets (IPs) to store*/
	/*maxTokens - maximum tokens in each bucket*/
	/*refillRate - number of tokens to add per second*/

	head := Bucket{}
	tail := Bucket{}
	head.Next = &tail
	tail.Prev = &head

	return &LRUCache{
		Cache:      make(map[string]*Bucket, capacity),
		Head:       &head,
		Tail:       &tail,
		Capacity:   capacity,
		MaxTokens:  maxTokens,
		RefillRate: refillRate,
	}
}

func (this *LRUCache) AllowRequest(ip string) bool {

	this.mu.Lock()
	defer this.mu.Unlock()

	node, exist := this.Cache[ip]
	if !exist {
		this.addBucket(ip)
	} else {
		// refill the bucket with tokens based on time elapsed since last visit
		this.refillBucket(ip)

		// if no tokens left, return 429 Too Many Requests
		if node.Tokens <= 0 {
			return false
		}

		// move the bucket to head of the list
		this.detachBucket(node)
		this.attachToHead(node)
	}
	// decrement a token from bucket
	this.decrementToken(ip)
	return true
}

func (this *LRUCache) addBucket(ip string) {

	INITIAL_TOKEN_COUNT := 10 // default token count
	if this.MaxTokens < INITIAL_TOKEN_COUNT {
		INITIAL_TOKEN_COUNT = this.MaxTokens
	}

	if node, isExist := this.Cache[ip]; isExist {
		// node value should be updated
		// node have to move to the next of head
		node.LastVisit = time.Now()
		this.detachBucket(node)
		this.attachToHead(node)
		return
	}

	if this.Capacity <= len(this.Cache) {
		// remove from tail
		lastBucket := this.Tail.Prev
		this.detachBucket(lastBucket)
		delete(this.Cache, lastBucket.IP)
	}

	b := &Bucket{
		IP:        ip,
		Tokens:    INITIAL_TOKEN_COUNT,
		LastVisit: time.Now(),
	}

	this.attachToHead(b)
	this.Cache[ip] = b
}

func (this *LRUCache) attachToHead(node *Bucket) {
	head := this.Head
	headNext := head.Next

	node.Prev = head
	node.Next = headNext
	head.Next = node
	headNext.Prev = node
}

func (this *LRUCache) detachBucket(node *Bucket) {
	prev := node.Prev
	next := node.Next
	prev.Next = next
	next.Prev = prev
}

func (this *LRUCache) decrementToken(ip string) {
	if this.Cache[ip].Tokens > 0 {
		this.Cache[ip].Tokens--
	}
}

func (this *LRUCache) refillBucket(ip string) {
	currentTime := time.Now()
	// add tokens to bucket , current time - bucket.LastVisit time in seconds * 5
	prevTime := this.Cache[ip].LastVisit
	tokensToBeAdded := int(currentTime.Sub(prevTime).Seconds() * float64(this.RefillRate))
	if tokensToBeAdded > 0 {
		this.Cache[ip].Tokens += tokensToBeAdded
		if this.Cache[ip].Tokens > this.MaxTokens {
			this.Cache[ip].Tokens = this.MaxTokens
		}
		// Update LastVisit only when tokens are actually added
		this.Cache[ip].LastVisit = currentTime
	}
}
