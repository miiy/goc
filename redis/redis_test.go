package redis

import (
	"context"
	"testing"
)

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
