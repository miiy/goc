package redis

import (
	"errors"
	"testing"
	"time"

	upstream "github.com/redis/go-redis/v9"
)

func TestIsNil(t *testing.T) {
	if !IsNil(upstream.Nil) {
		t.Fatal("redis.Nil was not recognized")
	}
	if IsNil(errors.New("other error")) {
		t.Fatal("non-nil Redis error was recognized as redis.Nil")
	}
}

func TestNewRedisDefaults(t *testing.T) {
	rdb, err := NewRedis(testOptions())
	client := requireClient(t, rdb, err)
	options := client.Options()

	upstreamClient := upstream.NewClient(&upstream.Options{Addr: "127.0.0.1:6379"})
	t.Cleanup(func() { _ = upstreamClient.Close() })
	upstreamOptions := upstreamClient.Options()

	assertEqual(t, "PoolSize", options.PoolSize, upstreamOptions.PoolSize)
	assertEqual(t, "MinIdleConns", options.MinIdleConns, upstreamOptions.MinIdleConns)
	assertEqual(t, "MaxIdleConns", options.MaxIdleConns, upstreamOptions.MaxIdleConns)
	assertEqual(t, "DialTimeout", options.DialTimeout, upstreamOptions.DialTimeout)
	assertEqual(t, "ReadTimeout", options.ReadTimeout, upstreamOptions.ReadTimeout)
	assertEqual(t, "WriteTimeout", options.WriteTimeout, upstreamOptions.WriteTimeout)
}

func TestNewRedisOptions(t *testing.T) {
	rdb, err := NewRedis(testOptions(),
		WithPoolSize(40),
		WithMinIdleConns(8),
		WithMaxIdleConns(16),
		WithDialTimeout(3*time.Second),
		WithReadTimeout(4*time.Second),
		WithWriteTimeout(5*time.Second),
	)
	client := requireClient(t, rdb, err)
	options := client.Options()

	assertEqual(t, "PoolSize", options.PoolSize, 40)
	assertEqual(t, "MinIdleConns", options.MinIdleConns, 8)
	assertEqual(t, "MaxIdleConns", options.MaxIdleConns, 16)
	assertEqual(t, "DialTimeout", options.DialTimeout, 3*time.Second)
	assertEqual(t, "ReadTimeout", options.ReadTimeout, 4*time.Second)
	assertEqual(t, "WriteTimeout", options.WriteTimeout, 5*time.Second)
}

func testOptions() *Options {
	return &Options{
		Addrs: []string{"127.0.0.1:6379"},
	}
}

func requireClient(t *testing.T, rdb UniversalClient, err error) *upstream.Client {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = rdb.Close() })

	client, ok := rdb.(*upstream.Client)
	if !ok {
		t.Fatalf("NewRedis() returned %T, want *redis.Client", rdb)
	}
	return client
}

func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
