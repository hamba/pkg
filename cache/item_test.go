package cache_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/hamba/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func TestItem_Bool(t *testing.T) {
	decoder := stringDecoder{}
	tests := []struct {
		item   cache.Item
		ok     bool
		expect bool
	}{
		{cache.Item{Decoder: decoder, Value: []byte("1")}, true, true},
		{cache.Item{Decoder: decoder, Err: errors.New("")}, false, false},
	}

	for i, tt := range tests {
		got, err := tt.item.Bool()
		if ok := err == nil; ok != tt.ok {
			if err != nil {
				assert.FailNow(t, "test %d, unexpected failure: %v", i, err)
			} else {
				assert.FailNow(t, "test %d, unexpected success", i)
			}
		}

		assert.Equal(t, tt.expect, got)
	}
}

func TestItem_Bytes(t *testing.T) {
	tests := []struct {
		item   cache.Item
		ok     bool
		expect []byte
	}{
		{cache.Item{Value: []byte{0x01}}, true, []byte{0x01}},
		{cache.Item{Err: errors.New("")}, false, nil},
	}

	for i, tt := range tests {
		got, err := tt.item.Bytes()
		if ok := err == nil; ok != tt.ok {
			if err != nil {
				assert.FailNow(t, "test %d, unexpected failure: %v", i, err)
			} else {
				assert.FailNow(t, "test %d, unexpected success", i)
			}
		}

		assert.Equal(t, tt.expect, got)
	}
}

func TestItem_String(t *testing.T) {
	tests := []struct {
		item   cache.Item
		ok     bool
		expect string
	}{
		{cache.Item{Value: []byte("hello")}, true, "hello"},
		{cache.Item{Err: errors.New("")}, false, ""},
	}

	for i, tt := range tests {
		got, err := tt.item.String()
		if ok := err == nil; ok != tt.ok {
			if err != nil {
				assert.FailNow(t, "test %d, unexpected failure: %v", i, err)
			} else {
				assert.FailNow(t, "test %d, unexpected success", i)
			}
		}

		assert.Equal(t, tt.expect, got)
	}
}

func TestItem_Int64(t *testing.T) {
	decoder := stringDecoder{}
	tests := []struct {
		item   cache.Item
		ok     bool
		expect int64
	}{
		{cache.Item{Decoder: decoder, Value: []byte("1")}, true, 1},
		{cache.Item{Decoder: decoder, Value: []byte("a")}, false, 0},
		{cache.Item{Decoder: decoder, Err: errors.New("")}, false, 0},
	}

	for i, tt := range tests {
		got, err := tt.item.Int64()
		if ok := err == nil; ok != tt.ok {
			if err != nil {
				assert.FailNow(t, "test %d, unexpected failure: %v", i, err)
			} else {
				assert.FailNow(t, "test %d, unexpected success", i)
			}
		}

		assert.Equal(t, tt.expect, got)
	}
}

func TestItem_Uint64(t *testing.T) {
	decoder := stringDecoder{}
	tests := []struct {
		item   cache.Item
		ok     bool
		expect uint64
	}{
		{cache.Item{Decoder: decoder, Value: []byte("1")}, true, 1},
		{cache.Item{Decoder: decoder, Value: []byte("a")}, false, 0},
		{cache.Item{Decoder: decoder, Err: errors.New("")}, false, 0},
	}

	for i, tt := range tests {
		got, err := tt.item.Uint64()
		if ok := (err == nil); ok != tt.ok {
			if err != nil {
				assert.FailNow(t, "test %d, unexpected failure: %v", i, err)
			} else {
				assert.FailNow(t, "test %d, unexpected success", i)
			}
		}

		assert.Equal(t, tt.expect, got)
	}
}

func TestItem_Float64(t *testing.T) {
	decoder := stringDecoder{}
	tests := []struct {
		item   cache.Item
		ok     bool
		expect float64
	}{
		{cache.Item{Decoder: decoder, Value: []byte("1.2")}, true, 1.2},
		{cache.Item{Decoder: decoder, Value: []byte("a")}, false, 0},
		{cache.Item{Decoder: decoder, Err: errors.New("")}, false, 0},
	}

	for i, tt := range tests {
		got, err := tt.item.Float64()
		if ok := (err == nil); ok != tt.ok {
			if err != nil {
				assert.FailNow(t, "test %d, unexpected failure: %v", i, err)
			} else {
				assert.FailNow(t, "test %d, unexpected success", i)
			}
		}

		assert.Equal(t, tt.expect, got)
	}
}

func TestItem_Err(t *testing.T) {
	expect := errors.New("Err")
	tests := []struct {
		item   cache.Item
		expect error
	}{
		{cache.Item{}, nil},
		{cache.Item{Err: expect}, expect},
	}

	for _, tt := range tests {
		err := tt.item.Error()

		assert.Equal(t, tt.expect, err)
	}
}

type stringDecoder struct{}

func (d stringDecoder) Bool(v []byte) (bool, error) {
	return string(v) == "1", nil
}

func (d stringDecoder) Int64(v []byte) (int64, error) {
	return strconv.ParseInt(string(v), 10, 64)
}

func (d stringDecoder) Uint64(v []byte) (uint64, error) {
	return strconv.ParseUint(string(v), 10, 64)
}

func (d stringDecoder) Float64(v []byte) (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}
