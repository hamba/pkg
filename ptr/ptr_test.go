package ptr_test

import (
	"testing"

	"github.com/hamba/pkg/v2/ptr"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	want := true

	got := ptr.Of(want)

	assert.Exactly(t, &want, got)
}
