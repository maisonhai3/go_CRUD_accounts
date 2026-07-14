package main

import (
	"context"
	"phase_1_worker_pool/worker_pool"
	"sync"
	"time"
)


func main(){
	// -- Workers --
	const N int = 3
	jobQ := make(chan *worker_pool.Job, N)
	var wg sync.WaitGroup
	
	m := worker_pool.NewManager(
		N,
		jobQ,
		&wg,
	)

	// -- Producers --


	// -- Shutdown --
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	
	// Gracefuly
	shutdownDone := make(chan struct{})

	go func(){
		m.GracefulShutdown()
		close(shutdownDone)
	}()
	select {
		case <- shutdownDone:
			return
		case <- ctx.Done():
			m.HardShutdown()
			return
	}
}