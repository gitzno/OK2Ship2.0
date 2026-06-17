package main

import (
	"export_module_service/internal/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.LoadConfig()

	//-- Init Rabbit MQ --
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

	// -- Init MinIO ---
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Lỗi khởi tạo MinIO Client: %v", err)
	}
	log.Println("Đã kết nối thành công tới cụm MinIO!")

	//--------Wiring------- Using ggwire
	
}
