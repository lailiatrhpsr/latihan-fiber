package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"latihan-fiber/pert3/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() *pgxpool.Pool {
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	user := config.GetEnv("DB_USER", "postgres")
	pass := config.GetEnv("DB_PASSWORD", "")
	name := config.GetEnv("DB_NAME", "praktikum_backend")
	ssl := config.GetEnv("DB_SSLMODE", "disable")
	maxConns := config.GetEnvInt("DB_MAX_CONNS", 10)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=%d",
		user, pass, host, port, name, ssl, maxConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Gagal membaca konfigurasi database pool: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Gagal inisialisasi connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Gagal terhubung ke database PostgreSQL: %v", err)
	}

	return pool
}