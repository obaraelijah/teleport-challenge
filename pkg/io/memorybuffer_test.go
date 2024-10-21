package io_test

import (
	"testing"

	"github.com/obaraelijah/teleport-challenge/pkg/io"
	"github.com/stretchr/testify/assert"
)

func Test_MemoryBuffer_InitialSizeZero(t *testing.T) {
	b := io.NewMemoryBuffer()

	assert.Equal(t, int64(0), b.Size())
}
