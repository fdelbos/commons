package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/fdelbos/commons/db"
	"github.com/georgysavva/scany/sqlscan"
	_ "github.com/mattn/go-sqlite3"
)

type (
	SqlConn struct {
		db *sql.DB
	}

	query struct {
		conn *sql.DB
		ctx  context.Context
	}
)

func NewConn(path string) (*SqlConn, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &SqlConn{db: db}, nil
}

func (conn *SqlConn) Query(ctx context.Context) db.Query {
	return queryFromCtx(ctx, conn.db)
}

func (conn *SqlConn) Close() error {
	if conn.db != nil {
		return conn.db.Close()
	}
	return nil
}

func (conn *SqlConn) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return tx(ctx, conn.db, fn)
}

func (conn *SqlConn) Lock(ctx context.Context, lockID db.AdvisoryLockID, fn func(ctx context.Context) error) error {
	log.Fatal("sqlite does not support advisory locks!")
	return nil // never called
}

func (q *query) Exec(query string, arguments ...interface{}) error {
	_, err := q.conn.ExecContext(q.ctx, query, arguments...)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return db.ErrNoRows
	}
	return err
}

func (q *query) Select(dest interface{}, query string, args ...interface{}) error {
	err := sqlscan.Select(q.ctx, q.conn, dest, query, args...)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return db.ErrNoRows
	}
	return err
}

func (q *query) Get(dest interface{}, query string, args ...interface{}) error {
	err := sqlscan.Get(q.ctx, q.conn, dest, query, args...)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return db.ErrNoRows
	}
	return err
}
