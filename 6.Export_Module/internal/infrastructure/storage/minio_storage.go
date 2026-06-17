package storage

import (
	"context"
	"errors"
	"export_module_service/internal/domain"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
}

func NewMinioStorage(client *minio.Client) domain.LogStorage {
	return &MinioStorage{client: client}
}

func (m *MinioStorage) GetLogStream(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return nil, errors.New("Not Implemented!")
}
