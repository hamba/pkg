package cache

// Decoder represents a byte decoder.
type Decoder interface {
	Bool([]byte) (bool, error)
	Int64([]byte) (int64, error)
	Uint64([]byte) (uint64, error)
	Float64([]byte) (float64, error)
}

// Item represents an item to be returned or stored in the cache
type Item struct {
	Decoder Decoder
	Value   []byte
	Err     error
}

// Bool gets the cache items Value as a bool, or and error.
func (i Item) Bool() (bool, error) {
	if i.Err != nil {
		return false, i.Err
	}

	return i.Decoder.Bool(i.Value)
}

// Bytes gets the cache items Value as bytes.
func (i Item) Bytes() ([]byte, error) {
	return i.Value, i.Err
}

// Bytes gets the cache items Value as a string.
func (i Item) String() (string, error) {
	if i.Err != nil {
		return "", i.Err
	}

	return string(i.Value), nil
}

// Int64 gets the cache items Value as an int64, or and error.
func (i Item) Int64() (int64, error) {
	if i.Err != nil {
		return 0, i.Err
	}

	return i.Decoder.Int64(i.Value)
}

// Uint64 gets the cache items Value as a uint64, or and error.
func (i Item) Uint64() (uint64, error) {
	if i.Err != nil {
		return 0, i.Err
	}

	return i.Decoder.Uint64(i.Value)
}

// Float64 gets the cache items Value as a float64, or and error.
func (i Item) Float64() (float64, error) {
	if i.Err != nil {
		return 0, i.Err
	}

	return i.Decoder.Float64(i.Value)
}

// Err returns the item error or nil.
func (i Item) Error() error {
	return i.Err
}
