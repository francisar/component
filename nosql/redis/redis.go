package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisClient struct {
	client *redis.Client
}

func (r *redisClient)Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd  {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *redisClient)Get(ctx context.Context, key string) *redis.StringCmd  {
	return r.client.Get(ctx, key)
}


func (r *redisClient)SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.client.SetNX(ctx, key, value, expiration)
}


func (r *redisClient)Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(ctx, key, expiration)
}

func (r *redisClient)Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Exists(ctx, keys...)
}

func (r *redisClient)SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.SetEX(ctx, key, value, expiration)
}

func (r *redisClient)TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.TTL(ctx, key)
}

func (r *redisClient)Del(ctx context.Context, keys ...string) *redis.IntCmd  {
	return r.client.Del(ctx, keys...)
}