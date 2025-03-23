package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute DB queries and transactions
type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(Queries *Queries) error) error
}

// SQLStore implements the Store interface
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) ExecTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}
