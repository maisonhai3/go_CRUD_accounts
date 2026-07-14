package phase1semaphores

type JobQueue struct {
	q chan *Job
}