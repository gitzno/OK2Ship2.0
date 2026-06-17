package main

import (
	"export_module_service/internal/domain"
	"export_module_service/internal/infrastructure/storage"
	"export_module_service/internal/usecase"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
)

var useCaseSet = wire.NewSet(
	usecase.NewPeelTestUC,
	wire.Bind(new(domain.MainUseCase), new(*usecase.NewPeelTestUC)),
)

func InitializeWorker(db *gorm.DB, minioClient *minio.Client) (domain.MainUseCase, error) {
	wire.Build(
		useCaseSet,
		storage.NewMinioStorage,
	)
	return nil, nil
}
