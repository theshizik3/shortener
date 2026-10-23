package stortener

import (
	"context"

	"github.com/theshizik3/shortener/internal/repository"
)

type Service struct {
	redisRepo    *repository.RedisRepo
	postgresRepo *repository.PostgresRepo
	code_len     int
}

func NewService(redisRepo *repository.RedisRepo, postgresRepo *repository.PostgresRepo, code_len int) *Service {
	return &Service{
		redisRepo:    redisRepo,
		postgresRepo: postgresRepo,
		code_len:     code_len,
	}
}

func (s *Service) GetUrl(ctx context.Context, code string) (string, error) {
	url, err := s.redisRepo.GetUrl(ctx, code)
	if err == nil {
		s.redisRepo.IncrClick(ctx, code)
		return url, nil
	}
	url, err = s.postgresRepo.GetUrl(ctx, code)
	if err == nil {
		s.redisRepo.IncrClick(ctx, code)
	}
	return url, err
}

func (s *Service) SetUrl(ctx context.Context, url string) (string, error) {
	code, err := generateCode(s.code_len)
	if err != nil {
		return "", err
	}
	err = s.postgresRepo.SetUrl(ctx, code, url)
	if err != nil {
		return "", err
	}
	return code, nil
}
