package settings

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	"github.com/Mboukhal/AGT_SSR/internal/db"
)

type contextKey string

const QueriesKey contextKey = "queries"

// WithQueries creates a middleware that injects queries into the request context
func WithQueries(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), QueriesKey, queries)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetQueries retrieves queries from the request context
func GetQueries(ctx context.Context) *db.Queries {
	queries, ok := ctx.Value(QueriesKey).(*db.Queries)
	if !ok {
		return nil
	}
	return queries
}

func OpenDB() (*sql.DB, error) {
	dbFile := os.Getenv("DATABASE_URL")
	if dbFile == "" {
		panic("DATABASE_URL is not set in .env")
	}

	// Extract directory from file path
	dbPath := os.Getenv("DATABASE_DIR")
	if dbPath == "" {
		panic("DATABASE_DIR is not set in .env")
	}

	// Ensure the directory for the database file exists
	if dbPath != "" {
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			createErr := os.MkdirAll(dbPath, os.ModePerm)
			if createErr != nil {
				panic("Failed to create database directory: " + createErr.Error())
			}
		}
	}

	// Note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-32000)"
	dbFile = "file:" + dbFile + pragmas

	// Open the database using the standard library
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	// log.Printf("Connected to database successfully")

	return db, nil
}
