package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func connectDB() *pgx.Conn {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), databaseURL)

	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	log.Println("Connected to database")

	return conn

}
