package worker_pool

import (
	"context"
	"sync"
)

type Manager struct {
	N int
	ctx context.Context
	jobQueue chan *Job
	wg *sync.WaitGroup
}

func NewManager(n int, q chan *Job, wg *sync.WaitGroup) *Manager {
	return &Manager{
		N: n,
		jobQueue: q,
		wg: wg,
		ctx: context.Background(),
	}
}

func (m *Manager) GracefulShutdown(){
	m.wg.Wait()
}

func (m *Manager) HardShutdown(){
	m.cancel()
}

func (m *Manager) cancel(){
}