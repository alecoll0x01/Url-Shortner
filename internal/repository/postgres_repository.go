package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Save(url URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO urls (code, original) VALUES ($1, $2)`,
		url.Code, url.OriginalURL,
	)
	return err
}

func (r *PostgresRepository) FindByCode(code string) (URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var url URL
	err := r.pool.QueryRow(ctx,
		`SELECT code, original, clicks FROM urls WHERE code = $1`,
		code,
	).Scan(&url.Code, &url.OriginalURL, &url.Clicks)

	if errors.Is(err, pgx.ErrNoRows) {
		return URL{}, ErrNotFound
	}
	return url, err
}

func (r *PostgresRepository) IncrementClicks(code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx,
		`UPDATE urls SET clicks = clicks + 1 WHERE code = $1`,
		code,
	)
	return err
}

func (r *PostgresRepository) ListAll() []URL {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx,
		`SELECT code, original, clicks FROM urls ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.Code, &url.OriginalURL, &url.Clicks); err == nil {
			urls = append(urls, url)
		}
	}
	return urls
}

func (r *PostgresRepository) Delete(code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.pool.Exec(ctx, `DELETE FROM urls WHERE code = $1`, code)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
