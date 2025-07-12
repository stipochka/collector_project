package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const janitoInterval = time.Minute * 1

type CacheValue struct {
	Value     any
	expiresAt time.Time
}

type TTLCache interface {
	Get(ctx context.Context, key string) (*CacheValue, bool)
	Set(ctx context.Context, key string, value interface{}) error
	Remove(ctx context.Context, key string)
	StartJanitor(ctx context.Context)
}

type InMemoryCache struct {
	data sync.Map
	ttl  time.Duration
}

func NewInMemoryCache(ttl time.Duration) *InMemoryCache {
	return &InMemoryCache{
		data: sync.Map{},
		ttl:  ttl,
	}
}

func (ic *InMemoryCache) Get(ctx context.Context, key string) (*CacheValue, bool) {
	if err := ctx.Err(); err != nil {
		return nil, false
	}

	value, exists := ic.data.Load(key)
	if !exists {
		return nil, false
	}

	cv, ok := value.(*CacheValue)
	if !ok {
		return nil, false
	}

	if time.Now().After(cv.expiresAt) {
		ic.data.Delete(key)
		return nil, false
	}

	return cv, ok

}

func (ic *InMemoryCache) Set(ctx context.Context, key string, value interface{}) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled")
	}

	ic.data.Store(key, &CacheValue{
		Value:     value,
		expiresAt: time.Now().Add(ic.ttl),
	})
	return nil
}

func (ic *InMemoryCache) Remove(_ context.Context, key string) {
	ic.data.Delete(key)
}

func (ic *InMemoryCache) StartJanitor(ctx context.Context) {

	clearTicker := time.NewTicker(janitoInterval)
	defer clearTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-clearTicker.C:
			ic.data.Range(func(key, value any) bool {
				cv, ok := value.(*CacheValue)
				if ok && time.Now().After(cv.expiresAt) {
					ic.data.Delete(key)
				}
				return true
			})
		}

	}
}
