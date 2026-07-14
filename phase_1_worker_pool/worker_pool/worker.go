package worker_pool

type WorkerStatus string
const (
	W WorkerStatus = "WORKING"
	P WorkerStatus = "PENDING"
)

type Worker struct {
	State string
	job *Job
}

func (w *Worker) Perform() {
}