package amqp

import (
	"context"
	"encoding/json"
	"export_module_service/internal/domain"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ExcelEvent struct {
	FilePath string `json:"file_path"`
	ItemCode string `json:"item_code"`
	LotNo    string `json:"lot_no"`
	Category string `json:"type"` // Lưu ý: JSON key là "type"
}

type Consumer struct {
	// Dùng registry để lưu danh sách các Use Case
	useCases map[string]domain.MainUseCase
}

// Khởi tạo Consumer với một danh sách các Use Case
func NewConsumer(uc domain.MainUseCase) *Consumer {
	return &Consumer{
		useCases: map[string]domain.MainUseCase{
			"peel_test": uc, // Map cứng hoặc dùng Registry injection
		},
	}
}

func (c *Consumer) Start(messages <-chan amqp.Delivery) {
	log.Println("[*] Worker đang chờ event xử lý Excel...")

	for msg := range messages {
		var event ExcelEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Lỗi JSON: %v", err)
			msg.Nack(false, false)
			continue
		}

		// 1. ROUTING: Tìm Use Case trong map
		handler, exists := c.useCases[event.Category]
		if !exists {
			log.Printf("[CẢNH BÁO] Không tìm thấy Use Case cho category: %s", event.Category)
			msg.Nack(false, false)
			continue
		}

		// 2. THỰC THI: Gọi phương thức ImportLogFile với đầy đủ tham số
		ctx := context.Background()
		log.Printf("=> Đang điều phối category [%s]...", event.Category)

		err := handler.ImportLogFile(ctx, event.FilePath, event.ItemCode, event.LotNo)

		// 3. KẾT QUẢ
		if err != nil {
			log.Printf("[X] Lỗi xử lý file %s: %v", event.FilePath, err)
			msg.Nack(false, false)
		} else {
			log.Printf("[V] Xử lý thành công %s", event.FilePath)
			msg.Ack(false)
		}
	}
}
