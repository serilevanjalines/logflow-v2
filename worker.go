package main

import (
	"context"
	"time"
	"log"

	"github.com/jackc/pgx/v5"
)

func logWorker(conn *pgx.Conn) {

	for {
		incomingLog := <-logQueue
		success := false

		log.Println("Worker received log:", incomingLog.ID)

		for attempt := 1; attempt <= 3; attempt++ {

			_, err := conn.Exec(
				context.Background(),
				`INSERT INTO logs (id, timestamp, service, level, message)
         VALUES ($1, $2, $3, $4, $5)`,
				incomingLog.ID,
				incomingLog.Timestamp,
				incomingLog.Service,
				incomingLog.Level,
				incomingLog.Message,
			)

			if err == nil {
				success = true
				log.Println("Log inserted into database:", incomingLog.ID)
				break
			}

			log.Printf(
				"Insert failed for %s (attempt %d/3): %v",
				incomingLog.ID,
				attempt,
				err,
			)

			if attempt < 3 {
				time.Sleep(time.Second)
			}
		}

		if !success {
			log.Println("Failed to insert log after 3 attempts, moving to failed queue:", incomingLog.ID)
			failedQueue <- incomingLog
		}

	}
}
