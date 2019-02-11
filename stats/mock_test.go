package stats_test

import (
	"testing"

	"github.com/hamba/pkg/stats"
	"github.com/stretchr/testify/assert"
)

func TestNewMockStatable(t *testing.T) {
	s := new(MockStats)
	sable := stats.NewMockStatable(s)

	assert.Implements(t, (*stats.Statable)(nil), sable)
}

func TestMockStatable_Statter(t *testing.T) {
	s := new(MockStats)
	sable := stats.NewMockStatable(s)

	assert.Equal(t, s, sable.Statter())
}
