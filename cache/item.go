package cache

// Decoder represents a value decoder.
type Decoder interface {
	Bool(interface{}) (bool, error)
	Bytes(interface{}) ([]byte, error)
	Int64(interface{}) (int64, error)
	Uint64(interface{}) (uint64, error)
	Float64(interface{}) (float64, error)
	String(interface{}) (string, error)
}

// Item represents an item to be returned or stored in the cache.
type Item struct {
	dec   Decoder
	Value interface{}
	Err   error
}

// NewItem creates a new Item.
func NewItem(d Decoder, v interface{}, err error) Item {
	return Item{
		dec:   d,
		Value: v,
		Err:   err,
	}
}

// Bool gets the cache items Value as a bool, or and error.
func (i Item) Bool() (bool, error) {
	if i.Err != nil {
		return false, i.Err
	}

	return i.dec.Bool(i.Value)
}

// Bytes gets the cache items Value as bytes.
func (i Item) Bytes() ([]byte, error) {
	if i.Err != nil {
		return nil, i.Err
	}

	return i.dec.Bytes(i.Value)
}

// Int64 gets the cache items Value as an int64, or and error.
func (i Item) Int64() (int64, error) {
	if i.Err != nil {
		return 0, i.Err
	}

	return i.dec.Int64(i.Value)
}

// Uint64 gets the cache items Value as a uint64, or and error.
func (i Item) Uint64() (uint64, error) {
	if i.Err != nil {
		return 0, i.Err
	}

	return i.dec.Uint64(i.Value)
}

// Float64 gets the cache items Value as a float64, or and error.
func (i Item) Float64() (float64, error) {
	if i.Err != nil {
		return 0, i.Err
	}

	return i.dec.Float64(i.Value)
}

// Bytes gets the cache items Value as a string.
func (i Item) String() (string, error) {
	if i.Err != nil {
		return "", i.Err
	}

	return i.dec.String(i.Value)
}
