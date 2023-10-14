package utils

import (
	"context"
	"sync"
	"time"
)

// Refresh returns a function that returns the result of fn, refreshed every interval.
func Refresh[T any](ctx context.Context, interval time.Duration, fn func(ctx context.Context) *T) func() *T {
	var (
		ticker = time.NewTicker(interval)
		mut    sync.RWMutex
		obj    T
	)

	update := func() {
		mut.Lock()
		defer mut.Unlock()
		res := fn(ctx)
		if res != nil {
			obj = *res
		}
	}

	go func() {
		defer ticker.Stop()
		update() // initial update
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				update()
			}
		}
	}()

	return func() *T {
		mut.RLock()
		defer mut.RUnlock()

		return &obj
	}
}

// Cron runs fn every interval until ctx is canceled.
func Cron(ctx context.Context, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fn()
		}
	}
}
