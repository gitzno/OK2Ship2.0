//go:build wireinject
// +build wireinject

package main

import (
	"export_module_service/internal/delivery/amqp"
	"export_module_service/internal/domain"
	"export_module_service/internal/infrastructure/storage"
	"export_module_service/internal/usecase"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
)

func InitializeWorker(minioClient *minio.Client) *amqp.Consumer {
	wire.Build(

		storage.NewMinioStorage,
		usecase.NewPeelTestUC,
		amqp.NewConsumer,

		wire.Bind(new(domain.MainUseCase), new(*usecase.PeelTestUC)),
	)
	return &amqp.Consumer{}
}
