package worker_pool

type JobStatus string
const (
	Pending JobStatus = "pending"
	Running JobStatus = "running"
	Suceed JobStatus = "suceed"
	Failed JobStatus = "failed"
)

type Job struct {
	ID     string
	Status JobStatus
}

func (j *Job) Perform() {}
