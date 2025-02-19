package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func OpenRedisClient(redisAddress string, password string, dbType int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: password,
		DB:       dbType,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SetCache(client *redis.Client, key string, value string, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

func GetCache(client *redis.Client, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func DeleteCache(client redis.Client, key string) (int64, error) {
	return client.Del(ctx, key).Result()
}
