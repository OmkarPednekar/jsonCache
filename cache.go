package jsonCache

import (
	"time"
)

type Node struct {
	key   string
	value []byte
	prev  *Node
	next  *Node
}

type List struct {
	head   *Node
	tail   *Node
	length int
}

type Cache struct {
	store map[string]CacheElement
	list  *List
	cap   int
}

type CacheElement struct {
	node   *Node
	expiry time.Time
}

func New(size int) *Cache {
	//initialize a lru list
	dummyHead, dummyTail := _initDummyNodes()
	list := &List{head: dummyHead, tail: dummyTail, length: 0}
	//initialize a new cache engine
	cache := &Cache{
		store: make(map[string]CacheElement),
		list:  list,
		cap: func() int {
			if size == 0 {
				return 500
			}
			return size
		}(),
	}

	return cache
}

func (c *Cache) Set(key string, value []byte, ttl int) {
	expiration := time.Now().Add(time.Duration(ttl) * time.Millisecond)
	_node := _newNode(key, value)
	var toStore = CacheElement{
		node:   _node,
		expiry: expiration,
	}
	//cache full
	if c.list.length >= c.cap {
		c._evict()
	}
	c._addNode(_node)
	c.store[key] = toStore
}

func (c *Cache) Get(key string) []byte {
	entry, found := c.store[key]
	if !found {
		return nil
	}
	if time.Now().After(entry.expiry) {
		delete(c.store, key)
		c._removeNode(entry.node)
		return nil
	}
	c._moveToHead(entry.node)
	return entry.node.value
}

func (c *Cache) _moveToHead(node *Node) {
	// If node is already at the head, no need to move it.
	if node == c.list.head.next {
		return
	}

	// Unlink the node from its current position
	node.prev.next = node.next
	node.next.prev = node.prev

	// Move node to the head
	head := c.list.head
	tmpNext := head.next
	head.next = node
	node.next = tmpNext
	node.prev = head
	if tmpNext != nil {
		tmpNext.prev = node
	}
}

func (c *Cache) _removeNode(toRemove *Node) {
	// Unlink the node
	toRemove.prev.next = toRemove.next
	if toRemove.next != nil {
		toRemove.next.prev = toRemove.prev
	}

	// Decrease the list length
	if c.list.length > 0 {
		c.list.length--
	}
}

func (c *Cache) _evict() {
	// Ensure there's something to evict
	if c.list.length == 0 {
		return
	}

	// Evict the tail node (least recently used)
	tail := c.list.tail
	tmpPrev := tail.prev

	// Unlink the tail node
	tmpPrev.prev.next = tmpPrev.next
	tail.prev = tmpPrev.prev

	// Decrease the list length
	if c.list.length > 0 {
		c.list.length--
	}
}

func (c *Cache) _addNode(node *Node) {
	// Add node at the head
	head := c.list.head
	tmpNext := head.next
	head.next = node
	node.next = tmpNext
	node.prev = head
	if tmpNext != nil {
		tmpNext.prev = node
	}
	c.list.length++
}

// helper functions
func _initDummyNodes() (*Node, *Node) {
	dummyHead := &Node{}
	dummyTail := &Node{}
	dummyHead.next = dummyTail
	dummyTail.prev = dummyHead
	return dummyHead, dummyTail
}

func _newNode(key string, value []byte) *Node {
	node := &Node{
		key:   key,
		value: value,
		prev:  nil,
		next:  nil,
	}
	return node
}
