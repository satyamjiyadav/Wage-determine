package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type cacheItem struct {
	value      []byte
	expiration time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewMemoryCache() *MemoryCache {
	mc := &MemoryCache{
		items: make(map[string]cacheItem),
	}
	// Background cleanup loop
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			mc.cleanup()
		}
	}()
	return mc
}

func (mc *MemoryCache) cleanup() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	now := time.Now()
	for k, item := range mc.items {
		if !item.expiration.IsZero() && now.After(item.expiration) {
			delete(mc.items, k)
		}
	}
}

func (mc *MemoryCache) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	mc.mu.RLock()
	item, found := mc.items[key]
	mc.mu.RUnlock()

	if !found {
		return false, nil
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		mc.mu.Lock()
		delete(mc.items, key)
		mc.mu.Unlock()
		return false, nil
	}

	err := json.Unmarshal(item.value, target)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (mc *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var exp time.Time
	if ttlSeconds > 0 {
		exp = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	}

	mc.mu.Lock()
	mc.items[key] = cacheItem{
		value:      bytes,
		expiration: exp,
	}
	mc.mu.Unlock()
	return nil
}

func (mc *MemoryCache) Invalidate(ctx context.Context, pattern string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	// Clear all for simplicity
	mc.items = make(map[string]cacheItem)
	return nil
}
