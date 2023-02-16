package redis

import (
	"github.com/go-redis/redis/v8"
	"strings"
)

func NewRedisClientFromOps(ops *Option) Client {

	switch ops.IsCluster {
	case true:
		redisOps := redis.Options{
			Password: ops.Password,
			Username: ops.UserName,
			DialTimeout: ops.Timeout,
			ReadTimeout: ops.Timeout,
			WriteTimeout: ops.Timeout,
			Network: "tcp",
			Addr:ops.Address,
		}
		client := redisClient{
			client: redis.NewClient(&redisOps),
		}
		return &client
	case false:
		redisOps := redis.ClusterOptions{
			Password: ops.Password,
			Username: ops.UserName,
			DialTimeout: ops.Timeout,
			ReadTimeout: ops.Timeout,
			WriteTimeout: ops.Timeout,
			Addrs:strings.Split(ops.Address, ","),
		}
		client := redisClusterClient{
			client: redis.NewClusterClient(&redisOps),
		}
		return &client
	}
	return nil
}

func NewRedisClient(client redis.Client) Client {
	newClient := redisClient{
		client: &client,
	}
	return &newClient
}