package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Log struct {
	Service string `json:"service"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type LogResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/logs", logHandler)

	log.Println("LogFlow listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {

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
	
	log.Printf(
		"service=%s level=%s message=%s",
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
