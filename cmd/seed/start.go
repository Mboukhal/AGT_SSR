package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Mboukhal/SvGoPg/cmd/seed/db"
	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("Starting the seeding process...")
	_ = godotenv.Load()
}

func dbSeed() {

	ctx := context.Background()

	conn, err := db.Setup(ctx)
	if err != nil {
		log.Fatal("Failed to set up database connection:", err)
	}
	defer conn.Close(context.Background())

	db.Test(conn)

	fmt.Println("Seeding completed successfully")
}

func main() {
	dbSeed()

}
