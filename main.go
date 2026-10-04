package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Log struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type LogResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {

	conn := connectDB()

	defer conn.Close(context.Background())

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		logHandler(w, r, conn)
	})

	log.Println("LogFlow listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func logHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var incomingLog Log

	err := json.NewDecoder(r.Body).Decode(&incomingLog)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if incomingLog.Service == "" || incomingLog.Level == "" || incomingLog.Message == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	incomingLog.Message = scrubPII(incomingLog.Message)
	incomingLog.ID = uuid.New().String()
	incomingLog.Timestamp = time.Now()

	_, err = conn.Exec(
		context.Background(),
		`INSERT INTO logs (id, timestamp, service, level, message)
	 VALUES ($1, $2, $3, $4, $5)`,
		incomingLog.ID,
		incomingLog.Timestamp,
		incomingLog.Service,
		incomingLog.Level,
		incomingLog.Message,
	)

	if err != nil {
		http.Error(w, "Failed to store log", http.StatusInternalServerError)
		return
	}

	log.Printf(
		"id=%s timestamp=%s service=%s level=%s message=%s",
		incomingLog.ID,
		incomingLog.Timestamp,
		incomingLog.Service,
		incomingLog.Level,
		incomingLog.Message,
	)

	var response LogResponse

	response.Status = "accepted"
	response.Message = "Log received"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
