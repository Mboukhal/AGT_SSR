package settings

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	sqlc "github.com/Mboukhal/AGT_SSR/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultDBMaxConns = int32(10)
	defaultDBMinConns = int32(0)
)

var defaultDBMaxConnLifetime = 30 * time.Minute

func buildPoolConfig() (*pgxpool.Config, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = defaultDBMaxConns
	poolConfig.MinConns = defaultDBMinConns
	poolConfig.MaxConnLifetime = defaultDBMaxConnLifetime

	if v := os.Getenv("DB_MAX_CONNS"); v != "" {
		maxConns, convErr := strconv.ParseInt(v, 10, 32)
		if convErr != nil || maxConns <= 0 {
			return nil, fmt.Errorf("invalid DB_MAX_CONNS value: %s", v)
		}
		poolConfig.MaxConns = int32(maxConns)
	}

	if v := os.Getenv("DB_MIN_CONNS"); v != "" {
		minConns, convErr := strconv.ParseInt(v, 10, 32)
		if convErr != nil || minConns < 0 {
			return nil, fmt.Errorf("invalid DB_MIN_CONNS value: %s", v)
		}
		poolConfig.MinConns = int32(minConns)
	}

	if v := os.Getenv("DB_MAX_CONN_LIFETIME"); v != "" {
		lifetime, convErr := time.ParseDuration(v)
		if convErr != nil {
			return nil, fmt.Errorf("invalid DB_MAX_CONN_LIFETIME value: %s", v)
		}
		poolConfig.MaxConnLifetime = lifetime
	}

	if poolConfig.MinConns > poolConfig.MaxConns {
		return nil, fmt.Errorf("DB_MIN_CONNS (%d) cannot be greater than DB_MAX_CONNS (%d)", poolConfig.MinConns, poolConfig.MaxConns)
	}

	return poolConfig, nil
}

func Setup(ctx context.Context) (*sqlc.Queries, *pgxpool.Pool, error) {
	poolConfig, err := buildPoolConfig()
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("DB pool config loaded: max_conns=%d min_conns=%d max_conn_lifetime=%s\n", poolConfig.MaxConns, poolConfig.MinConns, poolConfig.MaxConnLifetime)

	// connect using a pool to support concurrent requests without a single-connection bottleneck
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the database: %w", err)
	} else {
		fmt.Println("Successfully connected to the database")
	}

	if pingErr := pool.Ping(ctx); pingErr != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("failed to ping the database: %w", pingErr)
	}

	// Create queries instance
	q := sqlc.New(pool)
	if q == nil {
		pool.Close()
		return nil, nil, fmt.Errorf("failed to create queries instance")
	}
	return q, pool, nil
}
