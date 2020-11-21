package cache

import (
	"errors"
	"time"
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
	Get(k string) Item

	// GetMulti gets the items for the given keys.
	GetMulti(ks ...string) ([]Item, error)

	// Set sets the item in the cache.
	Set(k string, v interface{}, expire time.Duration) error

	// Add sets the item in the cache, but only if the key does not already exist.
	Add(k string, v interface{}, expire time.Duration) error

	// Replace sets the item in the cache, but only if the key already exists.
	Replace(k string, v interface{}, expire time.Duration) error

	// Delete deletes the item with the given key.
	Delete(k string) error

	// Inc increments a key by the Value.
	Inc(k string, v uint64) (int64, error)

	// Dec decrements a key by the Value.
	Dec(k string, v uint64) (int64, error)
}

type nullDecoder struct{}

func (d nullDecoder) Bool(v interface{}) (bool, error) {
	return false, nil
}

func (d nullDecoder) Bytes(v interface{}) ([]byte, error) {
	return []byte{}, nil
}

func (d nullDecoder) Int64(v interface{}) (int64, error) {
	return 0, nil
}

func (d nullDecoder) Uint64(v interface{}) (uint64, error) {
	return 0, nil
}

func (d nullDecoder) Float64(v interface{}) (float64, error) {
	return 0, nil
}

func (d nullDecoder) String(v interface{}) (string, error) {
	return "", nil
}

type nullCache struct{}

// Get gets the item for the given key.
func (c nullCache) Get(k string) Item {
	return NewItem(nullDecoder{}, nil, nil)
}

// GetMulti gets the items for the given keys.
func (c nullCache) GetMulti(ks ...string) ([]Item, error) {
	return []Item{}, nil
}

// Set sets the item in the cache.
func (c nullCache) Set(k string, v interface{}, expire time.Duration) error {
	return nil
}

// Add sets the item in the cache, but only if the key does not already exist.
func (c nullCache) Add(k string, v interface{}, expire time.Duration) error {
	return nil
}

// Replace sets the item in the cache, but only if the key already exists.
func (c nullCache) Replace(k string, v interface{}, expire time.Duration) error {
	return nil
}

// Delete deletes the item with the given key.
func (c nullCache) Delete(k string) error {
	return nil
}

// Inc increments a key by the Value.
func (c nullCache) Inc(k string, v uint64) (int64, error) {
	return 0, nil
}

// Dec decrements a key by the Value.
func (c nullCache) Dec(k string, v uint64) (int64, error) {
	return 0, nil
}
