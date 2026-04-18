package settings

import (
	"context"
	"fmt"
	"os"

	sqlc "github.com/Mboukhal/SvGoPg/internal/db"
	"github.com/jackc/pgx/v5"
)

func Setup(ctx context.Context) (*sqlc.Queries, *pgx.Conn, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	} else {
		fmt.Println("Using DATABASE_URL:", connStr)
	}

	// connect
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the database: %w", err)
	} else {
		fmt.Println("Successfully connected to the database")
	}

	// Create queries instance
	q := sqlc.New(conn)
	if q == nil {
		return nil, nil, fmt.Errorf("failed to create queries instance")
	}
	return q, conn, nil
}
