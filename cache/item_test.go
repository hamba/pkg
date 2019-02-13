package cache_test

import (
	"errors"
	"testing"

	"github.com/hamba/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItem_Bool(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Bool", []byte("1")).Return(true, nil)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("1"),
		Err:     nil,
	}

	got, err := item.Bool()

	assert.NoError(t, err)
	assert.Equal(t, true, got)
	dec.AssertExpectations(t)
}

func TestItem_BoolDecoderError(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Bool", []byte("1")).Return(false, errors.New("test"))
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("1"),
		Err:     nil,
	}

	_, err := item.Bool()

	assert.Error(t, err)
	dec.AssertExpectations(t)
}

func TestItem_BoolError(t *testing.T) {
	dec := new(MockDecoder)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("1"),
		Err:     errors.New("test"),
	}

	_, err := item.Bool()

	assert.Error(t, err)
}

func TestItem_Bytes(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Bytes", []byte{0x01}).Return([]byte{0x01}, nil)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte{0x01},
		Err:     nil,
	}

	got, err := item.Bytes()

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01}, got)
	dec.AssertExpectations(t)
}

func TestItem_BytesDecoderError(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Bytes", []byte{0x01}).Return([]byte{}, errors.New("test"))
	item := cache.Item{
		Decoder: dec,
		Value:   []byte{0x01},
		Err:     nil,
	}

	_, err := item.Bytes()

	assert.Error(t, err)
	dec.AssertExpectations(t)
}

func TestItem_BytesError(t *testing.T) {
	dec := new(MockDecoder)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte{0x01},
		Err:     errors.New("test"),
	}

	_, err := item.Bytes()

	assert.Error(t, err)
}

func TestItem_Int64(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Int64", []byte("2")).Return(int64(2), nil)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2"),
		Err:     nil,
	}

	got, err := item.Int64()

	assert.NoError(t, err)
	assert.Equal(t, int64(2), got)
	dec.AssertExpectations(t)
}

func TestItem_Int64DecoderError(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Int64", []byte("2")).Return(int64(0), errors.New("test"))
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2"),
		Err:     nil,
	}

	_, err := item.Int64()

	assert.Error(t, err)
	dec.AssertExpectations(t)
}

func TestItem_Int64Error(t *testing.T) {
	dec := new(MockDecoder)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2"),
		Err:     errors.New("test"),
	}

	_, err := item.Int64()

	assert.Error(t, err)
}

func TestItem_Uint64(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Uint64", []byte("2")).Return(uint64(2), nil)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2"),
		Err:     nil,
	}

	got, err := item.Uint64()

	assert.NoError(t, err)
	assert.Equal(t, uint64(2), got)
	dec.AssertExpectations(t)
}

func TestItem_Uint64DecoderError(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Uint64", []byte("2")).Return(uint64(0), errors.New("test"))
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2"),
		Err:     nil,
	}

	_, err := item.Uint64()

	assert.Error(t, err)
	dec.AssertExpectations(t)
}

func TestItem_Uint64Error(t *testing.T) {
	dec := new(MockDecoder)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2"),
		Err:     errors.New("test"),
	}

	_, err := item.Uint64()

	assert.Error(t, err)
}

func TestItem_Float64(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Float64", []byte("2.3")).Return(float64(2.3), nil)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2.3"),
		Err:     nil,
	}

	got, err := item.Float64()

	assert.NoError(t, err)
	assert.Equal(t, float64(2.3), got)
	dec.AssertExpectations(t)
}

func TestItem_Float64DecoderError(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("Float64", []byte("2.3")).Return(float64(0), errors.New("test"))
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2.3"),
		Err:     nil,
	}

	_, err := item.Float64()

	assert.Error(t, err)
	dec.AssertExpectations(t)
}

func TestItem_Float64Error(t *testing.T) {
	dec := new(MockDecoder)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("2.3"),
		Err:     errors.New("test"),
	}

	_, err := item.Float64()

	assert.Error(t, err)
}

func TestItem_String(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("String", []byte("test")).Return("test", nil)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("test"),
		Err:     nil,
	}

	got, err := item.String()

	assert.NoError(t, err)
	assert.Equal(t, "test", got)
	dec.AssertExpectations(t)
}

func TestItem_StringDecoderError(t *testing.T) {
	dec := new(MockDecoder)
	dec.On("String", []byte("test")).Return("", errors.New("test"))
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("test"),
		Err:     nil,
	}

	_, err := item.String()

	assert.Error(t, err)
	dec.AssertExpectations(t)
}

func TestItem_StringError(t *testing.T) {
	dec := new(MockDecoder)
	item := cache.Item{
		Decoder: dec,
		Value:   []byte("test"),
		Err:     errors.New("test"),
	}

	_, err := item.String()

	assert.Error(t, err)
}

type MockDecoder struct {
	mock.Mock
}

func (m *MockDecoder) Bool(v interface{}) (bool, error) {
	args := m.Called(v)

	return args.Bool(0), args.Error(1)
}

func (m *MockDecoder) Bytes(v interface{}) ([]byte, error) {
	args := m.Called(v)

	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockDecoder) Int64(v interface{}) (int64, error) {
	args := m.Called(v)

	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDecoder) Uint64(v interface{}) (uint64, error) {
	args := m.Called(v)

	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDecoder) Float64(v interface{}) (float64, error) {
	args := m.Called(v)

	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDecoder) String(v interface{}) (string, error) {
	args := m.Called(v)

	return args.String(0), args.Error(1)
}
