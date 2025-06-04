package db

import (
	"context"
	"fmt"

	"github.com/pdridh/k-line/db/sqlc"
)

// Wrapper for easy transactions
func (store *psqlStore) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.pool.Begin(ctx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
