package memory

import (
	"errors"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key expired")
)

type Cache interface {
	Get(key interface{}) (interface{}, error)
	Set(key interface{}, data interface{}, ttl *time.Time)
	Delete(key interface{})

	Start()
	Close() error
}
