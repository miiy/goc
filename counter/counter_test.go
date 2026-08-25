package counter

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

// TestCounterAggregates verifies that increments for the same key are combined.
func TestCounterAggregates(t *testing.T) {
	var got map[int64]int64
	counter, err := New(func(_ context.Context, increments map[int64]int64) error {
		got = clone(increments)
		return nil
	}, WithFlushInterval(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if !counter.Increment(1) || !counter.Increment(1) || !counter.Add(2, 3) {
		t.Fatal("counter unexpectedly rejected an increment")
	}
	if err := counter.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := map[int64]int64{1: 2, 2: 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("increments = %v, want %v", got, want)
	}
}

// TestCounterRetriesFailedFlush verifies that a failed periodic batch is retained.
func TestCounterRetriesFailedFlush(t *testing.T) {
	firstAttempt := make(chan struct{})
	var mu sync.Mutex
	var attempts int
	var got map[int64]int64
	counter, err := New(func(_ context.Context, increments map[int64]int64) error {
		mu.Lock()
		defer mu.Unlock()
		attempts++
		if attempts == 1 {
			close(firstAttempt)
			return errors.New("temporary failure")
		}
		got = clone(increments)
		return nil
	}, WithFlushInterval(5*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if !counter.Increment(7) {
		t.Fatal("counter unexpectedly rejected an increment")
	}
	select {
	case <-firstAttempt:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for the first flush")
	}
	if err := counter.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, map[int64]int64{7: 1}) {
		t.Fatalf("increments = %v, want map[7:1]", got)
	}
}

// TestCounterRejectsIncrementAfterClose verifies the counter lifecycle boundary.
func TestCounterRejectsIncrementAfterClose(t *testing.T) {
	counter, err := New(func(context.Context, map[int64]int64) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := counter.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if counter.Increment(1) {
		t.Fatal("increment was accepted after close")
	}
	if err := counter.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

// TestCounterValidatesOptions verifies invalid resource and timing options are rejected.
func TestCounterValidatesOptions(t *testing.T) {
	tests := []struct {
		name   string
		option Option
		want   error
	}{
		{name: "buffer size", option: WithBufferSize(0), want: ErrInvalidBufferSize},
		{name: "flush interval", option: WithFlushInterval(0), want: ErrInvalidFlushInterval},
		{name: "flush timeout", option: WithFlushTimeout(0), want: ErrInvalidFlushTimeout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(func(context.Context, map[int64]int64) error {
				return nil
			}, tt.option)
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

// clone copies increments retained by a test callback.
func clone[K comparable](increments map[K]int64) map[K]int64 {
	result := make(map[K]int64, len(increments))
	for key, count := range increments {
		result[key] = count
	}
	return result
}
