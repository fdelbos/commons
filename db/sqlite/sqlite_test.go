package sqlite_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fdelbos/commons/db"
	"github.com/fdelbos/commons/db/sqlite"
	"github.com/stretchr/testify/assert"
)

type (
	// TheTable is a table.
	TheTable struct {
		ID        int       `db:"id"`
		CreatedAt time.Time `db:"created_at"`
		Msg       string    `db:"msg"`
	}
)

func TestHSQLITE(t *testing.T) {
	dname, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(dname)

	dbURL := fmt.Sprintf("%s/db.sqlite3", dname)

	// create a new connection
	conn, err := sqlite.NewConn(dbURL)
	assert.NoError(t, err)
	defer conn.Close()

	// migrate the database
	err = sqlite.Migrate(dbURL, os.DirFS("test_migrations/."))
	assert.NoError(t, err)

	// insert a row
	ctx := context.Background()
	id := 0
	err = conn.Query(ctx).Get(&id, "insert into the_table (created_at, msg) values ($1, $2) returning id", "2021-01-01", "hello world")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	// select a row
	dest := TheTable{}
	err = conn.Query(ctx).Get(&dest, "select * from the_table where id = $1", id)
	assert.NoError(t, err)
	assert.Equal(t, id, dest.ID)
	assert.Equal(t, "2021-01-01", dest.CreatedAt.Format("2006-01-02"))
	assert.Equal(t, "hello world", dest.Msg)

	// select a row with a null value
	err = conn.Query(ctx).Get(&dest, "select * from the_table where id = $1", 999)
	assert.ErrorIs(t, err, db.ErrNoRows)

	// insert multiple rows
	for i := 0; i < 10; i++ {
		err = conn.Query(ctx).Exec("insert into the_table (created_at, msg) values ($1, $2)", "2021-01-01", "hello world")
		assert.NoError(t, err)
	}

	// select multiple rows
	dests := []TheTable{}
	err = conn.Query(ctx).Select(&dests, "select * from the_table")
	assert.NoError(t, err)
	assert.Equal(t, 11, len(dests))

}
