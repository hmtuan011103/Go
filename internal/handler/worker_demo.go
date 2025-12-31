package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gostructure/app/internal/worker"
	"github.com/gostructure/app/pkg/response"
)

// ============================================================================
// API ENDPOINT ĐỂ TEST WORKER POOL
// ============================================================================
//
// Endpoint: POST /api/v1/demo/worker-pool?workers=3&jobs=10
//
// Bạn có thể test bằng curl:
//   curl -X POST "http://localhost:8080/api/v1/demo/worker-pool?workers=3&jobs=10"
//
// ============================================================================

// WorkerPoolDemo xử lý request demo worker pool
func (h *Handler) WorkerPoolDemo(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	numWorkers, _ := strconv.Atoi(r.URL.Query().Get("workers"))
	numJobs, _ := strconv.Atoi(r.URL.Query().Get("jobs"))

	// Defaults
	if numWorkers <= 0 {
		numWorkers = 3
	}
	if numJobs <= 0 {
		numJobs = 10
	}

	// Giới hạn để tránh abuse
	if numWorkers > 10 {
		numWorkers = 10
	}
	if numJobs > 100 {
		numJobs = 100
	}

	startTime := time.Now()

	// Tạo processor function
	processor := func(job worker.Job) (interface{}, error) {
		// Mô phỏng công việc tốn thời gian
		workTime := time.Duration(50+job.ID%100) * time.Millisecond
		time.Sleep(workTime)
		return fmt.Sprintf("Job %d processed", job.ID), nil
	}

	// Tạo và chạy worker pool
	pool := worker.NewPool(numWorkers, numJobs, processor)
	ctx := context.Background()
	pool.Start(ctx)

	// Gửi jobs
	go func() {
		for i := 1; i <= numJobs; i++ {
			pool.Submit(worker.Job{
				ID:      i,
				Payload: fmt.Sprintf("Data %d", i),
			})
		}
		pool.Close()
	}()

	// Thu thập kết quả
	var results []map[string]interface{}
	workerStats := make(map[int]int)
	var totalProcessTime time.Duration
	var mu sync.Mutex

	for result := range pool.Results() {
		mu.Lock()
		results = append(results, map[string]interface{}{
			"job_id":       result.JobID,
			"worker_id":    result.WorkerID,
			"output":       result.Output,
			"process_time": result.ProcessTime.String(),
		})
		workerStats[result.WorkerID]++
		totalProcessTime += result.ProcessTime
		mu.Unlock()
	}

	totalTime := time.Since(startTime)

	// Trả về response
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"config": map[string]int{
			"workers": numWorkers,
			"jobs":    numJobs,
		},
		"performance": map[string]interface{}{
			"total_time":           totalTime.String(),
			"total_process_time":   totalProcessTime.String(),
			"average_time_per_job": (totalProcessTime / time.Duration(numJobs)).String(),
			"parallelism_factor":   fmt.Sprintf("%.2fx", float64(totalProcessTime)/float64(totalTime)),
		},
		"worker_distribution": workerStats,
		"results":             results,
	})
}

// ConcurrencyExplain trả về giải thích về concurrency
func (h *Handler) ConcurrencyExplain(w http.ResponseWriter, r *http.Request) {
	explanation := map[string]interface{}{
		"title": "Go Concurrency Patterns",
		"concepts": []map[string]string{
			{
				"name":        "Goroutine",
				"description": "Lightweight thread managed by Go runtime. Starts with only 2KB stack.",
				"example":     "go func() { /* code */ }()",
			},
			{
				"name":        "Channel",
				"description": "Typed conduit for communication between goroutines. Thread-safe by design.",
				"example":     "ch := make(chan int, 10) // buffered channel",
			},
			{
				"name":        "WaitGroup",
				"description": "Synchronization primitive to wait for goroutines to complete.",
				"example":     "var wg sync.WaitGroup; wg.Add(1); go func() { defer wg.Done() }(); wg.Wait()",
			},
			{
				"name":        "Select",
				"description": "Multiplexer for channel operations. Like switch but for channels.",
				"example":     "select { case msg := <-ch1: case ch2 <- value: case <-ctx.Done(): }",
			},
		},
		"patterns": []map[string]string{
			{
				"name":        "Worker Pool",
				"description": "Fixed number of goroutines processing jobs from a queue",
				"use_case":    "Batch processing, rate limiting, resource control",
			},
			{
				"name":        "Fan-out/Fan-in",
				"description": "Distribute work to multiple goroutines and collect results",
				"use_case":    "Parallel API calls, distributed computing",
			},
			{
				"name":        "Pipeline",
				"description": "Chain of processing stages connected by channels",
				"use_case":    "Data transformation, ETL processes",
			},
		},
		"demo_endpoint": "POST /api/v1/demo/worker-pool?workers=3&jobs=10",
	}

	response.JSON(w, http.StatusOK, explanation)
}
