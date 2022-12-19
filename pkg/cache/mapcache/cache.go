package mapcache

import (
	"errors"
	"sync"
	"time"
)

const TTLWithoutExpiration time.Duration = -1

var ErrItemNotFound = errors.New("cache: item not found")

type item struct {
	value     any
	createdAt time.Time
	ttl       time.Duration
}

type Cache struct {
	cache map[string]item
	sync.RWMutex
}

// New uses map to store key:value data in-memory.
func New() *Cache {
	c := &Cache{cache: make(map[string]item)}
	go c.setTTLTimer()

	return c
}

func (c *Cache) setTTLTimer() {
	for {
		c.Lock()
		for key, value := range c.cache {
			if value.ttl != TTLWithoutExpiration && time.Since(value.createdAt) > value.ttl {
				delete(c.cache, key)
			}
		}
		c.Unlock()

		<-time.After(time.Second)
	}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.Lock()
	c.cache[key] = item{
		value:     value,
		createdAt: time.Now(),
		ttl:       ttl,
	}
	c.Unlock()
}

func (c *Cache) Get(key string) (any, error) {
	c.RLock()
	it, ex := c.cache[key]
	c.RUnlock()

	if !ex {
		return nil, ErrItemNotFound
	}

	return it.value, nil
}

func (c *Cache) Del(key string) {
	c.Lock()
	delete(c.cache, key)
	c.Unlock()
}

func (c *Cache) TTL(key string) (time.Duration, error) {
	c.RLock()
	it, ex := c.cache[key]
	c.RUnlock()

	if !ex {
		return 0, ErrItemNotFound
	}

	return it.ttl - time.Since(it.createdAt), nil
}

func (c *Cache) Expire(key string, expiration time.Duration) error {
	c.RLock()
	it, ex := c.cache[key]
	c.RUnlock()

	if !ex {
		return ErrItemNotFound
	}

	c.Set(key, it.value, expiration)

	return nil
}
