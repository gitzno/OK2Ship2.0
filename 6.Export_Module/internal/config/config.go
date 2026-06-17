package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	WorkerCount    int
	RabbitMQURL    string
	RabbitMQQueue  string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioUseSSL    bool
	SQLServerConn  string
}

// LoadConfig tự động tìm và nạp file .env một cách thông minh
func LoadConfig() *AppConfig {
	// 1. Định nghĩa cờ (flag) để nhận đường dẫn thủ công (nếu có)
	envPath := flag.String("env", "", "Đường dẫn thủ công tới file .env")
	flag.Parse()

	var loaded bool

	// 2. KỊCH BẢN A: Người dùng truyền đường dẫn rõ ràng qua cờ -env
	if *envPath != "" {
		err := godotenv.Load(*envPath)
		if err == nil {
			log.Printf("✅ Đã load cấu hình thủ công từ: %s", *envPath)
			loaded = true
		} else {
			log.Printf("❌ Không thể nạp file từ %s. Lỗi: %v", *envPath, err)
		}
	}

	// 3. KỊCH BẢN B: Không truyền cờ, hệ thống TỰ ĐỘNG TÌM KIẾM
	if !loaded {
		// Danh sách các vị trí có khả năng chứa file .env
		pathsToTry := []string{
			".env",       // Thư mục hiện tại (ngang với nơi gõ lệnh go run)
			"../.env",    // Lùi 1 cấp (Thường dùng khi chạy test/debug trong folder con)
			"../../.env", // Lùi 2 cấp
		}

		for _, p := range pathsToTry {
			err := godotenv.Load(p)
			if err == nil {
				log.Printf("✅ Tự động tìm thấy và load cấu hình tại: %s", p)
				loaded = true
				break // Tìm thấy cái đầu tiên là dừng lặp ngay lập tức
			}
		}
	}

	// 4. KỊCH BẢN C: Vẫn không tìm thấy (VD: Chạy trên Production/Docker)
	if !loaded {
		log.Println("⚠️ Không tìm thấy file .env cục bộ. Hệ thống sẽ sử dụng biến môi trường của Hệ điều hành (OS/Docker ENV).")
	}
	port := getEnv("RABBITMQ_PORT", "5672")
	user := getEnv("RABBITMQ_DEFAULT_USER", "guest")
	pass := getEnv("RABBITMQ_DEFAULT_PASS", "guest")
	rabbitMQ_url := fmt.Sprintf("amqp://%v:%v@localhost:%v/", user, pass, port)

	port = getEnv("MINIO_HOT_API_PORT", "9000")
	minIO_url := fmt.Sprintf("http://localhost:%v/", port)

	// 5. Gom dữ liệu vào Struct
	return &AppConfig{
		WorkerCount:    getEnvAsInt("WORKER_COUNT", 5),
		RabbitMQURL:    rabbitMQ_url,
		RabbitMQQueue:  getEnv("RABBITMQ_QUEUE", "excel.import.queue"),
		MinioEndpoint:  minIO_url,
		MinioAccessKey: getEnv("MINIO_ROOT_USER_HOT", "admin_hot_default"),
		MinioSecretKey: getEnv("MINIO_ROOT_PASSWORD_HOT", "password123_hot_default"),
		SQLServerConn:  getEnv("SQL_SERVER_CONN", ""),
	}
}
func getURLMQ(port string, account string, password string) {

}
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if value, err := strconv.Atoi(valStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valStr); err == nil {
		return value
	}
	return fallback
}
