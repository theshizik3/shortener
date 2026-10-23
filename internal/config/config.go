package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis    *redis.Options
	Postgres *pgxpool.Config
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	r_addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	r := &redis.Options{
		Addr:     r_addr,
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	pg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(err)
	}

	maxConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_CONNS"))
	if err != nil {
		panic(err)
	}
	minConns, err := strconv.Atoi(os.Getenv("POSTGRES_MIN_CONNS"))
	if err != nil {
		panic(err)
	}

	pg.MaxConns = int32(maxConns)
	pg.MinConns = int32(minConns)

	return &Config{Redis: r, Postgres: pg}
}
