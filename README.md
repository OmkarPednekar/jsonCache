# jsonCache

A simple, in-memory key-value cache library for JSON data with support for TTL (Time-To-Live). This library provides fast caching for JSON responses, making it ideal for web applications where caching JSON data improves performance.

## Features

- **In-memory caching**: Fast, efficient storage of JSON data in memory.
- **TTL support**: Set an expiration time (TTL) for cache entries.
- **Automatic expired entry deletion**: Removes expired cache entries.
- **Thread-safety**: Designed to be used in concurrent environments.
- **LRU Eviction**: Least Recently Used eviction policy to ensure cache size remains within limits.

## Installation

You can easily add this package to your Go project by running:

```bash
go get github.com/OmkarPednekar/jsonCache
```

## Usage

Here’s a quick example of how to use `jsonCache`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/OmkarPednekar/jsonCache"
	"time"
)

func main() {
	// Initialize cache
	cache := jsonCache.New(2) // Set cache size to 2

	// Set cache with a TTL of 3 seconds
	cache.Set("key1", []byte(`{"message": "Hello, World!"}`), 3000)
	cache.Set("key2", []byte(`{"message": "Another value"}`), 3000)

	// Get from cache
	data := cache.Get("key1")
	if data != nil {
		var result map[string]string
		err := json.Unmarshal(data, &result)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		} else {
			// Print parsed JSON data
			fmt.Println("Cache hit for key1:", result)
		}
	} else {
		fmt.Println("Cache miss for key1")
	}

	// Set another value to trigger eviction (LRU)
	cache.Set("key3", []byte(`{"message": "Evicted value"}`), 3000)

	// Try to get evicted key
	data = cache.Get("key2") // This should be evicted due to LRU policy
	if data != nil {
		var result map[string]string
		err := json.Unmarshal(data, &result)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		} else {
			// Print parsed JSON data
			fmt.Println("Cache hit for key2:", result)
		}
	} else {
		fmt.Println("Cache miss for key2 (evicted)")
	}
}
```

### Output:

```go
Cache hit for key1: {"message": "Hello, World!"}
Cache miss for key2
Cache miss for key2 (evicted)
```
