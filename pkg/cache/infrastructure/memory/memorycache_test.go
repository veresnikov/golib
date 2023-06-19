package memory

import (
	"testing"
	"time"

	"github.com/veresnikov/golib/pkg/cache/application/memory"
)

type testCase struct {
	key  interface{}
	data interface{}
	ttl  *time.Time
}

var testCases = []testCase{
	{
		key:  struct{}{},
		data: struct{}{},
		ttl:  getTTL(time.Second),
	},
	{
		key:  struct{}{},
		data: struct{}{},
		ttl:  nil,
	},
}

func TestMemoryCache_GetSet(t *testing.T) {
	cache := &memoryCache{
		cache:           make(map[interface{}]value),
		stopChan:        make(chan bool),
		cleanupInterval: toOptionalDuration(time.Microsecond * 10),
	}
	go func() {
		cache.Start()
	}()
	defer func() {
		_ = cache.Close()
	}()
	for _, testData := range testCases {
		cache.Set(testData.key, testData.data, testData.ttl)
		data, err := cache.Get(testData.key)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if testData.data != data {
			t.Error("unexpected data value")
		}
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := &memoryCache{
		cache:           make(map[interface{}]value),
		stopChan:        make(chan bool),
		cleanupInterval: toOptionalDuration(time.Microsecond * 10),
	}
	go func() {
		cache.Start()
	}()
	defer func() {
		_ = cache.Close()
	}()
	for _, testData := range testCases {
		cache.Set(testData.key, testData.data, testData.ttl)
		cache.Delete(testData.key)
		_, err := cache.Get(testData.key)
		if err != nil {
			if err == memory.ErrKeyNotFound {
				continue
			}
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestMemoryCache_Cleanup(_ *testing.T) {
	cache := &memoryCache{
		cache:           make(map[interface{}]value),
		stopChan:        make(chan bool),
		cleanupInterval: toOptionalDuration(time.Microsecond * 10),
	}
	go func() {
		cache.Start()
	}()
	defer func() {
		_ = cache.Close()
	}()

	ttlDuration := time.Microsecond * 50
	for i := 0; i < 10; i++ {
		cache.Set(struct{}{}, i, getTTL(ttlDuration))
	}
	for {
		if len(cache.cache) == 0 {
			break
		}
	}
}

func getTTL(ttl time.Duration) *time.Time {
	result := time.Now().Add(ttl)
	return &result
}

func toOptionalDuration(duration time.Duration) *time.Duration {
	return &duration
}
