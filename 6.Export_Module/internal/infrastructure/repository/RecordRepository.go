package repository

import (
	"context"
	"export_module_service/internal/domain"

	"gorm.io/gorm"
)

type gormRepo struct {
	db *gorm.DB
}

func NewGormRepo(db *gorm.DB) domain.RecordRepository {
	return &gormRepo{db}
}

func (r *gormRepo) Upsert(
	ctx context.Context,
	res domain.Record,
) error {
	// 1. Dữ liệu dùng để tìm kiếm (Conflict)
	searchCondition := domain.Record{
		ItemCode: res.ItemCode,
		LotNo:    res.LotNo,
		Category: res.Category,
		Type:     res.Type,
	}

	// 2. Dữ liệu sẽ được cập nhật đè lên nếu tìm thấy (Update)
	updateData := domain.Record{
		DataLogfile: res.DataLogfile,
		DataUser:    res.DataUser,
		// GORM sẽ tự động gọi Trigger để update cột LastModified dưới DB
	}

	// 3. Thực hiện Upsert
	err := r.db.WithContext(ctx).
		Where(searchCondition).
		Assign(updateData).
		FirstOrCreate(&res).Error

	return err
}
