package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type (
	txBeginner interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
)

func tx(ctx context.Context, conn txBeginner, fn func(ctx context.Context) error) error {

	if q := ctx.Value(pgCtx); q != nil {
		if _, ok := q.(pgx.Tx); ok {
			log.Fatal().Msg("database context is already in a transaction")
		}
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("cant obtain transaction on the database")
		return err
	}

	defer func() {
		if !tx.Conn().IsClosed() {
			if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				log.Error().Err(err).Msg("error while rolling back")
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
