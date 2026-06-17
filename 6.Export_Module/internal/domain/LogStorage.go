package domain

import (
	"context"
	"io"
)

type LogStorage interface {
	GetLogStream(ctx context.Context, filePath string) (io.ReadCloser, error)
}
