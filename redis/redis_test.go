package redis

import (
	"context"
	"errors"
	"testing"

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

func TestNewRedis(t *testing.T) {
	rdb, err := NewRedis(&Options{
		Addrs:    []string{"127.0.0.1:6379"},
		DB:       0,
		Username: "",
		Password: "",
	})
	if err != nil {
		t.Error(err)
	}
	info, err := rdb.Info(context.Background()).Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
