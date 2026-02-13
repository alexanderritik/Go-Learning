package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(addr string) *RateLimiter{
	return &RateLimiter{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error){

	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.client.Expire(ctx, key, window)
	}

	if int(count) > limit {
		return false, nil
	}

	return  true, nil
}