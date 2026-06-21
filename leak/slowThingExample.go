package main

import (
	"context"
	"time"
)

func fetch(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ch := make(chan int, 1) // BUFFERED size 1: goroutine gửi xong là thoát, không bị leak

	go func() {
		ch <- slowThing() // slowThing() chạy ở đây, KHÔNG chặn select bên dưới
	}()

	select {
	case v := <-ch:
		return v, nil
	case <-ctx.Done():
		return 0, ctx.Err() // return sau 2s nếu slowThing() vẫn chưa xong
	}
}

// Consider this a 3rd party API
// They don't respect the context
func slowThing() int {
	time.Sleep(10 * time.Second)
	return 10
}
