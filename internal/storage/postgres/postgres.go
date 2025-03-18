package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PStorage struct {
	pool *pgxpool.Pool
	mu *sync.Mutex
}

func New(ctx context.Context, storagePath string) (*PStorage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.Connect(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &PStorage{
		pool: pool,
		mu: &sync.Mutex{},
	}, nil
}

func Close(ctx context.Context, storage *PStorage) {
	storage.pool.Close()
}