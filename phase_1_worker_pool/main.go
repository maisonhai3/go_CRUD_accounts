package main

import (
	"context"
	"phase_1_worker_pool/worker_pool"
	"sync"
	"time"
)


func main(){
	// -- Managers --
	const N int = 3
	jobQ := make(chan *worker_pool.Job, N)
	var wg sync.WaitGroup
	
	m := worker_pool.NewManager(
		N,
		jobQ,
		&wg,
	)
	
	// -- Workers --
	m.InitWorker()

	// -- Producers & Jobs--
	p := worker_pool.NewProducer(jobQ)

	var pWg sync.WaitGroup
	for range 9 {
		pWg.Go(
			func(){
				j := worker_pool.NewJob()
				p.PushJob(j) 
			})
	}

	go func(){
		pWg.Wait()
		p.CloseJobQueue()
	}()
	
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