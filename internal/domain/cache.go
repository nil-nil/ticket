package domain

import "errors"

var (
	ErrNotFoundInCache    = errors.New("key not found in cache")
	ErrCachedValueInvalid = errors.New("cached value is invalid")
	ErrCachePrefixInvalid = errors.New("not a valid cache prefix")
	ErrCacheKeyInvalid    = errors.New("not a valid cache key")
)

type CacheDriver interface {
	Set(key string, value interface{}) error
	Forget(key string) error
	Get(key string) (interface{}, error)
}

func NewCache[T any](prefix string, driver CacheDriver) (*Cache[T], error) {
	if prefix == "" {
		return nil, ErrCachePrefixInvalid
	}

	return &Cache[T]{
		prefix: prefix,
		driver: driver,
	}, nil
}

type Cache[T any] struct {
	prefix string
	driver CacheDriver
}

func (c *Cache[T]) Set(key string, value T) error {
	if key == "" {
		return ErrCacheKeyInvalid
	}
	return c.driver.Set(c.prefix+"."+key, value)
}

func (c *Cache[T]) Forget(key string) error {
	if key == "" {
		return ErrCacheKeyInvalid
	}
	return c.driver.Forget(c.prefix + "." + key)
}

func (c *Cache[T]) Get(key string) (T, error) {
	if key == "" {
		var zeroValue T
		return zeroValue, ErrCacheKeyInvalid
	}
	val, err := c.driver.Get(c.prefix + "." + key)
	if err != nil {
		var zeroValue T
		return zeroValue, err
	}

	data, ok := val.(T)
	if !ok {
		var zeroValue T
		return zeroValue, ErrCachedValueInvalid
	}

	return data, nil
}
