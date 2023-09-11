package domain_test

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestInvalidCacheKey(t *testing.T) {
	c, err := domain.NewCache[string]("", mockCache)
	assert.EqualError(t, err, domain.ErrCachePrefixInvalid.Error(), "invalid cache prefix should return a meaningful error")
	assert.Nil(t, c, "erroring NewCache should return a nil pointer")
}

func TestCache(t *testing.T) {
	c, err := domain.NewCache[string]("test", mockCache)
	assert.NoError(t, err, "valid new cache should not error")

	t.Run("TestInvalidCacheKey", func(t *testing.T) {
		err := c.Set("", "")
		assert.EqualError(t, err, domain.ErrCacheKeyInvalid.Error(), "Cache.Set(\"\", T) should return ErrInvalidCacheKey")

		v, err := c.Get("")
		assert.EqualError(t, err, domain.ErrCacheKeyInvalid.Error(), "Cache.Get(\"\") should return ErrInvalidCacheKey")
		assert.Zero(t, v, "erroring Cache.Get() should return zero value")

		err = c.Forget("")
		assert.EqualError(t, err, domain.ErrCacheKeyInvalid.Error(), "Cache.Forget(\"\") should return ErrInvalidCacheKey")
	})

	t.Run("TestSettingCache", func(t *testing.T) {
		err := c.Set("1", "examplestring")
		assert.NoError(t, err, "valid Cache.Set() should not error")
		assert.Equal(t, "examplestring", mockCache.cache["test.1"])
	})

	t.Run("TestGettingValidCache", func(t *testing.T) {
		val, err := c.Get("1")
		assert.NoError(t, err, "valid Cache.Get() should not error")
		assert.Equal(t, "examplestring", val)
	})

	t.Run("TestGettingInvalidTypedCache", func(t *testing.T) {
		mockCache.cache["test.2"] = 100
		val, err := c.Get("2")
		assert.EqualError(t, err, domain.ErrCachedValueInvalid.Error(), "invalid Cache.Get() should return meaningful error")
		assert.Zero(t, val, "errored Cache.Get() should return zero value")
	})

	t.Run("TestGettingEmptyCache", func(t *testing.T) {
		val, err := c.Get("3")
		assert.EqualError(t, err, domain.ErrNotFoundInCache.Error(), "empty Cache.Get() should return meaningful error")
		assert.Zero(t, val, "errored Cache.Get() should return zero value")
	})

	t.Run("TestForgettingCache", func(t *testing.T) {
		err := c.Forget("1")
		assert.NoError(t, err, "valid Cache.Forget() should not error")
		_, ok := mockCache.cache["test.1"]
		assert.False(t, ok, "forgotten cache item should be removed")
	})
}
