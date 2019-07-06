package main

import (
	"context"
	"errors"
	"testing"
)

func TestWorker(t *testing.T) {
	myerr := errors.New("test error")

	w := NewWorker(context.Background())
	w.Start(func(ctx context.Context) error {
		<-ctx.Done()
		return myerr
	})

	err := w.Stop()
	if err != myerr {
		t.Errorf("%v != %v", err, myerr)
	}
}
