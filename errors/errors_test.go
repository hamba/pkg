package errors_test

import (
	"errors"
	"testing"

	errorsx "github.com/hamba/pkg/v2/errors"
	"github.com/stretchr/testify/assert"
)

const testErr = errorsx.Error("test error")

func TestError(t *testing.T) {
	err := testErr

	assert.EqualError(t, err, "test error")
	assert.True(t, errors.Is(err, testErr))
}
