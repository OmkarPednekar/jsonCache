package jsonCache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheSetAndGet(t *testing.T) {
	cache := New(2)

	cache.Set("a", []byte("apple"), 1000)
	cache.Set("b", []byte("banana"), 1000)

	val := cache.Get("a")
	if string(val) != "apple" {
		t.Errorf("Expected 'apple', got '%s'", val)
	}

	val = cache.Get("b")
	if string(val) != "banana" {
		t.Errorf("Expected 'banana', got '%s'", val)
	}
}

func TestCacheExpiry(t *testing.T) {
	cache := New(2)

	cache.Set("a", []byte("apple"), 10) // 10ms TTL
	time.Sleep(20 * time.Millisecond)

	val := cache.Get("a")
	if val != nil {
		t.Errorf("Expected 'nil' after expiry, got '%s'", val)
	}
}

func TestCacheEviction(t *testing.T) {
	cache := New(2)

	cache.Set("a", []byte("apple"), 1000)
	cache.Set("b", []byte("banana"), 1000)
	cache.Set("c", []byte("cherry"), 1000) // Should evict "a"

	if cache.Get("a") != nil {
		t.Errorf("Expected 'a' to be evicted")
	}
	if string(cache.Get("b")) != "banana" {
		t.Errorf("Expected 'banana' to still be present")
	}
	if string(cache.Get("c")) != "cherry" {
		t.Errorf("Expected 'cherry' to still be present")
	}
}

func TestCacheUpdateExistingKey(t *testing.T) {
	cache := New(2)

	cache.Set("a", []byte("apple"), 1000)
	cache.Set("a", []byte("avocado"), 1000) // Update same key

	val := cache.Get("a")
	if cache.list.length != 1 {
		t.Errorf("Duplicate Nodes added with the same key ")
	}
	if string(val) != "avocado" {
		t.Errorf("Expected 'avocado', got '%s'", val)
	}
}

func TestCacheEvictionAndExpirationCombined(t *testing.T) {
	cache := New(2) // Set a small size for the cache to trigger eviction

	// Set some values with a TTL of 1 second
	cache.Set("key1", []byte("value1"), 1000)
	cache.Set("key2", []byte("value2"), 4000)

	// Wait for 2 seconds to make sure key1 expires
	time.Sleep(2 * time.Second)

	// Set a new value to cause eviction
	cache.Set("key3", []byte("value3"), 1000)

	// Test that expired key1 is evicted and key2 is still in the cache
	t.Run("Cache expiration and eviction combined", func(t *testing.T) {
		value := cache.Get("key1")
		assert.Nil(t, value, "Cache should return nil for expired key1")

		value = cache.Get("key2")
		assert.NotNil(t, value, "Cache should return value for key2")

		// key3 should still be in cache
		value = cache.Get("key3")
		assert.NotNil(t, value, "Cache should return value for key3")
	})
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := New(100)
	var wg sync.WaitGroup
	numOps := 1000

	t.Run("Concurrent set and get", func(t *testing.T) {
		// Writers
		for i := 0; i < numOps; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				key := fmt.Sprintf("key%d", i)
				cache.Set(key, []byte(fmt.Sprintf("val%d", i)), 5000)
			}(i)
		}

		// Readers
		for i := 0; i < numOps; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				key := fmt.Sprintf("key%d", i)
				_ = cache.Get(key) // We're just checking for races
			}(i)
		}

		wg.Wait()
	})
}
