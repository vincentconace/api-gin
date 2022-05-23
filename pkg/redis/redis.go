package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	Set(ctx context.Context, key string, data interface{}, timeDuration time.Duration) (string, error)
	Get(ctx context.Context, key string) (interface{}, error)
}

type redisClient struct {
	rd *redis.Client
}

var ctx = context.Background()

func NewRedisClient() Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &redisClient{rd: rdb}
}

func (r *redisClient) Set(ctx context.Context, key string, data interface{}, timeDuration time.Duration) (string, error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	datastring, err := r.rd.Set(ctx, key, string(dataByte), timeDuration*time.Hour).Result()
	if err != nil {
		fmt.Println(err)
	}

	return datastring, err
}

func (r *redisClient) Get(ctx context.Context, key string) (interface{}, error) {
	var data interface{}
	dataByte, err := r.rd.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(dataByte), &data)
	if err != nil {
		fmt.Println(err)
	}

	return data, err
}
