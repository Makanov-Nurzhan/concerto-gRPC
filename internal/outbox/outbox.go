package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID    string      // message_id (уникальный e2e)
	Topic string      // имя топика
	Key   string      // partition key
	Body  interface{} // payload в JSON
}

func Publish(ctx context.Context, tx *gorm.DB, m Message) error {
	if tx == nil {
		return errors.New("outbox.Publish: tx is nil")
	}
	if m.ID == "" || m.Topic == "" || m.Key == "" || m.Body == nil {
		return errors.New("outbox.Publish: required fields empty")
	}
	payload, err := json.Marshal(m.Body)
	if err != nil {
		return err
	}

	sqlDB, err := tx.DB()
	if err == nil {
		pingCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if pingErr := sqlDB.PingContext(pingCtx); pingErr != nil {
			time.Sleep(100 * time.Millisecond)
			pingCtx2, cancel2 := context.WithTimeout(ctx, 1*time.Second)
			defer cancel2()
			if pingErr2 := sqlDB.PingContext(pingCtx2); pingErr2 != nil {
				return fmt.Errorf("outbox.Publish: ping failed after retry: %w", pingErr2)
			}
		}
	}
	return tx.Exec(`
        INSERT INTO online_ko_outbox_messages (message_id, topic, kkey, payload, status, created_at)
        VALUES (?, ?, ?, ?, 'NEW', NOW())
    `, m.ID, m.Topic, m.Key, payload).Error
}
