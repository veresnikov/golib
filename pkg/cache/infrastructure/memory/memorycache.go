package memory

import (
	"sync"
	"time"

	"github.com/veresnikov/golib/pkg/cache/application/memory"
)

type value struct {
	data interface{}
	ttl  *time.Time
}

func (v *value) isExpire() bool {
	if v.ttl != nil {
		return v.ttl.Unix() < time.Now().Unix()
	}
	return false
}

func NewMemoryCache(cleanupInterval time.Duration) memory.Cache {
	return &memoryCache{
		cache:    make(map[interface{}]value),
		stopChan: make(chan bool),
	}
}

type memoryCache struct {
	cache map[interface{}]value
	mu    sync.RWMutex

	stopChan chan bool

	cleanupInterval time.Duration
}

func (c *memoryCache) Get(key interface{}) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.cache[key]
	if !ok {
		return nil, memory.ErrKeyNotFound
	}
	if v.isExpire() {
		return nil, memory.ErrKeyExpired
	}
	return v.data, nil
}

func (c *memoryCache) Set(key interface{}, data interface{}, ttl *time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = value{
		data: data,
		ttl:  ttl,
	}
}

func (c *memoryCache) Delete(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

func (c *memoryCache) Start() {
	go func() {
		for {
			select {
			case <-c.stopChan:
				return
			default:
				<-time.After(c.cleanupInterval)
				c.cleanup()
			}
		}
	}()
}

func (c *memoryCache) Close() error {
	c.stopChan <- true
	return nil
}

func (c *memoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.cache {
		if v.isExpire() {
			delete(c.cache, k)
		}
	}
}
