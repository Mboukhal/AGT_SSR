package main

import (
	"context"
	"fmt"

	db_seed "github.com/Mboukhal/SvGoPg/cmd/seed/db"
	"github.com/Mboukhal/SvGoPg/cmd/settings"
	"github.com/joho/godotenv"
)

func init() {
	fmt.Println("Starting the seeding process...")
	_ = godotenv.Load()
}

func dbSeed() {

	ctx := context.Background()

	// conn, err := db_seed.Setup(ctx)
	// if err != nil {
	// 	log.Fatal("Failed to set up database connection:", err)
	// }
	// defer conn.Close(context.Background())

	// // Create queries instance
	// q := sqlc.New(conn)

	q, conn, err := settings.Setup(ctx)
	if err != nil {
		fmt.Println("Failed to set up database connection:", err)
		return
	}
	defer conn.Close()
	if err = db_seed.LoadUsers(q, ctx); err != nil {
		println("Err:", err.Error())
	}
	// db_seed.Test(conn)

	fmt.Println("Seeding completed successfully")
}

func main() {
	dbSeed()

}
