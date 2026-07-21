package worker_pool

import "context"

type WorkerStatus string
const (
	W WorkerStatus = "WORKING"
	P WorkerStatus = "PENDING"
)
type Worker struct {
	State WorkerStatus
	q chan *Job
	job *Job
	ctx context.Context
}

func NewWorker(q chan *Job, mngCtx context.Context) *Worker{
	w := Worker{
		State: P,
		q: q,
	}
	return &w
}

func (w *Worker) Start() {
	for {
		select {
			case <- w.ctx.Done():
				return
			case job, ok := <- w.q:
				if !ok {
					return
				}
				job.Perform()
		}
	}
}

func (w *Worker) Perform(j *Job) {
	w.State = W
	j.Perform()
	w.State = P
}
