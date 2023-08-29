package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
	"strconv"
	"sync"
)

var (
	rClientInstance *redis.Client
	onceRC          sync.Once
)

func InitializeRedisClient(config *config.Config) *redis.Client {
	onceRC.Do(func() {
		address := config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port)

		rClientInstance = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		})

		_, err := rClientInstance.Ping(context.Background()).Result()
		if err != nil {
			panic(fmt.Sprintf("Could not connect to Redis: %v", err))
		}
	})
	return rClientInstance
}
