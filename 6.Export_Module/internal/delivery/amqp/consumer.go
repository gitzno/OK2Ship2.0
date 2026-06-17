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
	Category string `json:"type"`
}

type Consumer struct {
	handlers map[string]domain.LogImporterUseCase
}

func NewConsumer() *Consumer {
	return &Consumer{
		handlers: make(map[string]domain.LogImporterUseCase),
	}
}

func (c *Consumer) RegisterHandler(category string, uc domain.LogImporterUseCase) {
	c.handlers[category] = uc
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

		// 1. ROUTING: Tìm Use Case xử lý tương ứng với Category
		handler, exists := c.handlers[event.Category]
		if !exists {
			log.Printf("[CẢNH BÁO] Không tìm thấy Use Case cho category: %s", event.Category)
			// Đẩy vào Dead Letter Queue vì không biết xử lý sao
			msg.Nack(false, false)
			continue
		}

		// 2. THỰC THI: Gọi Use Case tương ứng
		ctx := context.Background()
		log.Printf("=> Đang điều phối category [%s] tới đúng handler...", event.Category)
		err := handler.ImportLogFile(ctx, event.FilePath)

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
