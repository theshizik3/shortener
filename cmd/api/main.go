package main

import (
	"context"

	"github.com/theshizik3/shortener/internal/config"
	"github.com/theshizik3/shortener/internal/flusher"
	"github.com/theshizik3/shortener/internal/httpapi"
	"github.com/theshizik3/shortener/internal/repository"
	"github.com/theshizik3/shortener/internal/shortener"
)

func main() {
	cfg := config.LoadConfig()

	redis_repo := repository.NewRedisRepo(cfg.Redis)
	postgres_repo, err := repository.NewPostgresRepo(context.Background(), cfg.Postgres)

	if err != nil {
		panic(err)
	}

	service := shortener.NewService(redis_repo, postgres_repo, cfg.CodeLen)
	fl := flusher.NewFlusher(cfg.FlushInterval, redis_repo, postgres_repo)

	ctrl := httpapi.NewController(service)

	r := httpapi.SetupRouter(ctrl)

	go fl.Run(context.Background())

	r.Run(":8080")
}
