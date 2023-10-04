package ristrettocache_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/infrastructure/ristrettocache"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	ristrettoCache, err := ristrettocache.NewCache(&ristrettocache.CacheConfig{Size: 8 * 512, ExpectedNumberOfItems: 1})
	assert.NoError(t, err, "creating ristretto cache shouldn't error")

	t.Run("not found error", func(t *testing.T) {
		val, err := ristrettoCache.Get("notexist")
		assert.Nil(t, val, "nil value should be returned for missing item")
		assert.EqualError(t, err, domain.ErrNotFoundInCache.Error(), "missing item should return meaningful error message")
	})

	t.Run("normal process", func(t *testing.T) {
		item := domain.User{ID: uuid.New(), FirstName: "Bob", LastName: "Test"}
		key := "testitem"
		err := ristrettoCache.Set(key, item)
		assert.NoError(t, err, "setting cache item should probably not error")

		// Allow the item to make it into the cache
		time.Sleep(1 * time.Millisecond)

		hit, err := ristrettoCache.Get(key)
		assert.NoError(t, err, "getting cache item should probably not error")

		u, ok := hit.(domain.User)
		assert.True(t, ok, "interface conversion should work to original type")
		assert.Equal(t, item, u, "set value and got value should match")

		err = ristrettoCache.Forget(key)
		assert.NoError(t, err, "forgetting cache item should probably not error")

		val, err := ristrettoCache.Get(key)
		assert.Nil(t, val, "nil value should be returned for forgotten item")
		assert.EqualError(t, err, domain.ErrNotFoundInCache.Error(), "forgotten item should return missing error message")
	})
}
