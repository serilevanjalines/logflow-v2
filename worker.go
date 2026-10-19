package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func logWorker(id int, pool *pgxpool.Pool){

	for {
		incomingLog := <-logQueue

		success := false

		log.Printf("Worker %d received log: %s", id, incomingLog.ID)

		for attempt := 1; attempt <= 3; attempt++ {

			_, err := pool.Exec(
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
				log.Printf("Worker %d inserted log: %s", id, incomingLog.ID)
				break
			}

			log.Printf(
				"Insert failed for %s (attempt %d/3): %v and the worker id is %d",
				incomingLog.ID,
				attempt,
				err,
				id,
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
