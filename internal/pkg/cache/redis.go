package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func Init(host, port, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("error client ping: %w", err)
	}
	return &Redis{client: client}, nil
}

type Redis struct {
	client *redis.Client
}

func (r *Redis) Set(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Set(ctx, key, key, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Get(ctx, key).Result()
	return strings.TrimSpace(val) != "", err
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
