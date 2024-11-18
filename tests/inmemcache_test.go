package tests

import (
	"github.com/StandyBee/go-inmemorycache/inmemcache"
	"testing"
	"time"
)

func TestInMemCache_SetAndGet(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	err := cache.Set("key1", 123, time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := cache.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != 123 {
		t.Fatalf("Get failed: %v, expected 123", val)
	}
}

func TestInMemCache_Set_EmptyKey(t *testing.T) {
	cache := inmemcache.NewInMemCache()
	err := cache.Set("", 0, time.Second)
	if err == nil {
		t.Fatalf("Set failed: %v", err)
	}
}

func TestInMemCache_Get_NonExistentKey(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	_, err := cache.Get("nonExistentKey")
	if err == nil {
		t.Fatalf("Get should return false for non-existent key")
	}
}

func TestInMemCache_Get_ExpiredKey(t *testing.T) {
	cache := inmemcache.NewInMemCache()
	err := cache.Set("key1", 123, time.Millisecond*10)

	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	time.Sleep(15 * time.Millisecond)

	_, err = cache.Get("key1")
	if err == nil || err.Error() != "key expired" {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestInMemCache_Delete(t *testing.T) {
	cache := inmemcache.NewInMemCache()

	err := cache.Set("key1", "value1", time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = cache.Delete("key1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = cache.Get("key1")
	if err == nil {
		t.Fatalf("Get should return err for deleted key")
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
			err := cache.Set(key, i, time.Second)
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
