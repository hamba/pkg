package reason_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hamba/pkg/v2/errors/reason"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtract(t *testing.T) {
	var errs error
	errs = errors.Join(errs, errors.New("test1"))
	errs = errors.Join(errs, reason.Error{Msg: "First Error"})
	errs = errors.Join(errs, fmt.Errorf("some error: %w", reason.Errorf("Second %s", "Error")))
	errs = errors.Join(errs, errors.New("test2"))

	reasons, errs := reason.Extract(errs)

	require.Error(t, errs)
	assert.Equal(t, "test1\ntest2", errs.Error())
	assert.Equal(t, []string{"First Error", "Second Error"}, reasons)
}
