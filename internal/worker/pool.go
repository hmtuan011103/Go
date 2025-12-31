package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// ============================================================================
// WORKER POOL PATTERN - Pattern kinh điển trong Go
// ============================================================================
//
// GIẢI THÍCH:
// -----------
// Worker Pool là pattern sử dụng một số lượng cố định các goroutine (workers)
// để xử lý nhiều tác vụ (jobs) từ một hàng đợi (queue/channel).
//
// TẠI SAO CẦN WORKER POOL?
// ------------------------
// 1. Giới hạn số goroutine chạy đồng thời (tránh tạo hàng ngàn goroutine)
// 2. Tái sử dụng goroutine thay vì tạo mới cho mỗi tác vụ
// 3. Kiểm soát tài nguyên (memory, CPU, connections)
// 4. Xử lý backpressure khi có quá nhiều jobs
//
// SƠ ĐỒ HOẠT ĐỘNG:
// ----------------
//
//                    ┌──────────────┐
//                    │   Producer   │ (Gửi jobs vào channel)
//                    └──────┬───────┘
//                           │
//                           ▼
//                    ┌──────────────┐
//                    │  Jobs Channel│ (Buffered channel - hàng đợi)
//                    └──────┬───────┘
//                           │
//          ┌────────────────┼────────────────┐
//          ▼                ▼                ▼
//    ┌──────────┐    ┌──────────┐    ┌──────────┐
//    │ Worker 1 │    │ Worker 2 │    │ Worker 3 │  (Goroutines)
//    └────┬─────┘    └────┬─────┘    └────┬─────┘
//         │               │               │
//         └───────────────┼───────────────┘
//                         ▼
//                  ┌──────────────┐
//                  │Results Channel│ (Thu thập kết quả)
//                  └──────┬───────┘
//                         ▼
//                  ┌──────────────┐
//                  │   Consumer   │ (Xử lý kết quả)
//                  └──────────────┘
//
// ============================================================================

// Job đại diện cho một tác vụ cần xử lý
type Job struct {
	ID      int         // Định danh của job
	Payload interface{} // Dữ liệu cần xử lý
}

// Result đại diện cho kết quả sau khi job được xử lý
type Result struct {
	JobID       int           // ID của job đã xử lý
	Output      interface{}   // Kết quả
	Error       error         // Lỗi nếu có
	WorkerID    int           // Worker nào đã xử lý
	ProcessTime time.Duration // Thời gian xử lý
}

// Pool là cấu trúc quản lý worker pool
type Pool struct {
	numWorkers int            // Số lượng workers
	jobs       chan Job       // Channel nhận jobs
	results    chan Result    // Channel gửi kết quả
	wg         sync.WaitGroup // WaitGroup để đợi tất cả workers hoàn thành
	processor  ProcessorFunc  // Hàm xử lý job
}

// ProcessorFunc là kiểu hàm để xử lý job
// Bạn có thể tùy chỉnh hàm này để xử lý các loại job khác nhau
type ProcessorFunc func(job Job) (interface{}, error)

// NewPool tạo một worker pool mới
// - numWorkers: Số lượng goroutine workers
// - bufferSize: Kích thước buffer của jobs channel
// - processor: Hàm xử lý job
func NewPool(numWorkers, bufferSize int, processor ProcessorFunc) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, bufferSize),    // Buffered channel cho jobs
		results:    make(chan Result, bufferSize), // Buffered channel cho results
		processor:  processor,
	}
}

// Start khởi động tất cả workers
// Mỗi worker là một goroutine chạy vòng lặp, lấy jobs từ channel và xử lý
func (p *Pool) Start(ctx context.Context) {
	log.Printf("🚀 Khởi động Worker Pool với %d workers", p.numWorkers)

	for i := 1; i <= p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i) // Mỗi worker chạy trong goroutine riêng
	}
}

// worker là hàm chạy trong mỗi goroutine
// Nó liên tục lấy jobs từ channel và xử lý cho đến khi channel đóng
func (p *Pool) worker(ctx context.Context, workerID int) {
	defer p.wg.Done()
	log.Printf("👷 Worker %d: Bắt đầu làm việc", workerID)

	for {
		select {
		case <-ctx.Done():
			// Context bị cancel - shutdown gracefully
			log.Printf("👷 Worker %d: Nhận tín hiệu dừng", workerID)
			return

		case job, ok := <-p.jobs:
			if !ok {
				// Channel đã đóng - không còn jobs
				log.Printf("👷 Worker %d: Jobs channel đóng, kết thúc", workerID)
				return
			}

			// Xử lý job và đo thời gian
			startTime := time.Now()
			output, err := p.processor(job)
			processTime := time.Since(startTime)

			// Gửi kết quả vào results channel
			p.results <- Result{
				JobID:       job.ID,
				Output:      output,
				Error:       err,
				WorkerID:    workerID,
				ProcessTime: processTime,
			}

			log.Printf("👷 Worker %d: Hoàn thành Job %d trong %v",
				workerID, job.ID, processTime)
		}
	}
}

// Submit gửi một job vào hàng đợi để xử lý
func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

// Results trả về channel để đọc kết quả
func (p *Pool) Results() <-chan Result {
	return p.results
}

// Close đóng jobs channel và đợi tất cả workers hoàn thành
func (p *Pool) Close() {
	log.Println("🔒 Đóng jobs channel...")
	close(p.jobs) // Đóng channel - workers sẽ nhận được tín hiệu kết thúc

	log.Println("⏳ Đợi tất cả workers hoàn thành...")
	p.wg.Wait() // Đợi tất cả workers kết thúc

	log.Println("🔒 Đóng results channel...")
	close(p.results)

	log.Println("✅ Worker Pool đã shutdown hoàn toàn")
}

// ============================================================================
// DEMO: Ví dụ sử dụng Worker Pool
// ============================================================================

// DemoResult chứa kết quả của demo
type DemoResult struct {
	TotalJobs       int
	SuccessfulJobs  int
	FailedJobs      int
	TotalTime       time.Duration
	AverageTime     time.Duration
	ResultsByWorker map[int]int
}

// RunDemo chạy demo worker pool
// - numWorkers: Số workers
// - numJobs: Số jobs cần xử lý
func RunDemo(numWorkers, numJobs int) DemoResult {
	log.Println("=" + fmt.Sprintf("%60s", "") + "=")
	log.Println("🎬 BẮT ĐẦU DEMO WORKER POOL")
	log.Printf("📊 Cấu hình: %d workers, %d jobs", numWorkers, numJobs)
	log.Println("=" + fmt.Sprintf("%60s", "") + "=")

	startTime := time.Now()

	// Tạo processor function - mô phỏng công việc tốn thời gian
	processor := func(job Job) (interface{}, error) {
		// Mô phỏng công việc tốn thời gian (100-300ms)
		workTime := time.Duration(100+job.ID%200) * time.Millisecond
		time.Sleep(workTime)

		// Mô phỏng tính toán
		result := fmt.Sprintf("Processed job %d with payload: %v", job.ID, job.Payload)
		return result, nil
	}

	// Tạo và khởi động pool
	pool := NewPool(numWorkers, numJobs, processor)
	ctx := context.Background()
	pool.Start(ctx)

	// ========================================
	// FAN-OUT: Gửi nhiều jobs đồng thời
	// ========================================
	log.Println("\n📤 FAN-OUT: Gửi jobs vào queue...")
	go func() {
		for i := 1; i <= numJobs; i++ {
			job := Job{
				ID:      i,
				Payload: fmt.Sprintf("Data for job %d", i),
			}
			pool.Submit(job)
			log.Printf("📤 Đã gửi Job %d vào queue", i)
		}
		// Đóng pool sau khi gửi hết jobs
		pool.Close()
	}()

	// ========================================
	// FAN-IN: Thu thập kết quả từ nhiều workers
	// ========================================
	log.Println("\n📥 FAN-IN: Thu thập kết quả...")

	demoResult := DemoResult{
		TotalJobs:       numJobs,
		ResultsByWorker: make(map[int]int),
	}
	var totalProcessTime time.Duration

	for result := range pool.Results() {
		if result.Error != nil {
			demoResult.FailedJobs++
			log.Printf("❌ Job %d thất bại: %v", result.JobID, result.Error)
		} else {
			demoResult.SuccessfulJobs++
			totalProcessTime += result.ProcessTime
			demoResult.ResultsByWorker[result.WorkerID]++
		}
	}

	demoResult.TotalTime = time.Since(startTime)
	if demoResult.SuccessfulJobs > 0 {
		demoResult.AverageTime = totalProcessTime / time.Duration(demoResult.SuccessfulJobs)
	}

	// In kết quả
	log.Println("\n" + "=" + fmt.Sprintf("%60s", "") + "=")
	log.Println("📊 KẾT QUẢ DEMO")
	log.Println("=" + fmt.Sprintf("%60s", "") + "=")
	log.Printf("✅ Jobs thành công: %d/%d", demoResult.SuccessfulJobs, demoResult.TotalJobs)
	log.Printf("❌ Jobs thất bại: %d", demoResult.FailedJobs)
	log.Printf("⏱️  Tổng thời gian: %v", demoResult.TotalTime)
	log.Printf("⏱️  Thời gian trung bình/job: %v", demoResult.AverageTime)
	log.Println("\n📈 Phân bố công việc theo worker:")
	for workerID, count := range demoResult.ResultsByWorker {
		log.Printf("   Worker %d: xử lý %d jobs", workerID, count)
	}

	return demoResult
}
