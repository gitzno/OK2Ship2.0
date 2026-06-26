package main

import (
	"export_module_service/internal/config"
	"export_module_service/internal/infrastructure/repository/db"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 0. Tải cấu hình
	cfg := config.LoadConfig()
	// 1. Kết nối tới database
	strCon := cfg.SQLServerConn
	conSql, err := db.NewConnectionSQLServer(strCon)

	// 2. Khởi tạo kết nối RabbitMQ
	conn, err := rabbitmq.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Lỗi kết nối RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Lỗi mở Channel: %v", err)
	}
	defer ch.Close()

	// Khai báo Queue
	q, err := ch.QueueDeclare(cfg.RabbitMQQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Lỗi khai báo queue: %v", err)
	}

	// 3. Khởi tạo MinIO Client
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Lỗi khởi tạo MinIO Client: %v", err)
	}
	log.Println("Đã kết nối thành công tới cụm MinIO!")

	// 4. Wiring bằng Google Wire
	// Hàm InitializeWorker được sinh ra bởi wire_gen.go
	consumer := InitializeWorker(minioClient, conSql)

	// 5. Bắt đầu lắng nghe message
	msgs, err := ch.Consume(q.Name, "worker-export", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Lỗi đăng ký consumer: %v", err)
	}

	// Chạy consumer trong goroutine để không chặn main thread
	go consumer.Start(msgs)

	log.Println(" Export Worker đã sẵn sàng. Nhấn CTRL+C để thoát.")

	// 6. Graceful Shutdown (Xử lý khi ứng dụng bị ngắt)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	log.Println("Đang dừng worker...")
}
