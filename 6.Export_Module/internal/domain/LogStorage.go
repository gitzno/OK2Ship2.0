package domain

import (
	"context"
	"io"
)

type LogStorage interface {
	GetLogStream(ctx context.Context, filePath string) (io.ReadCloser, error)

	ListFileLocation(ctx context.Context, filePath string) ([]string, error)

	UploadRawFile(ctx context.Context, filePath string, objectName string, data []byte, contentType string) (string, error)

	CopyFile(ctx context.Context, srcPath string, destPath string) error
}
