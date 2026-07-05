package cache_stampede

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)


func DirtyHandSingleFlight() {
	var balance atomic.Int64
	addBalance := func ()(any, error){
		balance.Add(100)
		time.Sleep(10 * time.Millisecond)
		return "ok", nil
	}

	// MAIN //
	var wg sync.WaitGroup
	var sfg singleflight.Group
	for range 10000 {
		wg.Go(func(){
			sfg.Do("add_100", addBalance)
		})
	}
	wg.Wait()
	fmt.Println(balance.Load())

	time.Sleep(1*time.Second)
	sfg.Do("add_100", addBalance)
	fmt.Println(balance.Load())
}
