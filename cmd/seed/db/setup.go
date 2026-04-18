package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func Setup(ctx context.Context) (*pgx.Conn, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	} else {
		fmt.Println("Using DATABASE_URL:", connStr)
	}
	// connect
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
	// defer conn.Close(context.Background())
	return conn, nil
}

func Test(conn *pgx.Conn) {

	// test query
	var now time.Time
	err := conn.QueryRow(context.Background(), "SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB time:", now.Format(time.RFC3339Nano))
}
