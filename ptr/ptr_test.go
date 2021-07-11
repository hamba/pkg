package ptr_test

import (
	"testing"

	"github.com/hamba/pkg/v2/ptr"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	want := true

	got := ptr.Bool(want)

	assert.Exactly(t, &want, got)
}

func TestFloat32(t *testing.T) {
	want := float32(1.0)

	got := ptr.Float32(want)

	assert.Exactly(t, &want, got)
}

func TestFloat64(t *testing.T) {
	want := float64(1.0)

	got := ptr.Float64(want)

	assert.Exactly(t, &want, got)
}

func TestInt(t *testing.T) {
	want := 1

	got := ptr.Int(want)

	assert.Exactly(t, &want, got)
}

func TestInt64(t *testing.T) {
	want := int64(1)

	got := ptr.Int64(want)

	assert.Exactly(t, &want, got)
}

func TestString(t *testing.T) {
	want := "foo"

	got := ptr.String(want)

	assert.Exactly(t, &want, got)
}
