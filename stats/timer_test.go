package stats_test

import (
	"context"
	"testing"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/mock"
)

func TestTimer(t *testing.T) {
	m := new(MockStats)
	m.On("Timing", "test", mock.Anything, float32(1.0), mock.Anything).Return(nil)

	ctx := stats.WithStatter(context.Background(), m)
	ti := stats.Time(ctx, "test", 1.0)
	ti.Done()

	m.AssertExpectations(t)
}
