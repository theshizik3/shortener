package shortener

import (
	"context"
	"fmt"
	"strings"

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

func (s *Service) ValidateCode(code string) error {
	if len(code) != s.code_len {
		return fmt.Errorf("invalid code length: want %d, got %d", s.code_len, len(code))
	}
	for i := 0; i < len(code); i++ {
		if !strings.ContainsRune(alphabet, rune(code[i])) {
			return fmt.Errorf("invalid code: must be base62")
		}
	}
	return nil
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

func (s *Service) DeleteUrl(ctx context.Context, code string) error {
	return s.postgresRepo.DelCode(ctx, code)
}

func (s *Service) GetClicks(ctx context.Context, code string) (int, error) {
	pg_clicks, err := s.postgresRepo.GetClicks(ctx, code)
	if err != nil {
		return 0, err
	}
	redis_clicks, err := s.redisRepo.GetClicks(ctx, code)
	if err != nil {
		return 0, err
	}
	return pg_clicks + redis_clicks, nil
}
