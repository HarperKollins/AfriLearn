package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type LatencyResult struct {
	Endpoint string
	Duration time.Duration
	Status   int
	Err      error
}

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "Base URL of AfriLearn API")
	apiKey := flag.String("key", "afr_live_demo_9f8e2b7a", "API Key for authentication")
	concurrency := flag.Int("c", 10, "Concurrency level (concurrent workers)")
	requests := flag.Int("n", 200, "Total number of requests per endpoint test")
	flag.Parse()

	fmt.Println("🚀 AfriLearn API Load Test & Latency Benchmark")
	fmt.Printf("   Base URL:     %s\n", *baseURL)
	fmt.Printf("   Concurrency:  %d workers\n", *concurrency)
	fmt.Printf("   Total Reqs:   %d per test\n", *requests)
	fmt.Println("──────────────────────────────────────────────────────")

	// 1. Health check first
	resp, err := http.Get(*baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ API server unreachable at %s: %v\n", *baseURL, err)
		fmt.Println("   Start the API server first: go run cmd/api/main.go")
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Health check returned status %d\n", resp.StatusCode)
		return
	}
	fmt.Println("✅ API Health check passed")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	endpoints := []struct {
		name    string
		method  string
		path    string
		payload string
	}{
		{
			name:   "GET Curriculum (WAEC Math)",
			method: "GET",
			path:   "/api/v1/curriculum/waec/mathematics",
		},
		{
			name:   "GET LLM Prompt (UNILAG Law)",
			method: "GET",
			path:   "/api/v1/curriculum/unilag/law/llm-prompt",
		},
		{
			name:   "GET RAG Embeddings (WAEC Physics)",
			method: "GET",
			path:   "/api/v1/curriculum/waec/physics/embeddings",
		},
		{
			name:   "GET Deep FTS Search ('quadratic')",
			method: "GET",
			path:   "/api/v1/search?q=quadratic",
		},
		{
			name:   "POST Vector Search ('photosynthesis')",
			method: "POST",
			path:   "/api/v1/search/vector",
			payload: `{"query": "photosynthesis in plants", "limit": 5}`,
		},
		{
			name:   "POST Query Brain ('physics topics')",
			method: "POST",
			path:   "/api/v1/query",
			payload: `{"query": "show me WAEC physics topics"}`,
		},
	}

	for _, ep := range endpoints {
		runEndpointBenchmark(client, *baseURL, *apiKey, ep.name, ep.method, ep.path, ep.payload, *concurrency, *requests)
	}

	fmt.Println("\n🎉 Load test suite completed!")
}

func runEndpointBenchmark(client *http.Client, baseURL, apiKey, name, method, path, payload string, concurrency, totalReqs int) {
	fmt.Printf("\n⚡ Benchmarking: %s\n", name)

	reqsPerWorker := totalReqs / concurrency
	results := make(chan LatencyResult, totalReqs)
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < reqsPerWorker; j++ {
				var bodyReader io.Reader
				if payload != "" {
					bodyReader = bytes.NewBufferString(payload)
				}

				req, err := http.NewRequest(method, baseURL+path, bodyReader)
				if err != nil {
					results <- LatencyResult{Endpoint: name, Err: err}
					continue
				}
				req.Header.Set("X-API-Key", apiKey)
				if payload != "" {
					req.Header.Set("Content-Type", "application/json")
				}

				t0 := time.Now()
				res, err := client.Do(req)
				dur := time.Since(t0)

				if err != nil {
					results <- LatencyResult{Endpoint: name, Duration: dur, Err: err}
					continue
				}

				io.Copy(io.Discard, res.Body)
				res.Body.Close()

				results <- LatencyResult{Endpoint: name, Duration: dur, Status: res.StatusCode}
			}
		}()
	}

	wg.Wait()
	totalDuration := time.Since(startTime)
	close(results)

	var durational []time.Duration
	successCount := 0
	errorCount := 0

	for r := range results {
		if r.Err != nil || r.Status != http.StatusOK {
			errorCount++
		} else {
			successCount++
			durational = append(durational, r.Duration)
		}
	}

	if len(durational) == 0 {
		fmt.Printf("   ❌ All requests failed! Errors: %d\n", errorCount)
		return
	}

	sort.Slice(durational, func(i, j int) bool {
		return durational[i] < durational[j]
	})

	min := durational[0]
	max := durational[len(durational)-1]
	p50 := durational[int(float64(len(durational))*0.50)]
	p90 := durational[int(float64(len(durational))*0.90)]
	p99 := durational[int(float64(len(durational))*0.99)]

	rps := float64(successCount) / totalDuration.Seconds()

	fmt.Printf("   Success:    %d/%d (Errors: %d)\n", successCount, totalReqs, errorCount)
	fmt.Printf("   Throughput: %.2f req/sec (Total Time: %v)\n", rps, totalDuration.Truncate(time.Millisecond))
	fmt.Printf("   Latency:    Min: %v | p50: %v | p90: %v | p99: %v | Max: %v\n",
		min.Truncate(time.Microsecond),
		p50.Truncate(time.Microsecond),
		p90.Truncate(time.Microsecond),
		p99.Truncate(time.Microsecond),
		max.Truncate(time.Microsecond),
	)
}

func prettyJSON(b []byte) string {
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	return out.String()
}
