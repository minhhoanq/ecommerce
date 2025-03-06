package repo

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/user_service/pkg/postgres"
)

type SQLStore struct {
	db postgres.Database
}

func New(db postgres.Database) Querier {
	return &SQLStore{
		db: db,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(Querier) error) error {
	// Bắt đầu transaction với GORM
	tx := store.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	db := postgres.Database{DB: tx}
	q := New(db)
	// Thực thi function với transaction
	err := fn(q)
	if err != nil {
		// Nếu có lỗi, rollback transaction
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}
