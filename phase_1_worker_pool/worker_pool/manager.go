package phase1semaphores

import (
	"context"
	"sync"
)

type Manager struct {
	N int
	ctx context.Context
	jobQueue JobQueue
	wg sync.WaitGroup
}

func (m *Manager) New (){
	ctx := context.Context(context.Background())
	m.ctx = ctx
}

func (m *Manager) GracefulShutdonw(){
	m.wg.Wait()
}

func (m *Manager) HardShutdown(){
	m.cancel()
}

func (m *Manager) cancel(){
}