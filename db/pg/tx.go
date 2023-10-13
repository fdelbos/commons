package pg

import (
	"context"
	"errors"
	"log"

	"github.com/fdelbos/commons/db"
	"github.com/jackc/pgx/v5"
)

type (
	txBeginner interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
)

func tx(ctx context.Context, conn txBeginner, fn func(ctx context.Context) error) error {

	if q := ctx.Value(pgCtx); q != nil {
		if _, ok := q.(pgx.Tx); ok {
			log.Fatal("db/pg database context is already in a transaction")
		}
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Printf("db/pg cant obtain transaction on the database: %v", err)
		return err
	}

	defer func() {
		if !tx.Conn().IsClosed() {
			if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				log.Printf("db/pg error while rolling back: %v", err)
			}
		}
	}()

	fnCtx := context.WithValue(ctx, pgCtx, tx)
	if err := fn(fnCtx); err != nil {
		return err
	} else {
		return tx.Commit(ctx)
	}
}

func lock(ctx context.Context, conn txBeginner, lockID db.AdvisoryLockID, fn func(ctx context.Context) error) error {
	return tx(ctx, conn, func(ctx context.Context) error {
		dest := struct {
			Locked bool `db:"locked"`
		}{}
		err := queryFromCtx(ctx, nil).Get(&dest, "SELECT pg_try_advisory_xact_lock($1) as locked", lockID)
		if err != nil {
			return err
		}
		if !dest.Locked {
			return db.ErrLockFailed
		}
		return fn(ctx)
	})
}
