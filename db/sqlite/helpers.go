package sqlite

import (
	"context"
	"database/sql"

	"github.com/fdelbos/commons/db"
)

type (
	ctxKey string
)

const (
	sqlCtx ctxKey = "internal/db/sqlite/ctx"
)

func queryFromCtx(ctx context.Context, defaultConn *sql.DB) db.Query {
	q := ctx.Value(sqlCtx)
	if q == nil {
		return &query{defaultConn, ctx}
	}
	if _, ok := q.(*sql.DB); ok {
		return &query{q.(*sql.DB), ctx}
	}
	panic("sqlite database context is not a Query object")
}
