package repository

import (
	"context"
	"fmt"
	"motico-api/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnectionPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.Database.MaxConnections)
	poolConfig.MaxIdleConns = int32(cfg.Database.MaxIdleConns)

	connMaxLifetime, err := cfg.Database.GetConnMaxLifetime()
	if err != nil {
		return nil, fmt.Errorf("error parsing conn_max_lifetime: %w", err)
	}
	poolConfig.MaxConnLifetime = connMaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return pool, nil
}

