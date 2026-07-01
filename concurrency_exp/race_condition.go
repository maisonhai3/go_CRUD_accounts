package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var balance int = 100

func withdraw(){
	mu.Lock()
	defer mu.Unlock()

	if balance >= 100{
		balance -= 100
		// fmt.Println("Withdraw: %i", 100)
		if balance < 0 {
			fmt.Println("Withdraw ERROR: %i", balance)
		}
	}
}

func MultipleWithdraw(){
	var wg sync.WaitGroup

	for range 10000{
		wg.Go(withdraw)
	}

	wg.Wait()
}