// test for the Cache implementation in pokecache.go
package poke

import (
	"testing"
	"time"
)

func TestCacheAddAndGet(t *testing.T) {
	cache := NewCache(2 * time.Second)

	key := "testKey"
	data := []byte("testData")

	cache.Add(key, data)

	retrievedData, found := cache.Get(key)
	if !found {
		t.Fatalf("Expected to find key %s in cache", key)
	}

	if string(retrievedData) != string(data) {
		t.Fatalf("Expected data %s, got %s", data, retrievedData)
	}
}
