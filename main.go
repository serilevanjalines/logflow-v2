package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

var logQueue = make(chan Log, 100)
var failedQueue = make(chan Log, 100)

func main() {

	conn := connectDB()

	defer conn.Close(context.Background())
	//when main eventually returns , deletes the connection to the database

	go logWorker(conn)
	//Start our worker goroutine and give it access to that DB connection

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		logHandler(w, r, conn)
	})

	http.HandleFunc("/logs/", func(w http.ResponseWriter, r *http.Request) {
		logHandler(w, r, conn)
	})

	log.Println("LogFlow listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func logHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {

	if strings.HasPrefix(r.URL.Path, "/logs/") && r.Method == http.MethodGet {
		getLogByID(w, r, conn)
		return
	}

	if r.Method == http.MethodGet {
		getLogs(w, r, conn)
		return
	}

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

	select {
	case logQueue <- incomingLog:

	default:
		http.Error(w, "Log queue is full", http.StatusServiceUnavailable)
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
	response.Message = "Log received for Processing"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getLogs(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {

	service := r.URL.Query().Get("service")
	level := r.URL.Query().Get("level")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	query := `
    SELECT id, timestamp, service, level, message
    FROM logs
`

	args := []any{}
	conditions := []string{}

	if service != "" {
		conditions = append(conditions, "service = $1")
		args = append(args, service)
	}

	if level != "" {
		conditions = append(conditions, fmt.Sprintf("level = $%d", len(args)+1))
		args = append(args, level)
	}

	if from != "" {
		fromTime, err := time.Parse(time.RFC3339, from)

		if err != nil {
			http.Error(w, "Invalid from timestamp", http.StatusBadRequest)
			return
		}

		conditions = append(
			conditions,
			fmt.Sprintf("timestamp >= $%d", len(args)+1),
		)

		args = append(args, fromTime)
	}

	if to != "" {
		toTime, err := time.Parse(time.RFC3339, to)

		if err != nil {
			http.Error(w, "Invalid to timestamp", http.StatusBadRequest)
			return
		}

		conditions = append(
			conditions,
			fmt.Sprintf("timestamp <= $%d", len(args)+1),
		)

		args = append(args, toTime)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := conn.Query(
		context.Background(),
		query,
		args...,
	)

	if err != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var logs []Log

	for rows.Next() {
		var logEntry Log

		err := rows.Scan(
			&logEntry.ID,
			&logEntry.Timestamp,
			&logEntry.Service,
			&logEntry.Level,
			&logEntry.Message,
		)

		if err != nil {
			http.Error(w, "Failed to read log", http.StatusInternalServerError)
			return
		}

		logs = append(logs, logEntry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func getLogByID(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {

	var foundLog Log

	id := strings.TrimPrefix(r.URL.Path, "/logs/")

	_, err1 := uuid.Parse(id)
	if err1 != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	row := conn.QueryRow(
		context.Background(),
		`SELECT id, timestamp, service, level, message
     FROM logs
     WHERE id = $1`,
		id,
	)

	err := row.Scan(
		&foundLog.ID,
		&foundLog.Timestamp,
		&foundLog.Service,
		&foundLog.Level,
		&foundLog.Message,
	)

	if err == pgx.ErrNoRows {
		http.Error(w, "Log not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundLog)
}
