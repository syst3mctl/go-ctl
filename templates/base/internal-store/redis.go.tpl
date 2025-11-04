package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func (rr *RedisRepository) Store(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	result := rr.client.Set(ctx, key, value, duration)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (rr *RedisRepository) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := rr.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (rr *RedisRepository) IsExist(ctx context.Context, key string) (int64, error) {
	result, err := rr.client.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (rr *RedisRepository) Delete(ctx context.Context, key string) error {
	err := rr.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rr *RedisRepository) SetHashValue(ctx context.Context, key string, value map[string]interface{}) error {
	err := rr.client.HSet(ctx, key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rr *RedisRepository) GetHMValue(ctx context.Context, key, field string) (string, error) {
	val, err := rr.client.HGet(ctx, key, field).Result()
	if err != nil {
		return "", nil
	}

	return val, nil
}

func NewRedisRepository(c *redis.Client) RedisRepository {
	return RedisRepository{
		client: c,
	}
}
