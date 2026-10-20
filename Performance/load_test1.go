package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type BenchmarkLog struct {
	Service string `json:"service"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 2000; i++ {

		logData := BenchmarkLog{
			Service: "payment",
			Level:   "INFO",
			Message: "benchmark test",
		}

		jsonData, err := json.Marshal(logData)
		if err != nil {
			fmt.Printf("Worker %d: failed to encode JSON: %v\n", id, err)
			continue
		}

		resp, err := http.Post(
			"http://localhost:8080/logs",
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			fmt.Printf("Worker %d: request failed: %v\n", id, err)
			continue
		}

		resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			fmt.Printf(
				"Worker %d: unexpected status: %d\n",
				id,
				resp.StatusCode,
			)
			continue
		}
	}
}

func LoadTest() {
	var wg sync.WaitGroup

	start := time.Now()

	wg.Add(5)

	for i := 0; i < 5; i++ {
		go worker(i, &wg)
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("Completed 10000 requests in %v\n", elapsed)
	fmt.Printf("Throughput: %.2f requests/sec\n", 10000/elapsed.Seconds())
}

func main() {
	LoadTest()
}