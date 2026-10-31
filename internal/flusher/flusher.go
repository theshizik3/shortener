package flusher

import (
	"context"
	"strings"
	"time"

	"github.com/theshizik3/shortener/internal/repository"
)

type Flusher struct {
	redisRepo    *repository.RedisRepo
	postgresRepo *repository.PostgresRepo
	t            time.Duration
}

func NewFlusher(t time.Duration, redisRepo *repository.RedisRepo, postgresRepo *repository.PostgresRepo) *Flusher {
	return &Flusher{t: t, redisRepo: redisRepo, postgresRepo: postgresRepo}
}

func (f *Flusher) flush(ctx context.Context) error {
	iter := f.redisRepo.ScanClicks(ctx)
	for iter.Next(ctx) {
		key := strings.TrimPrefix(iter.Val(), repository.ClicksKeyPrefix)
		clicks, err := f.redisRepo.GetDelClicks(ctx, key)
		if err != nil {
			return err
		}
		err = f.postgresRepo.IncrClicks(ctx, key, clicks)
	}
	return nil
}

func (f *Flusher) Run(ctx context.Context) {
	ticker := time.NewTicker(f.t)
	defer ticker.Stop()

	for range ticker.C {
		f.flush(ctx)
	}
}
