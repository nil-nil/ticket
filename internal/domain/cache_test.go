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
}
