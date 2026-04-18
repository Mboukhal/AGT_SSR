package main

import (
	"context"
	"fmt"
	"log"

	db_seed "github.com/Mboukhal/SvGoPg/cmd/seed/db"
	sqlc "github.com/Mboukhal/SvGoPg/internal/db"
	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("Starting the seeding process...")
	_ = godotenv.Load()
}

func dbSeed() {

	ctx := context.Background()

	conn, err := db_seed.Setup(ctx)
	if err != nil {
		log.Fatal("Failed to set up database connection:", err)
	}
	defer conn.Close(context.Background())

	// Create queries instance
	q := sqlc.New(conn)

	db_seed.LoadUsers(q, ctx)
	db_seed.Test(conn)

	fmt.Println("Seeding completed successfully")
}

func main() {
	dbSeed()

}
