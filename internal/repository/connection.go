package repository

import (
	"context"
	"fmt"
	"math"
	"motico-api/config"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnectionPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	user := url.UserPassword(cfg.Database.User, cfg.Database.Password)

	// Construir query parameters
	queryParams := fmt.Sprintf("sslmode=%s", cfg.Database.SSLMode)
	if cfg.Database.PoolMode != "" {
		queryParams += fmt.Sprintf("&pool_mode=%s", cfg.Database.PoolMode)
	}
	// Usar protocolo simple cuando se usa transaction pooling
	// El transaction pooler de Supabase no soporta prepared statements
	if cfg.Database.PoolMode == "transaction" {
		queryParams += "&prefer_simple_protocol=true"
	}

	connURL := &url.URL{
		Scheme:   "postgres",
		User:     user,
		Host:     fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port),
		Path:     "/" + cfg.Database.Name,
		RawQuery: queryParams,
	}

	connString := connURL.String()

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	if cfg.Database.MaxConnections < 0 || cfg.Database.MaxConnections > math.MaxInt32 {
		return nil, fmt.Errorf("max_connections value %d is out of valid range (0-%d)", cfg.Database.MaxConnections, math.MaxInt32)
	}
	// nolint:gosec // G115: Safe conversion - validated above to be within int32 range
	maxConns := int32(cfg.Database.MaxConnections)
	poolConfig.MaxConns = maxConns

	connMaxLifetime, err := cfg.Database.GetConnMaxLifetime()
	if err != nil {
		return nil, fmt.Errorf("error parsing conn_max_lifetime: %w", err)
	}
	poolConfig.MaxConnLifetime = connMaxLifetime

	// Usar modo de ejecución sin prepared statements cuando se usa transaction pooling
	// El transaction pooler de Supabase no soporta prepared statements
	if cfg.Database.PoolMode == "transaction" {
		poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec
		poolConfig.ConnConfig.StatementCacheCapacity = 0
	}

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
