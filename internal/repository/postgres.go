package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, cfg *pgxpool.Config) (*PostgresRepo, error) {
	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	if err != nil {
		return nil, err
	}

	r := &PostgresRepo{pool: pool}

	err = r.Init(context.Background())

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PostgresRepo) Init(ctx context.Context) error {
	_, err := r.pool.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS Links (
			id SERIAL PRIMARY KEY,
			code TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL UNIQUE,
			clicks BIGINT NOT NULL DEFAULT 0,
		)`,
	)
	return err
}

func (r *PostgresRepo) GetUrl(ctx context.Context, code string) (string, error) {
	var url string
	row := r.pool.QueryRow(ctx, "SELECT url FROM Links WHERE code = $1", code)
	err := row.Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (r *PostgresRepo) SetUrl(ctx context.Context, code, url string) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO Links (code, url) VALUES ($1, $2)", code, url)
	return err
}

func (r *PostgresRepo) DelCode(ctx context.Context, code string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM Links WHERE code = $1", code)
	return err
}

func (r *PostgresRepo) IncrClick(ctx context.Context, code string, count int) error {
	_, err := r.pool.Exec(ctx, "UPDATE Links SET clicks = clicks + $1 WHERE code = $2", count, code)
	return err
}
