package worker_pool

type JobQueue struct {
	q chan *Job
}