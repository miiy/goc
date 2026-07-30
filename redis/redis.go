package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"math"
	"runtime"
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

func NewRedis(o *Options) (redis.UniversalClient, error) {
	// go-redis default pollSize
	pollSize := 10 * runtime.GOMAXPROCS(0)
	// Set min idle connections
	minIdleConns := int(math.Floor(float64(pollSize / 3)))

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        o.Addrs,
		Username:     o.Username,
		Password:     o.Password,
		DB:           o.DB,
		PoolSize:     pollSize,
		MinIdleConns: minIdleConns,
	})

	return client, nil
}
