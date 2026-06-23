package streaming_exp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Transaction đại diện cho một giao dịch đơn lẻ từ ngân hàng
type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

// ProcessLargeTransactions xử lý luồng dữ liệu JSON khổng lồ từ HTTP Response
func ProcessLargeTransactions(url string) error {
	// 1. Gọi API và xử lý Connection Pool / Leak
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("lỗi kết nối API: %w", err)
	}
	// Đảm bảo đóng vòi nước khi hàm kết thúc để không rò rỉ tài nguyên mạng
	defer res.Body.Close()

	// 2. Khởi tạo Decoder gắn trực tiếp vào vòi nước res.Body
	decoder := json.NewDecoder(res.Body)

	// 3. Đọc Token đầu tiên: bóc lớp vỏ mảng '['
	// Hàm Token() sẽ bóc lấy dấu ngoặc vuông mở đầu tiên của mảng JSON
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("lỗi đọc token mở đầu (kỳ vọng '['): %w", err)
	}

	var totalSuccessAmount float64
	var processedCount int

	// 4. Vòng lặp múc từng gáo nước: Lặp chừng nào mảng vẫn còn phần tử
	for decoder.More() {
		var tx Transaction

		// Múc đúng 1 object JSON (ví dụ: {"id": "TXN-001"...}) và giải mã vào struct tx
		err := decoder.Decode(&tx)
		if err != nil {
			log.Printf("Bỏ qua giao dịch lỗi: %v", err)
			continue
		}

		// XỬ LÝ NGAY LẬP TỨC: Không lưu mảng vào RAM, chỉ cập nhật biến tổng
		if tx.Status == "SUCCESS" {
			totalSuccessAmount += tx.Amount
		}
		processedCount++

		// Xong vòng lặp này, vùng nhớ của tx sẽ sẵn sàng để ghi đè ở vòng lặp sau.
		// Object cũ không còn ai tham chiếu tới sẽ được Garbage Collector tự động dọn dẹp.
	}

	// 5. Đọc Token cuối cùng: bóc lớp vỏ mảng đóng ']' để hoàn tất
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("lỗi đọc token kết thúc (kỳ vọng ']'): %w", err)
	}

	// Hiển thị kết quả cuối cùng
	fmt.Printf("Đã xử lý thành công %d giao dịch.\n", processedCount)
	fmt.Printf("Tổng dòng tiền hợp lệ: %.2f\n", totalSuccessAmount)

	return nil
}
