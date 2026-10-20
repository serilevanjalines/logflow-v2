package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type BenchmarkLog struct {
	Service string `json:"service"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func worker(id int, wg *sync.WaitGroup,latencies chan<- time.Duration) {
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

		start := time.Now()

		resp, err := http.Post(
			"http://localhost:8080/logs",
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		elapsed := time.Since(start)

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

		latencies <- elapsed

	}
}

func percentile(sorted []time.Duration, percentile float64) time.Duration {
	index := int(percentile * float64(len(sorted)-1))
	return sorted[index]
}

func LoadTest() {
	var wg sync.WaitGroup

	latencies := make(chan time.Duration, 10000)

	start := time.Now()

	wg.Add(5)

	for i := 0; i < 5; i++ {
		go worker(i, &wg, latencies)
	}

	wg.Wait()

	close(latencies)

	totalElapsed := time.Since(start)

	var allLatencies []time.Duration

	for latency := range latencies {
		allLatencies = append(allLatencies, latency)
	}

	sort.Slice(allLatencies, func(i, j int) bool {
		return allLatencies[i] < allLatencies[j]
	})

	fmt.Printf("Completed %d requests in %v\n", len(allLatencies), totalElapsed)
	fmt.Printf("Throughput: %.2f requests/sec\n", float64(len(allLatencies))/totalElapsed.Seconds())

	fmt.Printf("p50 latency: %v\n", percentile(allLatencies, 0.50))
	fmt.Printf("p95 latency: %v\n", percentile(allLatencies, 0.95))
	fmt.Printf("p99 latency: %v\n", percentile(allLatencies, 0.99))

}

func main() {
	LoadTest()
}
