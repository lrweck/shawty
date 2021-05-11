package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	pool "github.com/jackc/pgx/v4/pgxpool"
	short "github.com/lrweck/shawty/shortener"
	"github.com/pkg/errors"
)

type pgRepo struct {
	conn    *pool.Pool
	timeout time.Duration
}

func newPgClient(pgURL string, timeout int) (*pool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	connPool, err := pool.Connect(ctx, pgURL)
	if err != nil {
		return nil, err
	}

	err = connPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return connPool, nil
}

// NewPGRepo - Creates a new PostgreSQl repository to store and consume data.
func NewPGRepo(pgURL string, timeout int) (short.RedirectRepository, error) {
	repo := &pgRepo{
		timeout: time.Duration(timeout) * time.Second,
	}
	conn, err := newPgClient(pgURL, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewPGRepo")
	}
	repo.conn = conn
	return repo, nil
}

// Find by code the redirect queried
func (r *pgRepo) Find(code string) (*short.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	redirect := &short.Redirect{}

	err := r.conn.QueryRow(ctx, findCodeSQL(), code).Scan(redirect.URL, redirect.Code, redirect.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrap(short.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

func findCodeSQL() string {
	return `SELECT url, code, created_at FROM REDIRECTS WHERE code = $1 LIMIT 1;`
}
func insertOneSQL() string {
	return `INSERT INTO REDIRECTS (code,url,created_at) VALUES ($1::text,$2::text,$3::timestamptz);`
}

// Store the user supplied redirect
func (r *pgRepo) Store(redirect *short.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.conn.Exec(ctx, insertOneSQL(), redirect.Code, redirect.URL, time.Unix(redirect.CreatedAt, 0))
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
