package main

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if res, err := GetPayment(ctx); err == nil {
		println(res)
	}
}
