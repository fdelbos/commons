package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/dchest/uniuri"
	"github.com/fdelbos/commons/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ctxKey string
)

const (
	pgCtx ctxKey = "internal/db/pg/ctx"
)

func queryFromCtx(ctx context.Context, defaultConn pgxInterface) db.Query {
	q := ctx.Value(pgCtx)
	if q == nil {
		return &query{defaultConn, ctx}
	}
	// return &query{q.(pgxInterface), ctx}

	switch v := q.(type) {
	case pgx.Tx:
		return &query{v, ctx}

	case *pgxpool.Pool:
		return &query{v, ctx}

	case *pgx.Conn:
		return &query{v, ctx}

	default:
		log.Fatalf("db/pg unknown type %T from context when looking for a query", v)
		return nil
	}
}

func CreateNewDB(ctx context.Context, posrgresURL, dbName string) error {
	pgConn, err := pgx.Connect(context.Background(), posrgresURL)
	if err != nil {
		log.Print("db/pg Unable to connect to postgres database: ", err)
		return err
	}

	defer pgConn.Close(ctx)

	// create the new database
	if _, err := pgConn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s ENCODING 'UTF8';", dbName)); err != nil {
		log.Print("db/pg Unable to create new database: ", err)
		return err
	}

	return nil
}

func DropDB(ctx context.Context, dbURL string) error {
	posrgresURL, err := db.ReplaceDBInURL(dbURL, "postgres")
	if err != nil {
		return err
	}

	dbName, err := db.DBName(dbURL)
	if err != nil {
		return err
	}

	pgConn, err := pgx.Connect(context.Background(), posrgresURL)
	if err != nil {
		log.Print("db/pg Unable to connect to postgres database: ", err)
		return err
	}

	defer pgConn.Close(ctx)

	// drop the database
	if _, err := pgConn.Exec(ctx, fmt.Sprintf("drop database if exists %s;", dbName)); err != nil {
		log.Print("db/pg Unable to drop the database: ", err)
		return err
	}
	return nil
}

func GenerateDB(ctx context.Context, dbURL, prefix string, migrate func(string) error) (string, error) {
	postgresURL, err := db.ReplaceDBInURL(dbURL, "postgres")
	if err != nil {
		return "", err
	}

	cloneName := fmt.Sprintf(
		"%s_%s",
		prefix,
		uniuri.NewLenChars(6, []byte("abcdefghijklmnopqrstuvwxyz")))

	CreateNewDB(ctx, postgresURL, cloneName)
	cloneURL, err := db.ReplaceDBInURL(dbURL, cloneName)
	if err != nil {
		return "", err
	}

	if err := migrate(cloneURL); err != nil {
		return "", err
	}
	return cloneURL, nil
}

func newPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Printf("db/pg Unable to parse DATABASE_URL: %s got error: %v", url, err)
		return nil, err
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Printf("db/pg Unable to create connection pool: %v", err)
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", config.ConnConfig.Host, config.ConnConfig.Port)
	log.Printf(`db/pg connected to host="%s" user="%s" db="%s"`,
		host,
		config.ConnConfig.User,
		config.ConnConfig.Database)

	if err := pool.Ping(context.Background()); err != nil {
		log.Printf("db/pg Unable to ping database: %v", err)
		return nil, err
	}

	return pool, nil
}

func newConn(ctx context.Context, url string) (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(url)
	if err != nil {
		log.Printf("db/pg Unable to parse DATABASE_URL: %s got error: %v", url, err)
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Print("db/pg Unable to create connection")
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Printf(`db/pg connected to host="%s" user="%s" db="%s"`,
		host,
		config.User,
		config.Database)

	if err := conn.Ping(context.Background()); err != nil {
		log.Printf("db/pg Unable to ping database: %v", err)
		return nil, err
	}

	return conn, nil
}
