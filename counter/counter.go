package counter

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	defaultBufferSize    = 1000
	defaultFlushInterval = 5 * time.Second
	defaultFlushTimeout  = 3 * time.Second
)

var (
	ErrInvalidBufferSize    = errors.New("counter buffer size must be greater than zero")
	ErrInvalidFlushInterval = errors.New("counter flush interval must be greater than zero")
	ErrInvalidFlushTimeout  = errors.New("counter flush timeout must be greater than zero")
)

// FlushFunc persists the increments accumulated for each key.
type FlushFunc[K comparable] func(ctx context.Context, increments map[K]int64) error

// Option configures a Counter.
type Option func(*options)

type options struct {
	bufferSize    int
	flushInterval time.Duration
	flushTimeout  time.Duration
	onError       func(error)
}

type increment[K comparable] struct {
	key   K
	delta int64
}

// Counter asynchronously aggregates increments by key and periodically flushes them.
type Counter[K comparable] struct {
	flusher FlushFunc[K]
	options options
	events  chan increment[K]
	stop    chan struct{}
	done    chan struct{}

	lifecycleMu sync.RWMutex
	closed      bool
	closeErrMu  sync.Mutex
	closeErr    error
}

// WithBufferSize configures the number of increments waiting to be aggregated.
func WithBufferSize(size int) Option {
	return func(options *options) {
		options.bufferSize = size
	}
}

// WithFlushInterval configures how often accumulated increments are flushed.
func WithFlushInterval(interval time.Duration) Option {
	return func(options *options) {
		options.flushInterval = interval
	}
}

// WithFlushTimeout configures the timeout for each flush operation.
func WithFlushTimeout(timeout time.Duration) Option {
	return func(options *options) {
		options.flushTimeout = timeout
	}
}

// WithErrorHandler configures a callback for periodic flush errors.
func WithErrorHandler(handler func(error)) Option {
	return func(options *options) {
		options.onError = handler
	}
}

// New creates and starts an asynchronous counter.
func New[K comparable](flusher FlushFunc[K], opts ...Option) (*Counter[K], error) {
	if flusher == nil {
		return nil, errors.New("counter flusher is required")
	}
	config := options{
		bufferSize:    defaultBufferSize,
		flushInterval: defaultFlushInterval,
		flushTimeout:  defaultFlushTimeout,
	}
	for _, option := range opts {
		option(&config)
	}
	if config.bufferSize <= 0 {
		return nil, ErrInvalidBufferSize
	}
	if config.flushInterval <= 0 {
		return nil, ErrInvalidFlushInterval
	}
	if config.flushTimeout <= 0 {
		return nil, ErrInvalidFlushTimeout
	}
	counter := &Counter[K]{
		flusher: flusher,
		options: config,
		events:  make(chan increment[K], config.bufferSize),
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
	go counter.process()
	return counter, nil
}

// Increment adds one to the pending count for key without blocking when the buffer is full.
func (c *Counter[K]) Increment(key K) bool {
	return c.Add(key, 1)
}

// Add adds delta to the pending count for key without blocking when the buffer is full.
func (c *Counter[K]) Add(key K, delta int64) bool {
	c.lifecycleMu.RLock()
	defer c.lifecycleMu.RUnlock()
	if c.closed {
		return false
	}
	select {
	case c.events <- increment[K]{key: key, delta: delta}:
		return true
	default:
		return false
	}
}

// Close stops accepting increments, flushes accepted increments, and waits for completion.
func (c *Counter[K]) Close(ctx context.Context) error {
	c.lifecycleMu.Lock()
	if !c.closed {
		c.closed = true
		close(c.stop)
	}
	c.lifecycleMu.Unlock()

	select {
	case <-c.done:
		c.closeErrMu.Lock()
		defer c.closeErrMu.Unlock()
		return c.closeErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

// process owns the aggregation map and serializes flush operations.
func (c *Counter[K]) process() {
	ticker := time.NewTicker(c.options.flushInterval)
	defer ticker.Stop()
	defer close(c.done)

	counts := make(map[K]int64)
	for {
		select {
		case event := <-c.events:
			counts[event.key] += event.delta
		case <-ticker.C:
			counts = c.flushAndRetainOnError(counts)
		case <-c.stop:
			counts = c.drain(counts)
			err := c.flush(counts)
			c.closeErrMu.Lock()
			c.closeErr = err
			c.closeErrMu.Unlock()
			return
		}
	}
}

// drain aggregates all increments accepted before Close acquired the lifecycle lock.
func (c *Counter[K]) drain(counts map[K]int64) map[K]int64 {
	for {
		select {
		case event := <-c.events:
			counts[event.key] += event.delta
		default:
			return counts
		}
	}
}

// flushAndRetainOnError clears a successful batch and retains a failed batch for retry.
func (c *Counter[K]) flushAndRetainOnError(counts map[K]int64) map[K]int64 {
	if err := c.flush(counts); err != nil {
		if c.options.onError != nil {
			c.options.onError(err)
		}
		return counts
	}
	return make(map[K]int64)
}

// flush persists a non-empty batch with the configured timeout.
func (c *Counter[K]) flush(counts map[K]int64) error {
	if len(counts) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.options.flushTimeout)
	defer cancel()
	return c.flusher(ctx, counts)
}
