package store

import (
	"context"
	"database/sql"
	"github.com/Kritvi0208/ShortEdge/model"
)

type URL interface {
	Create(ctx context.Context, url model.URL) error
	GetAll(ctx context.Context) ([]model.URL, error)
	GetByCode(ctx context.Context, code string) (model.URL, error)
	Update(ctx context.Context, code string, updated model.URL) error
	Delete(ctx context.Context, code string) error
}

type urlStore struct {
	db *sql.DB
}

func NewURLStore(db *sql.DB) URL {
	return &urlStore{db: db}
}

func (s *urlStore) Create(ctx context.Context, url model.URL) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO urls (code, long_url, created_at, visibility, expires_at) VALUES ($1, $2, $3, $4, $5)`,
		url.Code, url.LongURL, url.CreatedAt, url.Visibility, url.ExpiresAt)
	return err
}

func (s *urlStore) GetAll(ctx context.Context) ([]model.URL, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT code, long_url, created_at, visibility, expires_at FROM urls`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []model.URL
	for rows.Next() {
		var u model.URL
		err = rows.Scan(&u.Code, &u.LongURL, &u.CreatedAt, &u.Visibility, &u.ExpiresAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func (s *urlStore) GetByCode(ctx context.Context, code string) (model.URL, error) {
	var u model.URL
	err := s.db.QueryRowContext(ctx,
		`SELECT code, long_url, created_at, visibility, expires_at FROM urls WHERE code = $1`, code).
		Scan(&u.Code, &u.LongURL, &u.CreatedAt, &u.Visibility, &u.ExpiresAt)

	return u, err
}

func (s *urlStore) Update(ctx context.Context, code string, updated model.URL) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE urls SET long_url = $1, visibility = $2, expires_at = $3 WHERE code = $4`,
		updated.LongURL, updated.Visibility, updated.ExpiresAt, code)
	return err
}

func (s *urlStore) Delete(ctx context.Context, code string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM urls WHERE code = $1`, code)
	return err
}
