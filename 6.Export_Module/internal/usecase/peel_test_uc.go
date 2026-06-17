package usecase

import (
	"context"
	"export_module_service/internal/domain"
	"fmt"
)

type PeelTestUC struct {
	storage domain.LogStorage
}

func NewPeelTestUC(s domain.LogStorage) *PeelTestUC {
	return &PeelTestUC{storage: s}
}

// 1. Thực thi hàm ImportLogFile (Đọc logfile được người dùng up lên object storage -> lưu DB và ảnh di chuyển đi object storage)
func (uc *PeelTestUC) ImportLogFile(ctx context.Context, fileAddress string, itemCode string, lotNo string) error {

	stream, err := uc.storage.GetLogStream(ctx, fileAddress)
	if err != nil {
		return fmt.Errorf("không thể lấy log stream: %w", err)
	}
	defer stream.Close()

	// B. Xử lý logic bóc tách (Parsing)
	// (Ở đây bạn gọi thư viện excelize để đọc stream)
	// records, err := parseExcel(stream, itemCode, lotNo)

	// C. Gọi Repository để lưu vào Database
	// return uc.repo.InsertBatch(ctx, records)

	return nil
}

// 2. Thực thi hàm ExportData (Đọc từ DB -> Tạo file -> Đẩy lên MinIO)
func (uc *PeelTestUC) ExportData(ctx context.Context, itemCode string, lotNo string) error {
	// A. Lấy dữ liệu từ Database (ví dụ: lấy dữ liệu theo khoảng thời gian)
	// data, err := uc.repo.FetchDataForExport(ctx)

	// B. Tạo file Excel mới trong bộ nhớ
	// file, err := createExcelFile(data)

	// C. Upload file vừa tạo lên MinIO
	// _, err = uc.storage.UploadImage(ctx, fileBytes, ".xlsx")

	return nil
}
