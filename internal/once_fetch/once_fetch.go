package once_fetch

import (
	"context"
	"sync"
	"sync/atomic"
)

type Once[T any] struct {
	done atomic.Bool
	m    sync.Mutex
	got  T
	err  error
}

// Do is identical to sync.Once, except that it accepts a context
// and returns the result of calling f.
func (o *Once[T]) Do(ctx context.Context, f func(context.Context) (T, error)) (T, error) {
	if !o.done.Load() {
		o.doSlow(ctx, f)
	}

	return o.got, o.err
}

func (o *Once[T]) doSlow(ctx context.Context, f func(context.Context) (T, error)) {
	o.m.Lock()
	defer o.m.Unlock()
	if !o.done.Load() {
		defer o.done.Store(true)

		o.got, o.err = f(ctx)
	}
}
