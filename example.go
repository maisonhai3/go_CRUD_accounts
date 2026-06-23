package main

import (
	"context"
	"time"
)

func slowThing(ctx context.Context) int {
	time.Sleep(10 * time.Second)
	return 10
}

type result struct {
	val int
	err error
}

func fetch(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	var
	defer cancel()



	ch := make(chan result, 1) // buffered=1 prevents sender goroutine from leaking on timeout
	go func() {
		ch <- result{val: slowThing(ctx), err: nil}
	}()

	select {
	case r := <-ch:
		return r.val, r.err
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	defer fmt.print()
}
