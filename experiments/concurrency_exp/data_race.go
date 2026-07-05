package main

import (
	"fmt"
	"sync"
)


func DataRace(){
	var counter = 0
	var wg sync.WaitGroup

	for range 10000 {
		wg.Go(func() {
			mu.Lock()
			defer mu.Unlock()
			
			counter++
		})
	}

	wg.Wait()

	fmt.Print(counter)
	
}