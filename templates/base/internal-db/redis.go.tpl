package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisAddr, redisPassword, redisDB string) (*redis.Client, error) {
	ctx := context.Background()

	dbNum, err := strconv.Atoi(redisDB)
	if err != nil {
		dbNum = 0 // default database
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       dbNum,
	})

	status, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println(status)

	return client, nil
}
