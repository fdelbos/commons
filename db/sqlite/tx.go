package sqlite

import (
	"context"
	"database/sql"
	"log"
)

type (
	txBeginner interface {
		Begin() (*sql.Tx, error)
	}
)

func tx(ctx context.Context, conn txBeginner, fn func(ctx context.Context) error) error {

	if q := ctx.Value(sqlCtx); q != nil {
		if _, ok := q.(*sql.Tx); ok {
			log.Fatal("database context is already in a transaction")
		}
	}

	tx, err := conn.Begin()
	if err != nil {
		log.Printf("cant obtain transaction on the database: %v", err)
		return err
	}

	closed := false
	defer func() {
		if !closed {
			if err := tx.Rollback(); err != nil {
				log.Printf("error while rolling back: %v", err)
			}
		}
	}()

	fnCtx := context.WithValue(ctx, sqlCtx, tx)
	if err := fn(fnCtx); err != nil {
		return err
	} else if err := tx.Commit(); err != nil {
		return err
	} else {
		closed = true
		return nil
	}
}
