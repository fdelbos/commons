package pg

import (
	"context"
	"errors"

	"github.com/fdelbos/commons/db"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	pgPool struct {
		pool        *pgxpool.Pool
		databaseURL string
	}

	PgConn struct {
		conn        *pgx.Conn
		databaseURL string
	}

	pgxInterface interface {
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	}

	query struct {
		conn pgxInterface
		ctx  context.Context
	}
)

var (
	ErrNoConnectionInContext = errors.New("no connection found in context")
)

// NewConn creates a new connection to a postgres database.
// It must be closed when done.
func NewConn(url string) (*PgConn, error) {
	conn, err := newConn(context.Background(), url)
	if err != nil {
		return nil, err
	}
	res := &PgConn{
		conn:        conn,
		databaseURL: url,
	}

	return res, nil
}

func NewPool(url string) (*pgPool, error) {
	pool, err := newPool(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &pgPool{
		pool:        pool,
		databaseURL: url,
	}, nil
}

func (pg *pgPool) Query(ctx context.Context) db.Query {
	return queryFromCtx(ctx, pg.pool)
}

func (pg *pgPool) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return tx(ctx, pg.pool, fn)
}

func (pg *pgPool) Lock(ctx context.Context, lockID db.AdvisoryLockID, fn func(ctx context.Context) error) error {
	return lock(ctx, pg.pool, lockID, fn)
}

func (pg *PgConn) Lock(ctx context.Context, lockID db.AdvisoryLockID, fn func(ctx context.Context) error) error {
	return lock(ctx, pg.conn, lockID, fn)
}

func (pg *PgConn) Query(ctx context.Context) db.Query {
	return queryFromCtx(ctx, pg.conn)
}

func (pg *PgConn) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return tx(ctx, pg.conn, fn)
}

func (pg *PgConn) Close(ctx context.Context) error {
	return pg.conn.Close(ctx)
}

func (q *query) Exec(sql string, arguments ...interface{}) error {
	if q.conn == nil {
		return ErrNoConnectionInContext
	}
	_, err := q.conn.Exec(q.ctx, sql, arguments...)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return db.ErrNoRows
	}
	return err
}

func (q *query) Select(dest interface{}, sql string, args ...interface{}) error {
	if q.conn == nil {
		return ErrNoConnectionInContext
	}
	err := pgxscan.Select(q.ctx, q.conn, dest, sql, args...)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return db.ErrNoRows
	}
	return err
}

func (q *query) Get(dest interface{}, sql string, args ...interface{}) error {
	if q.conn == nil {
		return ErrNoConnectionInContext
	}
	err := pgxscan.Get(q.ctx, q.conn, dest, sql, args...)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return db.ErrNoRows
	}
	return err
}
