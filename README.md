# jsonCache

A simple, in-memory key-value cache library for JSON data with support for TTL (Time-To-Live). This library provides fast caching for JSON responses, making it ideal for web applications where caching JSON data improves performance.

## Features

- **In-memory caching**: Fast, efficient storage of JSON data in memory.
- **TTL support**: Set an expiration time (TTL) for cache entries.
- **Automatic expired entry deletion**: Removes expired cache entries.
- **Thread-safety**: Designed to be used in concurrent environments.
- **Eviction policies** (Coming soon): Plan to implement LRU, LFU, and FIFO eviction strategies.

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
	cache := jsonCache.New()

	// Set cache with a TTL of 3 seconds
	cache.Set("key", []byte(`{"message": "Hello, World!"}`), 3000)

	// Get from cache
	data := cache.Get("key")
	if data != nil {
		var result map[string]string
		err := json.Unmarshal(data, &result)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		} else {
			// Print parsed JSON data
			fmt.Println("Cache hit:", result)
		}
	} else {
		fmt.Println("Cache miss")
	}

	// Wait for cache to expire
	time.Sleep(4 * time.Second)

	// Try again after TTL expiry
	data = cache.Get("key")
	if data != nil {
		var result map[string]string
		err := json.Unmarshal(data, &result)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		} else {
			// Print parsed JSON data
			fmt.Println("Cache hit after expiry:", result)
		}
	} else {
		fmt.Println("Cache miss after expiry")
	}
}
```

### Output:

```go
Cache hit: {"message": "Hello, World!"}
Cache miss after expiry
```
