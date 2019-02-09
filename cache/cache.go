package cache

import (
	"errors"
	"time"
)

type key int

const (
	ctxKey key = iota
)

var (
	// ErrCacheMiss is returned if a Get failed because the item wasn't present.
	ErrCacheMiss = errors.New("cache: miss")

	// ErrNotStored is returned if conditional write (Add or Replace) failed because
	// the condition was not met.
	ErrNotStored = errors.New("cache: not stored")

	// Null is the null Cache instance.
	Null = &nullCache{}
)

// Cache represents a cache instance.
type Cache interface {
	// Get gets the item for the given key.
	Get(key string) *Item

	// GetMulti gets the items for the given keys.
	GetMulti(keys ...string) ([]*Item, error)

	// Set sets the item in the cache.
	Set(key string, value interface{}, expire time.Duration) error

	// Add sets the item in the cache, but only if the key does not already exist.
	Add(key string, value interface{}, expire time.Duration) error

	// Replace sets the item in the cache, but only if the key already exists.
	Replace(key string, value interface{}, expire time.Duration) error

	// Delete deletes the item with the given key.
	Delete(key string) error

	// Inc increments a key by the Value.
	Inc(key string, value uint64) (int64, error)

	// Dec decrements a key by the Value.
	Dec(key string, value uint64) (int64, error)
}

type nullDecoder struct{}

func (d nullDecoder) Bool(v []byte) (bool, error) {
	return false, nil
}

func (d nullDecoder) Int64(v []byte) (int64, error) {
	return 0, nil
}

func (d nullDecoder) Uint64(v []byte) (uint64, error) {
	return 0, nil
}

func (d nullDecoder) Float64(v []byte) (float64, error) {
	return 0, nil
}

type nullCache struct{}

// Get gets the item for the given key.
func (c nullCache) Get(key string) *Item {
	return &Item{Decoder: nullDecoder{}, Value: []byte{}}
}

// GetMulti gets the items for the given keys.
func (c nullCache) GetMulti(keys ...string) ([]*Item, error) {
	return []*Item{}, nil
}

// Set sets the item in the cache.
func (c nullCache) Set(key string, value interface{}, expire time.Duration) error {
	return nil
}

// Add sets the item in the cache, but only if the key does not already exist.
func (c nullCache) Add(key string, value interface{}, expire time.Duration) error {
	return nil
}

// Replace sets the item in the cache, but only if the key already exists.
func (c nullCache) Replace(key string, value interface{}, expire time.Duration) error {
	return nil
}

// Delete deletes the item with the given key.
func (c nullCache) Delete(key string) error {
	return nil
}

// Inc increments a key by the Value.
func (c nullCache) Inc(key string, value uint64) (int64, error) {
	return 0, nil
}

// Dec decrements a key by the Value.
func (c nullCache) Dec(key string, value uint64) (int64, error) {
	return 0, nil
}
