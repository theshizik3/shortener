package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

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
	return r.client.Incr(ctx, "clicks:"+code).Err()
}

func (r *RedisRepo) GetClicks(ctx context.Context, code string) (int, error) {
	return r.client.Get(ctx, "clicks:"+code).Int()
}

func (r *RedisRepo) GetDelClicks(ctx context.Context, code string) (int, error) {
	return r.client.GetDel(ctx, "clicks:"+code).Int()
}

func (r *RedisRepo) ScanClicks(ctx context.Context) *redis.ScanIterator {
	return r.client.Scan(ctx, 0, "clicks:*", 100).Iterator()
}
