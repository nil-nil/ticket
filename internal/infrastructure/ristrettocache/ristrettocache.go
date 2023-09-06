package ristrettocache

import (
	"github.com/dgraph-io/ristretto"
	"github.com/nil-nil/ticket/internal/domain"
)

const (
	Kilobyte = int64(1024)
	Megabyte = Kilobyte * 1024
	Gigabyte = Megabyte * 1024
)

type CacheConfig struct {
	Size                  int64
	ExpectedNumberOfItems int64
}

func NewCache(config *CacheConfig) (*RistrettoCache, error) {
	size := int64(256 * Megabyte)
	expectedNumberOfItems := int64(1000)
	if config != nil && config.Size != 0 {
		size = config.Size
	}
	if config != nil && config.ExpectedNumberOfItems != 0 {
		expectedNumberOfItems = config.ExpectedNumberOfItems
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: expectedNumberOfItems * 10,
		MaxCost:     size,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	return &RistrettoCache{
		ristretto: cache,
	}, nil
}

type RistrettoCache struct {
	ristretto *ristretto.Cache
}

func (r *RistrettoCache) Get(key string) (interface{}, error) {
	val, found := r.ristretto.Get(key)
	if !found {
		return nil, domain.ErrNotFoundInCache
	}

	return val, nil
}

func (r *RistrettoCache) Set(key string, value interface{}) error {
	r.ristretto.Set(key, value, 1)
	return nil
}

func (r *RistrettoCache) Forget(key string) error {
	r.ristretto.Del(key)
	return nil
}
