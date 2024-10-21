package io

import (
	"fmt"
	"sync"
)

type MemoryBuffer struct {
	content []byte
	closed  bool
	mutex   sync.RWMutex
}

// The default initial capacity of a memory buffer created using NewMemoryBuffer.
const DefaultInitialMemoryBufferCapacity = 4096

func NewMemoryBuffer() *MemoryBuffer {
	return NewMemoryBufferDetailed(DefaultInitialMemoryBufferCapacity)
}

// NewMemoryBufferDetailed creates and returns a MemoryBuffer with the given
// initialCapacity.
func NewMemoryBufferDetailed(initialCapacity int) *MemoryBuffer {
	return &MemoryBuffer{
		content: make([]byte, 0, initialCapacity),
	}
}

func (b *MemoryBuffer) Write(newContent []byte) (bytesWritten int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.closed {
		return 0, fmt.Errorf("cannot write to a closed MemoryBuffer")
	}

	b.content = append(b.content, newContent...)

	return len(newContent), nil
}

func (b *MemoryBuffer) ReadAt(outputBuffer []byte, offset int64) (bytesRead int, err error) {
	if len(outputBuffer) == 0 {
		return 0, nil
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return copy(outputBuffer, b.content[offset:]), nil
}
