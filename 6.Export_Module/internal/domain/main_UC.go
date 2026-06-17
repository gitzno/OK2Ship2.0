package domain

import "context"

type MainUseCase interface {
	ImportLogFile(ctx context.Context, fileAddress string, itemCode string, lotNo string) error
	ExportData(ctx context.Context, itemCode string, lotNo string) error
}
