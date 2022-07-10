package model

import (
	"context"
	redis "github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
	"time"
)

var (
	redisDB *redis.Client
)

func InitRedis() error {
	redisDB = redis.NewClient(&redis.Options{
		Addr:        viper.GetString("redis.addr"),
		Password:    viper.GetString("redis.password"),
		DB:          viper.GetInt("redis.db"),
		PoolTimeout: time.Duration(viper.GetInt("redis.pool_timeout")),
		PoolSize:    viper.GetInt("redis.pool_size"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisDB.Set(ctx, "test", "test", 10*time.Second)
	_, err := redisDB.Get(ctx, "test").Result()

	return err
}

func ZAddWithContext(id string, msg DialogMessage) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := redisDB.ZAdd(ctx, id, redis.Z{Score: float64(msg.CreatedAt.Unix()), Member: msg}).Result()
	return err

}

func ZGetAllWithContext(id string) ([]DialogMessage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := redisDB.ZCard(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	results, err := redisDB.ZRangeWithScores(ctx, id, 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]DialogMessage, 0)
	for _, msg := range results {
		messages = append(messages, msg.Member.(DialogMessage))
	}
	return messages, nil
}
