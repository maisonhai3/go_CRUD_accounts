package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

func CallBank(ctx context.Context, bankURL string) (string, error) {
	// Gắn context chứa tín hiệu cancel/timeout vào trực tiếp Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bankURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bank error: %d", resp.StatusCode)
	}

	return "Thành công từ: " + bankURL, nil
}

func GetPayment(ctx context.Context) (string, error) {
	const nAPI = 3
	var wg sync.WaitGroup
	wg.Add(nAPI)

	// API
	resCh := make(chan string, nAPI)
	apiList := []string{
		"https://httpbin.org/delay/5",
		"https://httpbin.org/status/500",
		"https://dummyjson.com/users/1",
	}
	for _, api := range apiList {
		go func() {
			defer wg.Done()
			if res, err := CallBank(ctx, api); err == nil {
				resCh <- res
			}
		}()
	}

	// Orchestrator
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()

	// Result
	select {
	case res := <-resCh: // Happy
		return res, nil
	case <-allDone: // Party happy
		return "", nil
	case <-ctx.Done(): // Timeout
		return "", ctx.Err()
	}
}
