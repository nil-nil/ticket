package domain

import "errors"

var (
	ErrNotFoundInCache = errors.New("key not found in cache")
)

type CacheDriver interface {
	Set(key string, value interface{}) error
	Forget(key string) error
	Get(key string) (interface{}, error)
}
