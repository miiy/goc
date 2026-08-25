package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

type UniversalClient = redis.UniversalClient

// IsNil reports whether a Redis command failed because its key does not exist.
func IsNil(err error) bool { return errors.Is(err, redis.Nil) }

type Options struct {
	Addrs    []string `yaml:"addrs"`
	DB       int      `yaml:"database"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

func NewRedis(o *Options, opts ...Option) (redis.UniversalClient, error) {
	clientOptions := &redis.UniversalOptions{
		Addrs:    o.Addrs,
		Username: o.Username,
		Password: o.Password,
		DB:       o.DB,
	}
	for _, opt := range opts {
		opt.apply(clientOptions)
	}

	return redis.NewUniversalClient(clientOptions), nil
}
