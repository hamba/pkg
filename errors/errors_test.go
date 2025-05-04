package errors_test

import (
	"testing"

	errorsx "github.com/hamba/pkg/v2/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const errTest = errorsx.Error("test error")

func TestError(t *testing.T) {
	err := errTest

	require.EqualError(t, err, "test error")
	assert.ErrorIs(t, err, errTest)
}
