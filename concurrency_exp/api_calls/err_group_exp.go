package api_calls

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

var apiList = []string{
	"https://httpbin.org/delay/5",
	"https://httpbin.org/status/500",
	"https://dummyjson.com/users/1",
}

func Request3rdAPI() (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// All-or-nothing
	erg, ctx := errgroup.WithContext(ctx)
	erg.SetLimit(2)

	respond := make([]string, 3)
	for i, api := range apiList {
		erg.Go(func() error {
			res, err := CallBank(ctx, api)
			if err != nil {
				return err
			}

			respond[i] = res // Avoid mutex
			return nil
		})
	}

	if err := erg.Wait(); err != nil {
		return nil, err
	}

	return "OK", nil
}
