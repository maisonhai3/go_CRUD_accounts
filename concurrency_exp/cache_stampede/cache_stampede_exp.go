package cache_stampede

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

func TrySingleFlight() {
	var calls int64

	fn := func() (any, error) {
		atomic.AddInt64(&calls, 1)
		time.Sleep(200 * time.Millisecond)
		return "ok", nil
	}

	var g singleflight.Group
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Do("k", fn)
		}()
	}
	wg.Wait()

	fmt.Println("calls:", atomic.LoadInt64(&calls)) // 1  ← gộp

	time.Sleep(250 * time.Millisecond) // gap > thời gian chạy
	g.Do("k", fn)
	fmt.Println("calls:", atomic.LoadInt64(&calls)) // 2  ← KHÔNG phải cache
}
