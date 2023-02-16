package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisClusterClient struct {
	client *redis.ClusterClient
}

func (r *redisClusterClient)Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd  {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *redisClusterClient)Get(ctx context.Context, key string) *redis.StringCmd  {
	return r.client.Get(ctx, key)
}


func (r *redisClusterClient)SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.client.SetNX(ctx, key, value, expiration)
}


func (r *redisClusterClient)Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(ctx, key, expiration)
}

func (r *redisClusterClient)Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Exists(ctx, keys...)
}

func (r *redisClusterClient)SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.SetEX(ctx, key, value, expiration)
}

func (r *redisClusterClient)TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.TTL(ctx, key)
}

func (r *redisClusterClient)Del(ctx context.Context, keys ...string) *redis.IntCmd  {
	return r.client.Del(ctx, keys...)
}