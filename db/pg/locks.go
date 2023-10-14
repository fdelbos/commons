package pg

import "context"

type (
	Locker struct {
		dbURL string
	}
)

const (
	tryLock = "SELECT pg_try_advisory_lock($1)"
)

func NewLocker(dbURL string) *Locker {
	return &Locker{
		dbURL: dbURL,
	}
}

func (l *Locker) TryLock(ctx context.Context, key int, fn func(context.Context) error) (bool, error) {
	return AdvisoryLock(ctx, l.dbURL, key, fn)
}

func AdvisoryLock(ctx context.Context, url string, key int, fn func(context.Context) error) (bool, error) {
	conn, err := NewConn(url)
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)
	res := false
	if err = conn.Query(ctx).Get(&res, tryLock, key); err != nil {
		return false, err
	}
	if !res {
		return false, nil
	}
	return true, fn(ctx)
}
