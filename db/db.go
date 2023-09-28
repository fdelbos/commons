package db

import (
	"context"
	"errors"
	"net/url"
)

type (
	AdvisoryLockID int

	Migrator func(string) error

	Query interface {
		// Exec is for INSERT, UPDATE, DELETE, CREATE, etc
		Exec(sql string, arguments ...any) error
		// Select is for multiple row results
		Select(dest interface{}, sql string, args ...any) error
		// Get is for single row results
		Get(dest interface{}, sql string, args ...any) error
	}

	DB interface {
		Query(ctx context.Context) Query
		Tx(ctx context.Context, fn func(ctx context.Context) error) error
		Lock(ctx context.Context, lockID AdvisoryLockID, fn func(ctx context.Context) error) error
	}
)

var (
	ErrNoRows     = errors.New("no rows in result set")
	ErrLockFailed = errors.New("lock failed")
)

func IsErrNoRows(err error) bool {
	return err == ErrNoRows
}

func ReplaceDBInURL(originalURL, newDb string) (string, error) {
	url, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}
	url.Path = newDb
	return url.String(), nil
}

func DBName(originalURL string) (string, error) {
	url, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}
	return url.Path[1:], nil
}
