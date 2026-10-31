package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const ClicksKeyPrefix = "clicks:"

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(cfg *redis.Options) *RedisRepo {
	client := redis.NewClient(cfg)
	return &RedisRepo{client: client}
}

func (r *RedisRepo) GetUrl(ctx context.Context, code string) (string, error) {
	return r.client.Get(ctx, code).Result()
}

func (r *RedisRepo) SetUrl(ctx context.Context, code, url string) error {
	return r.client.Set(ctx, code, url, 1*time.Hour).Err()
}

func (r *RedisRepo) IncrClick(ctx context.Context, code string) error {
	return r.client.Incr(ctx, ClicksKeyPrefix+code).Err()
}

func (r *RedisRepo) GetClicks(ctx context.Context, code string) (int, error) {
	clicks, err := r.client.Get(ctx, ClicksKeyPrefix+code).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return clicks, nil
}

func (r *RedisRepo) GetDelClicks(ctx context.Context, code string) (int, error) {
	clicks, err := r.client.GetDel(ctx, ClicksKeyPrefix+code).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return clicks, nil
}

func (r *RedisRepo) ScanClicks(ctx context.Context) *redis.ScanIterator {
	return r.client.Scan(ctx, 0, ClicksKeyPrefix+"*", 100).Iterator()
}
