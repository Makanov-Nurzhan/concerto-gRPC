package tx

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"time"
)

type mysqlTxManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) Manager {
	return &mysqlTxManager{db: db}
}
func (m mysqlTxManager) Do(ctx context.Context, fn func(tx *gorm.DB) error) error {
	var err error
	for i := 0; i < 3; i++ {
		err = m.db.WithContext(ctx).Transaction(fn)
		if err == nil {
			return nil
		}
		es := strings.ToLower(err.Error())
		if !strings.Contains(es, "deadlock") && !strings.Contains(es, "lock wait timeout") {
			return err
		}
		time.Sleep(time.Duration(50*(i+1)) * time.Millisecond)
	}
	return err
}
