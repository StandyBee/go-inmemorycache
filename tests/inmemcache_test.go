package tests

import (
	"github.com/StandyBee/go-inmemorycache/inmemcache"
	"testing"
)

func TestInMemCache_SetAndGet(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	err := cache.Set("key1", 123)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, found := cache.Get("key1")
	if !found {
		t.Fatalf("Get failed: %v", err)
	}

	if val != 123 {
		t.Fatalf("Get failed: %v, expected 123", val)
	}
}

func TestInMemCache_Set_EmptyKey(t *testing.T) {
	cache := inmemcache.NewInMemCache()
	err := cache.Set("", 0)
	if err == nil {
		t.Fatalf("Set failed: %v", err)
	}
}

func TestInMemCache_Get_NonExistentKey(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	_, found := cache.Get("nonExistentKey")
	if found {
		t.Fatalf("Get should return false for non-existent key")
	}
}

func TestInMemCache_Delete(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	err := cache.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = cache.Delete("key1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, found := cache.Get("key1")
	if found {
		t.Fatalf("Get should return false for deleted key")
	}
}

func TestInMemCache_Delete_NonExistentKey(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	err := cache.Delete("nonExistentKey")
	if err == nil {
		t.Fatalf("Delete should fail for non-existent key")
	}
}

func TestInMemCache_Concurrency(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	const numOps = 100
	done := make(chan struct{}, numOps*2)

	for i := 0; i < numOps; i++ {
		go func(i int) {
			key := "key" + string(rune(i))
			err := cache.Set(key, i)
			if err != nil {
				return
			}
			done <- struct{}{}
		}(i)

		go func(i int) {
			key := "key" + string(rune(i))
			cache.Get(key)
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < numOps*2; i++ {
		<-done
	}
}
