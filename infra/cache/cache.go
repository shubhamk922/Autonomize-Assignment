package cache

import "time"

// Cache interface abstracts any cache implementation
type Cache interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
}
