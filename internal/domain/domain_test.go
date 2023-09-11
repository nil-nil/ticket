package domain_test

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Get the next key in a map of keys and values
func nextKey[T any](values map[uint64]T) uint64 {
	var lastKey uint64
	for k := range values {
		if k > lastKey {
			lastKey = k
		}
	}

	return lastKey + 1
}

func TestNextKeyFunc(t *testing.T) {
	t.Run("test first key", func(t *testing.T) {
		testMap := make(map[uint64]string, 0)
		k := nextKey(testMap)
		assert.Equal(t, uint64(1), k, "first key should be 1")
	})

	t.Run("test first key", func(t *testing.T) {
		testMap := map[uint64]struct{}{
			10: {},
			50: {},
			51: {},
		}
		k := nextKey(testMap)
		assert.Equal(t, uint64(52), k, "test next key with gaps")
	})
}

type mockCacheDriver struct {
	cache map[string]interface{}
}

func (m *mockCacheDriver) Get(key string) (interface{}, error) {
	val, ok := m.cache[key]
	if !ok {
		return nil, domain.ErrNotFoundInCache
	}
	return val, nil
}

func (m *mockCacheDriver) Set(key string, value interface{}) error {
	m.cache[key] = value
	return nil
}

func (m *mockCacheDriver) Forget(key string) error {
	delete(m.cache, key)
	return nil
}

var mockCache = &mockCacheDriver{
	cache: make(map[string]interface{}, 0),
}
