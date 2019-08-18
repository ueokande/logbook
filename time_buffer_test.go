package main

import (
	"fmt"
	"testing"
	"time"
)

func testTimeBufferFlushBySize(t *testing.T) {
	b := NewTimeBuffer(1*time.Minute, 5)
	var counter [][]string

	done := make(chan struct{})
	go func() {
		for ss := range b.Flushed() {
			counter = append(counter, ss)
		}
		done <- struct{}{}
	}()
	for i := 0; i < 15; i++ {
		b.Write(fmt.Sprintf("%02d", i))
	}
	b.Close()

	<-done

	var i int
	if len(counter) != 3 {
		t.Fatalf("%d != %d", len(counter), 3)
	}
	for _, ss := range counter {
		if len(ss) != 5 {
			t.Fatalf("%d != %d", len(ss), 5)
		}
		for _, s := range ss {
			if s != fmt.Sprintf("%02d", i) {
				t.Fatalf("%s != %s", s, fmt.Sprintf("%02d", i))
			}
			i++
		}
	}
}

func testTimeBufferFlushByTime(t *testing.T) {
	b := NewTimeBuffer(100*time.Millisecond, 10)

	var counter [][]string

	done := make(chan struct{})
	go func() {
		for ss := range b.Flushed() {
			counter = append(counter, ss)
		}
		done <- struct{}{}
	}()
	for i := 0; i < 5; i++ {
		b.Write(fmt.Sprintf("%02d", i))
	}
	time.Sleep(188 * time.Millisecond)
	for i := 5; i < 10; i++ {
		b.Write(fmt.Sprintf("%02d", i))
	}
	time.Sleep(188 * time.Millisecond)

	b.Close()

	<-done

	var i int
	if len(counter) != 2 {
		t.Fatalf("%d != %d", len(counter), 2)
	}
	for _, ss := range counter {
		if len(ss) != 5 {
			t.Fatalf("%d != %d", len(ss), 5)
		}
		for _, s := range ss {
			if s != fmt.Sprintf("%02d", i) {
				t.Fatalf("%s != %s", s, fmt.Sprintf("%02d", i))
			}
			i++
		}
	}
}

func TestTimeBuffer(t *testing.T) {
	t.Run("flush by size", testTimeBuffer_flushBySize)
	t.Run("flush by time", testTimeBuffer_flushByTime)
}
