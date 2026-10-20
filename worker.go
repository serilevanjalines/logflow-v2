package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func processLog(incomingLog Log, id int, pool *pgxpool.Pool, redisClient *redis.Client) {
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
		logData, err := json.Marshal(incomingLog)
		if err != nil {
			log.Printf("Failed to encode failed log %s: %v", incomingLog.ID, err)
			return
		}

		err = redisClient.LPush(
			context.Background(),
			"failedQueue",
			logData,
		).Err()

		if err != nil {
			log.Printf("Failed to store log %s in failedQueue: %v", incomingLog.ID, err)
		}
	}
}

func logWorker(id int, pool *pgxpool.Pool, redisClient *redis.Client, shutdown <-chan struct{}, wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case <-shutdown:
			log.Printf("Worker %d shutting down", id)

			// Drain anything remaining in Redis.
			for {
				result, err := redisClient.RPop(
					context.Background(),
					"logQueue",
				).Result()

				if err == redis.Nil {
					return
				}

				if err != nil {
					log.Printf("Worker %d Redis error: %v", id, err)
					return
				}

				var incomingLog Log

				err = json.Unmarshal([]byte(result), &incomingLog)
				if err != nil {
					log.Printf("Worker %d failed to decode log: %v", id, err)
					continue
				}

				processLog(incomingLog, id, pool, redisClient)
			}

		default:
			result, err := redisClient.BRPop(
				context.Background(),
				time.Second,
				"logQueue",
			).Result()

			if err == redis.Nil {
				continue
			}

			if err != nil {
				log.Printf("Worker %d Redis error: %v", id, err)
				continue
			}

			var incomingLog Log

			err = json.Unmarshal([]byte(result[1]), &incomingLog)
			if err != nil {
				log.Printf("Worker %d failed to decode log: %v", id, err)
				continue
			}

			processLog(incomingLog, id, pool, redisClient)
		}
	}
}
