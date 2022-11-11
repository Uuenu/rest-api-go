package author

import (
	"context"
	"fmt"
	"rest-api-go/internal/author"
	"rest-api-go/pkg/client/postgresql"
	"rest-api-go/pkg/logging"
	"strings"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) *repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}

func (r *repository) Create(ctx context.Context, author *author.Author) (string, error) {
	q := `INSERT INTO author(name)
		  VALUES($1)
		  RETURNING id`
	r.logger.Trace(fmt.Sprintf("SQL Querry: %s", formatQuery(q)))

	if err := r.client.QueryRow(ctx, q, author.Name).Scan(&author.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL error: %s, Detail: %s, Where: %s, Code: %s, SQLstate: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}
		return "", fmt.Errorf("postgresql - Create - r.client.QueryRow().Scan(). error: %v", err)
	}
	return author.ID, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (author.Author, error) {
	q := `
		SELECT id, name FROM public.author WHERE id = $1
	`

	//TODO logger.Trace()

	var ath author.Author
	err := r.client.QueryRow(ctx, q, id).Scan(&ath.ID, &ath.Name)
	if err != nil {
		return author.Author{}, fmt.Errorf("postgresql - FindAll - r.client.Query error: %v", err)
	}

	return ath, nil
}

func (r *repository) FindAll(ctx context.Context) (a []author.Author, err error) {
	q := `
		SELECT id, name FROM public.author;
	`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("postgresql - FindAll - r.client.Query error: %v", err)
	}

	authors := make([]author.Author, 0)

	for rows.Next() {
		var ath author.Author

		if err := rows.Scan(&ath.ID, &ath.Name); err != nil {
			return nil, fmt.Errorf("postgresql - FindAll - rows.Scan. error: %v", err)
		}

		authors = append(authors, ath)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *repository) Update(ctx context.Context, author author.Author) error {
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}
