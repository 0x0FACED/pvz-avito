package database

import (
	"context"
	"fmt"

	"github.com/0x0FACED/pvz-avito/internal/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPool creates new *pgxpool.Pool instance and returns it.
// If cant get pgxpool.Config from config.DatabaseConfig - returns err.
// If cant pool.Ping() - return err.
func ConnectPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	pgxpoolConfig, err := pgxpoolConfig(cfg)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxpoolConfig)
	if err != nil {
		// TODO: handle err
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		// TODO: handle err
		return nil, err
	}

	return pool, nil
}

func pgxpoolConfig(cfg config.DatabaseConfig) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		// TODO: handle err
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MaxConns = cfg.MaxOpenConns
	config.MinConns = cfg.MaxIdleConns
	config.MaxConnLifetime = cfg.ConnMaxLifetime
	config.MaxConnIdleTime = cfg.ConnMaxIdleTime
	config.HealthCheckPeriod = cfg.HealthCheckPeriod

	config.ConnConfig.ConnectTimeout = cfg.ConnectionTimeout

	return config, nil
}
