package domain

import (
	"context"
)

type Record struct {
	ID          int    `gorm:"column:ID;primaryKey;autoIncrement"`
	ItemCode    string `gorm:"column:ItemCode"`
	LotNo       string `gorm:"column:LotNo"`
	Category    string `gorm:"column:Category"`
	Type        string `gorm:"column:Type"`
	DataLogfile string `gorm:"column:DataLogfile"`
	DataUser    string `gorm:"column:DataUser"`
}

func (Record) TableName() string {
	return "Record"
}

type RecordRepository interface {
	Upsert(
		ctx context.Context,
		res Record,
	) error
}
