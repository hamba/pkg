package request_test

import (
	"context"
	"testing"

	"github.com/hamba/pkg/v2/http/request"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	ctx := request.WithID(context.Background(), "my-id")

	got, ok := request.IDFrom(ctx)

	assert.True(t, ok)
	assert.Equal(t, "my-id", got)
}
