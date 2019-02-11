package log_test

import (
	"testing"

	"github.com/hamba/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestNewMockLoggable(t *testing.T) {
	l := new(MockLogger)
	lable := log.NewMockLoggable(l)

	assert.Implements(t, (*log.Loggable)(nil), lable)
}

func TestMockLoggable_Logger(t *testing.T) {
	l := new(MockLogger)
	lable := log.NewMockLoggable(l)

	assert.Equal(t, l, lable.Logger())
}
