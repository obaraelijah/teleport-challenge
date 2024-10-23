package io_test

import (
	"sync"
	"testing"

	"github.com/obaraelijah/teleport-challenge/pkg/io"
	"github.com/stretchr/testify/assert"
)

func Test_ByteStream(t *testing.T) {
	const iterationCount = 5
	payload := []byte("hello")
	output := make([]byte, 0, iterationCount*len(payload))

	buffer := io.NewMemoryBuffer()
	stream := io.NewByteStream(buffer)

	byteCount := 0
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for content := range stream.Stream() {
			byteCount += len(content)
			output = append(output, content...)
		}
		wg.Done()
	}()

	for i := 0; i < iterationCount; i++ {
		buffer.Write(payload)
	}

	buffer.Close()
	wg.Wait()

	assert.Equal(t, iterationCount*len(payload), byteCount)
	assert.Equal(t, []byte("hellohellohellohellohello"), output)
}

func Test_ByteStream_MultipleCallsToStream(t *testing.T) {
	buffer := io.NewMemoryBuffer()
	stream := io.NewByteStream(buffer)

	buffer.Write([]byte("hello"))
	assert.Equal(t, []byte("hello"), <-stream.Stream())

	buffer.Write([]byte("world"))
	assert.Equal(t, []byte("world"), <-stream.Stream())

	buffer.Close()
	assert.Nil(t, <-stream.Stream())
}

func Test_ByteStream_MultipleReaders(t *testing.T) {
	buffer := io.NewMemoryBuffer()
	stream1 := io.NewByteStream(buffer).Stream()
	stream2 := io.NewByteStream(buffer).Stream()

	buffer.Write([]byte("hello"))
	assert.Equal(t, []byte("hello"), <-stream1)
	assert.Equal(t, []byte("hello"), <-stream2)

	buffer.Write([]byte("world"))
	assert.Equal(t, []byte("world"), <-stream1)
	assert.Equal(t, []byte("world"), <-stream2)

	buffer.Close()
	assert.Nil(t, <-stream1)
	assert.Nil(t, <-stream2)
}
