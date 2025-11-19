package usecases

import "context"

type ImportData interface {
	ImportData(ctx context.Context) error
}
