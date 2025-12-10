package tx

import (
	"context"
	"gorm.io/gorm"
)

type Manager interface {
	Do(ctx context.Context, fn func(tx *gorm.DB) error) error
}
