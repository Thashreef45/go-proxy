package cache

import "time"

func NewCache(capacity int) *cache {
	h := &Node{}
	t := &Node{}
	h.Next = t
	t.Prev = h

	return &cache{
		cache:    make(map[string]*Node),
		Head:     h,
		Tail:     t,
		Capacity: capacity,
	}
}

func (c *cache) GetCachedData(key string) (*Node, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exist := c.cache[key]
	if !exist {
		return nil, false
	}
	if c.isExpired(node) {
		c.removeNode(node)
		delete(c.cache, key)
		return nil, false
	}

	c.removeNode(node)
	c.addToHead(node)
	return node, true
}

func (c *cache) isExpired(node *Node) bool {
	return time.Since(node.Value.CachedAt) > node.Value.ExpiredIn
	// return false
}

func (c *cache) write(key string, value Request) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exist := c.cache[key]; exist {
		c.removeNode(node)
		c.addToHead(node)
	} else {
		if len(c.cache) == c.Capacity {
			lastNode := c.Tail.Prev
			c.removeNode(lastNode)
			delete(c.cache, lastNode.Key)
		}

		rq := Request{
			Header:    value.Header,
			Body:      value.Body,
			CachedAt:  time.Now(),
			ExpiredIn: value.ExpiredIn,
		}
		newNode := &Node{
			Key:   key,
			Value: rq,
		}

		c.addToHead(newNode)
		c.cache[key] = newNode
	}
}

func (c *cache) addToHead(node *Node) {
	head := c.Head
	headNext := head.Next
	head.Next = node
	node.Prev = head
	node.Next = headNext
	headNext.Prev = node
}

func (c *cache) removeNode(node *Node) {
	prevNode := node.Prev
	nextNode := node.Next

	prevNode.Next = nextNode
	nextNode.Prev = prevNode
}
