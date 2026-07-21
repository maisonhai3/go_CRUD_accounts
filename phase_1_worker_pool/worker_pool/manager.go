package worker_pool

import (
	"context"
	"sync"
)

type Manager struct {
	N int
	ctx context.Context
	cancel context.CancelFunc
	
	jobQueue chan *Job
	wg *sync.WaitGroup
	workers [] *Worker
}

func NewManager(n int, q chan *Job, wg *sync.WaitGroup) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		N: n,
		jobQueue: q,
		wg: wg,
		ctx: ctx,
		cancel: cancel,
	}
}

func (m *Manager) InitWorker(){
	for range m.N {
		w := NewWorker(m.jobQueue, m.ctx)
		w.Start()
		m.workers = append(m.workers, w)
	}
}

func (m *Manager) GracefulShutdown(){
	m.wg.Wait()
}

func (m *Manager) HardShutdown(){
	m.cancel()
}