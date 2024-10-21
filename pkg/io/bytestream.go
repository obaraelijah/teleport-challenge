package io

// DefaultMaxBufferSize is the default maximum number of bytes that a ByteStream
// will read from the underlying buffer at once.
const DefaultMaxBufferSize = 1024

// ByteStream enables clients to stream bytes from an OutputBuffer.
type ByteStream struct {
	buffer      OutputBuffer
	channel     chan []byte
	maxReadSize int
}

// NewByteStream creates and returns a new ByteStream associated with the
// given buffer with a read buffer size of DefaultMaxBufferSize.
func NewByteStream(buffer OutputBuffer) *ByteStream {
	return NewByteStreamDetailed(buffer, DefaultMaxBufferSize)
}

// NewByteStream creates and returns a new ByteStream associated with the
// given buffer with a read buffer size of the given maxReadSize.
func NewByteStreamDetailed(buffer OutputBuffer, maxReadSize int) *ByteStream {
	return &ByteStream{
		channel:     make(chan []byte),
		buffer:      buffer,
		maxReadSize: maxReadSize,
	}
}
