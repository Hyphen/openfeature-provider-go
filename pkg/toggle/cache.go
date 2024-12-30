package toggle

import (
	"sync"
	"time"
)

type Cache struct {
	ttl    time.Duration
	keyGen func(ctx EvaluationContext) string
	store  sync.Map
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func newCache(config *CacheConfig) *Cache {
	ttl := DefaultCacheTTL
	if config.TTLSeconds > 0 {
		ttl = config.TTLSeconds
	}

	return &Cache{
		ttl:    time.Duration(ttl) * time.Second,
		keyGen: config.KeyGen,
	}
}

func (c *Cache) Get(ctx EvaluationContext) interface{} {
	key := c.keyGen(ctx)
	if item, ok := c.store.Load(key); ok {
		cacheItem := item.(cacheItem)
		if time.Now().Before(cacheItem.expiration) {
			return cacheItem.value
		}
		c.store.Delete(key)
	}
	return nil
}

func (c *Cache) Set(ctx EvaluationContext, value interface{}) {
	key := c.keyGen(ctx)
	c.store.Store(key, cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	})
}
