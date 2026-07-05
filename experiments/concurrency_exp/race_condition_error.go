package main

import (
	"fmt"
	"runtime"
	"sync"
)

var mu2 sync.Mutex
var balance2 int = 100

func withdrawError() {
	mu2.Lock()
	b := balance2
	mu2.Unlock()

	runtime.Gosched() // widen the TOCTOU window so goroutines interleave

	if b >= 100 {
		mu2.Lock()
		balance2 -= 100
		mu2.Unlock()
		if balance2 < 0 {
			fmt.Println("Withdraw ERROR:", balance2)
		}
	}
}

func MultipleWithdrawError() {
	var wg sync.WaitGroup

	for range 100000 {
		wg.Go(withdrawError)
	}

	wg.Wait()
}
