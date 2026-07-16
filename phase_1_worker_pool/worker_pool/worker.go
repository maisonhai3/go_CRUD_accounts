package worker_pool

type WorkerStatus string
const (
	W WorkerStatus = "WORKING"
	P WorkerStatus = "PENDING"
)
type Worker struct {
	State WorkerStatus
	q chan *Job
	job *Job
}

func NewWorker(q chan *Job) *Worker{
	w := Worker{
		State: P,
		q: q,
	}
	return &w
}

func (w *Worker) Perform() {
}

func (w *Worker) Start() {
	for job := range w.q{
		job.Perform()
	}
}