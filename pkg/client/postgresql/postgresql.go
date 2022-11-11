package postgresql

import (
	"context"
	"fmt"
	"rest-api-go/internal/config"
	"rest-api-go/pkg/utils"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return fmt.Errorf("postgressql - NewClient - pgxpool.Connect. error: %v", err)
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		return nil, fmt.Errorf("postgressql - NewClient - utilsDoWithTries. error: %v", err)
	}

	return pool, nil
}
