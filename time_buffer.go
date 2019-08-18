package main

import (
	"time"
)

// TimeBuffer is a string buffer to store strings until the buffer is flushed.
//
// It is flushed with size-base limit or time-based limit.
type TimeBuffer struct {
	buffer []string
	closed bool

	wch chan string
	ch  chan []string
}

// NewTimeBuffer returns a new TimeBuffer.  The buffer is flushed with the interval or the size.
func NewTimeBuffer(interval time.Duration, size int) *TimeBuffer {
	b := &TimeBuffer{
		ch:  make(chan []string),
		wch: make(chan string),
	}

	go func() {
		flush := func() {
			if len(b.buffer) == 0 {
				return
			}
			buf := b.buffer
			b.buffer = nil
			b.ch <- buf
		}
		for {
			select {
			case <-time.After(interval):
				flush()

			case s := <-b.wch:
				if b.closed {
					flush()
					close(b.ch)
					return
				}

				b.buffer = append(b.buffer, s)
				if len(b.buffer) >= size {
					flush()
				}
			}
		}
	}()
	return b
}

// Flushed returns the flush channel
func (b *TimeBuffer) Flushed() <-chan []string {
	return b.ch
}

// Write append a string to the buffer
func (b *TimeBuffer) Write(s string) {
	b.wch <- s
}

// Close flushes current buffer and close buffer
func (b *TimeBuffer) Close() {
	b.closed = true
	close(b.wch)
}
