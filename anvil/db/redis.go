package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
	"sync"
)

type Keystore struct {
	rc  RedisClient // Interface for abstraction
	ctx context.Context
}

type RedisClient interface {
	Ping(ctx context.Context) error
	Close() error // Expose Close method for cleanup
}

type RedisInstance struct {
	client *redis.Client
}

func (ri *RedisInstance) Ping(ctx context.Context) error {
	return ri.client.Ping(ctx).Err()
}

func (ri *RedisInstance) Close() error {
	return ri.client.Close()
}

var redisOnce sync.Once
var redisInstance RedisClient

func GetRedisInstance() RedisClient {
	redisOnce.Do(func() {
		cfg := config.GetConfigInstance(nil)
		address := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
		password := cfg.Redis.Password

		client := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       cfg.Redis.DB,
		})

		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err != nil {
			panic(fmt.Sprintf("Could not connect to Redis: %v", err))
		}

		redisInstance = &RedisInstance{client: client}
	})
	return redisInstance
}

func NewKeystore(ctx context.Context) *Keystore {
	return &Keystore{
		rc:  GetRedisInstance(),
		ctx: ctx,
	}
}

func (ks *Keystore) RedisClient() RedisClient {
	return ks.rc
}

func (ks *Keystore) Shutdown() error {
	redisClient := ks.RedisClient()
	if err := redisClient.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %v", err)
	}
	return nil
}
