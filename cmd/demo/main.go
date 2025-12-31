package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gostructure/app/internal/worker"
)

// ============================================================================
// DEMO CONCURRENCY TRONG GO
// ============================================================================
//
// Chạy file này để xem Worker Pool hoạt động:
//     go run cmd/demo/main.go
//
// ============================================================================

func main() {
	fmt.Println(`
	╔══════════════════════════════════════════════════════════════════════════════╗
	║                    🚀 DEMO: CONCURRENCY TRONG GO 🚀                          ║
	║                         Worker Pool Pattern                                   ║
	╚══════════════════════════════════════════════════════════════════════════════╝

	📚 KHÁI NIỆM CƠ BẢN:
	────────────────────
	1. GOROUTINE: Lightweight thread do Go runtime quản lý
	- Khởi tạo chỉ với 2KB stack (so với 1MB của OS thread)
	- Có thể chạy hàng triệu goroutines đồng thời

	2. CHANNEL: Cơ chế giao tiếp an toàn giữa goroutines
	- Tuân theo nguyên tắc: "Don't communicate by sharing memory;
		share memory by communicating"

	3. WORKER POOL: Pattern sử dụng N goroutines cố định để xử lý M jobs
	- Kiểm soát tài nguyên
	- Tránh tạo quá nhiều goroutines

	━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	`)

	// ========================================
	// SO SÁNH: TUẦN TỰ VS ĐỒNG THỜI
	// ========================================
	numJobs := 10

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📌 DEMO 1: XỬ LÝ TUẦN TỰ (Sequential)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	sequentialTime := runSequential(numJobs)

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📌 DEMO 2: XỬ LÝ ĐỒNG THỜI VỚI WORKER POOL (3 workers)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	result := worker.RunDemo(3, numJobs)

	// ========================================
	// SO SÁNH KẾT QUẢ
	// ========================================
	fmt.Println("\n" + `
╔══════════════════════════════════════════════════════════════════════════════╗
║                           📊 SO SÁNH KẾT QUẢ                                 ║
╚══════════════════════════════════════════════════════════════════════════════╝`)

	fmt.Printf(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ Phương pháp           │ Thời gian      │ Speedup                            │
├─────────────────────────────────────────────────────────────────────────────┤
│ Tuần tự               │ %-14v │ 1.0x (baseline)                    │
│ Worker Pool (3)       │ %-14v │ %.1fx nhanh hơn                    │
└─────────────────────────────────────────────────────────────────────────────┘
`, sequentialTime, result.TotalTime, float64(sequentialTime)/float64(result.TotalTime))

	fmt.Println(`
	💡 GIẢI THÍCH:
	─────────────
	• Với xử lý TUẦN TỰ: Jobs được xử lý lần lượt, job sau phải đợi job trước hoàn thành
	• Với WORKER POOL: 3 workers xử lý đồng thời, giảm thời gian đáng kể
	• Speedup lý thuyết tối đa = số workers (nếu jobs độc lập và cùng độ dài)

	🎯 KHI NÀO SỬ DỤNG WORKER POOL:
	──────────────────────────────
	✓ Xử lý batch nhiều items (ảnh, files, records)
	✓ Gọi nhiều API đồng thời
	✓ Xử lý queue messages
	✓ Web scraping nhiều URLs
	✓ Database bulk operations

	⚠️  LƯU Ý:
	─────────
	• Số workers tối ưu phụ thuộc vào loại công việc:
	- I/O bound (network, disk): Có thể dùng nhiều workers
	- CPU bound: Thường bằng số CPU cores (runtime.NumCPU())
	`)

	log.Println("✅ Demo hoàn thành!")
}

// runSequential chạy jobs tuần tự để so sánh
func runSequential(numJobs int) time.Duration {
	log.Println("🔄 Bắt đầu xử lý tuần tự...")
	start := time.Now()

	for i := 1; i <= numJobs; i++ {
		// Mô phỏng công việc (giống trong worker pool)
		workTime := time.Duration(100+i%200) * time.Millisecond
		time.Sleep(workTime)
		log.Printf("   Job %d hoàn thành (took %v)", i, workTime)
	}

	elapsed := time.Since(start)
	log.Printf("⏱️  Xử lý tuần tự hoàn thành trong: %v\n", elapsed)
	return elapsed
}
