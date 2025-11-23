package postgres

import (
	"context"
	"fmt"

	"github.com/Traunin/review-assigner/internal/infrastructure/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	*sqlc.Queries
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool:    pool,
		Queries: sqlc.New(pool),
	}
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := sqlc.New(tx)

	if err := fn(q); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
