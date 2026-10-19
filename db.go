package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectDB() *pgxpool.Pool {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	
	pool, err := pgxpool.New(context.Background(), databaseURL)

	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	log.Println("Connected to database")

	return pool

}
